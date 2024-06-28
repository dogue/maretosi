# Maretosi User Manual

## Overview

Maretosi is a static site generator with few features. It is a narrowly scoped tool by design, with the goal being to provide enough flexibility for general use without adding niche or esoteric features.

If you're looking for a fully featured static site generator, we recommend examining [Hugo](https://gohugo.io/) or [Zola](https://getzola.org).

If you value flexibility and user freedom, you may find Maretosi suitable to your needs.

## Command Line Flags

Maretosi is configured entirely at runtime via command line flags. Running without providing any flags will use the default values. This allows for convention over configuration when structuring your sources.

* `--input`, `-i` -- Sets the directory path containing your markdown content. If not provided, defaults to `<current directory>/content`
* `--output`, `-o` -- Sets the directory path for exporting rendered output files. If not provided, defaults to `<current directory>/public`
* `--assets`, `-a` -- Sets the directory path containing your static asset source files (JS, CSS, etc). If not provided, defaults to `<current directory>/assets`
* `--no-assets` -- Skips copying of static assets. Useful for barebones sites without CSS, Javascript, or other asset files.

## HTML Template

Maretosi ships with a default HTML template that provides a basic layout. However, if you wish to supply a custom template see the [Front Matter](#front-matter) section below.

The built-in template can be examined below. Template identifiers prefixed with double underscores are explained in the front matter section below.

```html
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
```

## Front Matter {#front-matter}

Maretosi supports a simple markdown front matter system for modifying the rendered HTML output per file. Front matter attributes are defined in either TOML or YAML.

TOML front matter is delimited by three or more `+` characters. YAML front matter uses `-` characters. The beginning and ending delimiters must match each other.

A typical TOML example might look like this:

```markdown
+++
title = "Page Title"
__template = "path/to/template.html"
+++

# My Blog

Lorem ipsum and so on...
```

And with YAML:

```markdown
---
title: "Page Title"
__template: "path/to/template.html"
---

# My Blog

Lorem ipsum and so on...
```

### Built-In Fields

Maretosi uses a few built-in front matter fields internally to modify the behavior of the the markdown renderer. These built-in fields are prefixed with `__` to avoid collisions with user fields.
* `__template` - (string) Specifies a path to an HTML template for this file. If not provided, the built-in template is used.
* `__syntax_theme` - (string) Specifies a [Chroma style](https://github.com/alecthomas/chroma/tree/master/styles) to use for highlighting code blocks. Defaults to "github-dark".
* `__line_numbers` - (bool) Specifies whether or not line numbers should be included in code blocks. False by default.
* `__body` - (string) Contains the HTML content rendered from the markdown input file. This field is not defined in the front matter, but is instead used to place the rendered content inside a template.


## Contributing

Bug reports and pull requests are welcome.
