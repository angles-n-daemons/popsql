package sys

import (
	"errors"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
)

// The sys package is a small utility package for managing the system
// tables in popsql. It provides constant values for reference, like the
// names of the system objects, as well as utilities for constructing the
// meta containers, to be used by the catalog manager.
type SystemTable struct {
	Table    *desc.Table
	Sequence *desc.Sequence
}

type SystemSchema struct {
	Tables    SystemTable
	Sequences SystemTable
}

func (m *SystemSchema) Objects() []any {
	return []any{
		m.Tables.Table,
		m.Tables.Sequence,
		m.Sequences.Table,
		m.Sequences.Sequence,
	}
}

const (
	// These ids are reserved for the system tables.
	// They are also generated on bootup, but these values are required
	// to read from the meta tables on bootup.
	TablesID    = 1
	SequencesID = 2

	Tables    = "__tables__"
	Sequences = "__sequences___"

	TablesSequence    = Tables + "_sequence"
	SequencesSequence = Sequences + "_sequence"

	idCol   = "id"
	nameCol = "name"
)

func InitSystemTable(id uint64, name, seqName string) *desc.Table {
	return &desc.Table{
		TID:   id,
		TName: name,
		Columns: []*desc.Column{
			desc.NewSequenceColumn(idCol, seqName),
			desc.NewColumn(nameCol, desc.STRING),
		},
		PrimaryKey: []string{idCol},
	}
}

func InitSystemMeta() *SystemSchema {
	metaTable := InitSystemTable(TablesID, Tables, TablesSequence)
	sequenceTable := InitSystemTable(SequencesID, Sequences, SequencesSequence)
	return &SystemSchema{
		Tables: SystemTable{
			Table:    metaTable,
			Sequence: desc.NewSequenceFromArgs(TablesID, TablesSequence, 2),
		},
		Sequences: SystemTable{
			Table:    sequenceTable,
			Sequence: desc.NewSequenceFromArgs(SequencesID, SequencesSequence, 2),
		},
	}
}

// FromSchema creates a new Meta from the given schema.
// It does this by reading the meta objects from the collections contained
// within the schema, erroring if any of the meta objects are missing.
func FromSchema(sc *schema.Schema) (*SystemSchema, error) {
	tables := schema.GetByName[*desc.Table](sc, Tables)
	sequences := schema.GetByName[*desc.Table](sc, Sequences)
	tablesSequence := schema.GetByName[*desc.Sequence](sc, TablesSequence)
	sequencesSequence := schema.GetByName[*desc.Sequence](sc, SequencesSequence)

	for _, v := range []any{tables, sequences, tablesSequence, sequencesSequence} {
		if v == nil {
			return nil, errors.New("could not find system object in schema")
		}
	}
	return &SystemSchema{
		Tables: SystemTable{
			Table:    tables,
			Sequence: tablesSequence,
		},
		Sequences: SystemTable{
			Table:    sequences,
			Sequence: sequencesSequence,
		},
	}, nil
}
