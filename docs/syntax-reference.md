# Template Syntax Reference

Complete reference for Wisp template syntax.

## Table of Contents

1. [Comments](#comments)
2. [Variables](#variables)
3. [Filters](#filters)
4. [Conditionals](#conditionals)
5. [Loops](#loops)
6. [Variable Assignment](#variable-assignment)
7. [Context Blocks](#context-blocks)
8. [Template Composition](#template-composition)
9. [Layout System](#layout-system)
10. [Capture](#capture)
11. [Break and Continue](#break-and-continue)

---

## Comments

Comments are not rendered in the output.

```liquid
{# This is a comment #}
{# Multi-line
   comments
   are supported #}
```

---

## Variables

### Simple Variables

Access variables with a leading dot:

```liquid
{% .name %}
{% .count %}
{% .is_active %}
```

### Nested Access

Access nested properties with dot notation:

```liquid
{% .user.name %}
{% .user.profile.avatar %}
{% .company.address.city %}
```

### Array/Map Indexing

Access array elements and map keys:

```liquid
{% .items[0] %}           {# First item #}
{% .items[-1] %}         {# Last item #}
{% .data[key] %}         {# Map key access #}
{% .matrix[0][1] %}      {# Nested array #}
```

### Chained Access

Combine all access patterns:

```liquid
{% .users[0].posts[0].title %}
{% .data[key].values[index] %}
```

---

## Filters

Filters transform values. Chain multiple filters with the pipe operator.

### Syntax

```liquid
{% .value | filter_name %}
{% .value | filter_name arg %}
{% .value | filter1 | filter2 | filter3 %}
```

### String Filters

| Filter | Description | Example |
|--------|-------------|---------|
| `upcase` | Convert to uppercase | `{% .name \| upcase %}` |
| `downcase` | Convert to lowercase | `{% .name \| downcase %}` |
| `capitalize` | Capitalize first letter | `{% .name \| capitalize %}` |
| `truncate n` | Truncate to n characters | `{% .text \| truncate 50 %}` |
| `strip` | Remove leading/trailing whitespace | `{% .text \| strip %}` |
| `lstrip` | Remove leading whitespace | `{% .text \| lstrip %}` |
| `rstrip` | Remove trailing whitespace | `{% .text \| rstrip %}` |
| `replace old new` | Replace substring | `{% .text \| replace "old" "new" %}` |
| `remove str` | Remove substring | `{% .text \| remove "x" %}` |
| `split str` | Split by delimiter | `{% .text \| split "," %}` |
| `join str` | Join array with separator | `{% .items \| join ", " %}` |
| `prepend str` | Prepend string | `{% .name \| prepend "Mr. " %}` |
| `append str` | Append string | `{% .name \| append " Jr." %}` |

### Numeric Filters

| Filter | Description | Example |
|--------|-------------|---------|
| `abs` | Absolute value | `{% .num \| abs %}` |
| `ceil` | Round up | `{% .num \| ceil %}` |
| `floor` | Round down | `{% .num \| floor %}` |
| `round` | Round to nearest | `{% .num \| round %}` |
| `plus n` | Add n | `{% .num \| plus 5 %}` |
| `minus n` | Subtract n | `{% .num \| minus 2 %}` |
| `times n` | Multiply by n | `{% .num \| times 2 %}` |
| `divided_by n` | Divide by n | `{% .num \| divided_by 2 %}` |
| `modulo n` | Modulo n | `{% .num \| modulo 3 %}` |

### Array Filters

| Filter | Description | Example |
|--------|-------------|---------|
| `first` | First element | `{% .items \| first %}` |
| `last` | Last element | `{% .items \| last %}` |
| `size` | Array length | `{% .items \| size %}` |
| `length` | Alias for size | `{% .items \| length %}` |
| `reverse` | Reverse array | `{% .items \| reverse %}` |
| `sort` | Sort array | `{% .items \| sort %}` |
| `uniq` | Remove duplicates | `{% .items \| uniq %}` |
| `map_field f` | Map field f | `{% .users \| map_field "name" %}` |

### Date Filters

| Filter | Description | Example |
|--------|-------------|---------|
| `date fmt` | Format date | `{% .date \| date "2006-01-02" %}` |
| `date_format fmt` | Format date | `{% .date \| date_format "Jan 2, 2006" %}` |

### URL Filters

| Filter | Description | Example |
|--------|-------------|---------|
| `url_encode` | URL encode | `{% .text \| url_encode %}` |
| `url_decode` | URL decode | `{% .text \| url_decode %}` |

### Utility Filters

| Filter | Description | Example |
|--------|-------------|---------|
| `default val` | Default if empty | `{% .val \| default "N/A" %}` |
| `json` | JSON encode | `{% .obj \| json %}` |
| `escape` | HTML escape | `{% .html \| escape %}` |
| `escape_once` | Escape only entities | `{% .html \| escape_once %}` |
| `raw` | Mark as safe | `{% .html \| raw %}` |

### Math Filters

| Filter | Description | Example |
|--------|-------------|---------|
| `min` | Minimum value | `{% .a \| min .b %}` |
| `max` | Maximum value | `{% .a \| max .b %}` |

---

## Conditionals

### If / Elsif / Else

```liquid
{% if .condition %}
    Content when true
{% elsif .other %}
    Content when other is true
{% else %}
    Content when all false
{% end %}
```

### Unless

`unless` is the opposite of `if`:

```liquid
{% unless .hide %}
    Content shown when .hide is false
{% end %}
```

Equivalent to:

```liquid
{% if not .hide %}
    Content shown when .hide is false
{% end %}
```

### Case / When

```liquid
{% case .status %}
    {% when "draft" %}
        Draft status
    {% when "published" %}
        Published
    {% when "archived" %}
        Archived
    {% else %}
        Unknown status
{% end %}
```

---

## Loops

### For Loop

```liquid
{% for .item in .items %}
    <li>{%.item%}</li>
{% end %}
```

With index:

```liquid
{% for .index, .item in .items %}
    <li>{%.index%}: {%.item%}</li>
{% end %}
```

### Range Loop

Loop over a numeric range:

```liquid
{% for .i in (range 1 5) %}
    {%.i%}
{% end %}
{# Outputs: 1 2 3 4 5 #}
```

### While Loop

```liquid
{% while .condition %}
    {%.value%}
    {% assign .value = .value | plus 1 %}
{% end %}
```

---

## Variable Assignment

### Assign

```liquid
{% assign .name = "value" %}
{% assign .count = .count | plus 1 %}
{% assign .user.name = "New Name" %}
```

---

## Context Blocks

### With

Isolate scope for a variable:

```liquid
{% with .user %}
    {% .name %}
    {% .email %}
{% end %}
```

### Cycle

Alternate between values:

```liquid
{% for .item in .items %}
    <div class="{% cycle "odd" "even" %}">
        {%.item%}
    </div>
{% end %}
```

### Increment/Decrement

```liquid
{% increment .counter %}
{% decrement .counter %}
```

---

## Template Composition

### Include

Include and evaluate another template:

```liquid
{% include "partials/header" %}
{% include "partials/footer" .data %}
```

### Render

Render with isolated scope:

```liquid
{% render "widgets/sidebar" .sidebar_data %}
```

### Component

Props-based component:

```liquid
{% component "Button" .buttonProps %}
{% component "Card" title=.title body=.body %}
```

---

## Layout System

### Extends

Child template extends a parent:

```liquid
{# child.html #}
{% extends "layouts/main" %}

{% block content %}
    Page content
{% endblock %}
```

### Block

Define or override a block:

```liquid
{% block title %}Default Title{% endblock %}

{% block content %}
    Default content
{% endblock %}
```

### Content

Provide content for parent blocks:

```liquid
{% content %}
    Main page content
{% endcontent %}
```

---

## Capture

Capture output to a variable:

```liquid
{% capture .output %}
    Captured content: {%.value%}
{% endcapture %}

{% .output %}  {# Use captured content #}
```

---

## Break and Continue

### Break

Exit loop early:

```liquid
{% for .item in .items %}
    {% if .item.last %}
        {% break %}
    {% end %}
    {%.item.name%}
{% end %}
```

### Continue

Skip to next iteration:

```liquid
{% for .item in .items %}
    {% if .item.skip %}
        {% continue %}
    {% end %}
    {%.item.name%}
{% end %}
```

---

## Raw Block

Output literal content without processing:

```liquid
{% raw %}
    This {%.will_not%} be processed
    {% if .ignored %}...
{% endraw %}
```

---

## Operators

### Comparison

```liquid
{% if .a == .b %}
{% if .a != .b %}
{% if .a > .b %}
{% if .a >= .b %}
{% if .a < .b %}
{% if .a <= .b %}
```

### Logical

```liquid
{% if .a and .b %}
{% if .a or .b %}
{% if not .a %}
```

### Containment

```liquid
{% if .item in .collection %}
{% if .item not in .collection %}
```
