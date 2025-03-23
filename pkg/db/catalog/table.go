package catalog

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sys/schema/desc"
)

// TODO: this should take a statement
func (m *Manager) CreateTable(t *desc.Table) error {
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

func (m *Manager) storeTable(metaTable *desc.Table, t *desc.Table) error {
	key := metaTable.Prefix().WithID(t.Key())
	tableBytes, err := t.Value()
	if err != nil {
		return fmt.Errorf("failed encoding table while saving to store %w", err)
	}
	err = m.Store.Put(key.Encode(), tableBytes)
	if err != nil {
		return fmt.Errorf("could not put table definition in store %w", err)
	}
	return nil
}
