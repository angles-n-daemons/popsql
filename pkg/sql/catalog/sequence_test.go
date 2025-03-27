package catalog_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/kv/memtable"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
	"github.com/angles-n-daemons/popsql/pkg/test/catalogT"
)

func TestCreateSequence(t *testing.T) {
	t.Run("basic use case", func(t *testing.T) {
		st := memtable.NewMemstore()
		m, err := catalog.NewManager(st)
		assert.NoError(t, err)

		// create the new sequence
		s := catalogT.Sequence()
		_, err = m.CreateSequence(s)
		assert.NoError(t, err)

		// verify the sequence is in the schema
		_, ok := m.GetSequence(s.Name)
		assert.True(t, ok)

		// verify that the sequence is in the store.
		// by loading a schema from it.
		sc, err := catalog.LoadSchema(st)
		assert.NoError(t, err)
		_, ok = sc.GetSequence(s.Name)
		assert.True(t, ok)
	})

	t.Run("sequence already exists", func(t *testing.T) {
		st := memtable.NewMemstore()
		m, err := catalog.NewManager(st)
		assert.NoError(t, err)
		s := catalogT.Sequence()

		// create the new sequence
		_, err = m.CreateSequence(s)
		assert.NoError(t, err)

		// try to create a sequence with the same name
		_, err = m.CreateSequence(s)
		assert.IsError(t, err, "sequence '%s' already exists", s.Name)
	})

}

func TestSequenceNext(t *testing.T) {
	t.Run("basic use case", func(t *testing.T) {
		st := memtable.NewMemstore()
		m, err := catalog.NewManager(st)
		assert.NoError(t, err)

		// create the new sequence
		s := catalogT.Sequence()
		_, err = m.CreateSequence(s)
		assert.NoError(t, err)
		assert.Equal(t, 0, s.V)

		// increment it a couple times
		val, err := m.SequenceNext(s)
		assert.NoError(t, err)
		assert.Equal(t, 1, val)
		val, err = m.SequenceNext(s)
		assert.NoError(t, err)
		assert.Equal(t, 2, val)
		assert.Equal(t, 2, s.V)

		// verify that the sequence value
		// by loading a schema with it.
		sc, err := catalog.LoadSchema(st)
		assert.NoError(t, err)
		s2, ok := sc.GetSequence(s.Name)
		assert.True(t, ok)
		assert.Equal(t, 2, s2.V)
	})

	t.Run("try to increment a sequence which doesn't exist", func(t *testing.T) {
		st := memtable.NewMemstore()
		m, err := catalog.NewManager(st)
		assert.NoError(t, err)

		// increment the sequence before it exists
		s := catalogT.Sequence()
		_, err = m.SequenceNext(s)
		assert.IsError(t, err, "sequence '%s' does not exist", s.Name)
	})

}
