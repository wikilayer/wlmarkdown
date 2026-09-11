# wlmarkdown

The WikiLayer markdown dialect: ordinary CommonMark with tables, strikethrough and
task lists, plus the constructs the dialect adds of its own.

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
> 44.8032201, 20.4726824
> Krunska 72, 11000 Beograd
```

Both numbers have to parse as numbers, or the quote stays a quote. They are handed
on as the source wrote them, digit for digit, because rounding a coordinate moves
the point.

## Limitations

A construct inside a callout is not recognised. A quote within a quote stays an
ordinary quote, so a map inside a note is a quote as well. This is written into the
corpus as a case of its own, rather than left to how the walk happens to be
arranged.

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
