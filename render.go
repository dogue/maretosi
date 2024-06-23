package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"io/fs"
	"os"

	"github.com/charmbracelet/log"
	"github.com/russross/blackfriday/v2"
)

//go:embed template.html
var TEMPLATE string

type TemplateData struct {
	Title string
	Body  template.HTML
}

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

func renderPage(inputFile, outputFile string) error {
	opts, body, err := extractFrontmatter(inputFile)
	if err != nil {
		return err
	}

	renderParams := blackfriday.HTMLRendererParameters{
		Flags: blackfriday.FootnoteReturnLinks,
	}
	renderer := blackfriday.NewHTMLRenderer(renderParams)

	var exts blackfriday.Option
	if disableExts {
		exts = blackfriday.WithNoExtensions()
	} else {
		exts = blackfriday.WithExtensions(blackfriday.CommonExtensions | blackfriday.Footnotes)
	}

	renderedBytes := blackfriday.Run([]byte(body), exts, blackfriday.WithRenderer(renderer))
	opts["body"] = template.HTML(renderedBytes)
	outputBuffer := bytes.Buffer{}

	t := template.New("maretosi")

	if templateFile, ok := opts["template"]; ok {
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

	if err = t.Execute(&outputBuffer, opts); err != nil {
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
