package execution

import (
	"time"

	"github.com/angles-n-daemons/popsql/pkg/db/kv"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/plan"
)

type Row []any

type Executor struct {
	Store   kv.Store
	Catalog *catalog.Manager
	State   *State
}

type Result struct {
	Columns  []string
	Rows     []Row
	Duration time.Duration
	Error    error
}

func Run(st kv.Store, cat *catalog.Manager, p plan.Plan) (*Result, error) {
	start := time.Now()
	result := func(cols []string, rows []Row) *Result {
		return &Result{
			Columns:  cols,
			Rows:     rows,
			Duration: time.Since(start),
		}
	}

	state, err := NewState(st, p)
	if err != nil {
		return nil, err
	}
	ex := &Executor{
		Store:   st,
		Catalog: cat,
		State:   state,
	}

	columns := p.Columns()
	rows := []Row{}
	for {
		row, err := Next(ex, p)
		if err != nil {
			return nil, err
		}
		if row == nil {
			break
		}
		rows = append(rows, row)
	}
	return result(columns, rows), nil
}

// Next executes the plan until the next resulting row is produced.
// It can be called on any plan node, and is used for recursively
// traversing the plan tree.
func Next(e *Executor, p plan.Plan) (Row, error) {
	return plan.VisitPlan(p, e)
}
