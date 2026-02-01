package lql

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	AND LogicalOperator = "^"
	OR  LogicalOperator = "^OR"
	NOT LogicalOperator = "!"
	XOR LogicalOperator = "^XOR"
)

const (
	EQUALS                ComparisonOperator = "="
	NOT_EQUALS            ComparisonOperator = "!="
	GREATER_THAN          ComparisonOperator = ">"
	LESS_THAN             ComparisonOperator = "<"
	GREATER_THAN_OR_EQUAL ComparisonOperator = ">="
	LESS_THAN_OR_EQUAL    ComparisonOperator = "<="
	STARTS_WITH           ComparisonOperator = "STARTSWITH"
	ENDS_WITH             ComparisonOperator = "ENDSWITH"
	IN                    ComparisonOperator = "IN"
	NOT_IN                ComparisonOperator = "NOTIN"
	CONTAINS              ComparisonOperator = "CONTAINS"
	DOES_NOT_CONTAIN      ComparisonOperator = "DOESNOTCONTAIN"
	MATCHES               ComparisonOperator = "MATCHES"
)

var logicalOperatorsOrdered = []LogicalOperator{
	XOR,
	OR,
	AND,
}

var comparisonOpsOrdered = []ComparisonOperator{
	DOES_NOT_CONTAIN,      // "DOESNOTCONTAIN" = 14 chars
	GREATER_THAN_OR_EQUAL, // ">=" = 2 chars
	LESS_THAN_OR_EQUAL,    // "<=" = 2 chars
	STARTS_WITH,           // "STARTSWITH" = 10 chars
	ENDS_WITH,             // "ENDSWITH" = 8 chars
	CONTAINS,              // "CONTAINS" = 8 chars
	MATCHES,               // "MATCHES" = 7 chars
	NOT_IN,                // "NOTIN" = 5 chars
	IN,                    // "IN" = 2 chars
	NOT_EQUALS,            // "!=" = 2 chars
	GREATER_THAN,          // ">" = 1 char
	LESS_THAN,             // "<" = 1 char
	EQUALS,                // "=" = 1 char
}

var logicalOperators = LogicalOperatorMapping{
	AND: func(results []bool) bool {
		for _, result := range results {
			if !result {
				return false
			}
		}
		return true
	},
	OR: func(results []bool) bool {
		for _, result := range results {
			if result {
				return true
			}
		}
		return false
	},
	XOR: func(results []bool) bool {
		count := 0
		for _, result := range results {
			if result {
				count++
			}
		}
		return count == 1
	},
}

var comparisonOperators = ComparisonOperatorMapping{
	EQUALS: func(expectedValues ...string) PredicateClosure {
		expected := expectedValues[0]
		return func(value interface{}) bool {
			return fmt.Sprint(value) == expected
		}
	},

	NOT_EQUALS: func(expectedValues ...string) PredicateClosure {
		expected := expectedValues[0]
		return func(value interface{}) bool {
			return fmt.Sprint(value) != expected
		}
	},

	IN: func(expectedValues ...string) PredicateClosure {
		set := make(map[string]bool)
		for _, v := range expectedValues {
			set[v] = true
		}
		return func(value interface{}) bool {
			return set[fmt.Sprint(value)]
		}
	},

	NOT_IN: func(expectedValues ...string) PredicateClosure {
		set := make(map[string]bool)
		for _, v := range expectedValues {
			set[v] = true
		}
		return func(value interface{}) bool {
			return !set[fmt.Sprint(value)]
		}
	},

	CONTAINS: func(expectedValues ...string) PredicateClosure {
		search := expectedValues[0]
		return func(value interface{}) bool {
			return strings.Contains(fmt.Sprint(value), search)
		}
	},

	DOES_NOT_CONTAIN: func(expectedValues ...string) PredicateClosure {
		search := expectedValues[0]
		return func(value interface{}) bool {
			return !strings.Contains(fmt.Sprint(value), search)
		}
	},

	STARTS_WITH: func(expectedValues ...string) PredicateClosure {
		prefix := expectedValues[0]
		return func(value interface{}) bool {
			return strings.HasPrefix(fmt.Sprint(value), prefix)
		}
	},

	ENDS_WITH: func(expectedValues ...string) PredicateClosure {
		suffix := expectedValues[0]
		return func(value interface{}) bool {
			return strings.HasSuffix(fmt.Sprint(value), suffix)
		}
	},

	GREATER_THAN: func(expectedValues ...string) PredicateClosure {
		expected := expectedValues[0]
		return func(value interface{}) bool {
			return fmt.Sprint(value) > expected
		}
	},

	LESS_THAN: func(expectedValues ...string) PredicateClosure {
		expected := expectedValues[0]
		return func(value interface{}) bool {
			return fmt.Sprint(value) < expected
		}
	},

	GREATER_THAN_OR_EQUAL: func(expectedValues ...string) PredicateClosure {
		expected := expectedValues[0]
		return func(value interface{}) bool {
			return fmt.Sprint(value) >= expected
		}
	},

	LESS_THAN_OR_EQUAL: func(expectedValues ...string) PredicateClosure {
		expected := expectedValues[0]
		return func(value interface{}) bool {
			return fmt.Sprint(value) <= expected
		}
	},

	MATCHES: func(expectedValues ...string) PredicateClosure {
		pattern := expectedValues[0]
		return func(value interface{}) bool {
			return regexp.MustCompile(pattern).MatchString(fmt.Sprint(value))
		}
	},
}
