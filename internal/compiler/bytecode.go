// internal/compiler/bytecode.go
package compiler

// Opcode is a single VM instruction opcode.
type Opcode uint8

// Instruction is a fixed-width 8-byte VM instruction.
// Field layout: A(2) + B(2) + Op(1) + Flags(1) + _(2) = 8 bytes.
type Instruction struct {
	A     uint16  // primary operand (const index, name index, jump target, arg count)
	B     uint16  // secondary operand (argc for FILTER)
	Op    Opcode
	Flags uint8   // modifier bits (e.g. escape flag)
	_     [2]byte // reserved
}

const (
	OP_HALT       Opcode = iota
	OP_PUSH_CONST        // A = index into Consts
	OP_PUSH_NIL
	OP_LOAD              // A = index into Names — scope lookup
	OP_GET_ATTR          // A = index into Names — pop obj, push obj.Names[A]
	OP_GET_INDEX         // pop key, pop obj, push obj[key]
	OP_OUTPUT            // pop value, HTML-escape and write (unless SafeHTML)
	OP_OUTPUT_RAW        // pop value, write verbatim (no escaping)
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD
	OP_CONCAT   // ~ operator: pop b, pop a, push a+b as string
	OP_EQ
	OP_NEQ
	OP_LT
	OP_LTE
	OP_GT
	OP_GTE
	OP_AND
	OP_OR
	OP_NOT
	OP_NEGATE           // unary minus
	OP_JUMP             // A = target instruction index (unconditional)
	OP_JUMP_FALSE       // A = target; pop value, jump if falsy
	OP_FILTER           // A = name index, B = argc; pop argc args then value, push result
)

// Bytecode is the compiled output for a single template.
// It is immutable after compilation and safe for concurrent use.
type Bytecode struct {
	Instrs []Instruction
	Consts []any    // constant pool: string | int64 | float64 | bool
	Names  []string // name pool: variable names, attribute names, filter names
}
