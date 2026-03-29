// pkg/grove/engine.go
package grove

import (
	"context"

	"grove/internal/compiler"
	"grove/internal/groverrors"
	"grove/internal/lexer"
	"grove/internal/parser"
	"grove/internal/vm"
)

// Option configures an Engine at creation time.
type Option func(*engineCfg)

type engineCfg struct {
	strictVariables bool
}

// WithStrictVariables makes undefined variable references return a RuntimeError.
// Default: false — undefined variables render as empty string.
func WithStrictVariables(strict bool) Option {
	return func(c *engineCfg) { c.strictVariables = strict }
}

// Engine is the Grove template engine. Create with New(). Safe for concurrent use.
type Engine struct {
	cfg     engineCfg
	globals map[string]any
	filters map[string]any // vm.FilterFn | *vm.FilterDef
}

// New creates a configured Engine. Register built-in filters here.
func New(opts ...Option) *Engine {
	e := &Engine{
		globals: make(map[string]any),
		filters: make(map[string]any),
	}
	for _, o := range opts {
		o(&e.cfg)
	}
	// Built-in filters
	e.filters["safe"] = vm.FilterFn(func(v vm.Value, _ []vm.Value) (vm.Value, error) {
		return vm.SafeHTMLVal(v.String()), nil
	})
	return e
}

// SetGlobal registers a value available in all render calls on this engine.
// Render-context data overrides globals; local scope overrides render context.
func (e *Engine) SetGlobal(key string, value any) {
	e.globals[key] = value
}

// RegisterFilter registers a custom filter function.
// fn may be a vm.FilterFn, func(Value, []Value)(Value, error), or *vm.FilterDef
// (created via grove.FilterFunc(fn, grove.FilterOutputsHTML())).
func (e *Engine) RegisterFilter(name string, fn any) {
	e.filters[name] = fn
}

// RenderTemplate compiles and renders an inline template string.
// This is the primary entry point for Plan 1.
// Restrictions: {% extends %} and {% import %} are ParseErrors in inline mode;
// use eng.Render() with a MemoryStore (Plan 4) for templates that need composition.
func (e *Engine) RenderTemplate(ctx context.Context, src string, data Data) (RenderResult, error) {
	// 1. Lex
	tokens, err := lexer.Tokenize(src)
	if err != nil {
		line := 0
		type liner interface{ LexLine() int }
		if le, ok := err.(liner); ok {
			line = le.LexLine()
		}
		return RenderResult{}, &groverrors.ParseError{
			Message: err.Error(),
			Line:    line,
		}
	}

	// 2. Parse (inline=true — forbids extends/import)
	prog, err := parser.Parse(tokens, true)
	if err != nil {
		return RenderResult{}, err // already *groverrors.ParseError
	}

	// 3. Compile
	bc, err := compiler.Compile(prog)
	if err != nil {
		return RenderResult{}, &groverrors.ParseError{Message: err.Error()}
	}

	// 4. Execute
	body, err := vm.Execute(ctx, bc, map[string]any(data), e)
	if err != nil {
		// Wrap vm-internal error into RuntimeError
		if _, ok := err.(*groverrors.RuntimeError); ok {
			return RenderResult{}, err
		}
		return RenderResult{}, &groverrors.RuntimeError{Message: err.Error()}
	}

	return RenderResult{Body: body}, nil
}

// ─── vm.EngineIface implementation ───────────────────────────────────────────

// LookupFilter resolves a filter by name. Returns (nil, false) if not registered.
func (e *Engine) LookupFilter(name string) (vm.FilterFn, bool) {
	v, ok := e.filters[name]
	if !ok {
		return nil, false
	}
	switch f := v.(type) {
	case vm.FilterFn:
		return f, true
	case func(vm.Value, []vm.Value) (vm.Value, error):
		return vm.FilterFn(f), true
	case *vm.FilterDef:
		return f.Fn, true
	}
	return nil, false
}

// StrictVariables reports whether undefined variable references should error.
func (e *Engine) StrictVariables() bool { return e.cfg.strictVariables }

// GlobalData returns the engine-level global variables.
func (e *Engine) GlobalData() map[string]any { return e.globals }
