package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"

	chromaHTML "github.com/alecthomas/chroma/v2/formatters/html"

	gmd "github.com/yuin/goldmark"
	ghl "github.com/yuin/goldmark-highlighting/v2"
	gme "github.com/yuin/goldmark/extension"
	gmp "github.com/yuin/goldmark/parser"
	gmh "github.com/yuin/goldmark/renderer/html"

	fm "github.com/adrg/frontmatter"
)

//go:embed template.html
var TEMPLATE string

type RenderErr = int

const (
	unknownErr RenderErr = iota
)

func walker(destination *[]string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to read path %q: %w", path, err.(*fs.PathError).Err)
		}

		if d.IsDir() {
			return nil
		}

		*destination = append(*destination, path)
		return nil
	}
}

func initRenderer(syntaxTheme string, lineNumbers bool) gmd.Markdown {
	exts := gmd.WithExtensions(
		gme.Table,
		gme.Strikethrough,
		gme.Footnote,
		ghl.NewHighlighting(
			ghl.WithStyle(syntaxTheme),
			ghl.WithFormatOptions(
				chromaHTML.WithLineNumbers(lineNumbers),
			),
		),
	)
	pOpts := gmd.WithParserOptions(gmp.WithAutoHeadingID())
	rOpts := gmd.WithRendererOptions(gmh.WithUnsafe())
	return gmd.New(exts, pOpts, rOpts)
}

func renderPage(inputFile, outputFile string) error {
	metadata := make(map[string]any)
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", inputFile, err.(*fs.PathError).Err)
	}

	body, err := fm.Parse(file, &metadata)
	if err != nil {
		return fmt.Errorf("failed to parse front matter from source file %q: %w", inputFile, err)
	}

	var syntaxTheme string
	var lineNumbers bool

	if theme, ok := metadata["__syntax_theme"]; ok {
		syntaxTheme = theme.(string)
	} else {
		syntaxTheme = "github-dark"
	}

	if numbers, ok := metadata["__line_numbers"]; ok {
		lineNumbers = numbers.(bool)
	} else {
		lineNumbers = false
	}

	renderer := initRenderer(syntaxTheme, lineNumbers)

	var renderedBytes bytes.Buffer
	renderer.Convert([]byte(body), &renderedBytes)
	metadata["__body"] = template.HTML(renderedBytes.String())

	outputBuf := bytes.Buffer{}
	t := template.New("maretosi")

	if templateFile, ok := metadata["__template"]; ok {
		file, err := os.ReadFile(templateFile.(string))
		if err != nil {
			return fmt.Errorf("failed to read template file %q: %w", templateFile, err)
		}

		if _, err := t.Parse(string(file)); err != nil {
			return fmt.Errorf("failed to parse template file %q: %w", templateFile, err)
		}
	} else {
		if _, err = t.Parse(TEMPLATE); err != nil {
			return fmt.Errorf("failed to parse built-in template. This is a bug. Please consider reporting it at https://github.com/dogue/maretosi/issues")
		}
	}

	if err = t.Execute(&outputBuf, metadata); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err = os.WriteFile(outputFile, outputBuf.Bytes(), fs.ModePerm); err != nil {
		return fmt.Errorf("failed to write output file %q: %w", outputFile, err.(*fs.PathError).Err)
	}

	return nil
}

func copyAsset(inputFile, outputFile string) error {
	fileData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read asset file %q: %w", inputFile, err.(*fs.PathError).Err)
	}

	err = os.WriteFile(outputFile, fileData, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write asset file %q: %w", outputFile, err.(*fs.PathError).Err)
	}

	return nil
}
