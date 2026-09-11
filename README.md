# wlmarkdown

The WikiLayer markdown dialect: GitHub-flavoured markdown, meaning CommonMark plus
tables, strikethrough, task lists and bare-URL linking, and then the constructs the
dialect adds of its own.

Today those are callouts and map embeds. A blockquote whose first line is exactly a
marker becomes a callout of that class:

```markdown
> [!WARNING]
> This cannot be undone.
```

| marker | class |
|---|---|
| `[!NOTE]` | `note` |
| `[!TIP]` | `tip` |
| `[!IMPORTANT]` | `important` |
| `[!WARNING]` | `warning` |
| `[!CAUTION]` | `caution` |

Anything else — a marker sharing its line with words, a marker nobody declared, a
quote without one — stays an ordinary quote.

A `[!MAP]` marker followed by a line of two numbers becomes a map embed. Whatever
follows the coordinates is its caption:

```markdown
> [!MAP]
> 44.7866, 20.4489
> Belgrade, the city centre
```

Both numbers have to parse as numbers, or the quote stays a quote. They are handed
on as the source wrote them, digit for digit, because rounding a coordinate moves
the point.

A link may name a node instead of a URL, under the scheme `page:` or `block:`:

```markdown
Read [the libraries page](page:42685) first.
```

The dialect reads the scheme and hands on whatever follows it, character for
character. A tail is as often a name as a number, `page:home` beside `page:42685`,
and which names exist is a question a store answers, along with whether anything is
there at all and what URL it turns into. None of those is a parser's to answer.

## Limitations

A quote within a quote is not looked into. The inner one stays an ordinary quote, so
a map inside a note, and a note inside a note, are quotes as well. Inline
constructs are a different matter: a link is found wherever it sits, a callout
included. Both are corpus cases, so neither can change without a case going red.

An autolink is not reported either. `<https://example.com/page>` stays whatever
CommonMark makes of it, and only a link written with brackets and a destination
comes back from `Recognise`.

Markdown as a whole is split on this, along the line of what the syntax is for.
GitHub, whose alert spelling this dialect borrows, states that "alerts cannot be
nested within other elements". Obsidian, which builds a knowledge base out of them,
says "you can nest callouts in multiple levels" and gives an example three deep.
Material for MkDocs nests admonitions by indentation and describes them as allowing
"the inclusion and nesting of arbitrary content".

This dialect follows GitHub, because the choice is not a symmetric one. Refusing
nesting stays reversible for as long as no page relies on it, while allowing it
commits every port to the same recursion and cannot be taken back from pages
already written.

## What this library does and does not do

It **recognises**. `New().Extensions()` hands over the goldmark extenders that parse
the dialect, and each node carries the one thing the source says: a callout its
class, a map its point.

It does not render, translate or resolve. A title for the callout, an icon, a
colour, a link target looked up in a store — all of that belongs to whoever holds
the page, because each of them answers differently on a web page and in an app.

Which markers exist is not among the things a caller sets. That table is what makes
the dialect this one rather than another, so `New()` is the only dialect there is.

## Use

`New().Recognise(source)` returns the flat list of dialect constructs found, in
document order:

```go
found := wlmarkdown.New().Recognise([]byte("> [!TIP]\n> Try the shorter form.\n"))
```

That list is what the corpus is written against, so every port of this library
answers the same questions with the same words.

To render, compose a goldmark of your own from the extenders and add a renderer for
each of the dialect's node kinds, `KindCallout` and `KindMapEmbed`:

```go
md := goldmark.New(
    goldmark.WithExtensions(wlmarkdown.New().Extensions()...),
    goldmark.WithRendererOptions(renderer.WithNodeRenderers(
        util.Prioritized(yourCalloutRenderer{}, 500),
        util.Prioritized(yourMapRenderer{}, 500),
    )),
)
```

Both renderers are yours to write, and a goldmark without them cannot render a
document that carries either node.

## Running it

```sh
make test    # the corpus, plus the wiring test
make lint    # go vet, gofmt, staticcheck, commentcensor
```

## The corpus

`corpus/dialect.yaml` holds the cases the dialect is defined by: a piece of
markdown and what must be recognised in it. The Go test reads that file directly,
and so does every port, which is what keeps them from drifting apart. A new
construct is added to the corpus once and is then asked of all of them.
