package execution

import (
	"encoding/json"

	"github.com/angles-n-daemons/popsql/pkg/sql/plan"
)

func (e *Executor) VisitScan(p *plan.Scan) ([]Row, error) {
	span := p.Table.Span()
	cur, err := e.Store.Scan(span.Start.Encode(), span.End.Encode())
	if err != nil {
		return nil, err
	}

	rawRows, err := cur.ReadAll()
	if err != nil {
		return nil, err
	}

	rows := make([]Row, len(rawRows))
	for i, rowBytes := range rawRows {
		rowMap := map[string]any{}
		err := json.Unmarshal(rowBytes, &rowMap)
		if err != nil {
			return nil, err
		}

		row := Row{}
		for _, col := range p.Table.GetColumns() {

			// condition for missing value?
			row = append(row, rowMap[col.Name])
		}

		rows[i] = row
	}
	return rows, nil
}
