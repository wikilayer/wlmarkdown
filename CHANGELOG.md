# Changelog

A Go package, `github.com/wikilayer/wlmarkdown`, built on
[goldmark](https://github.com/yuin/goldmark). It recognises the markdown dialect of
WikiLayer, a wiki whose pages are a tree of nodes: GitHub-flavoured markdown plus
callouts, map embeds, and links naming a node instead of a URL.

```markdown
> [!WARNING]
> This cannot be undone.

Start at [the front page](page:home), or at [one paragraph](block:50386) of it.
```

It recognises and does nothing else. What title that callout wears, which icon and
colour it gets, which address `page:home` resolves to: a web page answers each of
those one way and a phone app another, so each belongs to the application holding
the pages rather than to a parser.

Changes are documented here in the format of
[Keep a Changelog](https://keepachangelog.com/).

## 0.1.0 - 2026-09-12

First release.

### Added

- `New()` returns the dialect, and it is the only one there is. `Extensions()`
  hands over the goldmark extenders that parse it, `Parser()` a parser built from
  them, and `Recognise(source)` the flat list of what was found, in document order.
- Callouts: a blockquote whose first line is exactly `[!NOTE]`, `[!TIP]`,
  `[!IMPORTANT]`, `[!WARNING]` or `[!CAUTION]`. A marker sharing its line with
  words, or written in lower case, leaves an ordinary quote.
- Map embeds: `[!MAP]`, then a line holding latitude and longitude separated by a
  comma, and whatever follows as the caption. A coordinate is an optional sign, digits, and optionally a decimal point
  and more digits. It comes back as a string, digit for digit, because a coordinate
  rounded is a pin in the wrong street.
- Links naming a node under `page:` or `block:`. The scheme comes back named and
  the destination exactly as written. Whether `50386` is there at all, and what URL
  it becomes, can only be answered by whatever holds the pages.
- `Classes()` and `Schemes()` list every value that can arrive in `Found.Class` and
  `Found.Scheme`. If you keep a table of your own, an icon per class say, a test
  can check it against these instead of letting the two drift apart unwatched.
  A later version may add a value to either list; none of the values already there
  will change what it means.
- `corpus/rules.yaml` and `corpus/dialect.yaml`: what the dialect knows, and the
  cases that define it. They are meant to be read by an implementation in another
  language as much as by this one, which is what will keep the two answering alike.

### Worth knowing before you take it

- Nothing here renders. A document carrying a callout or a map needs a renderer of
  your own for `KindCallout` and `KindMapEmbed`; a goldmark without one cannot
  convert it.
- A callout inside a callout yields a single entry, the outer callout, and the
  inner marker stays among its words. A map inside a callout yields two entries,
  the callout and the map.
- An autolink is not reported: only a link written with brackets and a destination
  comes back.
- The version is 0.x because the shape is still settling. Nothing uses this library
  yet, and nobody has written the same dialect in another language against the
  corpus; both will have something to say about the API.
