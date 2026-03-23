package scope

import (
	"sync"
)

// Scope represents a variable scope with shared parent references.
type Scope struct {
	parent    *Scope
	variables map[string]interface{}
	functions map[string]interface{}
	isRoot    bool
	isolated  bool
}

// scopePool is a pool for reusing Scope objects.
var scopePool = sync.Pool{
	New: func() interface{} {
		return &Scope{
			variables: make(map[string]interface{}),
			functions: make(map[string]interface{}),
		}
	},
}

// NewScope creates a new root scope.
func NewScope() *Scope {
	s := scopePool.Get().(*Scope)
	s.parent = nil
	s.isRoot = true
	s.isolated = false
	// Clear the maps
	for k := range s.variables {
		delete(s.variables, k)
	}
	for k := range s.functions {
		delete(s.functions, k)
	}
	return s
}

// NewChildScope creates a new child scope with the given parent.
func NewChildScope(parent *Scope) *Scope {
	s := scopePool.Get().(*Scope)
	s.parent = parent
	s.isRoot = false
	s.isolated = false
	// Clear the maps
	for k := range s.variables {
		delete(s.variables, k)
	}
	for k := range s.functions {
		delete(s.functions, k)
	}
	return s
}

// NewIsolatedScope creates a new isolated scope (no parent access).
func NewIsolatedScope() *Scope {
	s := scopePool.Get().(*Scope)
	s.parent = nil
	s.isRoot = true
	s.isolated = true
	// Clear the maps
	for k := range s.variables {
		delete(s.variables, k)
	}
	for k := range s.functions {
		delete(s.functions, k)
	}
	return s
}

// Get retrieves a variable by name, walking up the scope chain.
func (s *Scope) Get(name string) (interface{}, bool) {
	// Check current scope first
	if val, ok := s.variables[name]; ok {
		return val, true
	}

	// Walk up the scope chain if not isolated
	if !s.isolated && s.parent != nil {
		return s.parent.Get(name)
	}

	return nil, false
}

// Set sets a variable in the current scope.
func (s *Scope) Set(name string, value interface{}) {
	s.variables[name] = value
}

// GetFunction retrieves a function by name, walking up the scope chain.
func (s *Scope) GetFunction(name string) (interface{}, bool) {
	// Check current scope first
	if fn, ok := s.functions[name]; ok {
		return fn, true
	}

	// Walk up the scope chain if not isolated
	if !s.isolated && s.parent != nil {
		return s.parent.GetFunction(name)
	}

	return nil, false
}

// SetFunction sets a function in the current scope.
func (s *Scope) SetFunction(name string, fn interface{}) {
	s.functions[name] = fn
}

// Delete removes a variable from the current scope.
func (s *Scope) Delete(name string) {
	delete(s.variables, name)
}

// DeleteFunction removes a function from the current scope.
func (s *Scope) DeleteFunction(name string) {
	delete(s.functions, name)
}

// IsIsolated returns true if this scope is isolated (no parent access).
func (s *Scope) IsIsolated() bool {
	return s.isolated
}

// IsRoot returns true if this is a root scope.
func (s *Scope) IsRoot() bool {
	return s.isRoot
}

// Parent returns the parent scope.
func (s *Scope) Parent() *Scope {
	return s.parent
}

// Release returns the scope to the pool for reuse.
func (s *Scope) Release() {
	// Clear the maps
	for k := range s.variables {
		delete(s.variables, k)
	}
	for k := range s.functions {
		delete(s.functions, k)
	}
	s.parent = nil
	s.isRoot = false
	s.isolated = false
	scopePool.Put(s)
}

// Variables returns a copy of all variables in the current scope.
func (s *Scope) Variables() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range s.variables {
		result[k] = v
	}
	return result
}

// Functions returns a copy of all functions in the current scope.
func (s *Scope) Functions() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range s.functions {
		result[k] = v
	}
	return result
}

// HasVariable checks if a variable exists in the current scope (not parent).
func (s *Scope) HasVariable(name string) bool {
	_, ok := s.variables[name]
	return ok
}

// HasFunction checks if a function exists in the current scope (not parent).
func (s *Scope) HasFunction(name string) bool {
	_, ok := s.functions[name]
	return ok
}

// Depth returns the depth of the scope chain.
func (s *Scope) Depth() int {
	depth := 0
	current := s
	for current != nil {
		depth++
		current = current.parent
	}
	return depth
}
