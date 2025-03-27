package catalog

import (
	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
)

// Bootstrap is the entry point for a fresh database. If nothing is in the store,
// it will automatically populate the two system tables, the meta table and the
// sequences table, and create their corresponding sequences.
// It needs to do this however using the direct functions on the desc as well
// as the store, since normal table and sequence creation depends on their existence.
// The order of operations for a bootstrap is as follows:
//   - Set the SystemObjects to be the initial values.
//   - Add the meta sequence to the desc and store.
//   - Add the sequences table sequence to the desc and store
//   - Add the meta table to the desc and store.
//   - Add the sequences table to the desc and store.
func (m *Manager) Bootstrap() error {
	metaTable := InitMetaTable()
	metaTableSequence := InitMetaTableSequence()
	sequencesTable := InitSequencesTable()
	sequencesTableSequence := InitSequencesTableSequence()
	m.Sys = &SystemObjects{
		MetaTable:              metaTable,
		MetaTableSequence:      metaTableSequence,
		SequencesTable:         sequencesTable,
		SequencesTableSequence: sequencesTableSequence,
	}

	err := m.bootstrapTable(metaTable)
	if err != nil {
		return err
	}

	err = m.bootstrapSequence(metaTableSequence)
	if err != nil {
		return err
	}

	err = m.bootstrapTable(sequencesTable)
	if err != nil {
		return err
	}

	err = m.bootstrapSequence(sequencesTableSequence)
	if err != nil {
		return err
	}

	return nil
}

// These bootstrap functions are identical to their corresponding AddTable /
// AddSequence with the exception that they do not call SequenceNext on the
// system sequences, which would increment and save values which did not yet
// exist.
func (m *Manager) bootstrapSequence(s *desc.Sequence) error {
	err := m.Schema.AddSequence(s)
	if err != nil {
		return err
	}
	err = m.StoreSequence(s)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) bootstrapTable(t *desc.Table) error {
	err := m.Schema.AddTable(t)
	if err != nil {
		return err
	}
	err = m.StoreTable(t)
	if err != nil {
		return err
	}
	return nil
}

// Below are the objects required for bootstrapping the system.
var metaTableStartKey = keys.New("1") // key string
var META_TABLE_START = metaTableStartKey.Encode()
var META_TABLE_END = metaTableStartKey.Next().Encode()

// Sequence Table Keys
var sequenceTableStartKey = keys.New("2") // key string
var SEQUENCE_TABLE_START = sequenceTableStartKey.Encode()
var SEQUENCE_TABLE_END = sequenceTableStartKey.Next().Encode()

// Meta Table Name
const MetaTableName = "__tables__"
const MetaTableID = 1
const MetaTableSequenceName = MetaTableName + "_sequence"

// Meta Sequence Name
const SequencesTableName = "__sequences__"
const SequencesTableID = 2
const SequencesTableSequenceName = SequencesTableName + "_sequence"

// Below are the sequence and table definitions for the meta and sequences tables.
// These are used to bootstrap the database when it is first created.

func InitMetaTable() *desc.Table {
	return &desc.Table{
		ID:   1,
		Name: MetaTableName,
		Columns: []*desc.Column{
			{
				Name:     "id",
				DataType: desc.NUMBER,
				Sequence: MetaTableSequenceName,
			},
			{
				Name:     "name",
				DataType: desc.STRING,
			},
		},
		PrimaryKey: []string{"id"},
	}
}

func InitSequencesTable() *desc.Table {
	return &desc.Table{
		ID:   2,
		Name: SequencesTableName,
		Columns: []*desc.Column{
			{
				Name:     "id",
				DataType: desc.NUMBER,
				Sequence: SequencesTableSequenceName,
			},
			{
				Name:     "name",
				DataType: desc.STRING,
			},
		},
		PrimaryKey: []string{"id"},
	}
}

func InitMetaTableSequence() *desc.Sequence {
	return &desc.Sequence{
		ID:   1,
		Name: MetaTableSequenceName,
		V:    2, // skip to 2 because the first two are reserved for the meta and sequences tables
	}
}

func InitSequencesTableSequence() *desc.Sequence {
	return &desc.Sequence{
		ID:   2,
		Name: SequencesTableSequenceName,
		V:    2, // skip to 3 because the first two are reserved for the meta and sequences tables
	}
}
