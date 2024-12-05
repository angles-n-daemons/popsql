package schema

import "fmt"

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
		table, err := NewTableFromBytes(tableBytes)
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

func (s *Schema) DropTable(key string) error {
	if key == RootTable.Key() {
		return ErrDropMetaTable
	}
	_, ok := s.Tables[key]
	if !ok {
		return fmt.Errorf("could not delete table '%s'", key)
	}
	delete(s.Tables, key)
	return nil
}
