// Wisp Template Engine - Example HTTP Server
//
// An interactive web server that demonstrates the Wisp template engine
// by showing the output of each pipeline stage: lexer, parser, and renderer.
//
// Usage:
//
//	go run ./examples/server
//
// Then open http://localhost:8080 in your browser.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"

	"template-wisp/internal/ast"
	"template-wisp/internal/lexer"
	"template-wisp/internal/parser"
	"template-wisp/pkg/engine"
)

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.SetLevel(log.DebugLevel)
	log.SetReportTimestamp(true)

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/lexer", handleLexer)
	mux.HandleFunc("/parser", handleParser)
	mux.HandleFunc("/render", handleRender)
	mux.HandleFunc("/api/tokens", handleAPITokens)
	mux.HandleFunc("/api/ast", handleAPIAST)
	mux.HandleFunc("/api/render", handleAPIRender)

	addr := ":" + port
	server := &http.Server{
		Addr:         addr,
		Handler:      logMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info("Wisp Example Server starting", "addr", "http://localhost:"+port)
	log.Info("Routes available",
		"/", "Interactive UI",
		"/lexer", "Lexer output view",
		"/parser", "Parser AST view",
		"/render", "Rendered output view",
		"/api/tokens", "JSON token stream",
		"/api/ast", "JSON AST",
		"/api/render", "JSON render result",
	)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed", "error", err)
	}
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
		)
	})
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, indexHTML)
}

func handleLexer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, lexerHTML)
}

func handleParser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, parserHTML)
}

func handleRender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, renderHTML)
}

// API handlers

