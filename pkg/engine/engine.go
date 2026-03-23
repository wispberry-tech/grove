// Package engine provides the public API for the Wisp template engine.
package engine

import (
	"fmt"
	"os"
	"sync"

	"template-wisp/internal/ast"
	"template-wisp/internal/evaluator"
	"template-wisp/internal/lexer"
	"template-wisp/internal/parser"
	"template-wisp/internal/scope"
	"template-wisp/internal/store"
)

// Engine is the main entry point for rendering Wisp templates.
type Engine struct {
	store      store.TemplateStore
	filters    map[string]interface{}
	mu         sync.RWMutex
	autoEscape bool                    // whether to auto-escape HTML in output (default: true)
	maxIter    int                     // max loop iterations (0 = unlimited)
	cache      map[string]*ast.Program // parsed template cache
}

// New creates a new Wisp template engine with auto-escaping enabled.
func New() *Engine {
	e := &Engine{
		filters:    make(map[string]interface{}),
		autoEscape: true, // safe by default
		cache:      make(map[string]*ast.Program),
	}
	e.registerBuiltinFilters()
	return e
}

// NewUnsafe creates a new engine without auto-escaping (for non-HTML output).
func NewUnsafe() *Engine {
	e := &Engine{
		filters:    make(map[string]interface{}),
		autoEscape: false,
		cache:      make(map[string]*ast.Program),
	}
	e.registerBuiltinFilters()
	return e
}

// SetAutoEscape enables or disables HTML auto-escaping.
func (e *Engine) SetAutoEscape(enabled bool) {
	e.autoEscape = enabled
}

// SetMaxIterations sets the maximum number of loop iterations.
// Set to 0 for unlimited (not recommended for user templates).
func (e *Engine) SetMaxIterations(max int) {
	e.maxIter = max
}

// NewWithStore creates a new engine with a custom template store.
func NewWithStore(store store.TemplateStore) *Engine {
	e := New()
	e.store = store
	return e
}

// Validate checks if a template string is syntactically correct.
func (e *Engine) Validate(template string) error {
	l := lexer.NewLexer(template)
	p := parser.NewParser(l)
	p.ParseProgram()

	if len(p.Errors()) > 0 {
		return fmt.Errorf("parse errors: %v", p.Errors())
	}

	return nil
}

// RenderString renders a template string with the given data.
func (e *Engine) RenderString(template string, data map[string]interface{}) (string, error) {
	// Try to get parsed program from cache
	e.mu.RLock()
	program, ok := e.cache[template]
	e.mu.RUnlock()

	if !ok {
		// Parse and cache
		l := lexer.NewLexer(template)
		p := parser.NewParser(l)
		program = p.ParseProgram()

		if len(p.Errors()) > 0 {
			return "", fmt.Errorf("parse errors: %v", p.Errors())
		}

		e.mu.Lock()
		e.cache[template] = program
		e.mu.Unlock()
	}

	s := scope.NewScope()
	defer s.Release()

	// Set data variables in scope
	for k, v := range data {
		s.Set(k, v)
	}

	// Register filters as functions
	e.mu.RLock()
	for name, fn := range e.filters {
		s.SetFunction(name, fn)
	}
	e.mu.RUnlock()

	eval := evaluator.NewEvaluator(s)
	eval.SetAutoEscape(e.autoEscape)
	eval.SetMaxIterations(e.maxIter)

	// Set up template loader if store is available
	if e.store != nil {
		eval.SetTemplateFn(func(name string) (string, error) {
			content, err := e.store.ReadTemplate(name)
			if err != nil {
				return "", err
			}
			return string(content), nil
		})
	}

	return eval.Evaluate(program)
}

// ClearCache clears the parsed template cache.
func (e *Engine) ClearCache() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cache = make(map[string]*ast.Program)
}

// RenderFile renders a template file with the given data.
func (e *Engine) RenderFile(filename string, data map[string]interface{}) (string, error) {
	if e.store == nil {
		// Default: read from filesystem
		content, err := os.ReadFile(filename)
		if err != nil {
			return "", fmt.Errorf("failed to read template file %s: %w", filename, err)
		}
		return e.RenderString(string(content), data)
	}

	content, err := e.store.ReadTemplate(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", filename, err)
	}

	return e.RenderString(string(content), data)
}

// RegisterFilter registers a custom filter function.
func (e *Engine) RegisterFilter(name string, fn interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.filters[name] = fn
}

// RegisterTemplate registers a template in the store.
func (e *Engine) RegisterTemplate(name string, content string) {
	if e.store == nil {
		e.store = store.NewMemoryStore()
	}
	if ms, ok := e.store.(*store.MemoryStore); ok {
		ms.Register(name, content)
	}
}
