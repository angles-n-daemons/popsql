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

func InitSchema() *Schema {
	schema := NewSchema()
	schema.AddTable(Tables)
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
	id, err := t.ID()
	if err != nil {
		return err
	}
	if _, ok := s.Tables[id]; ok {
		return fmt.Errorf("table with name '%s' already exists", t.Name)
	}
	return nil
}

func (s *Schema) GetTable(id string) (*Table, error) {
	table, ok := s.Tables[id]
	if !ok {
		return nil, fmt.Errorf("could not find table with id '%s'", id)
	}
	return table, nil
}

var Tables = &Table{
	Name: SchemaTableName,
	Columns: []*Column{
		{
			Name:     "name",
			DataType: STRING,
		},
	},
	PrimaryKey: []string{"name"},
}
