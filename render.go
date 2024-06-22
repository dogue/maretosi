package main

import (
	"bytes"
	_ "embed"
	// "fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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
			slog.Error(err.Error())
		}

		if d.IsDir() {
			return nil
		}

		*destination = append(*destination, path)
		return nil
	}
}

func render_all() (err error) {
	var pages []string
	filepath.WalkDir(input_dir, walker(&pages))

	for _, page := range pages {
		render_page(page)
	}

	return
}

func render_page(input_page string) (err error) {
	output_page := strings.Replace(input_page, filepath.Base(input_dir), filepath.Base(output_dir), 1)
	output_page = strings.Replace(output_page, ".md", ".html", 1)
	output_page_parent := filepath.Dir(output_page)

	err = os.MkdirAll(output_page_parent, fs.ModePerm)
	if err != nil {
		return
	}

	input_bytes, err := os.ReadFile(input_page)
	if err != nil {
		return
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
		return
	}

	if err = t.Execute(&out_buf, data); err != nil {
		return
	}

	err = os.WriteFile(output_page, out_buf.Bytes(), fs.ModePerm)
	return
}

func copy_assets() (err error) {
	var assets []string
	err = filepath.WalkDir(assets_dir, walker(&assets))

	for _, asset := range assets {
		copy_asset(asset)
	}

	return
}

func copy_asset(input_path string) (err error) {
	output_path := strings.Replace(input_path, filepath.Base(assets_dir), filepath.Join(filepath.Base(output_dir), "static"), 1)
	output_path_parent := filepath.Dir(output_path)

	err = os.MkdirAll(output_path_parent, fs.ModePerm)
	if err != nil {
		return
	}

	input, err := os.ReadFile(input_path)
	if err != nil {
		return
	}

	err = os.WriteFile(output_path, input, fs.ModePerm)
	return
}
