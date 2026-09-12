package wlmarkdown

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Found struct {
	Kind        string `yaml:"kind"`
	Class       string `yaml:"class,omitempty"`
	Lat         string `yaml:"lat,omitempty"`
	Lng         string `yaml:"lng,omitempty"`
	Caption     string `yaml:"caption,omitempty"`
	Scheme      string `yaml:"scheme,omitempty"`
	Destination string `yaml:"destination,omitempty"`
	Text        string `yaml:"text,omitempty"`
}

func (d Dialect) Recognise(source []byte) []Found {
	reader := text.NewReader(source)
	doc := d.Parser().Parse(reader)

	found := []Found{}
	walk(doc, func(n ast.Node) ast.WalkStatus {
		switch spoken := n.(type) {
		case *Callout:
			found = append(found, Found{
				Kind:  "callout",
				Class: spoken.Class,
				Text:  d.plainText(spoken, source),
			})
			return ast.WalkContinue
		case *MapEmbed:
			found = append(found, Found{
				Kind:    "map",
				Lat:     spoken.Lat,
				Lng:     spoken.Lng,
				Caption: spoken.Caption,
			})
			return ast.WalkSkipChildren
		case *ast.Link:
			destination := string(spoken.Destination)
			found = append(found, Found{
				Kind:        "link",
				Scheme:      d.schemeIn(destination),
				Destination: destination,
				Text:        d.plainText(spoken, source),
			})
			return ast.WalkSkipChildren
		}
		return ast.WalkContinue
	})
	return found
}

func (d Dialect) schemeIn(destination string) string {
	for _, scheme := range d.refSchemes {
		if strings.HasPrefix(destination, scheme+":") {
			return scheme
		}
	}
	return ""
}

func (d Dialect) plainText(from ast.Node, source []byte) string {
	var written strings.Builder
	walk(from, func(n ast.Node) ast.WalkStatus {
		if n.Type() == ast.TypeBlock {
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
	return d.squeezed(written.String())
}

func (d Dialect) squeezed(written string) string {
	var said strings.Builder
	spaced := true
	for _, letter := range written {
		if strings.ContainsRune(d.blanks, letter) {
			if !spaced {
				said.WriteString(" ")
			}
			spaced = true
			continue
		}
		said.WriteRune(letter)
		spaced = false
	}
	return strings.TrimSuffix(said.String(), " ")
}
