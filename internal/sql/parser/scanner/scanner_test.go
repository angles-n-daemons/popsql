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

var tokenpoem = athousandrandomtokens()

func assertTokensEqual(t *testing.T, expected Token, actual Token) {
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
	tokens, err := Scan("SELECT 'hi', 'bye'")
	if err != nil {
		t.Fatal(err)
	}
	for i, expected := range []Token{
		SimpleToken(SELECT, "SELECT"),
		NewToken(STRING, "hi", "hi"),
		SimpleToken(COMMA, ","),
		NewToken(STRING, "bye", "bye"),
	} {
		assertTokensEqual(t, expected, tokens[i])
	}
}

func TestScanningMethods(t *testing.T) {
	ergnorelen := func(tokens []Token, err error) int {
		fmt.Println(err)
		return len(tokens)
	}
	fmt.Println("Scan", ergnorelen(Scan(tokenpoem)))
	fmt.Println("ScanWithMap", ergnorelen(ScanWithMap(tokenpoem)))
	fmt.Println("ScanWithTrie", ergnorelen(ScanWithTrie(tokenpoem)))
	keywordtrie.walk(0)

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

func BenchmarkScanIfStatements(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Scan(tokenpoem)
	}
}

func BenchmarkScanWithMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanWithMap(tokenpoem)
	}
}

func BenchmarkScanWithTrie(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanWithTrie(tokenpoem)
	}
}

func BenchmarkScanSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanSequential(tokenpoem)
	}
}
