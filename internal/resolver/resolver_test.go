package resolver

import (
	"testing"

	"template-wisp/internal/ast"
	"template-wisp/internal/lexer"
	"template-wisp/internal/parser"
	"template-wisp/internal/scope"
)

func TestResolveIdentifier(t *testing.T) {
	s := scope.NewScope()
	defer s.Release()

	s.Set("name", "John")
	s.Set("age", 30)

	r := NewResolver(s)

	// Test resolving name
	ident := &ast.Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "name"}, Value: "name"}
	val, err := r.ResolveIdentifier(ident)
	if err != nil {
		t.Errorf("ResolveIdentifier failed: %v", err)
	}
	if val != "John" {
		t.Errorf("Expected 'John', got %v", val)
	}

	// Test resolving age
	ident = &ast.Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "age"}, Value: "age"}
	val, err = r.ResolveIdentifier(ident)
	if err != nil {
		t.Errorf("ResolveIdentifier failed: %v", err)
	}
	if val != 30 {
		t.Errorf("Expected 30, got %v", val)
	}

	// Test resolving undefined variable
	ident = &ast.Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "undefined"}, Value: "undefined"}
	_, err = r.ResolveIdentifier(ident)
	if err == nil {
		t.Error("Expected error for undefined variable")
	}
}

func TestResolveDotExpression(t *testing.T) {
	s := scope.NewScope()
	defer s.Release()

	// Create a user struct
	type User struct {
		Name string
		Age  int
	}
	user := User{Name: "Alice", Age: 25}
	s.Set("user", user)

	r := NewResolver(s)

	// Test resolving .user.Name
	dot := &ast.DotExpression{
		Token: lexer.Token{Type: lexer.DOT, Literal: "."},
		Field: &ast.Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "user"}, Value: "user"},
		Chain: []*ast.Identifier{
			{Token: lexer.Token{Type: lexer.IDENT, Literal: "Name"}, Value: "Name"},
		},
	}

	val, err := r.ResolveDotExpression(dot)
	if err != nil {
		t.Errorf("ResolveDotExpression failed: %v", err)
	}
	if val != "Alice" {
		t.Errorf("Expected 'Alice', got %v", val)
	}

	// Test resolving .user.Age
	dot.Chain = []*ast.Identifier{
		{Token: lexer.Token{Type: lexer.IDENT, Literal: "Age"}, Value: "Age"},
	}
	val, err = r.ResolveDotExpression(dot)
	if err != nil {
		t.Errorf("ResolveDotExpression failed: %v", err)
	}
	if val != 25 {
		t.Errorf("Expected 25, got %v", val)
	}
}

func TestAccessMember(t *testing.T) {
	r := NewResolver(nil)

	// Test struct field access
	type Person struct {
		Name string
	}
	person := Person{Name: "Bob"}
	val, err := r.AccessMember(person, "Name")
	if err != nil {
		t.Errorf("AccessMember failed: %v", err)
	}
	if val != "Bob" {
		t.Errorf("Expected 'Bob', got %v", val)
	}

	// Test map access
	m := map[string]interface{}{
		"key": "value",
	}
	val, err = r.AccessMember(m, "key")
	if err != nil {
		t.Errorf("AccessMember failed: %v", err)
	}
	if val != "value" {
		t.Errorf("Expected 'value', got %v", val)
	}

	// Test nil object
	_, err = r.AccessMember(nil, "field")
	if err == nil {
		t.Error("Expected error for nil object")
	}
}

func TestAccessIndex(t *testing.T) {
	r := NewResolver(nil)

	// Test slice access
	slice := []interface{}{1, 2, 3, 4, 5}
	val, err := r.AccessIndex(slice, 2)
	if err != nil {
		t.Errorf("AccessIndex failed: %v", err)
	}
	if val != 3 {
		t.Errorf("Expected 3, got %v", val)
	}

	// Test map access
	m := map[string]interface{}{
		"key": "value",
	}
	val, err = r.AccessIndex(m, "key")
	if err != nil {
		t.Errorf("AccessIndex failed: %v", err)
	}
	if val != "value" {
		t.Errorf("Expected 'value', got %v", val)
	}

	// Test string access
	str := "hello"
	val, err = r.AccessIndex(str, 1)
	if err != nil {
		t.Errorf("AccessIndex failed: %v", err)
	}
	if val != "e" {
		t.Errorf("Expected 'e', got %v", val)
	}

	// Test out of range
	_, err = r.AccessIndex(slice, 10)
	if err == nil {
		t.Error("Expected error for out of range index")
	}
}

