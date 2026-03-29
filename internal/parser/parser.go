// internal/parser/parser.go
package parser

import (
	"fmt"
	"strconv"

	"grove/internal/ast"
	"grove/internal/groverrors"
	"grove/internal/lexer"
)

// Parse converts a token stream into an AST.
// inline=true forbids {% extends %} and {% import %} (used by RenderTemplate).
func Parse(tokens []lexer.Token, inline bool) (*ast.Program, error) {
	p := &parser{tokens: tokens, inline: inline}
	return p.parseProgram()
}

type parser struct {
	tokens []lexer.Token
	pos    int
	inline bool
}

// ─── Program ──────────────────────────────────────────────────────────────────

func (p *parser) parseProgram() (*ast.Program, error) {
	prog := &ast.Program{}
	for !p.atEOF() {
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			prog.Body = append(prog.Body, node)
		}
	}
	return prog, nil
}

func (p *parser) parseNode() (ast.Node, error) {
	tk := p.peek()
	switch tk.Kind {
	case lexer.TK_TEXT:
		p.advance()
		return &ast.TextNode{Value: tk.Value, Line: tk.Line}, nil
	case lexer.TK_OUTPUT_START:
		return p.parseOutput()
	case lexer.TK_TAG_START:
		return p.parseTag()
	case lexer.TK_EOF:
		return nil, nil
	default:
		return nil, p.errorf(tk.Line, tk.Col, "unexpected token %q", tk.Value)
	}
}

// ─── Output {{ expr }} ────────────────────────────────────────────────────────

func (p *parser) parseOutput() (*ast.OutputNode, error) {
	start := p.advance() // consume OUTPUT_START
	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	end := p.peek()
	if end.Kind != lexer.TK_OUTPUT_END {
		return nil, p.errorf(end.Line, end.Col, "expected }}, got %q", end.Value)
	}
	p.advance() // consume OUTPUT_END
	return &ast.OutputNode{
		Expr:       expr,
		StripLeft:  start.StripLeft,
		StripRight: end.StripRight,
		Line:       start.Line,
	}, nil
}

// ─── Tags {% name ... %} ──────────────────────────────────────────────────────

func (p *parser) parseTag() (ast.Node, error) {
	tagStart := p.advance() // consume TAG_START
	name := p.peek()
	if name.Kind != lexer.TK_IDENT {
		return nil, p.errorf(name.Line, name.Col, "expected tag name after {%%")
	}

	switch name.Value {
	case "raw":
		// {% raw %} was already handled at the lexer level and emitted as TK_TEXT.
		// If we reach here, the lexer emitted {% raw %} as TAG_START + IDENT + TAG_END
		// (i.e. an empty raw block or parser called incorrectly). Consume and return empty.
		p.advance() // consume "raw"
		if p.peek().Kind != lexer.TK_TAG_END {
			return nil, p.errorf(p.peek().Line, p.peek().Col, "expected %%} after raw")
		}
		p.advance() // consume TAG_END
		// Collect TK_TEXT until {% endraw %}
		return p.consumeUntilEndraw(tagStart)

	case "extends":
		if p.inline {
			return nil, &groverrors.ParseError{
				Line:    name.Line,
				Column:  name.Col,
				Message: "extends not allowed in inline templates",
			}
		}
		return p.consumeTagRemainder(name.Value, tagStart)

	case "import":
		if p.inline {
			return nil, &groverrors.ParseError{
				Line:    name.Line,
				Column:  name.Col,
				Message: "import not allowed in inline templates",
			}
		}
		return p.consumeTagRemainder(name.Value, tagStart)

	default:
		return p.consumeTagRemainder(name.Value, tagStart)
	}
}

