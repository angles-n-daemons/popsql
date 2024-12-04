package schema_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
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

		if table.Name != "mytable" {
			t.Fatalf("expected table name '%s', got '%s'", "mytable", table.Name)
		}
		if !slices.Equal(table.Columns, []*schema.Column{a, b}) {
			t.Fatal("columns weren't equal")
		}
		if !slices.Equal(table.PrimaryKey, []string{"a"}) {
			t.Fatal("primary keys not equal")
		}
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
			if err != nil {
				t.Fatal(err)
			}

			if len(table.Columns) != 2 {
				t.Fatal("failed to add internal key column")
			}
			if len(table.PrimaryKey) != 1 {
				t.Fatal("failed to add or create primary key")
			}
			if table.PrimaryKey[0] != schema.ReservedInternalKeyName {
				t.Fatalf("unexpected primary key name '%s'", table.PrimaryKey[0])
			}
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
		if err == nil {
			t.Fatal("didn't error like expected")
		}
		if err.Error() != "could not find key column 'b' while creating table 'mytable'" {
			t.Fatalf("incorrect error '%s'", err)
		}
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
		if err != nil {
			t.Fatal(err)
		}
		if len(table.Columns) != 3 {
			t.Fatalf("columns length '%d' incorrect", len(table.Columns))
		}
		if !table.Columns[2].Equal(expected) {
			t.Fatal("column at expected index was different")
		}
	})
	t.Run("duplicate column", func(t *testing.T) {
		table := testTable()
		err := table.AddColumn("b", scanner.DATATYPE_BOOLEAN)
		if err == nil {
			t.Fatal("didn't error like expected")
		}
		if err.Error() != "a column with the name 'b' already exists on table 'mytable'" {
			t.Fatalf("wrong error '%s'", err)
		}
	})
	t.Run("NewColumn fails", func(t *testing.T) {
		table := testTable()
		err := table.AddColumn("c", scanner.BANG)
		if err == nil {
			t.Fatal("didn't error like expected")
		}
		if err.Error() != "unrecognized data type BANG" {
			t.Fatalf("wrong error '%s'", err)
		}
	})
}

func TestTableGetColumn(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		table := testTable()
		column := table.GetColumn("a")
		expected, err := schema.NewColumn("a", scanner.DATATYPE_NUMBER)
		if err != nil {
			t.Fatal(err)
		}
		if !expected.Equal(column) {
			t.Fatalf("columns not equal: %v, %v", expected, column)
		}
	})
	t.Run("column doesn't exist", func(t *testing.T) {
		table := testTable()
		column := table.GetColumn("c")
		if column != nil {
			t.Fatalf("found column '%s' when expecting nil response", column.Name)
		}
	})
}

func TestTableEqual(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		table1 := testTable()
		table2 := testTable()
		if table1 == table2 {
			t.Fatal("test conditions incorrect, using same addresses")
		}
		if !table1.Equal(table2) {
			t.Fatalf("tables not equal")
		}
	})
	t.Run("other is nil", func(t *testing.T) {
		table1 := testTable()
		var table2 *schema.Table
		if table1.Equal(table2) {
			t.Fatalf("expected tables to not be equal")
		}
	})
	t.Run("primary keys not equal", func(t *testing.T) {
		table1 := testTable()
		table2 := testTable()
		table2.PrimaryKey = []string{"b"}
		if table1.Equal(table2) {
			t.Fatalf("expected tables to not be equal")
		}
	})
	t.Run("t has more columns", func(t *testing.T) {
		table1 := testTable()
		table2 := testTable()
		newCol, err := schema.NewColumn("third", scanner.DATATYPE_STRING)
		if err != nil {
			t.Fatal(err)
		}
		table1.Columns = append(table1.Columns, newCol)
		if table1.Equal(table2) {
			t.Fatal("expected tables not to be equal")
		}
	})
	t.Run("other has more columns", func(t *testing.T) {
		table1 := testTable()
		table2 := testTable()
		newCol, err := schema.NewColumn("third", scanner.DATATYPE_STRING)
		if err != nil {
			t.Fatal(err)
		}
		table2.Columns = append(table2.Columns, newCol)
		if table1.Equal(table2) {
			t.Fatal("expected tables not to be equal")
		}
	})
	t.Run("columns are different", func(t *testing.T) {
		a, err := schema.NewColumn("a", scanner.DATATYPE_NUMBER)
		if err != nil {
			t.Fatal(err)
		}
		b, err := schema.NewColumn("b", scanner.DATATYPE_NUMBER)
		if err != nil {
			t.Fatal(err)
		}
		c, err := schema.NewColumn("c", scanner.DATATYPE_NUMBER)
		if err != nil {
			t.Fatal(err)
		}

		table1 := testTableFromArgs("", []*schema.Column{a, b}, []string{"a"})
		table2 := testTableFromArgs("", []*schema.Column{a, c}, []string{"a"})
		if table1.Equal(table2) {
			t.Fatal("expected tables not to be equal")
		}
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
			if prefix.String() != test.expected {
				t.Fatalf("expected %s but got %s", test.expected, prefix)
			}
		})
		fmt.Println(test)
	}
}

func TestTableID(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected string
	}{
		{"chuck", "chuck"},
		{"jim", "jim"},
	} {
		t.Run(fmt.Sprintf("name=%s, expected=%s", test.name, test.expected), func(t *testing.T) {
			id, err := testTableFromArgs(test.name, nil, nil).ID()
			if err != nil {
				t.Fatal(err)
			}
			if id != test.expected {
				t.Fatalf("expected %s but got %s", test.expected, id)
			}
		})
		fmt.Println(test)
	}
}
