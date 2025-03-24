package ast

import (
	"errors"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

func Identifier(t scanner.Token) (string, error) {
	if t.Type != scanner.IDENTIFIER {
		errors.New("expected identifier token")
	}
	return t.Lexeme, nil
}
