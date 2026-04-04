# Examples Expansion Design — v2

## Goal

Overhaul Grove's four existing examples (blog, store, docs, email) from shallow demos into realistic, feature-rich mini-applications. Each example serves two audiences: developers evaluating Grove and developers learning from reference code. The approach is feature-driven — each example owns specific Grove capabilities and demonstrates them thoroughly — but shaped by realism so nothing feels artificial.

## Design Principles

1. **Every link goes somewhere real.** Breadcrumbs, tag links, category navigation, author bylines — all connect to actual pages with real content.
2. **Data loaded from JSON files.** Each example has a `data/` directory with JSON files that the Go code loads at startup. Datasets are large enough to feel real (10+ blog posts, 20+ products, 12+ doc pages).
3. **Simple interactivity via query params.** Filtering, sorting, and search work via GET requests and query parameters. No JavaScript required. Stateful features (like cart) use cookies.
4. **Feature ownership with natural overlap.** Each Grove feature has one example that demonstrates it thoroughly (its "primary home"). Features appear elsewhere when natural but aren't forced.
5. **Each example is a standalone `go run main.go` HTTP server** with its own `go.mod` and a `replace` directive pointing to the root module.

## Feature Ownership Matrix

Each example owns specific Grove features (marked **primary**). Features appear naturally in other examples (marked **use**) but each has one thorough demonstration.

| Feature | Blog | Store | Docs | Email |
|---|:---:|:---:|:---:|:---:|
| `extends` / `block` | **primary** | use | use | use |
| `super()` | | | **primary** | |
| `component` / `slot` / `fill` / `props` | **primary** | use | | |
| `for` / `empty` | use | **primary** | use | use |
| `if` / `elif` / `else` | use | use | **primary** | use |
| `macro` / `import` | | **primary** | use | use |
| `render` (partials with params) | | | **primary** | |
| `include` | use | | | |
| `set` / `let` | use | use | use | **primary** |
| `capture` | | | | **primary** |
| `range()` | | **primary** | | |
| Ternary `? :` | | **primary** | | |
| Filters (pipe chains) | **primary** | use | use | use |
| Custom filters (Go-registered) | | **primary** | | use |
| `asset` / `meta` / `hoist` | use | use | use | **primary** |
| `safe` filter | use | | | **primary** |
| List/map literals | | **primary** | use | |
| String concatenation `~` | use | use | **primary** | |
| Global variables | use | use | use | use |
| `GroveResolve` interface | **primary** | use | use | use |
| Sandboxing (tag/filter whitelists) | | | **primary** | |

## Example 1: Blog

### What it is

A multi-author tech blog with posts, tags, authors, and archive pages.

### Data (`data/` JSON files)

- **posts.json** — 10-12 posts across 3-4 categories (Go, Templates, Web Dev, Releases). Each post has: title, slug, author_slug, date, tags (array of slugs), summary, body (HTML), draft (bool).
- **authors.json** — 3 authors with: name, slug, bio, avatar_url, role.
- **tags.json** — 6-8 tags with: name, slug, color.

### Routes

| Route | Description |
|---|---|
| `GET /` | Homepage with recent posts |
| `GET /posts?tag=go&author=jane` | Filtered post listing via query params |
| `GET /post/{slug}` | Single post with full content |
| `GET /tags` | All tags with post counts |
| `GET /tag/{slug}` | Posts filtered by a specific tag |
| `GET /author/{slug}` | Author page with bio and their posts |

### Navigation

- Tag badges on post cards link to `/tag/{slug}` — real filtered results
- Author bylines link to `/author/{slug}` — real author profile pages
- Breadcrumbs: Home → Post (on post pages), Home → Tag → (tag name) (on tag pages), Home → Author → (author name) (on author pages)
- "Related posts" at bottom of post page — other posts sharing tags
- Nav bar links to Home, Tags

### Template structure

```
templates/
├── base.grov                    # Master layout: blocks for title, content
├── index.grov                   # Homepage: recent posts grid
├── post.grov                    # Single post: full content, author card, related posts
├── post-list.grov               # Filtered post listing (used by /posts, /tag, /author)
├── tag-list.grov                # All tags overview with post counts
├── author.grov                  # Author profile + their posts
└── components/
    ├── nav.grov                 # Site navigation
    ├── footer.grov              # Site footer
    ├── card.grov                # Blog post card (title, summary, tags, author, date)
    ├── tag-badge.grov           # Tag with color and link
    ├── author-card.grov         # Author avatar, name, bio
    ├── pagination.grov          # Prev/next or numbered pagination
    └── breadcrumbs.grov         # Breadcrumb trail
```

