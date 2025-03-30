package schema_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
)

type MockCollectible struct {
	id   uint64
	name string
}

func (m *MockCollectible) WithID(id uint64) {
	m.id = id
}

func (m *MockCollectible) ID() uint64 {
	return m.id
}

func (m *MockCollectible) Key() string {
	return m.name
}

func (m *MockCollectible) Name() string {
	return m.name
}

func (c *MockCollectible) Equal(o *MockCollectible) bool {
	if o == nil {
		return false
	}
	return c.ID() == o.ID() && c.Name() == o.Name()
}

func NewMockCollectible(id uint64, name string) *MockCollectible {
	return &MockCollectible{id: id, name: name}
}

func TestNewCollection(t *testing.T) {
	c := schema.NewCollection[*MockCollectible]()
	assert.NotNil(t, c)
	assert.True(t, c.Empty())
	// verify that the id maps are initialized by checking no panic.
	c.Get(0)
	c.GetByName("")
}

func TestCollectionAdd(t *testing.T) {
	// happy path
	c := schema.NewCollection[*MockCollectible]()
	mc := NewMockCollectible(1, "one")
	err := c.Add(mc)
	assert.NoError(t, err)

	// collision on id
	mc2 := NewMockCollectible(1, "two")
	err = c.Add(mc2)
	assert.IsError(t, err, "%T with id '%d' already exists", mc, mc.ID())

	// collision on name
	mc3 := NewMockCollectible(2, "one")
	err = c.Add(mc3)
	assert.IsError(t, err, "%T with name '%s' already exists", mc, mc.Name())

	// add works after removal
	err = c.Remove(mc.ID())
	assert.NoError(t, err)
	err = c.Add(mc2)
	assert.NoError(t, err)
	err = c.Add(mc3)
	assert.NoError(t, err)
}

func TestCollectionGet(t *testing.T) {
	c := schema.NewCollection[*MockCollectible]()
	mc := NewMockCollectible(1, "one")
	err := c.Add(mc)
	assert.NoError(t, err)

	// happy path
	mc2 := c.Get(mc.ID())
	assert.Equal(t, mc, mc2)

	// missing id
	mc3 := c.Get(0)
	assert.Nil(t, mc3)

	// remove and check
	err = c.Remove(mc.ID())
	assert.NoError(t, err)
	mc4 := c.Get(mc.ID())
	assert.Nil(t, mc4)
}

func TestCollectionGetByName(t *testing.T) {
	c := schema.NewCollection[*MockCollectible]()
	mc := NewMockCollectible(1, "one")
	err := c.Add(mc)
	assert.NoError(t, err)

	// happy path
	mc2 := c.GetByName(mc.Name())
	assert.Equal(t, mc, mc2)

	// missing name
	mc3 := c.GetByName("not a name")
	assert.Nil(t, mc3)

	// remove and check
	err = c.Remove(mc.ID())
	assert.NoError(t, err)
	mc4 := c.GetByName(mc.Name())
	assert.Nil(t, mc4)
}

func TestCollectionRemove(t *testing.T) {
	c := schema.NewCollection[*MockCollectible]()
	mc := NewMockCollectible(1, "one")
	err := c.Add(mc)
	assert.NoError(t, err)

	// happy path
	err = c.Remove(mc.ID())
	assert.NoError(t, err)

	// missing from id map
	mc2 := c.Get(mc.ID())
	assert.Nil(t, mc2)

	// missing from name map
	mc3 := c.Get(mc.ID())
	assert.Nil(t, mc3)

	// missing id
	err = c.Remove(0)
	assert.IsError(t, err, "could not delete %T with id '%d'", mc, 0)
}

func TestCollectionEmpty(t *testing.T) {
	// start empty
	c := schema.NewCollection[*MockCollectible]()
	assert.True(t, c.Empty())

	// add some stuff
	mc := NewMockCollectible(1, "one")
	err := c.Add(mc)
	assert.NoError(t, err)
	mc2 := NewMockCollectible(2, "two")
	err = c.Add(mc2)
	assert.NoError(t, err)

	// verify not empty
	assert.False(t, c.Empty())

	// remove one thing and check still not empty
	err = c.Remove(mc.ID())
	assert.NoError(t, err)
	assert.False(t, c.Empty())

	// remove the final thing and check empty
	err = c.Remove(mc2.ID())
	assert.NoError(t, err)
	assert.True(t, c.Empty())
}

func TestCollectionSize(t *testing.T) {
	// start empty
	c := schema.NewCollection[*MockCollectible]()
	assert.Equal(t, 0, c.Size())

	// add some stuff
	mc := NewMockCollectible(1, "one")
	err := c.Add(mc)
	assert.NoError(t, err)
	mc2 := NewMockCollectible(2, "two")
	err = c.Add(mc2)
	assert.NoError(t, err)

	// verify not empty
	assert.Equal(t, 2, c.Size())

	// remove one thing and check still not empty
	err = c.Remove(mc.ID())
	assert.NoError(t, err)
	assert.Equal(t, 1, c.Size())

	// remove the final thing and check empty
	err = c.Remove(mc2.ID())
	assert.NoError(t, err)
	assert.Equal(t, 0, c.Size())

	// another remove doesn't do anything to the size
	err = c.Remove(0)
	assert.IsError(t, err, "could not delete %T with id '%d'", mc, 0)
	assert.Equal(t, 0, c.Size())

}

func TestCollectionEqual(t *testing.T) {
	// nil collection check.
	c := schema.NewCollection[*MockCollectible]()
	assert.NotEqual(t, c, nil)

	// empty collection check
	c2 := schema.NewCollection[*MockCollectible]()
	assert.Equal(t, c, c2)

	// when c1 has more elements than c2.
	mc := NewMockCollectible(1, "one")
	err := c.Add(mc)
	assert.NoError(t, err)
	assert.NotEqual(t, c, c2)

	// when c2 gets the same elements as c1.
	err = c2.Add(mc)
	assert.NoError(t, err)
	assert.Equal(t, c, c2)

	// when c1 and c2 have the same size but different elements.
	mc2 := NewMockCollectible(2, "two")
	mc3 := NewMockCollectible(3, "three")
	err = c.Add(mc2)
	assert.NoError(t, err)
	err = c2.Add(mc3)
	assert.NoError(t, err)
	assert.NotEqual(t, c, c2)

	// when the differing elements are removed.
	err = c.Remove(mc2.ID())
	assert.NoError(t, err)
	err = c2.Remove(mc3.ID())
	assert.NoError(t, err)
	assert.Equal(t, c, c2)
}
