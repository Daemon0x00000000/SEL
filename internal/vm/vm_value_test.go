package vm

import (
	"testing"
)

// ============================================================================
// Value.Compare Tests
// ============================================================================

func TestValueCompare(t *testing.T) {
	tests := []struct {
		name     string
		v1       Value
		v2       Value
		expected int
	}{
		// Bool
		{"bool: true == true", Value{Type: TYPE_BOOL, Bool: true}, Value{Type: TYPE_BOOL, Bool: true}, 0},
		{"bool: false == false", Value{Type: TYPE_BOOL, Bool: false}, Value{Type: TYPE_BOOL, Bool: false}, 0},
		{"bool: false < true", Value{Type: TYPE_BOOL, Bool: false}, Value{Type: TYPE_BOOL, Bool: true}, -1},
		{"bool: true > false", Value{Type: TYPE_BOOL, Bool: true}, Value{Type: TYPE_BOOL, Bool: false}, 1},

		// Int8
		{"int8: 5 == 5", Value{Type: TYPE_INT8, Int8: 5}, Value{Type: TYPE_INT8, Int8: 5}, 0},
		{"int8: 10 > 5", Value{Type: TYPE_INT8, Int8: 10}, Value{Type: TYPE_INT8, Int8: 5}, 1},
		{"int8: 3 < 7", Value{Type: TYPE_INT8, Int8: 3}, Value{Type: TYPE_INT8, Int8: 7}, -1},
		{"int8: -5 < 5", Value{Type: TYPE_INT8, Int8: -5}, Value{Type: TYPE_INT8, Int8: 5}, -1},

		// Int16
		{"int16: 1000 == 1000", Value{Type: TYPE_INT16, Int16: 1000}, Value{Type: TYPE_INT16, Int16: 1000}, 0},
		{"int16: 2000 > 1000", Value{Type: TYPE_INT16, Int16: 2000}, Value{Type: TYPE_INT16, Int16: 1000}, 1},
		{"int16: 500 < 1000", Value{Type: TYPE_INT16, Int16: 500}, Value{Type: TYPE_INT16, Int16: 1000}, -1},

		// Int32
		{"int32: 100000 == 100000", Value{Type: TYPE_INT32, Int32: 100000}, Value{Type: TYPE_INT32, Int32: 100000}, 0},
		{"int32: 200000 > 100000", Value{Type: TYPE_INT32, Int32: 200000}, Value{Type: TYPE_INT32, Int32: 100000}, 1},
		{"int32: 50000 < 100000", Value{Type: TYPE_INT32, Int32: 50000}, Value{Type: TYPE_INT32, Int32: 100000}, -1},

		// String
		{"string: equal", Value{Type: TYPE_STRING, String: "hello"}, Value{Type: TYPE_STRING, String: "hello"}, 0},
		{"string: world > hello", Value{Type: TYPE_STRING, String: "world"}, Value{Type: TYPE_STRING, String: "hello"}, 1},
		{"string: apple < banana", Value{Type: TYPE_STRING, String: "apple"}, Value{Type: TYPE_STRING, String: "banana"}, -1},
		{"string: empty < hello", Value{Type: TYPE_STRING, String: ""}, Value{Type: TYPE_STRING, String: "hello"}, -1},

		// Array - Equal
		{
			"array: equal",
			Value{Type: TYPE_ARRAY, Array: []Value{{Type: TYPE_INT8, Int8: 1}, {Type: TYPE_INT8, Int8: 2}}},
			Value{Type: TYPE_ARRAY, Array: []Value{{Type: TYPE_INT8, Int8: 1}, {Type: TYPE_INT8, Int8: 2}}},
			0,
		},
		// Array - First element greater
		{
			"array: [2,1] > [1,2]",
			Value{Type: TYPE_ARRAY, Array: []Value{{Type: TYPE_INT8, Int8: 2}, {Type: TYPE_INT8, Int8: 1}}},
			Value{Type: TYPE_ARRAY, Array: []Value{{Type: TYPE_INT8, Int8: 1}, {Type: TYPE_INT8, Int8: 2}}},
			1,
		},
		// Array - First element lesser
		{
			"array: [1,5] < [2,1]",
			Value{Type: TYPE_ARRAY, Array: []Value{{Type: TYPE_INT8, Int8: 1}, {Type: TYPE_INT8, Int8: 5}}},
			Value{Type: TYPE_ARRAY, Array: []Value{{Type: TYPE_INT8, Int8: 2}, {Type: TYPE_INT8, Int8: 1}}},
			-1,
		},
		// Array - Length comparison
		{
			"array: [1] < [1,2]",
			Value{Type: TYPE_ARRAY, Array: []Value{{Type: TYPE_INT8, Int8: 1}}},
			Value{Type: TYPE_ARRAY, Array: []Value{{Type: TYPE_INT8, Int8: 1}, {Type: TYPE_INT8, Int8: 2}}},
			-1,
		},
		{
			"array: [1,2,3] > [1,2]",
			Value{Type: TYPE_ARRAY, Array: []Value{{Type: TYPE_INT8, Int8: 1}, {Type: TYPE_INT8, Int8: 2}, {Type: TYPE_INT8, Int8: 3}}},
			Value{Type: TYPE_ARRAY, Array: []Value{{Type: TYPE_INT8, Int8: 1}, {Type: TYPE_INT8, Int8: 2}}},
			1,
		},
		{
			"array: empty arrays",
			Value{Type: TYPE_ARRAY, Array: []Value{}},
			Value{Type: TYPE_ARRAY, Array: []Value{}},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.v1.Compare(tt.v2)
			if err != nil {
				t.Fatalf("Compare failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestValueCompare_TypeMismatch(t *testing.T) {
	tests := []struct {
		name string
		v1   Value
		v2   Value
	}{
		{"int8 vs string", Value{Type: TYPE_INT8, Int8: 42}, Value{Type: TYPE_STRING, String: "42"}},
		{"bool vs int8", Value{Type: TYPE_BOOL, Bool: true}, Value{Type: TYPE_INT8, Int8: 1}},
		{"array vs string", Value{Type: TYPE_ARRAY, Array: []Value{}}, Value{Type: TYPE_STRING, String: "[]"}},
		{"int8 vs int16", Value{Type: TYPE_INT8, Int8: 42}, Value{Type: TYPE_INT16, Int16: 42}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.v1.Compare(tt.v2)
			if err == nil {
				t.Fatal("Expected error for type mismatch, got nil")
			}
		})
	}
}

func TestValueCompare_NestedArray(t *testing.T) {
	v1 := Value{Type: TYPE_ARRAY, Array: []Value{
		{Type: TYPE_ARRAY, Array: []Value{
			{Type: TYPE_INT8, Int8: 1},
			{Type: TYPE_INT8, Int8: 2},
		}},
	}}

	v2 := Value{Type: TYPE_ARRAY, Array: []Value{
		{Type: TYPE_ARRAY, Array: []Value{
			{Type: TYPE_INT8, Int8: 1},
			{Type: TYPE_INT8, Int8: 2},
		}},
	}}

	result, err := v1.Compare(v2)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if result != 0 {
		t.Errorf("Expected 0 for equal nested arrays, got %d", result)
	}
}

func TestValueCompare_ArrayMixedElements(t *testing.T) {
	v1 := Value{Type: TYPE_ARRAY, Array: []Value{
		{Type: TYPE_INT8, Int8: 1},
		{Type: TYPE_INT8, Int8: 2},
		{Type: TYPE_INT8, Int8: 3},
	}}

	v2 := Value{Type: TYPE_ARRAY, Array: []Value{
		{Type: TYPE_INT8, Int8: 1},
		{Type: TYPE_INT8, Int8: 2},
		{Type: TYPE_INT8, Int8: 4},
	}}

	result, err := v1.Compare(v2)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if result != -1 {
		t.Errorf("Expected -1 (3 < 4), got %d", result)
	}
}

// ============================================================================
// Utility Functions Tests
// ============================================================================

func TestCmpGeneric(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{5, 10, -1},
		{10, 5, 1},
		{5, 5, 0},
		{-5, 5, -1},
		{-10, -5, -1},
	}

	for _, tt := range tests {
		if got := cmpGeneric(tt.a, tt.b); got != tt.expected {
			t.Errorf("cmpGeneric(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
		}
	}
}

// ============================================================================
// Type Inference Tests
// ============================================================================

func TestDetermineIntType(t *testing.T) {
	tests := []struct {
		input        int
		expectedType Type
	}{
		{-128, TYPE_INT8},
		{127, TYPE_INT8},
		{0, TYPE_INT8},
		{42, TYPE_INT8},
		{-32768, TYPE_INT16},
		{32767, TYPE_INT16},
		{200, TYPE_INT16},
		{-200, TYPE_INT16},
		{100000, TYPE_INT32},
		{-100000, TYPE_INT32},
		{2147483647, TYPE_INT32},
	}

	for _, tt := range tests {
		typ, _ := determineIntType(tt.input)
		if typ != tt.expectedType {
			t.Errorf("determineIntType(%d) = %v, want %v", tt.input, typ, tt.expectedType)
		}
	}
}

func TestInferRuntimeValue(t *testing.T) {
	tests := []struct {
		name     string
		bytecode []byte
		check    func(*testing.T, Value)
	}{
		{
			name:     "int8",
			bytecode: []byte{byte(TYPE_INT8), 1, 42},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_INT8 || v.Int8 != 42 {
					t.Errorf("Expected int8(42), got %+v", v)
				}
			},
		},
		{
			name:     "int16",
			bytecode: []byte{byte(TYPE_INT16), 2, 0x03, 0xE8}, // 1000
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_INT16 || v.Int16 != 1000 {
					t.Errorf("Expected int16(1000), got %+v", v)
				}
			},
		},
		{
			name:     "int32",
			bytecode: []byte{byte(TYPE_INT32), 4, 0x00, 0x01, 0x86, 0xA0}, // 100000
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_INT32 || v.Int32 != 100000 {
					t.Errorf("Expected int32(100000), got %+v", v)
				}
			},
		},
		{
			name:     "string",
			bytecode: []byte{byte(TYPE_STRING), 5, 'h', 'e', 'l', 'l', 'o'},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_STRING || v.String != "hello" {
					t.Errorf("Expected string(hello), got %+v", v)
				}
			},
		},
		{
			name:     "bool true",
			bytecode: []byte{byte(TYPE_BOOL), 1, 1},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_BOOL || v.Bool != true {
					t.Errorf("Expected bool(true), got %+v", v)
				}
			},
		},
		{
			name:     "bool false",
			bytecode: []byte{byte(TYPE_BOOL), 1, 0},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_BOOL || v.Bool != false {
					t.Errorf("Expected bool(false), got %+v", v)
				}
			},
		},
		{
			name: "array",
			bytecode: []byte{
				byte(TYPE_ARRAY), 2,
				byte(TYPE_INT8), 1, 1,
				byte(TYPE_INT8), 1, 2,
			},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_ARRAY || len(v.Array) != 2 {
					t.Errorf("Expected array[2], got %+v", v)
				}
				if v.Array[0].Int8 != 1 || v.Array[1].Int8 != 2 {
					t.Errorf("Expected [1,2], got %+v", v.Array)
				}
			},
		},
		{
			name: "nested array",
			bytecode: []byte{
				byte(TYPE_ARRAY), 2,
				byte(TYPE_ARRAY), 2,
				byte(TYPE_INT8), 1, 1,
				byte(TYPE_INT8), 1, 2,
				byte(TYPE_ARRAY), 2,
				byte(TYPE_INT8), 1, 3,
				byte(TYPE_INT8), 1, 4,
			},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_ARRAY || len(v.Array) != 2 {
					t.Fatalf("Expected array[2], got %+v", v)
				}
				if v.Array[0].Type != TYPE_ARRAY || len(v.Array[0].Array) != 2 {
					t.Errorf("Expected nested array[2], got %+v", v.Array[0])
				}
				if v.Array[0].Array[0].Int8 != 1 {
					t.Errorf("Expected [0][0]=1, got %d", v.Array[0].Array[0].Int8)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewVM(tt.bytecode, nil)
			val, err := vm.inferRuntimeValue()
			if err != nil {
				t.Fatalf("inferRuntimeValue failed: %v", err)
			}
			tt.check(t, val)
		})
	}
}

