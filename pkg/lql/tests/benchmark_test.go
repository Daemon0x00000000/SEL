package lql_test

import (
	"testing"

	"github.com/Daemon0x0000000/lql/pkg/lql"
)

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkParse_Simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lql.Parse("field=value")
	}
}

func BenchmarkParse_WithOr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lql.Parse("a=1^ORb=2")
	}
}

func BenchmarkParse_WithAnd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lql.Parse("a=1^b=2")
	}
}

func BenchmarkParse_Complex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lql.Parse("sys_id=123^OR(testINhello,world^ORme=test)")
	}
}

func BenchmarkParse_VeryComplex(b *testing.B) {
	expr := "a=1^ORb=2^ORc=3^ORd=4^ORe=5^ORf=6^ORg=7^ORh=8^ORi=9^ORj=10"
	for i := 0; i < b.N; i++ {
		lql.Parse(expr)
	}
}

func BenchmarkParse_DeepNested(b *testing.B) {
	expr := "(((a=1^ORb=2)^ORc=3)^ORd=4)"
	for i := 0; i < b.N; i++ {
		lql.Parse(expr)
	}
}

func BenchmarkEval_Simple(b *testing.B) {
	data := map[lql.Field]interface{}{"field": "value"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ast, _ := lql.Parse("field=value")
		_ = ast.Eval(data)
	}
}

func BenchmarkEval_Complex(b *testing.B) {
	data := map[lql.Field]interface{}{
		"sys_id": "123",
		"test":   "hello",
		"me":     "test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ast, _ := lql.Parse("sys_id=123^OR(testINhello,world^ORme=test)")
		_ = ast.Eval(data)
	}
}

func BenchmarkEval_Complex_EvalOnly(b *testing.B) {
	ast, _ := lql.Parse("sys_id=123^OR(testINhello,world^ORme=test)")
	data := map[lql.Field]interface{}{
		"sys_id": "123",
		"test":   "hello",
		"me":     "test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ast.Eval(data)
	}
}

func BenchmarkEval_ManyFields(b *testing.B) {
	data := map[lql.Field]interface{}{
		"a": "1", "b": "2", "c": "3", "d": "4", "e": "5",
		"f": "6", "g": "7", "h": "8", "i": "9", "j": "10",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ast, _ := lql.Parse("a=1^ORb=2^ORc=3^ORd=4^ORe=5^ORf=6^ORg=7^ORh=8^ORi=9^ORj=10")
		_ = ast.Eval(data)
	}
}

func BenchmarkEval_ManyFields_EvalOnly(b *testing.B) {
	ast, _ := lql.Parse("a=1^ORb=2^ORc=3^ORd=4^ORe=5^ORf=6^ORg=7^ORh=8^ORi=9^ORj=10")
	data := map[lql.Field]interface{}{
		"a": "1", "b": "2", "c": "3", "d": "4", "e": "5",
		"f": "6", "g": "7", "h": "8", "i": "9", "j": "10",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ast.Eval(data)
	}
}
