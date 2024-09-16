package ast

// this file was generated by tools/gen_ast.go

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

type walkFunc func(Expr) error

type Expr interface{}

type ExprVisitor[T any] interface {
	VisitBinaryExpr(*Binary) (*T, error)
	VisitLiteralExpr(*Literal) (*T, error)
	VisitUnaryExpr(*Unary) (*T, error)
	VisitAssignmentExpr(*Assignment) (*T, error)
	VisitReferenceExpr(*Reference) (*T, error)
}

func VisitExpr[T any](expr Expr, visitor ExprVisitor[T]) (*T, error) {
	switch typedExpr := expr.(type) {
	case *Binary:
		return visitor.VisitBinaryExpr(typedExpr)
	case *Literal:
		return visitor.VisitLiteralExpr(typedExpr)
	case *Unary:
		return visitor.VisitUnaryExpr(typedExpr)
	case *Assignment:
		return visitor.VisitAssignmentExpr(typedExpr)
	case *Reference:
		return visitor.VisitReferenceExpr(typedExpr)
	default:
		return nil, fmt.Errorf("unable to visit type %T", typedExpr)
	}
}

type Binary struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

type Literal struct {
	Value scanner.Token
}

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

type Assignment struct {
	Name  scanner.Token
	Value Expr
}

type Reference struct {
	Names []*scanner.Token
}

type StmtVisitor[T any] interface {
	VisitSelectStmt(*Select) (*T, error)
	VisitInsertStmt(*Insert) (*T, error)
}

func VisitStmt[T any](expr Stmt, visitor StmtVisitor[T]) (*T, error) {
	switch typedStmt := expr.(type) {
	case *Select:
		return visitor.VisitSelectStmt(typedStmt)
	case *Insert:
		return visitor.VisitInsertStmt(typedStmt)
	default:
		return nil, fmt.Errorf("unable to visit type %T", typedStmt)
	}
}

type Select struct {
	Terms []Expr
	From  *Reference
	Where Expr
}

type Insert struct {
	Table   *Reference
	Columns []*Reference
	Values  [][]Expr
}

type Stmt interface{}