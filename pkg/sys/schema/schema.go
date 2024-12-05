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
	id := t.Key()
	if _, ok := s.Tables[id]; ok {
		return fmt.Errorf("table '%s' already exists", t.Name)
	}
	return nil
}

func (s *Schema) GetTable(id string) (*Table, error) {
	table, ok := s.Tables[id]
	if !ok {
		return nil, fmt.Errorf("could not find table with key '%s'", id)
	}
	return table, nil
}

func (s *Schema) DropTable(id string) error {
	_, ok := s.Tables[id]
	if !ok {
		return fmt.Errorf("could not delete table with key '%s'", id)
	}
	delete(s.Tables, id)
	return nil
}

var RootTable = &Table{
	Name: SchemaTableName,
	Columns: []*Column{
		{
			Name:     "name",
			DataType: STRING,
		},
	},
	PrimaryKey: []string{"name"},
}
