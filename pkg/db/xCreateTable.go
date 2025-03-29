package db

import (
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
)

func (e *Engine) CreateTable(stmt *ast.CreateTable) error {
	dt, err := NewTableFromStmt(stmt)
	if err != nil {
		return err
	}

	if dt.PrimaryKey == nil || len(dt.PrimaryKey) == 0 {
		seq, err := dt.AddInternalPrimaryKey()
		if err != nil {
			return err
		}
		_, err = catalog.Create(e.Catalog, seq)
		if err != nil {
			return err
		}
	}
	_, err = catalog.Create(e.Catalog, dt)
	return err
}

// NewTableFromStmt creates a new table from a create statement.
// The table will NOT have an ID to start, as it will be assigned
// by the catalog when the table is created.
func NewTableFromStmt(stmt *ast.CreateTable) (*desc.Table, error) {
	columns := make([]*desc.Column, len(stmt.Columns))

	for i, colSpec := range stmt.Columns {
		column, err := desc.NewColumnFromStmt(colSpec)
		if err != nil {
			return nil, err
		}
		columns[i] = column
	}
	// TODO: primary key parsing
	// TODO: validate primary key
	return desc.NewTable(stmt.Name.Lexeme, columns, []string{})
}
