# Grove Documentation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create complete project documentation — README, syntax reference, API reference, filter catalog, and example walkthrough — as Markdown files viewable on GitHub.

**Architecture:** 11 Markdown files: a root README.md and 10 docs pages in `docs/`. Each page is self-contained with cross-links. No build step, no static site generator. Content is derived from the actual codebase — public API in `pkg/grove/`, filter implementations in `internal/filters/`, test patterns in `pkg/grove/*_test.go`, and the blog example in `examples/blog/`.

**Tech Stack:** Markdown (GitHub-flavored). Go code examples use `go` fenced blocks. Template examples use `jinja2` fenced blocks (closest GitHub syntax highlighting for `.grov` files).

**Important codebase notes:**
- Go module: `grove` — import path for the public package is `grove/pkg/grove`
- The `Resolvable` interface method is `WispyResolve` (legacy name from before project rename)
- The blog example still references "Wispy" in some places — docs should use "Grove" consistently
- Template files use `.grov` extension but the blog example store references them as `.html` — this is a store naming convention, not a file extension requirement

---

## File Map

| File | Action | Content |
|------|--------|---------|
| `README.md` | Create | Project landing page — tagline, install, quick example, feature bullets, link to docs |
| `docs/index.md` | Create | Docs home — overview paragraph, linked TOC |
| `docs/getting-started.md` | Create | Go integration guide — install, engine setup, stores, data, options, errors |
| `docs/template-syntax.md` | Create | Complete syntax reference — delimiters, variables, expressions, control flow, assignment |
| `docs/template-inheritance.md` | Create | Inheritance system — extends, block, super(), multi-level chains |
| `docs/components.md` | Create | Component model — component, props, slots, fills, scope rules |
| `docs/macros-and-includes.md` | Create | Composition — include, render, macro, call, caller(), import |
| `docs/filters.md` | Create | Complete filter catalog — all 42 built-in filters with signatures and examples |
| `docs/web-primitives.md` | Create | Web features — asset, meta, hoist, RenderResult, auto-escaping |
| `docs/api-reference.md` | Create | Go API — Engine, options, stores, filters, values, errors |
| `docs/examples.md` | Create | Blog example walkthrough — structure, templates, Go code, running it |

---

## Task 1: README.md

**Files:**
- Create: `README.md`

- [ ] **Step 1: Write README.md**

```markdown
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
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add project README"
```

---

## Task 2: docs/index.md

**Files:**
- Create: `docs/index.md`

- [ ] **Step 1: Write docs/index.md**

```markdown
# Grove Documentation

Grove is a bytecode-compiled template engine for Go. Templates are lexed, parsed into an AST, compiled to bytecode, and executed on a stack-based VM. The engine is safe for concurrent use — compiled bytecode is immutable and shared across goroutines, and VM instances are pooled.

## Contents

| Page | Description |
|------|-------------|
| [Getting Started](getting-started.md) | Install Grove, configure an engine, render your first template |
| [Template Syntax](template-syntax.md) | Variables, expressions, operators, control flow, loops, assignment, literals |
| [Template Inheritance](template-inheritance.md) | Base layouts with `extends`, `block`, and `super()` |
| [Components](components.md) | Reusable templates with `props`, `slot`, and `fill` |
| [Macros & Includes](macros-and-includes.md) | Template functions with `macro`, and composition with `include`, `render`, `import` |
| [Filters](filters.md) | All 42 built-in filters — string, collection, numeric, HTML, type conversion |
| [Web Primitives](web-primitives.md) | `asset`, `meta`, `hoist` tags and `RenderResult` integration |
| [API Reference](api-reference.md) | Go types, methods, options, stores, custom filters, error types |
| [Examples](examples.md) | Walkthrough of the blog example app |
```

- [ ] **Step 2: Commit**

```bash
git add docs/index.md
git commit -m "docs: add documentation index"
```

---

## Task 3: docs/getting-started.md

**Files:**
- Create: `docs/getting-started.md`

- [ ] **Step 1: Write docs/getting-started.md**

```markdown
# Getting Started

## Installation

```bash
go get grove
```

Import the package:

```go
import "grove/pkg/grove"
```

## Rendering an Inline Template

The simplest way to use Grove — create an engine and render a template string directly:

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
		`Hello, {{ name }}! You have {{ count }} messages.`,
		grove.Data{
			"name":  "Alice",
			"count": 3,
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Body)
	// Output: Hello, Alice! You have 3 messages.
}
```

`grove.Data` is an alias for `map[string]any`. Pass any Go values — strings, numbers, booleans, slices, maps, or custom types.

`RenderTemplate` returns a `RenderResult`. The `Body` field contains the rendered output. Other fields (`Assets`, `Meta`, `Hoisted`, `Warnings`) are used by [web primitives](web-primitives.md).

## File-Based Templates

For real applications, store templates on disk using `FileSystemStore`:

```go
store := grove.NewFileSystemStore("./templates")
eng := grove.New(grove.WithStore(store))

result, err := eng.Render(
	context.Background(),
	"index.html",    // loads ./templates/index.html
	grove.Data{"title": "Home"},
)
```

Template names are forward-slash paths relative to the store root. `FileSystemStore` rejects absolute paths and `..` traversal for security.

`Render` loads the template from the store by name, compiles it (with LRU caching), and executes it. Use `Render` instead of `RenderTemplate` when working with stored templates — it's required for `extends`, `include`, `render`, `import`, and `component` tags.

## In-Memory Templates

For testing or dynamic templates, use `MemoryStore`:

```go
store := grove.NewMemoryStore()
store.Set("greeting.html", `Hello, {{ name }}!`)
store.Set("base.html", `<html>{% block content %}{% endblock %}</html>`)

eng := grove.New(grove.WithStore(store))

result, _ := eng.Render(ctx, "greeting.html", grove.Data{"name": "Bob"})
fmt.Println(result.Body) // Hello, Bob!
```

`MemoryStore` is thread-safe. You can add templates with `Set` at any time.

## Passing Data

### Maps and slices

Pass nested maps and slices — templates access them with dot notation and bracket indexing:

```go
data := grove.Data{
	"user": map[string]any{
		"name": "Alice",
		"tags": []any{"admin", "editor"},
	},
}
```

```jinja2
{{ user.name }}      {# Alice #}
{{ user.tags[0] }}   {# admin #}
```

### Custom Go types

Implement the `Resolvable` interface to expose specific fields from Go structs:

```go
type User struct {
	Name     string
	Email    string  // not exposed to templates
	Internal int     // not exposed to templates
}

func (u User) WispyResolve(key string) (any, bool) {
	switch key {
	case "name":
		return u.Name, true
	}
	return nil, false
}
```

```jinja2
{{ user.name }}   {# works — returns "Alice" #}
{{ user.email }}  {# empty — not exposed by WispyResolve #}
```

Only keys returned by `WispyResolve` are accessible in templates. This lets you control exactly what data is visible to template authors.

## Global Variables

Register variables available in every render call:

```go
eng := grove.New()
eng.SetGlobal("site_name", "My Blog")
eng.SetGlobal("current_year", "2026")
```

```jinja2
<footer>© {{ current_year }} {{ site_name }}</footer>
```

