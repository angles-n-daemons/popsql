package schema

import (
	"encoding/json"
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/db/sql/catalog/desc"
)

// Collection is a generic collection utility, made for managing
// system objects. It is parameterized by a desc.Object type,
// which represents any object that would need to be cached
// in the schema. In practice, this consists of descriptors
// for things like tables and sequences.
type Collection[V desc.Any[V]] struct {
	byID   map[uint64]V
	byName map[string]V
}

func NewCollection[V desc.Any[V]]() *Collection[V] {
	return &Collection[V]{
		byID:   make(map[uint64]V),
		byName: make(map[string]V),
	}
}

func CollectionFromBytes[V desc.Any[V]](bytesArr [][]byte) (*Collection[V], error) {
	c := NewCollection[V]()
	for _, b := range bytesArr {
		var v V
		err := json.Unmarshal(b, &v)
		if err != nil {
			return nil, err
		}
		c.Add(v)
	}
	return c, nil
}

func (c *Collection[V]) Add(v V) error {
	id := v.ID()
	name := v.Name()
	if _, ok := c.byID[id]; ok {
		return fmt.Errorf("%T with id '%d' already exists", *new(V), id)
	}
	if _, ok := c.byName[name]; ok {
		return fmt.Errorf("%T with name '%s' already exists", *new(V), name)
	}
	c.byName[name] = v
	c.byID[id] = v
	return nil
}

func (c *Collection[V]) All() []V {
	results := []V{}
	for _, v := range c.byID {
		results = append(results, v)
	}
	return results
}

func (c *Collection[V]) Get(id uint64) V {
	return c.byID[id]
}

func (c *Collection[V]) GetByName(name string) V {
	return c.byName[name]
}

func (c *Collection[V]) Remove(id uint64) error {
	v, ok := c.byID[id]
	if !ok {
		return fmt.Errorf("could not delete %T with id '%d'", *new(V), id)
	}
	delete(c.byName, v.Name())
	delete(c.byID, id)
	return nil
}

func (c *Collection[V]) Size() int {
	return len(c.byName)
}

func (c *Collection[V]) Empty() bool {
	return len(c.byName) == 0
}

func (c *Collection[V]) Equal(o *Collection[V]) bool {
	if o == nil {
		return false
	}
	if c.Size() != o.Size() {
		return false
	}
	for _, v := range c.byID {
		ov := o.Get(v.ID())
		if !v.Equal(ov) {
			return false
		}
	}
	return true
}
