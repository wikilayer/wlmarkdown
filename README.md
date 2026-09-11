# wlmarkdown

The WikiLayer markdown dialect: ordinary CommonMark with tables, strikethrough and
task lists, plus the constructs the dialect adds of its own.

Today that is callouts. A blockquote whose first line is exactly a marker becomes a
callout of that class:

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

## What this library does and does not do

It **recognises**. `Dialect()` returns a goldmark that parses the dialect, and the
callout node carries the one thing the source says: its class.

It does not render, translate or resolve. A title for the callout, an icon, a
colour, a link target looked up in a store — all of that belongs to whoever holds
the page, because each of them answers differently on a web page and in an app.

## Use

```go
source := []byte("> [!TIP]\n> Try the shorter form.\n")

var out bytes.Buffer
if err := wlmarkdown.Dialect().Convert(source, &out); err != nil {
    return err
}
```

`Recognise(source)` returns the flat list of dialect constructs found, in document
order. It is what the corpus is written against, so every port of this library
answers the same questions with the same words.

## The corpus

`corpus/dialect.yaml` holds the cases the dialect is defined by: a piece of
markdown and what must be recognised in it. The Go test reads that file directly,
and so does every port, which is what keeps them from drifting apart. A new
construct is added to the corpus once and is then asked of all of them.
