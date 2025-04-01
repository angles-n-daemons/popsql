package catalog

import (
	"encoding/json"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/sys"
)

type Manager struct {
	Schema *schema.Schema
	Store  kv.Store
}

var (
	tTable    = &desc.Table{}
	tSequence = &desc.Sequence{}
)

// Create is a manager function for adding new system objects to
// the catalog.
func Create[V desc.Any[V]](m *Manager, v V) error {
	// add it to the underlying schema.
	err := schema.Add(m.Schema, v)
	if err != nil {
		return err
	}
	// save it to the storage engine.
	err = save(m, v)
	if err != nil {
		return nil
	}

	return nil
}

// NextDescriptorID is a utility function for getting the next
// available id for a type of descriptor in the system.
func NextDescriptorID[V desc.Any[V]](m *Manager, v V) (uint64, error) {
	// get the sequence fot this type.
	s := getSystemSequence[V](m.Schema)

	return SequenceNext(m, s)
}

// SequenceNext is used to get the next value in a sequence. It
// also writes the value of the sequence to disk, ignoring any
// transaction semantics.
func SequenceNext(m *Manager, s *desc.Sequence) (uint64, error) {
	// Get the next value in the sequence.
	next := s.Next()

	// Update the sequence in the store.
	err := save(m, s)
	if err != nil {
		return 0, err
	}
	return next, nil
}

// save exists to store a collectible in the underlying store.
// It's used both by Add for new objects, and on its own to save
// changes to existing objects.
func save[V desc.Any[V]](m *Manager, v V) error {
	// Get the system table so that we can save the object.
	sysTable := getSystemTable[V](m.Schema)
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return m.Store.Put(sysTable.Prefix().WithID(v.Key()).Encode(), b)
}

func getSystemTable[V desc.Any[V]](sc *schema.Schema) *desc.Table {
	return schema.Get[*desc.Table](sc, getSystemID[V]())
}

func getSystemSequence[V desc.Any[V]](sc *schema.Schema) *desc.Sequence {
	return schema.Get[*desc.Sequence](sc, getSystemID[V]())
}

func getSystemID[V desc.Any[V]]() uint64 {
	switch any(*new(V)).(type) {
	case *desc.Table:
		return sys.TablesID
	case *desc.Sequence:
		return sys.SequencesID
	}
	return 0
}
