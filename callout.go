package wlmarkdown

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var KindCallout = ast.NewNodeKind("wikilayer.Callout")

type Callout struct {
	ast.BaseBlock
	Class string
}

func (n *Callout) Kind() ast.NodeKind { return KindCallout }

func (n *Callout) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Class": n.Class}, nil)
}

type calloutExtension struct {
	classByMarker map[string]string
}

func (e *calloutExtension) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(parser.WithASTTransformers(
		util.Prioritized(&calloutTransformer{classByMarker: e.classByMarker}, 100),
	))
}

type calloutTransformer struct {
	classByMarker map[string]string
}

type calloutTarget struct {
	quote        *ast.Blockquote
	paragraph    *ast.Paragraph
	class        string
	firstLineEnd int
}

func (t *calloutTransformer) Transform(doc *ast.Document, reader text.Reader, _ parser.Context) {
	source := reader.Source()
	targets := calloutTargets(doc, source, t.classByMarker)
	for _, found := range targets {
		callout := &Callout{Class: found.class}
		stripMarker(found.paragraph, found.firstLineEnd)
		if found.paragraph.FirstChild() == nil {
			found.quote.RemoveChild(found.quote, found.paragraph)
		}
		for child := found.quote.FirstChild(); child != nil; child = found.quote.FirstChild() {
			found.quote.RemoveChild(found.quote, child)
			callout.AppendChild(callout, child)
		}
		if parent := found.quote.Parent(); parent != nil {
			parent.ReplaceChild(parent, found.quote, callout)
		}
	}
}

func calloutTargets(doc *ast.Document, source []byte, classByMarker map[string]string) []calloutTarget {
	var targets []calloutTarget
	walk(doc, func(n ast.Node) ast.WalkStatus {
		quote, ok := n.(*ast.Blockquote)
		if !ok {
			return ast.WalkContinue
		}
		paragraph, ok := quote.FirstChild().(*ast.Paragraph)
		if !ok || paragraph.Lines().Len() == 0 {
			return ast.WalkSkipChildren
		}
		firstLine := paragraph.Lines().At(0)
		marker := strings.TrimRight(string(firstLine.Value(source)), "\n")
		class, ok := classByMarker[marker]
		if !ok {
			return ast.WalkSkipChildren
		}
		targets = append(targets, calloutTarget{
			quote:        quote,
			paragraph:    paragraph,
			class:        class,
			firstLineEnd: firstLine.Stop,
		})
		return ast.WalkSkipChildren
	})
	return targets
}

func stripMarker(paragraph *ast.Paragraph, lineEnd int) {
	for child := paragraph.FirstChild(); child != nil; {
		next := child.NextSibling()
		written, ok := child.(*ast.Text)
		if !ok || written.Segment.Stop > lineEnd {
			return
		}
		paragraph.RemoveChild(paragraph, child)
		child = next
	}
}