### Primary Grove features demonstrated

- **Template inheritance** (`extends`/`block`) — base layout with page-specific overrides
- **Components with slots** (`component`/`slot`/`fill`/`props`) — card component with tag slot, author-card with optional bio slot
- **Filter chains** — e.g., `| truncate(120)`, `| default("Unknown")`, `| length`
- **`GroveResolve` interface** — Post, Author, Tag structs all implement it, demonstrating the full pattern for Go-to-template data binding
- **`for`/`empty`** — Post loops with empty state for no results
- **`include`** — Simple partials without isolated scope

### Go code highlights

- JSON loading at startup with proper struct unmarshaling
- Query param parsing for tag/author filtering
- Helper functions to resolve relationships (post → author, post → tags)
- Pagination logic for post listings

---

## Example 2: Store

### What it is

A product catalog for a fictional outdoor gear shop with category browsing, filtering/sorting, search, and a cookie-based cart.

### Data (`data/` JSON files)

- **products.json** — 20-25 products across 5 categories. Each product has: name, slug, price (cents), sale_price (cents, 0 if not on sale), description, body (HTML), image_url, category_slug, rating (float), review_count, colors (array), sizes (array), in_stock (bool), featured (bool).
- **categories.json** — 5 categories (Camping, Hiking, Cycling, Running, Climbing). Each has: name, slug, description.

### Routes

| Route | Description |
|---|---|
| `GET /` | Homepage with featured products and category grid |
| `GET /products` | All products |
| `GET /products?category=hiking&sort=price-asc&min_price=2000&max_price=10000` | Filtered + sorted listing |
| `GET /category/{slug}` | Category page with description + products |
| `GET /product/{slug}` | Product detail page |
| `GET /cart` | Cart page (reads cart from cookie) |
| `GET /cart/add?product={slug}&qty=1` | Add to cart, redirects back |
| `GET /cart/remove?product={slug}` | Remove from cart, redirects back |
| `GET /search?q=tent` | Text search across product names and descriptions |

### Navigation