// ============================================================================
// Conversion Tests
// ============================================================================

func TestConvertInterfaceToValue(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		check func(*testing.T, Value)
	}{
		{
			name:  "string",
			input: "hello",
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_STRING || v.String != "hello" {
					t.Errorf("Expected string(hello), got %+v", v)
				}
			},
		},
		{
			name:  "int",
			input: 42,
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_INT8 || v.Int8 != 42 {
					t.Errorf("Expected int8(42), got %+v", v)
				}
			},
		},
		{
			name:  "bool",
			input: true,
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_BOOL || v.Bool != true {
					t.Errorf("Expected bool(true), got %+v", v)
				}
			},
		},
		{
			name:  "array",
			input: []interface{}{1, 2, 3},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_ARRAY || len(v.Array) != 3 {
					t.Errorf("Expected array[3], got %+v", v)
				}
			},
		},
		{
			name:  "nested array",
			input: []interface{}{[]interface{}{1, 2}, []interface{}{3, 4}},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_ARRAY || len(v.Array) != 2 {
					t.Errorf("Expected array[2], got %+v", v)
				}
				if v.Array[0].Type != TYPE_ARRAY {
					t.Errorf("Expected nested TYPE_ARRAY, got %v", v.Array[0].Type)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewVM([]byte{}, nil)
			val, err := vm.convertInterfaceToValue(tt.input)
			if err != nil {
				t.Fatalf("convertInterfaceToValue failed: %v", err)
			}
			tt.check(t, val)
		})
	}
}
