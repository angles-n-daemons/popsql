package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
	"github.com/angles-n-daemons/popsql/pkg/testing/assert"
)

func TestSchemaFromBytes(t *testing.T) {
	tt := testTableFromArgs("tt", nil, nil)
	tf := testTableFromArgs("tf", nil, nil)
	sc := schema.NewSchema()
	tablesBytes := [][]byte{}

	for _, table := range []*schema.Table{tt, tf} {
		b, err := json.Marshal(table)
		assert.NoError(t, err)
		tablesBytes = append(tablesBytes, b)
		err = sc.AddTable(tt)
		assert.NoError(t, err)
	}

	generated, err := schema.SchemaFromBytes(tablesBytes)
	assert.NoError(t, err)
	assert.Equal(t, sc, generated)
}

func TestSchemaAddTable(t *testing.T) {
	sc := schema.NewSchema()
	expected := testTableFromArgs("tt", nil, nil)
	err := sc.AddTable(expected)
	assert.NoError(t, err)
	actual, err := sc.GetTable(expected.Key())

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestSchemaAddExistingTable(t *testing.T) {
	sc := schema.NewSchema()
	table1 := testTableFromArgs("tt", nil, nil)
	err := sc.AddTable(table1)
	assert.NoError(t, err)

	table2 := testTableFromArgs("tt", nil, nil)
	err = sc.AddTable(table2)
	assert.IsError(t, err, "table 'tt' already exists")

}

func TestSchemaGetTable(t *testing.T) {
	sc := schema.NewSchema()
	expected := testTableFromArgs("tt", nil, nil)
	err := sc.AddTable(expected)
	assert.NoError(t, err)
	actual, err := sc.GetTable(expected.Key())

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestSchemaGetMissingTable(t *testing.T) {
	sc := schema.NewSchema()
	table, err := sc.GetTable("doesntexist")
	assert.Nil(t, table)
	assert.IsError(t, err, "could not find table 'doesntexist'")
}

func TestSchemaDropTable(t *testing.T) {
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

func TestSchemaDropMissingTable(t *testing.T) {
	sc := schema.NewSchema()
	err := sc.DropTable("doesntexist")
	assert.IsError(t, err, "could not delete table 'doesntexist'")
}

func TestSchemaEqual(t *testing.T) {
	assert.Nil(t, 5)
}
