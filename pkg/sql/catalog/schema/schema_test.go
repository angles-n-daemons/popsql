package schema_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
	"github.com/angles-n-daemons/popsql/pkg/test/catalogT"
)

func TestEmpty(t *testing.T) {
	t.Run("actually empty", func(t *testing.T) {
		sc := schema.New()
		assert.True(t, sc.Empty())
	})
	t.Run("has sequence, no tables", func(t *testing.T) {
		sc := schema.New()
		err := sc.AddSequence(catalogT.Sequence())
		assert.NoError(t, err)
		assert.False(t, sc.Empty())
	})

	t.Run("has table, no sequence", func(t *testing.T) {
		sc := schema.New()
		err := sc.AddTable(catalogT.Table())
		assert.NoError(t, err)
		assert.False(t, sc.Empty())
	})

	t.Run("has sequence and table", func(t *testing.T) {
		sc := schema.New()

		err := sc.AddTable(catalogT.Table())
		assert.NoError(t, err)
		err = sc.AddSequence(catalogT.Sequence())
		assert.NoError(t, err)

		assert.False(t, sc.Empty())
	})
	t.Run("has tables, everything removed, empty is true", func(t *testing.T) {
		sc := schema.New()

		// add stuff
		tb := catalogT.Table()
		err := sc.AddTable(tb)
		assert.NoError(t, err)
		s := catalogT.Sequence()
		err = sc.AddSequence(s)
		assert.NoError(t, err)

		// remove stuff
		err = sc.RemoveSequence(s.Name)
		assert.NoError(t, err)
		err = sc.RemoveTable(tb.Name)
		assert.NoError(t, err)

		assert.True(t, sc.Empty())
	})
}

func TestEqual(t *testing.T) {
	t.Run("nil other", func(t *testing.T) {
		sc := schema.New()
		assert.NotEqual(t, sc, nil)
	})

	// This tests demonstrates how a schema with multiple tables and
	// sequences can be equal if the internal types are equivalent.
	t.Run("same tables and sequences", func(t *testing.T) {
		sc := schema.New()
		sc2 := schema.New()

		tb1 := catalogT.Table()
		tb2 := catalogT.Table()

		s1 := catalogT.Sequence()
		s2 := catalogT.Sequence()

		assert.NoError(t, sc.AddTable(tb1))
		assert.NoError(t, sc.AddTable(tb2))
		assert.NoError(t, sc.AddSequence(s1))
		assert.NoError(t, sc.AddSequence(s2))

		assert.NoError(t, sc2.AddTable(tb1))
		assert.NoError(t, sc2.AddTable(tb2))
		assert.NoError(t, sc2.AddSequence(s1))
		assert.NoError(t, sc2.AddSequence(s2))

		assert.Equal(t, sc, sc2)
	})

	t.Run("empty schemas equal", func(t *testing.T) {
		sc := schema.New()
		sc2 := schema.New()
		assert.Equal(t, sc, sc2)
	})

	t.Run("different number of tables", func(t *testing.T) {
		sc := schema.New()
		sc2 := schema.New()

		tb1 := catalogT.Table()
		tb2 := catalogT.Table()

		// add one table to sc
		assert.NoError(t, sc.AddTable(tb1))
		// add both tables to sc2
		assert.NoError(t, sc2.AddTable(tb1))
		assert.NoError(t, sc2.AddTable(tb2))

		assert.NotEqual(t, sc, sc2)
	})

	t.Run("different number of sequences", func(t *testing.T) {
		sc := schema.New()
		sc2 := schema.New()

		s1 := catalogT.Sequence()
		s2 := catalogT.Sequence()

		// add both sequences to sc
		assert.NoError(t, sc.AddSequence(s1))
		assert.NoError(t, sc.AddSequence(s2))

		// only add one sequence to sc2
		assert.NoError(t, sc2.AddSequence(s1))

		assert.NotEqual(t, sc, sc2)
	})

	t.Run("different tables not equal", func(t *testing.T) {
		sc := schema.New()
		sc2 := schema.New()

		tb1 := catalogT.Table()
		tb2 := catalogT.Table()

		// add tb1 to sc
		assert.NoError(t, sc.AddTable(tb1))
		// add tb2 to sc2
		assert.NoError(t, sc2.AddTable(tb2))

		assert.NotEqual(t, sc, sc2)
	})

	t.Run("different sequences not equal", func(t *testing.T) {
		sc := schema.New()
		sc2 := schema.New()

		s1 := catalogT.Sequence()
		s2 := catalogT.Sequence()

		// add s1 to sc
		assert.NoError(t, sc.AddSequence(s1))

		// add s2 to sc
		assert.NoError(t, sc2.AddSequence(s2))

		assert.NotEqual(t, sc, sc2)
	})

	t.Run("tables with same key are different", func(t *testing.T) {
		sc := schema.New()
		sc2 := schema.New()

		tom := catalogT.TableWithName("tom")
		tom2 := catalogT.TableWithName("tom")

		// add tb1 to sc
		assert.NoError(t, sc.AddTable(tom))
		// add tb2 to sc2
		assert.NoError(t, sc2.AddTable(tom2))

		assert.NotEqual(t, sc, sc2)
	})

	t.Run("sequences with same key are different", func(t *testing.T) {
		sc := schema.New()
		sc2 := schema.New()

		gerry := catalogT.SequenceWithName("gerry")
		gerry2 := catalogT.SequenceWithName("gerry")

		// add tb1 to sc
		assert.NoError(t, sc.AddSequence(gerry))
		// add tb2 to sc2
		assert.NoError(t, sc2.AddSequence(gerry2))

		assert.NotEqual(t, sc, sc2)
	})
}
