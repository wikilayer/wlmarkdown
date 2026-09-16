package wlmarkdown

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Strip keeps reader-visible words, including code and map captions, while
// omitting unreadable map embeds and collapsing markdown structure to spaces.
func Strip(source []byte) string {
	dialect := New()
	doc := dialect.Parser().Parse(text.NewReader(source))

	var written strings.Builder
	walk(doc, func(n ast.Node) ast.WalkStatus {
		switch node := n.(type) {
		case *MapEmbed:
			written.WriteString(" ")
			written.WriteString(Strip([]byte(node.Caption)))
			written.WriteString(" ")
			return ast.WalkSkipChildren
		case *Unreadable:
			return ast.WalkSkipChildren
		case *ast.FencedCodeBlock:
			writeLines(&written, node.Lines(), source)
			return ast.WalkSkipChildren
		case *ast.CodeBlock:
			writeLines(&written, node.Lines(), source)
			return ast.WalkSkipChildren
		case *ast.Text:
			written.Write(node.Segment.Value(source))
			if node.SoftLineBreak() || node.HardLineBreak() {
				written.WriteString(" ")
			}
		}
		if n.Type() == ast.TypeBlock {
			written.WriteString(" ")
		}
		return ast.WalkContinue
	})
	return dialect.squeezed(written.String())
}

func writeLines(into *strings.Builder, lines *text.Segments, source []byte) {
	into.WriteString(" ")
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		into.Write(line.Value(source))
	}
	into.WriteString(" ")
}
