package catalog

var SYSTEM = "system"
var USER = "user"

type Schema struct {
	Tables        map[string]*Table
	TableIdByName map[string]string
}

func (n *Schema) AddTable(t *Table) error {
	// error if table already exists
	return nil
}

func (n *Schema) GetTable(id string) (*Table, error) {
	return nil, nil
}

func (n *Schema) GetTableByName(id string) (*Table, error) {
	return nil, nil
}

func InitSchema() *Schema {
	schema := &Schema{
		Tables:        map[string]*Table{},
		TableIdByName: map[string]string{},
	}
	var tables = &Table{
		Name: "tables",
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

func SchemaFromBytes(tables [][]byte) (*Schema, error) {
	var schema = InitSchema()

	return schema, nil
}
