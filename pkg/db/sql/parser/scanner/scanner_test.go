package scanner

import (
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
	return strings.Join(poem, ` `)
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
