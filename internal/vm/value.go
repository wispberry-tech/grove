// internal/vm/value.go
package vm

import (
	"fmt"
	"reflect"
	"strconv"
)

// ValueType identifies the runtime type of a Value.
type ValueType uint8

const (
	TypeNil        ValueType = iota
	TypeBool                 // ival: 0=false, 1=true
	TypeInt                  // ival: int64
	TypeFloat                // fval: float64
	TypeString               // sval: string
	TypeSafeHTML             // sval: trusted HTML, bypass auto-escape
	TypeList                 // oval: []Value
	TypeMap                  // oval: map[string]any (Go map, accessed via key lookup)
	TypeResolvable           // oval: Resolvable
)

// Value is the runtime value type. Zero value is Nil.
type Value struct {
	typ  ValueType
	ival int64
	fval float64
	sval string
	oval any
}

// Nil is the zero Value.
var Nil = Value{}

// Resolvable is implemented by Go types that expose specific fields to templates.
type Resolvable interface {
	GroveResolve(key string) (any, bool)
}

// ─── Constructors ─────────────────────────────────────────────────────────────

func BoolVal(b bool) Value {
	v := Value{typ: TypeBool}
	if b {
		v.ival = 1
	}
	return v
}

func IntVal(n int64) Value     { return Value{typ: TypeInt, ival: n} }
func FloatVal(f float64) Value { return Value{typ: TypeFloat, fval: f} }
func StringVal(s string) Value { return Value{typ: TypeString, sval: s} }
func SafeHTMLVal(s string) Value { return Value{typ: TypeSafeHTML, sval: s} }
func ListVal(items []Value) Value { return Value{typ: TypeList, oval: items} }
func MapVal(m map[string]any) Value { return Value{typ: TypeMap, oval: m} }
func ResolvableVal(r Resolvable) Value { return Value{typ: TypeResolvable, oval: r} }

// ─── String representation ────────────────────────────────────────────────────

// String returns the string representation for template output.
func (v Value) String() string {
	switch v.typ {
	case TypeNil:
		return ""
	case TypeBool:
		if v.ival != 0 {
			return "true"
		}
		return "false"
	case TypeInt:
		return strconv.FormatInt(v.ival, 10)
	case TypeFloat:
		// Format without trailing zeros; use shortest representation
		s := strconv.FormatFloat(v.fval, 'f', -1, 64)
		return s
	case TypeString, TypeSafeHTML:
		return v.sval
	case TypeList:
		return fmt.Sprintf("%v", v.oval)
	case TypeMap:
		return fmt.Sprintf("%v", v.oval)
	}
	return ""
}

// IsSafeHTML reports whether this value carries trusted HTML.
func (v Value) IsSafeHTML() bool { return v.typ == TypeSafeHTML }

// IsNil reports whether this is the nil value.
func (v Value) IsNil() bool { return v.typ == TypeNil }

// ─── Type coercions ───────────────────────────────────────────────────────────

// Truthy follows Jinja2/Python-style truthiness:
// nil=false, bool=value, int=nonzero, float=nonzero, string=nonempty, list=nonempty
func Truthy(v Value) bool {
	switch v.typ {
	case TypeNil:
		return false
	case TypeBool:
		return v.ival != 0
	case TypeInt:
		return v.ival != 0
	case TypeFloat:
		return v.fval != 0
	case TypeString, TypeSafeHTML:
		return v.sval != ""
	case TypeList:
		if lst, ok := v.oval.([]Value); ok {
			return len(lst) > 0
		}
		return false
	case TypeMap:
		if m, ok := v.oval.(map[string]any); ok {
			return len(m) > 0
		}
		return false
	case TypeResolvable:
		return v.oval != nil
	}
	return false
}

// ToInt64 converts v to int64. Returns (0, false) if not convertible.
func (v Value) ToInt64() (int64, bool) {
	switch v.typ {
	case TypeInt:
		return v.ival, true
	case TypeFloat:
		return int64(v.fval), true
	case TypeBool:
		return v.ival, true
	case TypeString:
		n, err := strconv.ParseInt(v.sval, 10, 64)
		return n, err == nil
	}
	return 0, false
}

// ToFloat64 converts v to float64.
func (v Value) ToFloat64() (float64, bool) {
	switch v.typ {
	case TypeFloat:
		return v.fval, true
	case TypeInt:
		return float64(v.ival), true
	case TypeString:
		f, err := strconv.ParseFloat(v.sval, 64)
		return f, err == nil
	}
	return 0, false
}

// ─── Arithmetic helpers ───────────────────────────────────────────────────────

