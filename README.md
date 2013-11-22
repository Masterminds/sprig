# Sprig: Template functions for Go templates

The Go language comes with a [built-in template
language](http://golang.org/pkg/text/template/), but not
very many template functions. This library provides a group of commonly
used template functions.

It is inspired by the template functions found in
[Twig](http://twig.sensiolabs.org/documentation).

## Usage

### Load the Sprig library

To load the Sprig `FuncMap`:

```go

import (
  "github.com/Masterminds/sprig"
  "html/template"
)

// This example illustrates that the FuncMap *must* be set before the
// templates themselves are loaded.
tpl := template.Must(
  template.New("base").FuncMap(sprig.FuncMap()).ParseGlob("*.html")
)


```

### Call the functions inside of templates

