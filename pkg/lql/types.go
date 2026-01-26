package lql

type Field string
type LogicalOperator string
type ComparisonOperator string

// []bool est un tableau de 2 valeurs
type LogicalOperatorFunc func(results []bool) bool
type PredicateClosure func(value interface{}) bool
type OperationProvider func(expectedValues ...string) PredicateClosure
type OperationCompute func(value interface{})

type OperationResult int8

const (
	NotResolved OperationResult = -1
	False       OperationResult = 0
	True        OperationResult = 1
)

func (res *OperationResult) Compute(predicate PredicateClosure) OperationCompute {
	return func(value interface{}) {
		if predicate(value) {
			*res = True
		} else {
			*res = False
		}
	}
}

type LogicalOperatorMapping map[LogicalOperator]LogicalOperatorFunc
type ComparisonOperatorMapping map[ComparisonOperator]OperationProvider
