package vm

import (
	"testing"
)

// ============================================================================
// Stack Operations Tests
// ============================================================================

func TestStackOperations(t *testing.T) {
	t.Run("pop", func(t *testing.T) {
		vm := NewVM([]byte{}, nil)
		vm.push(Value{Type: TYPE_INT8, Int8: 42})

		val, err := vm.pop()
		if err != nil {
			t.Fatalf("pop failed: %v", err)
		}
		if val.Int8 != 42 {
			t.Errorf("Expected 42, got %d", val.Int8)
		}
		if len(vm.dataStack) != 0 {
			t.Errorf("Expected empty stack after pop, got %d elements", len(vm.dataStack))
		}
	})

	t.Run("pop underflow", func(t *testing.T) {
		vm := NewVM([]byte{}, nil)
		_, err := vm.pop()
		if err == nil {
			t.Fatal("Expected error for pop on empty stack, got nil")
		}
	})

	t.Run("push", func(t *testing.T) {
		vm := NewVM([]byte{}, nil)
		vm.push(Value{Type: TYPE_INT8, Int8: 1})
		vm.push(Value{Type: TYPE_INT8, Int8: 2})
		vm.push(Value{Type: TYPE_INT8, Int8: 3})

		if len(vm.dataStack) != 3 {
			t.Errorf("Expected 3 elements on stack, got %d", len(vm.dataStack))
		}
		if vm.dataStack[0].Int8 != 1 {
			t.Errorf("Expected first element to be 1, got %d", vm.dataStack[0].Int8)
		}
	})

	t.Run("popN", func(t *testing.T) {
		vm := NewVM([]byte{}, nil)
		vm.push(Value{Type: TYPE_INT8, Int8: 1})
		vm.push(Value{Type: TYPE_INT8, Int8: 2})
		vm.push(Value{Type: TYPE_INT8, Int8: 3})

		vals, err := vm.popN(2)
		if err != nil {
			t.Fatalf("popN failed: %v", err)
		}
		if len(vals) != 2 {
			t.Errorf("Expected 2 values, got %d", len(vals))
		}
		if vals[0].Int8 != 2 || vals[1].Int8 != 3 {
			t.Errorf("Expected [2,3], got [%d,%d]", vals[0].Int8, vals[1].Int8)
		}
		if len(vm.dataStack) != 1 {
			t.Errorf("Expected 1 element left on stack, got %d", len(vm.dataStack))
		}
	})

	t.Run("popN underflow", func(t *testing.T) {
		vm := NewVM([]byte{}, nil)
		vm.push(Value{Type: TYPE_INT8, Int8: 1})

		_, err := vm.popN(5)
		if err == nil {
			t.Fatal("Expected error for popN with insufficient elements, got nil")
		}
	})
}
