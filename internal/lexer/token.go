package lexer

// TokenType represents the type of a token.
type TokenType string

const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Literals
	IDENT  TokenType = "IDENT" // identifier
	STRING TokenType = "STRING"
	NUMBER TokenType = "NUMBER"

	// Operators
	ASSIGN_OP TokenType = "="
	PLUS      TokenType = "+"
	MINUS     TokenType = "-"
	BANG      TokenType = "!"
	ASTERISK  TokenType = "*"
	SLASH     TokenType = "/"
	EQ        TokenType = "=="
	NOT_EQ    TokenType = "!="
	LT        TokenType = "<"
	LTE       TokenType = "<="
	GT        TokenType = ">"
	GTE       TokenType = ">="

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"

	// Keywords
	LET        TokenType = "LET"
	IF         TokenType = "IF"
	ELSE       TokenType = "ELSE"
	ELSIF      TokenType = "ELSIF"
	UNLESS     TokenType = "UNLESS"
	END        TokenType = "END"
	TRUE       TokenType = "TRUE"
	FALSE      TokenType = "FALSE"
	RETURN     TokenType = "RETURN"
	FUNCTION   TokenType = "FUNCTION"
	ASSIGN     TokenType = "ASSIGN"
	FOR        TokenType = "FOR"
	WHILE      TokenType = "WHILE"
	RANGE      TokenType = "RANGE"
	CASE       TokenType = "CASE"
	WHEN       TokenType = "WHEN"
	WITH       TokenType = "WITH"
	CYCLE      TokenType = "CYCLE"
	INCREMENT  TokenType = "INCREMENT"
	DECREMENT  TokenType = "DECREMENT"
	BREAK      TokenType = "BREAK"
	CONTINUE   TokenType = "CONTINUE"
	INCLUDE    TokenType = "INCLUDE"
	RENDER     TokenType = "RENDER"
	COMPONENT  TokenType = "COMPONENT"
	EXTENDS    TokenType = "EXTENDS"
	BLOCK      TokenType = "BLOCK"
	CONTENT    TokenType = "CONTENT"
	RAW        TokenType = "RAW"
	COMMENT    TokenType = "COMMENT"
	AS         TokenType = "AS"
	IN         TokenType = "IN"
	ENDRAW     TokenType = "ENDRAW"
	ENDCOMMENT TokenType = "ENDCOMMENT"

	// Our unique tokens for the unified bracket syntax
	DOT        TokenType = "DOT"        // .
	PIPE       TokenType = "PIPE"       // |
	LBRACE_PCT TokenType = "LBRACE_PCT" // {%
	RBRACE_PCT TokenType = "RBRACE_PCT" // %}
	TEXT       TokenType = "TEXT"       // literal text content
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// NewToken creates a new token.
func NewToken(t TokenType, l string, line, column int) Token {
	return Token{Type: t, Literal: l, Line: line, Column: column}
}
