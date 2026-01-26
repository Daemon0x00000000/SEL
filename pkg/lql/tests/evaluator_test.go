package lql_test

import (
	"testing"

	"github.com/Daemon0x0000000/lql/pkg/lql"
)

// =============================================================================
// Tests des noeuds
// =============================================================================

func TestComparisonNode_String(t *testing.T) {
	op := lql.OperationResult(lql.True)
	node := lql.NewComparisonNode(&op, "field", lql.EQUALS, []string{"value"})

	str := node.String()
	if str == "" {
		t.Error("ComparisonNode.String() should not be empty")
	}
}

func TestLogicalNode_String(t *testing.T) {
	ast, err := lql.Parse("a=1^ORb=2")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	str := ast.String()
	if str == "" {
		t.Error("LogicalNode.String() should not be empty")
	}
}

func TestOperationResult_Compute(t *testing.T) {
	tests := []struct {
		name      string
		predicate lql.PredicateClosure
		value     interface{}
		want      lql.OperationResult
	}{
		{
			"predicate returns true",
			func(v interface{}) bool { return v == "match" },
			"match",
			lql.True,
		},
		{
			"predicate returns false",
			func(v interface{}) bool { return v == "match" },
			"no-match",
			lql.False,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := lql.OperationResult(lql.NotResolved)
			compute := res.Compute(tt.predicate)
			compute(tt.value)
			if res != tt.want {
				t.Errorf("Compute result = %v, want %v", res, tt.want)
			}
		})
	}
}
