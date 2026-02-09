package ast

import (
	"testing"

	"github.com/Daemon0x00000000/lql/internal/vm"
)

func TestAST_Compile_SimpleComparison(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"equals", "field=value"},
		{"not equals", "field!=value"},
		{"greater than", "age>18"},
		{"less than", "age<65"},
		{"gte", "age>=18"},
		{"lte", "age<=65"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytecode := compileAndCheck(t, tt.expr)

			// Bytecode should contain:
			// - LOAD_GLOBAL for field
			// - PUSH for value
			// - Comparison operator
			if len(bytecode) < 5 {
				t.Errorf("bytecode too short: %d bytes", len(bytecode))
			}
		})
	}
}

func TestAST_Compile_StringOperators(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"contains", "textCONTAINShello"},
		{"starts with", "textSTARTSWITHhello"},
		{"ends with", "textENDSWITHworld"},
		{"in", "statusINa,b,c"},
		{"not in", "status!INx,y,z"},
		// TODO: MATCHES operator not fully implemented yet
		// {"matches", "emailMATCHES'^[a-z]+@'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytecode := compileAndCheck(t, tt.expr)

			// Should have reasonable size
			if len(bytecode) == 0 {
				t.Error("bytecode should not be empty")
			}
		})
	}
}

func TestAST_Compile_LogicalOperators(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"and", "a=1^b=2"},
		{"or", "a=1^ORb=2"},
		{"xor", "a=1^XORb=2"},
		{"not", "!(a=1)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytecode := compileAndCheck(t, tt.expr)

			// Logical ops should have bytecode from both sides + operator
			if tt.name != "not" && len(bytecode) < 10 {
				t.Errorf("bytecode too short for logical op: %d bytes", len(bytecode))
			}
		})
	}
}

func TestAST_Compile_NestedExpressions(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"simple parens", "(a=1)"},
		{"nested parens", "((a=1))"},
		{"parens with or", "(a=1^ORb=2)"},
		{"left nested", "(a=1^b=2)^c=3"},
		{"right nested", "a=1^(b=2^c=3)"},
		{"both nested", "(a=1^b=2)^(c=3^d=4)"},
		{"deep nesting", "((a=1^b=2)^(c=3^d=4))^e=5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytecode := compileAndCheck(t, tt.expr)

			// Nested expressions should produce longer bytecode
			if len(bytecode) < 5 {
				t.Errorf("bytecode too short: %d bytes", len(bytecode))
			}
		})
	}
}

func TestAST_Compile_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{
			"servicenow style",
			"sys_id=123^OR(active=true^category=incident)",
		},
		{
			"multiple ors",
			"a=1^ORb=2^ORc=3^ORd=4",
		},
		{
			"multiple ands",
			"a=1^b=2^c=3^d=4",
		},
		{
			"mixed operators",
			"a=1^b=2^ORc=3^d=4",
		},
		{
			"with not",
			"!(a=1)^ORb=2",
		},
		{
			"complex nested",
			"(a=1^ORb=2)^(c=3^ORd=4)^OR(e=5^f=6)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytecode := compileAndCheck(t, tt.expr)

			// Complex expressions should produce substantial bytecode
			if len(bytecode) < 15 {
				t.Errorf("bytecode too short for complex expression: %d bytes", len(bytecode))
			}
		})
	}
}

func TestAST_Compile_INOperatorWithArray(t *testing.T) {
	tests := []struct {
		name   string
		expr   string
		values int // expected number of values in array
	}{
		{"single value", "fieldINa", 1},
		{"two values", "fieldINa,b", 2},
		{"three values", "fieldINa,b,c", 3},
		{"five values", "fieldINa,b,c,d,e", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytecode := compileAndCheck(t, tt.expr)

			// IN operator should generate:
			// - LOAD_GLOBAL
			// - PUSH array with N values
			// - OP_IN
			if len(bytecode) < 5 {
				t.Errorf("bytecode too short: %d bytes", len(bytecode))
			}
		})
	}
}

func TestAST_Compile_WithQuotedValues(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"simple quote", "field='value'"},
		{"quote with space", "field='value with spaces'"},
		{"quote with comma", "field='a,b,c'"},
		{"escaped quote", "field='it\\'s here'"},
		{"escaped newline", "field='line1\\nline2'"},
		{"quoted in IN", "fieldIN'a','b','c'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytecode := compileAndCheck(t, tt.expr)

			if len(bytecode) == 0 {
				t.Error("bytecode should not be empty")
			}
		})
	}
}

