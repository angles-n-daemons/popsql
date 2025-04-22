package desc_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/db/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
	"github.com/angles-n-daemons/popsql/pkg/test/catalogT"
)

func TestNewSequence(t *testing.T) {
	name := "test_seq"
	seq := desc.NewSequence(name)

	assert.Equal(t, seq.SName, name)
	assert.Equal(t, seq.SID, uint64(0))
	assert.Equal(t, seq.V, uint64(0))
}

func TestSequenceEqual(t *testing.T) {
	// Test equal sequences
	seq1 := catalogT.SequenceWithName("seq1")
	seq2 := catalogT.CopySequence(seq1)
	assert.True(t, seq1.Equal(seq2))

	// Test different ID
	seq3 := catalogT.CopySequence(seq1)
	seq3.SID = seq1.SID + 1
	assert.False(t, seq1.Equal(seq3))

	// Test different Name
	seq4 := catalogT.CopySequence(seq1)
	seq4.SName = "seq4"
	assert.False(t, seq1.Equal(seq4))

	// Test different V
	seq5 := catalogT.CopySequence(seq1)
	seq5.V = seq1.V + 10
	assert.False(t, seq1.Equal(seq5))

	// Test nil comparison
	assert.False(t, seq1.Equal(nil))
}

func TestSequenceNext(t *testing.T) {
	seq := catalogT.Sequence()
	initialV := seq.V

	// First call should increment
	val := seq.Next()
	assert.Equal(t, val, initialV+1)
	assert.Equal(t, seq.V, initialV+1)

	// Second call should increment again
	val = seq.Next()
	assert.Equal(t, val, initialV+2)
	assert.Equal(t, seq.V, initialV+2)
}

func TestSequenceKey(t *testing.T) {
	seq := catalogT.Sequence()

	// Test the key is the string representation of the ID
	expected := strconv.FormatUint(seq.SID, 10)
	assert.Equal(t, seq.Key(), expected)

	// Test with specific ID
	seq.SID = 123
	assert.Equal(t, seq.Key(), "123")
}

func TestSequenceValue(t *testing.T) {
	seq := catalogT.Sequence()

	bytes, err := seq.Value()
	assert.NoError(t, err)

	var unmarshaled desc.Sequence
	err = json.Unmarshal(bytes, &unmarshaled)
	assert.NoError(t, err)

	assert.True(t, seq.Equal(&unmarshaled))
}
