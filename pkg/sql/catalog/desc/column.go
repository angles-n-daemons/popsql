package desc

import "github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"

type Column struct {
	Name     string
	DataType DataType
	Sequence string
}

func NewColumn(name string, dt DataType) *Column {
	return &Column{
		Name:     name,
		DataType: dt,
	}
}

// SequenceColumn is a utility function which turns a name
// and a scanned token into a desc column.
func SequenceColumn(name string, tokenType scanner.TokenType) (*Column, error) {
	datatype, err := GetDataType(tokenType)
	if err != nil {
		return nil, err
	}
	return &Column{
		Name:     name,
		DataType: datatype,
	}, nil
}

func NewSequenceColumn(name string, seq string) *Column {
	return &Column{
		Name:     name,
		DataType: NUMBER,
		Sequence: seq,
	}
}

func InternalKeyColumn(seq string) *Column {
	return NewSequenceColumn(ReservedInternalColumnName, seq)
}

func (c *Column) Equal(o *Column) bool {
	if o == nil {
		return false
	}
	return c.Name == o.Name && c.DataType == o.DataType
}
