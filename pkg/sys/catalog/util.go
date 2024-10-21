package catalog

import (
	"fmt"
	"strings"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

const NAME_VALIDATION_REGEXP = "[A-z0-9_]"

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

func nextString(s string) string {
	i := len(s) - 1
	for i >= 0 && s[i] == 'z' {
		i--
	}

	if i == -1 {
		return s + "a"
	}

	j := 0
	return strings.Map(func(r rune) rune {
		if j == i {
			r += 1
		}
		j++
		return r
	}, s)
}
