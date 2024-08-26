package data

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

const MAX_UINT32 = ^uint32(0)
const MAX_HEIGHT = 32

type skiplistnode struct {
	val  int
	next []*skiplistnode
}

type Skiplist struct {
	Size   uint32
	height int8
	heads  []*skiplistnode
	// curmaxheight int
}

func NewSkiplist() *Skiplist {
	return &Skiplist{
		height: MAX_HEIGHT,
		heads:  make([]*skiplistnode, MAX_HEIGHT),
	}
}

// Important properties of a skip list
//   - The ith head is always >= the i-1th
//   - If the ith head is not nil, the i-1th head is also not nil
//   - While searching, if we set a node to update at i, we will
//     then set a node to update at i - 1
//   - If a head is nil at level i, the update value at i will also
//     be nil

/* Put takes a value and tries to insert it into the skiplist.
 * It can error if the skiplist is full or if the value already exists in the list.
 */
func (s *Skiplist) Put(val int) error {
	if s.Size >= MAX_UINT32 {
		return errors.New("cannot put element in skiplist, at maximum size.")
	}

	found, updates := s.search(val)
	if found {
		return errors.New("cannot put element in skiplist, element already exists")
	}

	// create the new node with a randomized height
	height := genHeight(MAX_HEIGHT)
	node := &skiplistnode{
		val:  val,
		next: make([]*skiplistnode, height),
	}

	// for each level in the nodes height, insert the node
	// into that level's list
	for i := 0; i < height; i++ {
		if s.heads[i] == nil {
			// if the head is nil at this level, the level is empty
			s.heads[i] = node
		} else if updates[i] == nil {
			// if the update value is nil, we never found a value
			// < val so we insert the node before the head
			node.next[i] = s.heads[i]
			s.heads[i] = node
		} else {
			// otherwise, we insert the node after update
			node.next[i] = updates[i].next[i]
			updates[i].next[i] = node
		}
	}

	s.Size++
	return nil
}

// Get finds the element in the skiplist if it exists, otherwise returns nil
func (s *Skiplist) Get(val int) *skiplistnode {
	found, nodes := s.search(val)
	if found {
		return nodes[0]
	} else {
		return nil
	}
}

// search is an internal function, leveraged by both Put and Get
// it searches through the list for a value, returning a search array
// of nodes preceeding or equal to the node value.
// if val exists in the list, it will be in the returned slice.
func (s *Skiplist) search(val int) (bool, []*skiplistnode) {
	// Find the highest head which is less than val
	level := s.height - 1
	var search *skiplistnode
	for level >= 0 {
		cand := s.heads[level]
		if cand != nil && cand.val < val {
			search = cand
			break
		}
		level--
	}

	// Keep a list of which directly preceed val or are equal to it
	nodes := make([]*skiplistnode, s.height)
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
		if next == nil || next.val > val {
			nodes[level] = search
			// reached the bottom of the list
			if level == 0 {
				break
			} else {
				level--
				continue
			}
		}
		search = next
	}

	found := nodes[0] != nil && nodes[0].val == val
	return found, nodes
}

func genHeight(maxHeight int) int {
	var height int = 1
	for rand.Intn(2) == 1 && height < maxHeight {
		height++
	}
	return height
}

// DebugPrint is a simple helper function for visualizing the skiplist
// Its output looks like the below example:
// [- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 95 --- --- --- --- --- --- --- --- --- --- --- 156 --- --- --- --- --- ---]
// [8 -- -- -- -- -- -- -- -- -- -- -- 56 -- -- -- -- -- -- -- -- -- -- -- 94 95 --- 106 --- --- --- --- --- --- 141 --- --- 156 --- --- --- --- --- ---]
// [8 -- -- -- -- -- -- -- -- 37 -- 47 56 -- -- -- -- -- -- -- -- -- -- -- 94 95 --- 106 --- 118 124 --- --- --- 141 --- --- 156 --- --- --- --- --- ---]
// [8 -- 13 -- -- -- -- 31 33 37 45 47 56 -- -- -- -- -- 81 85 -- -- -- -- 94 95 --- 106 --- 118 124 --- --- 140 141 --- --- 156 --- --- --- --- 190 ---]
// [8 11 13 15 25 26 29 31 33 37 45 47 56 58 59 66 74 78 81 85 87 88 89 90 94 95 100 106 111 118 124 128 137 140 141 147 153 156 159 162 163 187 190 194]
func (s *Skiplist) DebugPrint(levels int) {
	lists := make([][]string, 5)
	for i := 0; i < levels; i++ {
		lists[i] = []string{}
	}

	node := s.heads[0]
	for node != nil {
		height := len(node.next)
		str := strconv.Itoa(node.val)
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
