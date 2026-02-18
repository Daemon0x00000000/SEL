package vm

import (
	"fmt"
	"strings"
)

// handlers slice for O(1) access
var handlers = []func(*VM) error{
	PUSH:          (*VM).pushHandler,
	POP:           (*VM).popHandler,
	STORE_GLOBAL:  (*VM).storeGlobalHandler,
	LOAD_GLOBAL:   (*VM).loadGlobalHandler,
	CALL_NATIVE:   (*VM).callNativeHandler,
	OP_EQ:         (*VM).eqHandler,
	OP_GT:         (*VM).gtHandler,
	OP_LT:         (*VM).ltHandler,
	OP_GTE:        (*VM).gteHandler,
	OP_LTE:        (*VM).lteHandler,
	OP_STARTSWITH: (*VM).startsWithHandler,
	OP_ENDSWITH:   (*VM).endsWithHandler,
	OP_CONTAINS:   (*VM).containsHandler,
	OP_IN:         (*VM).inHandler,
	OP_AND:        (*VM).andHandler,
	OP_OR:         (*VM).orHandler,
	OP_XOR:        (*VM).xorHandler,
	OP_NOT:        (*VM).notHandler,
}

// PUSH
func (vm *VM) pushHandler() error {
	//[OP_CODE: 1 byte][type: 1 byte][length: 1 byte][data: length bytes]
	val, err := vm.inferRuntimeValue()
	if err != nil {
		return err
	}
	vm.push(val)

	return nil
}

// POP
func (vm *VM) popHandler() error {
	_, err := vm.pop()
	if err != nil {
		return err
	}

	return nil
}

// STORE_GLOBAL
func (vm *VM) storeGlobalHandler() error {
	//[OP_CODE: 1 byte][length: 1 byte][data: length bytes][type: 1 byte][length: 1 byte][data: length bytes]
	length := int(vm.bytecode[vm.pc])
	vm.pc++ // skip length
	key := string(vm.bytecode[vm.pc : vm.pc+length])
	vm.pc += length // skip key

	val, err := vm.inferRuntimeValue()
	if err != nil {
		return err
	}

	vm.globals[key] = val
	return nil
}

// LOAD_GLOBAL
func (vm *VM) loadGlobalHandler() error {
	length := int(vm.bytecode[vm.pc])
	vm.pc++ // skip length
	key := string(vm.bytecode[vm.pc : vm.pc+length])
	vm.pc += length // skip key

	val, exists := vm.globals[key]
	if !exists {
		return fmt.Errorf("undefined global variable: %s", key)
	}

	vm.push(val)
	return nil
}

// CALL_NATIVE
func (vm *VM) callNativeHandler() error {
	// [OP_CODE: 1 byte][index: 1 byte][args_count: 1 byte] // Max 256 stack pops
	index := int(vm.bytecode[vm.pc])

	vm.pc++ // skip index

	argsCount := int(vm.bytecode[vm.pc])
	vm.pc++ // skip pop count

	if len(vm.nativeFuncs) <= index {
		return fmt.Errorf("native function index out of bounds: %d", index)
	}

	args, err := vm.popN(argsCount)
	if err != nil {
		return err
	}

	val, err := vm.nativeFuncs[index](args)
	if err != nil {
		return fmt.Errorf("native function failed: %v", err)
	}

	vm.push(val)
	return nil
}

// OP_EQ
func (vm *VM) eqHandler() error {
	if len(vm.dataStack) < 2 {
		return fmt.Errorf("not enough values on stack for comparison")
	}

	right, err := vm.pop()
	if err != nil {
		return err
	}
	left, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := left.Compare(right)
	if err != nil {
		return err
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: result == 0})
	return nil
}

// OP_GT
func (vm *VM) gtHandler() error {
	right, err := vm.pop()
	if err != nil {
		return err
	}
	left, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := left.Compare(right)
	if err != nil {
		return err
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: result > 0})
	return nil
}

// OP_LT
func (vm *VM) ltHandler() error {
	right, err := vm.pop()
	if err != nil {
		return err
	}
	left, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := left.Compare(right)
	if err != nil {
		return err
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: result < 0})
	return nil
}

