// Command system-spec provides CLI tools for working with system specifications.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/plexusone/system-spec/graph"
	"github.com/plexusone/system-spec/render"
	"github.com/plexusone/system-spec/schema"
	"github.com/plexusone/system-spec/spec"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "validate":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: system-spec validate [--strict] <file.json>")
			os.Exit(1)
		}
		strict := false
		filePath := ""
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "--strict" || os.Args[i] == "-s" {
				strict = true
			} else if !strings.HasPrefix(os.Args[i], "-") {
				filePath = os.Args[i]
			}
		}
		if filePath == "" {
			fmt.Fprintln(os.Stderr, "usage: system-spec validate [--strict] <file.json>")
			os.Exit(1)
		}
		if err := cmdValidate(filePath, strict); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "lint":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: system-spec lint <file.json>")
			os.Exit(1)
		}
		if err := cmdLint(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "render":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "usage: system-spec render <file.json> --format <format>")
			os.Exit(1)
		}
		format := ""
		for i := 3; i < len(os.Args); i++ {
			if os.Args[i] == "--format" || os.Args[i] == "-f" {
				if i+1 < len(os.Args) {
					format = os.Args[i+1]
				}
			}
		}
		if format == "" {
			fmt.Fprintln(os.Stderr, "error: --format is required")
			os.Exit(1)
		}
		if err := cmdRender(os.Args[2], format); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "schema":
		if err := cmdSchema(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "graph":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: system-spec graph <file.json>")
			os.Exit(1)
		}
		if err := cmdGraph(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "skill":
		if len(os.Args) < 3 {
			listSkills()
			os.Exit(0)
		}
		if err := cmdSkill(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`system-spec - System topology specification tool

Usage:
  system-spec <command> [arguments]

Commands:
  validate [--strict] <file.json>   Validate a system specification file
  lint <file.json>                  Check for potential issues (warnings)
  render <file.json> --format <fmt> Render to diagram format (d2, mermaid, cytoscape, sigma, dot)
  graph <file.json>                 Output intermediate graph as JSON
  schema                            Output JSON Schema for system-spec
  skill <name>                      Output AI agent skill instructions (create, analyze)

Flags:
  --strict, -s    Treat lint warnings as errors (for validate command)

Examples:
  system-spec validate system.json
  system-spec validate --strict system.json  # Fail if any lint warnings
  system-spec lint system.json               # Show warnings only
  system-spec render system.json --format d2 > system.d2
  system-spec render system.json --format mermaid > system.mmd
  system-spec render system.json --format cytoscape > system.cyto.json
  system-spec schema > system.schema.json
  system-spec skill create          # Get instructions for creating a spec`)
}

func cmdValidate(path string, strict bool) error {
	sys, err := spec.LoadFromFile(path)
	if err != nil {
		return err
	}

	fmt.Printf("valid: %s (%d services)\n", sys.Name, len(sys.Services))

	// Print service summary
	for name, svc := range sys.Services {
		fmt.Printf("  - %s: %s\n", name, svc.Image.FullName())
		if svc.Repo != nil {
			fmt.Printf("      repo: %s\n", svc.Repo.URL)
		}
		if len(svc.Connections) > 0 {
			fmt.Printf("      connections: %d\n", len(svc.Connections))
		}
	}

	// Run lint checks
	warnings := sys.Lint()
	if len(warnings) > 0 {
		fmt.Printf("\nwarnings: %d\n", len(warnings))
		for _, w := range warnings {
			if w.Service != "" {
				fmt.Printf("  [%s] %s: %s\n", w.Code, w.Service, w.Message)
			} else {
				fmt.Printf("  [%s] %s\n", w.Code, w.Message)
			}
		}

		if strict {
			return fmt.Errorf("%d lint warnings (strict mode)", len(warnings))
		}
	}

	return nil
}

func cmdLint(path string) error {
	sys, err := spec.LoadFromFile(path)
	if err != nil {
		return err
	}

	warnings := sys.Lint()

	if len(warnings) == 0 {
		fmt.Printf("no warnings: %s (%d services)\n", sys.Name, len(sys.Services))
		return nil
	}

	fmt.Printf("%s: %d warnings\n\n", sys.Name, len(warnings))

	// Group warnings by code for better readability
	byCode := make(map[string][]spec.Warning)
	for _, w := range warnings {
		byCode[w.Code] = append(byCode[w.Code], w)
	}

	// Print warnings grouped by code
	for code, ws := range byCode {
		fmt.Printf("[%s] (%d)\n", code, len(ws))
		for _, w := range ws {
			if w.Service != "" {
				fmt.Printf("  %s: %s\n", w.Service, w.Message)
			} else {
				fmt.Printf("  %s\n", w.Message)
			}
		}
		fmt.Println()
	}

	return nil
}

func cmdRender(path string, format string) error {
	sys, err := spec.LoadFromFile(path)
	if err != nil {
		return err
	}

	g := graph.FromSystem(sys)
	renderers := render.NewRenderers()

	r := renderers.Get(render.Format(format))
	if r == nil {
		return fmt.Errorf("unknown format: %s (supported: d2, mermaid, cytoscape, sigma, dot)", format)
	}

	output, err := r.Render(g)
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	return nil
}

func cmdSchema() error {
	data, err := schema.Generate()
	if err != nil {
		return err
	}

	fmt.Print(string(data))
	return nil
}

func cmdGraph(path string) error {
	sys, err := spec.LoadFromFile(path)
	if err != nil {
		return err
	}

	g := graph.FromSystem(sys)

	output, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	return nil
}
