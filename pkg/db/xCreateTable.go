package db

import (
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
)

var ReservedInternalKeyName = "__key"

func (e *Engine) CreateTable(stmt *ast.CreateTable) error {
	dt, err := NewTableFromStmt(stmt)
	if err != nil {
		return err
	}

	if dt.PrimaryKey == nil || len(dt.PrimaryKey) == 0 {
		pkeyCol, err := e.createTableSequence(dt)
		if err != nil {
			return err
		}

		dt.Columns = append(dt.Columns, pkeyCol)
		dt.PrimaryKey = []string{ReservedInternalKeyName}
	}

	dt.TID, err = catalog.NextID(e.Catalog, dt)
	if err != nil {
		return err
	}
	err = catalog.Create(e.Catalog, dt)
	return err
}

func (e *Engine) createTableSequence(t *desc.Table) (*desc.Column, error) {
	seqName := t.TName + "_seq"
	seq := desc.NewSequence(seqName)
	id, err := catalog.NextID(e.Catalog, seq)
	seq.SID = id
	if err != nil {
		return nil, err
	}

	err = catalog.Create(e.Catalog, seq)
	if err != nil {
		return nil, err
	}
	return desc.NewSequenceColumn(ReservedInternalKeyName, seqName), nil
}

// NewTableFromStmt creates a new table from a create statement.
// The table will NOT have an ID to start, as it will be assigned
// by the catalog when the table is created.
func NewTableFromStmt(stmt *ast.CreateTable) (*desc.Table, error) {
	columns := make([]*desc.Column, len(stmt.Columns))

	for i, colSpec := range stmt.Columns {
		column, err := NewColumnFromStmt(colSpec)
		if err != nil {
			return nil, err
		}
		columns[i] = column
	}
	// TODO: primary key parsing
	// TODO: validate primary key
	return desc.NewTable(stmt.Name.Lexeme, columns, []string{})
}

// NewColumnFromStmt is a utility function which turns a ColumnSpec into a desc.
func NewColumnFromStmt(col *ast.ColumnSpec) (*desc.Column, error) {
	// TODO: error handling:
	//  - check name type
	dt, err := desc.GetDataType(col.DataType.Type)
	if err != nil {
		return nil, err
	}

	name, err := ast.Identifier(col.Name)
	if err != nil {
		return nil, err
	}
	return desc.NewColumn(name, dt), nil
}
