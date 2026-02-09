package ast

import (
	"strings"
	"testing"

	"github.com/Daemon0x00000000/lql/internal/vm"
)

func TestComparisonNode_String(t *testing.T) {
	node := &ComparisonNode{
		left:        "field",
		right:       "value",
		operator:    vm.OP_EQ,
		operatorStr: EQUALS,
	}

	str := node.String()
	if str == "" {
		t.Error("ComparisonNode.String() should not be empty")
	}
	if !strings.Contains(str, "field") {
		t.Error("String should contain field name")
	}
	if !strings.Contains(str, "value") {
		t.Error("String should contain value")
	}
}

func TestComparisonNode_Compile(t *testing.T) {
	tests := []struct {
		name     string
		node     *ComparisonNode
		minBytes int
	}{
		{
			"simple equals",
			&ComparisonNode{
				left:        "field",
				right:       "value",
				operator:    vm.OP_EQ,
				operatorStr: EQUALS,
			},
			5, // LOAD_GLOBAL + PUSH + OP_EQ
		},
		{
			"with array (IN)",
			&ComparisonNode{
				left:        "status",
				right:       []interface{}{"a", "b", "c"},
				operator:    vm.OP_IN,
				operatorStr: IN,
			},
			10, // LOAD_GLOBAL + PUSH array + OP_IN
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytecode, err := tt.node.compile()
			assertNoError(t, err)

			if len(bytecode) < tt.minBytes {
				t.Errorf("expected at least %d bytes, got %d", tt.minBytes, len(bytecode))
			}
		})
	}
}

func TestLogicalNode_Compile(t *testing.T) {
	// Create simple comparison nodes for testing
	leftNode := &ComparisonNode{
		left:        "a",
		right:       "1",
		operator:    vm.OP_EQ,
		operatorStr: EQUALS,
	}
	rightNode := &ComparisonNode{
		left:        "b",
		right:       "2",
		operator:    vm.OP_EQ,
		operatorStr: EQUALS,
	}

	tests := []struct {
		name     string
		operator vm.OpCode
		opStr    LogicalOperator
	}{
		{"AND", vm.OP_AND, AND},
		{"OR", vm.OP_OR, OR},
		{"XOR", vm.OP_XOR, XOR},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &LogicalNode{
				left:        leftNode,
				right:       rightNode,
				operator:    tt.operator,
				operatorStr: tt.opStr,
			}

			bytecode, err := node.compile()
			assertNoError(t, err)

			// Should have bytecode from both children plus the logical operator
			if len(bytecode) < 10 {
				t.Errorf("bytecode too short: %d bytes", len(bytecode))
			}

			// Last byte should be the logical operator
			lastOpcode := vm.OpCode(bytecode[len(bytecode)-1])
			if lastOpcode != tt.operator {
				t.Errorf("expected operator %v, got %v", tt.operator, lastOpcode)
			}
		})
	}
}

func TestNotNode_Compile(t *testing.T) {
	operandNode := &ComparisonNode{
		left:        "field",
		right:       "value",
		operator:    vm.OP_EQ,
		operatorStr: EQUALS,
	}

	node := &NotNode{
		operand: operandNode,
	}

	bytecode, err := node.compile()
	assertNoError(t, err)

	// Should have operand bytecode plus OP_NOT
	if len(bytecode) < 5 {
		t.Errorf("bytecode too short: %d bytes", len(bytecode))
	}

	// Last byte should be OP_NOT
	lastOpcode := vm.OpCode(bytecode[len(bytecode)-1])
	if lastOpcode != vm.OP_NOT {
		t.Errorf("expected OP_NOT, got %v", lastOpcode)
	}
}

func TestTreeString(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"simple comparison", "field=value"},
		{"with OR", "a=1^ORb=2"},
		{"with NOT", "!(a=1)"},
		{"nested", "(a=1^ORb=2)^c=3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := Parse(tt.expr)
			assertNoError(t, err)

			str := treeString(ast.root, "", true)
			if str == "" {
				t.Error("treeString should not be empty")
			}

			// Should contain tree drawing characters
			if !strings.ContainsAny(str, "├└│") {
				t.Error("treeString should contain tree drawing characters")
			}
		})
	}
}

