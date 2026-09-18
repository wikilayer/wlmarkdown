# wlmarkdown

[![Tests](https://github.com/wikilayer/wlmarkdown/actions/workflows/tests.yml/badge.svg)](https://github.com/wikilayer/wlmarkdown/actions/workflows/tests.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/wikilayer/wlmarkdown.svg)](https://pkg.go.dev/github.com/wikilayer/wlmarkdown)

The leading implementation of the WikiLayer markdown dialect. It adds callouts,
map embeds, and `page:` and `block:` links to goldmark's GitHub-flavoured Markdown.
The Swift and Kotlin ports consume the same rules and test corpora.

Install the package:

```sh
go get github.com/wikilayer/wlmarkdown@v0.7.1
```

Recognise structured constructs or extract reader-visible text:

```go
dialect := wlmarkdown.New()
found := dialect.Recognise([]byte("> [!TIP]\n> Try the shorter form.\n"))
plain := wlmarkdown.Strip([]byte("Read **this** before `make test`."))
```

`Recognise` returns a flat list in document order. Each `Found` describes a
callout, map, unreadable map, or link. `Strip` removes markdown syntax for search,
previews, and indexing while keeping code searchable.

## What it recognises

A blockquote whose first line is exactly `[!NOTE]`, `[!TIP]`, `[!IMPORTANT]`,
`[!WARNING]`, or `[!CAUTION]` is a callout. A marker sharing its line with words or
written in another case leaves an ordinary quote.

A `[!MAP]` marker followed by a line of two numbers is a map embed. The remaining
lines of that paragraph are its caption. Coordinates contain digits with an
optional sign and fraction; their spelling is preserved. A latitude may go as far
as 90 and a longitude as far as 180. A pair outside those bounds becomes an
`Unreadable` node carrying the words the author wrote.

A written link may name a node under the `page:` or `block:` scheme. The dialect
reports the destination but does not resolve it against a store. Bare URLs are
linkified by GFM but, like angle-bracket autolinks, are omitted from `Found`.

Read `Markers`, `Classes`, and `Schemes` instead of copying their current values
into an application. `DeclinedIn` reports marked quotes the dialect left unchanged
when the parse uses a `parser.Context`:

```go
pc := parser.NewContext()
p := wlmarkdown.New().Parser()
p.Parse(text.NewReader(source), parser.WithContext(pc))
declined := wlmarkdown.DeclinedIn(pc)
```

## Rendering with goldmark

The library recognises constructs but does not render, decorate, or resolve them.
Build a goldmark instance from its extensions and register a renderer for every
kind returned by `Kinds`:

```go
md := goldmark.New(
    goldmark.WithExtensions(wlmarkdown.New().Extensions()...),
    goldmark.WithRendererOptions(renderer.WithNodeRenderers(
        util.Prioritized(yourCalloutRenderer{}, 500),
        util.Prioritized(yourMapRenderer{}, 500),
        util.Prioritized(yourUnreadableRenderer{}, 500),
    )),
)
```

A goldmark converter without those renderers panics when it reaches a custom node.
An application may instead replace the nodes in its own AST transformer. Every
transformer in this package runs below priority 200, so a transformer registered at
200 or above sees all dialect nodes.

## Nesting

A callout inside another callout remains part of the outer callout. A map directly
inside a callout is still recognised, and links inside callouts are reported. These
priorities are part of the shared corpus rather than renderer policy.

## The corpus

`corpus/rules.yaml` defines the dialect's markers, coordinate alphabet, blanks, and
link schemes. `corpus/dialect.yaml` defines structured recognition, and
`corpus/plain_text.yaml` defines `Strip` and the ports' plain-text functions. New
rules and cases are added here first; every port runs copies of all three.

Bare-URL linking lies outside what the flat corpus can observe. Goldmark and the
Kotlin port enable it; swift-markdown offers no equivalent option.

## Documentation

The public API is published in the [Go Reference](https://pkg.go.dev/github.com/wikilayer/wlmarkdown).

## Development

```sh
make test-build  # compile the package and tests
make test        # run the shared corpus and wiring tests
make lint        # commentcensor, go vet, gofmt, and staticcheck
make build       # all checks and the package build
```

Releases are published by the repository's
[Release workflow](https://github.com/wikilayer/wlmarkdown/actions/workflows/release.yml),
after it repeats the complete build.

## Lines of Code

<picture>
  <source media="(prefers-color-scheme: dark)" srcset=".github/loc-history-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset=".github/loc-history-light.svg">
  <img alt="Lines of code over time" src=".github/loc-history.svg">
</picture>
