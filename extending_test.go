package wlmarkdown_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	source := "> [!NOTE]\n> Where it happened.\n>\n> > [!MAP]\n> > 44.7866, 20.4489\n" +
		"\n> [!MAP]\n> 999, 20.4489\n"
	md.Parser().Parse(text.NewReader([]byte(source)))

	require.NotEmpty(t, wlmarkdown.Kinds(), "with no kinds on offer this test cannot fail")
	for _, kind := range wlmarkdown.Kinds() {
		assert.True(t, recorded.kinds[kind],
			"a transformer at priority %d met no %s, so a host cannot dress one there",
			afterTheDialect, kind)
	}
}
