package wlmarkdown

import "github.com/yuin/goldmark/ast"

var KindUnreadable = ast.NewNodeKind("wlmarkdown.Unreadable")

type Unreadable struct {
	ast.BaseBlock
}

func (n *Unreadable) Kind() ast.NodeKind { return KindUnreadable }

func (n *Unreadable) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
