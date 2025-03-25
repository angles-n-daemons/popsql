package desc

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

var ReservedInternalKeyName = "__key"

var ErrNilColumns = errors.New("nil columns passed into NewTable")

type Table struct {
	ID         uint64
	Name       string
	Columns    []*Column
	PrimaryKey []string
}

func NewTable(id uint64, name string, columns []*Column, pkey []string) (*Table, error) {
	if columns == nil {
		return nil, ErrNilColumns
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
		ID:         id,
		Name:       name,
		Columns:    columns,
		PrimaryKey: pkey,
	}, nil
}

func NewTableFromBytes(tableBytes []byte) (*Table, error) {
	var table *Table
	err := json.Unmarshal(tableBytes, &table)
	return table, err
}

func (t *Table) AddColumn(name string, tokenType scanner.TokenType) error {
	if t.GetColumn(name) != nil {
		return fmt.Errorf(
			"a column with the name '%s' already exists on table '%s'",
			name,
			t.Name,
		)
	}

	column, err := NewColumn(name, tokenType)
	if err != nil {
		return err
	}

	t.Columns = append(t.Columns, column)
	return nil
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
	if t.ID != o.ID {
		return false
	}
	if t.Name != o.Name {
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
	return keys.New(t.Name)
}

func (t *Table) PrefixEnd() *keys.Key {
	return t.Prefix().Next()
}

// Utility functions for the desc table
func (t *Table) Key() string {
	return strconv.FormatUint(t.ID, 10)
}

func (t *Table) Value() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Table) AddInternalPrimaryKey() (*Sequence, error) {
	s := NewSequence(t.DefaultSequenceName())
	if len(t.PrimaryKey) != 0 {
		return nil, fmt.Errorf("table '%s' already has a primary key", t.Name)
	}
	t.Columns = append(t.Columns, NewSequenceColumn(ReservedInternalKeyName, NUMBER, s.Name))
	t.PrimaryKey = []string{ReservedInternalKeyName}
	return s, nil
}

func (t *Table) DefaultSequenceName() string {
	return fmt.Sprintf("%s_sequence", t.Name)
}
