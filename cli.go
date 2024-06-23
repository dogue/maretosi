package main

import (
	"context"
	// "fmt"
	"log"

	// "fmt"

	// "fmt"
	// "errors"
	// "fmt"
	// "log"
	"os"
	"path/filepath"

	// "strings"

	"github.com/urfave/cli/v3"
)

func run() {
	cwd, _ := os.Getwd()
	cwd_base := filepath.Base(cwd)

	app := &cli.Command{
		Name:      "maretosi",
		Usage:     "render some markdown files to static html",
		Version:   VERSION,
		Copyright: "©️ 2024 dogue <https://github.com/dogue>",
		UsageText: "maretosi [options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				Usage:       "markdown source directory",
				Value:       filepath.Join(cwd, "content"),
				Destination: &inputDir,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "html destination directory",
				Value:       "public",
				Destination: &outputDir,
			},
			&cli.StringFlag{
				Name:        "title",
				Aliases:     []string{"t"},
				Usage:       "site title",
				Value:       cwd_base,
				Destination: &siteTitle,
			},
			&cli.StringFlag{
				Name:        "assets",
				Aliases:     []string{"a"},
				Usage:       "static assets source directory",
				Value:       filepath.Join(cwd, "assets"),
				Destination: &assetsDir,
			},
			&cli.BoolFlag{
				Name:        "no-ext",
				Usage:       "disable markdown extensions",
				Value:       false,
				Destination: &disableExts,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) (err error) {
			renderAll()
			copyAssets()
			return
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
