// Wisp is a secure, high-performance HTML templating engine for Go.
//
// Usage:
//
//	wisp <command> [arguments]
//
// The commands are:
//
//	render    Render a template with data
//	validate  Validate template syntax
//	version   Print version information
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"template-wisp/pkg/engine"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "render":
		cmdRender(os.Args[2:])
	case "validate":
		cmdValidate(os.Args[2:])
	case "version":
		fmt.Printf("wisp %s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Wisp Template Engine

Usage:
  wisp <command> [arguments]

Commands:
  render <template> [data]   Render a template with optional JSON data
  validate <template>        Validate template syntax
  version                    Print version information
  help                       Show this help message

Examples:
  wisp render template.wisp '{"name": "World"}'
  wisp validate template.wisp
  echo '{% .name %}' | wisp render - '{"name": "Alice"}'`)
}

func cmdRender(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: wisp render <template> [data]")
		os.Exit(1)
	}

	templatePath := args[0]

	// Read template content
	var templateContent string
	if templatePath == "-" {
		// Read from stdin
		data, err := os.ReadFile("/dev/stdin")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
			os.Exit(1)
		}
		templateContent = string(data)
	} else {
		data, err := os.ReadFile(templatePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading template: %v\n", err)
			os.Exit(1)
		}
		templateContent = string(data)
	}

	// Parse data
	dataMap := make(map[string]interface{})
	if len(args) > 1 {
		dataJSON := strings.Join(args[1:], " ")
		if err := json.Unmarshal([]byte(dataJSON), &dataMap); err != nil {
			fmt.Fprintf(os.Stderr, "error parsing data JSON: %v\n", err)
			os.Exit(1)
		}
	}

	// Render
	e := engine.New()
	result, err := e.RenderString(templateContent, dataMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}

func cmdValidate(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: wisp validate <template>")
		os.Exit(1)
	}

	templatePath := args[0]

	// Read template content
	var templateContent string
	if templatePath == "-" {
		data, err := os.ReadFile("/dev/stdin")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
			os.Exit(1)
		}
		templateContent = string(data)
	} else {
		data, err := os.ReadFile(templatePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading template: %v\n", err)
			os.Exit(1)
		}
		templateContent = string(data)
	}

	// Validate by parsing only (no evaluation)
	e := engine.New()
	err := e.Validate(templateContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("template is valid")
}
