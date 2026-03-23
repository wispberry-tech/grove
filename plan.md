# Wisp Template Engine - Design Plan

## Overview

Wisp is a secure, high-performance HTML templating engine for Go that provides a simplified syntax inspired by both Liquid and Go's pointer/member access patterns. It's designed for web applications, static site generation, and any scenario requiring secure, reusable template rendering.

### Design Goals

- **Simplicity**: Clear unified bracket syntax with explicit function/variable distinction
- **Security**: Built-in sandboxing, auto-escaping, and resource limiting
- **Performance**: Optimized compilation with caching and bytecode compilation
- **Compatibility**: Extensible and framework-agnostic
- **Developer Experience**: Clear syntax, excellent error messages, and development tools

## Syntax Specification

### Unified Syntax

**Variables & Data Access**: `{% .name %}` (leading dot for variable access)
**Function Calls**: `{% . | functionName %}` (leading dot, pipe, function name)
**Tags & Control Structures**: `{% tagName %}` (no dot, no pipe)

```liquid
{# Variable access - leading dot #}
{% .name %}

{# Nested variable access #}
{% .user.name %}

{# Array/map indexing #}
{% .items[0] %}
{% .items[key] %}

{# Function call - leading dot, pipe, function name #}
{% .createdAt | date %}

{# Function call with arguments #}
{% .price | format "%s" | currency "USD" %}

{# Tag - no dot, no pipe #}
{% if .condition %}{% end %}
```

### Member Access

```liquid
{# Basic access #}
{% .name %}

{# Nested access #}
{% .user.name %}

{# Array indexing #}
{% .users[0].name %}

{# Map/Hash indexing #}
{% .items[key] %}

{# Combinations #}
{% .users[index].profile.settings.theme %}
```

### Control Flow Tags

```liquid
{# Conditional #}
{% if .condition %}
  {% .output %}
{% elsif .altCondition %}
  {% .altOutput %}
{% else %}
  {% .defaultOutput %}
{% end %}

{# Unless - negated if #}
{% unless .shouldHide %}
  {% .visibleContent %}
{% end %}

{# Case/When #}
{% case .value %}
  {% when "a" %}
    {% .aContent %}
  {% when "b" %}
    {% .bContent %}
  {% else %}
    {% .defaultContent %}
{% end %}

{# Loops #}
{% for .item in .items %}
  {% .item.name %}
{% end %}

{% range .start .end %}
  {% .i %}
{% end %}

{% while .condition %}
  {% .value %}
{% end %}

{# Additional Control Flow Tags #}

{# Context block - isolate variable scope #}
{% with .user as .currentUser %}
  {% .currentUser.name %}
{% endwith %}

{# Cycle - alternating values #}
{% cycle .items 'odd' 'even' %}

{# Increment/Decrement #}
{% increment .counter %}
{% decrement .counter %}

{# Raw - literal rendering #}
{% raw %}
  <div>{{ .mustache }} {{ .variable }}</div>
{% endraw %}

{# Comment - multi-line #}
{% comment %}
  This is a comment and will not be output
{% endcomment %}

{# Break and Continue in loops #}
{% for .item in .items %}
  {% if .item.skip %}{% continue %}{% end %}
  {% if .item.break %}{% break %}{% end %}
  {% .item.name %}
{% end %}
```

### Variable Manipulation

```liquid
{# Assign variables #}
{% assign .count = 0 %}
{% assign .user.name = "Jane" %}

{# Capture output #}
{% capture .output %}
  {% .name %}'s content
{% endcapture %}

{# Comments #}
{# This is not output to the template #}
```

### Template Composition

```liquid
{# Include template (evaluate and output) #}
{% include "partials/header" %}
{% include "sidebar" .user %}

{# Component-like inclusion #}
{% component "Button" .buttonOptions %}

{# Render template (partial evaluation) #}
{% render "widget" .data %}
```

### Layout System

```liquid
{# Master layout - single extends #}
{% extends "layouts/base" %}

{# Block definition - defines override point #}
{% block title %}{{ .pageTitle }}{% endblock %}

{# Content block - for main content #}
{% content %}
  {% include "main" %}
{% endcontent %}

{# Block override - provides content #}
{% block scripts %}
  {% include "scripts" %}
{% endblock %}
```

## Architecture

### Core Components

#### 1. **Parser Layer**

