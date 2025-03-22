package catalog

import (
	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
)

// System Table Keys
var metaTableStartKey = keys.New(MetaTableName)
var META_TABLE_START = metaTableStartKey.String()
var META_TABLE_END = metaTableStartKey.Next().String()

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
