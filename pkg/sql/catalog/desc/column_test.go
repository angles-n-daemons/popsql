package desc_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
)

func TestNewColumn(t *testing.T) {
	column, err := desc.SequenceColumn("name", scanner.DATATYPE_STRING)
	if err != nil {
		t.Fatal(err)
	}
	expected := &desc.Column{
		Name:     "name",
		DataType: desc.STRING,
	}
	assert.Equal(t, expected, column)
}

func TestNewColumnWrongDataType(t *testing.T) {
	column, err := desc.SequenceColumn("test", scanner.STRING)
	assert.Nil(t, column)
	assert.IsError(t, err, "unrecognized data type STRING")
}

func TestColumnEqual(t *testing.T) {
	// check nil condition
	column, err := desc.SequenceColumn("name", scanner.DATATYPE_STRING)
	if err != nil {
		t.Fatal(err)
	}

	if column.Equal(nil) {
		t.Fatal("compared non-nil column to nil and comparison was equal")
	}

	// check non-nil
	expected := &desc.Column{
		Name:     "name",
		DataType: desc.STRING,
	}
	if !column.Equal(expected) {
		t.Fatalf("expected %v and %v to be equal", column, expected)
	}
}

func TestNewSequenceColumn(t *testing.T) {
	column := desc.NewSequenceColumn("name", "seq")
	expected := &desc.Column{
		Name:     "name",
		DataType: desc.NUMBER,
		Sequence: "seq",
	}
	assert.Equal(t, expected, column)
}
