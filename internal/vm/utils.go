package vm

import (
	"encoding/binary"
	"fmt"
	"math"
)

func (vm *VM) inferRuntimeValue() (Value, error) {
	codeType := Type(vm.bytecode[vm.pc])
	vm.pc++ // skip type

	length := int(vm.bytecode[vm.pc])
	vm.pc++ // skip length

	switch codeType {
	case TYPE_INT8:
		val := int8(vm.bytecode[vm.pc])
		vm.pc++ // skip data
		return Value{Type: TYPE_INT8, Int8: val}, nil

	case TYPE_INT16:
		val := binary.BigEndian.Uint16(vm.bytecode[vm.pc : vm.pc+2])
		vm.pc += 2 // skip data
		return Value{Type: TYPE_INT16, Int16: int16(val)}, nil

	case TYPE_INT32:
		val := binary.BigEndian.Uint32(vm.bytecode[vm.pc : vm.pc+4])
		vm.pc += 4 // skip data
		return Value{Type: TYPE_INT32, Int32: int32(val)}, nil

	case TYPE_STRING:
		strBytes := vm.bytecode[vm.pc : vm.pc+length]
		vm.pc += length // skip data
		return Value{Type: TYPE_STRING, String: string(strBytes)}, nil

	case TYPE_BOOL:
		boolVal := vm.bytecode[vm.pc] != 0
		vm.pc++ // skip data
		return Value{Type: TYPE_BOOL, Bool: boolVal}, nil

	case TYPE_ARRAY:
		// [TYPE_ARRAY][length: 1 byte][elements...]
		array := make([]Value, length)
		for i := 0; i < length; i++ {
			elem, err := vm.inferRuntimeValue()
			if err != nil {
				return Value{}, err
			}
			array[i] = elem
		}
		return Value{Type: TYPE_ARRAY, Array: array}, nil

	default:
		return Value{}, fmt.Errorf("unsupported type in runtime value inference: 0x%02x", codeType)
	}
}

func determineIntType(v int) (Type, interface{}) {
	if v >= math.MinInt8 && v <= math.MaxInt8 {
		return TYPE_INT8, int8(v)
	}
	if v >= math.MinInt16 && v <= math.MaxInt16 {
		return TYPE_INT16, int16(v)
	}
	return TYPE_INT32, int32(v)
}

// pop
func (vm *VM) pop() (Value, error) {
	if len(vm.dataStack) == 0 {
		return Value{}, fmt.Errorf("stack underflow")
	}
	val := vm.dataStack[len(vm.dataStack)-1]
	vm.dataStack = vm.dataStack[:len(vm.dataStack)-1]
	return val, nil
}

// pop n elements
func (vm *VM) popN(n int) ([]Value, error) {
	if len(vm.dataStack) < n {
		return nil, fmt.Errorf("stack underflow: need %d, have %d", n, len(vm.dataStack))
	}

	args := make([]Value, n)
	copy(args, vm.dataStack[len(vm.dataStack)-n:])

	vm.dataStack = vm.dataStack[:len(vm.dataStack)-n]

	return args, nil
}

// push element
func (vm *VM) push(val Value) {
	vm.dataStack = append(vm.dataStack, val)
}

func (vm *VM) convertInterfaceToValue(val interface{}) (Value, error) {
	switch v := val.(type) {
	case string:
		return Value{Type: TYPE_STRING, String: v}, nil
	case int:
		typ, intVal := determineIntType(v)
		switch typ {
		case TYPE_INT8:
			return Value{Type: TYPE_INT8, Int8: intVal.(int8)}, nil
		case TYPE_INT16:
			return Value{Type: TYPE_INT16, Int16: intVal.(int16)}, nil
		case TYPE_INT32:
			return Value{Type: TYPE_INT32, Int32: intVal.(int32)}, nil
		}
	case bool:
		return Value{Type: TYPE_BOOL, Bool: v}, nil
	case []interface{}:
		array := make([]Value, len(v))
		for i, item := range v {
			converted, err := vm.convertInterfaceToValue(item)
			if err != nil {
				return Value{}, err
			}
			array[i] = converted
		}
		return Value{Type: TYPE_ARRAY, Array: array}, nil
	}
	return Value{}, fmt.Errorf("unsupported type: %T", val)
}

func SerializeLoadGlobal(name string) []byte {
	length := byte(len(name))
	bytes := []byte{byte(LOAD_GLOBAL), length}
	bytes = append(bytes, []byte(name)...)
	return bytes
}

func SerializePush(val interface{}) ([]byte, error) {
	valueBytes, err := serializeValue(val)
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(PUSH)}, valueBytes...), nil
}

func SerializeOperator(op OpCode) []byte {
	return []byte{byte(op)}
}

// format : [type][len][data]
func serializeValue(val interface{}) ([]byte, error) {
	typ, err := inferType(val)
	if err != nil {
		return nil, err
	}

	switch v := val.(type) {
	case int:
		_, intVal := determineIntType(v)
		switch typ {
		case TYPE_INT8:
			return []byte{byte(typ), 1, byte(intVal.(int8))}, nil
		case TYPE_INT16:
			buf := make([]byte, 2)
			binary.BigEndian.PutUint16(buf, uint16(intVal.(int16)))
			return append([]byte{byte(typ), 2}, buf...), nil
		case TYPE_INT32:
			buf := make([]byte, 4)
			binary.BigEndian.PutUint32(buf, uint32(intVal.(int32)))
			return append([]byte{byte(typ), 4}, buf...), nil
		}

	case string:
		payload := []byte(v)
		return append([]byte{byte(typ), byte(len(payload))}, payload...), nil

	case bool:
		var b byte
		if v {
			b = 1
		}
		return []byte{byte(typ), 1, b}, nil

	case []interface{}:
		// array nested — récursif
		buf := []byte{byte(typ), byte(len(v))}
		for _, elem := range v {
			elemBytes, err := serializeValue(elem)
			if err != nil {
				return nil, err
			}
			buf = append(buf, elemBytes...)
		}
		return buf, nil
	}

	return nil, fmt.Errorf("unsupported type: %T", val)
}

func inferType(val interface{}) (Type, error) {
	switch v := val.(type) {
	case int:
		typ, _ := determineIntType(v)
		return typ, nil
	case string:
		return TYPE_STRING, nil
	case bool:
		return TYPE_BOOL, nil
	case []interface{}:
		return TYPE_ARRAY, nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", val)
	}
}
