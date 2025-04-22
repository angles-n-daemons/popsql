package execution

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/db/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/parser/scanner"
)

// The executor also acts as visitor for expressions.
// I need to figure out scoping here, but that's for another day.

func Eval(e *Executor, expr ast.Expr) (any, error) {
	return ast.VisitExpr(expr, e)
}

func (e *Executor) VisitBinaryExpr(expr *ast.Binary) (any, error) {
	lhs, err := Eval(e, expr.Left)
	if err != nil {
		return nil, err
	}
	rhs, err := Eval(e, expr.Right)
	if err != nil {
		return nil, err
	}
	return evalBinaryExpr(expr.Operator, lhs, rhs)
}

func evalBinaryExpr(op *scanner.Token, left, right any) (any, error) {
	switch op.Type {
	case scanner.PLUS:
		switch left.(type) {
		case string:
			leftStr := left.(string)
			if rightStr, ok := right.(string); ok {
				return leftStr + rightStr, nil
			}
		case float64:
			return arithmetic(op, left, right)
		}
		return nil, fmt.Errorf("cannot add values of type %T and %T", left, right)
	case scanner.MINUS, scanner.STAR, scanner.SLASH:
		return arithmetic(op, left, right)
	case scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL:
		return compare(op, left, right)
	case scanner.EQUAL_EQUAL, scanner.BANG_EQUAL:
		return equality(op, left, right)
	default:
		return nil, fmt.Errorf("unsupported binary operator: %s", op)
	}
}

func arithmetic(op *scanner.Token, left, right any) (any, error) {
	a, aok := left.(float64)
	b, bok := right.(float64)
	if !(aok && bok) {
		return nil, fmt.Errorf("cannot do arithmetic on values of type %T and %T", a, b)
	}
	switch op.Type {
	case scanner.PLUS:
		return a + b, nil
	case scanner.MINUS:
		return a - b, nil
	case scanner.STAR:
		return a * b, nil
	case scanner.SLASH:
		return a / b, nil
	default:
		return nil, fmt.Errorf("unsupported arithmetic operator: %s", op)
	}
}

func compare(op *scanner.Token, left, right any) (any, error) {
	a, aok := left.(float64)
	b, bok := right.(float64)
	if !(aok && bok) {
		return nil, fmt.Errorf("cannot do comparison on values of type %T and %T", a, b)
	}
	switch op.Type {
	case scanner.GREATER:
		return a > b, nil
	case scanner.GREATER_EQUAL:
		return a >= b, nil
	case scanner.LESS:
		return a < b, nil
	case scanner.LESS_EQUAL:
		return a <= b, nil
	default:
		return nil, fmt.Errorf("unsupported comparison operator: %s", op)
	}
}

func equality(op *scanner.Token, left, right any) (any, error) {
	switch op.Type {
	case scanner.EQUAL_EQUAL:
		return left == right, nil
	case scanner.BANG_EQUAL:
		return left != right, nil
	}
	return nil, fmt.Errorf("unsupported equality operator: %s", op)
}

func (e *Executor) VisitLiteralExpr(expr *ast.Literal) (any, error) {
	return expr.Value.Literal, nil
}
func (e *Executor) VisitUnaryExpr(expr *ast.Unary) (any, error) {
	right, err := Eval(e, expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case scanner.BANG:
		if boolVal, ok := right.(bool); ok {
			return !boolVal, nil
		}
	case scanner.MINUS:
		if num, ok := right.(float64); ok {
			return -num, nil
		}

	}
	return nil, fmt.Errorf("cannot perform operaion %s on value of type %T", expr.Operator.Type, right)
}
func (e *Executor) VisitIdentifierExpr(*ast.Identifier) (any, error) {
	return nil, fmt.Errorf("not implemented")
}
func (e *Executor) VisitColumnSpecExpr(spec *ast.ColumnSpec) (any, error) {
	return nil, fmt.Errorf("the executor should not see a column spec: name '%s', type '%s'", spec.Name.Name.Lexeme, spec.DataType.Lexeme)
}
