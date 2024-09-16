package parser

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

func Parse(s string) (ast.Stmt, error) {
	tokens, err := scanner.Scan(s)
	if err != nil {
		return nil, err
	}
	stmt, i, err := statement(tokens, 0)
	if err != nil {
		return nil, err
	}
	if !isAtEnd(tokens, i) {
		return nil, fmt.Errorf("finished parsing without consuming all input")
	}
	return stmt, err
}

type expressionSig func([]*scanner.Token, int) (ast.Expr, int, error)

func statement(tokens []*scanner.Token, i int) (ast.Stmt, int, error) {
	switch tokens[i].Type {
	case scanner.SELECT:
		return selectStmt(tokens, i+1)
	case scanner.INSERT:
		return insertStmt(tokens, i+1)
	default:
		return nil, i, fmt.Errorf("unexpected token %s looking for statement", tokens[i].Type)
	}
}

func selectStmt(tokens []*scanner.Token, i int) (ast.Stmt, int, error) {
	terms, i, err := expressionList(tokens, i)
	if err != nil {
		return nil, i, err
	}
	stmt := &ast.Select{Terms: terms}
	if match(tokens, i, scanner.FROM) {
		if len(tokens) <= i+1 {
			return nil, i, fmt.Errorf("reached end of input looking for 'from' expression")
		}

		var from *ast.Reference
		from, i, err = reference(tokens, i+1)
		if err != nil {
			return nil, i, err
		}

		stmt.From = from
	}
	if match(tokens, i, scanner.WHERE) {
		if len(tokens) <= i+1 {
			return nil, i, fmt.Errorf("reached end of input looking for 'where' expression")
		}

		var where ast.Expr
		where, i, err = expression(tokens, i+1)
		if err != nil {
			return nil, i, err
		}

		stmt.Where = where
	}
	return stmt, i + 1, nil
}

func insertStmt(tokens []*scanner.Token, i int) (ast.Stmt, int, error) {
	err := assertNext(tokens, i, scanner.INTO)
	if err != nil {
		return nil, i, err
	}

	table, i, err := reference(tokens, i+1)
	if err != nil {
		return nil, i, err
	}

	err = assertNext(tokens, i, scanner.LEFT_PAREN)
	if err != nil {
		return nil, i, err
	}
	columns, i, err := referenceList(tokens, i+1)
	if err != nil {
		return nil, i, err
	}
	err = assertNext(tokens, i, scanner.RIGHT_PAREN)
	if err != nil {
		return nil, i, err
	}

	err = assertNext(tokens, i+1, scanner.VALUES)
	if err != nil {
		return nil, i, err
	}

	values, i, err := tupleList(tokens, i+2)
	if err != nil {
		return nil, i, err
	}

	return &ast.Insert{Table: table, Columns: columns, Values: values}, i, nil
}

func tupleList(tokens []*scanner.Token, i int) ([][]ast.Expr, int, error) {
	var tup []ast.Expr
	var err error
	list := [][]ast.Expr{}
	for {
		tup, i, err = tuple(tokens, i)
		if err != nil {
			return nil, i, err
		}
		list = append(list, tup)
		if match(tokens, i, scanner.COMMA) {
			i++
		} else {
			break
		}
	}
	return list, i, nil
}

func tuple(tokens []*scanner.Token, i int) ([]ast.Expr, int, error) {
	err := assertNext(tokens, i, scanner.LEFT_PAREN)
	if err != nil {
		return nil, i, err
	}
	list, i, err := expressionList(tokens, i+1)
	err = assertNext(tokens, i, scanner.RIGHT_PAREN)
	if err != nil {
		return nil, i, err
	}
	return list, i + 1, nil
}

// expression list requries a minimum of one element
func expressionList(tokens []*scanner.Token, i int) ([]ast.Expr, int, error) {
	var expr ast.Expr
	var err error
	list := []ast.Expr{}
	for {
		expr, i, err = expression(tokens, i)
		if err != nil {
			return nil, i, err
		}
		list = append(list, expr)
		if match(tokens, i, scanner.COMMA) {
			i++
		} else {
			break
		}
	}
	return list, i, nil
}

