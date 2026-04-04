package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	grove "grove/pkg/grove"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// User represents a registered user.
type User struct {
	Name  string
	Email string
}

func (u User) GroveResolve(key string) (any, bool) {
	switch key {
	case "name":
		return u.Name, true
	case "email":
		return u.Email, true
	}
	return nil, false
}

// OrderItem is a single line item.
type OrderItem struct {
	Name     string
	Quantity int
	Price    int // cents
}

func (oi OrderItem) GroveResolve(key string) (any, bool) {
	switch key {
	case "name":
		return oi.Name, true
	case "quantity":
		return oi.Quantity, true
	case "price":
		return oi.Price, true
	case "line_total":
		return oi.Price * oi.Quantity, true
	}
	return nil, false
}

// Order represents a placed order.
type Order struct {
	ID    string
	Items []OrderItem
	Total int // cents
}

func (o Order) GroveResolve(key string) (any, bool) {
	switch key {
	case "id":
		return o.ID, true
	case "items":
		out := make([]any, len(o.Items))
		for i, item := range o.Items {
			out[i] = item
		}
		return out, true
	case "total":
		return o.Total, true
	}
	return nil, false
}

var sampleUser = User{Name: "Alice", Email: "alice@example.com"}

var sampleOrder = Order{
	ID: "ORD-20260404",
	Items: []OrderItem{
		{Name: "Wireless Headphones", Quantity: 1, Price: 5999},
		{Name: "Running Shoes", Quantity: 2, Price: 6499},
	},
	Total: 18997,
}

var emptyOrder = Order{
	ID:    "ORD-EMPTY",
	Items: []OrderItem{},
	Total: 0,
}

// templateSources maps template names to their source text.
var templateSources = map[string]string{
	"base-email.grov": `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <style>
    body { margin: 0; padding: 0; background: #f4f4f7; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
    .wrapper { max-width: 600px; margin: 0 auto; background: #ffffff; }
  </style>
</head>
<body>
  {% block preheader %}{% endblock %}
  <div class="wrapper">
    <div style="background: #1a1a2e; padding: 24px; text-align: center;">
      <span style="color: #e94560; font-size: 24px; font-weight: bold;">Grove Store</span>
    </div>
    <div style="padding: 32px 24px;">
      {% block body %}{% endblock %}
    </div>
    <div style="background: #f4f4f7; padding: 24px; text-align: center; color: #888; font-size: 12px;">
      {% block footer %}
        <p>© 2026 Grove Store. You received this email because you have an account with us.</p>
      {% endblock %}
    </div>
  </div>
</body>
</html>`,

	"helpers.grov": `{% macro button(text, href, color) %}
  {% if not color %}{% set color = "#e94560" %}{% endif %}
  <a href="{{ href }}" style="display: inline-block; padding: 12px 24px; background: {{ color }}; color: #ffffff; text-decoration: none; border-radius: 6px; font-weight: 600;">{{ text }}</a>
{% endmacro %}

{% macro divider() %}
  <hr style="border: none; border-top: 1px solid #eee; margin: 24px 0;">
{% endmacro %}

{% macro spacer(height) %}
  {% if not height %}{% set height = 16 %}{% endif %}
  <div style="height: {{ height }}px;"></div>
{% endmacro %}`,

	"welcome.grov": `{% extends "base-email.grov" %}
{% import "helpers.grov" as h %}

{% block preheader %}
  {% capture preheader_text %}Welcome aboard, {{ user.name }}! Here's how to get started.{% endcapture %}
  {% hoist target="preheader" %}{{ preheader_text }}{% endhoist %}
  <div style="display: none; max-height: 0; overflow: hidden;">
    {{ preheader_text }}
  </div>
{% endblock %}

{% block body %}
  <h1 style="margin: 0 0 16px; color: #1a1a2e;">Welcome, {{ user.name }}!</h1>
  <p style="color: #555; line-height: 1.6;">
    Thanks for joining Grove Store. We're excited to have you on board.
  </p>

  {% capture greeting_block %}
    <div style="background: #f0fdf4; border-radius: 8px; padding: 16px; margin: 16px 0;">
      <strong>Your account:</strong> {{ user.email }}
    </div>
  {% endcapture %}
  {{ greeting_block | safe }}

  {{ h.divider() }}
  <p style="color: #555; line-height: 1.6;">Ready to start shopping? Check out our latest products:</p>
  {{ h.spacer(8) }}
  {{ h.button("Browse Products", "https://example.com/products") }}
{% endblock %}`,

	"order-confirmation.grov": `{% extends "base-email.grov" %}
{% import "helpers.grov" as h %}

{% block preheader %}
  {% capture preheader_text %}Order {{ order.id }} confirmed — thanks for your purchase!{% endcapture %}
  {% hoist target="preheader" %}{{ preheader_text }}{% endhoist %}
  <div style="display: none; max-height: 0; overflow: hidden;">
    {{ preheader_text }}
  </div>
{% endblock %}

{% block body %}
  <h1 style="margin: 0 0 16px; color: #1a1a2e;">Order Confirmed!</h1>
  <p style="color: #555;">Hi {{ user.name | default("Customer") }}, your order <strong>{{ order.id | upper }}</strong> has been placed.</p>

  {{ h.divider() }}

  <table style="width: 100%; border-collapse: collapse;">
    <thead>
      <tr>
        <th style="text-align: left; padding: 8px 0; border-bottom: 2px solid #eee;">Item</th>
        <th style="text-align: center; padding: 8px 0; border-bottom: 2px solid #eee;">Qty</th>
        <th style="text-align: right; padding: 8px 0; border-bottom: 2px solid #eee;">Price</th>
      </tr>
    </thead>
    <tbody>
      {% for item in order.items %}
        <tr>
          <td style="padding: 8px 0; border-bottom: 1px solid #f4f4f7;">{{ item.name }}</td>
          <td style="text-align: center; padding: 8px 0; border-bottom: 1px solid #f4f4f7;">{{ item.quantity }}</td>
          <td style="text-align: right; padding: 8px 0; border-bottom: 1px solid #f4f4f7;">{{ item.line_total | currency }}</td>
        </tr>
      {% empty %}
        <tr>
          <td colspan="3" style="padding: 16px 0; text-align: center; color: #888;">No items in this order.</td>
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
{% endblock %}`,

	"password-reset.grov": `{% extends "base-email.grov" %}
{% import "helpers.grov" as h %}

{% block body %}
  <h1 style="margin: 0 0 16px; color: #1a1a2e;">Reset Your Password</h1>
  <p style="color: #555; line-height: 1.6;">
    Hi {{ user.name | default("there") }}, we received a request to reset your password. Click the button below to choose a new one:
  </p>
  {{ h.spacer(8) }}
  {{ h.button("Reset Password", "https://example.com/reset?token=abc123", "#0f3460") }}
  {{ h.spacer(16) }}
  <p style="color: #888; font-size: 13px;">
    If you didn't request this, you can safely ignore this email. The link expires in 24 hours.
  </p>
{% endblock %}

{% block footer %}
  <p>© 2026 Grove Store. For security, this email was sent to {{ user.email | default("your address") }}.</p>
{% endblock %}`,
}

