package catalog_test

import (
	"encoding/json"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/kv/memtable"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
	"github.com/angles-n-daemons/popsql/pkg/test/catalogT"
)

func checkManagerTable(t *testing.T, m *catalog.Manager, tb *desc.Table) {
	scT, ok := m.Schema.GetTable(tb.Name)
	assert.True(t, ok)
	assert.Equal(t, tb, scT)

	scT = catalogT.ReadTable(t, m.Store, m.TableKey(tb).Encode())
	assert.Equal(t, tb, scT)
}

func checkManagerSequence(t *testing.T, m *catalog.Manager, s *desc.Sequence) {
	scS, ok := m.Schema.GetSequence(s.Name)

	assert.True(t, ok)
	assert.Equal(t, s, scS)

	scS = catalogT.ReadSequence(t, m.Store, m.SequenceKey(s).Encode())
	assert.Equal(t, s, scS)
}

func TestNewManager(t *testing.T) {
	t.Run("test empty use case", func(t *testing.T) {
		st := memtable.NewMemstore()
		m, err := catalog.NewManager(st)
		assert.NoError(t, err)

		checkManagerTable(t, m, catalog.InitMetaTable())
		checkManagerSequence(t, m, catalog.InitMetaTableSequence())
		checkManagerTable(t, m, catalog.InitSequencesTable())
		checkManagerSequence(t, m, catalog.InitSequencesTableSequence())
	})

	t.Run("test partially populated store", func(t *testing.T) {
		// mt := catalogT.CopyTable(catalog.InitMetaTable)
		// mts := catalogT.CopySequence(catalog.InitMetaTableSequence)
	})

	t.Run("test already populated with system objects", func(t *testing.T) {
		st := memtable.NewMemstore()
		m, err := catalog.NewManager(st)
		assert.NoError(t, err)

		assert.Equal(t, m.Sys.MetaTableSequence.V, 2)
		assert.Equal(t, m.Sys.MetaTableSequence.V, 2)

		// fiddle around with the sequence values.
		m.SequenceNext(m.Sys.MetaTableSequence)
		m.SequenceNext(m.Sys.MetaTableSequence)
		m.SequenceNext(m.Sys.SequencesTableSequence)

		// verify reiniting a manager from the above store has the updated values.
		assert.Equal(t, 4, m.Sys.MetaTableSequence.V)
		assert.Equal(t, 3, m.Sys.SequencesTableSequence.V)

		// verify an empty store starts with the values back at 2
		st2 := memtable.NewMemstore()
		m2, err := catalog.NewManager(st2)
		assert.NoError(t, err)

		assert.Equal(t, 2, m2.Sys.MetaTableSequence.V)
		assert.Equal(t, 2, m2.Sys.SequencesTableSequence.V)
	})
}

func TestLoadSchema(t *testing.T) {
	t.Run("empty store", func(t *testing.T) {
		st := memtable.NewMemstore()
		sc, err := catalog.LoadSchema(st)
		assert.NoError(t, err)
		assert.True(t, sc.Empty())
	})

	t.Run("only system data", func(t *testing.T) {
		st := memtable.NewMemstore()
		m, err := catalog.NewManager(st)
		assert.NoError(t, err)

		// load a new schema from the store.
		sc, err := catalog.LoadSchema(st)
		assert.NoError(t, err)

		_, ok := sc.GetTable(m.Sys.MetaTable.Name)
		assert.True(t, ok)
	})

	t.Run("system and user data", func(t *testing.T) {
		st := memtable.NewMemstore()
		m, err := catalog.NewManager(st)
		assert.NoError(t, err)

		// add a user table
		tb := catalogT.Table()
		assert.NoError(t, m.CreateTable(tb))

		// load a new schema from the store.
		sc, err := catalog.LoadSchema(st)
		assert.NoError(t, err)

		_, ok := sc.GetTable(m.Sys.MetaTable.Name)
		assert.True(t, ok)
		_, ok = sc.GetTable(tb.Name)
		assert.True(t, ok)
	})

	t.Run("only user data", func(t *testing.T) {
		st := memtable.NewMemstore()
		tb := catalogT.Table()
		key := catalog.InitMetaTable().Prefix().WithID(tb.Key())
		b, err := json.Marshal(tb)
		assert.NoError(t, err)
		st.Put(key.Encode(), b)

		// load a new schema from the store.
		sc, err := catalog.LoadSchema(st)

		// system table not there, user table should be there.
		_, ok := sc.GetTable(catalog.InitMetaTable().Name)
		assert.False(t, ok)
		_, ok = sc.GetTable(tb.Name)
		assert.True(t, ok)
	})

}
