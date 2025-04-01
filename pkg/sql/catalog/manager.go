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
	Meta   *sys.SystemSchema
}

// Create is a manager function for adding new system objects to the catalog.
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

// NextDescriptorID is a utility function for getting the next available
// id for a type of descriptor in the system.
func NextDescriptorID[V desc.Any[V]](m *Manager, v V) (uint64, error) {
	// get the sequence fot this type.
	s := getSystemSequence[V](m.Meta)

	// Get the next value in the sequence.
	next := s.Next()

	// Update the sequence in the store.
	err := save(m, s)
	if err != nil {
		return 0, err
	}
	return next, nil
}

// SequenceNext is used to get the next value in a sequence.
// It also writes the value of the sequence to disk, ignoring any transaction
// semantics.
func SequenceNext[V desc.Any[V]](m *Manager, v V) (uint64, error) {
	// get the sequence fot this type.
	s := getSystemSequence[V](m.Meta)

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
// It's used both by Add for new objects, and on its own to save changes to
// existing objects.
func save[V desc.Any[V]](m *Manager, v V) error {
	// Get the system table so that we can save the object.
	sysTable := getSystemTable[V](m.Meta)
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return m.Store.Put(sysTable.Prefix().WithID(v.Key()).Encode(), b)
}

func getSystemTable[V desc.Any[V]](mt *sys.SystemSchema) *desc.Table {
	var zero V
	switch any(zero).(type) {
	case *desc.Table:
		return mt.Tables.Table
	case *desc.Sequence:
		return mt.Sequences.Table
	}
	return nil
}

func getSystemSequence[V desc.Any[V]](mt *sys.SystemSchema) *desc.Sequence {
	var zero V
	switch any(zero).(type) {
	case *desc.Table:
		return mt.Tables.Sequence
	case *desc.Sequence:
		return mt.Sequences.Sequence
	}
	return nil
}