// consumeTagRemainder skips to TAG_END and emits a TagNode (stub for unimplemented tags).
func (p *parser) consumeTagRemainder(name string, tagStart lexer.Token) (ast.Node, error) {
	p.advance() // consume tag name
	for p.peek().Kind != lexer.TK_TAG_END && !p.atEOF() {
		p.advance()
	}
	if p.peek().Kind == lexer.TK_TAG_END {
		p.advance() // consume TAG_END
	}
	return &ast.TagNode{Name: name, Line: tagStart.Line}, nil
}

// consumeUntilEndraw is the fallback for raw blocks not handled at lex time.
func (p *parser) consumeUntilEndraw(tagStart lexer.Token) (ast.Node, error) {
	var content string
	for !p.atEOF() {
		tk := p.peek()
		if tk.Kind == lexer.TK_TAG_START {
			// Look ahead for endraw
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Kind == lexer.TK_IDENT &&
				p.tokens[p.pos+1].Value == "endraw" {
				p.advance() // TAG_START
				p.advance() // endraw
				if p.peek().Kind == lexer.TK_TAG_END {
					p.advance() // TAG_END
				}
				return &ast.RawNode{Value: content, Line: tagStart.Line}, nil
			}
		}
		if tk.Kind == lexer.TK_TEXT {
			content += tk.Value
		}
		p.advance()
	}
	return nil, p.errorf(tagStart.Line, tagStart.Col, "unclosed raw block")
}

// ─── Expression parsing (Pratt) ───────────────────────────────────────────────

// parseExpr parses an expression with minimum precedence minPrec.
func (p *parser) parseExpr(minPrec int) (ast.Node, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for {
		tk := p.peek()
		prec, isInfix := infixPrec(tk.Kind)
		if !isInfix || prec <= minPrec {
			break
		}

		switch tk.Kind {
		case lexer.TK_IF:
			// Inline ternary: consequence if condition else alternative
			// left is already the consequence
			p.advance() // consume if
			cond, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			if p.peek().Kind != lexer.TK_ELSE {
				return nil, p.errorf(p.peek().Line, p.peek().Col, "expected 'else' in ternary expression")
			}
			p.advance() // consume else
			alt, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			left = &ast.TernaryExpr{
				Condition:   cond,
				Consequence: left,
				Alternative: alt,
				Line:        tk.Line,
			}

		case lexer.TK_PIPE:
			p.advance() // consume |
			left, err = p.parseFilter(left)
			if err != nil {
				return nil, err
			}

		case lexer.TK_DOT:
			p.advance() // consume .
			attr := p.peek()
			if attr.Kind != lexer.TK_IDENT {
				return nil, p.errorf(attr.Line, attr.Col, "expected attribute name after .")
			}
			p.advance()
			left = &ast.AttributeAccess{Object: left, Key: attr.Value, Line: attr.Line}

		case lexer.TK_LBRACKET:
			p.advance() // consume [
			idx, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			if p.peek().Kind != lexer.TK_RBRACKET {
				return nil, p.errorf(p.peek().Line, p.peek().Col, "expected ]")
			}
			p.advance() // consume ]
			left = &ast.IndexAccess{Object: left, Key: idx, Line: tk.Line}

		default:
			// Binary operator
			p.advance() // consume operator
			right, err := p.parseExpr(prec)
			if err != nil {
				return nil, err
			}
			left = &ast.BinaryExpr{Op: tk.Value, Left: left, Right: right, Line: tk.Line}
		}
	}
	return left, nil
}

func (p *parser) parseUnary() (ast.Node, error) {
	tk := p.peek()
	switch tk.Kind {
	case lexer.TK_NOT:
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: "not", Operand: operand, Line: tk.Line}, nil
	case lexer.TK_MINUS:
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: "-", Operand: operand, Line: tk.Line}, nil
	}
	return p.parsePrimary()
}

