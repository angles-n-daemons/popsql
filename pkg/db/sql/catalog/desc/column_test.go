package desc_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/db/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
)

func TestNewColumn(t *testing.T) {
	column := desc.NewColumn("name", desc.STRING)
	expected := &desc.Column{
		Name:     "name",
		DataType: desc.STRING,
	}
	assert.Equal(t, expected, column)
}

func TestColumnEqual(t *testing.T) {
	// check nil condition
	column := desc.NewColumn("name", desc.STRING)

	if column.Equal(nil) {
		t.Fatal("compared non-nil column to nil and comparison was equal")
	}

	// check non-nil
	expected := &desc.Column{
		Name:     "name",
		DataType: desc.STRING,
	}
	assert.Equal(t, expected, column)
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
