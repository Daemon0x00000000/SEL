package vm

import (
	"testing"
)

// ============================================================================
// Serialization Tests
// ============================================================================

func TestSerializationFunctions(t *testing.T) {
	t.Run("SerializeLoadGlobal", func(t *testing.T) {
		result := SerializeLoadGlobal("myvar")
		expected := []byte{byte(LOAD_GLOBAL), 5, 'm', 'y', 'v', 'a', 'r'}
		assertBytecodeEqual(t, result, expected)
	})

	t.Run("SerializePush", func(t *testing.T) {
		tests := []struct {
			name     string
			input    interface{}
			expected []byte
		}{
			{
				name:     "int8",
				input:    42,
				expected: []byte{byte(PUSH), byte(TYPE_INT8), 1, 42},
			},
			{
				name:     "string",
				input:    "hi",
				expected: []byte{byte(PUSH), byte(TYPE_STRING), 2, 'h', 'i'},
			},
			{
				name:     "bool",
				input:    true,
				expected: []byte{byte(PUSH), byte(TYPE_BOOL), 1, 1},
			},
			{
				name:  "array",
				input: []interface{}{1, 2},
				expected: []byte{
					byte(PUSH), byte(TYPE_ARRAY), 2,
					byte(TYPE_INT8), 1, 1,
					byte(TYPE_INT8), 1, 2,
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := SerializePush(tt.input)
				if err != nil {
					t.Fatalf("SerializePush failed: %v", err)
				}
				assertBytecodeEqual(t, result, tt.expected)
			})
		}
	})

	t.Run("SerializeOperator", func(t *testing.T) {
		tests := []struct {
			op       OpCode
			expected []byte
		}{
			{OP_EQ, []byte{byte(OP_EQ)}},
			{OP_GT, []byte{byte(OP_GT)}},
			{OP_AND, []byte{byte(OP_AND)}},
			{OP_NOT, []byte{byte(OP_NOT)}},
		}

		for _, tt := range tests {
			result := SerializeOperator(tt.op)
			assertBytecodeEqual(t, result, tt.expected)
		}
	})
}

func TestInferType(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected Type
	}{
		{42, TYPE_INT8},
		{1000, TYPE_INT16},
		{100000, TYPE_INT32},
		{"hello", TYPE_STRING},
		{true, TYPE_BOOL},
		{[]interface{}{1, 2}, TYPE_ARRAY},
	}

	for _, tt := range tests {
		typ, err := inferType(tt.input)
		if err != nil {
			t.Fatalf("inferType failed: %v", err)
		}
		if typ != tt.expected {
			t.Errorf("inferType(%v) = %v, want %v", tt.input, typ, tt.expected)
		}
	}
}

func TestSerializeValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []byte
	}{
		{
			name:     "int8",
			input:    42,
			expected: []byte{byte(TYPE_INT8), 1, 42},
		},
		{
			name:     "string",
			input:    "ab",
			expected: []byte{byte(TYPE_STRING), 2, 'a', 'b'},
		},
		{
			name:     "bool",
			input:    false,
			expected: []byte{byte(TYPE_BOOL), 1, 0},
		},
		{
			name:  "array",
			input: []interface{}{1, 2},
			expected: []byte{
				byte(TYPE_ARRAY), 2,
				byte(TYPE_INT8), 1, 1,
				byte(TYPE_INT8), 1, 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := serializeValue(tt.input)
			if err != nil {
				t.Fatalf("serializeValue failed: %v", err)
			}
			assertBytecodeEqual(t, result, tt.expected)
		})
	}
}
