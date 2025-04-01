package desc

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
)

type Table struct {
	TID        uint64
	TName      string
	Columns    []*Column
	PrimaryKey []string
}

func NewTable(name string, columns []*Column, pkey []string) (*Table, error) {
	return NewTableWithID(0, name, columns, pkey)
}

func NewTableWithID(id uint64, name string, columns []*Column, pkey []string) (*Table, error) {
	if columns == nil {
		columns = []*Column{}
	}
	if pkey == nil {
		pkey = []string{}
	}
	for _, key := range pkey {
		found := false
		for _, column := range columns {
			if column.Name == key {
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("could not find key column '%s' while creating table '%s'", key, name)
		}
	}
	return &Table{
		TID:        id,
		TName:      name,
		Columns:    columns,
		PrimaryKey: pkey,
	}, nil
}

func NewTableFromBytes(tableBytes []byte) (*Table, error) {
	var table *Table
	err := json.Unmarshal(tableBytes, &table)
	return table, err
}

func (t *Table) WithID(id uint64) {
	t.TID = id
}

func (t *Table) GetColumn(name string) *Column {
	for _, column := range t.Columns {
		if column.Name == name {
			return column
		}
	}
	return nil
}

func (t *Table) Equal(o *Table) bool {
	if o == nil {
		return false
	}
	if t.TID != o.TID {
		return false
	}
	if t.TName != o.TName {
		return false
	}
	if len(t.Columns) != len(o.Columns) {
		return false
	}
	for i, column := range t.Columns {
		if !column.Equal(o.Columns[i]) {
			return false
		}
	}
	if !slices.Equal(t.PrimaryKey, o.PrimaryKey) {
		return false
	}
	return true
}

func (t *Table) Prefix() *keys.Key {
	return keys.New(t.Key())
}

func (t *Table) Span() *keys.Span {
	p := t.Prefix()
	return &keys.Span{
		Start: p,
		End:   p.Next(),
	}
}

func (t *Table) Name() string {
	return t.TName
}

func (t *Table) ID() uint64 {
	return t.TID
}

// Utility functions for the desc table
func (t *Table) Key() string {
	return strconv.FormatUint(t.TID, 10)
}

func (t *Table) Value() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Table) DefaultSequenceName() string {
	return t.TName + "_seq"
}
