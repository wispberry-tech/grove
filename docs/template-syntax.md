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