- **Lexer**: Ragel-based tokenization for unified bracket handling
- **Parser**: Context-aware parsing with scope tracking
- **AST Builder**: Constructs abstract syntax tree with proper scoping

#### 2. **Compiler Layer**

- **Template Compiler**: Converts AST to optimized Go code
- **Scope Analyzer**: Detects scoping issues and variable access
- **Cache Manager**: Implements template caching and hot-reloading

#### 3. **Runtime Layer**

- **Evaluator**: Executes compiled templates with proper scoping
- **Renderer**: Handles output generation and escaping
- **Helper System**: Provides built-in functions and filters

#### 4. **Security Layer**

- **Sandbox**: Restricts access to system resources
- **Escaper**: HTML escaping and sanitization
- **Limits**: Memory, CPU, and output size controls

#### 5. **Storage Layer**

- **Template Store**: Interface for custom storage implementations
- **Filesystem Store**: Default file-based template loading
- **Cache Store**: Template caching implementation

### Data Flow

```
Template File
    ↓
Lexer (Tokenization)
    ↓
Parser (AST Construction)
    ↓
Scope Analyzer (Scoping Validation)
    ↓
Compiler (Go Code Generation)
    ↓
Template Cache (Optional)
    ↓
Runtime Evaluator (Execution)
    ↓
Output (HTML)
```

## Implementation Details

### Scoping Strategy

#### Nested Scope Management

```go
type Scope struct {
    parent   *Scope  // Shared reference to parent scope (read-only)
    variables map[string]interface{}
    functions map[string]interface{}
    isRoot   bool
    isolated bool  // Security mode - no parent access
}
```

**Scope Rules:**

1. Each `{% tag %}` creates a new scope with shared parent reference
2. Parent scopes are accessible for reading (no copying needed)
3. Only current scope can be modified
4. Variable lookup walks up the scope chain until found
5. In isolated mode, only current scope is accessible (security restriction)

**Example:**

```liquid
{# Root scope #}
{% assign .rootVar = "hello" %}

{% for .i, .item in .items %}
  {# Loop scope - sees .rootVar, adds .i, .item #}
  {% .i %} - {% .item.name %}

  {% for .j, .subitem in .item.subitems %}
    {# Nested loop scope - sees all scopes #}
    {% .i %}.{% .j %} - {% .subitem.name %}
  {% endfor %}
{% endfor %}
```

### Template Loading

#### Import Mechanisms

**Include (Evaluate & Output):**

```go
// Evaluate template and output its contents
func (e *Engine) Include(name string, context map[string]interface{}) (string, error)
```

**Render (Partial Evaluation with parameters):**

```go
// Render template with parameters and output
func (e *Engine) Render(name string, params map[string]interface{}) (string, error)
```

#### Layout Inheritance

**Extends System:**

```go
// Single extends per template
func (e *Engine) RenderLayout(name string, content string, context map[string]interface{}) (string, error)
```

**Block Resolution:**

```go
// Block definition and override
type Block struct {
    name      string
    content   string
    isDefault bool
}
```

### Security Model

#### Sandbox Environment

**Restricted Operations:**

- Filesystem access: Blocked
- Network access: Blocked
- System calls: Blocked
- Dangerous functions: Restricted

**Resource Limits:**

```go
type Limits struct {
    MaxOutputSize  int64   // Maximum template output size in bytes
    MaxExecutionTime time.Duration // Maximum time to render a template
    MaxMemoryUsage int64   // Maximum memory usage during template rendering
    MaxIterations  int     // Maximum loop iterations to prevent infinite loops
    MaxDepth       int     // Maximum scope nesting depth
}
```

**Safety Features:**

- HTML auto-escaping by default
- XSS protection through context-aware escaping
- SQL injection prevention (when used in SQL contexts)
- Input validation and sanitization
- Isolated scope mode for untrusted data
- Template validation before execution

### Caching Strategy

#### Template Caching

**Cache Levels:**

1. **Source Cache**: Store parsed templates
2. **Compiled Cache**: Store generated Go code
3. **Runtime Cache**: Store compiled executable templates

**Cache Invalidation:**

```go
type CachePolicy struct {
    Enable      bool
    AutoRefresh bool
    TTL         time.Duration
    Version     string
}
```

## API Design

### Core Engine API

