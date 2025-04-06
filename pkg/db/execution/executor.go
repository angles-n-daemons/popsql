package execution

import (
	"fmt"
	"time"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/plan"
)

type Row []any

type Executor struct {
	Store   kv.Store
	Catalog *catalog.Manager
}

type Result struct {
	Columns  []string
	Rows     []Row
	Duration time.Duration
}

func NewExecutor(st kv.Store, cat *catalog.Manager) *Executor {
	return &Executor{
		Store:   st,
		Catalog: cat,
	}
}

func (e *Executor) Execute(p plan.Plan) (*Result, error) {
	start := time.Now()
	result := func(cols []string, rows []Row) *Result {
		return &Result{
			Columns:  cols,
			Rows:     rows,
			Duration: time.Since(start),
		}
	}

	columns := []string{}
	rows, err := plan.VisitPlan(p, e)
	if err != nil {
		return nil, err
	}
	switch tp := p.(type) {
	case *plan.CreateTable:
		columns = []string{"table"}
		rows = []Row{{tp.Table.Name()}}
	case *plan.Insert:
		columns = []string{"count"}
		rows = []Row{{len(rows)}}
	case *plan.Scan:
		columns = tp.Columns()
	default:
		return nil, fmt.Errorf("unable to execute plan of type %T", tp)
	}
	return result(columns, rows), nil
}
