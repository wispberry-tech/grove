# Wisp Template Engine - Development Plan

This file tracks the detailed development progress for the Wisp template engine based on the design specifications in plan.md.

## Development Status Overview

- **Current Phase**: Development Complete (Optional Items Remaining)
- **Overall Progress**: ~97% complete
- **Last Updated**: 2026-03-23

### Test Status Summary
```
✅ internal/lexer     - PASS (2/2 tests)
✅ internal/parser    - PASS (8/8 tests)
✅ internal/resolver  - PASS (9/9 tests)
✅ internal/scope     - PASS (11/11 tests)
✅ internal/evaluator - PASS (23/23 tests)
✅ pkg/engine         - PASS (18/18 tests)
Total: 71/71 tests passing (100%)
```

## Current Implementation Status

### ✅ Completed Components

#### Lexer (100%) | Parser (100%) | AST (100%) | Scope (100%) | Resolver (100%)
All core infrastructure complete and fully tested.

- **Lexer**: 52 token types, 35 keywords, context-aware tokenization (TEXT vs Wisp tokens)
- **Parser**: Pratt parser (1495 lines), handles all Wisp `{% %}` syntax
- **AST**: 35 node types covering all template constructs
- **Scope**: Chain-based lookup with pooling, isolation support
- **Resolver**: Variable resolution, member access, array/map indexing, type coercion, operator application

#### Evaluator (100%)
- All control flow: if/elsif/else, unless, case/when
- All loops: for, while, range with break/continue
- Context blocks: with, cycle, increment, decrement
- Template composition: include (shared scope), render (isolated scope), component (props-based)
- Layout inheritance: extends/block/content with overrides
- Raw and comment blocks with content capture
- **HTML auto-escaping** (enabled by default, SafeString bypass)
- **Resource limiting** (max loop iterations for while/range)

#### Public API (100%)
- Engine with `RenderString`/`RenderFile`/`Validate`
- TemplateStore interface, MemoryStore, FileSystemStore
- `RegisterFilter`, `RegisterTemplate`, `ClearCache`
- **Template caching** (parsed AST cached per template string)
- **Auto-escape control** (`SetAutoEscape`, `NewUnsafe`)
- **Resource limits** (`SetMaxIterations`)

#### Built-in Filters (100%) - 40 total
- String (13): capitalize, upcase, downcase, truncate, strip, lstrip, rstrip, replace, remove, split, join, prepend, append
- Numeric (9): abs, ceil, floor, round, plus, minus, times, divided_by, modulo
- Array (8): first, last, size, length, reverse, sort, uniq, map_field
- General (5): default, json, raw, escape, escape_once
- Date (2): date, date_format
- URL (2): url_encode, url_decode
- Math (2): min, max

#### Security Features (90%)
- **HTML auto-escaping** by default (XSS protection via `html.EscapeString`)
- **SafeString** type for bypassing escaping
- **raw/escape/escape_once** filters
- **Max iteration limiting** (prevents infinite loop DoS)
- **Scope isolation** for render/component templates
- Missing: Full sandbox (filesystem/network blocking), output size limits

#### Performance Features (100%)
- **Parsed template caching** (AST cached per template string)
- **Scope pooling** (`sync.Pool` for scope reuse)

#### Example HTTP Server (100%)
- Interactive web UI at http://localhost:8080
- Routes: /, /lexer, /parser, /render
- API: /api/tokens, /api/ast, /api/render
- Colored logging with charmbracelet/log
- Live pipeline visualization (tokens + AST + rendered output)

#### CLI Tool (100%)
- `wisp render <template> [data]` - Render template with JSON data (supports stdin)
- `wisp validate <template>` - Validate template syntax
- `wisp version` - Print version
- `wisp help` - Usage information

#### Circular Include Detection (100%)
- Tracks include chain depth during template resolution
- Prevents infinite recursion from circular includes

### ❌ Not Started
- Fuzz testing
- Framework integration guides

### ✅ Completed Documentation
- Getting started guide (`docs/getting-started.md`)
- Template syntax reference (`docs/syntax-reference.md`)
- API documentation (`docs/api.md`)

### ✅ Benchmarks
- Performance benchmarks (`pkg/engine/benchmark_test.go`)
- Simple, conditional, loop, nested, filter, and complex template benchmarks

## Phase Completion

### Phase 1: Core Infrastructure ✅
- [x] Lexer with unified bracket tokenization
- [x] Parser with Pratt parsing for all Wisp syntax
- [x] AST with 35 node types
- [x] Scope management with chain lookup and pooling
- [x] Variable resolution with member/index access

### Phase 2: Control Flow ✅
- [x] If/elsif/else conditionals
- [x] Unless (negated if)
- [x] Case/when statements
- [x] For loops with index/value binding
- [x] While loops with iteration limits
- [x] Range loops
- [x] Break and continue
- [x] With context blocks
- [x] Cycle, increment, decrement

### Phase 3: Template Composition ✅
- [x] Include mechanism (shared scope)
- [x] Render mechanism (isolated scope)
- [x] Component system (props-based)
- [x] Template caching
- [x] Raw and comment blocks

### Phase 4: Layout System & Functions ✅
- [x] Layout inheritance with extends/block/content
- [x] 40 built-in filters
- [x] Helper functions (date, URL, math)
- [x] Custom filter registration

### Phase 5: Security & Performance ✅
- [x] HTML auto-escaping (XSS protection)
- [x] Resource limiting (max loop iterations)
- [x] Template caching (parsed AST)
- [x] Scope isolation for untrusted templates
- [x] Circular include detection
- [ ] Full sandbox (filesystem/network blocking) - deferred

### Phase 6: Testing & Tooling ✅
- [x] 71 unit tests, all passing
- [x] Performance benchmarks (12 benchmark tests)
- [x] CLI tool (render, validate, version, help)
- [x] Example HTTP server with interactive UI
- [ ] Fuzz testing - not yet

### Phase 7: Documentation ✅
- [x] Example HTTP server with interactive UI
- [x] Circular include detection
- [x] API documentation
- [x] Getting started guide
- [x] Template syntax reference
- [ ] Security best practices (optional)
- [ ] Framework integration guides (optional)

## Milestones

1. Foundation ✅
2. Control Flow ✅
3. Template Composition ✅
4. Layout System ✅
5. Security & Performance ✅
6. Testing & Tooling ✅
7. Documentation ✅

## Optional Remaining Items
- Fuzz testing
- Framework integration guides (Rails, Gin, Echo, Fiber)
- Security best practices guide

## Dependencies

- Go 1.26+
- github.com/stretchr/testify (testing)
- github.com/charmbracelet/log (colored logging, indirect)

---

*Last updated: 2026-03-23*
