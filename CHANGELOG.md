# Changelog

wlmarkdown recognises the WikiLayer markdown dialect: GitHub-flavoured markdown
plus callouts, map embeds, and links written against a node id. It recognises and
does nothing else, because a title, an icon, a colour and a resolved address are
answered differently on a web page and in an app, and belong to whoever holds the
pages.

Changes are documented here in the format of
[Keep a Changelog](https://keepachangelog.com/).

## [0.1.0] - 2026-09-12

First release.

### Added

- `New()` returns the dialect, and it is the only one there is. `Extensions()`
  hands over the goldmark extenders that parse it, `Parser()` a parser built from
  them, and `Recognise(source)` the flat list of what was found, in document order.
- Callouts: a blockquote whose first line is exactly `[!NOTE]`, `[!TIP]`,
  `[!IMPORTANT]`, `[!WARNING]` or `[!CAUTION]`. A marker sharing its line with
  words, or written in lower case, leaves an ordinary quote.
- Map embeds: `[!MAP]`, a line of two coordinates, and whatever follows as the
  caption. A coordinate is an optional sign, digits, and optionally a point and
  more digits; the digits are handed on as written, because rounding moves the
  point.
- Links naming a node under `page:` or `block:`. The scheme comes back named and
  the destination exactly as written, because whether anything is there, and what
  URL it becomes, are a store's questions.
- `Classes()` and `Schemes()` name the values that can come back in `Found.Class`
  and `Found.Scheme`, so a caller can hold its own tables against them rather than
  keeping a copy that nothing compares. Both lists grow in a minor version.
- `corpus/rules.yaml` and `corpus/dialect.yaml`: what the dialect knows, and the
  cases it is defined by. Every port reads both, which is what keeps ports from
  drifting apart.

### Worth knowing before you take it

- Nothing here renders. A document carrying a callout or a map needs a renderer of
  your own for `KindCallout` and `KindMapEmbed`; a goldmark without one cannot
  convert it.
- A callout inside a callout is one callout. A map inside a callout is found.
- An autolink is not reported: only a link written with brackets and a destination
  comes back.
- The version is 0.x because the shape is still settling. Nothing consumes this
  library yet, and no port has been written against the corpus; both will have
  something to say about the API.