```go
type Engine struct {
    store          TemplateStore
    compiler       *Compiler
    evaluator      *Evaluator
    sandbox        *Sandbox
    cache          *Cache
    options        *Options
}

// Render template from string
func (e *Engine) RenderString(template string, context map[string]interface{}) (string, error)

// Render template from file
func (e *Engine) RenderFile(filename string, context map[string]interface{}) (string, error)

// Template compilation
func (e *Engine) Compile(template string) (*Template, error)

// Cache management
func (e *Engine) ClearCache() error
func (e *Engine) LoadTemplate(name string) error
```

### Context Management

```go
// Context with chain support
func (e *Engine) NewContext() *Context
func (c *Context) Set(name string, value interface{})
func (c *Context) SetChain(chain ...*Context) *Context
func (c *Context) Get(name string) (interface{}, bool)
func (c *Context) GetScope() *Scope
```

### Security Controls

```go
// Security options
func (e *Engine) EnableSandbox()
func (e *Engine) DisableSandbox()
func (e *Engine) SetLimits(limits *Limits)
func (e *Engine) EnableAutoEscape(enable bool)
```

### Template Loading

```go
// Template store interface
type TemplateStore interface {
    ReadTemplate(name string) ([]byte, error)
    ListTemplates() ([]string, error)
    WatchForChanges(callback func(string))
}

// Filesystem store (default)
func NewFileTemplateStore(directory string) TemplateStore
```

## Performance Considerations

### Optimization Strategies

**Compilation Benefits:**

- Templates compile to optimized Go code
- No interpretation overhead during rendering
- Minimal runtime reflection

**Caching Benefits:**

- Compiled templates cached in memory
- Hot-reload support for development
- Multiple compilation modes (dev, prod)

**Benchmark Goals (Realistic):**

- 300 templates/second rendering (achievable with caching)
- 2ms per template compilation (production mode)
- 0.05ms per variable access (with scope chain optimization)
- 95% cache hit rate in production

### Compilation Modes

```go
type CompilationMode int

const (
    Development CompilationMode = iota  // Fast compilation, minimal optimization
    Production                          // Optimized compilation, full caching
    Debug                               // Full validation, detailed error reporting
)
```

### Memory Management

**Efficient Scope Handling:**

```go
// Shared scope chain for parent access (read-only reference)
type ScopeChain struct {
    scopes []*Scope
}

// Pooling for frequently used scopes
var scopePool = sync.Pool{
    New: func() interface{} {
        return &Scope{
            variables: make(map[string]interface{}),
            functions: make(map[string]interface{}),
        }
    },
}
```

## Development Roadmap

### Phase 1: Core Infrastructure & Syntax (Weeks 1-3)

- [ ] Project structure and basic setup
- [ ] Lexer and parser implementation with unified brackets
- [ ] Scope management with shared parent references
- [ ] Variable resolution and member access
- [ ] Basic tag support (if, assign, unless)
- [ ] Syntax validation and error reporting

### Phase 2: Control Flow & Loops (Weeks 4-5)

- [ ] Loop structures (for, range, while)
- [ ] Case/when statements
- [ ] Unless statements
- [ ] Nested scope handling
- [ ] With context blocks
- [ ] Cycle, increment, decrement tags
- [ ] Break and continue in loops

### Phase 3: Template Composition (Weeks 6-7)

- [ ] Include mechanism (evaluate and output)
- [ ] Render mechanism (partial evaluation with params)
- [ ] Component system (reusable templates)
- [ ] Template caching (memory-based)
- [ ] Raw and comment blocks

### Phase 4: Layout System & Functions (Weeks 8-9)

- [ ] Layout inheritance (extends)
- [ ] Block definition and override
- [ ] Content block handling
- [ ] Multi-layout support
- [ ] Built-in filters (string, numeric, array)
- [ ] Helper functions
- [ ] Custom filter registration
- [ ] Function system

### Phase 5: Security & Performance (Weeks 10-11)

- [ ] Sandbox implementation (resource limiting)
- [ ] HTML auto-escaping by default
- [ ] XSS and SQL injection prevention
- [ ] Template validation before execution
- [ ] Bytecode caching
- [ ] Performance monitoring and metrics
- [ ] Compilation modes (dev, prod, debug)

### Phase 6: Testing & Tooling (Weeks 12-13)

