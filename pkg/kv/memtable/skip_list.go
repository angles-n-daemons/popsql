package memtable

import (
	"cmp"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const MAX_UINT32 = ^uint32(0)
const MAX_HEIGHT = 32

type SkiplistNode[K cmp.Ordered, V any] struct {
	Key  K
	Val  V
	next []*SkiplistNode[K, V]
}

func (node *SkiplistNode[K, V]) Next() *SkiplistNode[K, V] {
	return node.next[0]
}

/*
 * A Skiplist is an efficiently sorted data structure.
 * It's desirable because its performance is similar to that
 * of a balanced tree, but it's simple to reason about and impl.
 */

// Important properties of a skip list:
//   - The ith head is always >= the i-1th
//   - If the ith head is not nil, the i-1th head is also not nil
//   - While searching, if we set a node to update at i, we will
//     then set a node to update at i - 1
//   - If a head is nil at level i, the update value at i will also
//     be nil
type Skiplist[K cmp.Ordered, V any] struct {
	Size   uint32
	height int8
	heads  []*SkiplistNode[K, V]
	rng    *rand.Rand
}

func NewSkiplist[K cmp.Ordered, V any]() *Skiplist[K, V] {
	return &Skiplist[K, V]{
		Size:   0,
		height: MAX_HEIGHT,
		heads:  make([]*SkiplistNode[K, V], MAX_HEIGHT),
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func NewSkiplistWithRandSource[K cmp.Ordered, V any](source rand.Source) *Skiplist[K, V] {
	return &Skiplist[K, V]{
		Size:   0,
		height: MAX_HEIGHT,
		heads:  make([]*SkiplistNode[K, V], MAX_HEIGHT),
		rng:    rand.New(source),
	}
}

func (list *Skiplist[K, V]) Head() *SkiplistNode[K, V] {
	return list.heads[0]
}

/* Put takes a value and tries to insert it into the skiplist.
 * This function returns a boolean which is true if the element is new.
 * It can error if the skiplist is full.
 */
func (list *Skiplist[K, V]) Put(key K, val V) (bool, error) {
	if list.Size >= MAX_UINT32 {
		return false, errors.New("cannot put element in skiplist, at maximum size.")
	}

	node, prevs := list.Search(key)
	if node != nil {
		// if the node already exists, we change its value
		node.Val = val
		return false, nil
	}

	// create the new node with a randomized height
	height := list.genHeight(MAX_HEIGHT)
	node = &SkiplistNode[K, V]{
		Key:  key,
		Val:  val,
		next: make([]*SkiplistNode[K, V], height),
	}

	// for each level in the nodes height, insert the node
	// into that level's list
	for i := 0; i < height; i++ {
		if list.heads[i] == nil {
			// if the head is nil at this level, the level is empty
			list.heads[i] = node
		} else if prevs[i] == nil {
			// if the update value is nil, we never found a value
			// < val so we insert the node before the head
			node.next[i] = list.heads[i]
			list.heads[i] = node
		} else {
			// otherwise, we insert the node after update
			node.next[i] = prevs[i].next[i]
			prevs[i].next[i] = node
		}
	}

	list.Size++
	return true, nil
}

// Get finds the element in the skiplist if it exists, otherwise returns nil
func (list *Skiplist[K, V]) Get(key K) *SkiplistNode[K, V] {
	node, _ := list.Search(key)
	return node
}

// Delete removes the element with the specified key from the list if it exists
func (list *Skiplist[K, V]) Delete(key K) *SkiplistNode[K, V] {
	node, prevs := list.Search(key)
	// If we didn't find the node, return nil
	if node == nil {
		return nil
	}

	// set the next pointer for the previous nodes
	// to the node's next pointer (where applicable)
	for i := 0; i < len(node.next); i++ {
		if prevs[i] == nil {
			list.heads[i] = node.next[i]
		} else {
			prevs[i].next[i] = node.next[i]
		}
	}

	list.Size--
	return node
}

// Search is an internal function, leveraged by Put, Get and Delete
// it searches through the list for a value, returning a Search array
// of nodes preceeding or equal to the node value.
// if the key exists, it will be returned in addition to the Search array
func (list *Skiplist[K, V]) Search(key K) (*SkiplistNode[K, V], []*SkiplistNode[K, V]) {
	// Find the highest head which is less than val
	level := list.height - 1
	var search *SkiplistNode[K, V]
	// Keep a list of which directly preceed val or are equal to it
	prevs := make([]*SkiplistNode[K, V], list.height)

	// Special case for when the head is the node we're looking for
	if list.heads[0] != nil && list.heads[0].Key == key {
		return list.heads[0], prevs
	}

	// Start the search at the first head whose key is less than
	// the one we are looking for
	for level >= 0 {
		cand := list.heads[level]
		if cand != nil && cand.Key < key {
			search = cand
			break
		}
		level--
	}

	// Run the search at each subsequent level below
	// For each level, continue traversing the list until either:
	//   * the next node is greater than val
	//   * the next node is nil
	// On these conditions, drop to the next level down and continue
	// If the level is 0, exit the loop
	for search != nil {
		next := search.next[level]
		// if the next value is greater at this level, or it is nil
		// we can continue the search one level down
		if next == nil || next.Key >= key {
			prevs[level] = search
			// reached the bottom of the list
			if level == 0 {
				// If we found the right node, return it
				// We do this at the bottom level so that
				// all the previous values can be properly set
				if next != nil && next.Key == key {
					return next, prevs
				}
				break
			} else {
				level--
				continue
			}
		}
		search = next
	}

	return nil, prevs
}

func (list *Skiplist[K, V]) genHeight(maxHeight int) int {
	var height int = 1
	for list.rng.Intn(2) == 1 && height < maxHeight {
		height++
	}
	return height
}

func (list *Skiplist[K, V]) DebugGetRow(level int) ([]K, error) {
	if level > len(list.heads) {
		return nil, fmt.Errorf("cannot get level %d of skiplist with height %d", level, len(list.heads))
	}
	keys := []K{}
	node := list.heads[level]
	for node != nil {
		keys = append(keys, node.Key)
		node = node.next[level]
	}
	return keys, nil
}

// DebugPrint is a simple helper function for visualizing the skiplist
// Its output looks like the below example:
// [- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 95 --- --- --- --- --- --- --- --- --- --- --- 156 --- --- --- --- --- ---]
// [8 -- -- -- -- -- -- -- -- -- -- -- 56 -- -- -- -- -- -- -- -- -- -- -- 94 95 --- 106 --- --- --- --- --- --- 141 --- --- 156 --- --- --- --- --- ---]
// [8 -- -- -- -- -- -- -- -- 37 -- 47 56 -- -- -- -- -- -- -- -- -- -- -- 94 95 --- 106 --- 118 124 --- --- --- 141 --- --- 156 --- --- --- --- --- ---]
// [8 -- 13 -- -- -- -- 31 33 37 45 47 56 -- -- -- -- -- 81 85 -- -- -- -- 94 95 --- 106 --- 118 124 --- --- 140 141 --- --- 156 --- --- --- --- 190 ---]
// [8 11 13 15 25 26 29 31 33 37 45 47 56 58 59 66 74 78 81 85 87 88 89 90 94 95 100 106 111 118 124 128 137 140 141 147 153 156 159 162 163 187 190 194]
func DebugPrintList[K cmp.Ordered, V any](list *Skiplist[K, V], levels int) {
	lists := make([][]string, 5)
	for i := 0; i < levels; i++ {
		lists[i] = []string{}
	}

	node := list.heads[0]
	for node != nil {
		height := len(node.next)
		str := fmt.Sprintf("%v", node.Key)
		for i := 0; i < levels; i++ {
			if i < height {
				lists[i] = append(lists[i], str)
			} else {
				lists[i] = append(lists[i], strings.Repeat("-", len(str)))
			}
		}
		node = node.next[0]
	}

	for i := levels - 1; i >= 0; i-- {
		fmt.Println(lists[i])
	}
}

func DebugPrintLevels[K cmp.Ordered, V any](s *Skiplist[K, V], levels int) {
	for i := levels - 1; i >= 0; i-- {
		fmt.Println(s.DebugGetRow(i))
	}
}
