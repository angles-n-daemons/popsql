package schema_test

import (
	"fmt"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
	"github.com/angles-n-daemons/popsql/pkg/test/schematest"
)

// Global counter to ensure each test-created table has a unique ID and (if needed) a unique name.
var tableIDCounter uint64

// Test table creates a simple table with two columns, a a number and b a string.
func testTableFromArgs(name string, columns []*schema.Column, pkey []string) *schema.Table {
	tableIDCounter++

	if name == "" {
		// Give each table a unique name if none was provided
		name = fmt.Sprintf("mytable%d", tableIDCounter)
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
		tableIDCounter, // unique ID for each created table
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
		tableIDCounter++
		table, err := schema.NewTable(
			tableIDCounter,
			"mytable",
			[]*schema.Column{a, b},
			[]string{"a"},
		)
		assert.NoError(t, err)
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
			tableIDCounter++
			table, err := schema.NewTable(
				tableIDCounter,
				"mytable_missing_pk",
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
		tableIDCounter++
		_, err = schema.NewTable(
			tableIDCounter,
			"mytable_invalid_pk",
			[]*schema.Column{a},
			[]string{"b"},
		)
		assert.IsError(t, err, "could not find key column 'b' while creating table 'mytable_invalid_pk'")
	})
}

func TestTableAddColumn(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		table := schematest.TestTable()
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
		table := schematest.TestTable()
		err := table.AddColumn("b", scanner.DATATYPE_BOOLEAN)
		assert.IsError(t, err, fmt.Sprintf("a column with the name 'b' already exists on table '%s'", table.Name))
	})

	t.Run("NewColumn fails", func(t *testing.T) {
		table := schematest.TestTable()
		err := table.AddColumn("c", scanner.BANG)
		assert.IsError(t, err, "unrecognized data type BANG")
	})
}

func TestTableGetColumn(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		table := schematest.TestTable()
		column := table.GetColumn("a")
		expected, err := schema.NewColumn("a", scanner.DATATYPE_NUMBER)
		assert.NoError(t, err)
		assert.Equal(t, expected, column)
	})
	t.Run("column doesn't exist", func(t *testing.T) {
		table := schematest.TestTable()
		column := table.GetColumn("c")
		assert.Nil(t, column)
	})
}

func TestTableEqual(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		table1 := schematest.TestTable()
		table2 := schematest.CopyTable(table1)
		assert.Equal(t, table1, table2)
	})

	t.Run("ids are different", func(t *testing.T) {
		table1 := schematest.TestTable()
		table2 := schematest.CopyTable(table1)
		table2.ID = table1.ID + 1
		assert.NotEqual(t, table1, table2)
	})

	t.Run("names are different", func(t *testing.T) {
		table1 := schematest.TestTable()
		table2 := schematest.CopyTable(table1)
		table2.Name = "different"
		assert.NotEqual(t, table1, table2)
	})

	t.Run("other is nil", func(t *testing.T) {
		table1 := schematest.TestTable()
		var table2 *schema.Table
		assert.NotEqual(t, table1, table2)
	})

	t.Run("primary keys not equal", func(t *testing.T) {
		table1 := schematest.TestTable()
		table2 := schematest.CopyTable(table1)
		table2.PrimaryKey = []string{"b"}
		assert.NotEqual(t, table1, table2)
	})

	t.Run("t has more columns", func(t *testing.T) {
		table1 := schematest.TestTable()
		table2 := schematest.CopyTable(table1)
		newCol, err := schema.NewColumn("third", scanner.DATATYPE_STRING)
		assert.NoError(t, err)
		table1.Columns = append(table1.Columns, newCol)
		assert.NotEqual(t, table1, table2)
	})

	t.Run("other has more columns", func(t *testing.T) {
		table1 := schematest.TestTable()
		table2 := schematest.CopyTable(table1)
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

		table1 := schematest.NewTable(&schema.Table{ID: 1, Columns: []*schema.Column{a, b}, PrimaryKey: []string{"a"}})
		table2 := schematest.CopyTable(table1)
		table2.Columns = []*schema.Column{a, c}
		assert.NotEqual(t, table1, table2)
	})
}

func TestTableSerialization(t *testing.T) {
	original := schematest.TestTable()
	bytes, err := original.Value()
	assert.NoError(t, err)

	table, err := schema.NewTableFromBytes(bytes)
	assert.NoError(t, err)
	assert.Equal(t, original, table)
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
			table := testTableFromArgs(test.name, nil, nil)
			prefix := table.Prefix()
			assert.Equal(t, prefix.Encode(), test.expected)
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
			table := testTableFromArgs(test.name, nil, nil)
			prefixEnd := table.PrefixEnd()
			assert.Equal(t, prefixEnd.Encode(), test.expected)
		})
	}
}

func TestTableKey(t *testing.T) {
	table := &schema.Table{
		ID: 123,
	}
	assert.Equal(t, table.Key(), "123")
}
