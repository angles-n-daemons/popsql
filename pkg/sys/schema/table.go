package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

var ReservedInternalKeyName = "___zkey"

var ErrNilColumns = errors.New("nil columns passed into NewTable")

type Table struct {
	Name       string
	Columns    []*Column
	PrimaryKey []string
}

func NewTable(name string, columns []*Column, pkey []string) (*Table, error) {
	if columns == nil {
		return nil, ErrNilColumns
	}
	if len(pkey) == 0 {
		// if a primary key wasn't found, create an internal one
		pkey = []string{ReservedInternalKeyName}
		pkeyColumn, err := NewColumn(ReservedInternalKeyName, scanner.DATATYPE_STRING)
		if err != nil {
			return nil, err
		}
		columns = append(columns, pkeyColumn)
	} else {
		// otherwise, verify that the names used for the primary key exist
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
	}
	return &Table{
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

func (t *Table) Equal(other *Table) bool {
	if other == nil {
		return false
	}
	if len(t.Columns) != len(other.Columns) {
		return false
	}
	for i, column := range t.Columns {
		if !column.Equal(other.Columns[i]) {
			return false
		}
	}
	if !slices.Equal(t.PrimaryKey, other.PrimaryKey) {
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

// Utility functions for the schema table
func (t *Table) Key() string {
	return t.Name
}

func (t *Table) Value() ([]byte, error) {
	return json.Marshal(t)
}
