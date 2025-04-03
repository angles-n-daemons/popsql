package executor

import (
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/sys"
	"github.com/angles-n-daemons/popsql/pkg/sql/plan"
)

func (e *Executor) VisitCreateTable(p *plan.CreateTable) ([]Row, error) {
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
	return nil, err
}

func (e *Executor) createTableSequence(t *desc.Table) (*desc.Column, error) {
	seqName := t.TName + "_seq"
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
	return desc.NewSequenceColumn(sys.ReservedInternalKeyName, seqName), nil
}