func referenceList(tokens []*scanner.Token, i int) ([]*ast.Reference, int, error) {
	var ref *ast.Reference
	var err error
	list := []*ast.Reference{}
	for {
		ref, i, err = reference(tokens, i)
		if err != nil {
			return nil, i, err
		}
		list = append(list, ref)
		if match(tokens, i, scanner.COMMA) {
			i++
		} else {
			break
		}
	}
	return list, i, nil
}

func expression(tokens []*scanner.Token, i int) (ast.Expr, int, error) {
	if isAtEnd(tokens, i) {
		return nil, i, fmt.Errorf("reached end of input parsing expression")
	}
	return equality(tokens, i)
}

func equality(tokens []*scanner.Token, i int) (ast.Expr, int, error) {
	return binary(
		tokens,
		i,
		comparison,
		scanner.BANG_EQUAL,
		scanner.EQUAL_EQUAL,
	)
}

func comparison(tokens []*scanner.Token, i int) (ast.Expr, int, error) {
	return binary(
		tokens,
		i,
		term,
		scanner.GREATER,
		scanner.GREATER_EQUAL,
		scanner.LESS,
		scanner.LESS_EQUAL,
	)
}

func term(tokens []*scanner.Token, i int) (ast.Expr, int, error) {
	return binary(
		tokens,
		i,
		factor,
		scanner.PLUS,
		scanner.MINUS,
	)
}

func factor(tokens []*scanner.Token, i int) (ast.Expr, int, error) {
	return binary(
		tokens,
		i,
		unary,
		scanner.STAR,
		scanner.SLASH,
	)
}

func unary(tokens []*scanner.Token, i int) (ast.Expr, int, error) {
	if i >= len(tokens) {
		return nil, i, fmt.Errorf("reached end of input parsing expression")
	}
	var expr ast.Expr
	var err error
	if match(tokens, i, scanner.BANG, scanner.MINUS) {
		operator := tokens[i]
		expr, i, err = unary(tokens, i+1)
		if err != nil {
			return nil, i, err
		}
		return &ast.Unary{Operator: *operator, Right: expr}, i, nil
	}
	return primary(tokens, i)
}

func binary(
	tokens []*scanner.Token, i int, next expressionSig, operators ...scanner.TokenType,
) (ast.Expr, int, error) {
	if i >= len(tokens) {
		return nil, i, fmt.Errorf("reached end of input parsing expression")
	}
	expr, i, err := next(tokens, i)
	if err != nil {
		return nil, i, err
	}
	for match(tokens, i, operators...) {
		operator := *tokens[i]
		i++
		var right ast.Expr
		right, i, err = next(tokens, i)
		if err != nil {
			return nil, i, err
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
		return nil, i, fmt.Errorf("reached end of input parsing expression")
	}
	switch tokens[i].Type {
	case scanner.NUMBER, scanner.STRING:
		return &ast.Literal{Value: *tokens[i]}, i + 1, nil
	case scanner.IDENTIFIER, scanner.STAR:
		return reference(tokens, i)
	case scanner.LEFT_PAREN:
		expr, i, err = expression(tokens, i)
		if err != nil {
			return nil, i, err
		}
		if !match(tokens, i, scanner.RIGHT_PAREN) {
			return nil, i, fmt.Errorf("expected ')' after expression")
		}
		return expr, i, nil
	default:
		return nil, i, fmt.Errorf("unexpected token %s found while parsing primary", tokens[i].Type)
	}
}

func reference(tokens []*scanner.Token, i int) (*ast.Reference, int, error) {
	names := []*scanner.Token{tokens[i]}
	for match(tokens, i+1, scanner.DOT) {
		i += 2
		if !match(tokens, i, scanner.IDENTIFIER, scanner.STAR) {
			return nil, i, fmt.Errorf("unexpected token %s parsing reference", tokens[i].Type)
		}
		names = append(names, tokens[i])
	}
	return &ast.Reference{Names: names}, i + 1, nil
}

func assertNext(tokens []*scanner.Token, i int, ttype scanner.TokenType) error {
	if isAtEnd(tokens, i) {
		return fmt.Errorf("reached end of input looking for '%s'", ttype)
	}

	if tokens[i].Type != ttype {
		return fmt.Errorf("expected token '%s' but got '%s'", ttype, tokens[i].Type)
	}
	return nil
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