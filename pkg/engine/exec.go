package engine

import (
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/sys"
)

func (db *Engine) PrintRows(rows []any) {

}

func (db *Engine) VisitSelectStmt(stmt *ast.Select) (*any, error) {
	return nil, nil
}
func (db *Engine) VisitInsertStmt(stmt *ast.Insert) (*any, error) { return nil, nil }
func (db *Engine) VisitCreateStmt(stmt *ast.Create) (*any, error) {
	columns := []sys.Column{}
	name := stmt.Name.Lexeme
	for _, column := range stmt.Columns {
		dataType, err := sys.GetDataType(column.DataType)
		if err != nil {
			return nil, err
		}
		column := sys.Column{
			Space:    sys.USER,
			Table:    name,
			Name:     column.Name.Lexeme,
			DataType: dataType,
		}
		columns = append(columns, column)
	}
	db.CreateTable(sys.Table{Space: sys.USER, Name: name, Columns: columns})
	return nil, nil
}
