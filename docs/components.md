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
