package wlmarkdown

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Found struct {
	Kind    string `yaml:"kind"`
	Class   string `yaml:"class,omitempty"`
	Lat     string `yaml:"lat,omitempty"`
	Lng     string `yaml:"lng,omitempty"`
	Caption string `yaml:"caption,omitempty"`
	Text    string `yaml:"text,omitempty"`
}

func (d Dialect) Recognise(source []byte) []Found {
	reader := text.NewReader(source)
	doc := d.Markdown().Parser().Parse(reader)

	found := []Found{}
	walk(doc, func(n ast.Node) ast.WalkStatus {
		switch spoken := n.(type) {
		case *Callout:
			found = append(found, Found{
				Kind:  "callout",
				Class: spoken.Class,
				Text:  plainText(spoken, source),
			})
			return ast.WalkSkipChildren
		case *MapEmbed:
			found = append(found, Found{
				Kind:    "map",
				Lat:     spoken.Lat,
				Lng:     spoken.Lng,
				Caption: spoken.Caption,
			})
			return ast.WalkSkipChildren
		}
		return ast.WalkContinue
	})
	return found
}

func plainText(from ast.Node, source []byte) string {
	var written strings.Builder
	walk(from, func(n ast.Node) ast.WalkStatus {
		if n.Type() == ast.TypeBlock && needsGap(written.String()) {
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

func needsGap(written string) bool {
	return written != "" && !strings.HasSuffix(written, " ")
}
