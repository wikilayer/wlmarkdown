package wlmarkdown

import "github.com/yuin/goldmark/ast"

func walk(from ast.Node, visit func(ast.Node) ast.WalkStatus) {
	_ = ast.Walk(from, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		return visit(n), nil
	})
}
