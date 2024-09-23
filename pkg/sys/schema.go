package sys

// sys.runtime?
// sys.schema?

// runtime.Schema
// runtime.Table
// runtime.Column

// runtime.Table.SERIAL

// schema.NewSchema()
// schema.Schema
// schema.Catalog

// sys.schema
//  - Schema
//  - Table
//    - KeyPrefix
//    - KeyPrefixEnd
//    - GetColumn
//    - Key
//  - Column

// sys.runtime
//  - Record { table, column, payload }
//    - Validate()

import (
	"fmt"
	"strings"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

const USER = "user"
const SYSTEM = "system"

func NewSchema() *Schema {
	return &Schema{
		System: SystemTables{
			Tables: Table{
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
			},
			Columns: Table{
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
			},
		},
		Tables: map[string]Table{},
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

type Table struct {
	Space      string
	Name       string
	Columns    []Column `json:"-"`
	PrimaryKey []string
}

func (t *Table) KeyPrefix() string {
	return fmt.Sprintf("%s/table/%s", t.Space, t.Name)
}

func (t *Table) KeyPrefixEnd() string {
	prefix := t.KeyPrefix()
	return nextString(prefix)
}

func (t *Table) Key() string {
	return fmt.Sprintf("%s:%s", t.Space, t.Name)
}

func nextString(s string) string {
	i := len(s) - 1
	for i >= 0 && s[i] == 'z' {
		i--
	}

	if i == -1 {
		return s + "a"
	}

	j := 0
	return strings.Map(func(r rune) rune {
		if j == i {
			r += 1
		}
		j++
		return r
	}, s)
}
func (t *Table) GetColumn(name string) (*Column, error) {
	for _, column := range t.Columns {
		if column.Name == name {
			return &column, nil
		}
	}
	return nil, fmt.Errorf("Unable to find column %s on table %s", name, t.Name)
}

func (t *Table) GetColumnFromRef(ref *ast.Reference) (*Column, error) {
	if len(ref.Names) == 0 {
		return nil, fmt.Errorf("Reference name too short: %d", len(ref.Names))
	}
	if len(ref.Names) > 1 {
		return nil, fmt.Errorf("Reference name too long: %d", len(ref.Names))
	}
	name := ref.Names[0].Lexeme
	return t.GetColumn(name)
}

type Column struct {
	Space    string
	Table    string
	Name     string
	DataType DataType
	// define an order?
}

func (c *Column) TableKeyPrefix() string {
	return fmt.Sprintf("%s/table/%s", c.Space, c.Table)
}

func (c *Column) Key() string {
	return fmt.Sprintf("%s:%s:%s", c.Space, c.Table, c.Name)
}

type DataType int

const (
	UNKNOWN = iota
	STRING
	NUMBER
	BOOLEAN
)

func GetDataType(token scanner.Token) (DataType, error) {
	switch token.Type {
	case scanner.DATATYPE_BOOLEAN:
		return BOOLEAN, nil
	case scanner.DATATYPE_STRING:
		return STRING, nil
	case scanner.DATATYPE_NUMBER:
		return NUMBER, nil
	default:
		return UNKNOWN, fmt.Errorf("unrecognized data type %s", token.Type)
	}
}

type Register interface {
	Key() string
}

type Row []any
