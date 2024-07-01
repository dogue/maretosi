package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
)

//go:embed template.html
var TEMPLATE string

// var templates map[string]*template.Template
var templates *template.Template

func parseTemplates(inputFiles []string) error {
	// templates = make(map[string]*template.Template)
	templates = template.New("maretosi")

	for _, templateFile := range inputFiles {
		// name := extractSubpath(templateFile, templDir)
		// t := template.New(name)

		file, err := os.ReadFile(templateFile)
		if err != nil {
			return fmt.Errorf("failed to read template file %q: %w", templateFile, err)
		}

		if _, err := templates.Parse(string(file)); err != nil {
			return fmt.Errorf("failed to parse template file %q: %w", templateFile, err)
		}

		// templates[name] = t
	}

	// t := template.New("__built_in")
	if _, err := templates.Parse(TEMPLATE); err != nil {
		return fmt.Errorf("failed to parse built-in template. This is a bug. Please consider reporting it at https://github.com/dogue/maretosi/issues")
	}

	// templates["__built_in"] = t

	return nil
}
