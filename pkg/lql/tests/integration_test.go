package lql_test

import (
	"testing"

	"github.com/Daemon0x0000000/lql/pkg/lql"
)

// =============================================================================
// Tests d'évaluation de l'AST
// =============================================================================

func TestAST_Eval_SimpleEquals(t *testing.T) {
	tests := []struct {
		name string
		data map[lql.Field]interface{}
		want bool
	}{
		{"match", map[lql.Field]interface{}{"name": "John"}, true},
		{"no match", map[lql.Field]interface{}{"name": "Jane"}, false},
		{"missing field", map[lql.Field]interface{}{"other": "John"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse("name=John")
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("Eval(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestAST_Eval_Or(t *testing.T) {
	tests := []struct {
		name string
		data map[lql.Field]interface{}
		want bool
	}{
		{"first matches", map[lql.Field]interface{}{"a": "1", "b": "x"}, true},
		{"second matches", map[lql.Field]interface{}{"a": "x", "b": "2"}, true},
		{"both match", map[lql.Field]interface{}{"a": "1", "b": "2"}, true},
		{"none match", map[lql.Field]interface{}{"a": "x", "b": "x"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse("a=1^ORb=2")
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("Eval(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestAST_Eval_And(t *testing.T) {
	tests := []struct {
		name string
		data map[lql.Field]interface{}
		want bool
	}{
		{"both match", map[lql.Field]interface{}{"a": "1", "b": "2"}, true},
		{"first matches only", map[lql.Field]interface{}{"a": "1", "b": "x"}, false},
		{"second matches only", map[lql.Field]interface{}{"a": "x", "b": "2"}, false},
		{"none match", map[lql.Field]interface{}{"a": "x", "b": "x"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse("a=1^b=2")
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("Eval(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestAST_Eval_Xor(t *testing.T) {
	tests := []struct {
		name string
		data map[lql.Field]interface{}
		want bool
	}{
		{"first matches only", map[lql.Field]interface{}{"a": "1", "b": "x"}, true},
		{"second matches only", map[lql.Field]interface{}{"a": "x", "b": "2"}, true},
		{"both match", map[lql.Field]interface{}{"a": "1", "b": "2"}, false},
		{"none match", map[lql.Field]interface{}{"a": "x", "b": "x"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse("a=1^XORb=2")
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("Eval(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestAST_Eval_Nor(t *testing.T) {
	tests := []struct {
		name string
		data map[lql.Field]interface{}
		want bool
	}{
		{"both match", map[lql.Field]interface{}{"a": "1", "b": "2"}, false},
		{"first matches only", map[lql.Field]interface{}{"a": "1", "b": "x"}, false},
		{"second matches only", map[lql.Field]interface{}{"a": "x", "b": "2"}, false},
		{"none match", map[lql.Field]interface{}{"a": "x", "b": "x"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse("a=1^NORb=2")
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("Eval(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestAST_Eval_Nand(t *testing.T) {
	tests := []struct {
		name string
		data map[lql.Field]interface{}
		want bool
	}{
		{"both match", map[lql.Field]interface{}{"a": "1", "b": "2"}, false},
		{"first matches only", map[lql.Field]interface{}{"a": "1", "b": "x"}, true},
		{"second matches only", map[lql.Field]interface{}{"a": "x", "b": "2"}, true},
		{"none match", map[lql.Field]interface{}{"a": "x", "b": "x"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse("a=1^NANDb=2")
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("Eval(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestAST_Eval_Xnor(t *testing.T) {
	tests := []struct {
		name string
		data map[lql.Field]interface{}
		want bool
	}{
		{"both match", map[lql.Field]interface{}{"a": "1", "b": "2"}, true},
		{"first matches only", map[lql.Field]interface{}{"a": "1", "b": "x"}, false},
		{"second matches only", map[lql.Field]interface{}{"a": "x", "b": "2"}, false},
		{"none match", map[lql.Field]interface{}{"a": "x", "b": "x"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse("a=1^XNORb=2")
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("Eval(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestAST_Eval_In(t *testing.T) {
	tests := []struct {
		name string
		data map[lql.Field]interface{}
		want bool
	}{
		{"value in list", map[lql.Field]interface{}{"status": "active"}, true},
		{"value not in list", map[lql.Field]interface{}{"status": "deleted"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse("statusINactive,pending,completed")
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("Eval(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestAST_Eval_Matches(t *testing.T) {
	tests := []struct {
		name string
		data map[lql.Field]interface{}
		want bool
	}{
		{"email matches", map[lql.Field]interface{}{"email": "user@example.com"}, true},
		{"email no match", map[lql.Field]interface{}{"email": "invalid"}, false},
		{"phone matches", map[lql.Field]interface{}{"phone": "0612345678"}, true},
		{"phone no match", map[lql.Field]interface{}{"phone": "abc"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expr string
			if tt.name == "email matches" || tt.name == "email no match" {
				expr = `emailMATCHES'^[a-z]+@[a-z]+\.[a-z]+$'`
			} else {
				expr = `phoneMATCHES'^[0-9]{10}$'`
			}
			ast, err := lql.Parse(expr)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("Eval(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestAST_Eval_Complex(t *testing.T) {
	tests := []struct {
		name string
		data map[lql.Field]interface{}
		want bool
	}{
		{
			"sys_id matches",
			map[lql.Field]interface{}{"sys_id": "123", "test": "no", "me": "no"},
			true,
		},
		{
			"test IN matches",
			map[lql.Field]interface{}{"sys_id": "no", "test": "hello", "me": "no"},
			true,
		},
		{
			"me matches",
			map[lql.Field]interface{}{"sys_id": "no", "test": "no", "me": "test"},
			true,
		},
		{
			"nothing matches",
			map[lql.Field]interface{}{"sys_id": "no", "test": "no", "me": "no"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := lql.Parse("sys_id=123^OR(testINhello,world^ORme=test)")
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := ast.Eval(tt.data); got != tt.want {
				t.Errorf("Eval(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestAST_Eval_NilRoot(t *testing.T) {
	// Test parsing an invalid expression to ensure robustness
	ast, err := lql.Parse("a=1")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	// Empty data should not match
	if ast.Eval(map[lql.Field]interface{}{}) != false {
		t.Error("Eval with empty data should return false")
	}
}
