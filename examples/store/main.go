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

// Product represents an item in the store.
type Product struct {
	Name        string
	Slug        string
	Price       int // cents
	SalePrice   int // cents; 0 = not on sale
	Description string
	ImageURL    string
	Category    string
	Rating      float64
	ReviewCount int
	Colors      []string
	InStock     bool
}

func (p Product) GroveResolve(key string) (any, bool) {
	switch key {
	case "name":
		return p.Name, true
	case "slug":
		return p.Slug, true
	case "price":
		return p.Price, true
	case "sale_price":
		return p.SalePrice, true
	case "on_sale":
		return p.SalePrice > 0, true
	case "description":
		return p.Description, true
	case "image_url":
		return p.ImageURL, true
	case "category":
		return p.Category, true
	case "rating":
		return p.Rating, true
	case "review_count":
		return p.ReviewCount, true
	case "colors":
		out := make([]any, len(p.Colors))
		for i, c := range p.Colors {
			out[i] = c
		}
		return out, true
	case "in_stock":
		return p.InStock, true
	}
	return nil, false
}

// CartItem pairs a product with a quantity.
type CartItem struct {
	Product  Product
	Quantity int
}

func (ci CartItem) GroveResolve(key string) (any, bool) {
	switch key {
	case "product":
		return ci.Product, true
	case "quantity":
		return ci.Quantity, true
	case "line_total":
		price := ci.Product.Price
		if ci.Product.SalePrice > 0 {
			price = ci.Product.SalePrice
		}
		return price * ci.Quantity, true
	}
	return nil, false
}

var products = []Product{
	{
		Name:        "Wireless Headphones",
		Slug:        "wireless-headphones",
		Price:       7999,
		SalePrice:   5999,
		Description: "Premium over-ear headphones with active noise cancellation, 30-hour battery life, and a comfortable fit for all-day listening.",
		ImageURL:    "https://placehold.co/400x300/251917/2E6740?text=Headphones",
		Category:    "Electronics",
		Rating:      4.5,
		ReviewCount: 128,
		Colors:      []string{"Black", "Silver", "Navy"},
		InStock:     true,
	},
	{
		Name:        "Mechanical Keyboard",
		Slug:        "mechanical-keyboard",
		Price:       12999,
		SalePrice:   0,
		Description: "Compact 75% layout with hot-swappable switches, RGB backlighting, and a solid aluminum frame.",
		ImageURL:    "https://placehold.co/400x300/2E6740/EEEBE3?text=Keyboard",
		Category:    "Electronics",
		Rating:      4.8,
		ReviewCount: 64,
		Colors:      []string{"White", "Black"},
		InStock:     true,
	},
	{
		Name:        "Running Shoes",
		Slug:        "running-shoes",
		Price:       8999,
		SalePrice:   6499,
		Description: "Lightweight and responsive running shoes with a breathable mesh upper and cushioned sole.",
		ImageURL:    "https://placehold.co/400x300/251917/2E6740?text=Shoes",
		Category:    "Footwear",
		Rating:      4.2,
		ReviewCount: 203,
		Colors:      []string{"Red", "Blue", "Green", "Black"},
		InStock:     true,
	},
	{
		Name:        "Desk Lamp",
		Slug:        "desk-lamp",
		Price:       3499,
		SalePrice:   0,
		Description: "Adjustable LED desk lamp with five brightness levels and a built-in USB charging port.",
		ImageURL:    "https://placehold.co/400x300/3D2E2A/EEEBE3?text=Lamp",
		Category:    "Home",
		Rating:      4.0,
		ReviewCount: 42,
		Colors:      []string{"White", "Black"},
		InStock:     false,
	},
}

var cart = []CartItem{
	{Product: products[0], Quantity: 1},
	{Product: products[2], Quantity: 2},
}

func main() {
	_, thisFile, _, _ := runtime.Caller(0)
	templateDir := filepath.Join(filepath.Dir(thisFile), "templates")

	store := grove.NewFileSystemStore(templateDir)
	eng := grove.New(grove.WithStore(store))
	eng.SetGlobal("site_name", "Grove Store")
	eng.SetGlobal("current_year", "2026")

	// Custom filter: format cents as "$12.99"
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
	r.Get("/product/{slug}", productHandler(eng))
	r.Get("/cart", cartHandler(eng))

	staticDir := filepath.Join(filepath.Dir(thisFile), "static")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	fmt.Println("Grove Store listening on http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", r))
}

func indexHandler(eng *grove.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productsAny := make([]any, len(products))
		for i, p := range products {
			productsAny[i] = p
		}
		result, err := eng.Render(r.Context(), "index.grov", grove.Data{
			"products": productsAny,
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeResult(w, result)
	}
}

func productHandler(eng *grove.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		var found *Product
		for i := range products {
			if products[i].Slug == slug {
				found = &products[i]
				break
			}
		}
		if found == nil {
			http.NotFound(w, r)
			return
		}
		result, err := eng.Render(r.Context(), "product.grov", grove.Data{
			"product": *found,
			"breadcrumbs": []any{
				map[string]any{"label": "Home", "href": "/"},
				map[string]any{"label": found.Category, "href": "/"},
				map[string]any{"label": found.Name, "href": ""},
			},
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeResult(w, result)
	}
}

func cartHandler(eng *grove.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartAny := make([]any, len(cart))
		for i, ci := range cart {
			cartAny[i] = ci
		}
		result, err := eng.Render(r.Context(), "cart.grov", grove.Data{
			"items": cartAny,
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeResult(w, result)
	}
}

func writeResult(w http.ResponseWriter, result grove.RenderResult) {
	body := result.Body
	body = strings.Replace(body, "<!-- HEAD_ASSETS -->", result.HeadHTML(), 1)

	var meta strings.Builder
	for name, content := range result.Meta {
		if strings.HasPrefix(name, "og:") || strings.HasPrefix(name, "property:") {
			meta.WriteString(fmt.Sprintf(`  <meta property="%s" content="%s">`+"\n", name, content))
		} else {
			meta.WriteString(fmt.Sprintf(`  <meta name="%s" content="%s">`+"\n", name, content))
		}
	}
	body = strings.Replace(body, "<!-- HEAD_META -->", meta.String(), 1)
	body = strings.Replace(body, "<!-- FOOT_ASSETS -->", result.FootHTML(), 1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, body)
}

var (
	_ interface{ GroveResolve(string) (any, bool) } = Product{}
	_ interface{ GroveResolve(string) (any, bool) } = CartItem{}
)