func (p *parser) parsePrimary() (ast.Node, error) {
	tk := p.advance()
	switch tk.Kind {
	case lexer.TK_NIL:
		return &ast.NilLiteral{Line: tk.Line}, nil
	case lexer.TK_TRUE:
		return &ast.BoolLiteral{Value: true, Line: tk.Line}, nil
	case lexer.TK_FALSE:
		return &ast.BoolLiteral{Value: false, Line: tk.Line}, nil
	case lexer.TK_STRING:
		return &ast.StringLiteral{Value: tk.Value, Line: tk.Line}, nil
	case lexer.TK_INT:
		n, err := strconv.ParseInt(tk.Value, 10, 64)
		if err != nil {
			return nil, p.errorf(tk.Line, tk.Col, "invalid integer: %s", tk.Value)
		}
		return &ast.IntLiteral{Value: n, Line: tk.Line}, nil
	case lexer.TK_FLOAT:
		f, err := strconv.ParseFloat(tk.Value, 64)
		if err != nil {
			return nil, p.errorf(tk.Line, tk.Col, "invalid float: %s", tk.Value)
		}
		return &ast.FloatLiteral{Value: f, Line: tk.Line}, nil
	case lexer.TK_IDENT:
		return &ast.Identifier{Name: tk.Value, Line: tk.Line}, nil
	case lexer.TK_LPAREN:
		expr, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		if p.peek().Kind != lexer.TK_RPAREN {
			return nil, p.errorf(p.peek().Line, p.peek().Col, "expected )")
		}
		p.advance()
		return expr, nil
	default:
		return nil, p.errorf(tk.Line, tk.Col, "unexpected token in expression: %q", tk.Value)
	}
}

func (p *parser) parseFilter(value ast.Node) (ast.Node, error) {
	name := p.peek()
	if name.Kind != lexer.TK_IDENT {
		return nil, p.errorf(name.Line, name.Col, "expected filter name after |")
	}
	p.advance()

	var args []ast.Node
	if p.peek().Kind == lexer.TK_LPAREN {
		p.advance() // consume (
		for p.peek().Kind != lexer.TK_RPAREN && !p.atEOF() {
			arg, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
			if p.peek().Kind == lexer.TK_COMMA {
				p.advance()
			}
		}
		if p.peek().Kind != lexer.TK_RPAREN {
			return nil, p.errorf(p.peek().Line, p.peek().Col, "expected ) after filter arguments")
		}
		p.advance() // consume )
	}

	return &ast.FilterExpr{
		Value:  value,
		Filter: name.Value,
		Args:   args,
		Line:   name.Line,
	}, nil
}

// infixPrec returns the infix precedence for a token kind.
// Higher precedence = binds tighter.
func infixPrec(k lexer.TokenKind) (int, bool) {
	switch k {
	case lexer.TK_IF:
		return 5, true // inline ternary — lowest
	case lexer.TK_OR:
		return 10, true
	case lexer.TK_AND:
		return 20, true
	case lexer.TK_EQ, lexer.TK_NEQ, lexer.TK_LT, lexer.TK_LTE, lexer.TK_GT, lexer.TK_GTE:
		return 40, true
	case lexer.TK_TILDE:
		return 50, true
	case lexer.TK_PLUS, lexer.TK_MINUS:
		return 60, true
	case lexer.TK_STAR, lexer.TK_SLASH, lexer.TK_PERCENT:
		return 70, true
	case lexer.TK_PIPE:
		return 90, true
	case lexer.TK_DOT, lexer.TK_LBRACKET:
		return 100, true
	}
	return 0, false
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func (p *parser) peek() lexer.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return lexer.Token{Kind: lexer.TK_EOF}
}

func (p *parser) advance() lexer.Token {
	tk := p.peek()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return tk
}

func (p *parser) atEOF() bool {
	return p.pos >= len(p.tokens) || p.tokens[p.pos].Kind == lexer.TK_EOF
}

func (p *parser) errorf(line, col int, format string, args ...any) *groverrors.ParseError {
	return &groverrors.ParseError{
		Line:    line,
		Column:  col,
		Message: fmt.Sprintf(format, args...),
	}
}
