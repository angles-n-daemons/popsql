package catalog

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

const NAME_VALIDATION_REGEX = "[A-z0-9_]"

func ValidateName(name string) error {
	return nil
}

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
