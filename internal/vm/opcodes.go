package vm

const (
	PUSH          OpCode = 0x00
	POP           OpCode = 0x01
	STORE_GLOBAL  OpCode = 0x02
	LOAD_GLOBAL   OpCode = 0x03
	CALL_NATIVE   OpCode = 0x04
	OP_EQ         OpCode = 0x05
	OP_GT         OpCode = 0x06
	OP_LT         OpCode = 0x07
	OP_GTE        OpCode = 0x08
	OP_LTE        OpCode = 0x09
	OP_STARTSWITH OpCode = 0x0A
	OP_ENDSWITH   OpCode = 0x0B
	OP_CONTAINS   OpCode = 0x0C
	OP_IN         OpCode = 0x0D
	OP_AND        OpCode = 0x0E
	OP_OR         OpCode = 0x0F
	OP_XOR        OpCode = 0x10
	OP_NOT        OpCode = 0x11
)

func (op OpCode) isLogical() bool {
	return op >= OP_AND && op <= OP_NOT
}

func (op OpCode) isComparison() bool {
	return op >= OP_EQ && op <= OP_IN
}
