package vm

import (
	"testing"
)

// ============================================================================
// VM Core Tests
// ============================================================================

func TestNewVM(t *testing.T) {
	bytecode := []byte{byte(PUSH), byte(TYPE_INT8), 1, 42}
	nativeFuncs := []NativeFunc{}

	vm := NewVM(bytecode, nativeFuncs)

	if vm == nil {
		t.Fatal("NewVM returned nil")
	}
	if len(vm.bytecode) != len(bytecode) {
		t.Errorf("Expected bytecode length %d, got %d", len(bytecode), len(vm.bytecode))
	}
	if vm.pc != 0 {
		t.Errorf("Expected pc to be 0, got %d", vm.pc)
	}
	if len(vm.globals) != 0 {
		t.Errorf("Expected empty globals, got %d entries", len(vm.globals))
	}
	if len(vm.dataStack) != 0 {
		t.Errorf("Expected empty dataStack, got %d entries", len(vm.dataStack))
	}
}

func TestLoadRecords(t *testing.T) {
	tests := []struct {
		name    string
		records map[string]interface{}
		checks  []func(*testing.T, *VM)
	}{
		{
			name:    "string",
			records: map[string]interface{}{"name": "Alice"},
			checks: []func(*testing.T, *VM){
				func(t *testing.T, vm *VM) {
					val, exists := vm.globals["name"]
					if !exists {
						t.Fatal("Expected 'name' to be in globals")
					}
					if val.Type != TYPE_STRING || val.String != "Alice" {
						t.Errorf("Expected string(Alice), got %+v", val)
					}
				},
			},
		},
		{
			name:    "int8",
			records: map[string]interface{}{"age": 42},
			checks: []func(*testing.T, *VM){
				func(t *testing.T, vm *VM) {
					val := vm.globals["age"]
					if val.Type != TYPE_INT8 || val.Int8 != 42 {
						t.Errorf("Expected int8(42), got %+v", val)
					}
				},
			},
		},
		{
			name:    "int16",
			records: map[string]interface{}{"port": 8080},
			checks: []func(*testing.T, *VM){
				func(t *testing.T, vm *VM) {
					val := vm.globals["port"]
					if val.Type != TYPE_INT16 || val.Int16 != 8080 {
						t.Errorf("Expected int16(8080), got %+v", val)
					}
				},
			},
		},
		{
			name:    "int32",
			records: map[string]interface{}{"population": 1000000},
			checks: []func(*testing.T, *VM){
				func(t *testing.T, vm *VM) {
					val := vm.globals["population"]
					if val.Type != TYPE_INT32 || val.Int32 != 1000000 {
						t.Errorf("Expected int32(1000000), got %+v", val)
					}
				},
			},
		},
		{
			name:    "bool",
			records: map[string]interface{}{"active": true},
			checks: []func(*testing.T, *VM){
				func(t *testing.T, vm *VM) {
					val := vm.globals["active"]
					if val.Type != TYPE_BOOL || val.Bool != true {
						t.Errorf("Expected bool(true), got %+v", val)
					}
				},
			},
		},
		{
			name:    "array",
			records: map[string]interface{}{"tags": []interface{}{"admin", "user"}},
			checks: []func(*testing.T, *VM){
				func(t *testing.T, vm *VM) {
					val := vm.globals["tags"]
					if val.Type != TYPE_ARRAY || len(val.Array) != 2 {
						t.Errorf("Expected array[2], got %+v", val)
					}
					if val.Array[0].String != "admin" || val.Array[1].String != "user" {
						t.Errorf("Expected [admin, user], got %+v", val.Array)
					}
				},
			},
		},
		{
			name:    "array with ints",
			records: map[string]interface{}{"numbers": []interface{}{1, 2, 3}},
			checks: []func(*testing.T, *VM){
				func(t *testing.T, vm *VM) {
					val := vm.globals["numbers"]
					if val.Type != TYPE_ARRAY || len(val.Array) != 3 {
						t.Errorf("Expected array[3], got %+v", val)
					}
					for i, expected := range []int8{1, 2, 3} {
						if val.Array[i].Int8 != expected {
							t.Errorf("Expected numbers[%d]=%d, got %d", i, expected, val.Array[i].Int8)
						}
					}
				},
			},
		},
		{
			name:    "multiple records",
			records: map[string]interface{}{"name": "Bob", "age": 30, "active": false},
			checks: []func(*testing.T, *VM){
				func(t *testing.T, vm *VM) {
					if len(vm.globals) != 3 {
						t.Errorf("Expected 3 globals, got %d", len(vm.globals))
					}
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewVM([]byte{}, nil)
			err := vm.LoadRecords(tt.records)
			if err != nil {
				t.Fatalf("LoadRecords failed: %v", err)
			}
			for _, check := range tt.checks {
				check(t, vm)
			}
		})
	}
}

func TestReset(t *testing.T) {
	vm := NewVM([]byte{1, 2, 3}, nil)
	vm.pc = 2
	vm.globals["test"] = Value{Type: TYPE_STRING, String: "value"}
	vm.dataStack = append(vm.dataStack, Value{Type: TYPE_INT8, Int8: 42})

	vm.Reset()

	if vm.pc != 0 {
		t.Errorf("Expected pc to be 0 after reset, got %d", vm.pc)
	}
	if len(vm.globals) != 0 {
		t.Errorf("Expected empty globals after reset, got %d entries", len(vm.globals))
	}
	if len(vm.dataStack) != 0 {
		t.Errorf("Expected empty dataStack after reset, got %d entries", len(vm.dataStack))
	}
}

func TestExecute_SimplePush(t *testing.T) {
	bytecode := []byte{byte(PUSH), byte(TYPE_INT8), 1, 42}
	vm := NewVM(bytecode, nil)

	if err := vm.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	assertStackValue(t, vm, TYPE_INT8, func(v Value) bool {
		return v.Int8 == 42
	})
}

func TestDataStack(t *testing.T) {
	vm := NewVM([]byte{}, nil)
	vm.dataStack = append(vm.dataStack, Value{Type: TYPE_INT8, Int8: 1})
	vm.dataStack = append(vm.dataStack, Value{Type: TYPE_INT8, Int8: 2})

	stack := vm.DataStack()

	if len(stack) != 2 {
		t.Errorf("Expected 2 values on stack, got %d", len(stack))
	}
	if stack[0].Int8 != 1 || stack[1].Int8 != 2 {
		t.Error("Stack values don't match expected values")
	}
}