// emailPreviews defines which data each email template uses.
var emailPreviews = map[string]grove.Data{
	"welcome.grov":            {"user": sampleUser},
	"order-confirmation.grov": {"user": sampleUser, "order": sampleOrder},
	"password-reset.grov":     {"user": sampleUser},
}

func main() {
	ms := grove.NewMemoryStore()
	for name, src := range templateSources {
		ms.Set(name, src)
	}

	eng := grove.New(grove.WithStore(ms))
	eng.SetGlobal("current_year", "2026")

	// Currency filter: format cents as "$12.99"
	eng.RegisterFilter("currency", grove.FilterFn(func(v grove.Value, args []grove.Value) (grove.Value, error) {
		cents, _ := v.ToInt64()
		dollars := cents / 100
		remainder := cents % 100
		return grove.StringValue(fmt.Sprintf("$%d.%02d", dollars, remainder)), nil
	}))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", indexHandler(eng))
	r.Get("/preview/{name}", previewHandler(eng))
	r.Get("/source/{name}", sourceHandler())

	fmt.Println("Grove Email Renderer listening on http://localhost:3002")
	log.Fatal(http.ListenAndServe(":3002", r))
}

func indexHandler(eng *grove.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		names := []string{"welcome.grov", "order-confirmation.grov", "password-reset.grov"}
		var links []any
		for _, n := range names {
			links = append(links, map[string]any{
				"name":  n,
				"label": strings.TrimSuffix(n, ".grov"),
			})
		}

		src := `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Grove Email Renderer</title>
</head>
<body style="margin: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f8f9fa; min-height: 100vh;">
  <div style="max-width: 600px; margin: 0 auto; padding: 2rem 1rem;">
    <h1 style="color: #1a1a2e;">Grove Email Renderer</h1>
    <p style="color: #666;">Preview HTML email templates stored in a MemoryStore.</p>
    <div style="display: grid; gap: 1rem; margin-top: 1.5rem;">
      {% for link in links %}
        <div style="background: #fff; border: 1px solid #ddd; border-radius: 8px; padding: 1rem; display: flex; justify-content: space-between; align-items: center;">
          <strong>{{ link.label }}</strong>
          <div style="display: flex; gap: 0.75rem;">
            <a href="/preview/{{ link.name }}" style="color: #e94560; text-decoration: none; font-weight: 600;">Preview</a>
            <a href="/source/{{ link.name }}" style="color: #0f3460; text-decoration: none; font-weight: 600;">Source</a>
          </div>
        </div>
      {% endfor %}
    </div>
  </div>
</body>
</html>`

		result, err := eng.RenderTemplate(r.Context(), src, grove.Data{"links": links})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, result.Body)
	}
}

func previewHandler(eng *grove.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		data, ok := emailPreviews[name]
		if !ok {
			http.NotFound(w, r)
			return
		}
		result, err := eng.Render(r.Context(), name, data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, result.Body)
	}
}

func sourceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		src, ok := templateSources[name]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, src)
	}
}

var (
	_ interface{ GroveResolve(string) (any, bool) } = User{}
	_ interface{ GroveResolve(string) (any, bool) } = OrderItem{}
	_ interface{ GroveResolve(string) (any, bool) } = Order{}
)
