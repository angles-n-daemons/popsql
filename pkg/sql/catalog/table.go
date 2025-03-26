package catalog

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
)

func (m *Manager) CreateTable(t *desc.Table) error {
	// create an id for the new table.
	id, err := m.SequenceNext(m.Sys.MetaTableSequence)
	if err != nil {
		return err
	}
	t.ID = id

	if len(t.PrimaryKey) == 0 {
		s, err := t.AddInternalPrimaryKey()
		if err != nil {
			return err
		}
		s, err = m.CreateSequence(s)
		if err != nil {
			return err
		}
	}

	// save the table in the memory schema.
	err = m.Schema.AddTable(t)
	if err != nil {
		return err
	}

	// attempt to store the table.
	err = m.StoreTable(t)
	if err != nil {
		// if storing the table fails, back the changes out of the schema.
		err = m.Schema.RemoveTable(t.Key())
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("could not put table descriptor into the store %w", err)
	}
	return nil
}

func (m *Manager) StoreTable(t *desc.Table) error {
	key := m.Sys.MetaTable.Prefix().WithID(t.Key())
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
