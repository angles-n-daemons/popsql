package plan

import (
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
)

type PlanVisitor[T any] interface {
	VisitCreateTable(*CreateTable) (T, error)
	VisitInsert(*Insert) (T, error)
}

func VisitPlan[T any](plan Plan, visitor PlanVisitor[T]) (T, error) {
	switch typedPlan := plan.(type) {
	case *CreateTable:
		return visitor.VisitCreateTable(typedPlan)
	case *Insert:
		return visitor.VisitInsert(typedPlan)
	default:
		return *new(T), nil
	}
}

type Plan interface {
	isPlan()
}

type CreateTable struct {
	Table *desc.Table
}

func (p *CreateTable) isPlan() {}

type Insert struct {
	Table   *desc.Table
	Columns []*desc.Column
	Values  [][]ast.Expr
}

func (p *Insert) isPlan() {}