// FromAny wraps a Go value into a VM Value.
func FromAny(v any) Value {
	if v == nil {
		return Nil
	}
	switch x := v.(type) {
	case bool:
		return BoolVal(x)
	case int:
		return IntVal(int64(x))
	case int8:
		return IntVal(int64(x))
	case int16:
		return IntVal(int64(x))
	case int32:
		return IntVal(int64(x))
	case int64:
		return IntVal(x)
	case uint:
		return IntVal(int64(x))
	case uint64:
		return IntVal(int64(x))
	case float32:
		return FloatVal(float64(x))
	case float64:
		return FloatVal(x)
	case string:
		return StringVal(x)
	case Value:
		return x
	case Resolvable:
		return ResolvableVal(x)
	case []any:
		vals := make([]Value, len(x))
		for i, elem := range x {
			vals[i] = FromAny(elem)
		}
		return ListVal(vals)
	case []string:
		vals := make([]Value, len(x))
		for i, s := range x {
			vals[i] = StringVal(s)
		}
		return ListVal(vals)
	case []int:
		vals := make([]Value, len(x))
		for i, n := range x {
			vals[i] = IntVal(int64(n))
		}
		return ListVal(vals)
	case map[string]any:
		return MapVal(x)
	default:
		// Try Resolvable via interface assertion
		if r, ok := v.(Resolvable); ok {
			return ResolvableVal(r)
		}
		// Handle named map types (e.g. grove.Data which is map[string]any)
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Map && rv.Type().Key().Kind() == reflect.String {
			m := make(map[string]any, rv.Len())
			for _, k := range rv.MapKeys() {
				m[k.String()] = rv.MapIndex(k).Interface()
			}
			return MapVal(m)
		}
		return StringVal(fmt.Sprintf("%v", v))
	}
}

// GetAttr resolves obj.name. Returns (Nil, error) if not found.
func GetAttr(obj Value, name string, strict bool) (Value, error) {
	switch obj.typ {
	case TypeMap:
		m, _ := obj.oval.(map[string]any)
		if v, ok := m[name]; ok {
			return FromAny(v), nil
		}
		if strict {
			return Nil, fmt.Errorf("undefined attribute %q", name)
		}
		return Nil, nil
	case TypeResolvable:
		r, _ := obj.oval.(Resolvable)
		if v, ok := r.GroveResolve(name); ok {
			return FromAny(v), nil
		}
		if strict {
			return Nil, fmt.Errorf("undefined attribute %q", name)
		}
		return Nil, nil
	case TypeNil:
		if strict {
			return Nil, fmt.Errorf("cannot access .%s on nil", name)
		}
		return Nil, nil
	}
	if strict {
		return Nil, fmt.Errorf("cannot access .%s on %T", name, obj.oval)
	}
	return Nil, nil
}

// GetIndex resolves obj[key].
func GetIndex(obj, key Value) (Value, error) {
	switch obj.typ {
	case TypeList:
		lst, _ := obj.oval.([]Value)
		idx, ok := key.ToInt64()
		if !ok {
			return Nil, fmt.Errorf("list index must be integer, got %s", key.String())
		}
		if idx < 0 || idx >= int64(len(lst)) {
			return Nil, nil
		}
		return lst[idx], nil
	case TypeMap:
		m, _ := obj.oval.(map[string]any)
		k := key.String()
		if v, ok := m[k]; ok {
			return FromAny(v), nil
		}
		return Nil, nil
	}
	return Nil, fmt.Errorf("cannot index %T", obj.oval)
}

// ─── Filter support ───────────────────────────────────────────────────────────

// FilterFn is the function signature for filter implementations.
type FilterFn func(v Value, args []Value) (Value, error)

// FilterDef bundles a FilterFn with metadata.
type FilterDef struct {
	Fn          FilterFn
	OutputsHTML bool
}

// FilterOption modifies a FilterDef.
type FilterOption func(*FilterDef)

// NewFilterDef creates a FilterDef from fn with optional options.
func NewFilterDef(fn FilterFn, opts ...FilterOption) *FilterDef {
	d := &FilterDef{Fn: fn}
	for _, o := range opts {
		o(d)
	}
	return d
}

// OptionOutputsHTML marks a filter as returning SafeHTML (skips auto-escape).
func OptionOutputsHTML() FilterOption {
	return func(d *FilterDef) { d.OutputsHTML = true }
}

// FilterSet is a named collection of filters for bulk registration.
type FilterSet map[string]any

// EngineIface is the callback interface the VM uses to call back into the Engine.
type EngineIface interface {
	LookupFilter(name string) (FilterFn, bool)
	StrictVariables() bool
	GlobalData() map[string]any
}

// ArgInt reads args[i] as an integer, returning def if out of range or not convertible.
func ArgInt(args []Value, i, def int) int {
	if i >= len(args) {
		return def
	}
	if n, ok := args[i].ToInt64(); ok {
		return int(n)
	}
	return def
}
