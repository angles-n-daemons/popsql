package execution

import (
	"encoding/json"

	"github.com/angles-n-daemons/popsql/pkg/sql/plan"
)

func (e *Executor) VisitScan(p *plan.Scan) (Row, error) {
	cur := e.State.cursors[p.ID]
	rowBytes, err := cur.Next()
	if err != nil {
		return nil, err
	}
	if rowBytes == nil {
		return nil, nil
	}

	rowMap := map[string]any{}
	err = json.Unmarshal(rowBytes, &rowMap)
	if err != nil {
		return nil, err
	}

	row := Row{}
	for _, col := range p.Table.GetColumns() {
		// condition for missing value?
		row = append(row, rowMap[col.Name])
	}

	return row, nil
}
