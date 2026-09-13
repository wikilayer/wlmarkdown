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

An entry names the call or the field it is about and no more; the
[package reference](https://pkg.go.dev/github.com/wikilayer/wlmarkdown) carries the
signatures, and the README an example of each, including what the two files under
`corpus/` hold. Which versions of goldmark and of Go it is built against is in
`go.mod`, where it cannot go stale. This file says only what changed between
versions and what that asks of you.

The version is 0.x because the shape is still settling: every application built on
this library so far has asked for something in the API rather than working around
its absence, so a minor may still change an answer you relied on. Read the entry
before taking one.

[The Swift port](https://github.com/wikilayer/wlmarkdown-swift) answers the same
corpus, and the two agree on all of it. One difference lies outside what the corpus
can ask: a bare URL is a link here and plain words there, because swift-markdown
offers no way to switch linkifying on. That port's changelog carries the same note.

Changes are documented here in the format of
[Keep a Changelog](https://keepachangelog.com/).

## 0.6.0 - 2026-09-13

### Added

- `Kinds()` hands back the `ast.NodeKind` of every node this dialect builds: today
  `KindCallout`, `KindMapEmbed` and `KindUnreadable`. A host composing its own
  goldmark registers a renderer per kind, and a kind it misses is not a missing box
  but a panic on the first page carrying that construct. 0.5.0 asked you to register
  one for `KindUnreadable` and did not say that much: if you took that release and
  did not, the first unreadable point a reader writes takes the page down. Hold your
  registrations against this list and a kind added in a later minor is a red test
  rather than a broken page. A host that instead replaces these nodes with its own, from a
  transformer at priority 200 or above as 0.1.1 describes, reads the same list to
  know what it can meet there. A later minor may add a kind; none already in the
  list will change what it means.

### Changed

- `corpus/dialect.yaml` gained a `declined` key per case, and new cases along with it.
  If you run the corpus yourself, a decoder that refuses unknown keys has to learn
  this one.
- What `DeclinedIn` reports, said properly, because 0.2.0 named two ways to arrive
  there and its own rule has always allowed three: a map whose second line does not
  read as a point, a callout written inside another callout's quote, and a marker of
  any kind inside an ordinary quote, which stops the dialect before it reads what is
  under it. A point outside 90 or 180 is not among them — since 0.5.0 that is an
  `Unreadable` block, something made rather than turned down, and this call has been
  silent about it since. Nothing changed here in 0.6.0: the call answered this way
  before the entry described it correctly, and the corpus now asks the same of both
  ports case by case.
- The minor moves with [the Swift port](https://github.com/wikilayer/wlmarkdown-swift),
  which fixes three answers of its own for this release. Between them these fixes
  close every difference the corpus can reach, and both ports now answer all of it
  alike. Matching major and minor said that much from 0.4.0 on and were wrong to;
  what makes it true now is that the corpus asks about what was turned down as well.

### Fixed

- **Breaking for anyone printing it:** the words of an `Unreadable` block keep the
  markdown they were written in.

  ```
  > [!MAP]
  > 999, 20
  > [the street](page:1)
  ```

  `Found.Text` on the `unreadable` entry held `[!MAP] 999, 20 the street`, the link
  flattened to its text; it now holds `[!MAP] 999, 20 [the street](page:1)`. It is
  one line either way: the quote's lines arrive joined by a single space with every
  run of blanks squeezed to one, as a callout's words always have. The marker and the
  coordinates are part of them on purpose, because the block exists to show the author
  what they typed.

  Draw that text through a markdown renderer, or escape it before it reaches the page:
  printed raw it shows the reader brackets and parentheses, and printed as HTML it
  hands whatever the author typed to the browser. A caption has always come back this
  way, and the two now match. In the tree the same words are the node's own children
  and always were; nothing there changed.

  The Swift port reads the source directly and never lost markdown this way. What it
  lost instead was everything below the first paragraph of the quote, which its own
  changelog covers.
- A callout's `Found.Text` no longer carries the words of an unreadable map written
  inside it. For

  ```
  > [!NOTE]
  > Where to find us.
  >
  > > [!MAP]
  > > 999, 20
  ```

  the callout said `Where to find us. [!MAP] 999, 20` and now says `Where to find
  us.`; those words arrive in the `unreadable` entry that follows it, where a working
  map's have always arrived. A search index or a preview built from callout text gets
  shorter here. The tree is unchanged: the block is still a child of the callout and
  renders where it stands.

## 0.5.0 - 2026-09-13

### Changed

- A point nowhere on Earth is no longer a place. `corpus/rules.yaml` now names how
  far a coordinate may go, `90` and `180`, beside the alphabets it already spells
  out, and both are compared digit by digit rather than through a float, so a number
  too long for one is judged by the same rule. A quote written as a map whose point
  is outside that comes back as `kind: unreadable` carrying the words as written,
  and in the tree it is an `Unreadable` node rather than a `MapEmbed`.

  Register a renderer for `KindUnreadable` alongside the two you already have. What
  it looks like is yours, as a callout's colour is; that it is visible is the point.
  Before this, `> [!MAP]\n> 999, 999` drew a map of a place the page does not name,
  and the reader who wrote the coordinates had nothing to tell them so.

  The bound is the last place there is, not the first one missing: `-90, 180` is a
  point at the pole and stays a map. Four cases in `corpus/dialect.yaml` hold that
  line — a latitude past the pole, a longitude past the meridian, the pole itself,
  and a run of digits no float could hold — so a port that answers any of them
  differently goes red rather than surprising a reader. Reading digits rather than
  parsing a number is what 0.3.0 pinned, and that has not changed: the run of digits
  too long for any float that it named is inside the bound and is still a place. How
  long a coordinate is was never the question, and a float's opinion of it is not one
  either.

## 0.4.0 - 2026-09-13

### Fixed

- `DeclinedIn` no longer answers a second document with the first one's refusals.
  The refusals were written into the parse context only when there were any, so a
  context parsed with twice kept what the earlier document turned down. Build a
  fresh context per parse if you like; you no longer have to.

### Changed

- The version moves with the Swift port, which is fixing a difference of its own at
  the same time. Matching major and minor say the two ports answer the corpus alike,
  which is what the corpus is run on both sides for, and they say it only while they
  move together. They do not promise more than the corpus reaches: the Swift port's
  changelog names the differences it cannot.

## 0.3.0 - 2026-09-13

### Added

- `DeclinedIn(context)`: the quotes this dialect turned down, one entry per quote
  with the marker it carried. 0.2.0 named the two ways to arrive there and left the
  finding to you, through `Markers()` and a walk of your own; this replaces that
  walk. Doing it from outside meant writing the dialect's rule for what opens a
  construct a second time, in your code, and two copies of a rule drift. Pass the
  `parser.Context` you parsed with, which is the `Parser()` route rather than
  `Recognise`; the README has it in full.
- Four cases in the corpus pinning down where a coordinate stops being one: exponent
  notation, a number ending on its point and one opening on it are refused, and a
  run of digits too long for any float is still a coordinate, because the dialect
  reads digits rather than parsing a number. Every port answers these now; 0.1.0
  and 0.2.0 already behaved this way and nothing said it out loud.

## 0.2.0 - 2026-09-13

### Added

- `Markers()`, every marker this dialect opens a construct with — the five callout
  markers and `[!MAP]` — sorted, and spelled the way a document spells them:
  `[!NOTE]`, not the class `note` that `Classes()` lists. Walk the parsed document
  and a blockquote that still stands as a blockquote, whose first line is exactly
  one of these — exactly, because a marker sharing its line with words was never a
  candidate — is a construct the dialect declined to make. There are two ways to
  arrive there: a map whose second line is missing or does not read as a pair of
  coordinates, and a callout written as a quote inside another callout's quote,
  whose inner blockquote this dialect leaves standing. Reading that first line off
  the node is yours to do — the list spares you writing the six markers out a
  second time, not the walk. As with `Classes()` and
  `Schemes()`, a later version may add a marker; none already in the list will
  change what it opens.

### Changed

- The two node kinds now print as `wlmarkdown.Callout` and `wlmarkdown.MapEmbed`
  instead of `wikilayer.Callout` and `wikilayer.MapEmbed`. **If you compare
  `Kind().String()` against either of the old strings, or keep a recorded AST dump
  in your tests, that text has to change.** Nothing else does: the kinds still
  compare as they did, so code that switches on `KindCallout` or `KindMapEmbed` is
  untouched.

  The old names belonged to the application holding the pages rather than to this
  library, and such an application, replacing these nodes with its own, was left
  registering a kind under the name it wanted for them. Goldmark allows two kinds
  spelled alike and neither misbehaves, so this cost no correctness — only that a
  dump stopped saying which of the two you were looking at.

## 0.1.1 - 2026-09-13

### Added

- A bound on where this library's own work sits in goldmark's order. Goldmark runs
  AST transformers from the lowest priority up, and every transformer this library
  registers runs below 200. So if you swap the node the dialect built — a `Callout`
  or a `MapEmbed`, the two carrying `KindCallout` and `KindMapEmbed` — for an AST
  node of your own, to give a callout a title in your reader's language or a map an
  embed URL, register your transformer at 200 or above and it will meet that node
  already built. Nothing parses differently than in 0.1.0; the bound held
  there too, but it was a number you had to read out of the source, and now it is a
  promise with a test behind it.

### Fixed

- A sign written outside ASCII in `corpus/rules.yaml` — a typographic minus, say —
  was read a byte at a time and so went unrecognised here, while a port of this
  dialect to another language, walking characters rather than bytes, would have
  honoured it. That file spells out which characters may open a coordinate, and the
  two it names today, `+` and `-`, behave exactly as before; what changed is that
  the list is read character by character the way the file means it, so two ports
  cannot part company over a sign someone adds to it later.

## 0.1.0 - 2026-09-12

First release.

### Added

- `New()` returns the dialect, and it is the only one there is. `Extensions()`
  hands over the goldmark extenders that parse it, `Parser()` a parser built from
  them, and `Recognise(source)` the flat list of what was found, in document order.
- Callouts: a blockquote whose first line is exactly `[!NOTE]`, `[!TIP]`,
  `[!IMPORTANT]`, `[!WARNING]` or `[!CAUTION]`. A marker sharing its line with
  words, or written in lower case, leaves an ordinary quote.
- Map embeds: `[!MAP]` alone on the first line, as a callout marker must be, then a
  line holding latitude and longitude separated by a comma, and whatever follows as
  the caption. A coordinate is an optional sign,
  digits, and optionally a decimal point and more digits. It comes back as a string,
  digit for digit, because a coordinate rounded is a pin in the wrong street.
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

What it will not do and what nesting yields are standing properties rather than
changes, and the README carries them.
