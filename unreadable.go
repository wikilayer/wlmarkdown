package wlmarkdown

import "github.com/yuin/goldmark/ast"

// KindUnreadable identifies a map whose coordinates lie outside the Earth.
var KindUnreadable = ast.NewNodeKind("wlmarkdown.Unreadable")

// Unreadable preserves the contents of a map that cannot be placed.
type Unreadable struct {
	ast.BaseBlock
}

// Kind returns KindUnreadable.
func (n *Unreadable) Kind() ast.NodeKind { return KindUnreadable }

// Dump writes the node in goldmark's diagnostic tree format.
func (n *Unreadable) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
