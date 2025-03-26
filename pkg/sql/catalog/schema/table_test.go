package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
	"github.com/angles-n-daemons/popsql/pkg/test/catalogT"
)

func TestLoadTables(t *testing.T) {
	tt := catalogT.Table()
	tf := catalogT.CopyTable(tt)
	tf.Name = "tf"
	tablesBytes := [][]byte{}
	sc := schema.New()

	for _, table := range []*desc.Table{tt, tf} {
		sc.AddTable(table)
		b, err := json.Marshal(table)
		assert.NoError(t, err)
		tablesBytes = append(tablesBytes, b)
	}

	sc2 := schema.New()
	err := sc2.LoadTables(tablesBytes)
	assert.NoError(t, err)
	assert.Equal(t, sc, sc2)
}

func TestAddTable(t *testing.T) {
	sc := schema.New()
	expected := catalogT.Table()
	err := sc.AddTable(expected)
	assert.NoError(t, err)
	actual, ok := sc.GetTable(expected.Name)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestAddExistingTable(t *testing.T) {
	sc := schema.New()
	table1 := catalogT.TableWithID(1)
	err := sc.AddTable(table1)
	assert.NoError(t, err)

	table2 := catalogT.CopyTable(table1)
	err = sc.AddTable(table2)
	assert.IsError(t, err, "table 'table_1' already exists")

}

func TestGetTable(t *testing.T) {
	sc := schema.New()
	expected := catalogT.Table()
	err := sc.AddTable(expected)
	assert.NoError(t, err)
	actual, ok := sc.GetTable(expected.Name)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestGetMissingTable(t *testing.T) {
	sc := schema.New()
	table, ok := sc.GetTable("doesntexist")
	assert.Nil(t, table)
	assert.False(t, ok)
}

func TestDropTable(t *testing.T) {
	sc := schema.New()
	table := catalogT.Table()

	err := sc.AddTable(table)
	assert.NoError(t, err)

	err = sc.RemoveTable(table.Name)
	assert.NoError(t, err)

	retrieved, ok := sc.GetTable(table.Name)
	assert.Nil(t, retrieved)
	assert.False(t, ok)
}

func TestDropMissingTable(t *testing.T) {
	sc := schema.New()
	err := sc.RemoveTable("doesntexist")
	assert.IsError(t, err, "could not delete table 'doesntexist'")
}