Globals have the lowest priority. Render data overrides globals, and local variables (from `set`, `let`, `for`) override render data.

## Engine Options

| Option | Description |
|--------|-------------|
| `WithStore(store)` | Set the template store for named template loading |
| `WithStrictVariables(true)` | Return a `RuntimeError` for undefined variable access (default: render as empty) |
| `WithCacheSize(n)` | Set LRU cache size for compiled bytecode (default: 512) |
| `WithSandbox(cfg)` | Restrict allowed tags, filters, and loop iterations |

```go
eng := grove.New(
	grove.WithStore(grove.NewFileSystemStore("./templates")),
	grove.WithStrictVariables(true),
	grove.WithCacheSize(1024),
	grove.WithSandbox(grove.SandboxConfig{
		AllowedTags:    []string{"if", "for", "set", "component"},
		AllowedFilters: []string{"upper", "lower", "escape", "safe"},
		MaxLoopIter:    10000,
	}),
)
```

See [API Reference](api-reference.md) for full details on each option.

## Error Handling

Grove returns two error types:

**`ParseError`** — syntax errors detected during compilation. Includes `Template`, `Line`, and `Column` fields:

```go
result, err := eng.RenderTemplate(ctx, `{% if %}oops{% endif %}`, nil)
if err != nil {
	var pe grove.ParseError
	if errors.As(err, &pe) {
		fmt.Printf("Parse error at line %d: %s\n", pe.Line, pe.Error())
	}
}
```

**`RuntimeError`** — errors during template execution (division by zero, missing required props, strict mode undefined variables):

```go
result, err := eng.Render(ctx, "page.html", data)
if err != nil {
	var re grove.RuntimeError
	if errors.As(err, &re) {
		fmt.Printf("Runtime error: %s\n", re.Error())
	}
}
```
```

- [ ] **Step 2: Commit**

```bash
git add docs/getting-started.md
git commit -m "docs: add getting started guide"
```

---

## Task 4: docs/template-syntax.md

**Files:**
- Create: `docs/template-syntax.md`

- [ ] **Step 1: Write docs/template-syntax.md**

```markdown
# Template Syntax

## Delimiters

Grove uses three delimiter pairs:

| Delimiter | Purpose | Example |
|-----------|---------|---------|
| `{{ }}` | Output expression | `{{ name }}` |
| `{% %}` | Tags (control flow, assignment) | `{% if active %}` |
| `{# #}` | Comments (not rendered) | `{# TODO: fix this #}` |

### Whitespace control

Add `-` inside any delimiter to strip adjacent whitespace:

```jinja2
{%- if active -%}   {# strips whitespace on both sides #}
{{- name -}}         {# strips whitespace on both sides #}
```

`-` on the left strips all preceding whitespace (back to previous output). `-` on the right strips all following whitespace (up to next output).

## Variables

Access data passed to the template:

```jinja2
{{ name }}              {# simple variable #}
{{ user.name }}         {# dot access #}
{{ user["name"] }}      {# bracket access (equivalent) #}
{{ items[0] }}          {# index access #}
{{ users[0].address.city }}  {# chained access #}
```

Undefined variables render as empty string by default. With `WithStrictVariables(true)`, they return a `RuntimeError`.

## Expressions

### Operators

Ordered by precedence (highest to lowest):

| Precedence | Operator | Description |
|------------|----------|-------------|
| 1 | `.`, `[]`, `()` | Attribute access, index, function call |
| 2 | `\|` | Filter pipe |
| 3 | `not` | Logical negation |
| 4 | `*`, `/`, `%` | Multiplication, division, modulo |
| 5 | `+`, `-`, `~` | Addition, subtraction, string concatenation |
| 6 | `<`, `<=`, `>`, `>=`, `==`, `!=` | Comparison |
| 7 | `and` | Logical AND |
| 8 | `or` | Logical OR |
| 9 | `? :` | Ternary |

### Arithmetic

```jinja2
{{ price * quantity }}       {# multiplication #}
{{ total / count }}          {# division #}
{{ index % 2 }}              {# modulo #}
{{ base + tax }}             {# addition #}
{{ "Hello" ~ " " ~ name }}  {# string concatenation #}
```

### Comparison and logic

```jinja2
{{ age >= 18 }}          {# true/false #}
{{ role == "admin" }}
{{ active and verified }}
{{ banned or suspended }}
{{ not disabled }}
```

### Ternary expressions

```jinja2
{{ active ? "yes" : "no" }}
{{ user ? user.name : "Anonymous" }}
```

Ternary nests right-to-left (like JavaScript):

```jinja2
{{ a ? "A" : b ? "B" : "C" }}
{# equivalent to: a ? "A" : (b ? "B" : "C") #}
```

Filters bind tighter than `?`, so use parentheses if filtering the condition:

```jinja2
{{ (name | length) ? name : "unnamed" }}
```

## List Literals

```jinja2
{% set colors = ["red", "green", "blue"] %}
{% set matrix = [[1, 2], [3, 4]] %}
{% set empty = [] %}

{{ colors[0] }}          {# red #}
{{ matrix[1][0] }}       {# 3 #}
{{ colors | join(", ") }} {# red, green, blue #}
```

Trailing commas are allowed: `["a", "b",]`.

## Map Literals

```jinja2
{% set theme = {bg: "#fff", fg: "#000", border: "#ccc"} %}
{% set nested = {card: {padding: "1rem", shadow: true}} %}
{% set empty = {} %}

{{ theme.bg }}           {# #fff #}
{{ theme["fg"] }}        {# #000 #}
{{ nested.card.padding }} {# 1rem #}
```

Keys are unquoted identifiers. Trailing commas are allowed. Maps preserve insertion order — iterating with `for` or using `keys`/`values` filters returns entries in declaration order.

Maps and lists nest freely:

```jinja2
{% set data = {
  users: [
    {name: "Alice", role: "admin"},
    {name: "Bob", role: "editor"}
  ]
} %}
{{ data.users[0].name }}  {# Alice #}
```

## Filters

Filters transform values using pipe syntax:

```jinja2
{{ name | upper }}                    {# ALICE #}
{{ name | lower | title }}            {# Alice (chained) #}
{{ text | truncate(100) }}            {# with arguments #}
{{ text | replace("old", "new") }}    {# multiple arguments #}
```

See [Filters](filters.md) for the complete catalog of 42 built-in filters.

## Control Flow

### if / elif / else

```jinja2
{% if user.admin %}
  <span class="badge">Admin</span>
{% elif user.role == "editor" %}
  <span class="badge">Editor</span>
{% else %}
  <span class="badge">Member</span>
{% endif %}
```

**Truthy/falsy rules:** `nil`, `false`, `0`, `""` (empty string), empty lists `[]`, and empty maps `{}` are falsy. Everything else is truthy.

### for loops

Iterate over lists:

```jinja2
{% for item in items %}
  <li>{{ item }}</li>
{% endfor %}
```

With an `{% empty %}` fallback for empty collections:

```jinja2
{% for post in posts %}
  <article>{{ post.title }}</article>
{% empty %}
  <p>No posts yet.</p>
{% endfor %}
```

Iterate with index:

