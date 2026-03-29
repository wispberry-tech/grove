// internal/ast/node.go
package ast

// Node is the base interface for all AST nodes.
type Node interface{ groveNode() }

// Program is the root node.
type Program struct{ Body []Node }

func (*Program) groveNode() {}

// ─── Statement nodes ──────────────────────────────────────────────────────────

// TextNode holds raw text content (no interpolation).
type TextNode struct {
	Value string
	Line  int
}

func (*TextNode) groveNode() {}

// OutputNode holds an {{ expression }} to be evaluated and printed.
type OutputNode struct {
	Expr       Node
	StripLeft  bool
	StripRight bool
	Line       int
}

func (*OutputNode) groveNode() {}

// RawNode holds content from {% raw %}...{% endraw %} — printed verbatim.
type RawNode struct {
	Value string
	Line  int
}

func (*RawNode) groveNode() {}

// TagNode is an unrecognised or deferred tag (e.g. if/for/extends).
// The parser uses this as a placeholder for tags handled in later plans,
// and to reject banned tags (extends/import) in inline mode.
type TagNode struct {
	Name string
	Line int
}

func (*TagNode) groveNode() {}

// ─── Expression nodes ─────────────────────────────────────────────────────────

// NilLiteral is the nil/null literal.
type NilLiteral struct{ Line int }

func (*NilLiteral) groveNode() {}

// BoolLiteral is true or false.
type BoolLiteral struct {
	Value bool
	Line  int
}

func (*BoolLiteral) groveNode() {}

// IntLiteral is an integer literal.
type IntLiteral struct {
	Value int64
	Line  int
}

func (*IntLiteral) groveNode() {}

// FloatLiteral is a floating-point literal.
type FloatLiteral struct {
	Value float64
	Line  int
}

func (*FloatLiteral) groveNode() {}

// StringLiteral is a quoted string literal.
type StringLiteral struct {
	Value string
	Line  int
}

func (*StringLiteral) groveNode() {}

// Identifier is a variable reference.
type Identifier struct {
	Name string
	Line int
}

func (*Identifier) groveNode() {}

// AttributeAccess is obj.key — resolves key on obj.
type AttributeAccess struct {
	Object Node
	Key    string
	Line   int
}

func (*AttributeAccess) groveNode() {}

// IndexAccess is obj[key] — integer or string key.
type IndexAccess struct {
	Object Node
	Key    Node
	Line   int
}

func (*IndexAccess) groveNode() {}

// BinaryExpr is left op right.
// Op is one of: + - * / % ~ == != < <= > >= and or
type BinaryExpr struct {
	Op    string
	Left  Node
	Right Node
	Line  int
}

func (*BinaryExpr) groveNode() {}

// UnaryExpr is op operand.
// Op is one of: not -
type UnaryExpr struct {
	Op      string
	Operand Node
	Line    int
}

func (*UnaryExpr) groveNode() {}

// TernaryExpr is: Consequence if Condition else Alternative
// (Grove syntax: `value if cond else fallback`)
type TernaryExpr struct {
	Condition   Node
	Consequence Node
	Alternative Node
	Line        int
}

func (*TernaryExpr) groveNode() {}

// FilterExpr applies Filter(Args...) to Value.
// e.g. name | truncate(20, "…") → FilterExpr{Value: Identifier{name}, Filter: "truncate", Args: [20, "…"]}
type FilterExpr struct {
	Value  Node
	Filter string
	Args   []Node
	Line   int
}

func (*FilterExpr) groveNode() {}
