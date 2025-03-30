package desc_test

import (
	"fmt"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
	"github.com/angles-n-daemons/popsql/pkg/test/catalogT"
)

func TestNewTable(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		a := desc.NewColumn("a", desc.NUMBER)
		b := desc.NewColumn("b", desc.STRING)
		table, err := desc.NewTable("mytable", []*desc.Column{a, b}, []string{"a"})
		assert.NoError(t, err)
		assert.Equal(t, table.TName, "mytable")
		assert.Equal(t, table.Columns, []*desc.Column{a, b})
		assert.Equal(t, table.PrimaryKey, []string{"a"})
	})

	t.Run("invalid primary key", func(t *testing.T) {
		a := desc.NewColumn("a", desc.NUMBER)
		_, err := desc.NewTable("mytable_invalid_pk", []*desc.Column{a}, []string{"b"})
		assert.IsError(t, err, "could not find key column 'b' while creating table 'mytable_invalid_pk'")
	})
}

func TestTableGetColumn(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		table := catalogT.Table()
		column := table.GetColumn("a")
		expected := desc.NewColumn("a", desc.NUMBER)
		assert.Equal(t, expected, column)
	})
	t.Run("column doesn't exist", func(t *testing.T) {
		table := catalogT.Table()
		column := table.GetColumn("c")
		assert.Nil(t, column)
	})
}

func TestTableEqual(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		table1 := catalogT.Table()
		table2 := catalogT.CopyTable(table1)
		assert.Equal(t, table1, table2)
	})

	t.Run("ids are different", func(t *testing.T) {
		table1 := catalogT.Table()
		table2 := catalogT.CopyTable(table1)
		table2.TID = table1.TID + 1
		assert.NotEqual(t, table1, table2)
	})

	t.Run("names are different", func(t *testing.T) {
		table1 := catalogT.Table()
		table2 := catalogT.CopyTable(table1)
		table2.TName = "different"
		assert.NotEqual(t, table1, table2)
	})

	t.Run("other is nil", func(t *testing.T) {
		table1 := catalogT.Table()
		var table2 *desc.Table
		assert.NotEqual(t, table1, table2)
	})

	t.Run("primary keys not equal", func(t *testing.T) {
		table1 := catalogT.Table()
		table2 := catalogT.CopyTable(table1)
		table2.PrimaryKey = []string{"b"}
		assert.NotEqual(t, table1, table2)
	})

	t.Run("t has more columns", func(t *testing.T) {
		table1 := catalogT.Table()
		table2 := catalogT.CopyTable(table1)
		newCol := desc.NewColumn("third", desc.STRING)
		table1.Columns = append(table1.Columns, newCol)
		assert.NotEqual(t, table1, table2)
	})

	t.Run("other has more columns", func(t *testing.T) {
		table1 := catalogT.Table()
		table2 := catalogT.CopyTable(table1)
		newCol := desc.NewColumn("third", desc.STRING)
		table2.Columns = append(table2.Columns, newCol)
		assert.NotEqual(t, table1, table2)
	})

	t.Run("columns are different", func(t *testing.T) {
		a := desc.NewColumn("a", desc.NUMBER)
		b := desc.NewColumn("b", desc.NUMBER)
		c := desc.NewColumn("c", desc.NUMBER)

		table1 := catalogT.NewTable(&desc.Table{TID: 1, Columns: []*desc.Column{a, b}, PrimaryKey: []string{"a"}})
		table2 := catalogT.CopyTable(table1)
		table2.Columns = []*desc.Column{a, c}
		assert.NotEqual(t, table1, table2)
	})
}

func TestTableSerialization(t *testing.T) {
	original := catalogT.Table()
	bytes, err := original.Value()
	assert.NoError(t, err)

	table, err := desc.NewTableFromBytes(bytes)
	assert.NoError(t, err)
	assert.Equal(t, original, table)
}

func TestTablePrefix(t *testing.T) {
	for _, test := range []struct {
		id       uint64
		expected string
	}{
		{1, "1/"},
		{2, "2/"},
	} {
		t.Run(fmt.Sprintf("id=%d, expected=%s", test.id, test.expected), func(t *testing.T) {
			table := catalogT.TableWithID(test.id)
			prefix := table.Prefix()
			assert.Equal(t, prefix.Encode(), test.expected)
		})
	}
}

func TestTableKey(t *testing.T) {
	table := &desc.Table{
		TID: 123,
	}
	assert.Equal(t, table.Key(), "123")
}
