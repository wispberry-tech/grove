# Grove

A bytecode-compiled template engine for Go with components, inheritance, and web primitives.

## Install

```bash
go get grove
```

## Quick Example

```go
package main

import (
	"context"
	"fmt"
	"grove/pkg/grove"
)

func main() {
	eng := grove.New()
	result, err := eng.RenderTemplate(
		context.Background(),
		`Hello, {{ name | title }}!`,
		grove.Data{"name": "world"},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Body) // Hello, World!
}
```

## Features

- **Bytecode compilation** — templates compile to bytecode and run on a stack-based VM. Compiled bytecode is immutable and shared across goroutines.
- **Template inheritance** — `extends`, `block`, and `super()` for layered layouts.
- **Components** — reusable templates with `props`, `slot`, and `fill`. Fills see the caller's scope, not the component's.
- **Macros** — `macro`, `call`, `caller()`, and `import` for reusable template functions.
- **40+ built-in filters** — string, collection, numeric, HTML, and type conversion filters with pipe syntax.
- **Web primitives** — `asset`, `meta`, and `hoist` tags collect resources during rendering. `RenderResult` returns them for assembly into the final HTML response.
- **Auto-escaping** — HTML output is escaped by default. Use the `safe` filter or `raw` blocks for trusted content.
- **Sandboxing** — restrict allowed tags, filters, and loop iterations per engine.
- **List and map literals** — `[1, 2, 3]` and `{key: "value"}` inline data structures.
- **Ternary expressions** — `condition ? truthy : falsy` with right-to-left nesting.

## Documentation

Full documentation is in the [`docs/`](docs/index.md) directory:

- [Getting Started](docs/getting-started.md) — install, configure, render your first template
- [Template Syntax](docs/template-syntax.md) — variables, expressions, control flow, assignment
- [Template Inheritance](docs/template-inheritance.md) — extends, block, super()
- [Components](docs/components.md) — props, slots, fills
- [Macros & Includes](docs/macros-and-includes.md) — macro, include, render, import
- [Filters](docs/filters.md) — all 42 built-in filters
- [Web Primitives](docs/web-primitives.md) — asset, meta, hoist, RenderResult
- [API Reference](docs/api-reference.md) — Go types, methods, and configuration
- [Examples](docs/examples.md) — blog app walkthrough
