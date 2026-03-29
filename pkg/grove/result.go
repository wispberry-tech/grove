// pkg/grove/result.go
package grove

import "strings"

// RenderResult holds the output of a render operation.
type RenderResult struct {
	Body   string
	Assets AssetBundle
	Meta   map[string]any
}

// AssetBundle holds collected CSS/JS assets (populated in Plan 6).
type AssetBundle struct {
	Scripts  []Asset
	Styles   []Asset
	Preloads []Asset
}

// Asset represents a single CSS or JS reference (populated in Plan 6).
type Asset struct {
	Src     string
	Content string
	Attrs   map[string]string
}

// InjectAssets inserts collected assets before </head>. No-op until Plan 6.
func (r RenderResult) InjectAssets() string {
	if len(r.Assets.Scripts) == 0 && len(r.Assets.Styles) == 0 {
		return r.Body
	}
	idx := strings.Index(r.Body, "</head>")
	if idx < 0 {
		return r.Body
	}
	// Plan 6 implements full injection
	return r.Body
}
