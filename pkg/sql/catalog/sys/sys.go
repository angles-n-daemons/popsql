package sys

import (
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
)

// These ids are reserved for the system tables. They are also
// generated on bootup, but these values are required to read
// from the meta tables on bootup.
const (
	TablesID    = 1
	SequencesID = 2

	Tables    = "__tables__"
	Sequences = "__sequences__"

	TablesSequence    = Tables + "_seq"
	SequencesSequence = Sequences + "_seq"

	idCol   = "id"
	nameCol = "name"
)

func InitSchema() (*schema.Schema, error) {
	sc := schema.NewSchema()
	metaTable, metaTableSeq := InitSystemTable(TablesID, Tables, TablesSequence)

	// setup the sequence table with the value column.
	sequenceTable, sequenceTableSeq := InitSystemTable(SequencesID, Sequences, SequencesSequence)
	sequenceTable.Columns = append(sequenceTable.Columns, desc.NewColumn("value", desc.STRING))
	for _, tab := range []*desc.Table{metaTable, sequenceTable} {
		err := schema.Add(sc, tab)
		if err != nil {
			return nil, err
		}

	}
	for _, seq := range []*desc.Sequence{metaTableSeq, sequenceTableSeq} {
		err := schema.Add(sc, seq)
		if err != nil {
			return nil, err
		}
	}

	return sc, nil
}

func InitSystemTable(id uint64, name, seqName string) (*desc.Table, *desc.Sequence) {
	tab := &desc.Table{
		TID:   id,
		TName: name,
		Columns: []*desc.Column{
			desc.NewSequenceColumn(idCol, seqName),
			desc.NewColumn(nameCol, desc.STRING),
		},
		PrimaryKey: []string{idCol},
	}
	seq := desc.NewSequenceFromArgs(id, tab.DefaultSequenceName(), 2)
	return tab, seq
}
