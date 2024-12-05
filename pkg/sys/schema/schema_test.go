package schema_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
	"github.com/angles-n-daemons/popsql/pkg/testing/assert"
)

func SchemaFromBytes(t *testing.T) {

}

func TestAddTable(t *testing.T) {
	sc := schema.NewSchema()
	expected := testTableFromArgs("tt", nil, nil)
	err := sc.AddTable(expected)
	assert.NoError(t, err)
	actual, err := sc.GetTable(expected.Key())

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestAddExistingTable(t *testing.T) {
	sc := schema.NewSchema()
	table1 := testTableFromArgs("tt", nil, nil)
	err := sc.AddTable(table1)
	assert.NoError(t, err)

	table2 := testTableFromArgs("tt", nil, nil)
	err = sc.AddTable(table2)
	assert.IsError(t, err, "table 'tt' already exists")

}

func TestGetTable(t *testing.T) {
	sc := schema.NewSchema()
	expected := testTableFromArgs("tt", nil, nil)
	err := sc.AddTable(expected)
	assert.NoError(t, err)
	actual, err := sc.GetTable(expected.Key())

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetMissingTable(t *testing.T) {
	sc := schema.NewSchema()
	table, err := sc.GetTable("doesntexist")
	assert.Nil(t, table)
	assert.IsError(t, err, "could not find table 'doesntexist'")
}

func TestDropTable(t *testing.T) {
	sc := schema.NewSchema()
	table := testTableFromArgs("tt", nil, nil)

	err := sc.AddTable(table)
	assert.NoError(t, err)

	err = sc.DropTable(table.Key())
	assert.NoError(t, err)

	retrieved, err := sc.GetTable(table.Key())
	assert.Nil(t, retrieved)
	assert.IsError(t, err, "could not find table 'tt'")
}

func TestDropMissingTable(t *testing.T) {
	sc := schema.NewSchema()
	err := sc.DropTable("doesntexist")
	assert.IsError(t, err, "could not delete table 'doesntexist'")
}
