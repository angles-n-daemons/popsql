package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema/desc"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
	"github.com/angles-n-daemons/popsql/pkg/test/schematest"
)

func TestSchemaFromBytes(t *testing.T) {
	tt := schematest.TestTable()
	tf := schematest.CopyTable(tt)
	tf.Name = "tf"
	tablesBytes := [][]byte{}
	sc := schema.NewSchema()

	for _, table := range []*desc.Table{tt, tf} {
		sc.AddTable(table)
		b, err := json.Marshal(table)
		assert.NoError(t, err)
		tablesBytes = append(tablesBytes, b)
	}

	sc2 := schema.NewSchema()
	err := sc2.LoadTables(tablesBytes)
	assert.NoError(t, err)
	assert.Equal(t, sc, sc2)
}

func TestSchemaAddTable(t *testing.T) {
	sc := schema.NewSchema()
	expected := schematest.TestTable()
	err := sc.AddTable(expected)
	assert.NoError(t, err)
	actual, ok := sc.GetTable(expected.Name)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestSchemaAddExistingTable(t *testing.T) {
	sc := schema.NewSchema()
	table1 := schematest.TableWithID(1)
	err := sc.AddTable(table1)
	assert.NoError(t, err)

	table2 := schematest.CopyTable(table1)
	err = sc.AddTable(table2)
	assert.IsError(t, err, "table 'table_1' already exists")

}

func TestSchemaGetTable(t *testing.T) {
	sc := schema.NewSchema()
	expected := schematest.TestTable()
	err := sc.AddTable(expected)
	assert.NoError(t, err)
	actual, ok := sc.GetTable(expected.Name)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestSchemaGetMissingTable(t *testing.T) {
	sc := schema.NewSchema()
	table, ok := sc.GetTable("doesntexist")
	assert.Nil(t, table)
	assert.False(t, ok)
}

func TestSchemaDropTable(t *testing.T) {
	sc := schema.NewSchema()
	table := schematest.TestTable()

	err := sc.AddTable(table)
	assert.NoError(t, err)

	err = sc.DropTable(table.Name)
	assert.NoError(t, err)

	retrieved, ok := sc.GetTable(table.Name)
	assert.Nil(t, retrieved)
	assert.False(t, ok)
}

func TestSchemaDropMissingTable(t *testing.T) {
	sc := schema.NewSchema()
	err := sc.DropTable("doesntexist")
	assert.IsError(t, err, "could not delete table 'doesntexist'")
}
