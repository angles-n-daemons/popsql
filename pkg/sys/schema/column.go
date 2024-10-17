package schema

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

type Column struct {
	Name     string
	DataType DataType
}

type ColumnRegister struct {
	Schema Schema
	Table  Table
	Column Column
}

func (c *ColumnRegister) Key() (string, error) {
	return "", nil
}

func (c *ColumnRegister) Value() ([]byte, error) {
	return nil, nil
}

type DataType int

const (
	UNKNOWN DataType = iota
	STRING
	NUMBER
	BOOLEAN
)

func GetDataType(token *scanner.Token) (DataType, error) {
	switch token.Type {
	case scanner.DATATYPE_BOOLEAN:
		return BOOLEAN, nil
	case scanner.DATATYPE_STRING:
		return STRING, nil
	case scanner.DATATYPE_NUMBER:
		return NUMBER, nil
	default:
		return UNKNOWN, fmt.Errorf("unrecognized data type %s", token.Type)
	}
}