func TestApplyOperator(t *testing.T) {
	r := NewResolver(nil)

	// Test addition
	val, err := r.ApplyOperator("+", 5, 3)
	if err != nil {
		t.Errorf("ApplyOperator failed: %v", err)
	}
	if val != int64(8) {
		t.Errorf("Expected 8, got %v", val)
	}

	// Test string concatenation
	val, err = r.ApplyOperator("+", "hello", " world")
	if err != nil {
		t.Errorf("ApplyOperator failed: %v", err)
	}
	if val != "hello world" {
		t.Errorf("Expected 'hello world', got %v", val)
	}

	// Test subtraction
	val, err = r.ApplyOperator("-", 10, 3)
	if err != nil {
		t.Errorf("ApplyOperator failed: %v", err)
	}
	if val != int64(7) {
		t.Errorf("Expected 7, got %v", val)
	}

	// Test multiplication
	val, err = r.ApplyOperator("*", 4, 5)
	if err != nil {
		t.Errorf("ApplyOperator failed: %v", err)
	}
	if val != int64(20) {
		t.Errorf("Expected 20, got %v", val)
	}

	// Test division
	val, err = r.ApplyOperator("/", 10.0, 2.0)
	if err != nil {
		t.Errorf("ApplyOperator failed: %v", err)
	}
	if val != 5.0 {
		t.Errorf("Expected 5.0, got %v", val)
	}

	// Test equality
	val, err = r.ApplyOperator("==", 5, 5)
	if err != nil {
		t.Errorf("ApplyOperator failed: %v", err)
	}
	if val != true {
		t.Errorf("Expected true, got %v", val)
	}

	// Test less than
	val, err = r.ApplyOperator("<", 3, 5)
	if err != nil {
		t.Errorf("ApplyOperator failed: %v", err)
	}
	if val != true {
		t.Errorf("Expected true, got %v", val)
	}
}

func TestToNumber(t *testing.T) {
	r := NewResolver(nil)

	// Test int
	num, err := r.ToNumber(42)
	if err != nil {
		t.Errorf("ToNumber failed: %v", err)
	}
	if num != 42.0 {
		t.Errorf("Expected 42.0, got %v", num)
	}

	// Test float
	num, err = r.ToNumber(3.14)
	if err != nil {
		t.Errorf("ToNumber failed: %v", err)
	}
	if num != 3.14 {
		t.Errorf("Expected 3.14, got %v", num)
	}

	// Test string
	num, err = r.ToNumber("123")
	if err != nil {
		t.Errorf("ToNumber failed: %v", err)
	}
	if num != 123.0 {
		t.Errorf("Expected 123.0, got %v", num)
	}

	// Test bool
	num, err = r.ToNumber(true)
	if err != nil {
		t.Errorf("ToNumber failed: %v", err)
	}
	if num != 1.0 {
		t.Errorf("Expected 1.0, got %v", num)
	}
}

func TestToString(t *testing.T) {
	r := NewResolver(nil)

	// Test int
	str, err := r.ToString(42)
	if err != nil {
		t.Errorf("ToString failed: %v", err)
	}
	if str != "42" {
		t.Errorf("Expected '42', got %v", str)
	}

	// Test float
	str, err = r.ToString(3.14)
	if err != nil {
		t.Errorf("ToString failed: %v", err)
	}
	if str != "3.14" {
		t.Errorf("Expected '3.14', got %v", str)
	}

	// Test bool
	str, err = r.ToString(true)
	if err != nil {
		t.Errorf("ToString failed: %v", err)
	}
	if str != "true" {
		t.Errorf("Expected 'true', got %v", str)
	}

	// Test string
	str, err = r.ToString("hello")
	if err != nil {
		t.Errorf("ToString failed: %v", err)
	}
	if str != "hello" {
		t.Errorf("Expected 'hello', got %v", str)
	}
}

func TestToBool(t *testing.T) {
	r := NewResolver(nil)

	// Test bool
	b, err := r.ToBool(true)
	if err != nil {
		t.Errorf("ToBool failed: %v", err)
	}
	if b != true {
		t.Errorf("Expected true, got %v", b)
	}

	// Test int
	b, err = r.ToBool(1)
	if err != nil {
		t.Errorf("ToBool failed: %v", err)
	}
	if b != true {
		t.Errorf("Expected true, got %v", b)
	}

	b, err = r.ToBool(0)
	if err != nil {
		t.Errorf("ToBool failed: %v", err)
	}
	if b != false {
		t.Errorf("Expected false, got %v", b)
	}

	// Test string
	b, err = r.ToBool("hello")
	if err != nil {
		t.Errorf("ToBool failed: %v", err)
	}
	if b != true {
		t.Errorf("Expected true, got %v", b)
	}

	b, err = r.ToBool("")
	if err != nil {
		t.Errorf("ToBool failed: %v", err)
	}
	if b != false {
		t.Errorf("Expected false, got %v", b)
	}
}

func TestResolveExpression(t *testing.T) {
	s := scope.NewScope()
	defer s.Release()

	s.Set("x", 10)
	s.Set("y", 20)

	r := NewResolver(s)

	// Test parsing and resolving an expression (using Wisp syntax)
	l := lexer.NewLexer("{% .x + .y %}")
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	val, err := r.ResolveExpression(stmt.Expression)
	if err != nil {
		t.Errorf("ResolveExpression failed: %v", err)
	}
	if val != int64(30) {
		t.Errorf("Expected 30, got %v", val)
	}
}
