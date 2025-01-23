# Maretosi User Manual

## Overview

Maretosi is a static site generator with few features. It is a narrowly scoped tool by design, with the goal being to provide enough flexibility for general use with a minimal feature set.

If you're looking for a more traditional "kitchen-sink" static site generator, we recommend examining [Hugo](https://gohugo.io/) or [Zola](https://getzola.org).

## Command Line Flags

Maretosi is configured entirely at runtime via command line flags. Running without providing any flags will use the default values. This allows for convention over configuration when structuring your sources.

* `--input`, `-i` -- Sets the directory path containing your markdown content. If not provided, defaults to `<current directory>/content`
* `--output`, `-o` -- Sets the directory path for exporting rendered output files. If not provided, defaults to `<current directory>/public`
* `--templates`, `-t` -- Sets the directory path containing your HTML template files. If not provided, defaults to `<current directory>/templates`
* `--assets`, `-a` -- Sets the directory path containing your static asset source files (JS, CSS, etc). If not provided, defaults to `<current directory>/assets`
* `--no-assets` -- Skips copying of static assets. Useful for barebones sites without CSS, Javascript, or other asset files.

## HTML Templates

Maretosi uses Go's standard library [HTML template](https://pkg.go.dev/html/template) package. All templates are parsed prior to any content being rendered and must include a `{{ define "<template name>" }}` directive. Maretosi includes a basic built-in template intended for quick prototyping, the contents of which can be examined at the bottom of this file.

### Front Matter

Maretosi supports a simple markdown front matter system for modifying the rendered HTML output per file. Front matter attributes are defined in either TOML or YAML. TOML front matter is delimited by three or more `+` characters. YAML front matter uses `-` characters. The beginning and ending delimiters must match each other.

A typical TOML example might look like this:

```markdown
+++
title = "Page Title"
__template = "blog"
+++

# My Blog

Lorem ipsum and so on...
```

And with YAML:

```markdown
---
title: "Page Title"
__template: "blog"
---

# My Blog

Lorem ipsum and so on...
```

### Template Data

All fields defined in the front matter are available inside templates. Any data that can be represented in TOML or YAML can be passed into your templates.

Maretosi uses a few built-in front matter fields internally to modify the behavior of the the markdown renderer. These built-in fields are prefixed with `__` to avoid collisions with user fields.

* `__template` - (string) Specifies a path to an HTML template for this file. If not provided, the built-in template is used.
* `__syntax_theme` - (string) Specifies a [Chroma style](https://github.com/alecthomas/chroma/tree/master/styles) to use for highlighting code blocks. Defaults to "github-dark".
* `__line_numbers` - (bool) Specifies whether or not line numbers should be included in code blocks. False by default.
* `__body` - (string) Contains the HTML content rendered from the markdown input file. This field is not defined in the front matter, but is instead used to place the rendered content inside a template.

## Contributing

Bug reports and pull requests are welcome. Please report issues [here](https://github.com/dogue/maretosi/issues).

## Built-In HTML Template

```html
{{ define "__built_in" }}
<!DOCTYPE html>
<html>

<head>
    <title>{{ .title }}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>

<body>
    {{ .__body }}
</body>

</html>
{{ end }}
```
