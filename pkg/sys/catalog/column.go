package catalog

import (
	"encoding/json"
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

type Column struct {
	Table    *Table
	Name     string
	DataType DataType
}

func (c *Column) ToRegister() (*ColumnRegister, error) {
	table, err := c.Table.ToRegister().Key()
	if err != nil {
		return nil, err
	}
	return &ColumnRegister{
		Table:    table,
		Name:     c.Name,
		DataType: c.DataType,
	}, nil
}

type ColumnRegister struct {
	Table    string
	Name     string
	DataType DataType
}

func (c *ColumnRegister) Key() (string, error) {
	return c.Table + "-" + c.Name, nil
}

func (c *ColumnRegister) Value() ([]byte, error) {
	return json.Marshal(c)
}

type DataType int

const (
	UNKNOWN DataType = iota
	STRING
	NUMBER
	BOOLEAN
)

func GetDataType(tokenType scanner.TokenType) (DataType, error) {
	switch tokenType {
	case scanner.DATATYPE_BOOLEAN:
		return BOOLEAN, nil
	case scanner.DATATYPE_STRING:
		return STRING, nil
	case scanner.DATATYPE_NUMBER:
		return NUMBER, nil
	default:
		return UNKNOWN, fmt.Errorf("unrecognized data type %s", tokenType)
	}
}
