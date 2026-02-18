package ast

import (
	"github.com/Daemon0x00000000/sel/internal/vm"
)

const (
	AND LogicalOperator = "^"
	OR  LogicalOperator = "^OR"
	XOR LogicalOperator = "^XOR"
)

const (
	EQUALS                ComparisonOperator = "="
	GREATER_THAN          ComparisonOperator = ">"
	LESS_THAN             ComparisonOperator = "<"
	GREATER_THAN_OR_EQUAL ComparisonOperator = ">="
	LESS_THAN_OR_EQUAL    ComparisonOperator = "<="
	STARTS_WITH           ComparisonOperator = "STARTSWITH"
	ENDS_WITH             ComparisonOperator = "ENDSWITH"
	IN                    ComparisonOperator = "IN"
	CONTAINS              ComparisonOperator = "CONTAINS"
	//MATCHES               ComparisonOperator = "MATCHES"
)

// !=, !IN, !CONTAINS, !MATCHES

var logicalOperatorsOrdered = []LogicalOperator{
	XOR,
	OR,
	AND,
}

var comparisonOpsOrdered = []ComparisonOperator{
	STARTS_WITH, // "STARTSWITH" = 10 chars
	ENDS_WITH,   // "ENDSWITH" = 8 chars
	CONTAINS,    // "CONTAINS" = 8 chars
	//MATCHES,               // "MATCHES" = 7 chars
	IN,                    // "IN" = 2 chars
	GREATER_THAN_OR_EQUAL, // ">=" = 2 chars
	LESS_THAN_OR_EQUAL,    // "<=" = 2 chars
	GREATER_THAN,          // ">" = 1 char
	LESS_THAN,             // "<" = 1 char
	EQUALS,                // "=" = 1 char
}

var logicalOperators = LogicalOperatorMapping{
	AND: vm.OP_AND,
	OR:  vm.OP_OR,
	XOR: vm.OP_XOR,
}

var comparisonOperators = ComparisonOperatorMapping{
	EQUALS:                vm.OP_EQ,
	IN:                    vm.OP_IN,
	CONTAINS:              vm.OP_CONTAINS,
	STARTS_WITH:           vm.OP_STARTSWITH,
	ENDS_WITH:             vm.OP_ENDSWITH,
	GREATER_THAN:          vm.OP_GT,
	LESS_THAN:             vm.OP_LT,
	GREATER_THAN_OR_EQUAL: vm.OP_GTE,
	LESS_THAN_OR_EQUAL:    vm.OP_LTE,
}
