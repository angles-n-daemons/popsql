package execution

import "github.com/angles-n-daemons/popsql/pkg/db/sql/plan"

func (e *Executor) VisitValues(p *plan.Values) (Row, error) {
	offset := e.State.valueOffset[p.ID]
	if offset >= len(p.Rows) {
		return nil, nil
	}
	row := Row{}
	for _, expr := range p.Rows[offset] {
		val, err := Eval(e, expr)
		if err != nil {
			return nil, err
		}
		row = append(row, val)
	}
	e.State.valueOffset[p.ID]++
	return row, nil
}
