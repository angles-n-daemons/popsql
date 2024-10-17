package schema

const USER = "user"
const SYSTEM = "system"

func NewSchema() *Schema {
	return &Schema{}
}

func SchemaFromBytes(tables [][]byte, columns [][]byte) (*Schema, error) {
	var schema = NewSchema()

	return schema, nil
}

func GetTable(space string, name string) (*Table, bool, error) {
	return nil, false, nil
}

type Schema struct {
	System *TableSpace
	User   *TableSpace
}

func (ts *Schema) AddTable(name string) {

}

func (s *Schema) GetTable(names []string) (*Table, error) {
	return nil, nil
}

type TableSpace struct {
	TableLookup   map[string]*Table
	TableNameToID map[string]string
}

func (ts *TableSpace) AddTable(name string) {

}

func (ts *TableSpace) GetTable(name string) {

}

var TABLE_TABLE = &Table{
	Space: SYSTEM,
	Name:  "tables",
	Columns: []*Column{
		{
			Name:     "scope",
			DataType: STRING,
		},
		{
			Name:     "name",
			DataType: STRING,
		},
	},
	PrimaryKey: []string{"scope", "table"},
}

var COLUMN_TABLE = &Table{
	Space: SYSTEM,
	Name:  "columns",
	Columns: []*Column{
		{
			Name:     "scope",
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
	PrimaryKey: []string{"scope", "table", "name"},
}
