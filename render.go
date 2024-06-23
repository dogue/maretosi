package main

import (
	"bytes"
	_ "embed"
	"errors"

	// "fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

func renderAll() {
	if _, err := os.Stat(inputDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("content source directory %q not found", inputDir)
		}
		log.Error(err)
	}

	var pages []string
	filepath.WalkDir(inputDir, walker(&pages))

	if len(pages) < 1 {
		log.Fatalf("no source content found in directory %q", inputDir)
	}

	for _, page := range pages {
		renderPage(page)
	}
}

func renderPage(inputPage string) {
	outputPage := strings.Replace(inputPage, filepath.Base(inputDir), filepath.Base(outputDir), 1)
	outputPage = strings.Replace(outputPage, ".md", ".html", 1)
	outputPageParent := filepath.Dir(outputPage)

	err := os.MkdirAll(outputPageParent, fs.ModePerm)
	if err != nil {
		log.Fatalf("failed to create output directory %q: %v", outputPageParent, err)
	}

	input_bytes, err := os.ReadFile(inputPage)
	if err != nil {
		log.Fatalf("failed to read source file %q: %v", inputPage, err)
	}

	var exts blackfriday.Option
	if disableExts {
		exts = blackfriday.WithNoExtensions()
	} else {
		exts = blackfriday.WithExtensions(blackfriday.CommonExtensions)
	}

	renderedBytes := blackfriday.Run(input_bytes, exts)
	out_buf := bytes.Buffer{}
	data := TemplateData{
		Title: siteTitle,
		Body:  template.HTML(renderedBytes),
	}

	t := template.New("maretosi")
	if _, err = t.Parse(TEMPLATE); err != nil {
		log.Fatalf("failed to parse HTML template: %v", err)
	}

	if err = t.Execute(&out_buf, data); err != nil {
		log.Fatalf("failed to execute HTML template: %v", err)
	}

	if err = os.WriteFile(outputPage, out_buf.Bytes(), fs.ModePerm); err != nil {
		log.Fatalf("failed to write output file %q: %v", outputPage, err)
	}
}

func copyAssets() {
	if _, err := os.Stat(assetsDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Warnf("assets source directory %q not found", inputDir)
		} else {
			log.Error(err)
		}

		return
	}

	var assets []string
	filepath.WalkDir(assetsDir, walker(&assets))

	for i, asset := range assets {
		if err := copyAsset(asset); err != nil {
			log.Errorf("error copying asset %q\n\tcopied %d of %d assets successfully", asset, i, len(assets))
			break
		}
	}

	return
}

func copyAsset(inputPath string) error {
	outputPath := strings.Replace(inputPath, filepath.Base(assetsDir), filepath.Join(filepath.Base(outputDir), "static"), 1)
	outputPathParent := filepath.Dir(outputPath)

	err := os.MkdirAll(outputPathParent, fs.ModePerm)
	if err != nil {
		log.Errorf("failed to create assets output directory %q: %v", outputPathParent, err)
		return err
	}

	input, err := os.ReadFile(inputPath)
	if err != nil {
		log.Errorf("failed to read asset file %q:\n\t%v", inputPath, err)
		return err
	}

	err = os.WriteFile(outputPath, input, fs.ModePerm)
	if err != nil {
		log.Errorf("failed to write asset file %q:\n\t%v", outputPath, err)
		return err
	}

	return nil
}
