package memtable_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/angles-n-daemons/popsql/internal/data/memtable"
)

const (
	randTo0 = 0
	randTo1 = -1
)

// source for which rand.Rand.Intn(2) returns 1
// height - 1 times before each 0
// useful for setting the height of a skiplist node manually
type mockHeightRandSource struct {
	height   int
	numCalls int
}

func (m *mockHeightRandSource) Int63() int64 {
	m.numCalls++
	if m.height == 0 {
		return randTo0
	}

	if m.numCalls%m.height == 0 {
		return randTo0
	} else {
		return randTo1
	}
}
func (m *mockHeightRandSource) Seed(seed int64) {}

func (m *mockHeightRandSource) changeHeight(height int) {
	m.height = height
	m.numCalls = 0
}

func TestMockRandHeight(t *testing.T) {
	for i := 1; i < 10; i++ {
		height := 1
		rng := rand.New(&mockHeightRandSource{height: i})
		for rng.Intn(2) == 1 {
			height++
			if height > 50 {
				break
			}
		}
		if i != height {
			t.Fatalf(
				"expected random height generator to return %d - 1 1s, but got %d",
				i,
				height,
			)
		}
	}
}

func skiplistFromArray(vals [][]int) (*memtable.Skiplist[int, int], error) {
	list := memtable.NewSkiplist[int, int]()
	for _, val := range vals {
		_, err := list.Put(val[0], val[1])
		if err != nil {
			return nil, fmt.Errorf(
				"list raised an error on Put: %v",
				err,
			)
		}
	}
	return list, nil
}

func assertListsEquivalent(t *testing.T, vals [][]int, list *memtable.Skiplist[int, int]) {
	for _, val := range vals {
		node := list.Get(val[0])
		if node == nil {
			t.Fatalf("expected Get to find key %d in list", val[0])
		}
		if node.Val != val[1] {
			t.Fatalf(
				"expected val %d for key %d on Get, but got %d",
				val[1], val[0], node.Val,
			)
		}
	}
}

func assertDeleteFromList(t *testing.T, vals [][]int, list *memtable.Skiplist[int, int]) {
	// check elements deleted
	for _, val := range vals {
		node := list.Get(val[0])
		if node != nil {
			t.Fatalf("expected Get to return nil on after deletion")
		}
		node = list.Delete(val[0])
		if node != nil {
			t.Fatalf("expected Delete to return nil on second attempt")
		}
	}
}

// helper function which asserts that the size of the vals array matches the
// skiplist, and that the elements found in the vals array can be found
// in the skiplist, and then deletes the elements when finished
func assertPutGetDelete(t *testing.T, vals [][]int) {
	list, err := skiplistFromArray(vals)
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != int(list.Size) {
		t.Fatalf(
			"expected list length to match vals %d but got length %d",
			len(vals), list.Size,
		)
	}
	assertListsEquivalent(t, vals, list)
	assertDeleteFromList(t, vals, list)
	if list.Size != 0 {
		t.Fatalf(
			"expected list to be empty after deleting elements, but got size %d",
			list.Size,
		)
	}
}

// test very simple skiplist use cases
func TestSkiplistBasic(t *testing.T) {
	vals := [][]int{
		{5, 1},
		{10, 3},
		{20, 100},
		{2, 50},
	}
	assertPutGetDelete(t, vals)
}

// test skiplist inserting incrementally larger values
func TestSkiplistIncreasing(t *testing.T) {
	vals := [][]int{}
	for i := 0; i < 32; i++ {
		vals = append(vals, []int{i, i})
	}
	assertPutGetDelete(t, vals)
}

// test skiplist with decreasing values
func TestSkiplistDecreasing(t *testing.T) {
	vals := [][]int{}
	for i := 31; i >= 0; i-- {
		vals = append(vals, []int{i, i})
	}
	assertPutGetDelete(t, vals)
}

// test skiplist with random values
func TestSkiplistRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	vals := [][]int{}
	for i := 0; i < 32; i++ {
		vals = append(vals, []int{rng.Int(), rng.Int()})
	}
	assertPutGetDelete(t, vals)
}

// test skiplist heights work appropriately
func TestSkiplistHeight(t *testing.T) {

}

// test size fluctuations
func TestSkiplistSize(t *testing.T) {
	// test start 0
	// test normal inserts
	// test if overwriting value
	// test deleting values
	// test miss deleting valuess
}

