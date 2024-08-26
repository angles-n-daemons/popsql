package data_test

import (
	"math/rand"
	"testing"

	"github.com/angles-n-daemons/popsql/internal/data"
)

// create a benchmark simulating likely usage

// test using a single head which is always the height of the tree
// -- report how many nodes of each height there are
// -- test alternative where there are multiple heads
// -- check insert performance
// -- check lookup performance

// test performance using generics vs not generics

// add unit tests

// question around heights
// oh if I use the approach where I always insert new head nodes if the node is smaller, does it drastically increase the size of the lists, number of nodes per level?
// I can use a random seed to figure this out

// test skiplist byte key, byte string

// what if I use int8 for height (performance for code writability)

// Things to test
// -- overwriting a value
// -- random 50 numbers, properly in order
// -- list with random values, can find
//    - elements at the head
//    - elements at the tail
//    - elements in the middle
// -- correctly returns nil for
//    - values in between values
//    - values before the list
//    - values at the end of the list

// Failure modes
// -- Puting into list with no value

// NON GENERIC
// cpu: Intel(R) Core(TM) i5-8257U CPU @ 1.40GHz
// BenchmarkSkiplistPerformance-8           1000000              3149 ns/op
// BenchmarkSkiplistReadHeavy-8              738738              5487 ns/op
// BenchmarkSkiplistWriteHeavy-8             522088              7030 ns/op
// BenchmarkSkiplistReadHits-8              1000000              1894 ns/op
// BenchmarkSkiplistReadMisses-8            1000000              2005 ns/op
//
// GENERIC
// BenchmarkSkiplistPerformance-8           1000000              3108 ns/op
// BenchmarkSkiplistReadHeavy-8              784930              5589 ns/op
// BenchmarkSkiplistWriteHeavy-8             477907              6719 ns/op
// BenchmarkSkiplistReadHits-8              1000000              1956 ns/op
// BenchmarkSkiplistReadMisses-8            1000000              1974 ns/op

func BenchmarkSkiplistPerformance(b *testing.B) {
	list := data.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		val := rng.Intn(1000000)
		list.Put(val, val)
		val = rng.Intn(1000000)
		list.Get(val)
	}
}

func BenchmarkSkiplistReadHeavy(b *testing.B) {
	list := data.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		val := rng.Intn(1000000)
		list.Put(val, val)
		for j := 0; j < 3; j++ {
			val = rng.Intn(1000000)
			list.Get(val)
		}
	}
}

func BenchmarkSkiplistWriteHeavy(b *testing.B) {
	list := data.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		var val int
		for j := 0; j < 3; j++ {
			val = rng.Intn(1000000)
			list.Put(val, val)
		}
		val = rng.Intn(1000000)
		list.Get(val)
	}
}

func BenchmarkSkiplistReadHits(b *testing.B) {
	list := data.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		val := rng.Intn(1000000)
		list.Put(val, val)
		list.Get(val)
	}
}

func BenchmarkSkiplistReadMisses(b *testing.B) {
	list := data.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		val := rng.Intn(1000000)
		list.Put(val, val)
		val = rng.Intn(1000000) + 1000000
		list.Get(val)
	}
}
