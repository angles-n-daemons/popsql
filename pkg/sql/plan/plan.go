package plan

import (
	"crypto/rand"
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
)

type PlanVisitor[T any] interface {
	VisitCreateTable(*CreateTable) (T, error)
	VisitInsert(*Insert) (T, error)
	VisitScan(*Scan) (T, error)
	VisitValues(*Values) (T, error)
}

func VisitPlan[T any](plan Plan, visitor PlanVisitor[T]) (T, error) {
	switch typedPlan := plan.(type) {
	case *CreateTable:
		return visitor.VisitCreateTable(typedPlan)
	case *Insert:
		return visitor.VisitInsert(typedPlan)
	case *Scan:
		return visitor.VisitScan(typedPlan)
	case *Values:
		return visitor.VisitValues(typedPlan)
	default:
		return *new(T), fmt.Errorf("Could not match plan of type %T", plan)
	}
}

type Plan interface {
	Columns() []string
}

type CreateTable struct {
	Table *desc.Table
}

func (p *CreateTable) Columns() []string { return []string{"table"} }

type Insert struct {
	Table  *desc.Table
	Cols   []*desc.Column
	Source *Values
	Offset int
}

func NewInsert(t *desc.Table, cols []*desc.Column, values *Values) *Insert {
	return &Insert{
		Table:  t,
		Cols:   cols,
		Source: values,
	}
}

func (p *Insert) Columns() []string {
	return []string{"id"}
}

type Scan struct {
	// Scan nodes need an ID so that multiple cursors to
	// the same table can be identified during execution.

	// The  ID itself only needs to be unique to the plan,
	// not globally in the process.
	ID    string
	Table *desc.Table
}

func NewScan(t *desc.Table) *Scan {
	return &Scan{ID: randomString(8), Table: t}
}

func (p *Scan) Columns() []string {
	columns := p.Table.GetColumns()
	cols := make([]string, len(columns))
	for i, col := range columns {
		cols[i] = col.Name
	}
	return cols
}

func randomString(length int) string {
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

func NewValues(rows [][]ast.Expr) *Values {
	return &Values{
		ID:   randomString(8),
		Rows: rows,
	}
}

type Values struct {
	// Values nodes are used to represent tuples in a variety
	// of statements. They, like the scan nodes, require an
	// id to keep track of the execution offset in reading
	// them.
	ID   string
	Rows [][]ast.Expr
}

func (p *Values) Columns() []string {
	columns := make([]string, len(p.Rows[0]))
	for i := range columns {
		columns[i] = fmt.Sprintf("col%d", i)
	}
	return columns
}
