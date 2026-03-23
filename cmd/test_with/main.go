package main

import (
	"fmt"
	"template-wisp/internal/lexer"
)

func main() {
	input := `{% with .user as .currentUser %}`
	fmt.Println("Input:", input)

	l := lexer.NewLexer(input)

	for i := 0; i < 20; i++ {
		tok := l.NextToken()
		fmt.Printf("Token %d: Type=%s, Literal=%q\n", i, tok.Type, tok.Literal)
		if tok.Type == lexer.EOF {
			break
		}
	}
}
