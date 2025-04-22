package plan

import (
	"errors"
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/db/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/catalog/schema"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/parser/scanner"
)

var Debug = false

type Planner struct {
	Schema *schema.Schema
}

func PlanQuery(sc *schema.Schema, stmt ast.Stmt) (Plan, error) {
	planner := &Planner{Schema: sc}
	plan, err := ast.VisitStmt(stmt, planner)
	if err != nil {
		return nil, err
	}
	if Debug {
		planStr, err := DebugPlan(plan)
		if err != nil {
			return nil, err
		}
		fmt.Println(planStr)
	}
	return plan, err
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
	return desc.NewTable(stmt.Name.Name.Lexeme, columns, []string{})
}

// NewColumnFromStmt is a utility function which turns a ColumnSpec into a desc.
func NewColumnFromStmt(col *ast.ColumnSpec) (*desc.Column, error) {
	// TODO: error handling:
	//  - check name type
	dt, err := desc.GetDataType(col.DataType.Type)
	if err != nil {
		return nil, err
	}

	name := col.Name.Name.Lexeme
	return desc.NewColumn(name, dt), nil
}

func (p *Planner) VisitInsertStmt(stmt *ast.Insert) (Plan, error) {
	tname := stmt.Table.Name.Lexeme

	dt := schema.GetByName[*desc.Table](p.Schema, tname)
	if dt == nil {
		return nil, fmt.Errorf("Could not find table with name %s", tname)
	}

	columns := make([]*desc.Column, len(stmt.Columns))
	if len(stmt.Columns) == 0 {
		// if no columns supplied, just use the table's columns
		columns = dt.Columns
	}
	for i, col := range stmt.Columns {
		columns[i] = dt.GetColumn(col.Name.Lexeme)
		if columns[i] == nil {
			return nil, fmt.Errorf("Could not find column in table %s with name %s", tname, col.Name.Lexeme)
		}
	}

	inputLen := len(columns)

	if inputLen == 0 {
		return nil, fmt.Errorf("Nothing to insert")
	}

	for i, tuple := range stmt.Values {
		if len(tuple) != inputLen {
			return nil, fmt.Errorf("Tuple %d has %d values, but %d columns were specified", i, len(tuple), inputLen)
		}
	}

	values := NewValues(stmt.Values)

	return NewInsert(dt, columns, values), nil
}

func (p *Planner) VisitSelectStmt(stmt *ast.Select) (Plan, error) {
	tname := stmt.From.Name.Lexeme

	dt := schema.GetByName[*desc.Table](p.Schema, tname)
	if dt == nil {
		return nil, fmt.Errorf("Could not find table with name %s", tname)
	}

	from := NewScan(dt)

	if len(stmt.Terms) == 1 && stmt.Terms[0].Name.Type == scanner.STAR {
		return from, nil
	}
	return from, errors.New("not implemented")
}
