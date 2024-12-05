package schema_test

import (
	"fmt"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
	"github.com/angles-n-daemons/popsql/pkg/testing/assert"
)

func testTable() *schema.Table {
	return testTableFromArgs("", nil, nil)
}

func testTableFromArgs(name string, columns []*schema.Column, pkey []string) *schema.Table {
	if name == "" {
		name = "mytable"
	}
	if columns == nil {
		a, err := schema.NewColumn("a", scanner.DATATYPE_NUMBER)
		if err != nil {
			panic(err)
		}
		b, err := schema.NewColumn("b", scanner.DATATYPE_STRING)
		if err != nil {
			panic(err)
		}
		columns = []*schema.Column{a, b}
	}
	if pkey == nil {
		pkey = []string{"a"}
	}
	table, err := schema.NewTable(
		name,
		columns,
		pkey,
	)
	if err != nil {
		panic(err)
	}
	return table
}

func TestNewTable(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		a, err := schema.NewColumn("a", scanner.DATATYPE_NUMBER)
		if err != nil {
			t.Fatal(err)
		}
		b, err := schema.NewColumn("b", scanner.DATATYPE_STRING)
		if err != nil {
			t.Fatal(err)
		}
		table, err := schema.NewTable(
			"mytable",
			[]*schema.Column{a, b},
			[]string{"a"},
		)

		assert.Equal(t, table.Name, "mytable")
		assert.Equal(t, table.Columns, []*schema.Column{a, b})
		assert.Equal(t, table.PrimaryKey, []string{"a"})
	})

	t.Run("missing primary key", func(t *testing.T) {
		a, err := schema.NewColumn("a", scanner.DATATYPE_NUMBER)
		if err != nil {
			t.Fatal(err)
		}
		for _, test := range [][]string{nil, {}} {
			table, err := schema.NewTable(
				"mytable",
				[]*schema.Column{a},
				test,
			)
			assert.NoError(t, err)
			assert.Equal(t, len(table.Columns), 2)
			assert.Equal(t, len(table.PrimaryKey), 1)
			assert.Equal(t, table.PrimaryKey[0], schema.ReservedInternalKeyName)
		}
	})

	t.Run("invalid primary key", func(t *testing.T) {
		a, err := schema.NewColumn("a", scanner.DATATYPE_NUMBER)
		if err != nil {
			t.Fatal(err)
		}
		_, err = schema.NewTable(
			"mytable",
			[]*schema.Column{a},
			[]string{"b"},
		)
		assert.IsError(t, err, "could not find key column 'b' while creating table 'mytable'")
	})
}

func TestTableAddColumn(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		table := testTable()
		fmt.Println(table)
		err := table.AddColumn("c", scanner.DATATYPE_BOOLEAN)
		if err != nil {
			t.Fatal(err)
		}
		expected, err := schema.NewColumn("c", scanner.DATATYPE_BOOLEAN)
		assert.NoError(t, err)
		assert.Equal(t, len(table.Columns), 3)
		assert.Equal(t, table.Columns[2], expected)
	})
	t.Run("duplicate column", func(t *testing.T) {
		table := testTable()
		err := table.AddColumn("b", scanner.DATATYPE_BOOLEAN)
		assert.IsError(t, err, "a column with the name 'b' already exists on table 'mytable'")
	})
	t.Run("NewColumn fails", func(t *testing.T) {
		table := testTable()
		err := table.AddColumn("c", scanner.BANG)
		assert.IsError(t, err, "unrecognized data type BANG")
	})
}

func TestTableGetColumn(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		table := testTable()
		column := table.GetColumn("a")
		expected, err := schema.NewColumn("a", scanner.DATATYPE_NUMBER)
		assert.NoError(t, err)
		assert.Equal(t, expected, column)
	})
	t.Run("column doesn't exist", func(t *testing.T) {
		table := testTable()
		column := table.GetColumn("c")
		assert.Nil(t, column)
	})
}

func TestTableEqual(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		table1 := testTable()
		table2 := testTable()
		assert.Equal(t, table1, table2)
	})
	t.Run("other is nil", func(t *testing.T) {
		table1 := testTable()
		var table2 *schema.Table
		assert.NotEqual(t, table1, table2)
	})
	t.Run("primary keys not equal", func(t *testing.T) {
		table1 := testTable()
		table2 := testTable()
		table2.PrimaryKey = []string{"b"}
		assert.NotEqual(t, table1, table2)
	})
	t.Run("t has more columns", func(t *testing.T) {
		table1 := testTable()
		table2 := testTable()
		newCol, err := schema.NewColumn("third", scanner.DATATYPE_STRING)
		assert.NoError(t, err)
		table1.Columns = append(table1.Columns, newCol)
		assert.NotEqual(t, table1, table2)
	})
	t.Run("other has more columns", func(t *testing.T) {
		table1 := testTable()
		table2 := testTable()
		newCol, err := schema.NewColumn("third", scanner.DATATYPE_STRING)
		assert.NoError(t, err)
		table2.Columns = append(table2.Columns, newCol)
		assert.NotEqual(t, table1, table2)
	})
	t.Run("columns are different", func(t *testing.T) {
		a, err := schema.NewColumn("a", scanner.DATATYPE_NUMBER)
		assert.NoError(t, err)
		b, err := schema.NewColumn("b", scanner.DATATYPE_NUMBER)
		assert.NoError(t, err)
		c, err := schema.NewColumn("c", scanner.DATATYPE_NUMBER)
		assert.NoError(t, err)

		table1 := testTableFromArgs("", []*schema.Column{a, b}, []string{"a"})
		table2 := testTableFromArgs("", []*schema.Column{a, c}, []string{"a"})
		assert.NotEqual(t, table1, table2)
	})
}

func TestTableSerialization(t *testing.T) {
	bytes, err := testTable().Value()
	if err != nil {
		t.Fatal(err)
	}
	table, err := schema.NewTableFromBytes(bytes)
	if !testTable().Equal(table) {
		t.Fatal("expected table to be equal to serialized and deserialized copy")
	}
}

func TestTablePrefix(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected string
	}{
		{"chuck", "chuck/"},
		{"jim", "jim/"},
	} {
		t.Run(fmt.Sprintf("name=%s, expected=%s", test.name, test.expected), func(t *testing.T) {
			prefix := testTableFromArgs(test.name, nil, nil).Prefix()
			assert.Equal(t, prefix.String(), test.expected)
		})
	}
}

func TestTablePrefixEnd(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected string
	}{
		{"chuck", "chuck/<END>"},
		{"jim", "jim/<END>"},
	} {
		t.Run(fmt.Sprintf("name=%s, expected=%s", test.name, test.expected), func(t *testing.T) {
			prefixEnd := testTableFromArgs(test.name, nil, nil).PrefixEnd()
			assert.Equal(t, prefixEnd.String(), test.expected)
		})
	}
}

func TestTableKey(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected string
	}{
		{"chuck", "chuck"},
		{"jim", "jim"},
	} {
		t.Run(fmt.Sprintf("name=%s, expected=%s", test.name, test.expected), func(t *testing.T) {
			id := testTableFromArgs(test.name, nil, nil).Key()
			assert.Equal(t, id, test.expected)
		})
	}
}