```jinja2
{% for i, item in items %}
  <li>{{ i }}: {{ item }}</li>
{% endfor %}
```

Iterate over maps (keys are sorted lexicographically):

```jinja2
{% for key, value in config %}
  {{ key }}: {{ value }}
{% endfor %}
```

#### Loop variables

Inside every `for` loop, a `loop` variable is automatically available:

| Variable | Description |
|----------|-------------|
| `loop.index` | 1-based position |
| `loop.index0` | 0-based position |
| `loop.first` | `true` if first iteration |
| `loop.last` | `true` if last iteration |
| `loop.length` | Total number of items |
| `loop.depth` | Nesting depth (1 for outermost loop) |
| `loop.parent` | Reference to the enclosing loop's `loop` variable |

```jinja2
{% for item in items %}
  {{ loop.index }}/{{ loop.length }}: {{ item }}
  {% if loop.first %}(first){% endif %}
  {% if loop.last %}(last){% endif %}
{% endfor %}
```

Nested loop example:

```jinja2
{% for row in rows %}
  {% for cell in row %}
    [{{ loop.parent.index }},{{ loop.index }}] = {{ cell }}
  {% endfor %}
{% endfor %}
```

### range

Generate numeric sequences:

```jinja2
{% for i in range(5) %}{{ i }}{% endfor %}
{# 0 1 2 3 4 #}

{% for i in range(1, 4) %}{{ i }}{% endfor %}
{# 1 2 3 #}

{% for i in range(10, 0, -2) %}{{ i }}{% endfor %}
{# 10 8 6 4 2 #}
```

## Variable Assignment

### set

Assign a single variable:

```jinja2
{% set greeting = "Hello, " ~ name %}
{% set total = price * quantity %}
{% set colors = ["red", "green", "blue"] %}
{{ greeting }}
```

Variables set inside a `for` loop persist after the loop ends.

### let

Assign multiple variables with optional conditionals:

```jinja2
{% let %}
  bg = "#d1ecf1"
  fg = "#0c5460"
  icon = "i"

  if type == "warning"
    bg = "#fff3cd"
    fg = "#856404"
    icon = "!"
  elif type == "error"
    bg = "#f8d7da"
    fg = "#721c24"
    icon = "x"
  end
{% endlet %}

<div style="background: {{ bg }}; color: {{ fg }}">
  {{ icon }} {{ message }}
</div>
```

**Rules:**
- Each line is `name = expression` (no `{% %}` delimiters inside the block)
- `if`, `elif`, `else`, `end` for conditionals (not `endif` — use `end`)
- Nested `if` blocks are allowed
- Expressions support the full syntax: filters, math, ternary, map/list literals
- Multi-line expressions work (e.g., a map literal spanning multiple lines) — the parser looks for `name =` to detect the next assignment
- Blank lines are ignored
- All variables are written to the outer scope (available after `{% endlet %}`)
- No output is produced inside the block

```jinja2
{% let %}
  themes = {
    warning: {bg: "#fff3cd", fg: "#856404"},
    error: {bg: "#f8d7da", fg: "#721c24"},
    info: {bg: "#d1ecf1", fg: "#0c5460"}
  }
  t = themes[type] | default(themes.info)
{% endlet %}
```

### capture

Render a block into a variable instead of outputting it:

```jinja2
{% capture greeting %}
  Hello, {{ name | title }}!
{% endcapture %}

{{ greeting | trim }}
```

The captured content is a string. You can filter or manipulate it after capture.

## Comments

```jinja2
{# This is a comment — not rendered in output #}

{# 
  Multi-line comments
  work too
#}
```

## Raw Blocks

Output template delimiters literally without parsing:

```jinja2
{% raw %}
  {{ this is not parsed }}
  {% neither is this %}
{% endraw %}
```
```

- [ ] **Step 2: Commit**

```bash
git add docs/template-syntax.md
git commit -m "docs: add template syntax reference"
```

---

## Task 5: docs/template-inheritance.md

**Files:**
- Create: `docs/template-inheritance.md`

- [ ] **Step 1: Write docs/template-inheritance.md**

```markdown
# Template Inheritance

Template inheritance lets you define a base layout and override specific sections in child templates.

## Base Template

A base template defines the page structure with `{% block %}` override points:

```jinja2
{# base.html #}
<!DOCTYPE html>
<html>
<head>
  <title>{% block title %}My Site{% endblock %}</title>
</head>
<body>
  <nav>...</nav>
  <main>
    {% block content %}{% endblock %}
  </main>
  <footer>
    {% block footer %}© 2026 My Site{% endblock %}
  </footer>
</body>
</html>
```

Blocks can have default content (like `title` and `footer` above) or be empty (like `content`). A base template renders on its own — blocks use their default content when not overridden.

## Child Template

A child template extends a parent with `{% extends %}` and overrides specific blocks:

```jinja2
{# home.html #}
{% extends "base.html" %}

{% block title %}Home — My Site{% endblock %}

{% block content %}
  <h1>Welcome</h1>
  <p>This replaces the content block.</p>
{% endblock %}
```

**Rules:**
- `{% extends %}` must be the first tag in the template
- Only `{% block %}` tags in the child are used — any content outside blocks is discarded
- Blocks not overridden keep the parent's default content
- `extends` requires a template store (`WithStore`) — it does not work with inline `RenderTemplate`

## super()

Include the parent block's content using `{{ super() }}`:

```jinja2
{# home.html #}
{% extends "base.html" %}

{% block title %}Home — {{ super() }}{% endblock %}
```

If the base template's `title` block contains `My Site`, this renders: `Home — My Site`.

## Multi-Level Inheritance

Inheritance chains to any depth. Each level can override blocks and call `super()`:

```jinja2
{# base.html #}
<html>
<body>
  {% block content %}base{% endblock %}
</body>
</html>
```

```jinja2
{# section.html #}
{% extends "base.html" %}

{% block content %}
  <div class="section">
    {% block inner %}section default{% endblock %}
  </div>
{% endblock %}
```

```jinja2
{# page.html #}
{% extends "section.html" %}

{% block inner %}page content{% endblock %}
```

Rendering `page.html` produces:

```html
<html>
<body>
  <div class="section">
    page content
  </div>
</body>
</html>
```

### super() chains

Each `super()` call reaches one level up. In a three-level chain:

```jinja2
{# base.html #}
{% block title %}Base{% endblock %}

{# mid.html #}
{% extends "base.html" %}
{% block title %}Mid:{{ super() }}{% endblock %}

{# leaf.html #}
{% extends "mid.html" %}
{% block title %}Leaf:{{ super() }}{% endblock %}
```

Rendering `leaf.html` produces: `Leaf:Mid:Base`.
```

- [ ] **Step 2: Commit**

```bash
git add docs/template-inheritance.md
git commit -m "docs: add template inheritance guide"
```

---

## Task 6: docs/components.md

**Files:**
- Create: `docs/components.md`

- [ ] **Step 1: Write docs/components.md**

```markdown
# Components

Components are reusable templates with a declared interface. They accept data through **props** and allow callers to inject content through **slots**.

## Using a Component

```jinja2
{% component "components/card.html" title="Hello" summary="A card" %}
  <p>This goes into the default slot.</p>
{% endcomponent %}
```

