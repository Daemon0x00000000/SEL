package lql

import (
	"fmt"
	"strings"
)

func Parse(expression string) (*AST, error) {
	ast := newAST()

	if err := validateParentheses(expression); err != nil {
		return nil, err
	}

	rootNode, err := ast.parseExpr(expression)
	if err != nil {
		return nil, err
	}

	ast.root = rootNode
	return ast, nil
}

func validateParentheses(expr string) error {
	depth := 0
	openPositions := []int{}

	for i, char := range expr {
		if char == '(' {
			depth++
			openPositions = append(openPositions, i)
		} else if char == ')' {
			depth--

			if depth < 0 {
				startFloor := max(0, i-5)
				endCeil := min(len(expr), i+5)
				return fmt.Errorf("unbalanced parentheses: closing ')' at position %d has no matching '('\nContext: ...%s...", i, expr[startFloor:endCeil])
			}

			if len(openPositions) > 0 {
				openPositions = openPositions[:len(openPositions)-1]
			}
		}
	}

	if depth > 0 {
		positions := make([]string, len(openPositions))
		for i, pos := range openPositions {
			positions[i] = fmt.Sprintf("%d (...%s...)", pos, expr[pos:min(len(expr), pos+5)])
		}
		return fmt.Errorf("unbalanced parentheses: %d unclosed '(' found at positions: %s", depth, strings.Join(positions, ", "))
	}

	return nil
}

func (ast *AST) parseExpr(expr string) (*LogicalNode, error) {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		depth := 0
		isEnglobing := true
		for i, c := range expr {
			if c == '(' {
				depth++
			} else if c == ')' {
				depth--
			}

			if depth == 0 && i < len(expr)-1 {
				isEnglobing = false
				break
			}
		}

		if isEnglobing {
			return ast.parseExpr(expr[1 : len(expr)-1])
		}
	}

	for _, logicalOp := range logicalOperatorsOrdered {
		idx := ast.findOperatorOutsideParens(expr, string(logicalOp))

		if idx != -1 {
			left := strings.TrimSpace(expr[:idx])
			right := strings.TrimSpace(expr[idx+len(logicalOp):])

			leftNode, err := ast.parseExpr(left)
			if err != nil {
				return nil, fmt.Errorf("error parsing left side: %w", err)
			}

			rightNode, err := ast.parseExpr(right)
			if err != nil {
				return nil, fmt.Errorf("error parsing right side: %w", err)
			}

			return &LogicalNode{
				logicalOperator: logicalOperators[logicalOp],
				children:        []Node{leftNode, rightNode},
				operator:        logicalOp,
			}, nil
		}
	}

	compNode, err := ast.parseTripletComparisonNode(expr)
	if err != nil {
		return nil, err
	}

	return &LogicalNode{
		logicalOperator: nil,
		children:        []Node{compNode},
	}, nil
}

func (ast *AST) parseTripletComparisonNode(expr string) (*ComparisonNode, error) {
	expr = strings.TrimSpace(expr)
	var opFound ComparisonOperator
	var opPos int = -1

	for _, op := range comparisonOpsOrdered {
		if idx := strings.Index(expr, string(op)); idx != -1 {
			opPos = idx
			opFound = op
			break
		}
	}

	if opPos == -1 {
		return nil, fmt.Errorf("no comparison operator found in: %s", expr)
	}

	left := strings.TrimSpace(expr[:opPos])
	right := strings.TrimSpace(expr[opPos+len(opFound):])

	values, err := parseValues(right)
	if err != nil {
		return nil, err
	}

	predicateClosure := comparisonOperators[opFound](values...)

	operationRes := OperationResult(NotResolved)

	if ast.fieldsOperations == nil {
		ast.fieldsOperations = make(map[Field][]OperationCompute)
	}

	ast.fieldsOperations[Field(left)] = append(
		ast.fieldsOperations[Field(left)],
		operationRes.Compute(predicateClosure),
	)

	return NewComparisonNode(
		&operationRes,
		Field(left),
		opFound,
		values,
	), nil
}

func (ast *AST) findOperatorOutsideParens(expr string, operator string) int {
	depth := 0
	inQuotes := false
	escaped := false

	for i := 0; i < len(expr); i++ {
		char := expr[i]

		if escaped {
			escaped = false
			continue
		}

		if char == '\\' && inQuotes {
			escaped = true
			continue
		}

		if char == '\'' {
			inQuotes = !inQuotes
			continue
		}

		if inQuotes {
			continue
		}

		if char == '(' {
			depth++
		} else if char == ')' {
			depth--
		}

		if depth == 0 && !inQuotes && strings.HasPrefix(expr[i:], operator) {
			return i
		}
	}

	return -1
}

func parseValues(input string) ([]string, error) {
	input = strings.TrimSpace(input)

	if !strings.Contains(input, "'") {
		values := strings.Split(input, ",")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}
		return values, nil
	}

	var values []string
	var current strings.Builder
	inQuotes := false
	escaped := false
	wasQuoted := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		switch {
		case escaped:
			switch char {
			case 'n':
				current.WriteByte('\n')
			case 't':
				current.WriteByte('\t')
			case 'r':
				current.WriteByte('\r')
			case '\\':
				current.WriteByte('\\')
			case '\'':
				current.WriteByte('\'')
			default:
				current.WriteByte(char)
			}
			escaped = false

		case char == '\\' && inQuotes:
			escaped = true

		case char == '\'':
			inQuotes = !inQuotes
			wasQuoted = true

		case char == ',' && !inQuotes:
			val := current.String()
			if !wasQuoted {
				val = strings.TrimSpace(val)
			}
			values = append(values, val)
			current.Reset()
			wasQuoted = false

		default:
			current.WriteByte(char)
		}
	}

	val := current.String()
	if !wasQuoted {
		val = strings.TrimSpace(val)
	}
	values = append(values, val)

	if inQuotes {
		return nil, fmt.Errorf("unclosed quote in: %s", input)
	}

	if escaped {
		return nil, fmt.Errorf("incomplete escape sequence at end of: %s", input)
	}

	return values, nil
}
