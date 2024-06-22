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

func render_all() {
	if _, err := os.Stat(input_dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("content source directory %q not found", input_dir)
		}
		log.Error(err)
	}

	var pages []string
	filepath.WalkDir(input_dir, walker(&pages))

	if len(pages) < 1 {
		log.Fatalf("no source content found in directory %q", input_dir)
	}

	for _, page := range pages {
		render_page(page)
	}
}

func render_page(input_page string) {
	output_page := strings.Replace(input_page, filepath.Base(input_dir), filepath.Base(output_dir), 1)
	output_page = strings.Replace(output_page, ".md", ".html", 1)
	output_page_parent := filepath.Dir(output_page)

	err := os.MkdirAll(output_page_parent, fs.ModePerm)
	if err != nil {
		log.Fatalf("failed to create output directory %q: %v", output_page_parent, err)
	}

	input_bytes, err := os.ReadFile(input_page)
	if err != nil {
		log.Fatalf("failed to read source file %q: %v", input_page, err)
	}

	var exts blackfriday.Option
	if disable_exts {
		exts = blackfriday.WithNoExtensions()
	} else {
		exts = blackfriday.WithExtensions(blackfriday.CommonExtensions)
	}

	rendered_bytes := blackfriday.Run(input_bytes, exts)
	out_buf := bytes.Buffer{}
	data := TemplateData{
		Title: site_title,
		Body:  template.HTML(rendered_bytes),
	}

	t := template.New("maretosi")
	if _, err = t.Parse(TEMPLATE); err != nil {
		log.Fatalf("failed to parse HTML template: %v", err)
	}

	if err = t.Execute(&out_buf, data); err != nil {
		log.Fatalf("failed to execute HTML template: %v", err)
	}

	if err = os.WriteFile(output_page, out_buf.Bytes(), fs.ModePerm); err != nil {
		log.Fatalf("failed to write output file %q: %v", output_page, err)
	}
}

func copy_assets() {
	if _, err := os.Stat(assets_dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Warnf("assets source directory %q not found", input_dir)
		} else {
			log.Error(err)
		}

		return
	}

	var assets []string
	filepath.WalkDir(assets_dir, walker(&assets))

	for i, asset := range assets {
		if err := copy_asset(asset); err != nil {
			log.Errorf("error copying asset %q\n\tcopied %d of %d assets successfully", asset, i, len(assets))
			break
		}
	}

	return
}

func copy_asset(input_path string) error {
	output_path := strings.Replace(input_path, filepath.Base(assets_dir), filepath.Join(filepath.Base(output_dir), "static"), 1)
	output_path_parent := filepath.Dir(output_path)

	err := os.MkdirAll(output_path_parent, fs.ModePerm)
	if err != nil {
		log.Errorf("failed to create assets output directory %q: %v", output_path_parent, err)
		return err
	}

	input, err := os.ReadFile(input_path)
	if err != nil {
		log.Errorf("failed to read asset file %q:\n\t%v", input_path, err)
		return err
	}

	err = os.WriteFile(output_path, input, fs.ModePerm)
	if err != nil {
		log.Errorf("failed to write asset file %q:\n\t%v", output_path, err)
		return err
	}

	return nil
}