func TestAST_Compile_BinaryOperatorStructure(t *testing.T) {
	// Test that binary logical operators compile left and right operands
	expr := "a=1^b=2"

	ast, err := Parse(expr)
	assertNoError(t, err)

	bytecode, err := ast.Compile()
	assertNoError(t, err)

	// Should contain opcodes for:
	// 1. LOAD_GLOBAL a
	// 2. PUSH 1
	// 3. OP_EQ
	// 4. LOAD_GLOBAL b
	// 5. PUSH 2
	// 6. OP_EQ
	// 7. OP_AND

	// At minimum we need these opcodes
	if len(bytecode) < 7 {
		t.Errorf("expected at least 7 bytes, got %d", len(bytecode))
	}

	// First opcode should be LOAD_GLOBAL
	if vm.OpCode(bytecode[0]) != vm.LOAD_GLOBAL {
		t.Errorf("first opcode should be LOAD_GLOBAL, got %v", vm.OpCode(bytecode[0]))
	}
}

func TestAST_Compile_UnaryNOTStructure(t *testing.T) {
	expr := "!(a=1)"

	ast, err := Parse(expr)
	assertNoError(t, err)

	bytecode, err := ast.Compile()
	assertNoError(t, err)

	// Should contain opcodes for:
	// 1. LOAD_GLOBAL a
	// 2. PUSH 1
	// 3. OP_EQ
	// 4. OP_NOT

	if len(bytecode) < 4 {
		t.Errorf("expected at least 4 bytes, got %d", len(bytecode))
	}

	// Last opcode should be OP_NOT
	lastOpcode := vm.OpCode(bytecode[len(bytecode)-1])
	if lastOpcode != vm.OP_NOT {
		t.Errorf("last opcode should be OP_NOT, got %v", lastOpcode)
	}
}

func TestAST_Compile_ErrorCases(t *testing.T) {
	// Test that invalid AST states are handled
	// (This is more of a defensive test)

	t.Run("empty AST", func(t *testing.T) {
		ast := &AST{
			root: nil,
		}

		_, err := ast.Compile()
		if err == nil {
			t.Error("expected error when compiling nil root")
		}
	})
}

func TestAST_Compile_Deterministic(t *testing.T) {
	// Test that compiling the same expression twice produces identical bytecode
	expr := "a=1^ORb=2"

	bytecode1 := compileAndCheck(t, expr)
	bytecode2 := compileAndCheck(t, expr)

	if len(bytecode1) != len(bytecode2) {
		t.Errorf("bytecode length mismatch: %d vs %d", len(bytecode1), len(bytecode2))
	}

	for i := range bytecode1 {
		if bytecode1[i] != bytecode2[i] {
			t.Errorf("bytecode differs at position %d: %d vs %d", i, bytecode1[i], bytecode2[i])
		}
	}
}

func TestAST_Compile_AllComparisonOperators(t *testing.T) {
	operators := []struct {
		name string
		expr string
	}{
		{"EQUALS", "field=value"},
		{"GREATER_THAN", "field>value"},
		{"LESS_THAN", "field<value"},
		{"GREATER_THAN_OR_EQUAL", "field>=value"},
		{"LESS_THAN_OR_EQUAL", "field<=value"},
		{"STARTS_WITH", "fieldSTARTSWITHvalue"},
		{"ENDS_WITH", "fieldENDSWITHvalue"},
		{"CONTAINS", "fieldCONTAINSvalue"},
		{"IN", "fieldINa,b,c"},
		// TODO: MATCHES operator not fully implemented yet
		// {"MATCHES", "fieldMATCHES'^test$'"},
	}

	for _, op := range operators {
		t.Run(op.name, func(t *testing.T) {
			bytecode := compileAndCheck(t, op.expr)

			// All comparison operators should produce valid bytecode
			if len(bytecode) < 5 {
				t.Errorf("%s: bytecode too short: %d bytes", op.name, len(bytecode))
			}
		})
	}
}

func TestAST_Compile_AllLogicalOperators(t *testing.T) {
	operators := []struct {
		name string
		expr string
	}{
		{"AND", "a=1^b=2"},
		{"OR", "a=1^ORb=2"},
		{"XOR", "a=1^XORb=2"},
		{"NOT", "!(a=1)"},
	}

	for _, op := range operators {
		t.Run(op.name, func(t *testing.T) {
			bytecode := compileAndCheck(t, op.expr)

			// All logical operators should produce valid bytecode
			if len(bytecode) < 5 {
				t.Errorf("%s: bytecode too short: %d bytes", op.name, len(bytecode))
			}
		})
	}
}
