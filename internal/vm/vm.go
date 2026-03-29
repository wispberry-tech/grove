// internal/vm/vm.go
package vm

import (
	"context"
	"fmt"
	"html"
	"strings"
	"sync"

	"grove/internal/compiler"
	"grove/internal/scope"
)

// VM is a stack-based bytecode executor. Instances are pooled; do not hold references.
type VM struct {
	stack [256]Value
	sp    int
	eng   EngineIface
	sc    *scope.Scope
	out   strings.Builder
}

var vmPool = sync.Pool{
	New: func() any {
		return &VM{}
	},
}

// Execute runs bc with data as the render context and returns the rendered string.
func Execute(ctx context.Context, bc *compiler.Bytecode, data map[string]any, eng EngineIface) (string, error) {
	v := vmPool.Get().(*VM)
	defer func() {
		v.out.Reset()
		v.sp = 0
		v.sc = nil
		v.eng = nil
		vmPool.Put(v)
	}()
	v.eng = eng

	// Build three-layer scope: local (empty) → render (data) → global
	globalSc := scope.New(nil)
	for k, val := range eng.GlobalData() {
		globalSc.Set(k, val)
	}
	renderSc := scope.New(globalSc)
	for k, val := range data {
		renderSc.Set(k, val)
	}
	v.sc = scope.New(renderSc) // local scope (for set, with, etc. — Plan 2)

	return v.run(ctx, bc)
}

func (v *VM) run(ctx context.Context, bc *compiler.Bytecode) (string, error) {
	ip := 0
	instrs := bc.Instrs
	for ip < len(instrs) {
		// Context cancellation check (also serves as yield point)
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		instr := instrs[ip]
		ip++

		switch instr.Op {
		case compiler.OP_HALT:
			return v.out.String(), nil

		case compiler.OP_PUSH_NIL:
			v.push(Nil)

		case compiler.OP_PUSH_CONST:
			v.push(fromConst(bc.Consts[instr.A]))

		case compiler.OP_LOAD:
			name := bc.Names[instr.A]
			val, found := v.sc.Get(name)
			if !found {
				if v.eng.StrictVariables() {
					return "", &runtimeErr{msg: fmt.Sprintf("undefined variable %q", name)}
				}
				v.push(Nil)
			} else {
				v.push(FromAny(val))
			}

		case compiler.OP_GET_ATTR:
			obj := v.pop()
			name := bc.Names[instr.A]
			result, err := GetAttr(obj, name, v.eng.StrictVariables())
			if err != nil {
				return "", &runtimeErr{msg: err.Error()}
			}
			v.push(result)

		case compiler.OP_GET_INDEX:
			key := v.pop()
			obj := v.pop()
			result, err := GetIndex(obj, key)
			if err != nil {
				return "", &runtimeErr{msg: err.Error()}
			}
			v.push(result)

		case compiler.OP_OUTPUT:
			val := v.pop()
			if val.typ == TypeSafeHTML {
				v.out.WriteString(val.sval)
			} else if val.typ != TypeNil {
				v.out.WriteString(html.EscapeString(val.String()))
			}
			// nil outputs nothing

		case compiler.OP_OUTPUT_RAW:
			val := v.pop()
			v.out.WriteString(val.String())

		case compiler.OP_ADD:
			b, a := v.pop(), v.pop()
			v.push(arithAdd(a, b))

		case compiler.OP_SUB:
			b, a := v.pop(), v.pop()
			v.push(arithSub(a, b))

		case compiler.OP_MUL:
			b, a := v.pop(), v.pop()
			v.push(arithMul(a, b))

		case compiler.OP_DIV:
			b, a := v.pop(), v.pop()
			result, err := arithDiv(a, b)
			if err != nil {
				return "", err
			}
			v.push(result)

		case compiler.OP_MOD:
			b, a := v.pop(), v.pop()
			result, err := arithMod(a, b)
			if err != nil {
				return "", err
			}
			v.push(result)

		case compiler.OP_CONCAT:
			b, a := v.pop(), v.pop()
			v.push(StringVal(a.String() + b.String()))

		case compiler.OP_EQ:
			b, a := v.pop(), v.pop()
			v.push(BoolVal(valEqual(a, b)))

		case compiler.OP_NEQ:
			b, a := v.pop(), v.pop()
			v.push(BoolVal(!valEqual(a, b)))

		case compiler.OP_LT:
			b, a := v.pop(), v.pop()
			r, err := valCompare(a, b)
			if err != nil {
				return "", err
			}
			v.push(BoolVal(r < 0))

		case compiler.OP_LTE:
			b, a := v.pop(), v.pop()
			r, err := valCompare(a, b)
			if err != nil {
				return "", err
			}
			v.push(BoolVal(r <= 0))

		case compiler.OP_GT:
			b, a := v.pop(), v.pop()
			r, err := valCompare(a, b)
			if err != nil {
				return "", err
			}
			v.push(BoolVal(r > 0))

		case compiler.OP_GTE:
			b, a := v.pop(), v.pop()
			r, err := valCompare(a, b)
			if err != nil {
				return "", err
			}
			v.push(BoolVal(r >= 0))

		case compiler.OP_AND:
			b, a := v.pop(), v.pop()
			v.push(BoolVal(Truthy(a) && Truthy(b)))

		case compiler.OP_OR:
			b, a := v.pop(), v.pop()
			v.push(BoolVal(Truthy(a) || Truthy(b)))

		case compiler.OP_NOT:
			a := v.pop()
			v.push(BoolVal(!Truthy(a)))

		case compiler.OP_NEGATE:
			a := v.pop()
			switch a.typ {
			case TypeInt:
				v.push(IntVal(-a.ival))
			case TypeFloat:
				v.push(FloatVal(-a.fval))
			default:
				v.push(IntVal(0))
			}

		case compiler.OP_JUMP:
			ip = int(instr.A)

		case compiler.OP_JUMP_FALSE:
			cond := v.pop()
			if !Truthy(cond) {
				ip = int(instr.A)
			}

		case compiler.OP_FILTER:
			name := bc.Names[instr.A]
			argc := int(instr.B)
			args := make([]Value, argc)
			for i := argc - 1; i >= 0; i-- {
				args[i] = v.pop()
			}
			val := v.pop()
			fn, ok := v.eng.LookupFilter(name)
			if !ok {
				return "", &runtimeErr{msg: fmt.Sprintf("unknown filter %q", name)}
			}
			result, err := fn(val, args)
			if err != nil {
				return "", &runtimeErr{msg: err.Error()}
			}
			v.push(result)

		default:
			return "", fmt.Errorf("vm: unknown opcode %d at ip=%d", instr.Op, ip-1)
		}
	}
	return v.out.String(), nil
}

