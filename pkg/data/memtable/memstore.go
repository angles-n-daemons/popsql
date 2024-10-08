package memtable

import "fmt"

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

func (m *Memstore) GetRange(start, end string) (*Memcursor, error) {
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

func (m *Memstore) Set(key string, value []byte) error {
	_, err := m.List.Put(key, value)
	return err
}

func (m *Memstore) SetRange(start, end string, value []byte) error {
	return fmt.Errorf("not implemented")
}