The first argument is the template path (loaded from the store). Remaining arguments are space-separated `key=value` props passed to the component.

`component` requires a template store — it does not work with inline `RenderTemplate`.

## Defining Props

Declare accepted props at the top of a component template with `{% props %}`:

```jinja2
{# components/button.html #}
{% props label, href="/", variant="primary" %}

<a href="{{ href }}" class="btn btn-{{ variant }}">{{ label }}</a>
```

- Props with a default value (like `href` and `variant`) are optional
- Props without a default (like `label`) are required — passing no value causes a `RuntimeError`
- Passing an unknown prop causes a `RuntimeError`
- If a component has no `{% props %}` declaration, it accepts any props without restriction

## Default Slot

Content between `{% component %}` and `{% endcomponent %}` fills the default slot:

```jinja2
{# components/box.html #}
<div class="box">
  {% slot %}No content provided{% endslot %}
</div>
```

```jinja2
{# Using it: #}
{% component "components/box.html" %}
  <p>This replaces "No content provided"</p>
{% endcomponent %}
```

The text inside `{% slot %}...{% endslot %}` is fallback content, rendered when the caller doesn't provide any.

## Named Slots

Components can define multiple injection points with named slots:

```jinja2
{# components/card.html #}
{% props title, summary %}

<article>
  <h2>{{ title }}</h2>
  <p>{{ summary }}</p>
  <div class="tags">
    {% slot "tags" %}{% endslot %}
  </div>
  <div class="actions">
    {% slot "actions" %}<a href="#">Read more</a>{% endslot %}
  </div>
</article>
```

Callers fill named slots with `{% fill %}`:

```jinja2
{% component "components/card.html" title="My Post" summary="A summary" %}
  {% fill "tags" %}
    <span class="tag">Go</span>
    <span class="tag">Templates</span>
  {% endfill %}
  {% fill "actions" %}
    <a href="/post/1">Read</a>
    <a href="/post/1/edit">Edit</a>
  {% endfill %}
{% endcomponent %}
```

Unfilled named slots render their fallback content.

## Scope Rules

This is the key design decision in Grove's component system:

- **Props** are available inside the component template. The component cannot see the caller's variables.
- **Fills see the caller's scope**, not the component's. This means you can use your page data inside a `{% fill %}` block without threading it through props.

```jinja2
{# page.html — caller's scope has "posts" #}
{% component "components/card.html" title="Recent" summary="Latest posts" %}
  {% fill "tags" %}
    {# This sees "posts" from the page, not from the card component #}
    {% for post in posts %}
      <span>{{ post.title }}</span>
    {% endfor %}
  {% endfill %}
{% endcomponent %}
```

## Nesting Components

Components can use other components:

```jinja2
{# components/post-list.html #}
{% props posts %}
{% for post in posts %}
  {% component "components/card.html" title=post.title summary=post.summary %}
    {% fill "tags" %}
      {% for tag in post.tags %}
        {% component "components/tag.html" label=tag.name color=tag.color %}{% endcomponent %}
      {% endfor %}
    {% endfill %}
  {% endcomponent %}
{% endfor %}
```

Components can also use template inheritance (`{% extends %}`).
```

- [ ] **Step 2: Commit**

```bash
git add docs/components.md
git commit -m "docs: add components guide"
```

---

## Task 7: docs/macros-and-includes.md

**Files:**
- Create: `docs/macros-and-includes.md`

- [ ] **Step 1: Write docs/macros-and-includes.md**

```markdown
# Macros & Includes

## include

Include a template inline. The included template shares the current scope:

```jinja2
{% include "partials/nav.html" %}
```

Pass additional variables:

```jinja2
{% include "partials/nav.html" section="about" active=true %}
```

The included template sees all variables from the current scope plus the explicitly passed ones.

## render

Like `include`, but with an isolated scope — only explicitly passed variables are visible:

```jinja2
{% render "partials/card.html" title="Widget" price=9.99 %}
```

The rendered template cannot access the caller's variables. Use `render` when you want self-contained fragments that don't depend on page context.

## include vs render

| | `include` | `render` |
|--|-----------|----------|
| **Scope** | Shared — sees caller's variables | Isolated — only passed variables |
| **Use when** | Partial needs page context | Fragment should be self-contained |
| **Example** | Navigation bar that needs `current_page` | Email template snippet |

Both require a template store (`WithStore`).

## macro

Define reusable template functions:

```jinja2
{% macro user_card(name, role="member") %}
  <div class="card">
    <strong>{{ name }}</strong>
    <span class="role">{{ role }}</span>
  </div>
{% endmacro %}
```

Call a macro like a function:

```jinja2
{{ user_card("Alice", "admin") }}
{{ user_card("Bob") }}
```

Macros support positional and keyword arguments:

```jinja2
{% macro link(href, text, target="_self") %}
  <a href="{{ href }}" target="{{ target }}">{{ text }}</a>
{% endmacro %}

{{ link("https://example.com", "Example") }}
{{ link("https://example.com", "Example", target="_blank") }}
```

**Macros have isolated scope** — they cannot access variables from the surrounding template. Only the arguments passed to the macro are available inside it.

## call and caller()

Use `{% call %}` to pass a block of content to a macro:

```jinja2
{% macro card(title) %}
  <div class="card">
    <h2>{{ title }}</h2>
    <div class="body">
      {{ caller() }}
    </div>
  </div>
{% endmacro %}

{% call card("Orders") %}
  <p>You have 3 pending orders.</p>
{% endcall %}
```

Inside the macro, `{{ caller() }}` renders the content from the `{% call %}` block. `caller()` can be called multiple times.

## import

Import macros from another template file into a namespace:

```jinja2
{% import "macros/ui.html" as ui %}

{{ ui.user_card("Alice") }}
{{ ui.link("https://example.com", "Click here") }}
```

`import` requires a template store. The imported template is executed, and any macros defined in it become available through the namespace.
```

- [ ] **Step 2: Commit**

```bash
git add docs/macros-and-includes.md
git commit -m "docs: add macros and includes guide"
```

---

## Task 8: docs/filters.md

**Files:**
- Create: `docs/filters.md`

- [ ] **Step 1: Write docs/filters.md**

```markdown
# Filters

Filters transform values using pipe syntax. They can be chained and accept arguments:

```jinja2
{{ name | upper }}                       {# ALICE #}
{{ name | trim | lower | title }}        {# Alice #}
{{ text | truncate(50, "…") }}           {# First 50 chars… #}
{{ items | sort | join(", ") }}          {# a, b, c #}
```

## String Filters

#### `upper`

`value | upper`

Converts string to uppercase.

```jinja2
{{ "hello" | upper }}  →  HELLO
```

#### `lower`

`value | lower`

Converts string to lowercase.

```jinja2
{{ "HELLO" | lower }}  →  hello
```

#### `title`

`value | title`

Capitalizes the first letter of each word.

```jinja2
{{ "hello world" | title }}  →  Hello World
```

#### `capitalize`

`value | capitalize`

Capitalizes the first letter, lowercases the rest.

```jinja2
{{ "hello WORLD" | capitalize }}  →  Hello world
```

