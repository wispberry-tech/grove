# AGENTS.md - Wisp Template Engine

This file provides guidelines for agentic coding agents working on the Wisp template engine codebase.

## Build and Test Commands

### Essential Commands
- **Build CLI**: `go build -o wisp ./cmd/wisp`
- **Run all tests**: `go test ./... -v`
- **Run single test**: `go test ./pkg/engine -v -run TestRenderString`
- **Run single test file**: `go test ./internal/lexer -v`
- **Run benchmarks**: `go test ./pkg/engine -bench=.`
- **Run example server**: `go run ./examples/server`
- **Format code**: `go fmt ./...`
- **Vet (lint)**: `go vet ./...`

### Test Package Selection
- `internal/lexer` - Tokenization logic (52 token types, 35 keywords)
- `internal/parser` - AST parsing (Pratt parser, 35 node types)
- `internal/ast` - AST node definitions
- `internal/resolver` - Variable resolution (member access, indexing)
- `internal/scope` - Variable scoping (sync.Pool for reuse)
- `internal/evaluator` - Template evaluation (control flow, filters)
- `internal/store` - Template storage (MemoryStore, FileStore)
- `pkg/engine` - Public API and integration tests

## Code Organization

### Directory Structure
```
/cmd/              - CLI tools (wisp, test_with)
/examples/         - Example applications (server)
/internal/         - Private packages (not exported externally)
  ast/            - AST node definitions
  evaluator/      - Template evaluation engine
  lexer/          - Tokenization
  parser/         - AST parsing
  resolver/       - Variable resolution
  scope/          - Variable scoping with sync.Pool
  store/          - Template storage
/pkg/             - Public API (exported)
  engine/         - Main Engine, filters, public interface
  registry/      - Tags, helpers, filters registry
/templates/       - Example template files
```

### Package Responsibilities
- **internal/lexer**: Tokenizes template strings into tokens. Uses context-aware tokenization to distinguish TEXT from Wisp `{% %}` blocks.
- **internal/parser**: Parses tokens into AST nodes using Pratt parsing. Handles all control flow, loops, filters, and template composition.
- **internal/resolver**: Resolves variable references in scopes. Supports dot notation, array indexing, and map key access.
- **internal/scope**: Manages nested variable scopes with sync.Pool for performance. Chain-based lookup with parent scope access.
- **internal/evaluator**: Evaluates AST nodes against data. Implements all template features: conditionals, loops, filters, includes, layouts.
- **internal/ast**: Defines all AST node types (Statement, Expression, etc.). 35 node types covering all constructs.
- **pkg/engine**: Public API (New, RenderString, RenderFile, SetX methods). Provides caching, auto-escaping, resource limits.
- **pkg/registry**: Built-in tags, filters, helpers registration.

## Code Style Guidelines

### Package Structure
- `internal/` - Private packages (not exported)
- `pkg/` - Public API (exported)
- `cmd/` - CLI tools
- `examples/` - Example applications

### Imports
Order: stdlib -> third-party -> internal/pkg, grouped with blank lines:

```go
import (
    "fmt"
    "os"
    "sync"

    "template-wisp/internal/ast"
    "template-wisp/internal/lexer"
)
```

### Naming Conventions
- **Packages**: lowercase, single word preferred (lexer, parser, scope, engine)
- **Exported functions**: PascalCase (NewLexer, RenderString, ParseProgram)
- **Internal functions**: camelCase (readChar, peekChar, skipWhitespace)
- **Struct fields**: camelCase (autoEscape, maxIter, cache)
- **Constants**: PascalCase for internal (signalNone), UPPER_CASE for exported (EOF)
- **Interfaces**: Simple nouns (Node, Statement, Expression, TemplateStore)

### Formatting
- Use standard Go formatting (`gofmt`)
- Tab indentation
- No trailing whitespace
- Blank lines between functions and logical sections

### Types
- Named types use PascalCase: Lexer, Parser, Engine, Scope
- Exported struct fields: PascalCase
- Unexported struct fields: camelCase
- Interface methods: PascalCase for exported

### Error Handling
- Return `(type, error)` for operations that can fail
- Use `fmt.Errorf` with descriptive messages:
  ```go
  return "", fmt.Errorf("parse errors: %v", p.Errors())
  ```
- Check errors immediately: `if err != nil { ... }`
- Return nil error on success
- Use wrapped errors with %w for context:
  ```go
  return "", fmt.Errorf("failed to read template %s: %w", filename, err)
  ```
- Return slices of errors for multiple validation errors:
  ```go
  if len(p.Errors()) > 0 {
      return fmt.Errorf("parse errors: %v", p.Errors())
  }
  ```

### Testing Conventions
- Test functions: `TestNameOfFunction`
- Use table-driven tests for multiple cases:
  ```go
  tests := []struct {
      name     string
      template string
      expected string
  }{ ... }
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) { ... })
  }
  ```
- Use `t.Fatalf()` for setup failures
- Use `t.Errorf()` for assertion failures
- Use testify/assert when appropriate

### Common Patterns
- **Constructors**: NewTypeName (NewLexer, NewParser, New)
- **Getters**: GetX (Get, GetFunction)
- **Setters**: SetX (SetAutoEscape, SetMaxIterations)
- **Validation**: Validate method returns error
- **Interfaces**: Define before implementations
- **Pooling**: Use sync.Pool for frequently created objects (scopes)
- **Dependency Injection**: Pass dependencies via constructors, not globals

### Public API Design
- Use `pkg/engine` for exported types and functions
- Provide both safe (New) and unsafe (NewUnsafe) constructors when applicable
- Cache parsed templates for performance
- Default to secure settings (auto-escape enabled)
- Provide configuration methods (SetX)
- Document all exported functions with package-level comments

## Security Considerations
- HTML auto-escaping enabled by default (XSS protection)
- Use SafeString type to bypass escaping when needed
- Limit loop iterations to prevent DoS via SetMaxIterations
- Track include depth to detect circular references
- Scope isolation for render/component templates

## Performance Guidelines
- Use sync.Pool for frequently allocated objects (scopes)
- Cache parsed templates (AST cached per template string)
- Use string.Builder for concatenation in hot paths
- Minimize allocations in hot paths (evaluator, resolver)
- Use RWMutex for cached template access (concurrent reads)

## Documentation Standards
- Exported functions need package-level comments explaining:
  - What the function does
  - Parameters and return values
  - Error conditions
  - Security implications if any
- Document reasoning for design choices
- Keep docs/getting-started.md, docs/syntax-reference.md, docs/api.md in sync
- Add examples showing typical usage patterns

## Architecture Overview

### Template Rendering Pipeline
1. **Lexical Analysis**: template string → tokens (lexer)
2. **Parsing**: tokens → AST (parser with Pratt parsing)
3. **Evaluation**: AST + data → rendered output (evaluator)

### Key Design Decisions
- **Unified bracket syntax**: `{% %}` for all template logic vs `{{ }}` for output
- **Auto-escaping default**: Security-first approach prevents XSS
- **Scope chaining**: Parent scope access enables nested context
- **Pratt parsing**: Handles operator precedence naturally
- **sync.Pool scoping**: Reduces allocation pressure in tight loops
