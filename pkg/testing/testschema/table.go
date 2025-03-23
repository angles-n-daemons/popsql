package testschema

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
)

// Global counter to ensure each test-created table has a unique ID and (if needed) a unique name.
var tableIDCounter uint64

func TestTable() *schema.Table {
	return NewTable(nil)
}

func NewTable(t *schema.Table) *schema.Table {
	if t == nil {
		t = &schema.Table{}
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