#### `trim`

`value | trim`

Strips leading and trailing whitespace.

```jinja2
{{ "  hello  " | trim }}  →  hello
```

#### `lstrip`

`value | lstrip`

Strips leading whitespace only.

```jinja2
{{ "  hello  " | lstrip }}  →  hello  
```

#### `rstrip`

`value | rstrip`

Strips trailing whitespace only.

```jinja2
{{ "  hello  " | rstrip }}  →    hello
```

#### `replace`

`value | replace(old, new)` or `value | replace(old, new, count)`

Replaces occurrences of `old` with `new`. Optional `count` limits replacements.

```jinja2
{{ "hello world" | replace("world", "Grove") }}  →  hello Grove
{{ "aaa" | replace("a", "b", 2) }}  →  bba
```

#### `truncate`

`value | truncate(length, suffix)`

Truncates string to `length` characters and appends `suffix`. Defaults: length=255, suffix="...".

```jinja2
{{ "Hello, World!" | truncate(5) }}  →  He...
{{ "Hello, World!" | truncate(8, "…") }}  →  Hello…
```

#### `center`

`value | center(width, fill)`

Centers string within `width` using `fill` character. Default fill: space.

```jinja2
{{ "hi" | center(10) }}  →      hi    
{{ "hi" | center(10, "-") }}  →  ----hi----
```

#### `ljust`

`value | ljust(width, fill)`

Left-justifies string within `width`. Default fill: space.

```jinja2
{{ "hi" | ljust(10, ".") }}  →  hi........
```

#### `rjust`

`value | rjust(width, fill)`

Right-justifies string within `width`. Default fill: space.

```jinja2
{{ "hi" | rjust(10, ".") }}  →  ........hi
```

#### `split`

`value | split(separator)`

Splits string into a list. Default separator: space.

```jinja2
{{ "a,b,c" | split(",") | join(" ") }}  →  a b c
```

#### `wordcount`

`value | wordcount`

Returns the number of words in a string.

```jinja2
{{ "hello beautiful world" | wordcount }}  →  3
```

## Collection Filters

#### `length`

`value | length`

Returns the length of a list, map, or string (by rune count for strings).

```jinja2
{{ [1, 2, 3] | length }}   →  3
{{ "hello" | length }}      →  5
{{ {a: 1, b: 2} | length }} →  2
```

#### `first`

`value | first`

Returns the first element of a list. Returns nil for empty lists.

```jinja2
{{ ["a", "b", "c"] | first }}  →  a
```

#### `last`

`value | last`

Returns the last element of a list. Returns nil for empty lists.

```jinja2
{{ ["a", "b", "c"] | last }}  →  c
```

#### `join`

`value | join(separator)`

Joins list elements into a string. Default separator: empty string.

```jinja2
{{ ["a", "b", "c"] | join(", ") }}  →  a, b, c
{{ [1, 2, 3] | join("-") }}  →  1-2-3
```

#### `sort`

`value | sort`

Sorts list elements as strings (stable sort).

```jinja2
{{ ["banana", "apple", "cherry"] | sort | join(", ") }}  →  apple, banana, cherry
```

#### `reverse`

`value | reverse`

Reverses a list or string.

```jinja2
{{ ["a", "b", "c"] | reverse | join("") }}  →  cba
{{ "hello" | reverse }}  →  olleh
```

#### `unique`

`value | unique`

Removes duplicate elements, preserving order.

```jinja2
{{ ["a", "b", "a", "c", "b"] | unique | join(", ") }}  →  a, b, c
```

#### `min`

`value | min`

Returns the minimum value in a list. Compares numerically if possible, otherwise as strings.

```jinja2
{{ [3, 1, 2] | min }}  →  1
```

#### `max`

`value | max`

Returns the maximum value in a list. Compares numerically if possible, otherwise as strings.

```jinja2
{{ [3, 1, 2] | max }}  →  3
```

#### `sum`

`value | sum`

Returns the sum of numeric values in a list.

```jinja2
{{ [1, 2, 3] | sum }}  →  6
{{ [1.5, 2.5] | sum }}  →  4
```

#### `map`

`value | map(attribute)`

Extracts an attribute from each item in a list.

```jinja2
{% set users = [{name: "Alice"}, {name: "Bob"}] %}
{{ users | map("name") | join(", ") }}  →  Alice, Bob
```

#### `batch`

`value | batch(size)`

Groups a list into batches (sub-lists) of the given size. Default size: 1.

```jinja2
{% for row in [1,2,3,4,5] | batch(2) %}
  {{ row | join(",") }}
{% endfor %}
{# 1,2 then 3,4 then 5 #}
```

#### `flatten`

`value | flatten`

Flattens nested lists one level deep.

```jinja2
{{ [[1, 2], [3, 4], [5]] | flatten | join(",") }}  →  1,2,3,4,5
```

#### `keys`

`value | keys`

Returns the keys of a map as a list. For map literals, returns keys in insertion order. For Go maps passed as data, returns keys sorted lexicographically.

```jinja2
{% set m = {b: 2, a: 1} %}
{{ m | keys | join(",") }}  →  b,a
```

#### `values`

`value | values`

Returns the values of a map as a list. For map literals, returns values in insertion order. For Go maps passed as data, returns values in sorted key order.

```jinja2
{% set m = {b: 2, a: 1} %}
{{ m | values | join(",") }}  →  2,1
```

## Numeric Filters

#### `abs`

`value | abs`

Returns the absolute value.

```jinja2
{{ -5 | abs }}  →  5
{{ -3.14 | abs }}  →  3.14
```

#### `round`

`value | round(precision)`

Rounds to the given precision. Default: 0. Returns an integer when precision is 0.

```jinja2
{{ 3.7 | round }}  →  4
{{ 3.14159 | round(2) }}  →  3.14
```

#### `ceil`

`value | ceil`

Returns the ceiling (rounds up to nearest integer).

```jinja2
{{ 3.2 | ceil }}  →  4
```

#### `floor`

`value | floor`

Returns the floor (rounds down to nearest integer).

```jinja2
{{ 3.8 | floor }}  →  3
```

#### `int`

`value | int`

Converts to integer.

```jinja2
{{ "42" | int }}  →  42
{{ 3.9 | int }}  →  3
```

#### `float`

`value | float`

Converts to float.

```jinja2
{{ "3.14" | float }}  →  3.14
{{ 42 | float }}  →  42
```

## Logic & Type Filters

#### `default`

`value | default(fallback)`

Returns `fallback` if the value is falsy (nil, false, 0, empty string, empty list, empty map).

```jinja2
{{ name | default("Anonymous") }}
{{ items | default([]) }}
```

#### `string`

`value | string`

Converts a value to its string representation.

```jinja2
{{ 42 | string }}  →  42
{{ true | string }}  →  true
```

#### `bool`

`value | bool`

Converts a value to boolean using truthy/falsy rules.

```jinja2
{{ "" | bool }}   →  false
{{ "hi" | bool }} →  true
{{ 0 | bool }}    →  false
{{ 1 | bool }}    →  true
```

## HTML Filters

#### `escape`

`value | escape`

