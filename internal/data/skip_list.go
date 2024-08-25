package data

import (
	"errors"
	"math/rand"
)

const MAX_UINT32 = ^uint32(0)
const MAX_HEIGHT = 32

type skiplistnode struct {
	val  int
	next []*skiplistnode
}

type SkipList struct {
	Size   uint32
	height int
	heads  []*skiplistnode
	// curmaxheight int
}

func NewSkipList() *SkipList {
	return &SkipList{
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

/*
Insert is broken into three sections:
 1. Start the search with the highest head node which is less than val.
 2. For each level, add the greatest node which less than val to an update list.
 3. Determine a height for the node and insert it after each update relevant.
*/
func (s *SkipList) Insert(val int) error {
	if s.Size >= MAX_UINT32 {
		return errors.New("cannot insert element to skiplist, at maximum size.")
	}

	// Start the search at the highest level where the head is
	// less than val
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

	// create an updates array for tracking the last seen node at each level
	// at each level starting here, we will add nodes less than val
	// if a node is the largest node less than val at multiple levels
	// that same node will be added multiple times to the updates array
	updates := make([]*skiplistnode, s.height)
	for search != nil {
		next := search.next[level]
		// if the next value is greater at this level, or it is nil
		// we can continue the search one level down
		if next == nil || next.val > val {
			updates[level] = search
			// reached the bottom of the list
			if level == 0 {
				break
			} else {
				level--
				continue
			}
		}
		if next.val == val {
			return errors.New("cannot insert into list, value already exists")
		}
		search = next
	}

	// insert the node in the various levels
	height := genHeight(MAX_HEIGHT)
	node := &skiplistnode{
		val:  val,
		next: make([]*skiplistnode, height),
	}

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

func (s *SkipList) Search(val int) *skiplistnode {
	level := s.height - 1
	for level >= 0 && s.heads[level] == nil {
		level--
	}
	if level == -1 {
		return nil
	}

	var search *skiplistnode = s.heads[level]
	for search != nil {
		next := search.next[level]
		if next == nil || next.val > val {
			if level == 0 {
				return nil
			} else {
				level--
				continue
			}
		}
		if next.val == val {
			return next
		}
		search = next
	}

	return nil
}

func (s *SkipList) GetLevelSlice(level int) []int {
	node := s.heads[level]
	result := []int{}
	for node != nil {
		result = append(result, node.val)
		node = node.next[level]
	}
	return result
}

func genHeight(maxHeight int) int {
	var height int = 1
	for rand.Intn(2) == 1 && height < maxHeight {
		height++
	}
	return height
}

type linkedlistnode struct {
	val  int
	next *linkedlistnode
}

type LinkedList struct {
	head *linkedlistnode
}

func (l *LinkedList) Insert(val int) error {
	node := &linkedlistnode{val: val, next: nil}
	if l.head == nil {
		l.head = node
		return nil
	}

	search := l.head
	next := search.next
	for next != nil && next.val < val {
		search = search.next
		next = search.next
	}
	if search.val == val {
		return errors.New("cannot insert into list, value already exists")
	}
	node.next = search.next
	search.next = node
	return nil
}

func (l *LinkedList) Search(val int) *linkedlistnode {
	search := l.head
	for search != nil {
		if search.val == val {
			return search
		}
		search = search.next
	}
	return nil
}
