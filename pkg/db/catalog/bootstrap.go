package catalog

import (
	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
)

// Bootstrap is the entry point for a fresh database. If nothing is in the store,
// it will automatically populate the two system tables, the meta table and the
// sequences table, and create their corresponding sequences.
// It needs to do this however using the direct functions on the schema as well
// as the store, since normal table and sequence creation depends on their existence.
// The order of operations for a bootstrap is as follows:
//   - Add the meta sequence to the schema and store.
//   - Add the sequences table sequence to the schema and store
//   - Add the meta table to the schema and store.
//   - Add the sequences table to the schema and store.
func (m *Manager) Bootstrap() error {
	err := m.bootstrapSequence(InitMetaTableSequence)
	if err != nil {
		return err
	}

	err = m.bootstrapSequence(InitSequencesTableSequence)
	if err != nil {
		return err
	}

	err = m.bootstrapTable(InitMetaTable)
	if err != nil {
		return err
	}

	err = m.bootstrapTable(InitSequencesTable)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) bootstrapSequence(s *schema.Sequence) error {
	err := m.Schema.AddSequence(s)
	if err != nil {
		return err
	}
	err = m.storeSequence(InitSequencesTable, s)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) bootstrapTable(t *schema.Table) error {
	err := m.Schema.AddTable(t)
	if err != nil {
		return err
	}
	err = m.storeTable(InitMetaTable, t)
	if err != nil {
		return err
	}
	return nil
}

// Below are the objects required for bootstrapping the system.
var metaTableStartKey = keys.New(MetaTableName)
var META_TABLE_START = metaTableStartKey.Encode()
var META_TABLE_END = metaTableStartKey.Next().Encode()

// Sequence Table Keys
var sequenceTableStartKey = keys.New(SequencesTableName)
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

var InitMetaTable = &schema.Table{
	ID:   1,
	Name: MetaTableName,
	Columns: []*schema.Column{
		{
			Name:     "id",
			DataType: schema.NUMBER,
			Sequence: 1,
		},
		{
			Name:     "name",
			DataType: schema.STRING,
		},
	},
	PrimaryKey: []string{"id"},
}

var InitSequencesTable = &schema.Table{
	ID:   2,
	Name: SequencesTableName,
	Columns: []*schema.Column{
		{
			Name:     "id",
			DataType: schema.NUMBER,
			Sequence: 2,
		},
		{
			Name:     "name",
			DataType: schema.STRING,
		},
	},
	PrimaryKey: []string{"id"},
}

var InitMetaTableSequence = &schema.Sequence{
	ID:   1,
	Name: MetaTableSequenceName,
	V:    3, // skip to 3 because the first two are reserved for the meta and sequences tables
}

var InitSequencesTableSequence = &schema.Sequence{
	ID:   2,
	Name: SequencesTableSequenceName,
	V:    3, // skip to 3 because the first two are reserved for the meta and sequences tables
}

func MetaTableKey(t *schema.Table) string {
	return keys.New(MetaTableName).WithID(t.Key()).Encode()
}

func SequenceKey(s *schema.Sequence) string {
	return keys.New(SequencesTableName).WithID(s.Key()).Encode()
}
