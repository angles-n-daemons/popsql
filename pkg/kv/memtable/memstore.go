package memtable

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/kv/store"
)

// Memstore is a struct which satisfies the Store interface
// and works entirely in memory. It's useful for testing the behavior of the system.
func NewMemstore() *Memstore {
	return &Memstore{
		List: NewSkiplist[string, []byte](),
	}
}

type Memstore struct {
	List *Skiplist[string, []byte]
}

func (m *Memstore) Get(key string) ([]byte, error) {
	node := m.List.Get(key)
	if node == nil {
		return nil, nil
	}
	return node.Val, nil
}

func (m *Memstore) GetRange(start, end string) (store.Cursor, error) {
	node, prevs := m.List.Search(start)
	if node == nil && prevs[0] != nil {
		node = prevs[0].Next()
	}
	if node == nil {
		node = m.List.Head()
	}
	return &Memcursor{
		Node: node,
		End:  end,
	}, nil
}

func (m *Memstore) Put(key string, value []byte) error {
	_, err := m.List.Put(key, value)
	return err
}

func (m *Memstore) PutRange(start, end string, value []byte) error {
	return fmt.Errorf("not implemented")
}
