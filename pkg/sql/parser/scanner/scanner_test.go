package scanner

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func athousandrandomtokens() string {
	rng := rand.New(rand.NewSource(0))
	keys := []string{}
	poem := make([]string, 1000)
	for key := range keywordLookup {
		keys = append(keys, key)
	}
	for i := 0; i < 1000; i++ {
		poem[i] = keys[rng.Intn(len(keys))]
	}
	return strings.Join(poem, " ")
}

func mixCase(s string) string {
	rng := rand.New(rand.NewSource(0))
	coin := rng.Intn(2)
	middle := len(s) / 2
	if coin == 1 {
		return s[:middle] + strings.ToLower(s[middle:])
	} else {
		return strings.ToLower(s[:middle])
	}
}

var tokenpoem = athousandrandomtokens()

func assertTokensEqual(t *testing.T, expected *Token, actual *Token) {
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
			"tokens unequal, expected literal %s, got %s",
			expected.Literal,
			actual.Literal,
		)
	}
}

func TestScannerBasic(t *testing.T) {
	tokens, err := Scan("SELECT 'hi', 'bye'")
	if err != nil {
		t.Fatal(err)
	}
	for i, expected := range []*Token{
		simpleToken(SELECT, "SELECT"),
		newToken(STRING, "hi", "hi"),
		simpleToken(COMMA, ","),
		newToken(STRING, "bye", "bye"),
	} {
		assertTokensEqual(t, expected, tokens[i])
	}
}

func TestScannerWithDataTypes(t *testing.T) {
	tokens, err := Scan("col1 INT")
	if err != nil {
		t.Fatal(err)
	}
	for i, expected := range []*Token{
		newToken(IDENTIFIER, "col1", "col1"),
		simpleToken(TYPE, "INT"),
	} {
		fmt.Println(expected)
		assertTokensEqual(t, expected, tokens[i])
	}
}

// BenchmarkScanIfStatements
// BenchmarkScanIfStatements-11                 354           3179018 ns/op
// BenchmarkScanWithMap
// BenchmarkScanWithMap-11                      564           2066271 ns/op
// BenchmarkScanWithTrie
// BenchmarkScanWithTrie-11                     580           2080027 ns/op
// BenchmarkScanSequential
// BenchmarkScanSequential-11                 21890             54766 ns/op
// PASS

// pkg: github.com/angles-n-daemons/popsql/internal/sql/parser/scanner
// BenchmarkScanTokenPoem
// BenchmarkScanTokenPoem-11          15246             76651 ns/op

func BenchmarkScanTokenPoem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Scan(tokenpoem)
	}
}

// Benchmark case insensitivity
// BenchmarkScanTokenPoem-11           7710            155490 ns/op case sensitive
// BenchmarkScanTokenPoem-11           6512            178254 ns/op case insensitive
