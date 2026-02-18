package ast

import "fmt"

type AST struct {
	root Node
}

func (ast *AST) String() string {
	return treeString(ast.root, "", true)
}

func (ast *AST) Compile() ([]byte, error) {
	if ast.root == nil {
		return nil, fmt.Errorf("cannot compile AST with nil root")
	}
	return ast.root.compile()
}

func newAST() *AST {
	return &AST{}
}
