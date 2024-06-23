package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func main() {
	app, err := parseCli()
	if err != nil {
		log.Fatal(err)
	}

	if app.IsSet("help") {
		os.Exit(0)
	}

	if err = validateDirs(); err != nil {
		log.Fatal(err)
	}

	var contentFiles []string
	contentWalker := walker(&contentFiles)
	if err := filepath.WalkDir(contentDir, contentWalker); err != nil {
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
	assetsWalker := walker(&assetsFiles)
	if err := filepath.WalkDir(assetsDir, assetsWalker); err != nil {
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
