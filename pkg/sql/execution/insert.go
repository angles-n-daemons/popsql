package execution

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
	"github.com/angles-n-daemons/popsql/pkg/sql/plan"
)

// return the rows of the new ids
func (e *Executor) VisitInsert(p *plan.Insert) (Row, error) {
	tup, err := Next(e, p.Source)
	if err != nil {
		return nil, err
	}
	if tup == nil {
		return nil, nil
	}

	data := map[string]any{}
	for i, val := range tup {
		data[p.Cols[i].Name] = val

	}

	key, err := e.getAndValidateKey(p.Table, data, tup)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	e.Store.Put(key.Encode(), b)
	return Row{key}, nil
}

// getAndValidate key has a few responsibilities, ultimately
// returning the key for the row, whether it existed in the data
// struct constructed by the executor, and any errors that may
// have occured during this process.
//
// 1. If the primary key for this table is ineternal, create it.
// 2. If the primary key fully exists within the data map, return it.
// 3. Otherwise, return an error.
func (e *Executor) getAndValidateKey(
	t *desc.Table, data map[string]any, tup Row,
) (*keys.Key, error) {
	key := t.Prefix()
	if primaryKeyInternal(t) {
		id, err := e.tableSequenceNext(t)
		if err != nil {
			return nil, err
		}
		// add the internal column to the data map
		data[desc.ReservedInternalColumnName] = id

		sid := strconv.FormatUint(id, 10)
		key = key.WithID(sid)
		return key, nil
	}

	for _, col := range t.PrimaryKey {
		v, ok := data[col]
		if !ok {
			return nil, fmt.Errorf("key column '%s' missing on insert of row '%v'", col, tup)
		}
		key = key.WithIDAddition(fmt.Sprintf("%v", v))
	}
	return key, nil
}

func primaryKeyInternal(t *desc.Table) bool {
	return len(t.PrimaryKey) == 1 && t.PrimaryKey[0] == desc.ReservedInternalColumnName
}

func (e *Executor) tableSequenceNext(t *desc.Table) (uint64, error) {
	seq := schema.GetByName[*desc.Sequence](e.Catalog.Schema, t.DefaultSequenceName())
	if seq == nil {
		return 0, fmt.Errorf("could not find key column for table %s", t.Name())
	}

	return catalog.SequenceNext(e.Catalog, seq)
}
