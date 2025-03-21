package catalog

import (
	"errors"
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
)

// Custom error for dropping the Meta table.
var ErrDropMetaTable = errors.New("cannot drop meta table")

// Manager is responsible for holding the entire schema as well as keeping it
// in sync with the underlying data store.
type Manager struct {
	Schema *schema.Schema
	Store  kv.Store
}

func (m *Manager) NewManager(store kv.Store) (*Manager, error) {
	return &Manager{
		Store: store,
	}, nil
}

func (m *Manager) LoadSchema() error {
	cur, err := m.Store.GetRange(META_TABLE_START, META_TABLE_END)
	if err != nil {
		return fmt.Errorf("failed to read the table catalog from the store %w", err)
	}

	tablesBytes, err := cur.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read the table catalog from a cursor %w", err)
	}

	schema, err := schema.SchemaFromBytes(tablesBytes)
	if err != nil {
		return err
	}

	m.Schema = schema
	return nil
}

func (m *Manager) CreateTable(t *schema.Table) error {
	// TODO: I need a way to generate an id for this table.
	err := m.Schema.AddTable(t)
	if err != nil {
		return err
	}

	tableBytes, err := t.Value()
	if err != nil {
		return fmt.Errorf("failed encoding table while saving to store %w", err)
	}
	err = m.Store.Put(t.Key(), tableBytes)
	if err != nil {
		err = m.Schema.DropTable(t.Key())
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("could not put table definition in store %w", err)
	}
	return nil
}

func (m *Manager) DropTable(t *schema.Table) error {
	return errors.New("not implemented")
}

func (m *Manager) storeTable(t *schema.Table) error {
	tableBytes, err := t.Value()
	if err != nil {
		return fmt.Errorf("failed encoding table while saving to store %w", err)
	}
	err = m.Store.Put(t.Key(), tableBytes)
	if err != nil {
		return fmt.Errorf("could not put table definition in store %w", err)
	}
	return nil
}

func (m *Manager) storeSequence(s *schema.Sequence) error {
	sequenceBytes, err := s.Value()
	if err != nil {
		return fmt.Errorf("failed encoding sequence while saving to store %w", err)
	}
	err = m.Store.Put(s.Key(), sequenceBytes)
	if err != nil {
		return fmt.Errorf("could not put sequence definition in store %w", err)
	}
	return nil
}
