package vm

import (
	"testing"
)

// ============================================================================
// String Operations Tests
// ============================================================================

func TestStringOperations(t *testing.T) {
	tests := []struct {
		name     string
		handler  func(*VM) error
		str      string
		pattern  string
		expected bool
	}{
		// STARTSWITH
		{"startswith: hello", (*VM).startsWithHandler, "hello world", "hello", true},
		{"startswith: world", (*VM).startsWithHandler, "hello world", "world", false},
		{"startswith: empty", (*VM).startsWithHandler, "hello", "", true},

		// ENDSWITH
		{"endswith: world", (*VM).endsWithHandler, "hello world", "world", true},
		{"endswith: hello", (*VM).endsWithHandler, "hello world", "hello", false},

		// CONTAINS
		{"contains: world", (*VM).containsHandler, "hello world", "world", true},
		{"contains: lo wo", (*VM).containsHandler, "hello world", "lo wo", true},
		{"contains: xyz", (*VM).containsHandler, "hello world", "xyz", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left := Value{Type: TYPE_STRING, String: tt.str}
			right := Value{Type: TYPE_STRING, String: tt.pattern}
			testBinaryHandler(t, tt.handler, left, right, TYPE_BOOL, func(v Value) bool {
				return v.Bool == tt.expected
			})
		})
	}
}

// ============================================================================
// IN Operator Tests
// ============================================================================

func TestInHandler(t *testing.T) {
	tests := []struct {
		name     string
		needle   Value
		haystack Value
		expected bool
	}{
		{
			name:   "scalar in array - found",
			needle: Value{Type: TYPE_INT8, Int8: 2},
			haystack: Value{Type: TYPE_ARRAY, Array: []Value{
				{Type: TYPE_INT8, Int8: 1},
				{Type: TYPE_INT8, Int8: 2},
				{Type: TYPE_INT8, Int8: 3},
			}},
			expected: true,
		},
		{
			name:   "scalar in array - not found",
			needle: Value{Type: TYPE_INT8, Int8: 5},
			haystack: Value{Type: TYPE_ARRAY, Array: []Value{
				{Type: TYPE_INT8, Int8: 1},
				{Type: TYPE_INT8, Int8: 2},
			}},
			expected: false,
		},
		{
			name: "array in array - any match",
			needle: Value{Type: TYPE_ARRAY, Array: []Value{
				{Type: TYPE_INT8, Int8: 2},
				{Type: TYPE_INT8, Int8: 3},
			}},
			haystack: Value{Type: TYPE_ARRAY, Array: []Value{
				{Type: TYPE_INT8, Int8: 1},
				{Type: TYPE_INT8, Int8: 2},
				{Type: TYPE_INT8, Int8: 4},
			}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testBinaryHandler(t, (*VM).inHandler, tt.needle, tt.haystack, TYPE_BOOL, func(v Value) bool {
				return v.Bool == tt.expected
			})
		})
	}
}
