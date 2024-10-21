package catalog_test

import (
	"fmt"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sys/catalog"
)

func TestNewColumn(t *testing.T) {
	column, err := catalog.NewColumn("name", scanner.DATATYPE_STRING)
	if err != nil {
		t.Fatal(err)
	}
	expected := &catalog.Column{
		"name",
		catalog.STRING,
	}
	if !column.Equal(expected) {
		t.Fatalf("expected %v and %v to be equal", column, expected)
	}
}

func TestNewColumnWrongDataType(t *testing.T) {
	column, err := catalog.NewColumn("test", scanner.STRING)
	if column != nil {
		t.Fatalf("expected column to be nil, got %v", column)
	}
	if err == nil {
		t.Fatalf("expected NewColumn to fail")
	}
	if err.Error() != fmt.Sprintf("unrecognized data type STRING") {
		t.Fatalf("wrong error message %s", err)
	}
}

func TestColumnEqual(t *testing.T) {
	// check nil condition
	column, err := catalog.NewColumn("name", scanner.DATATYPE_STRING)
	if err != nil {
		t.Fatal(err)
	}

	if column.Equal(nil) {
		t.Fatal("compared non-nil column to nil and comparison was equal")
	}

	// check non-nil
	expected := &catalog.Column{
		"name",
		catalog.STRING,
	}
	if !column.Equal(expected) {
		t.Fatalf("expected %v and %v to be equal", column, expected)
	}
}
