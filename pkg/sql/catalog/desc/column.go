package desc

import (
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

type Column struct {
	Name     string
	DataType DataType
	Sequence string
}

// NewColumn is a utility function which turns a name and a scanned token into
// a desc column.
func NewColumn(name string, tokenType scanner.TokenType) (*Column, error) {
	datatype, err := GetDataType(tokenType)
	if err != nil {
		return nil, err
	}
	return &Column{
		Name:     name,
		DataType: datatype,
	}, nil
}

func NewSequenceColumn(name string, dt DataType, seq string) *Column {
	return &Column{
		Name:     name,
		DataType: dt,
		Sequence: seq,
	}
}

// NewColumnFromStmt is a utility function which turns a ColumnSpec into a desc.
func NewColumnFromStmt(col *ast.ColumnSpec) (*Column, error) {
	// TODO: error handling:
	//  - check name type
	return NewColumn(col.Name.Lexeme, col.DataType.Type)
}

func (c *Column) Equal(o *Column) bool {
	if o == nil {
		return false
	}
	return c.Name == o.Name && c.DataType == o.DataType
}
