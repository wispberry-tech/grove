# API Documentation

## Engine

The `Engine` is the main entry point for template rendering.

### Constructor

```go
func New() *Engine
```

Creates a new Engine with default settings:
- Auto-escaping enabled
- Max iterations: 100000

```go
e := engine.New()
```

### Unsafe Engine

```go
func NewUnsafe() *Engine
```

Creates an engine with auto-escaping disabled.

```go
e := engine.NewUnsafe()
```

---

## Rendering

### RenderString

```go
func (e *Engine) RenderString(template string, data map[string]interface{}) (string, error)
```

Renders a template string with the given data context.

**Parameters:**
- `template`: Template string with Wisp syntax
- `data`: Map of variables available in the template

**Returns:**
- Rendered string
- Error if parsing or rendering fails

**Example:**

```go
result, err := e.RenderString(`Hello, {%.name%}!`, map[string]interface{}{
    "name": "World",
})
// result: "Hello, World!"
```

### RenderFile

```go
func (e *Engine) RenderFile(filename string, data map[string]interface{}) (string, error)
```

Renders a template from the registered template store.

**Parameters:**
- `filename`: Template name registered in the store
- `data`: Map of variables available in the template

**Returns:**
- Rendered string
- Error if template not found or rendering fails

**Example:**

```go
e.RegisterTemplate("greeting", `Hello, {%.name%}!`)
result, err := e.RenderFile("greeting", map[string]interface{}{
    "name": "World",
})
```

### Validate

```go
func (e *Engine) Validate(template string) error
```

Validates template syntax without rendering.

**Parameters:**
- `template`: Template string to validate

**Returns:**
- Error if template has syntax errors
- nil if valid

**Example:**

```go
err := e.Validate(`{% if .show %}{%.content%}{% end %}`)
if err != nil {
    // Handle syntax error
}
```

---

## Template Store

### SetStore

```go
func (e *Engine) SetStore(store TemplateStore)
```

Sets the template store for file-based template loading.

**Example:**

```go
e.SetStore(engine.NewFileStore("./templates"))
```

### RegisterTemplate

```go
func (e *Engine) RegisterTemplate(name, content string)
```

Registers a template for rendering by name.

**Parameters:**
- `name`: Template identifier
- `content`: Template content

**Example:**

```go
e.RegisterTemplate("header", `<header>{%.title%}</header>`)
e.RegisterTemplate("footer", `<footer>Copyright 2024</footer>`)
```

### ClearCache

```go
func (e *Engine) ClearCache()
```

Clears the parsed template cache. Call this after registering new templates or updating existing ones.

---

## Filters

### RegisterFilter

```go
func (e *Engine) RegisterFilter(name string, fn interface{})
```

Registers a custom filter function.

**Parameters:**
- `name`: Filter name used in templates (e.g., `"shout"` for `| shout`)
- `fn`: Filter function

**Filter Function Signatures:**

Single argument:
```go
func(input interface{}) interface{}
```

With arguments:
```go
func(input interface{}, args ...interface{}) interface{}
```

**Examples:**

```go
// Simple filter
e.RegisterFilter("shout", func(s interface{}) string {
    return fmt.Sprintf("!!! %s !!!", toString(s))
})

// Filter with arguments
e.RegisterFilter("wrap", func(s interface{}, before, after interface{}) string {
    return fmt.Sprintf("%s%s%s", toString(before), toString(s), toString(after))
})
```

Usage in templates:
```liquid
{% .message | shout %}
{% .text | wrap "<<" ">>" %}
```

---

## Security

### SetAutoEscape

```go
func (e *Engine) SetAutoEscape(enabled bool)
```

Enables or disables HTML auto-escaping.

**Parameters:**
- `enabled`: true for auto-escaping (default), false to disable

**Example:**

```go
e.SetAutoEscape(false)  // Disable escaping for HTML templates
```

### SetMaxIterations

```go
func (e *Engine) SetMaxIterations(max int)
```

Sets the maximum number of loop iterations to prevent infinite loops.

**Parameters:**
- `max`: Maximum iterations allowed

**Example:**

```go
e.SetMaxIterations(10000)
```

---

## SafeString

```go
type SafeString string
```

A string that will not be HTML-escaped when rendered.

### Create SafeString

```go
safe := engine.SafeString("<b>Bold</b>")
```

### Usage

```go
result, _ := e.RenderString(`{%.html%}`, map[string]interface{}{
    "html": engine.SafeString("<b>Bold</b>"),
})
// result: "<b>Bold</b>" (not escaped)
```

---

## TemplateStore Interface

```go
type TemplateStore interface {
    ReadTemplate(name string) (string, error)
}
```

Interface for custom template storage implementations.

### Implementations

#### MemoryStore

```go
func NewMemoryStore() *MemoryStore
```

Creates an in-memory template store.

```go
ms := engine.NewMemoryStore()
ms.Register("header", `<header>...</header>`)
e.SetStore(ms)
```

#### FileStore

```go
func NewFileStore(directory string) *FileStore
```

Creates a file system-based template store.

```go
fs := engine.NewFileStore("./templates")
e.SetStore(fs)
```

---

## Error Handling

Wisp returns errors for:

- **Parse errors**: Invalid template syntax
- **Runtime errors**: Missing variables, invalid operations
- **IO errors**: Template file not found

### Parse Errors

```go
result, err := e.RenderString(template, data)
if err != nil {
    if errs, ok := err.([]error); ok {
        for _, e := range errs {
            fmt.Println("Parse error:", e)
        }
    }
}
```

---

## CLI Commands

The `wisp` CLI tool provides:

### render

```bash
wisp render <template> [data]
```

Renders a template with JSON data.

```bash
# From argument
wisp render 'Hello, {%.name%}!' '{"name": "World"}'

# From stdin
echo '{"name": "World"}' | wisp render 'Hello, {%.name%}!'
```

### validate

```bash
wisp validate <template>
```

Validates template syntax.

```bash
wisp validate '{% if .show %}{%.content%}{% end %}'
```

### version

```bash
wisp version
```

Shows version information.
