package plan

import (
	"errors"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
)

type Planner struct{}

func PlanQuery(stmt ast.Stmt) (Plan, error) {
	planner := &Planner{}
	return ast.VisitStmt(stmt, planner)
}

func (p *Planner) VisitCreateTableStmt(stmt *ast.CreateTable) (Plan, error) {
	dt, err := NewTableFromStmt(stmt)
	if err != nil {
		return nil, err
	}

	return &CreateTable{
		Table: dt,
	}, nil
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

func (p *Planner) VisitInsertStmt(stmt *ast.Insert) (Plan, error) {
	return nil, errors.New("not implemented")
}

func (p *Planner) VisitSelectStmt(stmt *ast.Select) (Plan, error) {
	return nil, errors.New("not implemented")
}
