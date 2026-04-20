package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed skills/create.md
var skillCreate string

//go:embed skills/analyze.md
var skillAnalyze string

var skills = map[string]string{
	"create":  skillCreate,
	"analyze": skillAnalyze,
}

func cmdSkill(name string) error {
	// Handle aliases
	name = strings.TrimPrefix(name, "create-system-spec")
	name = strings.TrimPrefix(name, "analyze-system-spec")

	if name == "" || name == "create-system-spec" {
		name = "create"
	}
	if name == "analyze-system-spec" {
		name = "analyze"
	}

	content, ok := skills[name]
	if !ok {
		available := make([]string, 0, len(skills))
		for k := range skills {
			available = append(available, k)
		}
		return fmt.Errorf("unknown skill: %s (available: %s)", name, strings.Join(available, ", "))
	}

	fmt.Print(content)
	return nil
}

func listSkills() {
	fmt.Println("Available skills:")
	fmt.Println("  create   - Guide for creating a new system-spec")
	fmt.Println("  analyze  - Guide for analyzing/validating system-specs")
}
