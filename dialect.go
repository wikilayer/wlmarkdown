package wlmarkdown

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type Dialect struct {
	calloutClassByMarker map[string]string
	mapMarker            string
	refSchemes           []string
}

func New() Dialect {
	return Dialect{
		calloutClassByMarker: map[string]string{
			"[!NOTE]":      "note",
			"[!TIP]":       "tip",
			"[!IMPORTANT]": "important",
			"[!WARNING]":   "warning",
			"[!CAUTION]":   "caution",
		},
		mapMarker:  "[!MAP]",
		refSchemes: []string{"page", "block"},
	}
}

func (d Dialect) Extensions() []goldmark.Extender {
	return []goldmark.Extender{
		extension.GFM,
		&mapEmbedExtension{marker: d.mapMarker},
		&calloutExtension{classByMarker: d.calloutClassByMarker},
	}
}

func (d Dialect) Parser() parser.Parser {
	return goldmark.New(goldmark.WithExtensions(d.Extensions()...)).Parser()
}
