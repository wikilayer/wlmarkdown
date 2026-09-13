package wlmarkdown

import (
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var KindMapEmbed = ast.NewNodeKind("wlmarkdown.MapEmbed")

const mapEmbedPriority = calloutPriority + 10

type MapEmbed struct {
	ast.BaseBlock
	Lat     string
	Lng     string
	Caption string
}

func (n *MapEmbed) Kind() ast.NodeKind { return KindMapEmbed }

func (n *MapEmbed) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Lat":     n.Lat,
		"Lng":     n.Lng,
		"Caption": n.Caption,
	}, nil)
}

type mapEmbedExtension struct {
	marker string
	point  coordinate
}

func (e *mapEmbedExtension) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(parser.WithASTTransformers(
		util.Prioritized(
			&mapEmbedTransformer{marker: e.marker, point: e.point},
			mapEmbedPriority,
		),
	))
}

type mapEmbedTransformer struct {
	marker string
	point  coordinate
}

type mapEmbedTarget struct {
	quote *ast.Blockquote
	place MapEmbed
	how   howItReads
}

func (t *mapEmbedTransformer) Transform(doc *ast.Document, reader text.Reader, _ parser.Context) {
	source := reader.Source()
	for _, found := range mapEmbedTargets(doc, source, t.marker, t.point) {
		var made ast.Node = &MapEmbed{
			Lat:     found.place.Lat,
			Lng:     found.place.Lng,
			Caption: found.place.Caption,
		}
		if found.how == aPlaceNowhere {
			written := &Unreadable{}
			for child := found.quote.FirstChild(); child != nil; child = found.quote.FirstChild() {
				found.quote.RemoveChild(found.quote, child)
				written.AppendChild(written, child)
			}
			made = written
		}
		if parent := found.quote.Parent(); parent != nil {
			parent.ReplaceChild(parent, found.quote, made)
		}
	}
}

func mapEmbedTargets(
	doc *ast.Document,
	source []byte,
	marker string,
	point coordinate,
) []mapEmbedTarget {
	var targets []mapEmbedTarget
	walk(doc, func(n ast.Node) ast.WalkStatus {
		quote, ok := n.(*ast.Blockquote)
		if !ok {
			return ast.WalkContinue
		}
		place, how := placeInQuote(quote, source, marker, point)
		if how == noPlaceWritten {
			return ast.WalkSkipChildren
		}
		targets = append(targets, mapEmbedTarget{quote: quote, place: place, how: how})
		return ast.WalkSkipChildren
	})
	return targets
}

func placeInQuote(
	quote *ast.Blockquote,
	source []byte,
	marker string,
	point coordinate,
) (MapEmbed, howItReads) {
	const markerAndCoordinates = 2

	paragraph, ok := quote.FirstChild().(*ast.Paragraph)
	if !ok {
		return MapEmbed{}, noPlaceWritten
	}
	lines := paragraph.Lines()
	if lines.Len() < markerAndCoordinates {
		return MapEmbed{}, noPlaceWritten
	}
	if line(lines.At(0), source) != marker {
		return MapEmbed{}, noPlaceWritten
	}
	lat, lng, how := coordinates(line(lines.At(1), source), point)
	if how == noPlaceWritten {
		return MapEmbed{}, noPlaceWritten
	}
	return MapEmbed{Lat: lat, Lng: lng, Caption: caption(paragraph, source)}, how
}

func caption(paragraph *ast.Paragraph, source []byte) string {
	const afterMarkerAndCoordinates = 2

	lines := paragraph.Lines()
	var said []string
	for i := afterMarkerAndCoordinates; i < lines.Len(); i++ {
		spoken := strings.TrimSpace(line(lines.At(i), source))
		if spoken != "" {
			said = append(said, spoken)
		}
	}
	return strings.Join(said, " ")
}

type howItReads int

const (
	noPlaceWritten howItReads = iota
	aPlace
	aPlaceNowhere
)

func coordinates(spoken string, written coordinate) (lat, lng string, how howItReads) {
	const latitudeAndLongitude = 2

	parts := strings.SplitN(spoken, ",", latitudeAndLongitude)
	if len(parts) != latitudeAndLongitude {
		return "", "", noPlaceWritten
	}
	lat = strings.TrimSpace(parts[0])
	lng = strings.TrimSpace(parts[1])
	if !written.reads(lat) || !written.reads(lng) {
		return "", "", noPlaceWritten
	}
	if !written.within(lat, written.LatitudeWithin) ||
		!written.within(lng, written.LongitudeWithin) {
		return lat, lng, aPlaceNowhere
	}
	return lat, lng, aPlace
}

func (c coordinate) within(spoken, bound string) bool {
	whole, fraction, pointed := strings.Cut(strings.TrimLeft(spoken, c.Signs), c.Point)
	whole = strings.TrimLeft(whole, "0")
	if whole == "" {
		return true
	}
	if len(whole) != len(bound) {
		return len(whole) < len(bound)
	}
	if whole != bound {
		return whole < bound
	}
	return !pointed || strings.Trim(fraction, "0") == ""
}

func (c coordinate) reads(spoken string) bool {
	rest := spoken
	if first, width := utf8.DecodeRuneInString(rest); strings.ContainsRune(c.Signs, first) {
		rest = rest[width:]
	}
	whole, fraction, pointed := strings.Cut(rest, c.Point)
	if !c.digitsOnly(whole) {
		return false
	}
	return !pointed || c.digitsOnly(fraction)
}

func (c coordinate) digitsOnly(spoken string) bool {
	if spoken == "" {
		return false
	}
	return strings.IndexFunc(spoken, func(letter rune) bool {
		return !strings.ContainsRune(c.Digits, letter)
	}) < 0
}

func line(segment text.Segment, source []byte) string {
	return strings.TrimRight(string(segment.Value(source)), "\n")
}
