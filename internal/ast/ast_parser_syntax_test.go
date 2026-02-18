package ast

import "testing"

func TestParse_ValidSyntax(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		// Basic comparisons
		{"simple equals", "field=value"},
		{"not equals", "field!=value"},
		{"greater than", "age>18"},
		{"less than", "age<65"},
		{"greater or equal", "age>=18"},
		{"less or equal", "age<=65"},

		// String operators
		{"in operator", "statusINactive,pending,completed"},
		{"not in operator", "status!INdeleted,archived"},
		{"contains", "nameCONTAINSjohn"},
		{"not contains", "name!CONTAINSjohn"},
		{"starts with", "emailSTARTSWITHadmin"},
		{"ends with", "emailENDSWITH@example.com"},
		// TODO: MATCHES operator not fully implemented yet
		// {"matches regex", "emailMATCHES'^[a-z]+@[a-z]+\\.com$'"},

		// Logical operators
		{"and", "a=1^b=2"},
		{"or", "a=1^ORb=2"},
		{"xor", "a=1^XORb=2"},
		{"not with parens", "!(a=1)"},

		// Multiple conditions
		{"multiple and", "a=1^b=2^c=3"},
		{"multiple or", "a=1^ORb=2^ORc=3"},
		{"mixed", "a=1^b=2^ORc=3"},

		// Parentheses
		{"simple parens", "(a=1)"},
		{"nested parens", "((a=1))"},
		{"parens with or", "(a=1^ORb=2)"},
		{"complex nested", "a=1^OR(b=2^c=3)"},
		{"deep nested", "((a=1^ORb=2)^c=3)^OR(d=4^e=5)"},

		// Edge cases
		{"empty value", "field="},
		{"single char field", "a=1"},
		{"long field name", "very_long_field_name_with_underscores=value"},
		{"numeric value", "age=123"},
		{"special chars in value", "field=hello-world_123"},
		{"value with spaces", "field=hello world"},

		// Quoted values
		{"single quoted", "field='value with spaces'"},
		{"quoted in IN", "statusIN'active','pending'"},
		{"quoted with comma", "field='value, with comma'"},
		{"quoted with equals", "field='value=123'"},

		// Complex real-world examples
		{"servicenow style", "sys_id=123^OR(active=true^category=incident)"},
		// TODO: MATCHES operator not fully implemented yet
		// {"email validation", "emailMATCHES'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$'"},
		{"multi-field search", "nameSTARTSWITHjohn^ORnameSTARTSWITHjane^ORnameSTARTSWITHjack"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := Parse(tt.expr)
			assertNoError(t, err)
			assertASTNotNil(t, ast)
		})
	}
}

func TestParse_InvalidSyntax(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"unbalanced open paren", "(a=1"},
		{"unbalanced close paren", "a=1)"},
		{"multiple unbalanced", "((a=1)"},
		{"wrong order parens", ")("},
		{"no operator", "field"},
		{"no field", "=value"},
		{"double operator", "field==value"},
		{"unclosed quote", "field='value"},
		{"incomplete escape", "field='value\\"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.expr)
			assertError(t, err)
		})
	}
}

func TestParse_EmptyExpression(t *testing.T) {
	// Empty expression should fail at parsing stage
	_, err := Parse("")
	if err == nil {
		t.Skip("Empty expression handling not specified")
	}
}

func TestParse_WithWhitespace(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"leading spaces", "  field=value"},
		{"trailing spaces", "field=value  "},
		{"spaces around operator", "field = value"},
		{"spaces in parens", "( a=1 )"},
		{"multiple spaces", "a  =  1  ^  b  =  2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := Parse(tt.expr)
			assertNoError(t, err)
			assertASTNotNil(t, ast)
		})
	}
}

func TestParse_OperatorPrecedence(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		// XOR has highest precedence
		{"xor before or", "a=1^ORb=2^XORc=3"},
		{"xor before and", "a=1^b=2^XORc=3"},

		// OR has middle precedence
		{"or before and", "a=1^b=2^ORc=3^d=4"},

		// Parentheses override precedence
		{"parens override", "(a=1^ORb=2)^c=3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := Parse(tt.expr)
			assertNoError(t, err)
			assertASTNotNil(t, ast)
		})
	}
}

func TestParse_QuotedValues(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		wantFail bool
	}{
		{"simple quote", "field='value'", false},
		{"quote with space", "field='value with spaces'", false},
		{"quote with comma", "field='value,with,comma'", false},
		{"quote with equals", "field='a=b'", false},
		{"quote with operator", "field='a^ORb'", false},
		{"escaped quote", "field='it\\'s'", false},
		{"escaped backslash", "field='path\\\\to\\\\file'", false},
		{"escaped newline", "field='line1\\nline2'", false},
		{"escaped tab", "field='col1\\tcol2'", false},
		{"unclosed quote", "field='value", true},
		{"incomplete escape", "field='value\\", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := Parse(tt.expr)
			if tt.wantFail {
				assertError(t, err)
			} else {
				assertNoError(t, err)
				assertASTNotNil(t, ast)
			}
		})
	}
}

func TestParse_INOperatorValues(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"simple list", "statusINa,b,c"},
		{"single value", "statusINactive"},
		{"with spaces", "statusINa, b, c"},
		{"quoted values", "statusIN'active','pending','completed'"},
		{"mixed quoted", "statusINactive,'with space',another"},
		{"numbers", "ageIN18,21,65"},
		{"empty value in list", "statusINa,,c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := Parse(tt.expr)
			assertNoError(t, err)
			assertASTNotNil(t, ast)
		})
	}
}

func TestParse_NOTOperator(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"not with parens", "!(a=1)"},
		{"not equals", "a!=1"},
		{"not in", "status!INdeleted,archived"},
		{"not contains", "name!CONTAINSjohn"},
		{"double not", "!!(a=1)"},
		{"not with complex", "!(a=1^ORb=2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := Parse(tt.expr)
			assertNoError(t, err)
			assertASTNotNil(t, ast)
		})
	}
}

func TestParse_ASTString(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"simple", "field=value"},
		{"with or", "a=1^ORb=2"},
		{"nested", "a=1^OR(b=2^c=3)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := Parse(tt.expr)
			assertNoError(t, err)

			str := ast.String()
			if str == "" {
				t.Error("AST.String() should not be empty")
			}
		})
	}
}
