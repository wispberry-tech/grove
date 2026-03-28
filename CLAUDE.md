# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=profile.cov ./...

# Run tests for a specific package
go test ./internal/evaluator/...

# Run benchmarks
go test -bench=. ./pkg/engine/...

# Run a single test by name
go test -run TestName ./pkg/engine/...

# Start the documentation/testing server
go run ./docs/server/main.go
```

## Architecture

**Template-Wisp** is a Liquid-inspired HTML templating engine for Go with a classic compiler pipeline:

```
Text Input → Lexer → Parser → AST → Evaluator → String Output
```

### Pipeline stages

| Package | Location | Responsibility |
|---|---|---|
| `lexer` | `internal/lexer/` | Converts raw template text into tokens (`{% %}`, `{# #}` delimiters) |
| `parser` | `internal/parser/` | Builds an AST (Program, Statement, Expression nodes) from the token stream |
| `ast` | `internal/ast/node.go` | Defines all AST node types (IfStatement, ForStatement, AssignStatement, etc.) |
| `evaluator` | `internal/evaluator/` | Walks the AST, produces output; owns the `SafeString` type and auto-escape logic |
| `scope` | `internal/scope/` | Manages variable scoping; uses `sync.Pool` for reuse; supports parent-chain lookup |
| `resolver` | `internal/resolver/` | Resolves variable/function names within a scope |
| `store` | `internal/store/store.go` | Pluggable template storage: `MemoryStore`, `FileSystemStore` |
| `engine` | `pkg/engine/` | Public API — wraps the pipeline, owns the parsed-AST cache and filter registry |

### Security model

HTML auto-escaping is **on by default**. The `SafeString` type in the evaluator marks trusted HTML that should not be escaped. Use `NewUnsafe` only when rendering non-HTML output. The `render` tag executes sub-templates in an isolated scope to prevent variable leakage.

### Template composition

- `include` — includes another template, sharing the current scope
- `render` — includes another template with an isolated scope (sandboxed)
- `component` — includes a template passing explicit props
- `extends` / `block` / `endblock` — layout inheritance

### Public API entry points

- `pkg/engine/engine.go` — `New()`, `NewUnsafe()`, `NewWithStore()`, `RenderString()`, `RenderFile()`
- `pkg/engine/filters.go` — 40+ built-in filters; also shows how to register custom filters

### Documentation server

`docs/server/main.go` doubles as an interactive testing server. Key routes:

- `GET /api/tokens` — tokenize a template and return the token stream
- `GET /api/ast` — parse and return the AST as JSON
- `GET /api/render` — render a template with provided data
- `GET /docs/` — serve documentation pages
