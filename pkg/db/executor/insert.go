package executor

import "github.com/angles-n-daemons/popsql/pkg/sql/plan"

func (e *Executor) VisitInsert(p *plan.Insert) ([]Row, error) {
	return nil, nil
}
