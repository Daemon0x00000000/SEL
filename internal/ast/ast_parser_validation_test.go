package ast

import (
	"strings"
	"testing"
)

func TestValidateParentheses_Balanced(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"no parens", "a=1"},
		{"simple parens", "(a=1)"},
		{"double parens", "((a=1))"},
		{"multiple groups", "(a=1)^(b=2)"},
		{"nested left", "((a=1)^b=2)"},
		{"nested right", "(a=1^(b=2))"},
		{"complex nested", "((a=1^b=2)^(c=3^d=4))"},
		{"empty expr", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateParentheses(tt.expr)
			if err != nil {
				t.Errorf("validateParentheses(%q) = %v, want nil", tt.expr, err)
			}
		})
	}
}

func TestValidateParentheses_Unbalanced(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		errContains string
	}{
		{"missing close", "(a=1", "unclosed"},
		{"missing open", "a=1)", "closing ')'"},
		{"double open", "((a=1)", "unclosed"},
		{"double close", "(a=1))", "closing ')'"},
		{"wrong order", ")(", "closing ')'"},
		{"nested wrong", "((a=1)", "unclosed"},
		{"multiple unclosed", "(((a=1)", "unclosed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateParentheses(tt.expr)
			if err == nil {
				t.Errorf("validateParentheses(%q) should fail", tt.expr)
				return
			}
			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("error should contain %q, got: %v", tt.errContains, err)
			}
		})
	}
}

func TestValidateParentheses_ErrorMessages(t *testing.T) {
	tests := []struct {
		name              string
		expr              string
		wantPositionInMsg bool
		wantContextInMsg  bool
	}{
		{"unclosed paren", "(a=1", true, true},
		{"extra closing", "a=1)", true, true},
		{"multiple unclosed", "((a=1)", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateParentheses(tt.expr)
			if err == nil {
				t.Fatal("expected error")
			}

			errMsg := err.Error()

			if tt.wantPositionInMsg && !strings.Contains(errMsg, "position") {
				t.Error("error message should mention position")
			}

			if tt.wantContextInMsg && !strings.Contains(errMsg, "Context:") {
				t.Error("error message should include context")
			}
		})
	}
}

func TestStripOuterParens(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
		stripped bool
	}{
		{"simple parens", "(a=1)", "a=1", true},
		{"double parens", "((a=1))", "(a=1)", true},
		{"no outer parens", "a=1", "a=1", false},
		{"non-matching", "(a=1)^(b=2)", "(a=1)^(b=2)", false},
		{"closing before end", "(a=1)^b=2", "(a=1)^b=2", false},
		{"empty parens", "()", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, stripped := stripOuterParens(tt.expr)
			if result != tt.expected {
				t.Errorf("stripOuterParens(%q) = %q, want %q", tt.expr, result, tt.expected)
			}
			if stripped != tt.stripped {
				t.Errorf("stripOuterParens(%q) stripped = %v, want %v", tt.expr, stripped, tt.stripped)
			}
		})
	}
}

func TestFindOperatorOutsideParens(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		operator string
		expected int
	}{
		{"simple", "a=1^ORb=2", "^OR", 3},
		{"not in parens", "(a=1^OR(b=2))^ORc=3", "^OR", 13},
		{"at start", "^ORa=1", "^OR", 0},
		{"not found", "a=1", "^OR", -1},
		{"in parens only", "(a=1^ORb=2)", "^OR", -1},
		{"multiple occurrences", "a=1^ORb=2^ORc=3", "^OR", 3},
		{"nested parens", "((a=1)^ORb=2)^ORc=3", "^OR", 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findOperatorOutsideParens(tt.expr, tt.operator)
			if result != tt.expected {
				t.Errorf("findOperatorOutsideParens(%q, %q) = %d, want %d",
					tt.expr, tt.operator, result, tt.expected)
			}
		})
	}
}

func TestFindOperatorOutsideParens_WithQuotes(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		operator string
		expected int
	}{
		{"in quotes", "field='value^OR'", "^OR", -1},
		{"after quotes", "field='value'^ORb=2", "^OR", 13},
		{"before quotes", "a=1^ORfield='value'", "^OR", 3},
		{"escaped quote", "field='it\\'s'^ORb=2", "^OR", 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findOperatorOutsideParens(tt.expr, tt.operator)
			if result != tt.expected {
				t.Errorf("findOperatorOutsideParens(%q, %q) = %d, want %d",
					tt.expr, tt.operator, result, tt.expected)
			}
		})
	}
}

func TestParseValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"simple list", "a,b,c", []string{"a", "b", "c"}},
		{"with spaces", "a, b, c", []string{"a", "b", "c"}},
		{"single value", "only", []string{"only"}},
		{"quoted", "'a','b','c'", []string{"a", "b", "c"}},
		{"quoted with spaces", "'a b','c d'", []string{"a b", "c d"}},
		{"mixed", "a,'b c',d", []string{"a", "b c", "d"}},
		{"empty value", "a,,c", []string{"a", "", "c"}},
		{"trailing comma", "a,b,", []string{"a", "b", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseValues(tt.input)
			assertNoError(t, err)

			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("value[%d]: got %q, want %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestParseValues_EscapeSequences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"escaped quote", "'it\\'s'", []string{"it's"}},
		{"escaped backslash", "'path\\\\to\\\\file'", []string{"path\\to\\file"}},
		{"escaped newline", "'line1\\nline2'", []string{"line1\nline2"}},
		{"escaped tab", "'col1\\tcol2'", []string{"col1\tcol2"}},
		{"escaped r", "'text\\rmore'", []string{"text\rmore"}},
		{"unknown escape", "'test\\x'", []string{"testx"}}, // unknown escapes pass through
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseValues(tt.input)
			assertNoError(t, err)

			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}

			if result[0] != tt.expected[0] {
				t.Errorf("got %q, want %q", result[0], tt.expected[0])
			}
		})
	}
}

func TestParseValues_Errors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		errContains string
	}{
		{"unclosed quote", "'value", "unclosed quote"},
		{"incomplete escape", "'value\\", "incomplete escape"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseValues(tt.input)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("error should contain %q, got: %v", tt.errContains, err)
			}
		})
	}
}

func TestParseValues_PreservesSpaces(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"leading space in quote", "' value'", " value"},
		{"trailing space in quote", "'value '", "value "},
		{"both spaces in quote", "' value '", " value "},
		{"no quotes trims", " value ", "value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseValues(tt.input)
			assertNoError(t, err)

			if result[0] != tt.expected {
				t.Errorf("got %q, want %q", result[0], tt.expected)
			}
		})
	}
}

func TestParseComparison(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr bool
	}{
		{"simple equals", "field=value", false},
		{"with IN", "fieldINa,b,c", false},
		{"no operator", "field", true},
		{"only operator", "=value", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseComparison(tt.expr)
			if tt.wantErr {
				assertError(t, err)
			} else {
				assertNoError(t, err)
			}
		})
	}
}

func TestParseExpr_RecursiveStripping(t *testing.T) {
	// Test that multiple layers of parens are stripped correctly
	tests := []struct {
		name string
		expr string
	}{
		{"single layer", "(a=1)"},
		{"double layer", "((a=1))"},
		{"triple layer", "(((a=1)))"},
		{"quadruple layer", "((((a=1))))"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseExpr(tt.expr)
			assertNoError(t, err)

			// Should parse to a comparison node, not wrapped in logical nodes
			_, isComparison := node.(*ComparisonNode)
			if !isComparison {
				t.Error("excessive parens should be stripped, resulting in ComparisonNode")
			}
		})
	}
}
