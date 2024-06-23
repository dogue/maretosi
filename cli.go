package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func parseCli() (*cli.Command, error) {
	cwd, _ := os.Getwd()
	// cwdBase := filepath.Base(cwd)

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
				DefaultText: "content",
				Destination: &contentDir,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "html destination directory",
				Value:       "public",
				DefaultText: "public",
				Destination: &outputDir,
			},
			// &cli.StringFlag{
			// 	Name:        "title",
			// 	Aliases:     []string{"t"},
			// 	Usage:       "site title",
			// 	Value:       cwd_base,
			// 	Destination: &siteTitle,
			// },
			&cli.StringFlag{
				Name:        "assets",
				Aliases:     []string{"a"},
				Usage:       "static assets source directory",
				Value:       filepath.Join(cwd, "assets"),
				DefaultText: "assets",
				Destination: &assetsDir,
			},
			&cli.BoolFlag{
				Name:        "no-assets",
				Usage:       "skip processing statis assets",
				Value:       false,
				Destination: &skipAssets,
			},
			&cli.BoolFlag{
				Name:        "no-ext",
				Usage:       "disable markdown extensions",
				Value:       false,
				Destination: &disableExts,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) (err error) {
			return
		},
	}

	return app, app.Run(context.Background(), os.Args)
}
