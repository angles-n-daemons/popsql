package sys

import (
	"fmt"
	"strings"

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
	Space   string
	Name    string
	Columns []Column `json:"-"`
}

func (t *Table) KeyPrefix() string {
	return fmt.Sprintf("%s/%s", t.Space, t.Name)
}
func (t *Table) KeyPrefixEnd() string {
	prefix := t.KeyPrefix()
	return nextString(prefix)
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

func (t *Table) Key() string {
	return t.Name
}

type Column struct {
	Space    string
	Table    string
	Name     string
	DataType DataType
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
