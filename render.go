package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	chromaHTML "github.com/alecthomas/chroma/v2/formatters/html"

	gmd "github.com/yuin/goldmark"
	ghl "github.com/yuin/goldmark-highlighting/v2"
	gme "github.com/yuin/goldmark/extension"
	gmp "github.com/yuin/goldmark/parser"
	gmh "github.com/yuin/goldmark/renderer/html"

	fm "github.com/adrg/frontmatter"
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
	// copy non markdown files as-is
	if filepath.Ext(inputFile) != ".md" {
		return copyFile(inputFile, outputFile)
	}

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
	var templ string

	if templateFile, ok := metadata["__template"]; ok {
		templ = templateFile.(string)
	} else {
		templ = "__built_in"
	}

	if templates == nil {
		panic("templates are nil")
	}

	if err = templates.ExecuteTemplate(&outputBuf, templ, metadata); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err = os.WriteFile(outputFile, outputBuf.Bytes(), fs.ModePerm); err != nil {
		return fmt.Errorf("failed to write output file %q: %w", outputFile, err.(*fs.PathError).Err)
	}

	return nil
}

func copyFile(inputFile, outputFile string) error {
	fileData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read source file %q: %w", inputFile, err.(*fs.PathError).Err)
	}

	err = os.WriteFile(outputFile, fileData, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write destination file %q: %w", outputFile, err.(*fs.PathError).Err)
	}

	return nil
}

func copyAsset(inputFile, outputFile string) error {
	return copyFile(inputFile, outputFile)
}
