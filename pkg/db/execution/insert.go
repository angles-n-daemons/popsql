package execution

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/sys"
	"github.com/angles-n-daemons/popsql/pkg/sql/plan"
)

// return the rows of the new ids
func (e *Executor) VisitInsert(p *plan.Insert) ([]Row, error) {
	keys := []*keys.Key{}
	rows := []map[string]any{}
	result := []Row{}

	// convert input to serializable data
	for _, tup := range p.Values {
		row := map[string]any{}
		rr := Row{}
		key := p.Table.Prefix()

		// I need to do this as well if the primary key isn't internal, but still shows up
		if primaryKeyInternal(p.Table) {
			id, err := catalog.TableSequenceNext(e.Catalog, p.Table)
			if err != nil {
				return nil, err
			}
			sid := strconv.FormatUint(id, 10)
			key = key.WithID(sid)
		}

		// TODO: validate types here
		for i, expr := range tup {
			val, err := Eval(expr)
			if err != nil {
				return nil, err
			}

			row[p.Columns[i].Name] = val
			rr = append(rr, val)

			if slices.Contains(p.Table.PrimaryKey, p.Columns[i].Name) {
				key = key.WithIDAddition(fmt.Sprintf("%v", val))
			}
		}
		rows = append(rows, row)
		result = append(result, rr)
		keys = append(keys, key)
	}

	// save data in the db
	for i, row := range rows {
		b, err := json.Marshal(row)
		if err != nil {
			return nil, err
		}
		e.Store.Put(keys[i].Encode(), b)
	}
	return result, nil
}

func primaryKeyInternal(t *desc.Table) bool {
	return len(t.PrimaryKey) == 1 && t.PrimaryKey[0] == sys.ReservedInternalKeyName
}