func TestLogicalNode_String(t *testing.T) {
	ast, err := Parse("a=1^ORb=2")
	assertNoError(t, err)

	str := ast.String()
	if str == "" {
		t.Error("LogicalNode.String() via AST should not be empty")
	}
}

func TestNestedNodes(t *testing.T) {
	// Test deeply nested logical operations
	left := &ComparisonNode{
		left:        "a",
		right:       "1",
		operator:    vm.OP_EQ,
		operatorStr: EQUALS,
	}

	middle := &ComparisonNode{
		left:        "b",
		right:       "2",
		operator:    vm.OP_EQ,
		operatorStr: EQUALS,
	}

	right := &ComparisonNode{
		left:        "c",
		right:       "3",
		operator:    vm.OP_EQ,
		operatorStr: EQUALS,
	}

	// (a=1 AND b=2) OR c=3
	innerNode := &LogicalNode{
		left:        left,
		right:       middle,
		operator:    vm.OP_AND,
		operatorStr: AND,
	}

	outerNode := &LogicalNode{
		left:        innerNode,
		right:       right,
		operator:    vm.OP_OR,
		operatorStr: OR,
	}

	bytecode, err := outerNode.compile()
	assertNoError(t, err)

	// Should compile successfully with reasonable size
	if len(bytecode) < 15 {
		t.Errorf("nested bytecode too short: %d bytes", len(bytecode))
	}
}

func TestNotNode_DoubleNegation(t *testing.T) {
	// Test !(!(a=1))
	innerComparison := &ComparisonNode{
		left:        "a",
		right:       "1",
		operator:    vm.OP_EQ,
		operatorStr: EQUALS,
	}

	innerNot := &NotNode{
		operand: innerComparison,
	}

	outerNot := &NotNode{
		operand: innerNot,
	}

	bytecode, err := outerNot.compile()
	assertNoError(t, err)

	// Should have comparison + 2x OP_NOT
	if len(bytecode) < 7 {
		t.Errorf("double NOT bytecode too short: %d bytes", len(bytecode))
	}

	// Last two bytes should be OP_NOT
	if vm.OpCode(bytecode[len(bytecode)-1]) != vm.OP_NOT {
		t.Error("expected OP_NOT at end")
	}
	if vm.OpCode(bytecode[len(bytecode)-2]) != vm.OP_NOT {
		t.Error("expected OP_NOT before end")
	}
}

func TestComparisonNode_WithDifferentTypes(t *testing.T) {
	tests := []struct {
		name  string
		right interface{}
	}{
		{"string", "value"},
		{"int", 123},
		{"array", []interface{}{"a", "b", "c"}},
		{"empty string", ""},
		{"single item array", []interface{}{"only"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &ComparisonNode{
				left:        "field",
				right:       tt.right,
				operator:    vm.OP_EQ,
				operatorStr: EQUALS,
			}

			bytecode, err := node.compile()
			assertNoError(t, err)
			assertBytecodeNotEmpty(t, bytecode)
		})
	}
}

func TestLogicalNode_AsymmetricChildren(t *testing.T) {
	// Left is simple comparison, right is nested logical
	leftNode := &ComparisonNode{
		left:        "a",
		right:       "1",
		operator:    vm.OP_EQ,
		operatorStr: EQUALS,
	}

	rightLeft := &ComparisonNode{
		left:        "b",
		right:       "2",
		operator:    vm.OP_EQ,
		operatorStr: EQUALS,
	}

	rightRight := &ComparisonNode{
		left:        "c",
		right:       "3",
		operator:    vm.OP_EQ,
		operatorStr: EQUALS,
	}

	rightNode := &LogicalNode{
		left:        rightLeft,
		right:       rightRight,
		operator:    vm.OP_AND,
		operatorStr: AND,
	}

	// a=1 OR (b=2 AND c=3)
	node := &LogicalNode{
		left:        leftNode,
		right:       rightNode,
		operator:    vm.OP_OR,
		operatorStr: OR,
	}

	bytecode, err := node.compile()
	assertNoError(t, err)

	// Should compile successfully
	if len(bytecode) < 15 {
		t.Errorf("asymmetric bytecode too short: %d bytes", len(bytecode))
	}
}
