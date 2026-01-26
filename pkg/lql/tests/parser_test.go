package lql_test

import (
	"testing"

	"github.com/Daemon0x0000000/lql/pkg/lql"
)

// =============================================================================
// Tests de validation des parenthèses
// =============================================================================

func TestValidateParentheses(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr bool
	}{
		{"balanced simple", "(a=1)", false},
		{"balanced nested", "((a=1))", false},
		{"balanced complex", "(a=1^OR(b=2^ANDc=3))", false},
		{"no parens", "a=1", false},
		{"unbalanced open", "(a=1", true},
		{"unbalanced close", "a=1)", true},
		{"multiple unbalanced", "((a=1)", true},
		{"wrong order", ")(", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := lql.Parse(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
			}
		})
	}
}

// =============================================================================
// Tests du parser
// =============================================================================

func TestParse_SimpleExpression(t *testing.T) {
	ast, err := lql.Parse("field=value")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if ast == nil {
		t.Fatal("AST is nil")
	}
}

func TestParse_ComplexExpression(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"simple equals", "field=value"},
		{"with OR", "a=1^ORb=2"},
		{"with AND", "a=1^b=2"},
		{"with XOR", "a=1^XORb=2"},
		{"with NOR", "a=1^NORb=2"},
		{"with NAND", "a=1^NANDb=2"},
		{"with XNOR", "a=1^XNORb=2"},
		{"nested parens", "(a=1^ORb=2)"},
		{"complex nested", "a=1^OR(b=2^c=3)"},
		{"IN operator", "fieldINa,b,c"},
		{"CONTAINS operator", "fieldCONTAINSvalue"},
		{"STARTSWITH operator", "fieldSTARTSWITHprefix"},
		{"ENDSWITH operator", "fieldENDSWITHsuffix"},
		{"NOT_EQUALS operator", "field!=value"},
		{"GREATER_THAN operator", "field>value"},
		{"LESS_THAN operator", "field<value"},
		{"MATCHES operator", `emailMATCHES'^[a-z]+@[a-z]+\.com$'`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", tt.expr, err)
			}
			if ast == nil {
				t.Errorf("Parse(%q) returned nil AST", tt.expr)
			}
		})
	}
}

func TestParse_InvalidExpression(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"unbalanced parens", "(a=1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := lql.Parse(tt.expr)
			if err == nil {
				t.Errorf("Parse(%q) should have failed", tt.expr)
			}
		})
	}
}
