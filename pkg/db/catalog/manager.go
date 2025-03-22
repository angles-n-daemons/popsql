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

func NewManager(store kv.Store) *Manager {
	return &Manager{
		Store: store,
	}
}

func (m *Manager) Init() error {
	err := m.LoadSchema()
	if err != nil {
		return err
	}
	// if meta table does not exist, bootstrap the system tables
	if _, ok := m.Schema.GetTable(MetaTableName); !ok {
		err = m.Bootstrap()
		if err != nil {
			return err
		}
	}
	return nil
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

	sc := schema.NewSchema()
	err = sc.LoadTables(tablesBytes)
	if err != nil {
		return err
	}

	cur, err = m.Store.GetRange(SEQUENCE_TABLE_START, SEQUENCE_TABLE_END)
	if err != nil {
		return fmt.Errorf("failed to read the table catalog from the store %w", err)
	}

	sequencesBytes, err := cur.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read the table catalog from a cursor %w", err)
	}
	err = sc.LoadSequences(sequencesBytes)
	if err != nil {
		return err
	}

	m.Schema = sc
	return nil
}

// TODO: this should take a statement
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

func (m *Manager) storeTable(metaTable *schema.Table, t *schema.Table) error {
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

func (m *Manager) storeSequence(sequenceTable *schema.Table, s *schema.Sequence) error {
	key := sequenceTable.Prefix().WithID(s.Key())
	sequenceBytes, err := s.Value()
	if err != nil {
		return fmt.Errorf("failed encoding sequence while saving to store %w", err)
	}
	err = m.Store.Put(key.Encode(), sequenceBytes)
	if err != nil {
		return fmt.Errorf("could not put sequence definition in store %w", err)
	}
	return nil
}
