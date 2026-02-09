package vm

import (
	"cmp"
	"fmt"
)

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
	case TYPE_BOOL:
		if v.Bool == other.Bool {
			return 0, nil
		} else if !v.Bool && other.Bool {
			return -1, nil
		}

		return 1, nil
	case TYPE_INT8:
		return cmpGeneric(v.Int8, other.Int8), nil

	case TYPE_INT16:
		return cmpGeneric(v.Int16, other.Int16), nil

	case TYPE_INT32:
		return cmpGeneric(v.Int32, other.Int32), nil

	case TYPE_STRING:
		return cmpGeneric(v.String, other.String), nil

	case TYPE_ARRAY:
		minLen := len(v.Array)
		if len(other.Array) < minLen {
			minLen = len(other.Array)
		}

		for i := 0; i < minLen; i++ {
			elem := v.Array[i]

			res, err := elem.Compare(other.Array[i])
			if err != nil { //  [1,3,-2] [1,2,0] = [0, 1, -1]
				return 0, err
			}

			if res != 0 {
				return res, nil // first diff
			}
		}

		if len(v.Array) < len(other.Array) {
			return -1, nil
		} else if len(v.Array) > len(other.Array) {
			return 1, nil
		}

		return 0, nil

	default:
		return 0, fmt.Errorf("type %v does not support comparison", v.Type)
	}
}

func cmpGeneric[T cmp.Ordered](a, b T) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}
