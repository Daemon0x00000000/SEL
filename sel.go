package sel

import (
	"fmt"

	iast "github.com/Daemon0x00000000/sel/internal/ast"
	"github.com/Daemon0x00000000/sel/internal/vm"
)

type Expression struct {
	vm *vm.VM
}

func (expr *Expression) Parse(expression string) error {
	ast, err := iast.Parse(expression)
	if err != nil {
		return err
	}
	bytes, err := ast.Compile()
	if err != nil {
		return err
	}

	// TODO: Get Native Funcs from AST
	expr.vm = vm.NewVM(bytes, []vm.NativeFunc{})
	return nil
}

func (expr *Expression) Eval(data map[string]interface{}) (bool, error) {
	if expr.vm == nil {
		return false, fmt.Errorf("expression not parsed yet")
	}
	expr.vm.Reset()

	err := expr.vm.LoadRecords(data)
	if err != nil {
		return false, err
	}

	err = expr.vm.Execute()
	if err != nil || len(expr.vm.DataStack()) != 1 {
		return false, err
	}
	return expr.vm.DataStack()[0].Bool, nil
}