// OP_GTE
func (vm *VM) gteHandler() error {
	right, err := vm.pop()
	if err != nil {
		return err
	}
	left, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := left.Compare(right)
	if err != nil {
		return err
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: result >= 0})
	return nil
}

// OP_LTE
func (vm *VM) lteHandler() error {
	right, err := vm.pop()
	if err != nil {
		return err
	}
	left, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := left.Compare(right)
	if err != nil {
		return err
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: result <= 0})
	return nil
}

// OP_STARTSWITH
func (vm *VM) startsWithHandler() error {
	prefix, err := vm.pop()
	if err != nil {
		return err
	}
	str, err := vm.pop()
	if err != nil {
		return err
	}

	if str.Type != TYPE_STRING || prefix.Type != TYPE_STRING {
		return fmt.Errorf("STARTSWITH requires string operands")
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: strings.HasPrefix(str.String, prefix.String)})
	return nil
}

// OP_ENDSWITH
func (vm *VM) endsWithHandler() error {
	suffix, err := vm.pop()
	if err != nil {
		return err
	}
	str, err := vm.pop()
	if err != nil {
		return err
	}

	if str.Type != TYPE_STRING || suffix.Type != TYPE_STRING {
		return fmt.Errorf("ENDSWITH requires string operands")
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: strings.HasSuffix(str.String, suffix.String)})
	return nil
}

// OP_CONTAINS
func (vm *VM) containsHandler() error {
	substr, err := vm.pop()
	if err != nil {
		return err
	}
	str, err := vm.pop()
	if err != nil {
		return err
	}

	if str.Type != TYPE_STRING || substr.Type != TYPE_STRING {
		return fmt.Errorf("CONTAINS requires string operands")
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: strings.Contains(str.String, substr.String)})
	return nil
}

// OP_IN
func (vm *VM) inHandler() error {
	haystack, err := vm.pop() // array (right operand)
	if err != nil {
		return err
	}
	needle, err := vm.pop() // needle (left operand)
	if err != nil {
		return err
	}

	if haystack.Type != TYPE_ARRAY {
		return fmt.Errorf("IN requires array as right operand")
	}

	// needle is a scalar
	if needle.Type != TYPE_ARRAY {
		found := false
		for _, item := range haystack.Array {
			val, err := needle.Compare(item)
			if err != nil {
				return err
			}
			if val == 0 {
				found = true
				break
			}
		}
		vm.push(Value{Type: TYPE_BOOL, Bool: found})
		return nil
	}

	// needle is an array
	found := false
	for _, needleItem := range needle.Array {
		for _, haystackItem := range haystack.Array {
			val, err := needleItem.Compare(haystackItem)
			if err != nil {
				return err
			}
			if val == 0 {
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: found})
	return nil
}

// NOT
func (vm *VM) notHandler() error {
	val, err := vm.pop()
	if err != nil {
		return err
	}

	if val.Type != TYPE_BOOL {
		return fmt.Errorf("NOT operation requires boolean type")
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: !val.Bool})
	return nil
}

// AND
func (vm *VM) andHandler() error {
	right, err := vm.pop()
	if err != nil {
		return err
	}
	left, err := vm.pop()
	if err != nil {
		return err
	}

	if left.Type != TYPE_BOOL || right.Type != TYPE_BOOL {
		return fmt.Errorf("AND requires boolean operands")
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: left.Bool && right.Bool})
	return nil
}

// OR
func (vm *VM) orHandler() error {
	right, err := vm.pop()
	if err != nil {
		return err
	}
	left, err := vm.pop()
	if err != nil {
		return err
	}

	if left.Type != TYPE_BOOL || right.Type != TYPE_BOOL {
		return fmt.Errorf("OR requires boolean operands")
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: left.Bool || right.Bool})
	return nil
}

// XOR
func (vm *VM) xorHandler() error {
	right, err := vm.pop()
	if err != nil {
		return err
	}
	left, err := vm.pop()
	if err != nil {
		return err
	}

	if left.Type != TYPE_BOOL || right.Type != TYPE_BOOL {
		return fmt.Errorf("XOR requires boolean operands")
	}

	vm.push(Value{Type: TYPE_BOOL, Bool: left.Bool != right.Bool})
	return nil
}
