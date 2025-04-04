package executor

import "github.com/angles-n-daemons/popsql/pkg/sql/plan"

// return the rows of the new ids
func (e *Executor) VisitInsert(p *plan.Insert) ([]Row, error) {
	return nil, nil
}
