package lql

import "fmt"

type AST struct {
	fieldsOperations map[Field][]OperationCompute
	root             *LogicalNode
}

func (ast *AST) String() string {
	return fmt.Sprintf("%v", ast.root)
}

func (ast *AST) Eval(data map[Field]interface{}) bool {
	if ast.root == nil {
		return false
	}

	for field, operations := range ast.fieldsOperations {
		value, ok := data[field]
		if !ok {
			continue
		}

		for _, compute := range operations {
			compute(value)
		}
	}

	return ast.root.Eval()
}

func newAST() *AST {
	return &AST{
		fieldsOperations: make(map[Field][]OperationCompute),
	}
}