- Category links in nav → real category pages with filtered products
- Breadcrumbs: Home → Category → Product (derived from product's actual category)
- "Related products" on detail page — same category, different product
- Sort/filter controls on listing pages with real query params
- Cart link in nav showing item count from cookie
- Search form in nav that submits to `/search`

### Template structure

```
templates/
├── base.grov                    # Layout: nav with categories + cart badge, footer
├── index.grov                   # Featured products + category cards
├── product-list.grov            # Filterable/sortable product grid
├── category.grov                # Category header + filtered products
├── product.grov                 # Full detail: pricing, options, breadcrumbs, related
├── cart.grov                    # Line items, totals, empty state
├── search.grov                  # Search results with query echo and count
├── components/
│   └── product-card.grov        # Product card component with price slot
└── macros/
    ├── pricing.grov             # Price display, sale badge, star rating, discount calc
    └── filters.grov             # Sort dropdown, active filter pills, price range
```

### Primary Grove features demonstrated

- **Macros and `import`** — pricing macros (price display, star rating, discount percentage), filter UI macros (sort controls, filter pills)
- **Custom Go-registered filters** — `currency` (cents → `$12.99`), `stars` (rating → star display)
- **`range()`** — Quantity selector on product page (1–10), pagination page numbers
- **Ternary expressions** — Conditional CSS classes (`product.on_sale ? "sale" : ""`), display toggling
- **List/map literals** — Building breadcrumb arrays, filter state objects inline in templates
- **`for`/`empty`** — Products, cart items, categories, search results with empty states
- **Arithmetic** — Cart subtotals, discount calculations, total with shipping threshold

### Go code highlights

- JSON loading with category-to-product relationship resolution
- Cookie-based cart: encode cart as JSON in a cookie, decode on each request
- Query param parsing for sort, category filter, price range, search query
- Text search (simple substring match across name + description)
- Redirect-after-POST pattern (GET-based cart add/remove → redirect to referrer)

---

## Example 3: Docs

### What it is

A documentation site that actually documents Grove's template syntax. The example is its own reference — someone reading the rendered docs learns how to use the engine that renders them.

### Data (`data/` JSON files or structured Go)

- **sections.json** — 4-5 sections: Getting Started, Template Syntax, Tags, Filters, Advanced. Each has: name, slug, description, order.
- **pages.json** — 12-15 pages total with accurate Grove documentation. Each has: title, slug, section_slug, body (HTML with code examples), order (within section).
- **filters.json** — All 40+ built-in filters with: name, description, category (String, Collection, Numeric, HTML, Date), example_input, example_output.

### Pages (accurate Grove documentation)

**Getting Started:**
- Installation — go get, module setup
- Quick Start — first template, rendering from Go
- Template Basics — variables, expressions, comments

**Template Syntax:**
- Variables & Expressions — interpolation, arithmetic, comparisons, ternary, string concat
- Filters — pipe syntax, chaining, all built-in filters with examples
- Control Flow — if/elif/else, for/empty, range, set/let

**Tags:**
- Template Inheritance — extends, block, super()
- Includes & Partials — include, render (with params)
- Components — component, slot, fill, props
- Macros — macro, import, call

**Advanced:**
- Asset Management — asset, meta, hoist
- Sandboxing — tag/filter whitelists, loop limits
- Custom Filters — RegisterFilter from Go
- Go Integration — GroveResolve, Engine API, RenderResult

### Routes

| Route | Description |
|---|---|
| `GET /` | Docs landing page (overview + section links) |
| `GET /docs/{section}` | Section index listing pages in that section |
| `GET /docs/{section}/{page}` | Individual doc page |
| `GET /docs/filters?q=upper&category=string` | Searchable filter reference |

### Navigation

- Sidebar with all sections and pages, current page highlighted, sections collapsible
- Breadcrumbs: Docs → Section → Page (every crumb links to a real page)
- Prev/next pagination derived from actual page ordering across sections
- Section index pages list all pages within that section

### Template structure

```
templates/
├── base.grov                         # Minimal HTML shell with asset/meta placeholders
├── docs-layout.grov                  # Extends base: sidebar + breadcrumbs, uses super()
├── pages/
│   ├── _default.grov                 # Generic doc page (most pages use this)
│   ├── variables-and-filters.grov    # Custom: interactive filter reference table
│   └── template-inheritance.grov     # Custom: explains what it's doing (meta)
├── partials/
│   ├── sidebar.grov                  # Rendered via {% render %} with explicit params
│   ├── breadcrumbs.grov              # Breadcrumb trail from page hierarchy
│   └── filter-table.grov            # Filter reference with category grouping
└── macros/
    ├── admonitions.grov              # note(), warning(), tip() macros
    └── code-example.grov             # Macro for displaying labeled code snippets
```

### Primary Grove features demonstrated

- **Multi-level inheritance with `super()`** — base → docs-layout → page (3 levels). docs-layout uses `super()` to extend the nav block from base while adding breadcrumbs.
- **`render` with explicit params** — sidebar partial receives `sections`, `all_pages`, `current_slug` as explicit parameters with isolated scope
- **Sandboxing** — Engine configured with `WithSandbox`: tag whitelist, filter whitelist, `MaxLoopIter`. Demonstrates the security model.
- **String concatenation `~`** — Building URLs: `"/docs/" ~ section_slug ~ "/" ~ page.slug`
- **`if`/`elif`/`else` chains** — Section-aware rendering, conditional active states, filter category matching
- **`import` for macro libraries** — Admonition macros and code example macros imported across pages

### The content twist

Doc pages contain accurate Grove syntax documentation with real code examples. The "Template Inheritance" page explains extends/block/super() — the very features rendering it. The "Filters" page lists real filters with real examples. This makes the example genuinely useful as a quick reference.

### Go code highlights

- JSON loading for sections, pages, filters
- Page lookup with section-aware routing
- Prev/next calculation across sections (ordered)
- Template selection: specific template if exists (`pages/{slug}.grov`), else `pages/_default.grov`
- Sandbox configuration with explicit tag/filter whitelists
- Filter reference query param handling (text search + category filter)

---

## Example 4: Email

### What it is

A transactional email template system for a fictional SaaS product ("Grove Cloud"), with a preview server for viewing rendered emails with different data scenarios.

### Data (`data/` JSON files)

- **users.json** — 3-4 user profiles with varying states: name, email, plan (free/pro/enterprise), joined_date.
- **orders.json** — 2-3 sample orders with line items: id, user_id, items (array of name/quantity/price), total, date.
- **scenarios.json** — Named data scenarios for previewing emails with different contexts (e.g., "new_user", "enterprise_user", "expired_token", "empty_order").

### Email templates

| Template | Description |
|---|---|
| `welcome.grov` | Welcome email with onboarding steps, CTA button |
| `order-confirmation.grov` | Receipt with line items, subtotals, totals |
| `password-reset.grov` | Reset link with expiry, security notice |
| `plan-change.grov` | Upgrade/downgrade notification with before/after comparison |
| `usage-alert.grov` | Approaching plan limits with usage bar and upgrade CTA |

### Routes

| Route | Description |
|---|---|
| `GET /` | Index listing all email templates with preview links |
| `GET /preview/{name}` | Rendered email preview with default data |
| `GET /preview/{name}?user=2&scenario=expired` | Preview with specific user/scenario |
| `GET /source/{name}` | Raw template source view |

### Template structure

```
templates/
├── base-email.grov              # Email-safe HTML layout (table-based, inline styles)
│                                # Blocks: preheader, body, footer
├── helpers.grov                 # Macro library: button(), divider(), spacer(),
│                                # heading(), usage-bar()
├── welcome.grov                 # Extends base, uses hoist for preheader
├── order-confirmation.grov      # Extends base, line item loop, arithmetic
├── password-reset.grov          # Extends base, uses capture for greeting block
├── plan-change.grov             # Extends base, before/after with let blocks
└── usage-alert.grov             # Extends base, usage-bar macro, conditional urgency
```

### Primary Grove features demonstrated

- **`hoist`** — Moving content to email preheader from within body templates. Each email hoists a contextual preheader snippet that differs per template.
- **`capture`** — Building complex content blocks (e.g., assembling a personalized greeting block with conditional logic, then inserting it in multiple places)
- **`let`/`set` scoping** — Local variable blocks for computing display values: savings amount, usage percentages, formatted dates, plan comparison data
- **`safe` filter** — Rendering pre-built HTML content without escaping (captured blocks, macro output)
- **`asset`/`meta`** — Collecting stylesheet references and meta info across nested email templates
- **Macro composition** — Macros calling other macros: button inside a card-like pattern, heading with divider

### What makes it feel real

- Multiple preview scenarios per template — different users, edge cases (empty order, expired token)
- Preheader text that differs per email type via `hoist`
- Realistic SaaS content: plan tiers, usage numbers, team context
- Email-safe HTML patterns: table-based layout, inline styles, 600px max-width
- Properly structured preview server for template development workflow

### Go code highlights

- JSON loading for users, orders, scenarios
- Scenario resolution: merge base data with scenario-specific overrides
- Custom `currency` filter (shared concept with store, different registration)
- FileSystemStore like other examples, but note in comments that MemoryStore is an alternative for DB-stored templates
- RenderResult processing: preheader extraction from hoisted content, asset injection

---

## Conventions (all examples)

- Each example is a standalone `go run main.go` HTTP server
- Each has its own `go.mod` with a `replace` directive pointing to the root module
- Data loaded from `data/*.json` files at startup
- `GroveResolve` implemented on all custom Go structs
- Template files use `.grov` extension
- External CSS in `static/` directory (not inline styles, except email which requires inline)
- Chi router for HTTP routing (consistent across examples)
- `writeResult()` helper for processing `RenderResult` (assets, meta, hoisted content)
- Global variables set via `eng.SetGlobal()`: `site_name`, `current_year`

## File structure per example

```
examples/{name}/
├── main.go                      # HTTP server, data loading, handlers
├── data/
│   └── *.json                   # Data files
├── templates/
│   ├── base.grov                # Base layout
│   ├── *.grov                   # Page templates
│   ├── components/              # (blog, store) Reusable components
│   ├── partials/                # (docs) Rendered partials
│   └── macros/                  # (store, docs, email) Macro libraries
├── static/
│   └── style.css                # Stylesheet (except email)
├── go.mod
└── go.sum
```
