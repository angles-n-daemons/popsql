package grammar

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser"
)

type walkFunc func(Expr) error

type Expr interface {
	Walk(walkFunc) error
}

type ExprVisitor[T any] interface {
	VisitBinaryExpr(Binary) (*T, error)
	VisitLiteralExpr(Literal) (*T, error)
	VisitUnaryExpr(Unary) (*T, error)
	VisitListExpr(List) (*T, error)
}

func Visit[T any](expr Expr, visitor ExprVisitor[T]) (*T, error) {
	switch typedExpr := expr.(type) {
	case Binary:
		return visitor.VisitBinaryExpr(typedExpr)
	case Literal:
		return visitor.VisitLiteralExpr(typedExpr)
	case Unary:
		return visitor.VisitUnaryExpr(typedExpr)
	case List:
		return visitor.VisitListExpr(typedExpr)
	default:
		return nil, fmt.Errorf("unable to visit type %T", typedExpr)
	}
}

type Binary struct {
	left     Expr
	operator parser.Token
	right    Expr
}

func (e Binary) Walk(f walkFunc) error {
	var err error
	err = e.left.Walk(f)

	if err != nil {
		return err
	}

	err = e.right.Walk(f)

	if err != nil {
		return err
	}

	return err
}

type Literal struct {
	value any
}

func (e Literal) Walk(f walkFunc) error {
	var err error
	return err
}

type Unary struct {
	operator parser.Token
	right    Expr
}

func (e Unary) Walk(f walkFunc) error {
	var err error
	err = e.right.Walk(f)

	if err != nil {
		return err
	}

	return err
}

type List struct {
	exprs []Expr
}

func (e List) Walk(f walkFunc) error {
	var err error
	for i := 0; i < len(e.exprs); i++ {
		err = e.exprs[i].Walk(f)

		if err != nil {
			return err
		}

	}
	return err
}
