package memtable_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/kv/store/memtable"
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

func skiplistFromArray(elems [][]int) (*memtable.Skiplist[int, int], error) {
	list := memtable.NewSkiplist[int, int]()
	for _, elem := range elems {
		_, err := list.Put(elem[0], elem[1])
		if err != nil {
			return nil, fmt.Errorf(
				"list raised an error on Put: %v",
				err,
			)
		}
	}
	return list, nil
}

func assertListsEquivalent(t *testing.T, elems [][]int, list *memtable.Skiplist[int, int]) {
	for _, elem := range elems {
		node := list.Get(elem[0])
		if node == nil {
			t.Fatalf("expected Get to find key %d in list", elem[0])
		}
		if node.Val != elem[1] {
			t.Fatalf(
				"expected val %d for key %d on Get, but got %d",
				elem[1], elem[0], node.Val,
			)
		}
	}
}

func assertDeleteFromList(t *testing.T, elems [][]int, list *memtable.Skiplist[int, int]) {
	// check elements deleted
	for _, elem := range elems {
		node := list.Delete(elem[0])
		if node == nil {
			t.Fatalf("expected Delete to return node when deleting")
		}
	}
	for _, elem := range elems {
		node := list.Get(elem[0])
		if node != nil {
			t.Fatalf("expected Get to return nil after deletion")
		}
		node = list.Delete(elem[0])
		if node != nil {
			t.Fatalf("expected Delete to return nil on second attempt")
		}
	}
}

// helper function which asserts that the size of the elems array matches the
// skiplist, and that the elements found in the elems array can be found
// in the skiplist, and then deletes the elements when finished
func assertPutGetDelete(t *testing.T, elems [][]int) {
	list, err := skiplistFromArray(elems)
	if err != nil {
		t.Fatal(err)
	}
	if len(elems) != int(list.Size) {
		t.Fatalf(
			"expected list length to match elems %d but got length %d",
			len(elems), list.Size,
		)
	}
	assertListsEquivalent(t, elems, list)
	assertDeleteFromList(t, elems, list)
	if list.Size != 0 {
		t.Fatalf(
			"expected list to be empty after deleting elements, but got size %d",
			list.Size,
		)
	}
}

// test very simple skiplist use cases
func TestSkiplistBasic(t *testing.T) {
	elems := [][]int{
		{5, 1},
		{10, 3},
		{20, 100},
		{2, 50},
	}
	assertPutGetDelete(t, elems)
}

// test skiplist inserting incrementally larger values
func TestSkiplistIncreasing(t *testing.T) {
	elems := [][]int{}
	for i := 0; i < 32; i++ {
		elems = append(elems, []int{i, i})
	}
	assertPutGetDelete(t, elems)
}

// test skiplist with decreasing values
func TestSkiplistDecreasing(t *testing.T) {
	elems := [][]int{}
	for i := 31; i >= 0; i-- {
		elems = append(elems, []int{i, i})
	}
	assertPutGetDelete(t, elems)
}

// test skiplist with random values
func TestSkiplistRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	elems := [][]int{}
	for i := 0; i < 32; i++ {
		elems = append(elems, []int{rng.Int(), rng.Int()})
	}
	assertPutGetDelete(t, elems)
}

