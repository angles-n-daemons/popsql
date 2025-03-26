package catalog

import (
	"errors"
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
)

// Custom error for dropping the Meta table.
var ErrDropMetaTable = errors.New("cannot drop meta table")

// The Manager is a composite object, with a variety of responsibilities.
//
// It first and foremost acts as an in-system repository for the active
// schema.
//
// Secondarily, it's responsible for all schema changes, including the logic
// required to execute them as well as the in-memory and storage persistence
// of them.
type Manager struct {
	Sys    *SystemObjects
	Schema *schema.Schema
	Store  kv.Store
}

// SystemObjects are the objects required for the catalog to read and write
// descriptors in the system.
type SystemObjects struct {
	MetaTable              *desc.Table
	MetaTableSequence      *desc.Sequence
	SequencesTable         *desc.Table
	SequencesTableSequence *desc.Sequence
}

func NewManager(st kv.Store) (*Manager, error) {
	// Setup the manager struct.
	m := &Manager{Store: st}

	// Initialize the manager for use.
	sc, err := LoadSchema(m.Store)
	if err != nil {
		return nil, err
	}
	m.Schema = sc

	if m.Schema.Empty() {
		// if meta table does not exist, bootstrap the system tables.
		err = m.Bootstrap()
		if err != nil {
			return nil, err
		}
	} else {
		// otherwise, populate the manager with the system tables and
		// sequences.
		m.Sys, err = LoadSystemObjects(sc)
	}

	return m, nil
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

func LoadSystemObjects(sc *schema.Schema) (*SystemObjects, error) {
	// get the meta table and seqquence
	mt, ok := sc.GetTable(MetaTableName)
	if !ok {
		return nil, errors.New("meta table not found")
	}
	mts, ok := sc.GetSequence(MetaTableSequenceName)
	if !ok {
		return nil, errors.New("meta table sequence not found")
	}

	// get the sequences table and sequence
	st, ok := sc.GetTable(SequencesTableName)
	if !ok {
		return nil, errors.New("meta table not found")
	}
	sts, ok := sc.GetSequence(SequencesTableSequenceName)
	if !ok {
		return nil, errors.New("meta table sequence not found")
	}
	return &SystemObjects{
		MetaTable:              mt,
		MetaTableSequence:      mts,
		SequencesTable:         st,
		SequencesTableSequence: sts,
	}, nil
}
