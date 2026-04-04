# Example CSS Refactor — Global Stylesheets with Brand Colors

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace inline styles in all example templates with global CSS files using CSS variables and utility classes, matching the brand colors from the logo (`#2E6740` green, `#251917` espresso, `#EEEBE3` cream).

**Architecture:** Each example (blog, store, docs) gets its own self-contained CSS file with shared CSS custom properties and common component classes. Templates reference these via the existing `{% asset %}` tag. Email templates keep inline styles (email client compatibility) but get brand color updates. Each Go server gets a static file handler to actually serve the CSS.

**Tech Stack:** CSS custom properties, Go `http.FileServer`, Grove `{% asset %}` tag

---

## File Structure

### Blog (`examples/blog/`)
- **Create:** `static/style.css` — CSS variables + component classes
- **Modify:** `main.go` — add static file server route
- **Modify:** `templates/base.grov` — replace body inline styles with classes
- **Modify:** `templates/index.grov` — replace inline styles with classes
- **Modify:** `templates/post.grov` — replace common inline styles with classes
- **Modify:** `templates/pages/styleguide.grov` — replace inline styles with classes
- **Modify:** `templates/components/nav.grov` — replace inline styles with classes
- **Modify:** `templates/components/footer.grov` — replace inline styles with classes
- **Modify:** `templates/components/card.grov` — replace inline styles with classes
- **Modify:** `templates/components/button.grov` — replace inline styles with classes
- **Modify:** `templates/components/tag.grov` — replace inline styles with classes
- **Modify:** `templates/components/alert.grov` — replace inline styles with classes

### Store (`examples/store/`)
- **Create:** `static/style.css` — CSS variables + component classes
- **Modify:** `main.go` — add static file server route + update placeholder URLs
- **Modify:** `templates/base.grov` — replace inline nav/footer/body with classes
- **Modify:** `templates/index.grov` — replace inline styles with classes
- **Modify:** `templates/product.grov` — replace common inline styles with classes, keep unique layout inline
- **Modify:** `templates/cart.grov` — replace common inline styles with classes, keep unique table layout inline
- **Modify:** `templates/components/product-card.grov` — replace inline styles with classes
- **Modify:** `templates/macros/pricing.grov` — update colors to brand

### Docs (`examples/docs/`)
- **Create:** `static/docs.css` — CSS variables + component classes + docs-specific classes
- **Modify:** `main.go` — add static file server route
- **Modify:** `templates/base.grov` — replace inline styles with classes
- **Modify:** `templates/docs-layout.grov` — replace inline styles with classes
- **Modify:** `templates/pages/_default.grov` — replace inline styles with classes
- **Modify:** `templates/pages/variables-and-filters.grov` — replace inline styles with classes
- **Modify:** `templates/partials/sidebar.grov` — replace inline styles with classes
- **Modify:** `templates/macros/admonitions.grov` — replace inline styles with classes

### Email (`examples/email/`)
- **Modify:** `templates/base-email.grov` — update hardcoded colors to brand (keep inline)
- **Modify:** `templates/index.grov` — update colors to brand (keep inline)
- **Modify:** `templates/welcome.grov` — update colors to brand (keep inline)
- **Modify:** `templates/order-confirmation.grov` — update colors to brand (keep inline)
- **Modify:** `templates/password-reset.grov` — update colors to brand (keep inline)
- **Modify:** `templates/helpers.grov` — update default button color to brand

---

### Task 1: Create blog CSS file and static file server

**Files:**
- Create: `examples/blog/static/style.css`
- Modify: `examples/blog/main.go`

- [ ] **Step 1: Create the static directory and CSS file**

Create `examples/blog/static/style.css`:

```css
/* Grove Blog — Global Stylesheet */

:root {
  /* Brand colors (from logo) */
  --color-primary: #2E6740;
  --color-primary-hover: #245533;
  --color-dark: #251917;
  --color-cream: #EEEBE3;
  --color-page-bg: #F7F5F0;
  --color-text: #3D2E2A;
  --color-text-muted: #7A6B66;
  --color-border: #D9D3CB;
  --color-green-light: #E8F0EA;
  --color-cream-dark: #DDD8CE;

  /* Alert colors */
  --color-info-bg: #E8F0EA;
  --color-info-text: #2E6740;
  --color-info-border: #2E6740;
  --color-warning-bg: #FFF3CD;
  --color-warning-text: #6B5210;
  --color-warning-border: #E6C547;
  --color-error-bg: #F8D7DA;
  --color-error-text: #6B1D24;
  --color-error-border: #E8A0A7;
  --color-success-bg: #D4EDDA;
  --color-success-text: #1B5E28;
  --color-success-border: #A3D4AE;

  /* Tag colors */
  --color-tag-green-bg: #E8F0EA;
  --color-tag-green-text: #2E6740;
  --color-tag-brown-bg: #F0EAE4;
  --color-tag-brown-text: #5C3D2E;
  --color-tag-red-bg: #FEE2E2;
  --color-tag-red-text: #991B1B;
  --color-tag-purple-bg: #EDE9FE;
  --color-tag-purple-text: #5B21B6;
  --color-tag-orange-bg: #FFEDD5;
  --color-tag-orange-text: #9A3412;
  --color-tag-gray-bg: #EDEBE7;
  --color-tag-gray-text: #4A3F3A;

  /* Spacing */
  --radius-sm: 4px;
  --radius-md: 6px;
  --radius-lg: 8px;
  --radius-pill: 999px;
}

/* ── Base ──────────────────────────────── */

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  color: var(--color-text);
  background: var(--color-page-bg);
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

a { color: var(--color-primary); }

/* ── Layout ────────────────────────────── */

.container {
  max-width: 960px;
  width: 100%;
  margin: 0 auto;
  padding: 2rem 1rem;
  flex: 1;
}

.grid { display: grid; gap: 1.5rem; }

.flex-row {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

/* ── Nav ───────────────────────────────── */

.nav {
  background: var(--color-dark);
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-brand {
  color: var(--color-primary);
  font-size: 1.4rem;
  font-weight: bold;
  text-decoration: none;
}

.nav-links {
  display: flex;
  gap: 1.5rem;
  align-items: center;
}

.nav-link {
  color: var(--color-cream);
  text-decoration: none;
}

/* ── Footer ────────────────────────────── */

.footer {
  background: var(--color-dark);
  color: var(--color-text-muted);
  padding: 2rem;
  text-align: center;
  margin-top: 3rem;
}

.footer p { margin: 0; }

.footer a { color: var(--color-primary); }

/* ── Card ──────────────────────────────── */

.card {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  background: var(--color-cream);
  transition: box-shadow 0.2s;
}

.card:hover {
  box-shadow: 0 2px 8px rgba(37, 25, 23, 0.1);
}

.card h2 { margin: 0 0 0.5rem; }

.card h2 a {
  color: var(--color-dark);
  text-decoration: none;
}

.card time {
  color: var(--color-text-muted);
  font-size: 0.85rem;
}

.card p {
  color: var(--color-text-muted);
  margin: 0.75rem 0;
}

/* ── Button ────────────────────────────── */

.btn {
  display: inline-block;
  padding: 0.5rem 1.25rem;
  border: 2px solid;
  border-radius: var(--radius-md);
  text-decoration: none;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
}

.btn-primary {
  background: var(--color-primary);
  color: var(--color-cream);
  border-color: var(--color-primary);
}

.btn-secondary {
  background: var(--color-dark);
  color: var(--color-cream);
  border-color: var(--color-dark);
}

.btn-outline {
  background: transparent;
  color: var(--color-primary);
  border-color: var(--color-primary);
}

.btn-default {
  background: var(--color-text-muted);
  color: var(--color-cream);
  border-color: var(--color-text-muted);
}

/* ── Tag ───────────────────────────────── */

.tag {
  display: inline-block;
  padding: 0.2rem 0.6rem;
  border-radius: var(--radius-pill);
  font-size: 0.75rem;
  font-weight: 600;
}

.tag-green  { background: var(--color-tag-green-bg);  color: var(--color-tag-green-text); }
.tag-brown  { background: var(--color-tag-brown-bg);  color: var(--color-tag-brown-text); }
.tag-red    { background: var(--color-tag-red-bg);    color: var(--color-tag-red-text); }
.tag-purple { background: var(--color-tag-purple-bg); color: var(--color-tag-purple-text); }
.tag-orange { background: var(--color-tag-orange-bg); color: var(--color-tag-orange-text); }
.tag-gray   { background: var(--color-tag-gray-bg);   color: var(--color-tag-gray-text); }

/* ── Alert ─────────────────────────────── */

.alert {
  padding: 1rem 1.25rem;
  border: 1px solid;
  border-radius: var(--radius-md);
  display: flex;
  gap: 0.75rem;
  align-items: flex-start;
}

.alert-icon {
  font-weight: bold;
  font-size: 1.1rem;
}

.alert-info    { background: var(--color-info-bg);    border-color: var(--color-info-border);    color: var(--color-info-text); }
.alert-warning { background: var(--color-warning-bg); border-color: var(--color-warning-border); color: var(--color-warning-text); }
.alert-error   { background: var(--color-error-bg);   border-color: var(--color-error-border);   color: var(--color-error-text); }
.alert-success { background: var(--color-success-bg); border-color: var(--color-success-border); color: var(--color-success-text); }

/* ── Article ───────────────────────────── */

.article {
  background: var(--color-cream);
  border-radius: var(--radius-lg);
  padding: 2rem;
}

.article-body { line-height: 1.7; }

.article time {
  color: var(--color-text-muted);
  font-size: 0.9rem;
}

/* ── Styleguide ────────────────────────── */

.section-heading {
  border-bottom: 2px solid var(--color-primary);
  padding-bottom: 0.5rem;
}

.section-desc { color: var(--color-text-muted); }

/* ── Utilities ─────────────────────────── */

.text-muted { color: var(--color-text-muted); }
.text-small { font-size: 0.85rem; }
.text-accent { color: var(--color-primary); }
.font-bold { font-weight: 600; }
.no-decor { text-decoration: none; }
```

