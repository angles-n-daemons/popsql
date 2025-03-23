package schematest

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema/desc"
)

// Global counter to ensure each test-created table has a unique ID and (if needed) a unique name.
var tableIDCounter uint64

func TableID() uint64 {
	tableIDCounter++
	return tableIDCounter
}

func TestTable() *desc.Table {
	return NewTable(nil)
}

func TableWithID(id uint64) *desc.Table {
	return NewTable(&desc.Table{ID: id})
}

// Testing utility, which takes any portional part of a table and fills it out.
func NewTable(t *desc.Table) *desc.Table {
	if t == nil {
		t = &desc.Table{}
	}

	if t.ID == 0 {
		t.ID = TableID()
	}

	if t.Name == "" {
		t.Name = fmt.Sprintf("table_%d", t.ID)
	}

	if t.Columns == nil {
		a, err := desc.NewColumn("a", scanner.DATATYPE_NUMBER)
		if err != nil {
			panic(err)
		}
		b, err := desc.NewColumn("b", scanner.DATATYPE_STRING)
		if err != nil {
			panic(err)
		}
		t.Columns = []*desc.Column{a, b}
		t.PrimaryKey = []string{"a"}
	}
	return t
}

func CopyTable(t *desc.Table) *desc.Table {
	t, err := desc.NewTable(t.ID, t.Name, t.Columns, t.PrimaryKey)
	if err != nil {
		panic(err)
	}
	return t
}
