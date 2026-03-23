# Getting Started with Wisp

Wisp is a secure, high-performance HTML templating engine for Go with a Liquid-inspired syntax.

## Installation

```bash
go get github.com/anomalyco/wisp
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/anomalyco/wisp/pkg/engine"
)

func main() {
    e := engine.New()
    
    template := `Hello, {%.name%}!`
    data := map[string]interface{}{"name": "World"}
    
    result, err := e.RenderString(template, data)
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // Output: Hello, World!
}
```

## Basic Syntax

### Variables

Access variables with a leading dot:

```liquid
{% .name %}
{% .user.email %}
{% .items[0] %}
{% .data[key] %}
```

### Conditionals

```liquid
{% if .show %}
    Content to show
{% elsif .alt %}
    Alternative content
{% else %}
    Default content
{% end %}
```

### Loops

```liquid
{% for .item in .items %}
    <li>{%.item.name%}</li>
{% end %}
```

### Filters

Apply filters with the pipe operator:

```liquid
{% .name | upcase %}
{% .price | times 1.1 %}
{% .date | date "2006-01-02" %}
```

## Rendering Templates

### From String

```go
result, err := e.RenderString(template, data)
```

### From File

```go
// Register a template store
e := engine.New()
e.SetStore(engine.NewFileStore("./templates"))

result, err := e.RenderFile("index.html", data)
```

### Register Templates Manually

```go
e := engine.New()
e.RegisterTemplate("header", `Header: {%.title%}`)

result, err := e.RenderFile("header", data)
```

## Filters

### String Filters

```liquid
{% .name | upcase %}           {# HELLO #}
{% .name | downcase %}         {# hello #}
{% .text | truncate 50 %}      {# Truncate to 50 chars #}
{% .name | replace "old" "new" %}
{% .tags | join ", " %}
```

### Numeric Filters

```liquid
{% .price | times 1.1 %}       {# Multiply #}
{% .price | plus 5 %}          {# Add #}
{% .price | minus 2 %}         {# Subtract #}
{% .value | abs %}             {# Absolute value #}
{% .value | round %}           {# Round to nearest #}
```

### Array Filters

```liquid
{% .items | first %}
{% .items | last %}
{% .items | size %}
{% .items | reverse %}
{% .items | sort %}
{% .items | uniq %}
```

### Escape Filters

```liquid
{% .html | escape %}           {# HTML escape #}
{% .html | escape_once %}      {# Escape only HTML entities #}
{% .raw | raw %}               {# Mark as safe (no escaping) #}
```

## Template Composition

### Include

Include another template (shares scope):

```liquid
{% include "partials/header" %}
{% include "sidebar" .user %}
```

### Render

Render with isolated scope:

```liquid
{% render "widget" .data %}
```

### Component

Props-based component system:

```liquid
{% component "Button" .buttonProps %}
```

## Layout System

### Extends

```liquid
{# child.html #}
{% extends "layouts/base" %}

{% block content %}
    Page content here
{% endblock %}
```

### Base Layout

```liquid
{# layouts/base.html #}
<html>
<head>{% block title %}Default Title{% endblock %}</head>
<body>
    {% block content %}{% endblock %}
</body>
</html>
```

## Security

### Auto-Escaping

HTML auto-escaping is enabled by default:

```go
e := engine.New()  // Auto-escape enabled by default
e.SetAutoEscape(false)  // Disable if needed
```

### Safe Strings

Mark content as safe (no escaping):

```go
import "github.com/anomalyco/wisp/pkg/engine"

safe := engine.SafeString("<b>Bold</b>")
result, _ := e.RenderString(`{%.content%}`, map[string]interface{}{"content": safe})
// Output: <b>Bold</b> (not escaped)
```

### Resource Limits

```go
e := engine.New()
e.SetMaxIterations(10000)  // Prevent infinite loops
```

## Custom Filters

```go
e := engine.New()
e.RegisterFilter("shout", func(input interface{}) string {
    return fmt.Sprintf("!!! %s !!!", toString(input))
})
```

```liquid
{% .message | shout %}
{# Output: !!! Hello World !!! #}
```

## Error Handling

```go
result, err := e.RenderString(template, data)
if err != nil {
    if parseErrs, ok := err.([]error); ok {
        for _, e := range parseErrs {
            fmt.Println("Parse error:", e)
        }
    }
}
```

## CLI Tool

```bash
# Render a template
echo '{"name": "World"}' | wisp render 'Hello, {%.name%}!'

# Validate template syntax
wisp validate '{% if .show %}{%.content%}{% end %}'

# Show version
wisp version
```

## Next Steps

- [Template Syntax Reference](./syntax-reference.md)
- [API Documentation](./api.md)
- [Security Best Practices](./security.md)
