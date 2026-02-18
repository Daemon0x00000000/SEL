package vm

import (
	"testing"
)

// ============================================================================
// Test Helpers
// ============================================================================

// Helper pour comparer des bytecodes
func assertBytecodeEqual(t *testing.T, got, expected []byte) {
	t.Helper()
	if len(got) != len(expected) {
		t.Fatalf("Length mismatch: expected %d, got %d", len(expected), len(got))
	}
	for i, b := range expected {
		if got[i] != b {
			t.Errorf("Byte mismatch at index %d: expected 0x%02x, got 0x%02x", i, b, got[i])
		}
	}
}

// Helper pour vérifier une valeur sur la stack
func assertStackValue(t *testing.T, vm *VM, expectedType Type, checker func(Value) bool) {
	t.Helper()
	if len(vm.dataStack) != 1 {
		t.Fatalf("Expected 1 value on stack, got %d", len(vm.dataStack))
	}
	val := vm.dataStack[0]
	if val.Type != expectedType {
		t.Errorf("Expected type %v, got %v", expectedType, val.Type)
	}
	if !checker(val) {
		t.Errorf("Value check failed for %+v", val)
	}
}

// Helper pour tester un handler avec setup de stack
func testBinaryHandler(t *testing.T, handler func(*VM) error, left, right Value, expectedType Type, checker func(Value) bool) {
	t.Helper()
	vm := NewVM([]byte{}, nil)
	vm.push(left)
	vm.push(right)

	err := handler(vm)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	assertStackValue(t, vm, expectedType, checker)
}
