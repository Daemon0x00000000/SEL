package ast

import (
	"fmt"
	"strings"

	"github.com/Daemon0x00000000/sel/internal/vm"
)

type Node interface {
	compile() ([]byte, error)
}

type LogicalNode struct {
	left        Node
	right       Node
	operator    vm.OpCode
	operatorStr LogicalOperator
}

func treeString(node Node, prefix string, isLast bool) string {
	var sb strings.Builder

	connector := "├── "
	if isLast {
		connector = "└── "
	}
	ext := "│   "
	if isLast {
		ext = "    "
	}

	switch n := node.(type) {
	case *LogicalNode:
		_, err := fmt.Fprintf(&sb, "%s%s%v\n", prefix, connector, n.operatorStr)
		if err != nil {
			return ""
		}
		sb.WriteString(treeString(n.left, prefix+ext, false))
		sb.WriteString(treeString(n.right, prefix+ext, true))

	case *NotNode:
		_, err := fmt.Fprintf(&sb, "%s%sNOT\n", prefix, connector)
		if err != nil {
			return ""
		}
		sb.WriteString(treeString(n.operand, prefix+ext, true))

	case *ComparisonNode:
		_, err := fmt.Fprintf(&sb, "%s%s%s %v %v\n", prefix, connector, n.left, n.operatorStr, n.right)
		if err != nil {
			return ""
		}
	}

	return sb.String()
}

func (n *LogicalNode) compile() ([]byte, error) {
	leftBytes, err := n.left.compile()
	if err != nil {
		return nil, err
	}
	rightBytes, err := n.right.compile()
	if err != nil {
		return nil, err
	}

	bytes := append(leftBytes, rightBytes...)
	return append(bytes, vm.SerializeOperator(n.operator)...), nil
}

type NotNode struct {
	operand Node // child
}

// NOT
func (n *NotNode) compile() ([]byte, error) {
	childBytes, err := n.operand.compile()
	if err != nil {
		return nil, err
	}
	return append(childBytes, vm.SerializeOperator(vm.OP_NOT)...), nil
}

type ComparisonNode struct {
	left        Field
	right       interface{}
	operator    vm.OpCode
	operatorStr ComparisonOperator
}

func (n *ComparisonNode) String() string {
	return fmt.Sprintf("%v %v %v", n.left, n.operator, n.right)
}

// LOAD_GLOBAL <length> <field>
// PUSH <type> <length> <data (right)>
// OPERATOR
func (n *ComparisonNode) compile() ([]byte, error) {
	bytes := vm.SerializeLoadGlobal(string(n.left))

	pushBytes, err := vm.SerializePush(n.right)
	if err != nil {
		return nil, err
	}
	bytes = append(bytes, pushBytes...)

	return append(bytes, vm.SerializeOperator(n.operator)...), nil
}
