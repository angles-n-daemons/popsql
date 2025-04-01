package ast

import (
	"errors"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
)

// Identifier returns the identifier from a token.
// To be used wherever tokens are expected to be strings, sql
// execution for example.
func Identifier(t scanner.Token) (string, error) {
	if t.Type != scanner.IDENTIFIER {
		return "", errors.New("expected identifier token")
	}
	return t.Lexeme, nil
}
