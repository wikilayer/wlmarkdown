package wlmarkdown

import (
	"slices"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Declined struct {
	Marker string
}

var declinedKey = parser.NewContextKey()

func DeclinedIn(pc parser.Context) []Declined {
	held, ok := pc.Get(declinedKey).([]Declined)
	if !ok {
		return nil
	}
	return held
}

const declinedPriority = mapEmbedPriority + 10

type declinedExtension struct {
	markers []string
}

func (e *declinedExtension) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(parser.WithASTTransformers(
		util.Prioritized(&declinedTransformer{markers: e.markers}, declinedPriority),
	))
}

type declinedTransformer struct {
	markers []string
}

func (t *declinedTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()
	var declined []Declined

	walk(doc, func(n ast.Node) ast.WalkStatus {
		quote, ok := n.(*ast.Blockquote)
		if !ok {
			return ast.WalkContinue
		}
		paragraph, ok := quote.FirstChild().(*ast.Paragraph)
		if !ok || paragraph.Lines().Len() == 0 {
			return ast.WalkContinue
		}
		opening := line(paragraph.Lines().At(0), source)
		if slices.Contains(t.markers, opening) {
			declined = append(declined, Declined{Marker: opening})
		}
		return ast.WalkContinue
	})

	pc.Set(declinedKey, declined)
}