- [ ] **Step 2: Add static file server to blog main.go**

In `examples/blog/main.go`, add a static file server route after the existing routes. Add this import and route:

Add `"os"` to the imports, then after the `r.Get("/styleguide", ...)` line, add:

```go
staticDir := filepath.Join(filepath.Dir(thisFile), "static")
r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
```

- [ ] **Step 3: Verify the blog builds**

Run: `cd examples/blog && go build ./...`
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add examples/blog/static/style.css examples/blog/main.go
git commit -m "feat(blog): add global CSS file with brand variables and static file server"
```

---

### Task 2: Update blog templates to use CSS classes

**Files:**
- Modify: `examples/blog/templates/base.grov`
- Modify: `examples/blog/templates/components/nav.grov`
- Modify: `examples/blog/templates/components/footer.grov`
- Modify: `examples/blog/templates/components/card.grov`
- Modify: `examples/blog/templates/components/button.grov`
- Modify: `examples/blog/templates/components/tag.grov`
- Modify: `examples/blog/templates/components/alert.grov`
- Modify: `examples/blog/templates/index.grov`
- Modify: `examples/blog/templates/post.grov`
- Modify: `examples/blog/templates/pages/styleguide.grov`

- [ ] **Step 1: Update base.grov**

Replace the entire file with:

```html
{% asset "/static/style.css" type="stylesheet" priority=10 %}
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{% block title %}Grove Blog{% endblock %}</title>
  <!-- HEAD_ASSETS -->
  <!-- HEAD_META -->
  <!-- HEAD_HOISTED -->
</head>
<body>
  {% component "components/nav.grov" site_name=site_name %}{% endcomponent %}
  <main class="container">
    {% block content %}{% endblock %}
  </main>
  {% component "components/footer.grov" year=current_year %}{% endcomponent %}
  <!-- FOOT_ASSETS -->
</body>
</html>
```

- [ ] **Step 2: Update nav.grov**

Replace the entire file with:

```html
{% props site_name %}
<nav class="nav">
  <a href="/" class="nav-brand">{{ site_name }}</a>
  <div class="nav-links">
    <a href="/" class="nav-link">Home</a>
    <a href="/styleguide" class="nav-link">Components</a>
    {% slot %}{% endslot %}
  </div>
</nav>
```

- [ ] **Step 3: Update footer.grov**

Replace the entire file with:

```html
{% props year %}
<footer class="footer">
  <p>© {{ year }} Grove Blog. Built with the <a href="#">Grove</a> template engine.</p>
</footer>
```

- [ ] **Step 4: Update card.grov**

Replace the entire file with:

```html
{% props title, summary, href="#", date="" %}
<article class="card">
  <h2>
    <a href="{{ href }}">{{ title }}</a>
  </h2>
  {% if date %}<time>{{ date }}</time>{% endif %}
  <p>{{ summary | truncate(120) }}</p>
  <div class="flex-row">
    {% slot "tags" %}{% endslot %}
  </div>
</article>
```

- [ ] **Step 5: Update button.grov**

Replace the entire file with:

```html
{% props label, href="/", variant="primary" %}
{% if variant == "primary" %}
  {% set cls = "btn btn-primary" %}
{% elif variant == "secondary" %}
  {% set cls = "btn btn-secondary" %}
{% elif variant == "outline" %}
  {% set cls = "btn btn-outline" %}
{% else %}
  {% set cls = "btn btn-default" %}
{% endif %}

<a href="{{ href }}" class="{{ cls }}">{{ label }}</a>
```

- [ ] **Step 6: Update tag.grov**

Replace the entire file with:

```html
{% props label, color="gray" %}
{% if color == "green" %}
  {% set cls = "tag tag-green" %}
{% elif color == "red" %}
  {% set cls = "tag tag-red" %}
{% elif color == "purple" %}
  {% set cls = "tag tag-purple" %}
{% elif color == "orange" %}
  {% set cls = "tag tag-orange" %}
{% elif color == "blue" %}
  {% set cls = "tag tag-brown" %}
{% else %}
  {% set cls = "tag tag-gray" %}
{% endif %}
<span class="{{ cls }}">{{ label }}</span>
```

Note: The old "blue" color maps to `tag-brown` (the brand's warm brown replaces cold blue).

- [ ] **Step 7: Update alert.grov**

Replace the entire file with:

```html
{% props type="info" %}
{% let %}
  icon = "ℹ"
  cls = "alert alert-info"

  if type == "warning"
    icon = "⚠"
    cls = "alert alert-warning"
  elif type == "error"
    icon = "✕"
    cls = "alert alert-error"
  elif type == "success"
    icon = "✓"
    cls = "alert alert-success"
  end
{% endlet %}
<div class="{{ cls }}">
  <span class="alert-icon">{{ icon }}</span>
  <div>{% slot %}{% endslot %}</div>
</div>
```

- [ ] **Step 8: Update index.grov**

Replace the entire file with:

```html
{% extends "base.grov" %}

{% block title %}Home — Grove Blog{% endblock %}

