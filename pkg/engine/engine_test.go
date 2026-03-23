package engine

import (
	"strings"
	"testing"

	"template-wisp/internal/store"
)

func TestRenderString(t *testing.T) {
	e := New()

	tests := []struct {
		name     string
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			name:     "variable access",
			template: `{% .name %}`,
			data:     map[string]interface{}{"name": "Alice"},
			expected: "Alice",
		},
		{
			name:     "if true",
			template: `{% if.show%}visible{%end%}`,
			data:     map[string]interface{}{"show": true},
			expected: "visible",
		},
		{
			name:     "if false with else",
			template: `{%if.show%}yes{%else%}no{%end%}`,
			data:     map[string]interface{}{"show": false},
			expected: "no",
		},
		{
			name:     "for loop",
			template: `{% for .item in .items %}{%.item%}{%end%}`,
			data:     map[string]interface{}{"items": []interface{}{"a", "b", "c"}},
			expected: "abc",
		},
		{
			name:     "text content",
			template: `Hello{%.name%}!`,
			data:     map[string]interface{}{"name": "World"},
			expected: "HelloWorld!",
		},
		{
			name:     "assign variable",
			template: `{%assign.x="hello"%}{%.x%}`,
			data:     map[string]interface{}{},
			expected: "hello",
		},
		{
			name:     "unless statement",
			template: `{%unless.hide%}shown{%end%}`,
			data:     map[string]interface{}{"hide": false},
			expected: "shown",
		},
		{
			name:     "comment block",
			template: `{%comment%}secret{%endcomment%}visible`,
			data:     map[string]interface{}{},
			expected: "visible",
		},
		{
			name:     "nested access",
			template: `{%.user.name%}`,
			data:     map[string]interface{}{"user": map[string]interface{}{"name": "Bob"}},
			expected: "Bob",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.RenderString(tt.template, tt.data)
			if err != nil {
				t.Fatalf("RenderString failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestRenderStringErrors(t *testing.T) {
	e := New()

	_, err := e.RenderString(`{%invalid_tag%}`, nil)
	if err == nil {
		t.Error("Expected parse error for invalid tag")
	}
}

func TestRegisterFilter(t *testing.T) {
	e := New()
	e.RegisterFilter("shout", func(s interface{}) string {
		return strings.ToUpper(toString(s)) + "!!!"
	})
}

func TestRegisterTemplate(t *testing.T) {
	e := New()
	e.RegisterTemplate("greeting", `Hello {%.name%}!`)

	result, err := e.RenderFile("greeting", map[string]interface{}{"name": "World"})
	if err != nil {
		t.Fatalf("RenderFile failed: %v", err)
	}
	if result != "Hello World!" {
		t.Errorf("Expected 'Hello World!', got %q", result)
	}
}

func TestMemoryStore(t *testing.T) {
	ms := store.NewMemoryStore()
	ms.Register("test", `Hello`)

	content, err := ms.ReadTemplate("test")
	if err != nil {
		t.Fatalf("ReadTemplate failed: %v", err)
	}
	if string(content) != "Hello" {
		t.Errorf("Expected 'Hello', got %q", string(content))
	}

	_, err = ms.ReadTemplate("missing")
	if err == nil {
		t.Error("Expected error for missing template")
	}

	names, err := ms.ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}
	if len(names) != 1 || names[0] != "test" {
		t.Errorf("Expected ['test'], got %v", names)
	}
}

func TestAutoEscaping(t *testing.T) {
	e := New()

	result, err := e.RenderString(`<p>{%.html%}</p>`, map[string]interface{}{
		"html": "<script>alert('xss')</script>",
	})
	if err != nil {
		t.Fatalf("RenderString failed: %v", err)
	}
	if strings.Contains(result, "<script>") {
		t.Errorf("HTML should be escaped, got %q", result)
	}
	if !strings.Contains(result, "&lt;script&gt;") {
		t.Errorf("Expected escaped HTML, got %q", result)
	}
}

func TestAutoEscapingDisabled(t *testing.T) {
	e := NewUnsafe()

	result, err := e.RenderString(`<p>{%.html%}</p>`, map[string]interface{}{
		"html": "<b>bold</b>",
	})
	if err != nil {
		t.Fatalf("RenderString failed: %v", err)
	}
	if result != "<p><b>bold</b></p>" {
		t.Errorf("Expected unescaped HTML, got %q", result)
	}
}

func TestTemplateCaching(t *testing.T) {
	e := New()

	template := `{%.name%}`
	data := map[string]interface{}{"name": "Alice"}

	result1, err := e.RenderString(template, data)
	if err != nil {
		t.Fatalf("First render failed: %v", err)
	}

	result2, err := e.RenderString(template, data)
	if err != nil {
		t.Fatalf("Second render failed: %v", err)
	}

	if result1 != result2 {
		t.Errorf("Cached result differs: %q vs %q", result1, result2)
	}

	e.ClearCache()
	result3, err := e.RenderString(template, data)
	if err != nil {
		t.Fatalf("Third render failed: %v", err)
	}

	if result1 != result3 {
		t.Errorf("Result after cache clear differs: %q vs %q", result1, result3)
	}
}

func TestMaxIterations(t *testing.T) {
	e := New()
	e.SetMaxIterations(100)

	_, err := e.RenderString(`{%assign.x=true%}{%while.x%}loop{%end%}`, nil)
	if err == nil {
		t.Error("Expected iteration limit error")
	}
	if !strings.Contains(err.Error(), "iteration limit") {
		t.Errorf("Expected iteration limit error, got %v", err)
	}
}

func TestBuiltinFilters(t *testing.T) {
	e := New()
	_ = e
}
