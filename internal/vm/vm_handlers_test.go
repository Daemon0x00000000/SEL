package vm

import (
	"testing"
)

// ============================================================================
// PUSH Handler Tests
// ============================================================================

func TestPushHandler(t *testing.T) {
	tests := []struct {
		name     string
		bytecode []byte
		check    func(*testing.T, Value)
	}{
		{
			name:     "push int8",
			bytecode: []byte{byte(PUSH), byte(TYPE_INT8), 1, 42},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_INT8 || v.Int8 != 42 {
					t.Errorf("Expected int8(42), got %+v", v)
				}
			},
		},
		{
			name:     "push string",
			bytecode: []byte{byte(PUSH), byte(TYPE_STRING), 5, 'h', 'e', 'l', 'l', 'o'},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_STRING || v.String != "hello" {
					t.Errorf("Expected string(hello), got %+v", v)
				}
			},
		},
		{
			name:     "push bool true",
			bytecode: []byte{byte(PUSH), byte(TYPE_BOOL), 1, 1},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_BOOL || v.Bool != true {
					t.Errorf("Expected bool(true), got %+v", v)
				}
			},
		},
		{
			name:     "push bool false",
			bytecode: []byte{byte(PUSH), byte(TYPE_BOOL), 1, 0},
			check: func(t *testing.T, v Value) {
				if v.Type != TYPE_BOOL || v.Bool != false {
					t.Errorf("Expected bool(false), got %+v", v)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewVM(tt.bytecode, nil)
			if err := vm.Execute(); err != nil {
				t.Fatalf("Execute failed: %v", err)
			}
			if len(vm.dataStack) != 1 {
				t.Fatalf("Expected 1 value on stack, got %d", len(vm.dataStack))
			}
			tt.check(t, vm.dataStack[0])
		})
	}
}

// ============================================================================
// POP Handler Tests
// ============================================================================

func TestPopHandler(t *testing.T) {
	bytecode := []byte{
		byte(PUSH), byte(TYPE_INT8), 1, 42,
		byte(POP),
	}
	vm := NewVM(bytecode, nil)
	if err := vm.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if len(vm.dataStack) != 0 {
		t.Errorf("Expected empty stack after POP, got %d values", len(vm.dataStack))
	}
}

func TestPopHandler_Underflow(t *testing.T) {
	vm := NewVM([]byte{byte(POP)}, nil)
	if err := vm.Execute(); err == nil {
		t.Fatal("Expected error for stack underflow, got nil")
	}
}

// ============================================================================
// STORE_GLOBAL & LOAD_GLOBAL Tests
// ============================================================================

func TestStoreGlobalHandler(t *testing.T) {
	bytecode := []byte{
		byte(STORE_GLOBAL), 5, 'm', 'y', 'v', 'a', 'r',
		byte(TYPE_STRING), 5, 'h', 'e', 'l', 'l', 'o',
	}
	vm := NewVM(bytecode, nil)
	if err := vm.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	val, exists := vm.globals["myvar"]
	if !exists {
		t.Fatal("Expected 'myvar' to be in globals")
	}
	if val.Type != TYPE_STRING || val.String != "hello" {
		t.Errorf("Expected string(hello), got %+v", val)
	}
}

func TestLoadGlobalHandler(t *testing.T) {
	vm := NewVM([]byte{}, nil)
	vm.globals["myvar"] = Value{Type: TYPE_INT8, Int8: 42}

	bytecode := []byte{byte(LOAD_GLOBAL), 5, 'm', 'y', 'v', 'a', 'r'}
	vm.bytecode = bytecode

	if err := vm.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	assertStackValue(t, vm, TYPE_INT8, func(v Value) bool {
		return v.Int8 == 42
	})
}

func TestLoadGlobalHandler_Undefined(t *testing.T) {
	bytecode := []byte{
		byte(LOAD_GLOBAL), 9, 'u', 'n', 'd', 'e', 'f', 'i', 'n', 'e', 'd',
	}
	vm := NewVM(bytecode, nil)
	if err := vm.Execute(); err == nil {
		t.Fatal("Expected error for undefined variable, got nil")
	}
}

// ============================================================================
// CALL_NATIVE Handler Tests
// ============================================================================

func TestCallNativeHandler(t *testing.T) {
	addFunc := func(args []Value) (Value, error) {
		return Value{Type: TYPE_INT8, Int8: args[0].Int8 + args[1].Int8}, nil
	}

	bytecode := []byte{
		byte(PUSH), byte(TYPE_INT8), 1, 10,
		byte(PUSH), byte(TYPE_INT8), 1, 20,
		byte(CALL_NATIVE), 0, 2,
	}

	vm := NewVM(bytecode, []NativeFunc{addFunc})
	if err := vm.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	assertStackValue(t, vm, TYPE_INT8, func(v Value) bool {
		return v.Int8 == 30
	})
}
