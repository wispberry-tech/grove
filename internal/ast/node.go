package ast

import (
	"strings"

	"template-wisp/internal/lexer"
)

// Node represents a node in the abstract syntax tree.
type Node interface {
	// TokenLiteral returns the literal value of the token associated with this node.
	TokenLiteral() string
	// String returns a string representation of the node.
	String() string
}

// Statement represents a statement node.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out string
	for _, s := range p.Statements {
		out += s.String()
	}
	return out
}

func (p *Program) statementNode() {}

// LetStatement represents a let statement: `{% let .name = "value" %}` or `{% assign .name = "value" %}`
type LetStatement struct {
	Token lexer.Token // the LET or ASSIGN token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out string

	out += ls.TokenLiteral() + " "
	out += ls.Name.String()
	out += " = "

	if ls.Value != nil {
		out += ls.Value.String()
	}

	out += ";"

	return out
}

// Identifier represents an identifier node: `.name` or `.user.name`
type Identifier struct {
	Token lexer.Token // the DOT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return "." + i.Value }

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// Boolean represents a boolean literal.
type Boolean struct {
	Token lexer.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// StringLiteral represents a string literal.
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// ArrayLiteral represents an array literal.
type ArrayLiteral struct {
	Token    lexer.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out string

	out += "["
	for i, el := range al.Elements {
		out += el.String()
		if i < len(al.Elements)-1 {
			out += ", "
		}
	}
	out += "]"

	return out
}

// HashLiteral represents a hash/map literal.
type HashLiteral struct {
	Token lexer.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out string

	out += "{"
	for k, v := range hl.Pairs {
		out += k.String()
		out += ": "
		out += v.String()
		if k != nil && v != nil && len(hl.Pairs) > 1 {
			// We can't easily get the index, so we'll just not add comma after last item
			// This is a simplified approach
			out += ", "
		}
	}
	out += "}"

	return out
}

// PrefixExpression represents a prefix expression.
type PrefixExpression struct {
	Token    lexer.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out string

	out += "("
	out += pe.Operator
	out += pe.Right.String()
	out += ")"

	return out
}

// InfixExpression represents an infix expression.
type InfixExpression struct {
	Token    lexer.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out string

	out += "("
	out += ie.Left.String()
	out += ie.Operator
	out += ie.Right.String()
	out += ")"

	return out
}

// IfStatement represents an if statement.
type IfStatement struct {
	Token       lexer.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out string

	out += "if "
	out += is.Condition.String()
	out += " "
	out += is.Consequence.String()

	if is.Alternative != nil {
		out += "else "
		out += is.Alternative.String()
	}

	return out
}

// BlockStatement represents a block of statements.
type BlockStatement struct {
	Token      lexer.Token // The { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) expressionNode()      {} // To allow BlockStatement to be used as an expression
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out string

	for _, s := range bs.Statements {
		out += s.String()
	}

	return out
}

// ExpressionStatement represents an expression statement: `{% .name %}` or `{% . | date %}`
type ExpressionStatement struct {
	Token      lexer.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

// ReturnStatement represents a return statement: `return <expression>;`
type ReturnStatement struct {
	Token       lexer.Token // The 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out string

	out += rs.TokenLiteral() + " "

	if rs.ReturnValue != nil {
		out += rs.ReturnValue.String()
	}

	out += ";"

	return out
}

// FunctionLiteral represents a function literal.
type FunctionLiteral struct {
	Token      lexer.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out string

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out += fl.TokenLiteral()
	out += "("
	out += strings.Join(params, ", ")
	out += ") "
	out += fl.Body.String()

	return out
}

// CallExpression represents a function call.
type CallExpression struct {
	Token     lexer.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out string

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out += ce.Function.String()
	out += "("
	out += strings.Join(args, ", ")
	out += ")"

	return out
}

// IndexExpression represents an index expression.
type IndexExpression struct {
	Token lexer.Token // The '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out string

	out += "("
	out += ie.Left.String()
	out += "["
	out += ie.Index.String()
	out += "])"

	return out
}

// IfExpression represents an if expression (ternary-like if-else as expression).
type IfExpression struct {
	Token       lexer.Token // The 'if' token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out string

	out += "if "
	out += ie.Condition.String()
	out += " "
	out += ie.Consequence.String()
	out += "else "
	out += ie.Alternative.String()

	return out
}

// DotExpression represents a dot expression for variable access: .name, .user.name
type DotExpression struct {
	Token lexer.Token   // The DOT token
	Field *Identifier   // The first field after the dot
	Chain []*Identifier // Additional chained fields
}

func (de *DotExpression) expressionNode()      {}
func (de *DotExpression) TokenLiteral() string { return de.Token.Literal }
func (de *DotExpression) String() string {
	var out string
	out += "." + de.Field.String()
	for _, field := range de.Chain {
		out += "." + field.String()
	}
	return out
}

// PipeExpression represents a pipe expression for function calls: . | date, . | format "%s"
type PipeExpression struct {
	Token     lexer.Token  // The PIPE token
	Function  *Identifier  // The function name
	Arguments []Expression // Optional arguments
}

func (pe *PipeExpression) expressionNode()      {}
func (pe *PipeExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PipeExpression) String() string {
	var out string
	out += ". | " + pe.Function.String()
	for _, arg := range pe.Arguments {
		out += " " + arg.String()
	}
	return out
}

// AssignStatement represents an assignment statement: {% assign .name = value %}
type AssignStatement struct {
	Token lexer.Token // The ASSIGN token
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string {
	var out string
	out += "assign ." + as.Name.String() + " = "
	if as.Value != nil {
		out += as.Value.String()
	}
	return out
}

// UnlessStatement represents an unless statement: {% unless .condition %}
type UnlessStatement struct {
	Token       lexer.Token // The UNLESS token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (us *UnlessStatement) statementNode()       {}
func (us *UnlessStatement) TokenLiteral() string { return us.Token.Literal }
func (us *UnlessStatement) String() string {
	var out string
	out += "unless " + us.Condition.String() + " "
	if us.Consequence != nil {
		out += us.Consequence.String()
	}
	if us.Alternative != nil {
		out += " else " + us.Alternative.String()
	}
	return out
}

// ForStatement represents a for loop: {% for .item in .items %}
type ForStatement struct {
	Token      lexer.Token // The FOR token
	IndexVar   *Identifier // Optional index variable
	LoopVar    *Identifier // Loop variable
	Collection Expression  // The collection to iterate over
	Body       *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out string
	out += "for "
	if fs.IndexVar != nil {
		out += fs.IndexVar.String() + ", "
	}
	out += fs.LoopVar.String() + " in " + fs.Collection.String() + " "
	if fs.Body != nil {
		out += fs.Body.String()
	}
	return out
}

// WhileStatement represents a while loop: {% while .condition %}
type WhileStatement struct {
	Token     lexer.Token // The WHILE token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out string
	out += "while " + ws.Condition.String() + " "
	if ws.Body != nil {
		out += ws.Body.String()
	}
	return out
}

// RangeStatement represents a range loop: {% range .start .end %}
type RangeStatement struct {
	Token lexer.Token // The RANGE token
	Start Expression
	End   Expression
	Body  *BlockStatement
}

func (rs *RangeStatement) statementNode()       {}
func (rs *RangeStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *RangeStatement) String() string {
	var out string
	out += "range " + rs.Start.String() + " " + rs.End.String() + " "
	if rs.Body != nil {
		out += rs.Body.String()
	}
	return out
}

// CaseStatement represents a case statement: {% case .value %}
type CaseStatement struct {
	Token lexer.Token // The CASE token
	Value Expression
	Body  *BlockStatement
}

func (cs *CaseStatement) statementNode()       {}
func (cs *CaseStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *CaseStatement) String() string {
	var out string
	out += "case " + cs.Value.String() + " "
	if cs.Body != nil {
		out += cs.Body.String()
	}
	return out
}

// WithStatement represents a with block: {% with .user as .currentUser %}
type WithStatement struct {
	Token  lexer.Token // The WITH token
	Source Expression
	Target *Identifier
	Body   *BlockStatement
}

func (ws *WithStatement) statementNode()       {}
func (ws *WithStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WithStatement) String() string {
	var out string
	out += "with " + ws.Source.String() + " as " + ws.Target.String() + " "
	if ws.Body != nil {
		out += ws.Body.String()
	}
	return out
}

// CycleStatement represents a cycle tag: {% cycle 'odd' 'even' %}
type CycleStatement struct {
	Token  lexer.Token // The CYCLE token
	Values []Expression
}

func (cs *CycleStatement) statementNode()       {}
func (cs *CycleStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *CycleStatement) String() string {
	var out string
	out += "cycle "
	for i, v := range cs.Values {
		out += v.String()
		if i < len(cs.Values)-1 {
			out += " "
		}
	}
	return out
}

// IncrementStatement represents an increment tag: {% increment .counter %}
type IncrementStatement struct {
	Token    lexer.Token // The INCREMENT token
	Variable *Identifier
}

func (is *IncrementStatement) statementNode()       {}
func (is *IncrementStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IncrementStatement) String() string {
	return "increment ." + is.Variable.String()
}

// DecrementStatement represents a decrement tag: {% decrement .counter %}
type DecrementStatement struct {
	Token    lexer.Token // The DECREMENT token
	Variable *Identifier
}

func (ds *DecrementStatement) statementNode()       {}
func (ds *DecrementStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DecrementStatement) String() string {
	return "decrement ." + ds.Variable.String()
}

// BreakStatement represents a break tag: {% break %}
type BreakStatement struct {
	Token lexer.Token // The BREAK token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return "break" }

// ContinueStatement represents a continue tag: {% continue %}
type ContinueStatement struct {
	Token lexer.Token // The CONTINUE token
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string       { return "continue" }

// IncludeStatement represents an include tag: {% include "template" %}
type IncludeStatement struct {
	Token    lexer.Token // The INCLUDE token
	Template *StringLiteral
	Context  Expression
}

func (is *IncludeStatement) statementNode()       {}
func (is *IncludeStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IncludeStatement) String() string {
	var out string
	out += "include " + is.Template.String()
	if is.Context != nil {
		out += " " + is.Context.String()
	}
	return out
}

// RenderStatement represents a render tag: {% render "template" .data %}
type RenderStatement struct {
	Token    lexer.Token // The RENDER token
	Template *StringLiteral
	Params   []Expression
}

func (rs *RenderStatement) statementNode()       {}
func (rs *RenderStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *RenderStatement) String() string {
	var out string
	out += "render " + rs.Template.String()
	for _, p := range rs.Params {
		out += " " + p.String()
	}
	return out
}

// ComponentStatement represents a component tag: {% component "Button" .props %}
type ComponentStatement struct {
	Token lexer.Token // The COMPONENT token
	Name  *StringLiteral
	Props []Expression
}

func (cs *ComponentStatement) statementNode()       {}
func (cs *ComponentStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ComponentStatement) String() string {
	var out string
	out += "component " + cs.Name.String()
	for _, p := range cs.Props {
		out += " " + p.String()
	}
	return out
}

// ExtendsStatement represents an extends tag: {% extends "layout" %}
type ExtendsStatement struct {
	Token  lexer.Token // The EXTENDS token
	Layout *StringLiteral
}

func (es *ExtendsStatement) statementNode()       {}
func (es *ExtendsStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExtendsStatement) String() string {
	return "extends " + es.Layout.String()
}

// BlockTagStatement represents a block tag: {% block name %}
type BlockTagStatement struct {
	Token lexer.Token // The BLOCK token
	Name  *Identifier
}

func (bts *BlockTagStatement) statementNode()       {}
func (bts *BlockTagStatement) TokenLiteral() string { return bts.Token.Literal }
func (bts *BlockTagStatement) String() string {
	return "block " + bts.Name.String()
}

// ContentStatement represents a content tag: {% content %}
type ContentStatement struct {
	Token lexer.Token // The CONTENT token
}

func (cs *ContentStatement) statementNode()       {}
func (cs *ContentStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContentStatement) String() string       { return "content" }

// EndStatement represents an end tag: {% end %}
type EndStatement struct {
	Token lexer.Token // The END token
}

func (es *EndStatement) statementNode()       {}
func (es *EndStatement) TokenLiteral() string { return es.Token.Literal }
func (es *EndStatement) String() string       { return "end" }

// ElseStatement represents an else tag: {% else %}
type ElseStatement struct {
	Token lexer.Token // The ELSE token
}

func (es *ElseStatement) statementNode()       {}
func (es *ElseStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ElseStatement) String() string       { return "else" }

// ElsifStatement represents an elsif tag: {% elsif .condition %}
type ElsifStatement struct {
	Token     lexer.Token // The ELSIF token
	Condition Expression
}

func (es *ElsifStatement) statementNode()       {}
func (es *ElsifStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ElsifStatement) String() string {
	return "elsif " + es.Condition.String()
}

// WhenStatement represents a when tag: {% when "value" %}
type WhenStatement struct {
	Token lexer.Token // The WHEN token
	Value Expression
}

func (ws *WhenStatement) statementNode()       {}
func (ws *WhenStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhenStatement) String() string {
	return "when " + ws.Value.String()
}

// RawStatement represents a raw block: {% raw %}
type RawStatement struct {
	Token   lexer.Token // The RAW token
	Content string      // Literal content between {% raw %} and {% endraw %}
}

func (rs *RawStatement) statementNode()       {}
func (rs *RawStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *RawStatement) String() string       { return "raw " + rs.Content }

// CommentStatement represents a comment block: {% comment %}
type CommentStatement struct {
	Token lexer.Token // The COMMENT token
}

func (cs *CommentStatement) statementNode()       {}
func (cs *CommentStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *CommentStatement) String() string       { return "comment" }

// TextContent represents literal text content between Wisp tags
type TextContent struct {
	Token lexer.Token // The text token
	Value string
}

func (tc *TextContent) statementNode()       {}
func (tc *TextContent) TokenLiteral() string { return tc.Token.Literal }
func (tc *TextContent) String() string       { return tc.Value }
