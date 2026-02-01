package vm

import (
	"fmt"
)

type VM struct {
	bytecode    []byte
	pc          int
	globals     map[string]Value
	dataStack   []Value
	nativeFuncs []NativeFunc // O(1) native funcs access with index
}

func (vm *VM) DataStack() []Value {
	return vm.dataStack
}

func NewVM(bytecode []byte, nativeFuncs []NativeFunc) *VM {
	return &VM{bytecode: bytecode, globals: make(map[string]Value), dataStack: make([]Value, 0), nativeFuncs: nativeFuncs}
}

func (vm *VM) LoadRecords(records map[string]interface{}) error {
	for key, val := range records {
		switch v := val.(type) {
		case string:
			vm.globals[key] = Value{Type: TYPE_STRING, String: v}
		case int:
			// infer precision
			typ, intVal := determineIntType(v)
			switch typ {
			case TYPE_INT8:
				vm.globals[key] = Value{Type: TYPE_INT8, Int8: intVal.(int8)}
			case TYPE_INT16:
				vm.globals[key] = Value{Type: TYPE_INT16, Int16: intVal.(int16)}
			case TYPE_INT32:
				vm.globals[key] = Value{Type: TYPE_INT32, Int32: intVal.(int32)}
			}
		case bool:
			vm.globals[key] = Value{Type: TYPE_BOOL, Bool: v}
		case []interface{}:
			array := make([]Value, len(v))
			for i, item := range v {
				converted, err := vm.convertInterfaceToValue(item)
				if err != nil {
					return fmt.Errorf("failed to convert array element %d: %v", i, err)
				}
				array[i] = converted
			}
			vm.globals[key] = Value{Type: TYPE_ARRAY, Array: array}
		default:
			return fmt.Errorf("unsupported type for field %s", key)
		}
	}
	return nil
}

// Reset the VM state to its initial state (native functions are not reset)
func (vm *VM) Reset() {
	vm.pc = 0
	vm.globals = make(map[string]Value)
	vm.dataStack = make([]Value, 0)
}

func (vm *VM) Execute() error {
	for vm.pc < len(vm.bytecode) {
		opCode := OpCode(vm.bytecode[vm.pc])
		vm.pc++ // skip op code
		handler := handlers[opCode]
		if handler == nil {
			return fmt.Errorf("unknown opcode: 0x%02x at pc=%d", opCode, vm.pc)
		}

		if err := handler(vm); err != nil {
			return err
		}
	}
	return nil
}
