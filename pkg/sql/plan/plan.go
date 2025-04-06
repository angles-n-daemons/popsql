package plan

import (
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
)

type PlanVisitor[T any] interface {
	VisitCreateTable(*CreateTable) (T, error)
	VisitInsert(*Insert) (T, error)
	VisitScan(*Scan) (T, error)
}

func VisitPlan[T any](plan Plan, visitor PlanVisitor[T]) (T, error) {
	switch typedPlan := plan.(type) {
	case *CreateTable:
		return visitor.VisitCreateTable(typedPlan)
	case *Insert:
		return visitor.VisitInsert(typedPlan)
	case *Scan:
		return visitor.VisitScan(typedPlan)
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

type Scan struct {
	Table *desc.Table
}

func (p *Scan) isPlan() {}

func (p *Scan) Columns() []string {
	columns := p.Table.GetColumns()
	cols := make([]string, len(columns))
	for i, col := range columns {
		cols[i] = col.Name
	}
	return cols
}
