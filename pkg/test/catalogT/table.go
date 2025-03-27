package catalogT

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

// Global counter to ensure each test-created table has a unique ID and (if needed) a unique name.
var tableIDCounter uint64

func TableID() uint64 {
	sequenceIDCounter++
	return sequenceIDCounter
}

func Table() *desc.Table {
	return NewTable(nil)
}

func TableWithID(id uint64) *desc.Table {
	return NewTable(&desc.Table{ID: id})
}

func TableWithName(name string) *desc.Table {
	return NewTable(&desc.Table{Name: name})
}

// Testing utility, which takes any portional part of a table and fills it out.
func NewTable(t *desc.Table) *desc.Table {
	if t == nil {
		t = &desc.Table{}
	}

	if t.ID == 0 {
		t.ID = SequenceID()
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
	tn, err := desc.NewTableWithID(t.ID, t.Name, t.Columns, t.PrimaryKey)
	if err != nil {
		panic(err)
	}
	return tn
}

func ReadTable(t *testing.T, st kv.Store, key string) *desc.Table {
	tableBytes, err := st.Get(key)
	if err != nil {
		t.Fatal(t)
	}
	var tb *desc.Table
	err = json.Unmarshal(tableBytes, &tb)
	if err != nil {
		t.Fatal(t)
	}
	return tb
}
