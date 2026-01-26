package lql_test

import (
	"testing"

	"github.com/Daemon0x0000000/lql/pkg/lql"
)

// =============================================================================
// Tests des opérateurs de comparaison
// =============================================================================

func TestComparisonOperator_Equals(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"string match", "name=hello", map[lql.Field]interface{}{"name": "hello"}, true},
		{"string no match", "name=hello", map[lql.Field]interface{}{"name": "world"}, false},
		{"int match", "age=123", map[lql.Field]interface{}{"age": 123}, true},
		{"int no match", "age=123", map[lql.Field]interface{}{"age": 456}, false},
		{"empty string match", "field=", map[lql.Field]interface{}{"field": ""}, true},
		{"empty vs non-empty", "field=", map[lql.Field]interface{}{"field": "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("EQUALS eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_NotEquals(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"string different", "name!=hello", map[lql.Field]interface{}{"name": "world"}, true},
		{"string same", "name!=hello", map[lql.Field]interface{}{"name": "hello"}, false},
		{"int different", "age!=123", map[lql.Field]interface{}{"age": 456}, true},
		{"int same", "age!=123", map[lql.Field]interface{}{"age": 123}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("NOT_EQUALS eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_In(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"value in set", "statusINa,b,c", map[lql.Field]interface{}{"status": "b"}, true},
		{"value not in set", "statusINa,b,c", map[lql.Field]interface{}{"status": "d"}, false},
		{"int in set", "ageIN1,2,3", map[lql.Field]interface{}{"age": 2}, true},
		{"single value match", "fieldINonly", map[lql.Field]interface{}{"field": "only"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("IN eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_NotIn(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"value not in set", "statusNOTINa,b,c", map[lql.Field]interface{}{"status": "d"}, true},
		{"value in set", "statusNOTINa,b,c", map[lql.Field]interface{}{"status": "b"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("NOT_IN eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_Contains(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"substring found", "textCONTAINSell", map[lql.Field]interface{}{"text": "hello"}, true},
		{"substring not found", "textCONTAINSxyz", map[lql.Field]interface{}{"text": "hello"}, false},
		{"exact match", "textCONTAINShello", map[lql.Field]interface{}{"text": "hello"}, true},
		{"search in empty", "textCONTAINSa", map[lql.Field]interface{}{"text": ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("CONTAINS eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_DoesNotContain(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"substring not found", "textDOESNOTCONTAINxyz", map[lql.Field]interface{}{"text": "hello"}, true},
		{"substring found", "textDOESNOTCONTAINell", map[lql.Field]interface{}{"text": "hello"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("DOES_NOT_CONTAIN eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_StartsWith(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"prefix match", "textSTARTSWITHhel", map[lql.Field]interface{}{"text": "hello"}, true},
		{"prefix no match", "textSTARTSWITHwor", map[lql.Field]interface{}{"text": "hello"}, false},
		{"exact match", "textSTARTSWITHhello", map[lql.Field]interface{}{"text": "hello"}, true},
		{"longer prefix", "textSTARTSWITHhello world", map[lql.Field]interface{}{"text": "hello"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("STARTS_WITH eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_EndsWith(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"suffix match", "textENDSWITHllo", map[lql.Field]interface{}{"text": "hello"}, true},
		{"suffix no match", "textENDSWITHwor", map[lql.Field]interface{}{"text": "hello"}, false},
		{"exact match", "textENDSWITHhello", map[lql.Field]interface{}{"text": "hello"}, true},
		{"longer suffix", "textENDSWITHsay hello", map[lql.Field]interface{}{"text": "hello"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("ENDS_WITH eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_GreaterThan(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"string greater", "field>a", map[lql.Field]interface{}{"field": "b"}, true},
		{"string less", "field>b", map[lql.Field]interface{}{"field": "a"}, false},
		{"string equal", "field>a", map[lql.Field]interface{}{"field": "a"}, false},
		{"number greater", "age>1", map[lql.Field]interface{}{"age": "2"}, true},
		{"number less", "age>2", map[lql.Field]interface{}{"age": "1"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("GREATER_THAN eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_LessThan(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"string less", "field<b", map[lql.Field]interface{}{"field": "a"}, true},
		{"string greater", "field<a", map[lql.Field]interface{}{"field": "b"}, false},
		{"string equal", "field<a", map[lql.Field]interface{}{"field": "a"}, false},
		{"number less", "age<2", map[lql.Field]interface{}{"age": "1"}, true},
		{"number greater", "age<1", map[lql.Field]interface{}{"age": "2"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("LESS_THAN eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_GreaterThanOrEqual(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"string greater", "field>=a", map[lql.Field]interface{}{"field": "b"}, true},
		{"string equal", "field>=a", map[lql.Field]interface{}{"field": "a"}, true},
		{"string less", "field>=b", map[lql.Field]interface{}{"field": "a"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("GREATER_THAN_OR_EQUAL eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_LessThanOrEqual(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"string less", "field<=b", map[lql.Field]interface{}{"field": "a"}, true},
		{"string equal", "field<=a", map[lql.Field]interface{}{"field": "a"}, true},
		{"string greater", "field<=a", map[lql.Field]interface{}{"field": "b"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("LESS_THAN_OR_EQUAL eval = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparisonOperator_Matches(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[lql.Field]interface{}
		want bool
	}{
		{"simple match", "textMATCHES'hello'", map[lql.Field]interface{}{"text": "hello world"}, true},
		{"regex digit", "codeMATCHES'^[0-9]+$'", map[lql.Field]interface{}{"code": "12345"}, true},
		{"regex digit no match", "codeMATCHES'^[0-9]+$'", map[lql.Field]interface{}{"code": "abc123"}, false},
		{"email pattern", `emailMATCHES'^[a-z]+@[a-z]+\.[a-z]+$'`, map[lql.Field]interface{}{"email": "user@example.com"}, true},
		{"email pattern no match", `emailMATCHES'^[a-z]+@[a-z]+\.[a-z]+$'`, map[lql.Field]interface{}{"email": "invalid.email"}, false},
		{"starts with pattern", "textMATCHES'^hello'", map[lql.Field]interface{}{"text": "hello world"}, true},
		{"ends with pattern", "textMATCHES'world$'", map[lql.Field]interface{}{"text": "hello world"}, true},
		{"case sensitive", "textMATCHES'Hello'", map[lql.Field]interface{}{"text": "hello"}, false},
		{"dot matches any", "textMATCHES'h.llo'", map[lql.Field]interface{}{"text": "hello"}, true},
		{"dot matches any no match", "textMATCHES'h.llo'", map[lql.Field]interface{}{"text": "hllo"}, false},
		{"alternation", "textMATCHES'cat|dog'", map[lql.Field]interface{}{"text": "I have a dog"}, true},
		{"alternation no match", "textMATCHES'cat|dog'", map[lql.Field]interface{}{"text": "I have a bird"}, false},
		{"quantifier star", "textMATCHES'go*gle'", map[lql.Field]interface{}{"text": "google"}, true},
		{"quantifier plus", "textMATCHES'go+gle'", map[lql.Field]interface{}{"text": "google"}, true},
		{"quantifier plus no match", "textMATCHES'go+gle'", map[lql.Field]interface{}{"text": "ggle"}, false},
		{"character class", "textMATCHES'[aeiou]'", map[lql.Field]interface{}{"text": "hello"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("MATCHES eval = %v, want %v", got, tt.want)
			}
		})
	}
}
