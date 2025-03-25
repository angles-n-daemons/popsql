package catalog

import (
	"errors"
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema/desc"
)

// Custom error for dropping the Meta table.
var ErrDropMetaTable = errors.New("cannot drop meta table")

// Catalog is responsible for holding the entire desc as well as keeping it
// in sync with the underlying data store.
// It is also responsible for desc management operations, such as creating
// and dropping tables and sequences.
type Catalog struct {
	metaTable              *desc.Table
	metaTableSequence      *desc.Sequence
	sequencesTable         *desc.Table
	sequencesTableSequence *desc.Sequence

	Schema *schema.Schema
	Store  kv.Store
}

func NewCatalog(st kv.Store) (*Catalog, error) {
	// Load the desc from the store.
	sc, err := LoadSchema(st)
	if err != nil {
		return nil, err
	}

	// Setup the manager struct.
	m := &Catalog{Store: st, Schema: sc}

	// Initialize the manager for use.
	err = m.Init()
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Init initializes the manager by loading the schema from the store.
// If the schema does not exist in the store, it will bootstrap the system
// tables.
func (m *Catalog) Init() error {
	sc, err := LoadSchema(m.Store)
	if err != nil {
		return err
	}
	m.Schema = sc

	if m.Schema.Empty() {
		// if meta table does not exist, bootstrap the system tables.
		err = m.Bootstrap()
		if err != nil {
			return err
		}
	} else {
		// otherwise, populate the manager with the system tables and
		// sequences.
		metaTable, ok := m.Schema.GetTable(MetaTableName)
		if !ok {
			return errors.New("meta table not found")
		}
		metaTableSequence, ok := m.Schema.GetSequence(MetaTableSequenceName)
		if !ok {
			return errors.New("meta table sequence not found")
		}

		// set the sequences table and sequence
		sequencesTable, ok := m.Schema.GetTable(SequencesTableName)
		if !ok {
			return errors.New("meta table not found")
		}
		sequencesTableSequence, ok := m.Schema.GetSequence(SequencesTableSequenceName)
		if !ok {
			return errors.New("meta table sequence not found")
		}

		m.metaTable = metaTable
		m.metaTableSequence = metaTableSequence
		m.sequencesTable = sequencesTable
		m.sequencesTableSequence = sequencesTableSequence
	}
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
