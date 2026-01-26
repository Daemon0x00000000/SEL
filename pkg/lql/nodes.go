package lql

import (
	"fmt"
	"strings"
)

type Node interface {
	Eval() bool
}

type LogicalNode struct {
	logicalOperator LogicalOperatorFunc
	children        []Node
	operator        LogicalOperator
}

func (n *LogicalNode) String() string {
	return n.treeHelper("", true)
}

func (n *LogicalNode) treeHelper(prefix string, isLast bool) string {
	var sb strings.Builder

	connector := "└── "
	if !isLast {
		connector = "├── "
	}

	if n.logicalOperator == nil {
		fmt.Fprintf(&sb, "%s%s%s\n", prefix, connector, n.children[0])
		return sb.String()
	}

	fmt.Fprintf(&sb, "%s%s%s\n", prefix, connector, n.operator)

	extension := "    "
	if !isLast {
		extension = "│   "
	}

	for i, child := range n.children {
		childIsLast := i == len(n.children)-1

		if logNode, ok := child.(*LogicalNode); ok {
			sb.WriteString(logNode.treeHelper(prefix+extension, childIsLast))
		} else if compNode, ok := child.(*ComparisonNode); ok {
			childConnector := "└── "
			if !childIsLast {
				childConnector = "├── "
			}
			fmt.Fprintf(&sb, "%s%s%s%s\n", prefix, extension, childConnector, compNode)
		}
	}

	return sb.String()
}

func (n *LogicalNode) Eval() bool {
	if n.logicalOperator == nil {
		return n.children[0].Eval()
	}

	results := make([]bool, len(n.children))
	for i, child := range n.children {
		results[i] = child.Eval()
	}

	return n.logicalOperator(results)
}

type ComparisonNode struct {
	operation *OperationResult
	field     Field
	operator  ComparisonOperator
	value     interface{}
}

func (cNode *ComparisonNode) String() string {
	return fmt.Sprintf("%v %v %v", cNode.field, cNode.operator, cNode.value)
}

func (cNode *ComparisonNode) Eval() bool {
	return *cNode.operation == True
}

func NewComparisonNode(operation *OperationResult, field Field, operator ComparisonOperator, value interface{}) *ComparisonNode {
	return &ComparisonNode{
		operation: operation,
		field:     field,
		operator:  operator,
		value:     value,
	}
}
