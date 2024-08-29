package scanner_test

import (
	"testing"

	scanner "github.com/angles-n-daemons/popsql/internal/sql/parser/scanner"
)

func assertTokensEqual(t *testing.T, expected scanner.Token, actual scanner.Token) {
	if expected.Type != actual.Type {
		t.Fatalf(
			"tokens unequal, expected type %s, got %s",
			expected.Type.String(),
			actual.Type.String(),
		)
	}
	if expected.Lexeme != actual.Lexeme {
		t.Fatalf(
			"tokens unequal, expected lexeme %s, got %s",
			expected.Lexeme,
			actual.Lexeme,
		)
	}
	if expected.Literal != actual.Literal {
		t.Fatalf(
			"tokens unequal, expected lexeme %s, got %s",
			expected.Literal,
			actual.Literal,
		)
	}
}

func TestScannerBasic(t *testing.T) {
	tokens, err := scanner.Scan("SELECT 'hi', 'bye'")
	if err != nil {
		t.Fatal(err)
	}
	for i, expected := range []scanner.Token{
		scanner.SimpleToken(scanner.SELECT, "SELECT"),
		scanner.NewToken(scanner.STRING, "hi", "hi"),
		scanner.SimpleToken(scanner.COMMA, ","),
		scanner.NewToken(scanner.STRING, "bye", "bye"),
	} {
		assertTokensEqual(t, expected, tokens[i])
	}
}
