package catalog

import "fmt"

const SchemaTableName = "__schema__"

type Schema struct {
	Tables        map[string]*Table
	TableIdByName map[string]string
}

func (n *Schema) AddTable(t *Table) error {
	// error if table already exists
	return nil
}

func (n *Schema) GetTable(id string) (*Table, error) {
	// error if doesn't exist
	return nil, nil
}

func (n *Schema) GetTableByName(id string) (*Table, error) {
	// error if doesn't exist
	return nil, nil
}

func InitSchema() *Schema {
	schema := &Schema{
		Tables:        map[string]*Table{},
		TableIdByName: map[string]string{},
	}
	var tables = &Table{
		Name: SchemaTableName,
		Columns: []*Column{
			{
				Name:     "namespace",
				DataType: STRING,
			},
			{
				Name:     "name",
				DataType: STRING,
			},
		},
	}

	schema.AddTable(tables)
	return &Schema{}
}

func SchemaFromBytes(tablesBytes [][]byte) (*Schema, error) {
	var schema = InitSchema()

	for _, tableBytes := range tablesBytes {
		fmt.Println(tableBytes)
	}

	return schema, nil
}
