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
// It is also responsible for schema management operations, such as creating
// and dropping tables and sequences.
type Manager struct {
	metaTable         *schema.Table
	metaTableSequence *schema.Sequence
	Schema            *schema.Schema
	Store             kv.Store
}

func NewManager(st kv.Store) (*Manager, error) {
	// Load the schema from the store.
	sc, err := LoadSchema(st)
	if err != nil {
		return nil, err
	}

	// Setup the manager struct.
	m := &Manager{Store: st, Schema: sc}

	// Initialize the manager for use.
	err = m.Init()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) Init() error {
	sc, err := LoadSchema(m.Store)
	if err != nil {
		return err
	}
	m.Schema = sc

	// if meta table does not exist, bootstrap the system tables.
	if m.Schema.Empty() {
		err = m.Bootstrap()
		if err != nil {
			return err
		}
	}

	// set the meta table and meta table sequence.
	metaTable, ok := m.Schema.GetTable(MetaTableName)
	if !ok {
		return errors.New("meta table not found")
	}
	metaTableSequence, ok := m.Schema.GetSequence(MetaTableSequenceName)
	if !ok {
		return errors.New("meta table sequence not found")
	}

	m.metaTable = metaTable
	m.metaTableSequence = metaTableSequence
	return nil
}

func LoadSchema(st kv.Store) (*schema.Schema, error) {
	cur, err := st.GetRange(META_TABLE_START, META_TABLE_END)
	if err != nil {
		return nil, fmt.Errorf("failed to read the table catalog from the store %w", err)
	}

	tablesBytes, err := cur.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read the table catalog from a cursor %w", err)
	}

	cur, err = st.GetRange(SEQUENCE_TABLE_START, SEQUENCE_TABLE_END)
	if err != nil {
		return nil, fmt.Errorf("failed to read the table catalog from the store %w", err)
	}

	sequencesBytes, err := cur.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read the table catalog from a cursor %w", err)
	}

	sc := schema.NewSchema()
	err = sc.LoadTables(tablesBytes)
	if err != nil {
		return nil, err
	}

	err = sc.LoadSequences(sequencesBytes)
	if err != nil {
		return nil, err
	}

	return sc, nil
}
