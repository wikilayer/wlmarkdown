package wlmarkdown

import (
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var KindMapEmbed = ast.NewNodeKind("wikilayer.MapEmbed")

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
	marker     string
	coordinate *regexp.Regexp
}

func (e *mapEmbedExtension) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(parser.WithASTTransformers(
		util.Prioritized(
			&mapEmbedTransformer{marker: e.marker, coordinate: e.coordinate},
			mapEmbedPriority,
		),
	))
}

type mapEmbedTransformer struct {
	marker     string
	coordinate *regexp.Regexp
}

type mapEmbedTarget struct {
	quote *ast.Blockquote
	place MapEmbed
}

func (t *mapEmbedTransformer) Transform(doc *ast.Document, reader text.Reader, _ parser.Context) {
	source := reader.Source()
	for _, found := range mapEmbedTargets(doc, source, t.marker, t.coordinate) {
		place := &MapEmbed{
			Lat:     found.place.Lat,
			Lng:     found.place.Lng,
			Caption: found.place.Caption,
		}
		if parent := found.quote.Parent(); parent != nil {
			parent.ReplaceChild(parent, found.quote, place)
		}
	}
}

func mapEmbedTargets(
	doc *ast.Document,
	source []byte,
	marker string,
	coordinate *regexp.Regexp,
) []mapEmbedTarget {
	var targets []mapEmbedTarget
	walk(doc, func(n ast.Node) ast.WalkStatus {
		quote, ok := n.(*ast.Blockquote)
		if !ok {
			return ast.WalkContinue
		}
		place, ok := placeInQuote(quote, source, marker, coordinate)
		if !ok {
			return ast.WalkSkipChildren
		}
		targets = append(targets, mapEmbedTarget{quote: quote, place: place})
		return ast.WalkSkipChildren
	})
	return targets
}

func placeInQuote(
	quote *ast.Blockquote,
	source []byte,
	marker string,
	coordinate *regexp.Regexp,
) (MapEmbed, bool) {
	const markerAndCoordinates = 2

	paragraph, ok := quote.FirstChild().(*ast.Paragraph)
	if !ok {
		return MapEmbed{}, false
	}
	lines := paragraph.Lines()
	if lines.Len() < markerAndCoordinates {
		return MapEmbed{}, false
	}
	if line(lines.At(0), source) != marker {
		return MapEmbed{}, false
	}
	lat, lng, ok := coordinates(line(lines.At(1), source), coordinate)
	if !ok {
		return MapEmbed{}, false
	}
	return MapEmbed{Lat: lat, Lng: lng, Caption: caption(paragraph, source)}, true
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

func coordinates(spoken string, coordinate *regexp.Regexp) (lat, lng string, ok bool) {
	const latitudeAndLongitude = 2

	parts := strings.SplitN(spoken, ",", latitudeAndLongitude)
	if len(parts) != latitudeAndLongitude {
		return "", "", false
	}
	lat = strings.TrimSpace(parts[0])
	lng = strings.TrimSpace(parts[1])
	if !coordinate.MatchString(lat) || !coordinate.MatchString(lng) {
		return "", "", false
	}
	return lat, lng, true
}

func line(segment text.Segment, source []byte) string {
	return strings.TrimRight(string(segment.Value(source)), "\n")
}
