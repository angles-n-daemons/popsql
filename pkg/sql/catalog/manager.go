package catalog

import (
	"encoding/json"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/meta"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
)

type Manager struct {
	Schema *schema.Schema
	Store  kv.Store
	Meta   *meta.Meta
}

// Create is a manager function for adding new system objects to the catalog.
func Create[V schema.Collectible[V]](m *Manager, v V) (uint64, error) {
	// Get the id for the new object.
	sysSeq := GetSystemSequence[V](m)
	id, err := SequenceNext(m, sysSeq)
	if err != nil {
		return 0, err
	}

	return createWithID(m, v, id)
}

func createWithID[V schema.Collectible[V]](m *Manager, v V, id uint64) (uint64, error) {
	// Set the ID of the new object.
	v.WithID(id)

	// add it to the underlying schema.
	err := schema.Add(m.Schema, v)
	if err != nil {
		return 0, err
	}
	// save it to the storage engine.
	err = Save(m, v)
	if err != nil {
		return 0, nil
	}

	return id, nil
}

// Save exists to store a collectible in the underlying store.
// It's used both by Add for new objects, and on its own to save changes to
// existing objects.
func Save[V schema.Collectible[V]](m *Manager, v V) error {
	// Get the system table so that we can save the object.
	sysTable := GetSystemTable[V](m)
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return m.Store.Put(sysTable.Prefix().WithID(v.Key()).Encode(), b)
}

func SequenceNext(m *Manager, s *desc.Sequence) (uint64, error) {
	// Get the next value in the sequence.
	next := s.Next()

	// Update the sequence in the store.
	err := Save(m, s)
	if err != nil {
		return 0, err
	}
	return next, nil
}

func GetSystemSequence[V schema.Collectible[V]](m *Manager) *desc.Sequence {
	var zero V
	switch any(zero).(type) {
	case *desc.Table:
		return m.Meta.Tables.Sequence
	case *desc.Sequence:
		return m.Meta.Sequences.Sequence
	}
	return nil
}

func GetSystemTable[V schema.Collectible[V]](m *Manager) *desc.Table {
	var zero V
	switch any(zero).(type) {
	case *desc.Table:
		return m.Meta.Tables.Table
	case *desc.Sequence:
		return m.Meta.Sequences.Table
	}
	return nil
}
