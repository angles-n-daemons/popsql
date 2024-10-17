package catalog

var SYSTEM = "system"
var USER = "user"

type Schema struct {
	System *Namespace
	User   *Namespace
}

func InitSchema() *Schema {
	var system = NewNamespace(SYSTEM)
	var user = NewNamespace(USER)

	var tables = &Table{
		Namespace: system,
		Name:      "tables",
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

	var columns = &Table{
		Namespace: system,
		Name:      "columns",
		Columns: []*Column{
			{
				Name:     "namespace",
				DataType: STRING,
			},
			{
				Name:     "table",
				DataType: STRING,
			},
			{
				Name:     "name",
				DataType: STRING,
			},
			{
				Name:     "datatype",
				DataType: STRING,
			},
		},
	}

	system.AddTable(tables)
	system.AddTable(columns)
	return &Schema{
		System: system,
		User:   user,
	}
}

func SchemaFromBytes(tables [][]byte, columns [][]byte) (*Schema, error) {
	var schema = InitSchema()

	return schema, nil
}
