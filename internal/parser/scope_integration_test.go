package parser

import (
	"testing"

	"template-wisp/internal/ast"
	"template-wisp/internal/lexer"
)

func TestScopeIntegration(t *testing.T) {
	input := `{% for .item in .items %}{% .item %}{% end %}`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Logf("parser has %d errors:", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Logf("  parser error: %s", err)
		}
		t.FailNow()
	}

	// With body parsing, program should have 1 statement (the ForStatement)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	// First statement should be a for loop
	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ForStatement. got=%T",
			program.Statements[0])
	}

	if forStmt.LoopVar.Value != "item" {
		t.Errorf("forStmt.LoopVar.Value not 'item'. got=%s", forStmt.LoopVar.Value)
	}

	// The for loop body should contain the expression statement
	if forStmt.Body == nil {
		t.Fatal("forStmt.Body is nil")
	}

	if len(forStmt.Body.Statements) != 1 {
		t.Fatalf("forStmt.Body.Statements does not contain 1 statement. got=%d",
			len(forStmt.Body.Statements))
	}

	// Body statement should be an expression statement with dot expression
	exprStmt, ok := forStmt.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("forStmt.Body.Statements[0] is not ExpressionStatement. got=%T",
			forStmt.Body.Statements[0])
	}

	dotExpr, ok := exprStmt.Expression.(*ast.DotExpression)
	if !ok {
		t.Fatalf("expression is not DotExpression. got=%T", exprStmt.Expression)
	}

	if dotExpr.Field.Value != "item" {
		t.Errorf("dotExpr.Field.Value not 'item'. got=%s", dotExpr.Field.Value)
	}

	// Check that scope was created and destroyed
	if p.CurrentScope() == nil {
		t.Error("Parser should have a current scope")
	}
}
