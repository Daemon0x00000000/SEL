package ast

import (
	"fmt"
	"strings"
)

func Parse(expression string) (*AST, error) {
	ast := newAST()

	if err := validateParentheses(expression); err != nil {
		return nil, err
	}

	rootNode, err := parseExpr(expression)
	if err != nil {
		return nil, err
	}

	ast.root = rootNode
	return ast, nil
}

func validateParentheses(expr string) error {
	depth := 0
	openPositions := make([]int, 0)

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
		// Build context around first unclosed paren
		firstPos := openPositions[0]
		startFloor := max(0, firstPos-5)
		endCeil := min(len(expr), firstPos+5)
		return fmt.Errorf("unbalanced parentheses: %d unclosed '(' found at positions: %s\nContext: ...%s...", depth, strings.Join(positions, ", "), expr[startFloor:endCeil])
	}

	return nil
}

func parseExpr(expr string) (Node, error) {
	expr = strings.TrimSpace(expr)

	// strip outer parentheses
	if inner, ok := stripOuterParens(expr); ok {
		return parseExpr(inner)
	}

	// Logical operators first (lower priority than NOT in parsing order)
	for _, logicalOp := range logicalOperatorsOrdered {
		idx := findOperatorOutsideParens(expr, string(logicalOp))
		if idx == -1 {
			continue
		}

		left := strings.TrimSpace(expr[:idx])
		right := strings.TrimSpace(expr[idx+len(logicalOp):])

		leftNode, err := parseExpr(left)
		if err != nil {
			return nil, fmt.Errorf("left of %s: %w", logicalOp, err)
		}
		rightNode, err := parseExpr(right)
		if err != nil {
			return nil, fmt.Errorf("right of %s: %w", logicalOp, err)
		}

		return &LogicalNode{
			operator:    logicalOperators[logicalOp],
			operatorStr: logicalOp,
			left:        leftNode,
			right:       rightNode,
		}, nil
	}

	// NOT - handle ! prefix (including !! for double NOT or !field=value)
	if strings.HasPrefix(expr, "!") {
		rest := strings.TrimSpace(expr[1:])
		node, err := parseExpr(rest)
		if err != nil {
			return nil, err
		}
		return &NotNode{operand: node}, nil
	}

	return parseComparison(expr)
}

func stripOuterParens(expr string) (string, bool) {
	if !strings.HasPrefix(expr, "(") || !strings.HasSuffix(expr, ")") {
		return expr, false
	}
	depth := 0
	for i, c := range expr {
		if c == '(' {
			depth++
		} else if c == ')' {
			depth--
		}
		if depth == 0 && i < len(expr)-1 {
			return expr, false
		}
	}
	return expr[1 : len(expr)-1], true
}

func parseComparison(expr string) (Node, error) {
	expr = strings.TrimSpace(expr)

	var opFound ComparisonOperator
	var opPos int = -1
	var isNegated bool

	for _, op := range comparisonOpsOrdered {
		if idx := strings.Index(expr, "!"+string(op)); idx != -1 {
			opPos = idx
			opFound = op
			isNegated = true
			break
		}
		if idx := strings.Index(expr, string(op)); idx != -1 {
			opPos = idx
			opFound = op
			break
		}
	}

	if opPos == -1 {
		return nil, fmt.Errorf("no comparison operator found in: %s", expr)
	}

	if opPos == 0 {
		return nil, fmt.Errorf("missing field before operator in: %s", expr)
	}

	opLen := len(opFound)
	if isNegated {
		opLen++
	}
	// double operator
	afterOp := opPos + opLen
	if afterOp < len(expr) {
		remaining := expr[afterOp:]
		for _, op := range comparisonOpsOrdered {
			if strings.HasPrefix(remaining, string(op)) {
				return nil, fmt.Errorf("double operator found in: %s", expr)
			}
		}
	}

	left := Field(strings.TrimSpace(expr[:opPos]))
	rawRight := strings.TrimSpace(expr[opPos+opLen:])

	values, err := parseValues(rawRight)
	if err != nil {
		return nil, err
	}

	var right interface{}
	if opFound == IN {
		arr := make([]interface{}, len(values))
		for i, v := range values {
			arr[i] = v
		}
		right = arr
	} else {
		if len(values) != 1 {
			return nil, fmt.Errorf("operator %s expects single value, got %d", opFound, len(values))
		}
		right = values[0]
	}

	node := &ComparisonNode{
		left:        left,
		operator:    comparisonOperators[opFound],
		operatorStr: opFound,
		right:       right,
	}

	if isNegated {
		return &NotNode{operand: node}, nil
	}
	return node, nil
}

func findOperatorOutsideParens(expr string, operator string) int {
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
		}

		if inQuotes {
			continue
		}

		if char == '(' {
			depth++
		} else if char == ')' {
			depth--
		}

		if depth == 0 && strings.HasPrefix(expr[i:], operator) {
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

	if escaped {
		return nil, fmt.Errorf("incomplete escape sequence at end of: %s", input)
	}

	if inQuotes {
		return nil, fmt.Errorf("unclosed quote in: %s", input)
	}

	return values, nil
}
