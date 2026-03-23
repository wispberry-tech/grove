package parser

import (
	"testing"

	"template-wisp/internal/ast"
	"template-wisp/internal/lexer"
)

func TestWispVariableAccess(t *testing.T) {
	input := `{% .name %}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %s", err)
		}
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
			program.Statements[0])
	}

	dotExpr, ok := stmt.Expression.(*ast.DotExpression)
	if !ok {
		t.Fatalf("expression is not DotExpression. got=%T", stmt.Expression)
	}

	if dotExpr.Field.Value != "name" {
		t.Errorf("dotExpr.Field.Value not 'name'. got=%s", dotExpr.Field.Value)
	}
}

func TestWispVariableAccessChained(t *testing.T) {
	input := `{% .user.name %}`

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

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
			program.Statements[0])
	}

	dotExpr, ok := stmt.Expression.(*ast.DotExpression)
	if !ok {
		t.Fatalf("expression is not DotExpression. got=%T", stmt.Expression)
	}

	if dotExpr.Field.Value != "user" {
		t.Errorf("dotExpr.Field.Value not 'user'. got=%s", dotExpr.Field.Value)
	}

	if len(dotExpr.Chain) != 1 {
		t.Fatalf("dotExpr.Chain does not contain 1 identifier. got=%d",
			len(dotExpr.Chain))
	}

	if dotExpr.Chain[0].Value != "name" {
		t.Errorf("dotExpr.Chain[0].Value not 'name'. got=%s", dotExpr.Chain[0].Value)
	}
}

func TestWispPipeExpression(t *testing.T) {
	input := `{% . | date %}`

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

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
			program.Statements[0])
	}

	pipeExpr, ok := stmt.Expression.(*ast.PipeExpression)
	if !ok {
		t.Fatalf("expression is not PipeExpression. got=%T", stmt.Expression)
	}

	if pipeExpr.Function.Value != "date" {
		t.Errorf("pipeExpr.Function.Value not 'date'. got=%s", pipeExpr.Function.Value)
	}
}

func TestWispAssignStatement(t *testing.T) {
	input := `{% assign .name = "John" %}`

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

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not AssignStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "name" {
		t.Errorf("stmt.Name.Value not 'name'. got=%s", stmt.Name.Value)
	}

	stringLit, ok := stmt.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not StringLiteral. got=%T", stmt.Value)
	}

	if stringLit.Value != "John" {
		t.Errorf("stringLit.Value not 'John'. got=%s", stringLit.Value)
	}
}

func TestWispIfStatement(t *testing.T) {
	input := `{% if .condition %}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %s", err)
		}
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not IfStatement. got=%T",
			program.Statements[0])
	}

	dotExpr, ok := stmt.Condition.(*ast.DotExpression)
	if !ok {
		t.Fatalf("stmt.Condition is not DotExpression. got=%T", stmt.Condition)
	}

	if dotExpr.Field.Value != "condition" {
		t.Errorf("dotExpr.Field.Value not 'condition'. got=%s", dotExpr.Field.Value)
	}
}

func TestWispForStatement(t *testing.T) {
	input := `{% for .item in .items %}`

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

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ForStatement. got=%T",
			program.Statements[0])
	}

	if stmt.LoopVar.Value != "item" {
		t.Errorf("stmt.LoopVar.Value not 'item'. got=%s", stmt.LoopVar.Value)
	}

	dotExpr, ok := stmt.Collection.(*ast.DotExpression)
	if !ok {
		t.Fatalf("stmt.Collection is not DotExpression. got=%T", stmt.Collection)
	}

	if dotExpr.Field.Value != "items" {
		t.Errorf("dotExpr.Field.Value not 'items'. got=%s", dotExpr.Field.Value)
	}
}

func TestWispEndStatement(t *testing.T) {
	input := `{% end %}`

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

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.EndStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not EndStatement. got=%T",
			program.Statements[0])
	}
}