// test overwriting an existing value
func TestPutDifferentPositions(t *testing.T) {
	vals := [][]int{
		{5, 5},  // from empty
		{10, 5}, // with one after
		{15, 5},
		{0, 5},  // before the head
		{2, 5},  // right after head
		{6, 5},  // in the middle of the list
		{12, 5}, // before the tail
		{20, 5}, // after the tail

		{0, -1},  // overwrite head
		{20, -1}, // overwrite tail
		{10, -1}, // overwrite middle of list
		{2, -1},  // overwrite element after head
		{15, -1}, // overwrite element after tail
	}
	list, err := skiplistFromArray(vals)
	if err != nil {
		t.Fatal(err)
	}
	// splice out elements that were overwritten
	// 1 (10), 2 (15), 3 (0), 4 (2), 7 (20)
	vals = append(append(vals[0:1], vals[4:7]...), vals[8:]...)
	assertListsEquivalent(t, vals, list)
}

// test finding values at random points
func TestSkiplistInsertPoints(t *testing.T) {
	// before list
	// after end of list
	// in middle
}

func TestSkiplistGetPoints(t *testing.T) {
	// key before first
	// first element
	// key exists in middle
	// key doesnt exist in middle
	// last element
	// key beyond last element
}

func TestSkiplistDelete(t *testing.T) {
	// deleting element that exists
	// deleting element which doesn't exist
}

// Things to test
// -- overwriting a value
// -- random 50 numbers, properly in order
// -- list with random values, can find
//    - elements at the head
//    - elements at the tail
//    - elements in the middle
// -- test height functionality works
//    - injecting randomness
// -- correctly returns nil for
//    - values in between values
//    - values before the list
//    - values at the end of the list
// -- support delete as well

// Failure modes
// -- Puting into list with no value

// NON GENERIC
// cpu: Intel(R) Core(TM) i5-8257U CPU @ 1.40GHz (home laptop)
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

// FIXED HEIGHT NEXTS
// goarch: arm64 (work laptop)
// BenchmarkSkiplistPerformance-11          1000000              1396 ns/op
// BenchmarkSkiplistReadHeavy-11            1000000              2482 ns/op
// BenchmarkSkiplistWriteHeavy-11           1000000              3265 ns/op
// BenchmarkSkiplistReadHits-11             1363813               917.7 ns/op
// BenchmarkSkiplistReadMisses-11           1378286               916.4 ns/op

// VARIABLE HEIGHT NEXTS
// BenchmarkSkiplistPerformance-11          1000000              1280 ns/op
// BenchmarkSkiplistReadHeavy-11            1000000              2409 ns/op
// BenchmarkSkiplistWriteHeavy-11           1000000              3804 ns/op
// BenchmarkSkiplistReadHits-11             1398126               997.9 ns/op
// BenchmarkSkiplistReadMisses-11           1410793               970.8 ns/op

// 15% speedup with multiple heads
// likely similar memory improvement
// MULTIPLE HEADS
// BenchmarkSkiplistPerformance-11          1000000              1192 ns/op
// BenchmarkSkiplistReadHeavy-11            1000000              2275 ns/op
// BenchmarkSkiplistWriteHeavy-11           1000000              3448 ns/op
// BenchmarkSkiplistReadHits-11             1442425               940.7 ns/op
// BenchmarkSkiplistReadMisses-11           1505054               918.5 ns/op
//
// SINGLE HEAD
// BenchmarkSkiplistPerformance-11          1000000              1400 ns/op
// BenchmarkSkiplistReadHeavy-11            1000000              2674 ns/op
// BenchmarkSkiplistWriteHeavy-11            956762              3879 ns/op
// BenchmarkSkiplistReadHits-11             1250410              1036 ns/op
// BenchmarkSkiplistReadMisses-11           1305763              1015 ns/op

func BenchmarkSkiplistPerformance(b *testing.B) {
	list := memtable.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		val := rng.Intn(1000000)
		list.Put(val, val)
		val = rng.Intn(1000000)
		list.Get(val)
	}
}

func BenchmarkSkiplistReadHeavy(b *testing.B) {
	list := memtable.NewSkiplist[int, int]()
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
	list := memtable.NewSkiplist[int, int]()
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
	list := memtable.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		val := rng.Intn(1000000)
		list.Put(val, val)
		list.Get(val)
	}
}

func BenchmarkSkiplistReadMisses(b *testing.B) {
	list := memtable.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		val := rng.Intn(1000000)
		list.Put(val, val)
		val = rng.Intn(1000000) + 1000000
		list.Get(val)
	}
}

func BenchmarkSkiplistWithDeletes(b *testing.B) {
	list := memtable.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	last := 0
	for i := 0; i < b.N; i++ {
		val := rng.Intn(1000000)
		list.Put(val, val)
		last = val
		val = rng.Intn(1000000) + 1000000
		list.Get(val)
		if val < 1000000/2 {
			list.Delete(last)
		}
	}
}
