package schematest

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
)

// Global counter to ensure each test-created table has a unique ID and (if needed) a unique name.
var tableIDCounter uint64

func TableID() uint64 {
	tableIDCounter++
	return tableIDCounter
}

func TestTable() *schema.Table {
	return NewTable(nil)
}

// Testing utility, which takes any portional part of a table and fills it out.
func NewTable(t *schema.Table) *schema.Table {
	if t == nil {
		t = &schema.Table{}
	}

	if t.ID == 0 {
		t.ID = TableID()
	}

	if t.Name == "" {
		t.Name = fmt.Sprintf("mytable%d", tableIDCounter)
	}

	if t.Columns == nil {
		a, err := schema.NewColumn("a", scanner.DATATYPE_NUMBER)
		if err != nil {
			panic(err)
		}
		b, err := schema.NewColumn("b", scanner.DATATYPE_STRING)
		if err != nil {
			panic(err)
		}
		t.Columns = []*schema.Column{a, b}
		t.PrimaryKey = []string{"a"}
	}
	return t
}

func CopyTable(t *schema.Table) *schema.Table {
	t, err := schema.NewTable(t.ID, t.Name, t.Columns, t.PrimaryKey)
	if err != nil {
		panic(err)
	}
	return t
}