{% block content %}
{% meta name="description" content="A blog built with the Grove template engine" %}

<h1 style="margin: 0 0 1.5rem;">Latest Posts</h1>
<div class="grid">
  {% for post in posts %}
    {% component "components/card.grov" title=post.title summary=post.summary href="/post/" ~ post.slug date=post.date %}
      {% fill "tags" %}
        {% for tag in post.tags %}
          {% component "components/tag.grov" label=tag.name color=tag.color %}{% endcomponent %}
        {% endfor %}
      {% endfill %}
    {% endcomponent %}
  {% endfor %}
</div>
{% endblock %}
```

- [ ] **Step 9: Update post.grov**

Replace the entire file with:

```html
{% extends "base.grov" %}

{% block title %}{{ post.title }} — Grove Blog{% endblock %}

{% block content %}
{% meta name="description" content="A blog post on Grove Blog" %}
{% meta property="og:title" content="Grove Blog Post" %}
{% asset "/static/highlight.css" type="stylesheet" %}

{% if post.draft %}
  {% component "components/alert.grov" type="warning" %}
    This post is a <strong>draft</strong> and is not yet published.
  {% endcomponent %}
{% endif %}

<article class="article" style="margin-top: 1rem;">
  <header style="margin-bottom: 1.5rem;">
    <h1 style="margin: 0 0 0.5rem;">{{ post.title }}</h1>
    <time>{{ post.date }}</time>
    <div class="flex-row" style="margin-top: 0.75rem;">
      {% for tag in post.tags %}
        {% component "components/tag.grov" label=tag.name color=tag.color %}{% endcomponent %}
      {% endfor %}
    </div>
  </header>
  <div class="article-body">
    {{ post.body | nl2br | safe }}
  </div>
</article>

<div style="margin-top: 1.5rem;">
  {% component "components/button.grov" label="← Back to posts" href="/" variant="secondary" %}{% endcomponent %}
</div>
{% endblock %}
```

- [ ] **Step 10: Update styleguide.grov**

Replace the entire file with:

```html
{% extends "base.grov" %}

{% block title %}Component Library — Grove Blog{% endblock %}

{% block content %}
{% meta name="description" content="A showcase of all reusable Grove components" %}

<h1 style="margin: 0 0 0.5rem;">Component Library</h1>
<p class="section-desc" style="margin: 0 0 2rem;">A showcase of every reusable component in the Grove Blog.</p>

