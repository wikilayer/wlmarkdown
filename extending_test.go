package wlmarkdown_test

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/wikilayer/wlmarkdown"
)

const afterTheDialect = 200

type kindRecorder struct {
	kinds map[ast.NodeKind]bool
}

func (r *kindRecorder) Transform(doc *ast.Document, _ text.Reader, _ parser.Context) {
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			r.kinds[n.Kind()] = true
		}
		return ast.WalkContinue, nil
	})
}

func TestAHostTransformerAtTwoHundredMeetsTheDialectsNodes(t *testing.T) {
	recorded := &kindRecorder{kinds: map[ast.NodeKind]bool{}}
	md := goldmark.New(
		goldmark.WithExtensions(wlmarkdown.New().Extensions()...),
		goldmark.WithParserOptions(parser.WithASTTransformers(
			util.Prioritized(recorded, afterTheDialect),
		)),
	)

	source := "> [!NOTE]\n> Where it happened.\n>\n> > [!MAP]\n> > 44.7866, 20.4489\n"
	md.Parser().Parse(text.NewReader([]byte(source)))

	if !recorded.kinds[wlmarkdown.KindCallout] {
		t.Errorf("a transformer at priority %d met no callout, so a host cannot dress one there",
			afterTheDialect)
	}
	if !recorded.kinds[wlmarkdown.KindMapEmbed] {
		t.Errorf("a transformer at priority %d met no map, so a host cannot dress one there",
			afterTheDialect)
	}
}
