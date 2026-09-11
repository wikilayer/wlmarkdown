package wlmarkdown

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
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

func (d Dialect) Markdown() goldmark.Markdown {
	return goldmark.New(goldmark.WithExtensions(
		extension.Table,
		extension.Strikethrough,
		extension.TaskList,
		&mapEmbedExtension{marker: d.mapMarker},
		&calloutExtension{classByMarker: d.calloutClassByMarker},
	))
}
