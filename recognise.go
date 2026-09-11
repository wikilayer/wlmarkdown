package markdown

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Found struct {
	Kind  string `yaml:"kind"`
	Class string `yaml:"class,omitempty"`
	Text  string `yaml:"text,omitempty"`
}

func Recognise(source []byte) []Found {
	reader := text.NewReader(source)
	doc := Dialect().Parser().Parse(reader)

	found := []Found{}
	walk(doc, func(n ast.Node) ast.WalkStatus {
		callout, ok := n.(*Callout)
		if !ok {
			return ast.WalkContinue
		}
		found = append(found, Found{
			Kind:  "callout",
			Class: callout.Class,
			Text:  plainText(callout, source),
		})
		return ast.WalkSkipChildren
	})
	return found
}

func plainText(from ast.Node, source []byte) string {
	var written strings.Builder
	walk(from, func(n ast.Node) ast.WalkStatus {
		if n.Type() == ast.TypeBlock && written.Len() > 0 {
			written.WriteString(" ")
		}
		if leaf, ok := n.(*ast.Text); ok {
			written.Write(leaf.Segment.Value(source))
			if leaf.SoftLineBreak() || leaf.HardLineBreak() {
				written.WriteString(" ")
			}
		}
		return ast.WalkContinue
	})
	return strings.TrimSpace(written.String())
}
