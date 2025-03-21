package schema_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
	"github.com/angles-n-daemons/popsql/pkg/testing/assert"
)

func TestNewColumn(t *testing.T) {
	column, err := schema.NewColumn("name", scanner.DATATYPE_STRING)
	if err != nil {
		t.Fatal(err)
	}
	expected := &schema.Column{
		Name:     "name",
		DataType: schema.STRING,
	}
	assert.Equal(t, expected, column)
}

func TestNewColumnWrongDataType(t *testing.T) {
	column, err := schema.NewColumn("test", scanner.STRING)
	assert.Nil(t, column)
	assert.IsError(t, err, "unrecognized data type STRING")
}

func TestColumnEqual(t *testing.T) {
	// check nil condition
	column, err := schema.NewColumn("name", scanner.DATATYPE_STRING)
	if err != nil {
		t.Fatal(err)
	}

	if column.Equal(nil) {
		t.Fatal("compared non-nil column to nil and comparison was equal")
	}

	// check non-nil
	expected := &schema.Column{
		Name:     "name",
		DataType: schema.STRING,
	}
	if !column.Equal(expected) {
		t.Fatalf("expected %v and %v to be equal", column, expected)
	}
}
