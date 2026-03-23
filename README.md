# Wisp Template Engine

A secure, high-performance HTML templating engine for Go inspired by Liquid with simplified syntax.

## Features

- **Simplified Syntax**: Unified bracket syntax with clear variable/function distinction
- **Security First**: HTML auto-escaping by default, XSS protection
- **High Performance**: Template caching, scope pooling, optimized evaluation
- **Template Composition**: Include, render, and component systems
- **Layout System**: Template inheritance with extends/block/content
- **40 Built-in Filters**: String, numeric, array, date, URL, and math operations
- **Control Flow**: if/elsif/else, unless, case/when, for/while/range loops
- **Context Management**: Nested scopes with shared parent references

## Quick Start

### Installation

```bash
go get template-wisp
```

### Basic Usage

```go
package main

import (
    "fmt"
    "template-wisp/pkg/engine"
)

func main() {
    e := engine.New()
    
    template := `{% .name %} is {% .age %} years old.`
    data := map[string]interface{}{
        "name": "Alice",
        "age":  30,
    }
    
    result, err := e.RenderString(template, data)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result)
    // Output: Alice is 30 years old.
}
```

### Using Filters

```go
template := `{% .title | upcase %}`
data := map[string]interface{}{
    "title": "hello world",
}

result, _ := e.RenderString(template, data)
// Output: HELLO WORLD
```

### Control Flow

```go
template := `
{% for .item in .items %}
  - {% .item.name %}: {% .item.price | currency %}
{% end %}
`

data := map[string]interface{}{
    "items": []map[string]interface{}{
        {"name": "Apple", "price": 1.99},
        {"name": "Orange", "price": 0.99},
    },
}
```

## Documentation

- **[API Documentation](docs/api.md)** - Complete API reference
- **[Getting Started](docs/getting-started.md)** - Detailed usage guide
- **[Syntax Reference](docs/syntax-reference.md)** - Full template syntax
- **[Design Plan](plan.md)** - Architecture and design decisions
- **[Development Status](development.md)** - Current implementation status

## Examples

See the [examples/](examples/) directory for working examples:

- **HTTP Server**: Interactive web server with live template rendering
  ```bash
  go run ./examples/server
  # Visit http://localhost:8080
  ```

## CLI Tool

The `wisp` CLI tool provides command-line template rendering:

```bash
# Build the CLI
go build -o wisp ./cmd/wisp

# Render a template with JSON data
echo '{"name":"Alice"}' | wisp render "{% .name %}"

# Validate template syntax
wisp validate "{% if .condition %}{% .value %}{% end %}"

# Show version
wisp version
```

## Built-in Filters

### String Filters
- `capitalize`, `upcase`, `downcase`, `truncate`, `strip`, `lstrip`, `rstrip`
- `replace`, `remove`, `split`, `join`, `prepend`, `append`

### Numeric Filters
- `abs`, `ceil`, `floor`, `round`
- `plus`, `minus`, `times`, `divided_by`, `modulo`

### Array Filters
- `first`, `last`, `size`, `length`, `reverse`, `sort`, `uniq`, `map_field`

### General Filters
- `default`, `json`, `raw`, `escape`, `escape_once`

### Date Filters
- `date`, `date_format`

### URL Filters
- `url_encode`, `url_decode`

### Math Filters
- `min`, `max`

## Security

### Auto-Escaping

HTML auto-escaping is enabled by default to prevent XSS attacks:

```go
e := engine.New()  // Auto-escape enabled
e.SetAutoEscape(false)  // Disable if needed
```

### Safe Strings

Use the `raw` filter to output unescaped HTML:

```html
{% .html_content | raw %}
```

### Resource Limits

Configure limits to prevent resource exhaustion:

```go
e := engine.New()
e.SetMaxIterations(10000)  // Max loop iterations
```

## Testing

Run the test suite:

```bash
go test ./... -v
```

Run benchmarks:

```bash
go test ./pkg/engine -bench=.
```

## Contributing

Contributions are welcome! Please read the design documentation and development status before submitting PRs.

## License

[See LICENSE file](LICENSE)

## Development Status

The Wisp template engine is in active development with core functionality complete.

- ✅ Core infrastructure (lexer, parser, AST, scope, resolver)
- ✅ All control flow and loops
- ✅ Template composition (include, render, component)
- ✅ Layout system (extends, block, content)
- ✅ 40 built-in filters
- ✅ Security features (auto-escaping, resource limiting)
- ✅ Performance optimizations (caching, pooling)
- ✅ CLI tool and example server

See [development.md](development.md) for detailed status.
