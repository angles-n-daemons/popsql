package schema

import (
	"encoding/json"
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
)

func (s *Schema) LoadTables(tablesBytes [][]byte) error {
	for _, tableBytes := range tablesBytes {
		var table *desc.Table
		err := json.Unmarshal(tableBytes, &table)
		if err != nil {
			return err
		}

		err = s.AddTable(table)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) AddTable(t *desc.Table) error {
	key := t.Name
	if _, ok := s.Tables[key]; ok {
		return fmt.Errorf("table '%s' already exists", t.Name)
	}
	s.Tables[key] = t
	return nil
}

func (s *Schema) GetTable(key string) (*desc.Table, bool) {
	table, ok := s.Tables[key]
	return table, ok
}

// RemoveTable attempts to drop the table with the given key.
// If the table does not exist, it returns an error.
func (s *Schema) RemoveTable(key string) error {
	_, ok := s.Tables[key]
	if !ok {
		return fmt.Errorf("could not delete table '%s'", key)
	}
	delete(s.Tables, key)
	return nil
}