// ─── Stack helpers ────────────────────────────────────────────────────────────

func (v *VM) push(val Value) {
	if v.sp >= len(v.stack) {
		panic("vm: stack overflow")
	}
	v.stack[v.sp] = val
	v.sp++
}

func (v *VM) pop() Value {
	v.sp--
	return v.stack[v.sp]
}

// ─── Arithmetic ───────────────────────────────────────────────────────────────

func fromConst(c any) Value {
	switch x := c.(type) {
	case bool:
		return BoolVal(x)
	case int64:
		return IntVal(x)
	case float64:
		return FloatVal(x)
	case string:
		return StringVal(x)
	}
	return Nil
}

func arithAdd(a, b Value) Value {
	if a.typ == TypeFloat || b.typ == TypeFloat {
		af, _ := a.ToFloat64()
		bf, _ := b.ToFloat64()
		return FloatVal(af + bf)
	}
	ai, aok := a.ToInt64()
	bi, bok := b.ToInt64()
	if aok && bok {
		return IntVal(ai + bi)
	}
	return StringVal(a.String() + b.String())
}

func arithSub(a, b Value) Value {
	if a.typ == TypeFloat || b.typ == TypeFloat {
		af, _ := a.ToFloat64()
		bf, _ := b.ToFloat64()
		return FloatVal(af - bf)
	}
	ai, _ := a.ToInt64()
	bi, _ := b.ToInt64()
	return IntVal(ai - bi)
}

func arithMul(a, b Value) Value {
	if a.typ == TypeFloat || b.typ == TypeFloat {
		af, _ := a.ToFloat64()
		bf, _ := b.ToFloat64()
		return FloatVal(af * bf)
	}
	ai, _ := a.ToInt64()
	bi, _ := b.ToInt64()
	return IntVal(ai * bi)
}

func arithDiv(a, b Value) (Value, error) {
	af, _ := a.ToFloat64()
	bf, _ := b.ToFloat64()
	if bf == 0 {
		return Nil, &runtimeErr{msg: "division by zero"}
	}
	result := af / bf
	// Return int if both operands were ints and result is whole
	if a.typ == TypeInt && b.typ == TypeInt && result == float64(int64(result)) {
		return IntVal(int64(result)), nil
	}
	return FloatVal(result), nil
}

func arithMod(a, b Value) (Value, error) {
	bi, bok := b.ToInt64()
	if !bok || bi == 0 {
		bf, _ := b.ToFloat64()
		if bf == 0 {
			return Nil, &runtimeErr{msg: "modulo by zero"}
		}
	}
	ai, _ := a.ToInt64()
	return IntVal(ai % bi), nil
}

// ─── Comparison ───────────────────────────────────────────────────────────────

func valEqual(a, b Value) bool {
	if a.typ != b.typ {
		// Cross-type numeric equality
		if (a.typ == TypeInt || a.typ == TypeFloat) && (b.typ == TypeInt || b.typ == TypeFloat) {
			af, _ := a.ToFloat64()
			bf, _ := b.ToFloat64()
			return af == bf
		}
		return false
	}
	switch a.typ {
	case TypeNil:
		return true
	case TypeBool:
		return a.ival == b.ival
	case TypeInt:
		return a.ival == b.ival
	case TypeFloat:
		return a.fval == b.fval
	case TypeString, TypeSafeHTML:
		return a.sval == b.sval
	}
	return false
}

// valCompare returns -1, 0, or 1 for a <=> b.
func valCompare(a, b Value) (int, error) {
	if (a.typ == TypeInt || a.typ == TypeFloat) && (b.typ == TypeInt || b.typ == TypeFloat) {
		af, _ := a.ToFloat64()
		bf, _ := b.ToFloat64()
		if af < bf {
			return -1, nil
		} else if af > bf {
			return 1, nil
		}
		return 0, nil
	}
	if a.typ == TypeString && b.typ == TypeString {
		if a.sval < b.sval {
			return -1, nil
		} else if a.sval > b.sval {
			return 1, nil
		}
		return 0, nil
	}
	return 0, &runtimeErr{msg: fmt.Sprintf("cannot compare %v and %v", a.typ, b.typ)}
}

// ─── Runtime error ────────────────────────────────────────────────────────────

type runtimeErr struct {
	msg string
}

func (e *runtimeErr) Error() string { return e.msg }
