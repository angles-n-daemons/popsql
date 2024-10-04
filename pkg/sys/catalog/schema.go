package catalog

const USER = "user"
const SYSTEM = "system"

func NewSchema() *Schema {
	return &Schema{
		System: SystemTables{
			Tables:  TABLE_TABLE,
			Columns: COLUMN_TABLE,
		},
	}
}

type Schema struct {
	System SystemTables
	Tables map[string]Table
}

type SystemTables struct {
	Tables  Table
	Columns Table
}

var TABLE_TABLE = Table{
	Space: SYSTEM,
	Name:  "tables",
	Columns: []Column{
		{
			Space:    SYSTEM,
			Table:    "tables",
			Name:     "scope",
			DataType: STRING,
		},
		{
			Space:    SYSTEM,
			Table:    "tables",
			Name:     "name",
			DataType: STRING,
		},
	},
	PrimaryKey: []string{"scope", "table"},
}

var COLUMN_TABLE = Table{
	Space: SYSTEM,
	Name:  "columns",
	Columns: []Column{
		{
			Space:    SYSTEM,
			Table:    "columns",
			Name:     "scope",
			DataType: STRING,
		},
		{
			Space:    SYSTEM,
			Table:    "columns",
			Name:     "table",
			DataType: STRING,
		},
		{
			Space:    SYSTEM,
			Table:    "columns",
			Name:     "name",
			DataType: STRING,
		},
		{
			Space:    SYSTEM,
			Table:    "columns",
			Name:     "datatype",
			DataType: STRING,
		},
	},
	PrimaryKey: []string{"scope", "table", "name"},
}
