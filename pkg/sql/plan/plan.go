package plan

import "github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"

type PlanVisitor[T any] interface {
	VisitCreateTable(*CreateTable) (T, error)
}

func VisitPlan[T any](plan Plan, visitor PlanVisitor[T]) (T, error) {
	switch typedPlan := plan.(type) {
	case *CreateTable:
		return visitor.VisitCreateTable(typedPlan)
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
