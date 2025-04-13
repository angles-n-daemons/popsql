package execution

import (
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/plan"
)

func (e *Executor) VisitCreateTable(p *plan.CreateTable) (Row, error) {
	if e.State.tableCreated {
		return nil, nil
	}

	dt := p.Table
	if dt.PrimaryKey == nil || len(dt.PrimaryKey) == 0 {
		pkeyCol, err := e.createTableSequence(dt)
		if err != nil {
			return nil, err
		}

		dt.Columns = append(dt.Columns, pkeyCol)
		dt.PrimaryKey = []string{pkeyCol.Name}
	}

	id, err := catalog.NextDescriptorID(e.Catalog, dt)
	if err != nil {
		return nil, err
	}
	dt.TID = id
	err = catalog.Create(e.Catalog, dt)

	// set the state variable to prevent re-creating the table.
	e.State.tableCreated = true
	return Row{dt.Name()}, err
}

func (e *Executor) createTableSequence(t *desc.Table) (*desc.Column, error) {
	seqName := t.DefaultSequenceName()
	seq := desc.NewSequence(seqName)
	id, err := catalog.NextDescriptorID(e.Catalog, seq)
	seq.SID = id
	if err != nil {
		return nil, err
	}

	err = catalog.Create(e.Catalog, seq)
	if err != nil {
		return nil, err
	}
	return desc.InternalKeyColumn(seqName), nil
}
