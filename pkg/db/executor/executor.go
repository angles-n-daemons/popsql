package executor

import (
	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/plan"
)

type Row []any

type Executor struct {
	Store   kv.Store
	Catalog *catalog.Manager
}

func New(st kv.Store, cat *catalog.Manager) *Executor {
	return &Executor{
		Store:   st,
		Catalog: cat,
	}
}

func (e *Executor) Execute(p plan.Plan) ([]Row, error) {
	return plan.VisitPlan(p, e)
}
