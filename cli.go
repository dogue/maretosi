package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func parseCli() (*cli.Command, error) {
	cwd, _ := os.Getwd()

	app := &cli.Command{
		Name:      "maretosi",
		Usage:     "render some markdown files to static html",
		Version:   VERSION,
		Copyright: "Copyright ©️ 2024 dogue <https://github.com/dogue>\nBSD-3 License",
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
			&cli.StringFlag{
				Name:        "templates",
				Aliases:     []string{"t"},
				Usage:       "html templates directory",
				Value:       "templates",
				Destination: &templDir,
			},
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
				Usage:       "skip processing static assets",
				Value:       false,
				Destination: &skipAssets,
			},
		},
		// prevent printing help when no flags passed
		Action: func(ctx context.Context, c *cli.Command) (err error) {
			return
		},
	}

	if app.IsSet("verbose") {
		verbose = true
	}

	return app, app.Run(context.Background(), os.Args)
}