- [ ] Comprehensive test suite
- [ ] Performance benchmarks
- [ ] Fuzz testing for security
- [ ] Cross-engine compatibility tests
- [ ] CLI tool development (wisp command)
- [ ] Template validation CLI

### Phase 7: Documentation & Examples (Weeks 14-15)

- [ ] API documentation
- [ ] Getting started guide
- [ ] Template syntax reference
- [ ] Security best practices
- [ ] Examples and tutorials
- [ ] Integration guides for popular frameworks

## Testing Strategy

### Test Categories

**Unit Tests:**

- Parser functionality (lexer, tokenization, AST generation)
- Scope management (creation, lookup, isolation, chaining)
- Variable resolution (member access, indexing, chaining)
- Tag processing (control flow, loops, variable manipulation)
- Function and filter execution (built-in and custom)
- Template compilation (AST to Go code generation)
- Security features (sandboxing, auto-escaping, validation)

**Integration Tests:**

- Full template rendering (end-to-end execution)
- Layout inheritance (extends, blocks, content)
- Template composition (include, render, component)
- Nested structures (loops, conditionals, with blocks)
- Data binding (structs, maps, slices, primitives)
- Helper functions (string, numeric, array, date, HTML)

**Performance Tests:**

- Rendering speed (templates per second)
- Compilation overhead (time per template)
- Memory usage (allocation and garbage collection)
- Concurrency handling (parallel template execution)
- Cache efficiency (hit rates, invalidation)
- Scope chain performance (variable lookup depth)

**Security Tests:**

- XSS protection (script injection prevention)
- SQL injection prevention (in SQL contexts)
- Input validation (malformed data handling)
- Resource exhaustion (memory, CPU, output limits)
- Sandbox enforcement (filesystem, network access blocking)
- Template validation (malicious template detection)

### Test Examples

```liquid
{# Test: Nested loops with scope isolation #}
{% for .i, .item in .items %}
  {# .i and .item available here #}
  {% for .j, .subitem in .item.subitems %}
    {# .i, .item, .j, .subitem all available #}
    {% .i %}.{% .j %}: {% .subitem.name %}
  {% endfor %}
  {# .j and .subitem not available here #}
{% endfor %}

{# Test: With context block #}
{% with .user as .currentUser %}
  {% .currentUser.name %} {# Accessible #}
{% endwith %}
{# .currentUser not available here #}

{# Test: Cycle tag #}
{% cycle .items 'odd' 'even' %}
{# Should output: odd, even, odd, even, ... #}

{# Test: Function call with pipe syntax #}
{% .createdAt | date "2006-01-02" %}

{# Test: Chained filters #}
{% .price | format "%.2f" | currency "USD" %}

{# Test: Raw block #}
{% raw %}
  This {{ .variable }} will not be processed
{% endraw %}

{# Test: Comment block #}
{% comment %}
  This entire block is a comment
{% endcomment %}

{# Test: Break and continue #}
{% for .item in .items %}
  {% if .item.skip %}{% continue %}{% end %}
  {% if .item.last %}{% break %}{% end %}
  {% .item.name %}
{% end %}

{# Test: Assign and capture #}
{% assign .count = 0 %}
{% capture .output %}
  Count is {{ .count }}
{% endcapture %}
{% .output %} {# Should output: "Count is 0" #}

{# Test: Include and render #}
{% include "header" .user %}
{% render "widget" .data prop1 .value prop2 .other %}
```

## File Structure

```
wisp/
├── cmd/
│   └── wisp/
│       ├── main.go
│       ├── test.go
│       └── benchmark.go
├── internal/
│   ├── lexer/
│   ├── parser/
│   ├── compiler/
│   ├── evaluator/
│   ├── scope/
│   ├── tags/
│   ├── filters/
│   ├── security/
│   └── store/
├── pkg/
│   ├── engine/
│   ├── context/
│   └── types/
├── templates/
│   ├── layouts/
│   ├── partials/
│   └── pages/
├── tests/
│   ├── fixtures/
│   └── benchmarks/
├── examples/
├── docs/
└── scripts/
```

## Future Enhancements

### Planned Features

- Custom tag registration
- Plugin system
- Template inheritance with partial blocks
- Debug mode with line numbers
- Template validation

## References

- Shopify Liquid: <https://shopify.github.io/liquid>
- Go Template System: <https://pkg.go.dev/text/template>
- Liquid Design Philosophy: Security, Simplicity, Performance

