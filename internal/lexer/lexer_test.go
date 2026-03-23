package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `{% .name %} {% . | date %} {% if .condition %}{% end %}`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LBRACE_PCT, "{%"},
		{DOT, "."},
		{IDENT, "name"},
		{RBRACE_PCT, "%}"},
		{LBRACE_PCT, "{%"},
		{DOT, "."},
		{PIPE, "|"},
		{IDENT, "date"},
		{RBRACE_PCT, "%}"},
		{LBRACE_PCT, "{%"},
		{IF, "if"},
		{DOT, "."},
		{IDENT, "condition"},
		{RBRACE_PCT, "%}"},
		{LBRACE_PCT, "{%"},
		{END, "end"},
		{RBRACE_PCT, "%}"},
		{EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
