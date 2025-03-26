package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
	"github.com/angles-n-daemons/popsql/pkg/test/catalogT"
)

func TestLoadSequences(t *testing.T) {
	s := catalogT.Sequence()
	sf := catalogT.CopySequence(s)
	sf.Name = "sf"
	sequencesBytes := [][]byte{}
	sc := schema.New()

	for _, sequence := range []*desc.Sequence{s, sf} {
		sc.AddSequence(sequence)
		b, err := json.Marshal(sequence)
		assert.NoError(t, err)
		sequencesBytes = append(sequencesBytes, b)
	}

	sc2 := schema.New()
	err := sc2.LoadSequences(sequencesBytes)
	assert.NoError(t, err)
	assert.Equal(t, sc, sc2)
}

func TestAddSequence(t *testing.T) {
	sc := schema.New()
	expected := catalogT.Sequence()
	err := sc.AddSequence(expected)
	assert.NoError(t, err)
	actual, ok := sc.GetSequence(expected.Name)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestGetSequence(t *testing.T) {
	sc := schema.New()
	expected := catalogT.Sequence()
	err := sc.AddSequence(expected)
	assert.NoError(t, err)
	actual, ok := sc.GetSequence(expected.Name)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestGetMissingSequence(t *testing.T) {
	sc := schema.New()
	table, ok := sc.GetSequence("doesntexist")
	assert.Nil(t, table)
	assert.False(t, ok)
}

func TestRemoveSequence(t *testing.T) {
	sc := schema.New()
	table := catalogT.Sequence()

	err := sc.AddSequence(table)
	assert.NoError(t, err)

	err = sc.RemoveSequence(table.Name)
	assert.NoError(t, err)

	retrieved, ok := sc.GetSequence(table.Name)
	assert.Nil(t, retrieved)
	assert.False(t, ok)
}

func TestRemoveMissingSequence(t *testing.T) {
	sc := schema.New()
	err := sc.RemoveSequence("doesntexist")
	assert.IsError(t, err, "could not delete table 'doesntexist'")
}
