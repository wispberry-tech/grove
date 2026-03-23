package lexer

import (
	"fmt"
	"testing"
)

func TestDebugTokens(t *testing.T) {
	input := `{% .user.name %}`

	l := NewLexer(input)

	fmt.Println("Input:", input)
	fmt.Println("Input bytes:", []byte(input))
	fmt.Println("Tokens:")
	for i := 0; i < 20; i++ {
		fmt.Printf("  Before token %d: Pos=%d, ReadPos=%d, Ch=%q\n", i, l.position, l.readPosition, string(l.ch))
		tok := l.NextToken()
		fmt.Printf("  Token %d: Type=%s, Literal=%q\n", i, tok.Type, tok.Literal)
		if tok.Type == EOF {
			break
		}
	}
}
