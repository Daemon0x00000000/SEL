package ast

import (
	"testing"

	"github.com/Daemon0x00000000/lql/internal/vm"
)

// assertNoError vérifie qu'il n'y a pas d'erreur
func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// assertError vérifie qu'il y a une erreur
func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// assertASTNotNil vérifie que l'AST n'est pas nil
func assertASTNotNil(t *testing.T, ast *AST) {
	t.Helper()
	if ast == nil {
		t.Fatal("AST should not be nil")
	}
}

// assertBytecodeNotEmpty vérifie que le bytecode n'est pas vide
func assertBytecodeNotEmpty(t *testing.T, bytecode []byte) {
	t.Helper()
	if len(bytecode) == 0 {
		t.Fatal("bytecode should not be empty")
	}
}

// compileAndCheck parse et compile une expression
func compileAndCheck(t *testing.T, expr string) []byte {
	t.Helper()

	ast, err := Parse(expr)
	assertNoError(t, err)
	assertASTNotNil(t, ast)

	bytecode, err := ast.Compile()
	assertNoError(t, err)
	assertBytecodeNotEmpty(t, bytecode)

	return bytecode
}

// executeInVM exécute du bytecode dans la VM avec les données fournies
func executeInVM(t *testing.T, bytecode []byte, data map[string]interface{}) bool {
	t.Helper()

	vmInstance := vm.NewVM(bytecode, nil)

	// Convert map[string]interface{} to map[vm.Field]interface{}
	vmData := make(map[string]interface{})
	for k, v := range data {
		vmData[k] = v
	}

	err := vmInstance.LoadRecords(vmData)
	assertNoError(t, err)

	err = vmInstance.Execute()
	assertNoError(t, err)

	stack := vmInstance.DataStack()
	if len(stack) != 1 {
		t.Fatalf("expected 1 value on stack, got %d", len(stack))
	}

	if stack[0].Type != vm.TYPE_BOOL {
		t.Fatalf("expected bool on stack, got type %v", stack[0].Type)
	}

	return stack[0].Bool
}

// testParseCompileEval teste le cycle complet pour une expression
func testParseCompileEval(t *testing.T, expr string, data map[string]interface{}, expected bool) {
	t.Helper()

	bytecode := compileAndCheck(t, expr)
	result := executeInVM(t, bytecode, data)

	if result != expected {
		t.Errorf("Eval(%q, %v) = %v, want %v", expr, data, result, expected)
	}
}
