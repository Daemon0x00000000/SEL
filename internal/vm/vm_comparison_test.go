package vm

import (
	"testing"
)

// ============================================================================
// Comparison Operators Tests
// ============================================================================

func TestComparisonHandlers(t *testing.T) {
	tests := []struct {
		name     string
		handler  func(*VM) error
		left     Value
		right    Value
		expected bool
	}{
		// OP_EQ
		{"eq: int8 equal", (*VM).eqHandler, Value{Type: TYPE_INT8, Int8: 42}, Value{Type: TYPE_INT8, Int8: 42}, true},
		{"eq: int8 not equal", (*VM).eqHandler, Value{Type: TYPE_INT8, Int8: 42}, Value{Type: TYPE_INT8, Int8: 43}, false},
		{"eq: string equal", (*VM).eqHandler, Value{Type: TYPE_STRING, String: "hello"}, Value{Type: TYPE_STRING, String: "hello"}, true},
		{"eq: string not equal", (*VM).eqHandler, Value{Type: TYPE_STRING, String: "hello"}, Value{Type: TYPE_STRING, String: "world"}, false},
		{"eq: bool equal", (*VM).eqHandler, Value{Type: TYPE_BOOL, Bool: true}, Value{Type: TYPE_BOOL, Bool: true}, true},

		// OP_GT
		{"gt: 10 > 5", (*VM).gtHandler, Value{Type: TYPE_INT8, Int8: 10}, Value{Type: TYPE_INT8, Int8: 5}, true},
		{"gt: 5 > 10", (*VM).gtHandler, Value{Type: TYPE_INT8, Int8: 5}, Value{Type: TYPE_INT8, Int8: 10}, false},
		{"gt: 5 > 5", (*VM).gtHandler, Value{Type: TYPE_INT8, Int8: 5}, Value{Type: TYPE_INT8, Int8: 5}, false},

		// OP_LT
		{"lt: 5 < 10", (*VM).ltHandler, Value{Type: TYPE_INT8, Int8: 5}, Value{Type: TYPE_INT8, Int8: 10}, true},
		{"lt: 10 < 5", (*VM).ltHandler, Value{Type: TYPE_INT8, Int8: 10}, Value{Type: TYPE_INT8, Int8: 5}, false},

		// OP_GTE
		{"gte: 10 >= 5", (*VM).gteHandler, Value{Type: TYPE_INT8, Int8: 10}, Value{Type: TYPE_INT8, Int8: 5}, true},
		{"gte: 5 >= 5", (*VM).gteHandler, Value{Type: TYPE_INT8, Int8: 5}, Value{Type: TYPE_INT8, Int8: 5}, true},
		{"gte: 5 >= 10", (*VM).gteHandler, Value{Type: TYPE_INT8, Int8: 5}, Value{Type: TYPE_INT8, Int8: 10}, false},

		// OP_LTE
		{"lte: 5 <= 10", (*VM).lteHandler, Value{Type: TYPE_INT8, Int8: 5}, Value{Type: TYPE_INT8, Int8: 10}, true},
		{"lte: 5 <= 5", (*VM).lteHandler, Value{Type: TYPE_INT8, Int8: 5}, Value{Type: TYPE_INT8, Int8: 5}, true},
		{"lte: 10 <= 5", (*VM).lteHandler, Value{Type: TYPE_INT8, Int8: 10}, Value{Type: TYPE_INT8, Int8: 5}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testBinaryHandler(t, tt.handler, tt.left, tt.right, TYPE_BOOL, func(v Value) bool {
				return v.Bool == tt.expected
			})
		})
	}
}
