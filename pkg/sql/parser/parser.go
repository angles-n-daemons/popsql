package parser

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

func Parse(s string) (ast.Expr, error) {
	tokens, err := scanner.Scan(s)
	if err != nil {
		return nil, err
	}
	expr, i, err := expression(tokens, 0)
	if err != nil {
		return nil, err
	}
	if !isAtEnd(tokens, i) {
		return nil, fmt.Errorf("finished parsing without consuming all input")
	}
	return expr, err
}

func expression(tokens []*scanner.Token, i int) (ast.Expr, int, error) {
	if i >= len(tokens) {
		return nil, 0, fmt.Errorf("reached end of input parsing expression")
	}
	return equality(tokens, i)
}

func equality(tokens []*scanner.Token, i int) (ast.Expr, int, error) {
	if i >= len(tokens) {
		return nil, 0, fmt.Errorf("reached end of input parsing expression")
	}
	expr, i, err := primary(tokens, i)
	if err != nil {
		return nil, 0, err
	}
	for match(tokens, i, scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := *tokens[i]
		i++
		var right ast.Expr
		right, i, err = primary(tokens, i)
		if err != nil {
			return nil, 0, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, i, nil
}

func primary(tokens []*scanner.Token, i int) (ast.Expr, int, error) {
	var expr ast.Expr
	var err error
	if i >= len(tokens) {
		return nil, 0, fmt.Errorf("reached end of input parsing expression")
	}
	switch tokens[i].Type {
	case scanner.STRING:
		return &ast.Literal{Value: tokens[i].Literal}, i + 1, nil
	case scanner.LEFT_PAREN:
		expr, i, err = expression(tokens, i)
		if err != nil {
			return nil, 0, err
		}
		if !match(tokens, i, scanner.RIGHT_PAREN) {
			return nil, 0, fmt.Errorf("expected ')' after expression")
		}
		return expr, i, nil
	default:
		return nil, 0, fmt.Errorf("unexpected token %s found while parsing primary", tokens[i].Type)
	}
}

func match(tokens []*scanner.Token, i int, types ...scanner.TokenType) bool {
	if isAtEnd(tokens, i) {
		return false
	}
	for _, ttype := range types {
		if tokens[i].Type == ttype {
			return true
		}
	}
	return false
}

func isAtEnd(tokens []*scanner.Token, i int) bool {
	return len(tokens) <= i
}
