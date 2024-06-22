# Maretosi User Manual

## Overview

Maretosi is a static site generator with few features and no bells or whistles. It is a narrowly scoped tool by design.

If you're looking for a fully featured static site generator, we recommend examining [Hugo](https://gohugo.io/) or [Zola](https://getzola.org).

If instead you're looking for a simple and reliable tool that does one thing well and you don't require all of the extra complexity of those generators, you may find Maretosi suitable to your needs.

## Quick Start

Maretosi is configured entirely at runtime via command line flags. Below is the built-in help text.

```text
NAME:
   maretosi - render some markdown files to static html
                                                                                                                                                                            
USAGE:
   maretosi [options]
                                                                                                                                                                            
VERSION:
   0.0.1
                                                                                                                                                                            
COMMANDS:
   help, h  Shows a list of commands or help for one command
                                                                                                                                                                            
GLOBAL OPTIONS:
   --input value, -i value   markdown source directory (default: "<current directory>/content")
   --output value, -o value  html destination directory (default: "public")
   --title value, -t value   site title (default: "maretosi")
   --assets value, -a value  static assets source directory (default: "<current directory>/assets")
   --no-ext                  disable markdown extensions (default: false)
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
                                                                                                                                                                            
COPYRIGHT:
   © 2024 dogue <https://github.com/dogue>
```

## Command Line Flags

* `--input`, `-i` -- Sets the directory path containing your markdown content. If not provided, defaults to `<current directory>/content`
* `--output`, `-o` -- Sets the directory path for exporting rendered output files. If not provided, defaults to `<current directory>/public`
* `--title`, `-t` -- Sets the HTML `<title>` element value in the rendered files. If not provided, defaults to `<current directory>`
* `--assets`, `-a` -- Sets the directory path containing your static asset files (JS, CSS, etc). If not provided, defaults to `<current directory>/assets`
* `--no-ext` -- Disables all markdown extensions. If not set, the `CommonExtensions` set from [blackfriday](https://pkg.go.dev/github.com/russross/blackfriday/v2#pkg-constants) is enabled

## HTML Template

Maretosi does very little templating. Below is the single template file used to generate all HTML pages. At this time there are no plans to make the HTML templates configurable, though the option is open if interest is shown.

```html
<!DOCTYPE html>
<html>

<head>
    <title>{{ .Title }}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="/static/style.css">
</head>

<body>
    {{ .Body }}
</body>

</html>
```

## Contributing

Bug reports and pull requests are welcome.
