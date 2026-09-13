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

Signatures and fields are not repeated here: the
[package reference](https://pkg.go.dev/github.com/wikilayer/wlmarkdown) carries
those, and the README an example of each, including what the two files under
`corpus/` hold. Which versions of goldmark and of Go it is built against is in
`go.mod`, where it cannot go stale. This file says only what changed between
versions and what that asks of you.

Changes are documented here in the format of
[Keep a Changelog](https://keepachangelog.com/).

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

What it will not do, what nesting yields and why the version is 0.x are standing
properties rather than changes, and the README carries them.
