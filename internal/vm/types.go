package vm

import "fmt"

type NativeFunc func(args []Value) (Value, error)

type OpCode byte

type Type byte

const (
	TYPE_BOOL   Type = 0x00
	TYPE_INT8   Type = 0x01
	TYPE_INT16  Type = 0x02
	TYPE_INT32  Type = 0x03
	TYPE_STRING Type = 0x04
	TYPE_ARRAY  Type = 0x05
)

type Handler func(*VM) error

type Value struct {
	Type   Type
	Int8   int8
	Int16  int16
	Int32  int32
	String string
	Bool   bool
	Array  []Value
}

func (v Value) Compare(other Value) (int, error) {
	if v.Type != other.Type {
		return 0, fmt.Errorf("cannot compare different types")
	}

	switch v.Type {
	case TYPE_INT8:
		if v.Int8 < other.Int8 {
			return -1, nil
		} else if v.Int8 > other.Int8 {
			return 1, nil
		}
		return 0, nil

	case TYPE_INT16:
		if v.Int16 < other.Int16 {
			return -1, nil
		} else if v.Int16 > other.Int16 {
			return 1, nil
		}
		return 0, nil

	case TYPE_INT32:
		if v.Int32 < other.Int32 {
			return -1, nil
		} else if v.Int32 > other.Int32 {
			return 1, nil
		}
		return 0, nil

	case TYPE_STRING:
		if v.String < other.String {
			return -1, nil
		} else if v.String > other.String {
			return 1, nil
		}
		return 0, nil

	default:
		return 0, fmt.Errorf("type %v does not support comparison", v.Type)
	}
}