type TokenInfo struct {
	Type    string `json:"type"`
	Literal string `json:"literal"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
}

func handleAPITokens(w http.ResponseWriter, r *http.Request) {
	template := r.URL.Query().Get("template")
	if template == "" {
		template = `{% .name %}`
	}

	log.Debug("Tokenizing template", "template", template)

	l := lexer.NewLexer(template)
	var tokens []TokenInfo
	for {
		tok := l.NextToken()
		tokens = append(tokens, TokenInfo{
			Type:    string(tok.Type),
			Literal: tok.Literal,
			Line:    tok.Line,
			Column:  tok.Column,
		})
		if tok.Type == lexer.EOF {
			break
		}
	}

	log.Debug("Tokenization complete", "count", len(tokens))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"template": template,
		"tokens":   tokens,
	})
}

type ASTNode struct {
	Type     string    `json:"type"`
	String   string    `json:"string"`
	Children []ASTNode `json:"children,omitempty"`
}

func handleAPIAST(w http.ResponseWriter, r *http.Request) {
	template := r.URL.Query().Get("template")
	if template == "" {
		template = `{% .name %}`
	}

	log.Debug("Parsing template", "template", template)

	l := lexer.NewLexer(template)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	var astNodes []ASTNode
	for _, stmt := range program.Statements {
		astNodes = append(astNodes, astToNode(stmt))
	}

	log.Debug("Parsing complete",
		"statements", len(program.Statements),
		"errors", len(p.Errors()),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"template": template,
		"nodes":    astNodes,
		"errors":   p.Errors(),
	})
}

type RenderResult struct {
	Template string `json:"template"`
	Data     string `json:"data"`
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
}

func handleAPIRender(w http.ResponseWriter, r *http.Request) {
	template := r.URL.Query().Get("template")
	if template == "" {
		template = `<h1>Hello {% .name %}!</h1>`
	}
	dataJSON := r.URL.Query().Get("data")
	if dataJSON == "" {
		dataJSON = `{"name": "World"}`
	}

	log.Debug("Rendering template", "template", template, "data", dataJSON)

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
		log.Warn("Failed to parse data JSON", "error", err)
		data = make(map[string]interface{})
	}

	// Use the engine for rendering
	e := engine.New()
	output, err := e.RenderString(template, data)

	result := RenderResult{
		Template: template,
		Data:     dataJSON,
	}

	if err != nil {
		result.Error = err.Error()
		log.Warn("Render error", "error", err)
	} else {
		result.Output = output
		log.Debug("Render complete", "output_len", len(output))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func astToNode(stmt ast.Statement) ASTNode {
	node := ASTNode{
		Type:   fmt.Sprintf("%T", stmt),
		String: stmt.String(),
	}

	switch s := stmt.(type) {
	case *ast.IfStatement:
		if s.Consequence != nil {
			for _, cs := range s.Consequence.Statements {
				node.Children = append(node.Children, astToNode(cs))
			}
		}
	case *ast.ForStatement:
		if s.Body != nil {
			for _, bs := range s.Body.Statements {
				node.Children = append(node.Children, astToNode(bs))
			}
		}
	case *ast.CaseStatement:
		if s.Body != nil {
			for _, bs := range s.Body.Statements {
				node.Children = append(node.Children, astToNode(bs))
			}
		}
	}

	return node
}

// HTML pages

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Wisp Template Engine</title>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: 'SF Mono', 'Fira Code', monospace; background: #0d1117; color: #c9d1d9; }
.container { max-width: 1200px; margin: 0 auto; padding: 20px; }
h1 { color: #58a6ff; margin-bottom: 8px; font-size: 28px; }
.subtitle { color: #8b949e; margin-bottom: 24px; }
nav { display: flex; gap: 12px; margin-bottom: 24px; }
nav a { color: #58a6ff; text-decoration: none; padding: 8px 16px; border: 1px solid #30363d; border-radius: 6px; }
nav a:hover { background: #161b22; }
.card { background: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 20px; margin-bottom: 20px; }
label { display: block; color: #8b949e; margin-bottom: 8px; font-size: 13px; }
textarea { width: 100%; height: 120px; background: #0d1117; color: #c9d1d9; border: 1px solid #30363d; border-radius: 6px; padding: 12px; font-family: inherit; font-size: 14px; resize: vertical; }
textarea:focus { outline: none; border-color: #58a6ff; }
button { background: #238636; color: white; border: none; padding: 10px 20px; border-radius: 6px; cursor: pointer; font-family: inherit; font-size: 14px; margin-top: 12px; }
button:hover { background: #2ea043; }
.output { background: #0d1117; border: 1px solid #30363d; border-radius: 6px; padding: 16px; margin-top: 16px; white-space: pre-wrap; min-height: 60px; font-size: 13px; }
.output .token { display: inline-block; margin: 2px 4px; padding: 2px 8px; border-radius: 4px; font-size: 12px; }
.tok-keyword { background: #1f6feb33; color: #58a6ff; }
.tok-string { background: #23863622; color: #7ee787; }
.tok-number { background: #a371f722; color: #d2a8ff; }
.tok-ident { background: #da363322; color: #ff7b72; }
.tok-dot { background: #d2992222; color: #e3b341; }
.tok-text { color: #8b949e; }
.tok-other { background: #30363d; color: #c9d1d9; }
.error { color: #f85149; }
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
@media (max-width: 800px) { .grid { grid-template-columns: 1fr; } }
.badge { display: inline-block; background: #1f6feb; color: white; padding: 2px 8px; border-radius: 10px; font-size: 11px; margin-left: 8px; }
</style>
</head>
<body>
<div class="container">
<h1>Wisp Template Engine <span class="badge">v0.1.0</span></h1>
<p class="subtitle">Interactive template pipeline explorer</p>
<nav>
<a href="/lexer">Lexer</a>
<a href="/parser">Parser</a>
<a href="/render">Render</a>
</nav>

<div class="card">
<label>Template</label>
<textarea id="template">{% for .item in .items %}
  <li>{% .item.name %}</li>
{% end %}</textarea>
<label>Data (JSON)</label>
<textarea id="data">{"items": [{"name": "Apple"}, {"name": "Banana"}, {"name": "Cherry"}]}</textarea>
<button onclick="runAll()">Run Pipeline</button>
</div>

<div class="grid">
<div class="card">
<label>Tokens (Lexer Output)</label>
<div id="tokens" class="output"></div>
</div>
<div class="card">
<label>AST (Parser Output)</label>
<div id="ast" class="output"></div>
</div>
</div>

<div class="card">
<label>Rendered Output</label>
<div id="rendered" class="output"></div>
</div>
</div>

<script>
async function runAll() {
    const template = document.getElementById('template').value;
    const data = document.getElementById('data').value;

    // Fetch tokens
    const tokRes = await fetch('/api/tokens?template=' + encodeURIComponent(template));
    const tokData = await tokRes.json();
    document.getElementById('tokens').innerHTML = tokData.tokens.map(t =>
        '<span class="token tok-' + tokenClass(t.type) + '" title="L' + t.line + ':' + t.column + '">' + escapeHtml(t.literal || t.type) + '</span>'
    ).join('');

    // Fetch AST
    const astRes = await fetch('/api/ast?template=' + encodeURIComponent(template));
    const astData = await astRes.json();
    let astHtml = '';
    if (astData.errors && astData.errors.length > 0) {
        astHtml = '<span class="error">Errors: ' + escapeHtml(astData.errors.join(', ')) + '</span>\n';
    }
    astHtml += astData.nodes.map(n => renderASTNode(n, 0)).join('');
    document.getElementById('ast').innerHTML = astHtml;

    // Fetch render
    const renRes = await fetch('/api/render?template=' + encodeURIComponent(template) + '&data=' + encodeURIComponent(data));
    const renData = await renRes.json();
    if (renData.error) {
        document.getElementById('rendered').innerHTML = '<span class="error">' + escapeHtml(renData.error) + '</span>';
    } else {
        document.getElementById('rendered').textContent = renData.output;
    }
}

function tokenClass(type) {
    if (['IF','ELSE','ELSIF','UNLESS','FOR','WHILE','RANGE','CASE','WHEN','WITH','END','BREAK','CONTINUE','INCLUDE','RENDER','COMPONENT','EXTENDS','BLOCK','CONTENT','RAW','COMMENT','ASSIGN','LET','RETURN','FUNCTION','AS','IN','CYCLE','INCREMENT','DECREMENT','ENDRAW','ENDCOMMENT','TRUE','FALSE'].includes(type)) return 'keyword';
    if (type === 'STRING') return 'string';
    if (type === 'NUMBER') return 'number';
    if (type === 'IDENT') return 'ident';
    if (type === 'DOT') return 'dot';
    if (type === 'TEXT') return 'text';
    return 'other';
}

function renderASTNode(node, depth) {
    const indent = '  '.repeat(depth);
    let html = indent + '<span class="tok-keyword">' + escapeHtml(node.type.split('.').pop()) + '</span>';
    if (node.string) html += ' <span class="tok-string">' + escapeHtml(node.string) + '</span>';
    html += '\n';
    if (node.children) {
        html += node.children.map(c => renderASTNode(c, depth + 1)).join('');
    }
    return html;
}

function escapeHtml(s) {
    return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}

runAll();
</script>
</body>
</html>`

const lexerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Wisp - Lexer</title>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: 'SF Mono', monospace; background: #0d1117; color: #c9d1d9; padding: 20px; }
h1 { color: #58a6ff; margin-bottom: 16px; }
a { color: #58a6ff; }
textarea { width: 100%; height: 100px; background: #0d1117; color: #c9d1d9; border: 1px solid #30363d; border-radius: 6px; padding: 12px; font-family: inherit; }
button { background: #238636; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; margin: 12px 0; }
table { width: 100%; border-collapse: collapse; margin-top: 12px; }
th, td { text-align: left; padding: 8px 12px; border: 1px solid #30363d; }
th { background: #161b22; color: #58a6ff; }
tr:nth-child(even) { background: #161b22; }
.kw { color: #58a6ff; } .str { color: #7ee787; } .num { color: #d2a8ff; } .id { color: #ff7b72; }
</style>
</head>
<body>
<a href="/">Back</a>
<h1>Lexer Output</h1>
<textarea id="t">{% for .item in .items %}{% .item %}{% end %}</textarea>
<button onclick="run()">Tokenize</button>
<table><thead><tr><th>#</th><th>Type</th><th>Literal</th><th>Line</th><th>Col</th></tr></thead><tbody id="o"></tbody></table>
<script>
async function run() {
    const r = await fetch('/api/tokens?template=' + encodeURIComponent(document.getElementById('t').value));
    const d = await r.json();
    document.getElementById('o').innerHTML = d.tokens.map((t,i) =>
        '<tr><td>'+i+'</td><td class="'+cls(t.type)+'">'+esc(t.type)+'</td><td>'+esc(t.literal)+'</td><td>'+t.line+'</td><td>'+t.column+'</td></tr>'
    ).join('');
}
function cls(t) { return ['IF','ELSE','FOR','END','BREAK','CONTINUE','INCLUDE','RENDER','COMPONENT','EXTENDS','BLOCK','CONTENT','RAW','COMMENT','ASSIGN','WHEN','WITH','CASE','UNLESS','WHILE','RANGE','CYCLE','INCREMENT','DECREMENT','ELSIF','LET'].includes(t)?'kw':t==='STRING'?'str':t==='NUMBER'?'num':t==='IDENT'?'id':''; }
function esc(s) { return s.replace(/&/g,'&amp;').replace(/</g,'&lt;'); }
run();
</script>
</body>
</html>`

const parserHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Wisp - Parser</title>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: 'SF Mono', monospace; background: #0d1117; color: #c9d1d9; padding: 20px; }
h1 { color: #58a6ff; margin-bottom: 16px; }
a { color: #58a6ff; }
textarea { width: 100%; height: 100px; background: #0d1117; color: #c9d1d9; border: 1px solid #30363d; border-radius: 6px; padding: 12px; font-family: inherit; }
button { background: #238636; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; margin: 12px 0; }
#out { background: #0d1117; border: 1px solid #30363d; border-radius: 6px; padding: 16px; white-space: pre; min-height: 100px; }
.err { color: #f85149; }
.tp { color: #58a6ff; } .str { color: #7ee787; }
</style>
</head>
<body>
<a href="/">Back</a>
<h1>Parser AST</h1>
<textarea id="t">{% if .show %}<h1>{% .title %}</h1>{% else %}<p>Hidden</p>{% end %}</textarea>
<button onclick="run()">Parse</button>
<div id="out"></div>
<script>
async function run() {
    const r = await fetch('/api/ast?template=' + encodeURIComponent(document.getElementById('t').value));
    const d = await r.json();
    let h = '';
    if (d.errors.length) h += '<span class="err">Errors: ' + d.errors.join(', ') + '</span>\n';
    h += d.nodes.map(n => rn(n,0)).join('');
    document.getElementById('out').innerHTML = h;
}
function rn(n,d) {
    const ind = '  '.repeat(d);
    let h = ind + '<span class="tp">' + esc(n.type.split('.').pop()) + '</span>';
    if (n.string) h += ' <span class="str">' + esc(n.string) + '</span>';
    h += '\n';
    if (n.children) h += n.children.map(c => rn(c,d+1)).join('');
    return h;
}
function esc(s) { return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); }
run();
</script>
</body>
</html>`

const renderHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Wisp - Render</title>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: 'SF Mono', monospace; background: #0d1117; color: #c9d1d9; padding: 20px; }
h1 { color: #58a6ff; margin-bottom: 16px; }
a { color: #58a6ff; }
label { display: block; color: #8b949e; margin: 12px 0 4px; font-size: 13px; }
textarea { width: 100%; height: 100px; background: #0d1117; color: #c9d1d9; border: 1px solid #30363d; border-radius: 6px; padding: 12px; font-family: inherit; }
button { background: #238636; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; margin: 12px 0; }
#out { background: #0d1117; border: 1px solid #30363d; border-radius: 6px; padding: 16px; white-space: pre-wrap; min-height: 100px; }
.err { color: #f85149; }
.examples { margin: 16px 0; }
.examples button { background: #30363d; font-size: 12px; padding: 4px 10px; margin: 2px; }
</style>
</head>
<body>
<a href="/">Back</a>
<h1>Render Template</h1>
<div class="examples">
<button onclick="ex('hello')">Hello World</button>
<button onclick="ex('loop')">For Loop</button>
<button onclick="ex('cond')">Conditional</button>
<button onclick="ex('case')">Case/When</button>
<button onclick="ex('nested')">Nested Access</button>
</div>
<label>Template</label>
<textarea id="t"><h1>Hello {% .name %}!</h1></textarea>
<label>Data (JSON)</label>
<textarea id="d">{"name": "World"}</textarea>
<button onclick="run()">Render</button>
<div id="out"></div>
<script>
const examples = {
    hello: { t: '<h1>Hello {% .name %}!</h1>', d: '{"name": "World"}' },
    loop: { t: '{% for .item in .items %}<li>{% .item %}</li>{% end %}', d: '{"items": ["Go", "Rust", "Zig"]}' },
    cond: { t: '{% if .loggedIn %}<p>Welcome, {% .user %}!</p>{% else %}<p>Please log in.</p>{% end %}', d: '{"loggedIn": true, "user": "Alice"}' },
    case: { t: '{% case .status %}{% when "active" %}Online{% when "away" %}Away{% else %}Offline{% end %}', d: '{"status": "active"}' },
    nested: { t: '{% .user.profile.name %} ({% .user.profile.age %})', d: '{"user": {"profile": {"name": "Bob", "age": 30}}}' },
};
function ex(name) {
    document.getElementById('t').value = examples[name].t;
    document.getElementById('d').value = examples[name].d;
    run();
}
async function run() {
    const t = document.getElementById('t').value;
    const d = document.getElementById('d').value;
    const r = await fetch('/api/render?template=' + encodeURIComponent(t) + '&data=' + encodeURIComponent(d));
    const res = await r.json();
    if (res.error) document.getElementById('out').innerHTML = '<span class="err">' + esc(res.error) + '</span>';
    else document.getElementById('out').textContent = res.output;
}
function esc(s) { return s.replace(/&/g,'&amp;').replace(/</g,'&lt;'); }
ex('hello');
</script>
</body>
</html>`
