package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"io/fs"
	"os"

	// "github.com/alecthomas/chroma/v2"
	chromaHTML "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/charmbracelet/log"

	gmd "github.com/yuin/goldmark"
	ghl "github.com/yuin/goldmark-highlighting/v2"
	gme "github.com/yuin/goldmark/extension"
	gmp "github.com/yuin/goldmark/parser"
	gmh "github.com/yuin/goldmark/renderer/html"

	fm "github.com/adrg/frontmatter"
)

//go:embed template.html
var TEMPLATE string

func walker(destination *[]string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error(err)
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
		return err
	}

	body, err := fm.Parse(file, &metadata)
	if err != nil {
		return err
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
	// renderedBytes := blackfriday.Run([]byte(body), exts, blackfriday.WithRenderer(renderer))
	metadata["__body"] = template.HTML(renderedBytes.String())
	outputBuffer := bytes.Buffer{}

	t := template.New("maretosi")

	if templateFile, ok := metadata["__template"]; ok {
		file, err := os.ReadFile(templateFile.(string))
		if err != nil {
			return err
		}

		if _, err := t.Parse(string(file)); err != nil {
			return err
		}
	} else {
		if _, err = t.Parse(TEMPLATE); err != nil {
			return err
		}
	}

	if err = t.Execute(&outputBuffer, metadata); err != nil {
		return err
	}

	if err = os.WriteFile(outputFile, outputBuffer.Bytes(), fs.ModePerm); err != nil {
		return err
	}

	return nil
}

func copyAsset(inputFile, outputFile string) error {
	fileData, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(outputFile, fileData, fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
