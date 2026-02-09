package vm

import (
	"testing"
)

// ============================================================================
// Logical Operators Tests
// ============================================================================

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		name     string
		handler  func(*VM) error
		left     bool
		right    bool
		expected bool
	}{
		// AND
		{"and: true && true", (*VM).andHandler, true, true, true},
		{"and: true && false", (*VM).andHandler, true, false, false},
		{"and: false && true", (*VM).andHandler, false, true, false},
		{"and: false && false", (*VM).andHandler, false, false, false},

		// OR
		{"or: true || true", (*VM).orHandler, true, true, true},
		{"or: true || false", (*VM).orHandler, true, false, true},
		{"or: false || true", (*VM).orHandler, false, true, true},
		{"or: false || false", (*VM).orHandler, false, false, false},

		// XOR
		{"xor: true ^ true", (*VM).xorHandler, true, true, false},
		{"xor: true ^ false", (*VM).xorHandler, true, false, true},
		{"xor: false ^ true", (*VM).xorHandler, false, true, true},
		{"xor: false ^ false", (*VM).xorHandler, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left := Value{Type: TYPE_BOOL, Bool: tt.left}
			right := Value{Type: TYPE_BOOL, Bool: tt.right}
			testBinaryHandler(t, tt.handler, left, right, TYPE_BOOL, func(v Value) bool {
				return v.Bool == tt.expected
			})
		})
	}
}

func TestNotHandler(t *testing.T) {
	tests := []struct {
		input    bool
		expected bool
	}{
		{true, false},
		{false, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			vm := NewVM([]byte{}, nil)
			vm.push(Value{Type: TYPE_BOOL, Bool: tt.input})

			if err := vm.notHandler(); err != nil {
				t.Fatalf("notHandler failed: %v", err)
			}

			assertStackValue(t, vm, TYPE_BOOL, func(v Value) bool {
				return v.Bool == tt.expected
			})
		})
	}
}

func TestNotHandler_TypeError(t *testing.T) {
	vm := NewVM([]byte{}, nil)
	vm.push(Value{Type: TYPE_INT8, Int8: 42})

	if err := vm.notHandler(); err == nil {
		t.Fatal("Expected error for NOT on non-boolean, got nil")
	}
}

// ============================================================================
// OpCode Methods Tests
// ============================================================================

func TestOpCodeMethods(t *testing.T) {
	logicalTests := []struct {
		op       OpCode
		expected bool
	}{
		{OP_AND, true},
		{OP_OR, true},
		{OP_XOR, true},
		{OP_NOT, true},
		{OP_EQ, false},
		{OP_GT, false},
		{PUSH, false},
	}

	for _, tt := range logicalTests {
		if got := tt.op.isLogical(); got != tt.expected {
			t.Errorf("OpCode(%d).isLogical() = %v, want %v", tt.op, got, tt.expected)
		}
	}

	comparisonTests := []struct {
		op       OpCode
		expected bool
	}{
		{OP_EQ, true},
		{OP_GT, true},
		{OP_LT, true},
		{OP_GTE, true},
		{OP_LTE, true},
		{OP_STARTSWITH, true},
		{OP_ENDSWITH, true},
		{OP_CONTAINS, true},
		{OP_IN, true},
		{OP_AND, false},
		{PUSH, false},
	}

	for _, tt := range comparisonTests {
		if got := tt.op.isComparison(); got != tt.expected {
			t.Errorf("OpCode(%d).isComparison() = %v, want %v", tt.op, got, tt.expected)
		}
	}
}
