package wlmarkdown

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func Dialect() goldmark.Markdown {
	return goldmark.New(goldmark.WithExtensions(
		extension.Table,
		extension.Strikethrough,
		extension.TaskList,
		&mapEmbedExtension{},
		&calloutExtension{},
	))
}
