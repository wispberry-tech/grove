package scope

import (
	"testing"
)

func TestNewScope(t *testing.T) {
	s := NewScope()
	defer s.Release()

	if !s.IsRoot() {
		t.Error("NewScope should create a root scope")
	}

	if s.IsIsolated() {
		t.Error("NewScope should not create an isolated scope")
	}

	if s.Parent() != nil {
		t.Error("Root scope should have no parent")
	}
}

func TestNewChildScope(t *testing.T) {
	parent := NewScope()
	defer parent.Release()

	child := NewChildScope(parent)
	defer child.Release()

	if child.IsRoot() {
		t.Error("Child scope should not be a root scope")
	}

	if child.IsIsolated() {
		t.Error("Child scope should not be isolated")
	}

	if child.Parent() != parent {
		t.Error("Child scope should have parent as parent")
	}
}

func TestNewIsolatedScope(t *testing.T) {
	s := NewIsolatedScope()
	defer s.Release()

	if !s.IsIsolated() {
		t.Error("NewIsolatedScope should create an isolated scope")
	}

	if !s.IsRoot() {
		t.Error("Isolated scope should be a root scope")
	}
}

func TestScopeGetSet(t *testing.T) {
	s := NewScope()
	defer s.Release()

	s.Set("name", "John")
	s.Set("age", 30)

	val, ok := s.Get("name")
	if !ok {
		t.Error("Get should return true for existing variable")
	}
	if val != "John" {
		t.Errorf("Get should return 'John', got %v", val)
	}

	val, ok = s.Get("age")
	if !ok {
		t.Error("Get should return true for existing variable")
	}
	if val != 30 {
		t.Errorf("Get should return 30, got %v", val)
	}

	_, ok = s.Get("nonexistent")
	if ok {
		t.Error("Get should return false for non-existent variable")
	}
}

func TestScopeChain(t *testing.T) {
	parent := NewScope()
	defer parent.Release()

	parent.Set("parentVar", "parent value")

	child := NewChildScope(parent)
	defer child.Release()

	child.Set("childVar", "child value")

	// Child should see both parent and child variables
	val, ok := child.Get("parentVar")
	if !ok {
		t.Error("Child should see parent variable")
	}
	if val != "parent value" {
		t.Errorf("Child should see parent value, got %v", val)
	}

	val, ok = child.Get("childVar")
	if !ok {
		t.Error("Child should see child variable")
	}
	if val != "child value" {
		t.Errorf("Child should see child value, got %v", val)
	}

	// Parent should not see child variables
	_, ok = parent.Get("childVar")
	if ok {
		t.Error("Parent should not see child variable")
	}
}

func TestIsolatedScope(t *testing.T) {
	parent := NewScope()
	defer parent.Release()

	parent.Set("parentVar", "parent value")

	isolated := NewIsolatedScope()
	defer isolated.Release()

	isolated.Set("isolatedVar", "isolated value")

	// Isolated scope should not see parent variables
	_, ok := isolated.Get("parentVar")
	if ok {
		t.Error("Isolated scope should not see parent variable")
	}

	// Isolated scope should see its own variables
	val, ok := isolated.Get("isolatedVar")
	if !ok {
		t.Error("Isolated scope should see its own variable")
	}
	if val != "isolated value" {
		t.Errorf("Isolated scope should see its own value, got %v", val)
	}
}

func TestScopeDelete(t *testing.T) {
	s := NewScope()
	defer s.Release()

	s.Set("name", "John")
	s.Set("age", 30)

	s.Delete("name")

	_, ok := s.Get("name")
	if ok {
		t.Error("Deleted variable should not be found")
	}

	val, ok := s.Get("age")
	if !ok {
		t.Error("Non-deleted variable should still exist")
	}
	if val != 30 {
		t.Errorf("Non-deleted variable should have correct value, got %v", val)
	}
}

func TestScopeFunctions(t *testing.T) {
	s := NewScope()
	defer s.Release()

	s.SetFunction("myFunc", func() {})

	fn, ok := s.GetFunction("myFunc")
	if !ok {
		t.Error("GetFunction should return true for existing function")
	}
	if fn == nil {
		t.Error("GetFunction should return non-nil function")
	}

	_, ok = s.GetFunction("nonexistent")
	if ok {
		t.Error("GetFunction should return false for non-existent function")
	}
}

func TestScopeDepth(t *testing.T) {
	root := NewScope()
	defer root.Release()

	if root.Depth() != 1 {
		t.Errorf("Root scope depth should be 1, got %d", root.Depth())
	}

	child := NewChildScope(root)
	defer child.Release()

	if child.Depth() != 2 {
		t.Errorf("Child scope depth should be 2, got %d", child.Depth())
	}

	grandchild := NewChildScope(child)
	defer grandchild.Release()

	if grandchild.Depth() != 3 {
		t.Errorf("Grandchild scope depth should be 3, got %d", grandchild.Depth())
	}
}

func TestScopeVariables(t *testing.T) {
	s := NewScope()
	defer s.Release()

	s.Set("name", "John")
	s.Set("age", 30)

	vars := s.Variables()

	if len(vars) != 2 {
		t.Errorf("Variables should return 2 variables, got %d", len(vars))
	}

	if vars["name"] != "John" {
		t.Errorf("Variables should contain name='John', got %v", vars["name"])
	}

	if vars["age"] != 30 {
		t.Errorf("Variables should contain age=30, got %v", vars["age"])
	}
}

func TestScopeHasVariable(t *testing.T) {
	s := NewScope()
	defer s.Release()

	s.Set("name", "John")

	if !s.HasVariable("name") {
		t.Error("HasVariable should return true for existing variable")
	}

	if s.HasVariable("nonexistent") {
		t.Error("HasVariable should return false for non-existent variable")
	}
}