// test skiplist heights work appropriately
func TestSkiplistHeight(t *testing.T) {
	elems := [][]int{
		// { key, val, height }
		{8, 2, 3},
		{3, 1, 1},
		{6, 3, 3},
		{2, 10, 3},
		{4, 8, 2},
		{7, 1, 4},
	}
	heightSource := &mockHeightRandSource{}
	list := memtable.NewSkiplistWithRandSource[int, int](heightSource)
	for _, elem := range elems {
		heightSource.changeHeight(elem[2])
		list.Put(elem[0], elem[1])
	}
	expected := [][]int{
		{2, 3, 4, 6, 7, 8},
		{2, 4, 6, 7, 8},
		{2, 6, 7, 8},
		{7},
	}
	for i, expectedRow := range expected {
		row, err := list.DebugGetRow(i)
		if err != nil {
			t.Fatal(err)
		}
		if len(expectedRow) != len(row) {
			t.Fatalf(
				"checking row %d, expected length %d, got length %d",
				i,
				len(expectedRow),
				len(row),
			)
		}
		for j, key := range expectedRow {
			if key != row[j] {
				t.Fatalf(
					"expected key %d at row %d depth %d, got %d",
					key,
					i,
					j,
					row[j],
				)
			}
		}
	}
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
	elems := [][]int{
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
	list, err := skiplistFromArray(elems)
	if err != nil {
		t.Fatal(err)
	}
	// splice out elements that were overwritten
	// 1 (10), 2 (15), 3 (0), 4 (2), 7 (20)
	expected := append(append(elems[0:1], elems[5:7]...), elems[8:]...)
	assertListsEquivalent(t, expected, list)
	// expect deduplication on count
	if list.Size != 8 {
		t.Fatalf(
			"expected list size to be %d after Puts, got %d",
			8,
			list.Size,
		)
	}
}

func TestSkiplistGetPoints(t *testing.T) {
	elems := [][]int{
		{5, 10},
		{10, 20},
		{15, 5},
		{0, 15},
		{20, 0},
	}
	list, err := skiplistFromArray(elems)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		key    int
		val    int
		inList bool
	}{
		// key, val, inList
		{-5, -1, false}, // before head
		{0, 15, true},   // head
		{2, -1, false},  // in between head and next node
		{5, 10, true},   // value after head
		{8, -1, false},  // value missing in the middle
		{10, 20, true},  // value in middle of list
		{15, 5, true},   // last value before tail
		{20, 0, true},   // tail
		{25, -1, false}, // after tail
	} {
		node := list.Get(test.key)
		if node == nil {
			if !test.inList {
				continue
			} else {
				t.Fatalf("expected Get key %d to return value, but returned nil", test.key)
			}
		}
		if !test.inList {
			t.Fatalf("expected Get key %d to return nil, but returned key %d", test.key, node.Key)
		}
		if node.Key != test.key {
			t.Fatalf("Get return incorrect key %d for key %d", node.Key, test.key)
		}
		if node.Val != test.val {
			t.Fatalf("Get return incorrect val %d, expected %d", node.Val, test.val)
		}
	}
	if list.Size != 5 {
		t.Fatalf(
			"expected list size to be %d after Gets, got %d",
			5,
			list.Size,
		)
	}
}

func TestSkiplistDelete(t *testing.T) {
	elems := [][]int{
		{0, 0},
		{5, 5},
		{10, 10},
		{15, 15},
		{20, 20},
		{25, 25},
		{30, 30},
		{35, 35},
	}
	list, err := skiplistFromArray(elems)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		key    int
		val    int
		inList bool
	}{
		// key, val, inList
		{-5, -1, false}, // before head
		{0, 0, true},    // head
		{2, -1, false},  // in between head and next node
		{5, 5, true},    // value after head
		{8, -1, false},  // value missing in the middle
		{10, 10, true},  // value in middle of list
		{30, 30, true},  // last value before tail
		{35, 35, true},  // tail
		{40, -1, false}, // after tail
	} {
		node := list.Delete(test.key)
		if node == nil {
			if !test.inList {
				continue
			} else {
				t.Fatalf("expected Get key %d to return value, but returned nil", test.key)
			}
		}
		if !test.inList {
			t.Fatalf("expected Get key %d to return nil, but returned key %d", test.key, node.Key)
		}
		if node.Key != test.key {
			t.Fatalf("Get return incorrect key %d for key %d", node.Key, test.key)
		}
		if node.Val != test.val {
			fmt.Println(node.Key, test.key)
			t.Fatalf("Get return incorrect val %d, expected %d", node.Val, test.val)
		}
	}
	if list.Size != 3 {
		t.Fatalf(
			"expected list size to be %d after Deletes, got %d",
			3,
			list.Size,
		)
	}
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
		elem := rng.Intn(1000000)
		list.Put(elem, elem)
		elem = rng.Intn(1000000)
		list.Get(elem)
	}
}

func BenchmarkSkiplistReadHeavy(b *testing.B) {
	list := memtable.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		elem := rng.Intn(1000000)
		list.Put(elem, elem)
		for j := 0; j < 3; j++ {
			elem = rng.Intn(1000000)
			list.Get(elem)
		}
	}
}

func BenchmarkSkiplistWriteHeavy(b *testing.B) {
	list := memtable.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		var elem int
		for j := 0; j < 3; j++ {
			elem = rng.Intn(1000000)
			list.Put(elem, elem)
		}
		elem = rng.Intn(1000000)
		list.Get(elem)
	}
}

func BenchmarkSkiplistReadHits(b *testing.B) {
	list := memtable.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		elem := rng.Intn(1000000)
		list.Put(elem, elem)
		list.Get(elem)
	}
}

func BenchmarkSkiplistReadMisses(b *testing.B) {
	list := memtable.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		elem := rng.Intn(1000000)
		list.Put(elem, elem)
		elem = rng.Intn(1000000) + 1000000
		list.Get(elem)
	}
}

func BenchmarkSkiplistWithDeletes(b *testing.B) {
	list := memtable.NewSkiplist[int, int]()
	rng := rand.New(rand.NewSource(1))
	last := 0
	for i := 0; i < b.N; i++ {
		elem := rng.Intn(1000000)
		list.Put(elem, elem)
		last = elem
		elem = rng.Intn(1000000) + 1000000
		list.Get(elem)
		if elem < 1000000/2 {
			list.Delete(last)
		}
	}
}
