package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
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

// emailPreviews defines which data each email template uses.
var emailPreviews = map[string]grove.Data{
	"welcome.grov":            {"user": sampleUser},
	"order-confirmation.grov": {"user": sampleUser, "order": sampleOrder},
	"password-reset.grov":     {"user": sampleUser},
}

func main() {
	_, thisFile, _, _ := runtime.Caller(0)
	templateDir := filepath.Join(filepath.Dir(thisFile), "templates")

	store := grove.NewFileSystemStore(templateDir)
	eng := grove.New(grove.WithStore(store))
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
	r.Get("/source/{name}", sourceHandler(store))

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

		result, err := eng.Render(r.Context(), "index.grov", grove.Data{"links": links})
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

func sourceHandler(store *grove.FileSystemStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		src, err := store.Load(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(src)
	}
}

var (
	_ interface{ GroveResolve(string) (any, bool) } = User{}
	_ interface{ GroveResolve(string) (any, bool) } = OrderItem{}
	_ interface{ GroveResolve(string) (any, bool) } = Order{}
)