HTML-escapes special characters. Returns SafeHTML (won't be double-escaped).

```jinja2
{{ "<b>bold</b>" | escape }}  →  &lt;b&gt;bold&lt;/b&gt;
```

Note: auto-escaping is on by default for all `{{ }}` output, so you rarely need this filter explicitly. It's useful when you want to escape a value *before* passing it to another filter.

#### `safe`

`value | safe`

Marks a value as trusted HTML, bypassing auto-escaping.

```jinja2
{{ html_content | safe }}
```

**Use with caution** — only apply `safe` to content you trust. Untrusted content marked as safe creates XSS vulnerabilities.

#### `striptags`

`value | striptags`

Removes all HTML tags.

```jinja2
{{ "<p>Hello <b>world</b></p>" | striptags }}  →  Hello world
```

#### `nl2br`

`value | nl2br`

Converts newlines to `<br>` tags. HTML-escapes the input first, then returns SafeHTML.

```jinja2
{{ "line one\nline two" | nl2br }}  →  line one<br>
line two
```

## Custom Filters

Register custom filters on an engine:

```go
eng := grove.New()

// Simple filter — no arguments
eng.RegisterFilter("shout", grove.FilterFn(
	func(v grove.Value, args []grove.Value) (grove.Value, error) {
		s := v.String() + "!!!"
		return grove.StringValue(s), nil
	},
))

// Filter with arguments
eng.RegisterFilter("repeat", grove.FilterFn(
	func(v grove.Value, args []grove.Value) (grove.Value, error) {
		n := grove.ArgInt(args, 0, 1)
		s := strings.Repeat(v.String(), n)
		return grove.StringValue(s), nil
	},
))

// Filter that outputs trusted HTML (bypasses auto-escape)
eng.RegisterFilter("bold", grove.FilterFunc(
	grove.FilterFn(func(v grove.Value, args []grove.Value) (grove.Value, error) {
		return grove.StringValue("<b>" + v.String() + "</b>"), nil
	}),
	grove.FilterOutputsHTML(),
))
```

```jinja2
{{ name | shout }}       →  Alice!!!
{{ "ha" | repeat(3) }}   →  hahaha
{{ name | bold }}        →  <b>Alice</b>  (not escaped)
```

See [API Reference](api-reference.md) for details on `FilterFn`, `FilterDef`, and `FilterFunc`.
```

- [ ] **Step 2: Commit**

```bash
git add docs/filters.md
git commit -m "docs: add filter catalog"
```

---

## Task 9: docs/web-primitives.md

**Files:**
- Create: `docs/web-primitives.md`

- [ ] **Step 1: Write docs/web-primitives.md**

```markdown
# Web Primitives

Grove templates can declare CSS/JS assets, meta tags, and hoisted content. These are collected during rendering — including across nested includes, components, and inherited templates — and returned in `RenderResult` for the application to assemble into the final HTML response.

## asset

Declare a stylesheet or script dependency:

```jinja2
{% asset "/static/style.css" type="stylesheet" %}
{% asset "/static/app.js" type="script" %}
```

With priority and HTML attributes:

```jinja2
{% asset "/static/main.css" type="stylesheet" priority=10 %}
{% asset "/static/analytics.js" type="script" defer=true async=true %}
```

**Rules:**
- `type` is required — typically `"stylesheet"` or `"script"`
- `priority` controls sort order (higher = earlier within its type group). Default: 0
- Additional attributes (`defer`, `async`, `crossorigin`, etc.) are passed through as HTML attributes
- Boolean attributes use `attr=true` — rendered as bare attributes (e.g., `defer`)
- Assets are deduplicated by `Src` — declaring the same URL twice results in one entry
- Assets declared in components and includes bubble up to the top-level `RenderResult`

`asset` requires a template store — it does not work with inline `RenderTemplate`.

## meta

Declare document metadata:

```jinja2
{% meta name="description" content="A great page" %}
{% meta property="og:title" content="My Page" %}
{% meta property="og:image" content="https://example.com/image.png" %}
```

**Rules:**
- `name` or `property` serves as the key
- `content` is the value
- Last-write-wins for duplicate keys — a `Warning` is added to `RenderResult.Warnings` on collision
- Meta tags from components and includes bubble up

## hoist

Capture rendered content and collect it into a named target instead of outputting it inline:

```jinja2
{% hoist target="head" %}
  <link rel="preload" href="/font.woff2" as="font" crossorigin>
{% endhoist %}

{% hoist target="head" %}
  <style>.hero { background: blue; }</style>
{% endhoist %}
```

**Rules:**
- `target` names the collection bucket (any string)
- Multiple hoists to the same target are concatenated in order
- Hoisted content is removed from `Body` and collected in `RenderResult.Hoisted`
- Hoisted content from components and includes bubbles up

## RenderResult

When you call `Render` or `RenderTemplate`, Grove returns a `RenderResult`:

```go
type RenderResult struct {
    Body     string              // Rendered HTML output
    Assets   []Asset             // Collected assets, deduplicated by Src
    Meta     map[string]string   // Collected meta tags (last-write-wins)
    Hoisted  map[string][]string // target → ordered fragments
    Warnings []Warning           // Non-fatal warnings (e.g., meta key collision)
}
```

### Helper methods

**`HeadHTML()`** — returns `<link rel="stylesheet">` tags for all stylesheet assets, sorted by descending priority:

```go
result.HeadHTML()
// <link rel="stylesheet" href="/static/main.css">
// <link rel="stylesheet" href="/static/theme.css">
```

**`FootHTML()`** — returns `<script>` tags for all script assets, sorted by descending priority:

```go
result.FootHTML()
// <script src="/static/app.js" defer></script>
```

**`GetHoisted(target)`** — returns concatenated content for a hoist target:

```go
result.GetHoisted("head")
// <link rel="preload" href="/font.woff2" as="font" crossorigin>
// <style>.hero { background: blue; }</style>
```

## Integration Pattern

A typical web application renders a template and then injects collected assets and meta into the response. Here's the complete pattern:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    result, err := eng.Render(r.Context(), "page.html", grove.Data{
        "title": "My Page",
    })
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    body := result.Body

    // Inject stylesheet assets into <head>
    body = strings.Replace(body, "<!-- HEAD_ASSETS -->", result.HeadHTML(), 1)

    // Inject script assets before </body>
    body = strings.Replace(body, "<!-- FOOT_ASSETS -->", result.FootHTML(), 1)

    // Build and inject meta tags
    var meta strings.Builder
    for name, content := range result.Meta {
        if strings.HasPrefix(name, "og:") {
            meta.WriteString(fmt.Sprintf(`  <meta property="%s" content="%s">`+"\n", name, content))
        } else {
            meta.WriteString(fmt.Sprintf(`  <meta name="%s" content="%s">`+"\n", name, content))
        }
    }
    body = strings.Replace(body, "<!-- HEAD_META -->", meta.String(), 1)

    // Inject hoisted content
    body = strings.Replace(body, "<!-- HEAD_HOISTED -->", result.GetHoisted("head"), 1)

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, body)
}
```

The base template uses placeholder comments that get replaced:

```jinja2
<head>
  <title>{% block title %}My Site{% endblock %}</title>
  <!-- HEAD_ASSETS -->
  <!-- HEAD_META -->
  <!-- HEAD_HOISTED -->
