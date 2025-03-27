package catalog_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
	"github.com/angles-n-daemons/popsql/pkg/test/catalogT"
)

func TestCreateTable(t *testing.T) {
	t.Run("basic use case", func(t *testing.T) {
		m := catalogT.Manager(t)
		tb := catalogT.Table()
		assert.NoError(t, m.CreateTable(tb))

		// verify table is in the schema
		scT, ok := m.Schema.GetTable(tb.Name)
		assert.True(t, ok)

		// verify id has been changed to whatever the manager assigned
		assert.Equal(t, scT.ID, m.Sys.MetaTableSequence.V)

		// verify table loads on a fresh schema load
		sc, err := catalog.LoadSchema(m.Store)
		assert.NoError(t, err)
		scT, ok = sc.GetTable(tb.Name)

		assert.Equal(t, scT.ID, m.Sys.MetaTableSequence.V)
		assert.Equal(t, tb, scT)
	})

	t.Run("table already exists", func(t *testing.T) {
		m := catalogT.Manager(t)
		tb := catalogT.Table()
		assert.NoError(t, m.CreateTable(tb))
		startingID := tb.ID
		assert.IsError(t, m.CreateTable(tb), "table '%s' already exists", tb.Name)

		// verify that the sequence is still permanently incremented
		tb2 := catalogT.Table()
		assert.NoError(t, m.CreateTable(tb2))
		assert.Equal(t, tb2.ID, startingID+2)
	})

	t.Run("with no primary key", func(t *testing.T) {
		m := catalogT.Manager(t)
		tb, err := desc.NewTable("hi", nil, nil)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(tb.PrimaryKey))
		assert.NoError(t, m.CreateTable(tb))
		assert.Equal(t, 1, len(tb.PrimaryKey))
	})
}
