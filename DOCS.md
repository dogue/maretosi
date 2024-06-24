# Maretosi User Manual

## Overview

Maretosi is a static site generator with few features. It is a narrowly scoped tool by design, with the goal being to provide enough flexibility for general use without adding niche or esoteric features.

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
   0.0.3
                                                                                                                                                                            
COMMANDS:
   help, h  Shows a list of commands or help for one command
                                                                                                                                                                            
GLOBAL OPTIONS:
   --input value, -i value   markdown source directory (default: "<current directory>/content")
   --output value, -o value  html destination directory (default: "public")
   --title value, -t value   site title (default: "maretosi")
   --assets value, -a value  static assets source directory (default: "<current directory>/assets")
   --no-assets               skip processing static assets (default: false)
   --no-ext                  disable markdown extensions (default: false)
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
                                                                                                                                                                            
COPYRIGHT:
   © 2024 dogue <https://github.com/dogue>
```

## Command Line Flags

* `--input`, `-i` -- Sets the directory path containing your markdown content. If not provided, defaults to `<current directory>/content`
* `--output`, `-o` -- Sets the directory path for exporting rendered output files. If not provided, defaults to `<current directory>/public`
* `--title`, `-t` -- Sets the HTML `<title>` element value in the rendered files. If not provided, defaults to `<current directory>`. Can be overriden per file by the `title` attribute in the front matter
* `--assets`, `-a` -- Sets the directory path containing your static asset files (JS, CSS, etc). If not provided, defaults to `<current directory>/assets`
* `--no-assets` -- Skips copying of static assets. Useful for barebones sites without CSS, Javascript, or other asset files.
* `--no-ext` -- Disables all markdown extensions. If not set, the `CommonExtensions` set from [blackfriday](https://pkg.go.dev/github.com/russross/blackfriday/v2#pkg-constants) is enabled

## HTML Template

Maretosi ships with a default HTML template that is designed to be general enough for most uses. However, if you wish to supply a custom template see the [Front Matter](#front-matter) section below.

```html
<!DOCTYPE html>
<html>

<head>
    <title>{{ .title }}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="{{ .stylesheet }}">
</head>

<body>
    {{ .body }}
</body>

</html>
```

## Front Matter {#front-matter}

Maretosi supports a simple markdown front matter system for modifying the rendered HTML output per file.

Front matter should be placed at the top of all of your content files, even if it is empty. It uses `~~~` as the delimiter. Not including the delimiters will result in an error during processing.

A typical example might look like this:

```markdown
~~~
title = "Site Title"
stylesheet = "/static/myblog.css"
template = "path/to/template.html"
~~~

# My Blog

Lorem ipsum and so on...
```

Whereas an empty front matter might look like this:

```markdown
~~~
~~~

# My Blog

Lorem ipsum and so forth...
```

### Abusing Front Matter

Internally, the front matter is extracted and parsed as TOML into a `map[string]interface{}`. While it's not the intended use, this means that you can add arbitrary data into the front matter and extract it in a custom template.

*content.md*
```md
~~~
title = "Custom Data"
template = "my_template.html"
foo = "bar"
~~~

Lorem ipsum something something
```

*my_template.html*
```
<!DOCTYPE html>
<html>
<head>
</head>
<body>
{{ .foo }}
</body>
</html>
```

## Contributing

Bug reports and pull requests are welcome.
