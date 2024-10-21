package catalog

import "github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"

type Column struct {
	Name     string
	DataType DataType
}

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

func (c *Column) Equal(o *Column) bool {
	if o == nil {
		return false
	}
	return c.Name == o.Name && c.DataType == o.DataType
}