{# ── Buttons ─────────────────────────────────────────── #}
<section style="margin-bottom: 3rem;">
  <h2 class="section-heading">Buttons</h2>
  <p class="section-desc">The <code>button.grov</code> component accepts <code>label</code>, <code>href</code>, and <code>variant</code> props.</p>
  <div style="display: flex; gap: 1rem; flex-wrap: wrap; margin-top: 1rem;">
    {% component "components/button.grov" label="Primary" href="#" variant="primary" %}{% endcomponent %}
    {% component "components/button.grov" label="Secondary" href="#" variant="secondary" %}{% endcomponent %}
    {% component "components/button.grov" label="Outline" href="#" variant="outline" %}{% endcomponent %}
    {% component "components/button.grov" label="Default" href="#" %}{% endcomponent %}
  </div>
</section>

{# ── Alerts ──────────────────────────────────────────── #}
<section style="margin-bottom: 3rem;">
  <h2 class="section-heading">Alerts</h2>
  <p class="section-desc">The <code>alert.grov</code> component accepts a <code>type</code> prop and a default slot for the message body.</p>
  <div class="grid" style="margin-top: 1rem;">
    {% component "components/alert.grov" type="info" %}
      This is an <strong>info</strong> alert — useful for general messages.
    {% endcomponent %}
    {% component "components/alert.grov" type="success" %}
      This is a <strong>success</strong> alert — something went well!
    {% endcomponent %}
    {% component "components/alert.grov" type="warning" %}
      This is a <strong>warning</strong> alert — proceed with caution.
    {% endcomponent %}
    {% component "components/alert.grov" type="error" %}
      This is an <strong>error</strong> alert — something went wrong.
    {% endcomponent %}
  </div>
</section>

{# ── Tags ────────────────────────────────────────────── #}
<section style="margin-bottom: 3rem;">
  <h2 class="section-heading">Tags</h2>
  <p class="section-desc">The <code>tag.grov</code> component accepts <code>label</code> and <code>color</code> props.</p>
  <div class="flex-row" style="margin-top: 1rem;">
    {% for color in tag_colors %}
      {% component "components/tag.grov" label=color | title color=color %}{% endcomponent %}
    {% endfor %}
    {% component "components/tag.grov" label="Default" %}{% endcomponent %}
  </div>
</section>

{# ── Cards ───────────────────────────────────────────── #}
<section style="margin-bottom: 3rem;">
  <h2 class="section-heading">Cards</h2>
  <p class="section-desc">The <code>card.grov</code> component accepts <code>title</code>, <code>summary</code>, <code>href</code>, <code>date</code> props and a <code>tags</code> slot.</p>
  <div class="grid" style="margin-top: 1rem;">
    {% component "components/card.grov" title="Card With Tags" summary="This card demonstrates the tags slot, which lets the caller inject tag components into the card footer." href="#" date="April 1, 2026" %}
      {% fill "tags" %}
        {% component "components/tag.grov" label="Grove" color="purple" %}{% endcomponent %}
        {% component "components/tag.grov" label="Example" color="blue" %}{% endcomponent %}
      {% endfill %}
    {% endcomponent %}
    {% component "components/card.grov" title="Card Without Tags" summary="This card has no tags filled, showing the empty default slot. The date prop is also omitted." href="#" %}{% endcomponent %}
    {% component "components/card.grov" title="Card With Long Summary" summary="This card has a very long summary to demonstrate the truncate filter in action. The summary text will be cut off at 120 characters and an ellipsis will be appended to indicate there is more content." href="#" date="March 15, 2026" %}{% endcomponent %}
  </div>
</section>

{# ── Navigation ──────────────────────────────────────── #}
<section style="margin-bottom: 3rem;">
  <h2 class="section-heading">Navigation</h2>
  <p class="section-desc">The <code>nav.grov</code> component accepts <code>site_name</code> and an optional default slot for extra links. The nav bar at the top of this page is the default usage. Here is one with extra slot content:</p>
  <div style="margin-top: 1rem; border-radius: 8px; overflow: hidden;">
    {% component "components/nav.grov" site_name="Custom Nav" %}
      <a href="#" class="text-accent no-decor font-bold">Extra Link</a>
    {% endcomponent %}
  </div>
</section>

{# ── Footer ──────────────────────────────────────────── #}
<section style="margin-bottom: 3rem;">
  <h2 class="section-heading">Footer</h2>
  <p class="section-desc">The <code>footer.grov</code> component accepts a <code>year</code> prop. The footer at the bottom of this page is the default usage.</p>
  <div style="margin-top: 1rem; border-radius: 8px; overflow: hidden;">
    {% component "components/footer.grov" year=current_year %}{% endcomponent %}
  </div>
</section>
{% endblock %}
```

- [ ] **Step 11: Run build check**

Run: `go build ./examples/blog/...`
Expected: no errors

- [ ] **Step 12: Commit**

```bash
git add examples/blog/templates/
git commit -m "feat(blog): replace inline styles with CSS classes across all templates"
```

---

### Task 3: Create store CSS file and static file server

**Files:**
- Create: `examples/store/static/style.css`
- Modify: `examples/store/main.go`

- [ ] **Step 1: Create the static directory and CSS file**

Create `examples/store/static/style.css`:

```css
/* Grove Store — Global Stylesheet */

:root {
  /* Brand colors (from logo) */
  --color-primary: #2E6740;
  --color-primary-hover: #245533;
  --color-dark: #251917;
  --color-cream: #EEEBE3;
  --color-page-bg: #F7F5F0;
  --color-text: #3D2E2A;
  --color-text-muted: #7A6B66;
  --color-border: #D9D3CB;
  --color-green-light: #E8F0EA;

  /* Tag colors */
  --color-tag-red-bg: #FEE2E2;
  --color-tag-red-text: #991B1B;
  --color-tag-gray-bg: #EDEBE7;
  --color-tag-gray-text: #4A3F3A;

  /* Spacing */
  --radius-sm: 4px;
  --radius-md: 6px;
  --radius-lg: 8px;
  --radius-pill: 999px;
}

/* ── Base ──────────────────────────────── */

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  color: var(--color-text);
  background: var(--color-page-bg);
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

a { color: var(--color-primary); }

/* ── Layout ────────────────────────────── */

.container {
  max-width: 1080px;
  width: 100%;
  margin: 0 auto;
  padding: 2rem 1rem;
  flex: 1;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 1.5rem;
}

/* ── Nav ───────────────────────────────── */

.nav {
  background: var(--color-dark);
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-brand {
  color: var(--color-primary);
  font-size: 1.4rem;
  font-weight: bold;
  text-decoration: none;
}

.nav-links {
  display: flex;
  gap: 1.5rem;
  align-items: center;
}

.nav-link {
  color: var(--color-cream);
  text-decoration: none;
}

/* ── Footer ────────────────────────────── */

.footer {
  background: var(--color-dark);
  color: var(--color-text-muted);
  padding: 2rem;
  text-align: center;
  margin-top: 3rem;
}

.footer p { margin: 0; }
.footer a { color: var(--color-primary); }

/* ── Card ──────────────────────────────── */

.card {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  background: var(--color-cream);
  transition: box-shadow 0.2s;
}

.card:hover {
  box-shadow: 0 2px 8px rgba(37, 25, 23, 0.1);
}

.card img {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.card-body { padding: 1rem; }

.card-body h3 { margin: 0 0 0.5rem; }

.card-body h3 a {
  color: var(--color-dark);
  text-decoration: none;
}

/* ── Badge ─────────────────────────────── */

.badge {
  display: inline-block;
  margin-top: 0.5rem;
  padding: 0.2rem 0.6rem;
  border-radius: var(--radius-pill);
  font-size: 0.75rem;
  font-weight: 600;
}

.badge-sale { background: var(--color-tag-red-bg); color: var(--color-tag-red-text); }
.badge-stock { background: var(--color-tag-gray-bg); color: var(--color-tag-gray-text); }

/* ── Button ────────────────────────────── */

.btn {
  display: inline-block;
  padding: 0.75rem 2rem;
  background: var(--color-primary);
  color: var(--color-cream);
  border: none;
  border-radius: var(--radius-md);
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  text-decoration: none;
}

.btn-sm {
  padding: 0.5rem 1.25rem;
  font-size: 0.9rem;
  border: 2px solid var(--color-dark);
  background: var(--color-dark);
}

/* ── Breadcrumb ────────────────────────── */

.breadcrumb {
  font-size: 0.9rem;
  color: var(--color-text-muted);
  margin-bottom: 1.5rem;
}

.breadcrumb a {
  color: var(--color-primary);
  text-decoration: none;
}

/* ── Article / Content ─────────────────── */

.panel {
  background: var(--color-cream);
  border-radius: var(--radius-lg);
  padding: 2rem;
}

/* ── Table ─────────────────────────────── */

.table-wrap {
  background: var(--color-cream);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.table {
  width: 100%;
  border-collapse: collapse;
}

.table thead { background: var(--color-page-bg); text-align: left; }
.table th { padding: 1rem; }
.table td { padding: 1rem; }
.table tbody tr { border-top: 1px solid var(--color-border); }

.table a {
  color: var(--color-dark);
  text-decoration: none;
  font-weight: 600;
}

/* ── Cart Summary ──────────────────────── */

.cart-summary {
  margin-top: 1.5rem;
  background: var(--color-cream);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  max-width: 360px;
  margin-left: auto;
}

.cart-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.cart-divider {
  border: none;
  border-top: 1px solid var(--color-border);
  margin: 0.75rem 0;
}

/* ── Pricing ───────────────────────────── */

.price-strike { text-decoration: line-through; color: var(--color-text-muted); }
.price-sale { color: var(--color-primary); font-weight: bold; }
.price-regular { font-weight: bold; }

.stars { color: #D4A843; letter-spacing: 2px; }
.review-count { color: var(--color-text-muted); font-size: 0.85rem; }

/* ── Utilities ─────────────────────────── */

.text-muted { color: var(--color-text-muted); }
.text-accent { color: var(--color-primary); }
.text-error { color: var(--color-tag-red-text); }
.text-success { color: #1B5E28; }
.font-bold { font-weight: 600; }
.no-decor { text-decoration: none; }
```

- [ ] **Step 2: Add static file server to store main.go and update placeholder image URLs**

In `examples/store/main.go`:

1. After the `r.Get("/cart", ...)` line, add:

```go
staticDir := filepath.Join(filepath.Dir(thisFile), "static")
r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
```

2. Update the placeholder image URLs to use brand colors. Replace these `ImageURL` values:

```go
// Wireless Headphones
ImageURL: "https://placehold.co/400x300/251917/2E6740?text=Headphones",
// Mechanical Keyboard
ImageURL: "https://placehold.co/400x300/2E6740/EEEBE3?text=Keyboard",
// Running Shoes
ImageURL: "https://placehold.co/400x300/251917/2E6740?text=Shoes",
// Desk Lamp
ImageURL: "https://placehold.co/400x300/3D2E2A/EEEBE3?text=Lamp",
```

- [ ] **Step 3: Verify the store builds**

Run: `go build ./examples/store/...`
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add examples/store/static/style.css examples/store/main.go
git commit -m "feat(store): add global CSS file with brand variables and static file server"
```

---

### Task 4: Update store templates to use CSS classes

**Files:**
- Modify: `examples/store/templates/base.grov`
- Modify: `examples/store/templates/index.grov`
- Modify: `examples/store/templates/product.grov`
- Modify: `examples/store/templates/cart.grov`
- Modify: `examples/store/templates/components/product-card.grov`
- Modify: `examples/store/templates/macros/pricing.grov`

- [ ] **Step 1: Update base.grov**

Replace the entire file with:

```html
{% asset "/static/style.css" type="stylesheet" priority=10 %}
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{% block title %}Grove Store{% endblock %}</title>
  <!-- HEAD_ASSETS -->
  <!-- HEAD_META -->
</head>
<body>
  <nav class="nav">
    <a href="/" class="nav-brand">{{ site_name }}</a>
    <div class="nav-links">
      <a href="/" class="nav-link">Products</a>
      <a href="/cart" class="nav-link">Cart</a>
    </div>
  </nav>
  <main class="container">
    {% block content %}{% endblock %}
  </main>
  <footer class="footer">
    <p>© {{ current_year }} Grove Store. Powered by the Grove template engine.</p>
  </footer>
  <!-- FOOT_ASSETS -->
</body>
</html>
```

- [ ] **Step 2: Update index.grov**

Replace the entire file with:

```html
{% extends "base.grov" %}
{% import "macros/pricing.grov" as pricing %}

{% block title %}Products — Grove Store{% endblock %}

{% block content %}
{% meta name="description" content="Browse our product catalog" %}

<h1 style="margin: 0 0 1.5rem;">Products</h1>
<div class="product-grid">
  {% for product in products %}
    {% capture price_html %}
      {{ pricing.price(product.price, product.sale_price) }}
    {% endcapture %}
    {% component "components/product-card.grov" name=product.name slug=product.slug image_url=product.image_url price_display=price_html %}
      {% fill "badge" %}
        {% if product.on_sale %}
          <span class="badge badge-sale">Sale!</span>
        {% endif %}
        {% if not product.in_stock %}
          <span class="badge badge-stock">Out of Stock</span>
        {% endif %}
      {% endfill %}
    {% endcomponent %}
  {% endfor %}
</div>
{% endblock %}
```

- [ ] **Step 3: Update product-card.grov**

Replace the entire file with:

```html
{% props name, slug, image_url, price_display %}
<div class="card">
  <img src="{{ image_url }}" alt="{{ name }}">
  <div class="card-body">
    <h3>
      <a href="/product/{{ slug }}">{{ name }}</a>
    </h3>
    <div>{{ price_display | safe }}</div>
    {% slot "badge" %}{% endslot %}
  </div>
</div>
```

- [ ] **Step 4: Update pricing.grov**

Replace the entire file with:

```html
{% macro price(amount, sale_amount) %}
  {% if sale_amount > 0 %}
    <span class="price-strike">{{ amount | currency }}</span>
    <span class="price-sale">{{ sale_amount | currency }}</span>
  {% else %}
    <span class="price-regular">{{ amount | currency }}</span>
  {% endif %}
{% endmacro %}

{% macro star_rating(rating, count) %}
  {% set full = rating | floor %}
  {% set half = rating - full >= 0.5 ? 1 : 0 %}
  <span class="stars">
    {% for i in range(1, full) %}★{% endfor %}{% if half %}½{% endif %}
  </span>
  <span class="review-count">({{ count }})</span>
{% endmacro %}
```

- [ ] **Step 5: Update product.grov**

Replace the entire file with:

```html
{% extends "base.grov" %}
{% import "macros/pricing.grov" as pricing %}

{% block title %}{{ product.name }} — Grove Store{% endblock %}

{% block content %}
{% meta name="description" content=product.description | truncate(160) %}
{% meta property="og:title" content=product.name %}

<nav class="breadcrumb">
  {% for crumb in breadcrumbs %}
    {% if crumb.href %}
      <a href="{{ crumb.href }}">{{ crumb.label }}</a> /
    {% else %}
      {{ crumb.label }}
    {% endif %}
  {% endfor %}
</nav>

<div class="panel" style="display: grid; grid-template-columns: 1fr 1fr; gap: 2rem;">
  <img src="{{ product.image_url }}" alt="{{ product.name }}" style="width: 100%; border-radius: 8px;">
  <div>
    <h1 style="margin: 0 0 0.5rem;">{{ product.name }}</h1>
    <div style="margin-bottom: 1rem;">
      {{ pricing.star_rating(product.rating, product.review_count) }}
    </div>
    <div style="font-size: 1.5rem; margin-bottom: 1rem;">
      {{ pricing.price(product.price, product.sale_price) }}
    </div>

    {% if product.on_sale %}
      {% let %}
        savings = product.price - product.sale_price
      {% endlet %}
      <p class="text-success font-bold">You save {{ savings | currency }}!</p>
    {% endif %}

    <p class="text-muted" style="line-height: 1.6;">{{ product.description }}</p>

    {% if product.colors | length > 0 %}
      <div style="margin: 1rem 0;">
        <strong>Colors:</strong>
        {% for color in product.colors %}
          <span style="display: inline-block; padding: 0.2rem 0.6rem; border: 1px solid var(--color-border); border-radius: 4px; margin: 0.25rem; font-size: 0.85rem;">{{ color }}</span>
        {% endfor %}
      </div>
    {% endif %}

    {% if product.in_stock %}
      <div style="margin: 1rem 0;">
        <label class="font-bold">Quantity:</label>
        <select style="padding: 0.4rem; border-radius: 4px; border: 1px solid var(--color-border);">
          {% for n in range(1, 10) %}
            <option value="{{ n }}">{{ n }}</option>
          {% endfor %}
        </select>
      </div>
      <button class="btn">Add to Cart</button>
    {% else %}
      <p class="text-error font-bold">Out of stock</p>
    {% endif %}
  </div>
</div>
{% endblock %}
```

- [ ] **Step 6: Update cart.grov**

Replace the entire file with:

```html
{% extends "base.grov" %}
{% import "macros/pricing.grov" as pricing %}

{% block title %}Cart — Grove Store{% endblock %}

{% block content %}
{% meta name="description" content="Your shopping cart" %}

<h1 style="margin: 0 0 1.5rem;">Shopping Cart</h1>

{% if items | length > 0 %}
  <div class="table-wrap">
    <table class="table">
      <thead>
        <tr>
          <th>Product</th>
          <th>Price</th>
          <th>Qty</th>
          <th style="text-align: right;">Total</th>
        </tr>
      </thead>
      <tbody>
        {% for item in items %}
          <tr>
            <td>
              <a href="/product/{{ item.product.slug }}">{{ item.product.name }}</a>
            </td>
            <td>
              {{ pricing.price(item.product.price, item.product.sale_price) }}
            </td>
            <td>{{ item.quantity }}</td>
            <td style="text-align: right;" class="font-bold">{{ item.line_total | currency }}</td>
          </tr>
        {% endfor %}
      </tbody>
    </table>
  </div>

  {% let %}
    subtotal = 0
  {% endlet %}
  {% for item in items %}
    {% set subtotal = subtotal + item.line_total %}
  {% endfor %}

  <div class="cart-summary">
    <div class="cart-row">
      <span>Subtotal</span>
      <span class="font-bold">{{ subtotal | currency }}</span>
    </div>
    <div class="cart-row text-muted">
      <span>Shipping</span>
      <span>{{ subtotal >= 5000 ? "Free" : "$4.99" }}</span>
    </div>
    <hr class="cart-divider">
    {% set total = subtotal >= 5000 ? subtotal : subtotal + 499 %}
    <div class="cart-row" style="font-size: 1.2rem; font-weight: bold;">
      <span>Total</span>
      <span>{{ total | currency }}</span>
    </div>
    <button class="btn" style="margin-top: 1rem; width: 100%;">Checkout</button>
  </div>
{% else %}
  <p class="text-muted">Your cart is empty.</p>
  <a href="/" class="text-accent no-decor font-bold">Continue shopping →</a>
{% endif %}
{% endblock %}
```

- [ ] **Step 7: Run build check**

Run: `go build ./examples/store/...`
Expected: no errors

- [ ] **Step 8: Commit**

```bash
git add examples/store/templates/
git commit -m "feat(store): replace inline styles with CSS classes across all templates"
```

---

### Task 5: Create docs CSS file and static file server

**Files:**
- Create: `examples/docs/static/docs.css`
- Modify: `examples/docs/main.go`

- [ ] **Step 1: Create the static directory and CSS file**

Create `examples/docs/static/docs.css`:

```css
/* Grove Docs — Global Stylesheet */

:root {
  /* Brand colors (from logo) */
  --color-primary: #2E6740;
  --color-primary-hover: #245533;
  --color-dark: #251917;
  --color-dark-secondary: #3D2E2A;
  --color-cream: #EEEBE3;
  --color-page-bg: #F7F5F0;
  --color-text: #3D2E2A;
  --color-text-muted: #7A6B66;
  --color-border: #D9D3CB;
  --color-green-light: #E8F0EA;

  /* Admonition colors */
  --color-note-bg: #E8F0EA;
  --color-note-border: #2E6740;
  --color-note-text: #2E6740;
  --color-warning-bg: #FEF3C7;
  --color-warning-border: #D4A843;
  --color-warning-text: #6B5210;
  --color-tip-bg: #D4EDDA;
  --color-tip-border: #1B5E28;
  --color-tip-text: #1B5E28;

  /* Spacing */
  --radius-sm: 4px;
  --radius-md: 6px;
  --radius-lg: 8px;
}

/* ── Base ──────────────────────────────── */

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  color: var(--color-text);
  background: var(--color-page-bg);
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

a { color: var(--color-primary); }

/* ── Nav ───────────────────────────────── */

.nav {
  background: var(--color-dark);
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  gap: 2rem;
}

.nav-brand {
  color: var(--color-primary);
  font-size: 1.4rem;
  font-weight: bold;
  text-decoration: none;
}

/* ── Breadcrumb Bar ────────────────────── */

.breadcrumb-bar {
  background: var(--color-dark-secondary);
  padding: 0.5rem 2rem;
  font-size: 0.85rem;
}

.breadcrumb-bar .crumb-root { color: var(--color-text-muted); }
.breadcrumb-bar .crumb-sep { color: var(--color-text-muted); }
.breadcrumb-bar .crumb-section { color: var(--color-cream); }
.breadcrumb-bar .crumb-active { color: var(--color-primary); }

/* ── Layout ────────────────────────────── */

.docs-layout {
  flex: 1;
  display: flex;
}

.docs-main {
  flex: 1;
  padding: 2rem;
  max-width: 740px;
}

.container {
  max-width: 960px;
  width: 100%;
  margin: 0 auto;
  padding: 2rem 1rem;
}

/* ── Sidebar ───────────────────────────── */

.sidebar {
  width: 220px;
  padding: 1.5rem 1rem;
}

.sidebar-heading {
  margin: 1.5rem 0 0.5rem;
  font-size: 0.85rem;
  text-transform: uppercase;
  color: var(--color-text-muted);
  letter-spacing: 0.05em;
}

.sidebar-link {
  display: block;
  padding: 0.4rem 0.75rem;
  margin: 2px 0;
  border-radius: var(--radius-sm);
  color: var(--color-dark);
  text-decoration: none;
  font-size: 0.9rem;
}

.sidebar-link-active {
  background: var(--color-primary);
  color: var(--color-cream);
}

/* ── Footer ────────────────────────────── */

.footer {
  background: var(--color-dark);
  color: var(--color-text-muted);
  padding: 2rem;
  text-align: center;
}

.footer p { margin: 0; }
.footer a { color: var(--color-primary); }

/* ── Article ───────────────────────────── */

.article { line-height: 1.7; }

.page-nav {
  display: flex;
  justify-content: space-between;
  margin-top: 3rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--color-border);
}

.page-nav a {
  color: var(--color-primary);
  text-decoration: none;
}

/* ── Admonitions ───────────────────────── */

.admonition {
  padding: 12px 16px;
  border-left: 4px solid;
  border-radius: 0 var(--radius-md) var(--radius-md) 0;
  margin: 16px 0;
}

.admonition strong { display: block; }
.admonition div { margin-top: 4px; }

.admonition-note   { background: var(--color-note-bg);    border-color: var(--color-note-border);    color: var(--color-note-text); }
.admonition-warning { background: var(--color-warning-bg); border-color: var(--color-warning-border); color: var(--color-warning-text); }
.admonition-tip     { background: var(--color-tip-bg);     border-color: var(--color-tip-border);     color: var(--color-tip-text); }

/* ── Code / Filter Tags ────────────────── */

.code-tag {
  background: var(--color-green-light);
  padding: 0.2rem 0.6rem;
  border-radius: var(--radius-sm);
  font-size: 0.85rem;
}

/* ── Pagination ────────────────────────── */

.pagination {
  display: flex;
  gap: 0.5rem;
  margin: 1rem 0;
}

.page-num {
  display: inline-block;
  width: 2rem;
  height: 2rem;
  line-height: 2rem;
  text-align: center;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}

.page-num-active {
  background: var(--color-primary);
  color: var(--color-cream);
  border-color: var(--color-primary);
}

/* ── Utilities ─────────────────────────── */

.text-muted { color: var(--color-text-muted); }
.text-small { font-size: 0.85rem; }
.flex-row { display: flex; flex-wrap: wrap; gap: 0.5rem; }
```

- [ ] **Step 2: Add static file server to docs main.go**

In `examples/docs/main.go`, after the `r.Get("/docs/{section}/{page}", ...)` line, add:

```go
staticDir := filepath.Join(filepath.Dir(thisFile), "static")
r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
```

- [ ] **Step 3: Verify the docs build**

Run: `go build ./examples/docs/...`
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add examples/docs/static/docs.css examples/docs/main.go
git commit -m "feat(docs): add global CSS file with brand variables and static file server"
```

---

### Task 6: Update docs templates to use CSS classes

**Files:**
- Modify: `examples/docs/templates/base.grov`
- Modify: `examples/docs/templates/docs-layout.grov`
- Modify: `examples/docs/templates/pages/_default.grov`
- Modify: `examples/docs/templates/pages/variables-and-filters.grov`
- Modify: `examples/docs/templates/partials/sidebar.grov`
- Modify: `examples/docs/templates/macros/admonitions.grov`

- [ ] **Step 1: Update base.grov**

Replace the entire file with:

```html
{% asset "/static/docs.css" type="stylesheet" priority=10 %}
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{% block title %}Grove Docs{% endblock %}</title>
  <!-- HEAD_ASSETS -->
  <!-- HEAD_META -->
  <!-- HEAD_HOISTED -->
</head>
<body>
  {% block nav %}
  <nav class="nav">
    <a href="/" class="nav-brand">{{ site_name }}</a>
  </nav>
  {% endblock %}
  <div class="docs-layout">
    {% block layout %}
      <main class="container">
        {% block content %}{% endblock %}
      </main>
    {% endblock %}
  </div>
  <footer class="footer">
    <p>© {{ current_year }} Grove Docs. Built with the Grove template engine.</p>
  </footer>
  <!-- FOOT_ASSETS -->
</body>
</html>
```

- [ ] **Step 2: Update docs-layout.grov**

Replace the entire file with:

```html
{% extends "base.grov" %}

{% block nav %}
  {{ super() }}
  <div class="breadcrumb-bar">
    <span class="crumb-root">Documentation</span>
    {% if page %}
      <span class="crumb-sep"> / </span>
      <span class="crumb-section">{{ page.section }}</span>
      <span class="crumb-sep"> / </span>
      <span class="crumb-active">{{ page.title }}</span>
    {% endif %}
  </div>
{% endblock %}

{% block layout %}
  {% render "partials/sidebar.grov" sections=sections all_pages=all_pages current_slug=page.slug %}
  <main class="docs-main">
    {% block content %}{% endblock %}
  </main>
{% endblock %}
```

- [ ] **Step 3: Update sidebar.grov**

Replace the entire file with:

```html
<nav class="sidebar">
  {% for section in sections %}
    <h3 class="sidebar-heading">{{ section }}</h3>
    {% for page in all_pages %}
      {% if page.section == section %}
        {% set section_slug = section | lower | replace(" ", "-") %}
        {% set href = "/docs/" ~ section_slug ~ "/" ~ page.slug %}
        {% if page.slug == current_slug %}
          <a href="{{ href }}" class="sidebar-link sidebar-link-active">{{ page.title }}</a>
        {% else %}
          <a href="{{ href }}" class="sidebar-link">{{ page.title }}</a>
        {% endif %}
      {% endif %}
    {% endfor %}
  {% endfor %}
</nav>
```

- [ ] **Step 4: Update admonitions.grov**

Replace the entire file with:

```html
{% macro note(message, title) %}
  {% if not title %}{% set title = "Note" %}{% endif %}
  <div class="admonition admonition-note">
    <strong>{{ title }}</strong>
    <div>{{ message | safe }}</div>
  </div>
{% endmacro %}

{% macro warning(message, title) %}
  {% if not title %}{% set title = "Warning" %}{% endif %}
  <div class="admonition admonition-warning">
    <strong>{{ title }}</strong>
    <div>{{ message | safe }}</div>
  </div>
{% endmacro %}

{% macro tip(message, title) %}
  {% if not title %}{% set title = "Tip" %}{% endif %}
  <div class="admonition admonition-tip">
    <strong>{{ title }}</strong>
    <div>{{ message | safe }}</div>
  </div>
{% endmacro %}
```

- [ ] **Step 5: Update _default.grov**

Replace the entire file with:

```html
{% extends "docs-layout.grov" %}
{% import "macros/admonitions.grov" as adm %}

{% block title %}{{ page.title }} — Grove Docs{% endblock %}

{% block content %}
{% meta name="description" content=page.title ~ " — Grove documentation" %}

{% let %}
  title = page.title
  section = page.section
{% endlet %}

<h1 style="margin: 0 0 0.5rem;">{{ title }}</h1>
<p class="text-muted text-small" style="margin: 0 0 2rem;">{{ section }}</p>

<article class="article">
  {{ page.body | safe }}
</article>

<div class="page-nav">
  {% if prev %}
    <a href="/docs/{{ prev.section | lower | replace(" ", "-") }}/{{ prev.slug }}">← {{ prev.title }}</a>
  {% else %}
    <span></span>
  {% endif %}
  {% if next %}
    <a href="/docs/{{ next.section | lower | replace(" ", "-") }}/{{ next.slug }}">{{ next.title }} →</a>
  {% endif %}
</div>
{% endblock %}
```

- [ ] **Step 6: Update variables-and-filters.grov**

Replace the entire file with:

```html
{% extends "docs-layout.grov" %}
{% import "macros/admonitions.grov" as adm %}

{% block title %}{{ page.title }} — Grove Docs{% endblock %}

{% block content %}
{% meta name="description" content="Variables and filters in Grove templates" %}

{% let %}
  title = page.title
  section = page.section
{% endlet %}

<h1 style="margin: 0 0 0.5rem;">{{ title }}</h1>
<p class="text-muted text-small" style="margin: 0 0 2rem;">{{ section }}</p>

<article class="article">
  {{ page.body | safe }}
</article>

{{ adm.tip("Filters can be chained: <code>{{ name | lower | truncate(20) }}</code>") }}

<h2 style="margin-top: 2rem;">Built-in Filters by Category</h2>

{% set filter_categories = [
  {"name": "String", "filters": ["upper", "lower", "title", "trim", "truncate", "replace", "split", "join"]},
  {"name": "Collection", "filters": ["length", "first", "last", "reverse", "sort", "unique", "map", "slice"]},
  {"name": "Numeric", "filters": ["abs", "floor", "ceil", "round"]},
  {"name": "HTML", "filters": ["escape", "safe", "nl2br", "striptags"]}
] %}

{% for category in filter_categories %}
  <h3 style="margin-top: 1.5rem;">{{ category.name }}</h3>
  <div class="flex-row">
    {% for filter in category.filters %}
      <code class="code-tag">{{ filter }}</code>
    {% empty %}
      <span class="text-muted">No filters in this category.</span>
    {% endfor %}
  </div>
{% endfor %}

{{ adm.note("See the <a href='#'>API reference</a> for full filter documentation.") }}

<h2 style="margin-top: 2rem;">Pagination Example</h2>
<p>Here are pages 1 through 5:</p>
<div class="pagination">
  {% for n in range(1, 5) %}
    <span class="page-num {{ n == 1 ? "page-num-active" : "" }}">{{ n }}</span>
  {% endfor %}
</div>

<div class="page-nav">
  {% if prev %}
    <a href="/docs/{{ prev.section | lower | replace(" ", "-") }}/{{ prev.slug }}">← {{ prev.title }}</a>
  {% else %}
    <span></span>
  {% endif %}
  {% if next %}
    <a href="/docs/{{ next.section | lower | replace(" ", "-") }}/{{ next.slug }}">{{ next.title }} →</a>
  {% endif %}
</div>
{% endblock %}
```

- [ ] **Step 7: Run build check**

Run: `go build ./examples/docs/...`
Expected: no errors

- [ ] **Step 8: Commit**

```bash
git add examples/docs/templates/
git commit -m "feat(docs): replace inline styles with CSS classes across all templates"
```

---

### Task 7: Update email templates to use brand colors

Email templates must keep all styles inline (email client compatibility). This task only updates the color values to match the brand.

**Files:**
- Modify: `examples/email/templates/base-email.grov`
- Modify: `examples/email/templates/index.grov`
- Modify: `examples/email/templates/welcome.grov`
- Modify: `examples/email/templates/order-confirmation.grov`
- Modify: `examples/email/templates/password-reset.grov`
- Modify: `examples/email/templates/helpers.grov`

- [ ] **Step 1: Update base-email.grov**

Replace the entire file with:

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <style>
    body { margin: 0; padding: 0; background: #F7F5F0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
    .wrapper { max-width: 600px; margin: 0 auto; background: #ffffff; }
  </style>
</head>
<body>
  {% block preheader %}{% endblock %}
  <div class="wrapper">
    <div style="background: #251917; padding: 24px; text-align: center;">
      <span style="color: #2E6740; font-size: 24px; font-weight: bold;">Grove Store</span>
    </div>
    <div style="padding: 32px 24px;">
      {% block body %}{% endblock %}
    </div>
    <div style="background: #F7F5F0; padding: 24px; text-align: center; color: #7A6B66; font-size: 12px;">
      {% block footer %}
        <p>&copy; 2026 Grove Store. You received this email because you have an account with us.</p>
      {% endblock %}
    </div>
  </div>
</body>
</html>
```

Color changes:
- Header bg: `#1a1a2e` → `#251917` (brand espresso)
- Brand text: `#e94560` → `#2E6740` (brand green)
- Footer bg: `#f4f4f7` → `#F7F5F0` (brand page bg)
- Body bg: `#f4f4f7` → `#F7F5F0`
- Footer text: `#888` → `#7A6B66` (brand muted)

- [ ] **Step 2: Update helpers.grov**

Replace the entire file with:

```html
{% macro button(text, href, color) %}
  {% if not color %}{% set color = "#2E6740" %}{% endif %}
  <a href="{{ href }}" style="display: inline-block; padding: 12px 24px; background: {{ color }}; color: #ffffff; text-decoration: none; border-radius: 6px; font-weight: 600;">{{ text }}</a>
{% endmacro %}

{% macro divider() %}
  <hr style="border: none; border-top: 1px solid #D9D3CB; margin: 24px 0;">
{% endmacro %}

{% macro spacer(height) %}
  {% if not height %}{% set height = 16 %}{% endif %}
  <div style="height: {{ height }}px;"></div>
{% endmacro %}
```

Color changes:
- Default button: `#e94560` → `#2E6740` (brand green)
- Divider: `#eee` → `#D9D3CB` (brand border)

- [ ] **Step 3: Update welcome.grov**

Replace the entire file with:

```html
{% extends "base-email.grov" %}
{% import "helpers.grov" as h %}

{% block preheader %}
  {% capture preheader_text %}Welcome aboard, {{ user.name }}! Here's how to get started.{% endcapture %}
  {% hoist target="preheader" %}{{ preheader_text }}{% endhoist %}
  <div style="display: none; max-height: 0; overflow: hidden;">
    {{ preheader_text }}
  </div>
{% endblock %}

{% block body %}
  <h1 style="margin: 0 0 16px; color: #251917;">Welcome, {{ user.name }}!</h1>
  <p style="color: #3D2E2A; line-height: 1.6;">
    Thanks for joining Grove Store. We're excited to have you on board.
  </p>

  {% capture greeting_block %}
    <div style="background: #E8F0EA; border-radius: 8px; padding: 16px; margin: 16px 0;">
      <strong>Your account:</strong> {{ user.email }}
    </div>
  {% endcapture %}
  {{ greeting_block | safe }}

  {{ h.divider() }}
  <p style="color: #3D2E2A; line-height: 1.6;">Ready to start shopping? Check out our latest products:</p>
  {{ h.spacer(8) }}
  {{ h.button("Browse Products", "https://example.com/products") }}
{% endblock %}
```

Color changes:
- H1: `#1a1a2e` → `#251917`
- Body text: `#555` → `#3D2E2A`
- Greeting bg: `#f0fdf4` → `#E8F0EA` (brand green light)

- [ ] **Step 4: Update order-confirmation.grov**

Replace the entire file with:

```html
{% extends "base-email.grov" %}
{% import "helpers.grov" as h %}

{% block preheader %}
  {% capture preheader_text %}Order {{ order.id }} confirmed — thanks for your purchase!{% endcapture %}
  {% hoist target="preheader" %}{{ preheader_text }}{% endhoist %}
  <div style="display: none; max-height: 0; overflow: hidden;">
    {{ preheader_text }}
  </div>
{% endblock %}

{% block body %}
  <h1 style="margin: 0 0 16px; color: #251917;">Order Confirmed!</h1>
  <p style="color: #3D2E2A;">Hi {{ user.name | default("Customer") }}, your order <strong>{{ order.id | upper }}</strong> has been placed.</p>

  {{ h.divider() }}

  <table style="width: 100%; border-collapse: collapse;">
    <thead>
      <tr>
        <th style="text-align: left; padding: 8px 0; border-bottom: 2px solid #D9D3CB;">Item</th>
        <th style="text-align: center; padding: 8px 0; border-bottom: 2px solid #D9D3CB;">Qty</th>
        <th style="text-align: right; padding: 8px 0; border-bottom: 2px solid #D9D3CB;">Price</th>
      </tr>
    </thead>
    <tbody>
      {% for item in order.items %}
        <tr>
          <td style="padding: 8px 0; border-bottom: 1px solid #F7F5F0;">{{ item.name }}</td>
          <td style="text-align: center; padding: 8px 0; border-bottom: 1px solid #F7F5F0;">{{ item.quantity }}</td>
          <td style="text-align: right; padding: 8px 0; border-bottom: 1px solid #F7F5F0;">{{ item.line_total | currency }}</td>
        </tr>
      {% empty %}
        <tr>
          <td colspan="3" style="padding: 16px 0; text-align: center; color: #7A6B66;">No items in this order.</td>
        </tr>
      {% endfor %}
    </tbody>
  </table>

  {{ h.spacer(8) }}
  <div style="text-align: right; font-size: 18px; font-weight: bold;">
    Total: {{ order.total | currency }}
  </div>

  {{ h.divider() }}
  {{ h.button("View Order", "https://example.com/orders/" ~ order.id) }}
{% endblock %}
```

Color changes:
- H1: `#1a1a2e` → `#251917`
- Body text: `#555` → `#3D2E2A`
- Table borders: `#eee` → `#D9D3CB`, `#f4f4f7` → `#F7F5F0`
- Empty text: `#888` → `#7A6B66`

- [ ] **Step 5: Update password-reset.grov**

Replace the entire file with:

```html
{% extends "base-email.grov" %}
{% import "helpers.grov" as h %}

{% block body %}
  <h1 style="margin: 0 0 16px; color: #251917;">Reset Your Password</h1>
  <p style="color: #3D2E2A; line-height: 1.6;">
    Hi {{ user.name | default("there") }}, we received a request to reset your password. Click the button below to choose a new one:
  </p>
  {{ h.spacer(8) }}
  {{ h.button("Reset Password", "https://example.com/reset?token=abc123", "#251917") }}
  {{ h.spacer(16) }}
  <p style="color: #7A6B66; font-size: 13px;">
    If you didn't request this, you can safely ignore this email. The link expires in 24 hours.
  </p>
{% endblock %}

{% block footer %}
  <p>&copy; 2026 Grove Store. For security, this email was sent to {{ user.email | default("your address") }}.</p>
{% endblock %}
```

Color changes:
- H1: `#1a1a2e` → `#251917`
- Body text: `#555` → `#3D2E2A`
- Secondary button: `#0f3460` → `#251917` (brand dark)
- Small text: `#888` → `#7A6B66`

- [ ] **Step 6: Update index.grov**

Replace the entire file with:

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Grove Email Renderer</title>
</head>
<body style="margin: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #F7F5F0; min-height: 100vh;">
  <div style="max-width: 600px; margin: 0 auto; padding: 2rem 1rem;">
    <h1 style="color: #251917;">Grove Email Renderer</h1>
    <p style="color: #7A6B66;">Preview HTML email templates stored in a FileSystemStore.</p>
    <div style="display: grid; gap: 1rem; margin-top: 1.5rem;">
      {% for link in links %}
        <div style="background: #EEEBE3; border: 1px solid #D9D3CB; border-radius: 8px; padding: 1rem; display: flex; justify-content: space-between; align-items: center;">
          <strong>{{ link.label }}</strong>
          <div style="display: flex; gap: 0.75rem;">
            <a href="/preview/{{ link.name }}" style="color: #2E6740; text-decoration: none; font-weight: 600;">Preview</a>
            <a href="/source/{{ link.name }}" style="color: #251917; text-decoration: none; font-weight: 600;">Source</a>
          </div>
        </div>
      {% endfor %}
    </div>
  </div>
</body>
</html>
```

Color changes:
- Body bg: `#f8f9fa` → `#F7F5F0`
- H1: `#1a1a2e` → `#251917`
- Desc text: `#666` → `#7A6B66`
- Card bg: `#fff` → `#EEEBE3`
- Card border: `#ddd` → `#D9D3CB`
- Preview link: `#e94560` → `#2E6740`
- Source link: `#0f3460` → `#251917`

- [ ] **Step 7: Verify build**

Run: `go build ./examples/email/...`
Expected: no errors

- [ ] **Step 8: Commit**

```bash
git add examples/email/templates/
git commit -m "feat(email): update all email templates to use brand colors"
```

---

### Task 8: Update .gitignore and run full build verification

**Files:**
- Modify: `.gitignore` (if needed)

- [ ] **Step 1: Run full build check**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 2: Run full test suite**

Run: `go clean -testcache && go test ./... -v`
Expected: all tests pass (template changes don't affect engine tests since tests use inline templates, not the example files)

- [ ] **Step 3: Final commit**

If any adjustments were needed:

```bash
git add -A
git commit -m "chore: final adjustments after CSS refactor"
```
