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

type DataType int

const (
	UNKNOWN DataType = iota
	STRING
	NUMBER
	BOOLEAN
)
