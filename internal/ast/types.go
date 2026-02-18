package ast

import "github.com/Daemon0x00000000/sel/internal/vm"

type Field string
type LogicalOperator string
type ComparisonOperator string

// []bool est un tableau de 2 valeurs
type LogicalOperatorFunc func(results []bool) bool
type PredicateClosure func(value interface{}) bool
type OperationProvider func(expectedValues ...string) PredicateClosure
type OperationCompute func(value interface{})

type LogicalOperatorMapping map[LogicalOperator]vm.OpCode
type ComparisonOperatorMapping map[ComparisonOperator]vm.OpCode
