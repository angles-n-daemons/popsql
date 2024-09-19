package memtable

import "fmt"

type MemStore struct {
	List *Skiplist[string, []byte]
}

func NewMemStore() *MemStore {
	return &MemStore{
		NewSkiplist[string, []byte](),
	}
}

func (m *MemStore) Put(key string, value []byte, overwrite bool) error {
	if !overwrite {
		val := m.List.Get(key)
		if val != nil {
			return fmt.Errorf("cannot overwrite value with key %s", key)
		}
	}

	_, err := m.List.Put(key, value)
	return err
}

func (m *MemStore) Get(start string) ([]byte, error) {
	// this abstraction may not live beyond testing the memory database
	node, _ := m.List.Search(start)
	if node != nil {
		return node.Val, nil
	}
	return nil, nil
}

func (m *MemStore) Scan(start, end string) ([][]byte, error) {
	// this abstraction may not live beyond testing the memory database
	_, prevs := m.List.Search(start)
	search := prevs[0]
	if search != nil && search.Key < start {
		search = search.Next()
	}

	results := [][]byte{}
	for search != nil && search.Key < end {
		results = append(results, search.Val)
		search = search.Next()
	}
	return results, nil
}
