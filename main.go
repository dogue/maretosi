package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

var contentDir string
var outputDir string
var assetsDir string
var templDir string
var skipAssets bool
var verbose bool

func main() {
	cli, err := parseCli()
	if err != nil {
		log.Fatal(err)
	}

	if cli.IsSet("help") {
		os.Exit(0)
	}

	if err = validateDirs(); err != nil {
		log.Fatal(err)
	}

	var templateFiles []string
	if err := filepath.WalkDir(templDir, walker(&templateFiles)); err != nil {
		log.Fatal(err)
	}

	if err := parseTemplates(templateFiles); err != nil {
		log.Fatal(err)
	}

	var contentFiles []string
	if err := filepath.WalkDir(contentDir, walker(&contentFiles)); err != nil {
		log.Fatal(err)
	}

	for _, contentPath := range contentFiles {
		subpath := extractSubpath(contentPath, contentDir)
		outputFile := filepath.Join(outputDir, mdToHtml(subpath))

		parentDir := filepath.Dir(outputFile)
		if err := os.MkdirAll(parentDir, fs.ModePerm); err != nil {
			log.Fatal(err)
		}

		err = renderPage(contentPath, outputFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if skipAssets {
		return
	}

	var assetsFiles []string
	if err := filepath.WalkDir(assetsDir, walker(&assetsFiles)); err != nil {
		log.Fatal(err)
	}

	for _, asset := range assetsFiles {
		subpath := extractSubpath(asset, assetsDir)
		outputFile := filepath.Join(outputDir, "static", subpath)

		parentDir := filepath.Dir(outputFile)
		if err := os.MkdirAll(parentDir, fs.ModePerm); err != nil {
			log.Fatal(err)
		}

		err := copyAsset(asset, outputFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}
