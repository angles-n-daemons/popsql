package schema

import (
	"encoding/json"
	"fmt"
)

const SchemaTableName = "__schema__"

type Schema struct {
	Tables map[string]*Table
}

func NewSchema() *Schema {
	schema := &Schema{
		Tables: map[string]*Table{},
	}
	return schema
}

func SchemaFromBytes(tablesBytes [][]byte) (*Schema, error) {
	var schema = &Schema{}

	for _, tableBytes := range tablesBytes {
		var table *Table
		err := json.Unmarshal(tableBytes, &table)
		if err != nil {
			return nil, err
		}
		schema.AddTable(table)
	}

	return schema, nil
}

func (s *Schema) AddTable(t *Table) error {
	key := t.Key()
	if _, ok := s.Tables[key]; ok {
		return fmt.Errorf("table '%s' already exists", t.Name)
	}
	return nil
}

func (s *Schema) GetTable(key string) (*Table, error) {
	table, ok := s.Tables[key]
	if !ok {
		return nil, fmt.Errorf("could not find table '%s'", key)
	}
	return table, nil
}

// Drop table attempts to drop the table with the given key.
// If the table does not exist, it returns an error.
func (s *Schema) DropTable(key string) error {
	_, ok := s.Tables[key]
	if !ok {
		return fmt.Errorf("could not delete table '%s'", key)
	}
	delete(s.Tables, key)
	return nil
}

// Equal is a simple comparator which says whether two schema references
// are logically equivalent. It does this by checking whether the references'
// internal maps are equivalent in size, and whether the values for each of
// their keys are equivalent.
func (s *Schema) Equal(other *Schema) bool {
	// if only one is nil, they cannot be equivalent
	if other == nil {
		return false
	}
	// if their internal maps are different sizes, they are not equivalent
	if len(s.Tables) != len(other.Tables) {
		return false
	}
	for key, table := range s.Tables {
		if !table.Equal(other.Tables[key]) {
			return false
		}
	}
	return true
}