</head>
<body>
  {% block content %}{% endblock %}
  <!-- FOOT_ASSETS -->
</body>
```

This pattern keeps template authors and application developers in their own domains — templates declare what they need, and the Go layer assembles it.

## Auto-Escaping

All `{{ }}` output is HTML-escaped by default. This prevents XSS from untrusted data:

```jinja2
{% set input = "<script>alert('xss')</script>" %}
{{ input }}
{# Output: &lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt; #}
```

To output trusted HTML, use the `safe` filter:

```jinja2
{{ trusted_html | safe }}
```

Or use `{% raw %}` blocks to output template syntax literally (no parsing or escaping):

```jinja2
{% raw %}{{ not parsed }}{% endraw %}
```

From Go code, use `SafeHTMLValue` to mark a value as pre-trusted:

```go
data := grove.Data{
    "content": grove.SafeHTMLValue("<p>Trusted HTML</p>"),
}
```

**Only mark content as safe when you trust the source.** Auto-escaping exists to protect against XSS — bypassing it with untrusted data creates vulnerabilities.
```

- [ ] **Step 2: Commit**

```bash
git add docs/web-primitives.md
git commit -m "docs: add web primitives guide"
```

---

## Task 10: docs/api-reference.md

**Files:**
- Create: `docs/api-reference.md`

- [ ] **Step 1: Write docs/api-reference.md**

```markdown
# API Reference

Import:

```go
import "grove/pkg/grove"
```

## Engine

```go
func New(opts ...Option) *Engine
```

Creates a new template engine. Safe for concurrent use — multiple goroutines can call render methods simultaneously.

### Rendering Methods

```go
func (e *Engine) RenderTemplate(ctx context.Context, src string, data Data) (RenderResult, error)
```

Compiles and renders an inline template string. Does not support `extends`, `include`, `render`, `import`, `component`, or `asset` tags (these require a store).

```go
func (e *Engine) Render(ctx context.Context, name string, data Data) (RenderResult, error)
```

Loads a named template from the store, compiles it (with caching), and renders it. Requires `WithStore`.

```go
func (e *Engine) RenderTo(ctx context.Context, name string, data Data, w io.Writer) error
```

Like `Render`, but streams output to an `io.Writer`. Does not return a `RenderResult` — use `Render` if you need access to collected assets, meta, or hoisted content.

```go
func (e *Engine) LoadTemplate(name string) (*compiler.Bytecode, error)
```

Compiles and caches a template without rendering it. Useful for pre-warming the cache.

### Engine Configuration

```go
func (e *Engine) SetGlobal(key string, value any)
```

Registers a global variable available in all renders. Globals have the lowest priority — render data overrides them, and template-local variables override render data.

```go
func (e *Engine) RegisterFilter(name string, fn any)
```

Registers a custom filter. `fn` can be a `FilterFn` or a `*FilterDef` (created with `FilterFunc`).

## Options

```go
func WithStore(s store.Store) Option
```

Sets the template store used by `Render`, `include`, `render`, `import`, and `component`.

```go
func WithStrictVariables(strict bool) Option
```

When true, accessing an undefined variable returns a `RuntimeError` instead of an empty string.

```go
func WithCacheSize(n int) Option
```

Sets the LRU cache capacity for compiled bytecode. Default: 512. Pass 0 to use the default.

```go
func WithSandbox(cfg SandboxConfig) Option
```

Applies security restrictions to all templates rendered by this engine.

## SandboxConfig

```go
type SandboxConfig struct {
    AllowedTags    []string  // nil = all allowed; non-nil = whitelist
    AllowedFilters []string  // nil = all allowed; non-nil = whitelist
    MaxLoopIter    int       // 0 = unlimited
}
```

- `AllowedTags`: when set, only listed tags are permitted. Others cause a `ParseError` at compile time.
- `AllowedFilters`: when set, only listed filters are permitted. Others cause a `ParseError` at compile time.
- `MaxLoopIter`: maximum total loop iterations across all loops in a single render. Exceeding this causes a `RuntimeError`.

```go
eng := grove.New(grove.WithSandbox(grove.SandboxConfig{
    AllowedTags:    []string{"if", "for", "set", "component"},
    AllowedFilters: []string{"upper", "lower", "escape", "safe", "default"},
    MaxLoopIter:    10000,
}))
```

## Data

```go
type Data map[string]any
```

The map type passed to render methods. Values can be any Go type: strings, numbers, booleans, slices (`[]any`), maps (`map[string]any`), or types implementing `Resolvable`.

## Resolvable

```go
type Resolvable interface {
    WispyResolve(key string) (any, bool)
}
```

Implement this interface on Go types to control which fields are accessible in templates. Only keys that return `(value, true)` are visible. All other field access returns empty (or errors in strict mode).

```go
type User struct {
    Name     string
    Email    string
    password string
}

func (u User) WispyResolve(key string) (any, bool) {
    switch key {
    case "name":
        return u.Name, true
    case "email":
        return u.Email, true
    }
    return nil, false
}
```

```jinja2
{{ user.name }}      {# "Alice" #}
{{ user.email }}     {# "alice@example.com" #}
{{ user.password }}  {# empty — not exposed #}
```

## Stores

### MemoryStore

```go
func NewMemoryStore() *MemoryStore
```

Creates an empty in-memory template store. Thread-safe.

```go
func (s *MemoryStore) Set(name, content string)
```

Adds or updates a template.

```go
store := grove.NewMemoryStore()
store.Set("base.html", `<html>{% block content %}{% endblock %}</html>`)
store.Set("page.html", `{% extends "base.html" %}{% block content %}Hello{% endblock %}`)
```

### FileSystemStore

```go
func NewFileSystemStore(root string) *FileSystemStore
```

Creates a store that loads templates from disk. Template names are forward-slash paths relative to `root`.

```go
store := grove.NewFileSystemStore("./templates")
eng := grove.New(grove.WithStore(store))

// Loads ./templates/pages/home.html
result, err := eng.Render(ctx, "pages/home.html", data)
```

**Security:** Rejects absolute paths and `..` traversal. Paths are cleaned and verified to stay within the root directory.

## RenderResult

```go
type RenderResult struct {
    Body     string
    Assets   []Asset
    Meta     map[string]string
    Hoisted  map[string][]string
    Warnings []Warning
}
```

See [Web Primitives](web-primitives.md) for detailed documentation on `RenderResult`, `Asset`, `Warning`, and the helper methods `HeadHTML()`, `FootHTML()`, and `GetHoisted()`.

## Filter Types

```go
type FilterFn func(v Value, args []Value) (Value, error)
```

The function signature for filter implementations. `v` is the piped value, `args` are any arguments passed in parentheses.

```go
type FilterDef struct { /* ... */ }
```

A filter with metadata. Create with `FilterFunc`:

```go
func FilterFunc(fn FilterFn, opts ...FilterOption) *FilterDef
```

```go
func FilterOutputsHTML() FilterOption
```

Marks a filter as returning trusted HTML, which bypasses auto-escaping.

```go
type FilterSet map[string]any
```

A named collection of filters for bulk registration.

## Value Types

```go
type Value /* opaque runtime value */
```

The template runtime value type. Used in filter functions.

```go
var Nil Value // zero value (nil)
```

```go
func StringValue(s string) Value
```

Wraps a Go string as a template `Value`.

```go
func SafeHTMLValue(s string) Value
```

Wraps trusted HTML as a `Value` — auto-escaping is skipped when this value is output.

```go
func ArgInt(args []Value, i, def int) int
```

Helper for filter implementations. Returns `args[i]` as an integer, or `def` if `i` is out of range.

## Error Types

### ParseError

```go
type ParseError struct {
    Template string
    Line     int
    Column   int
    // ...
}
```

Returned for syntax errors detected during compilation. `Template` is the template name (or empty for inline templates). `Line` and `Column` identify the source location.

### RuntimeError

```go
type RuntimeError struct {
    // ...
}
```

Returned for errors during template execution: division by zero, missing required props, strict mode undefined variables, sandbox loop limit exceeded.

Both error types implement the `error` interface. Use `errors.As` for type checking:

```go
var pe grove.ParseError
if errors.As(err, &pe) {
    fmt.Printf("Syntax error at line %d\n", pe.Line)
}
```
```

- [ ] **Step 2: Commit**

```bash
git add docs/api-reference.md
git commit -m "docs: add API reference"
```

---

## Task 11: docs/examples.md

**Files:**
- Create: `docs/examples.md`

- [ ] **Step 1: Write docs/examples.md**

```markdown
# Examples

## Blog Application

The `examples/blog/` directory contains a complete web application demonstrating Grove's features. It's a blog with posts, tags, components, template inheritance, and asset collection.

### Project Structure

```
examples/blog/
  main.go                           # Go web app (chi router)
  templates/
    base.grov                       # Root layout — nav, main, footer, asset placeholders
    index.grov                      # Homepage — extends base, lists posts
    post.grov                       # Post page — extends base, shows single post
    components/
      card.grov                     # Post card — props: title, summary, href, date; slot: tags
      nav.grov                      # Navigation bar — props: site_name; default slot
      footer.grov                   # Footer — props: year
      tag.grov                      # Color tag badge — props: label, color
      button.grov                   # Button link — props: label, href, variant
      alert.grov                    # Alert box — props: type; default slot; uses let block
    pages/
      styleguide.grov              # Component showcase page
```

### The Go Application

`main.go` sets up a Grove engine with a filesystem store and global variables:

```go
store := grove.NewFileSystemStore(templateDir)
eng := grove.New(grove.WithStore(store))
eng.SetGlobal("site_name", "Blog")
eng.SetGlobal("current_year", "2026")
```

Posts are Go structs implementing `Resolvable` to expose fields to templates:

```go
type Post struct {
    Title   string
    Slug    string
    Summary string
    Body    string
    Date    string
    Draft   bool
    Tags    []Tag
}

func (p Post) WispyResolve(key string) (any, bool) {
    switch key {
    case "title":
        return p.Title, true
    case "slug":
        return p.Slug, true
    // ... other fields
    }
    return nil, false
}
```

Handlers render templates and assemble the response by replacing placeholder comments with collected assets and meta:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    result, _ := eng.Render(r.Context(), "index.grov", grove.Data{
        "posts": postsAny,
    })

    body := result.Body
    body = strings.Replace(body, "<!-- HEAD_ASSETS -->", result.HeadHTML(), 1)
    body = strings.Replace(body, "<!-- FOOT_ASSETS -->", result.FootHTML(), 1)
    // ... meta tags, hoisted content
    w.Write([]byte(body))
}
```

### Base Layout

`base.grov` defines the HTML skeleton with blocks and asset placeholders:

```jinja2
{% asset "/static/style.css" type="stylesheet" priority=10 %}
<!DOCTYPE html>
<html lang="en">
<head>
  <title>{% block title %}Blog{% endblock %}</title>
  <!-- HEAD_ASSETS -->
  <!-- HEAD_META -->
  <!-- HEAD_HOISTED -->
</head>
<body>
  {% component "components/nav.html" site_name=site_name %}{% endcomponent %}
  <main>{% block content %}{% endblock %}</main>
  {% component "components/footer.html" year=current_year %}{% endcomponent %}
  <!-- FOOT_ASSETS -->
</body>
</html>
```

Every page inherits this layout. The base template declares a global stylesheet asset, uses components for nav and footer, and provides placeholder comments that the Go layer replaces.

### Page Templates

`index.grov` extends the base and iterates over posts using the card component:

```jinja2
{% extends "base.html" %}
{% block title %}Home — Blog{% endblock %}

{% block content %}
{% meta name="description" content="A blog built with Grove" %}

{% for post in posts %}
  {% component "components/card.html" title=post.title summary=post.summary href="/post/" ~ post.slug date=post.date %}
    {% fill "tags" %}
      {% for tag in post.tags %}
        {% component "components/tag.html" label=tag.name color=tag.color %}{% endcomponent %}
      {% endfor %}
    {% endfill %}
  {% endcomponent %}
{% endfor %}
{% endblock %}
```

This demonstrates nested components (tag inside card), slot fills, expression concatenation (`"/post/" ~ post.slug`), and meta tag declaration.

### Components

**card.grov** — shows props with defaults and a named slot:

```jinja2
{% props title, summary, href="#", date="" %}
<article>
  <h2><a href="{{ href }}">{{ title }}</a></h2>
  {% if date %}<time>{{ date }}</time>{% endif %}
  <p>{{ summary | truncate(120) }}</p>
  <div>{% slot "tags" %}{% endslot %}</div>
</article>
```

**alert.grov** — shows the `let` block with conditional variable assignment:

```jinja2
{% props type="info" %}
{% let %}
  bg = "#d1ecf1"
  fg = "#0c5460"
  icon = "i"

  if type == "warning"
    bg = "#fff3cd"
    fg = "#856404"
    icon = "!"
  elif type == "error"
    bg = "#f8d7da"
    fg = "#721c24"
    icon = "x"
  end
{% endlet %}
<div style="background: {{ bg }}; color: {{ fg }}">
  <span>{{ icon }}</span>
  <div>{% slot %}{% endslot %}</div>
</div>
```

**button.grov** — shows ternary expressions:

```jinja2
{% props label, href="/", variant="primary" %}
{% if variant == "primary" %}
  {% set bg = "#e94560" %}{% set fg = "#fff" %}
{% elif variant == "outline" %}
  {% set bg = "transparent" %}{% set fg = "#e94560" %}
{% else %}
  {% set bg = "#6c757d" %}{% set fg = "#fff" %}
{% endif %}

<a href="{{ href }}" style="background: {{ bg }}; color: {{ fg }}; border-color: {{ variant != "outline" ? bg : "#e94560" }}">{{ label }}</a>
```

### Running It

```bash
cd examples/blog
go run main.go
```

Open `http://localhost:3000` to see:
- **Home page** — list of post cards with tags
- **Post pages** — individual posts with draft warnings (alert component)
- **Component library** (`/styleguide`) — showcase of all components with variations
```

- [ ] **Step 2: Commit**

```bash
git add docs/examples.md
git commit -m "docs: add examples walkthrough"
```
