package memtable

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/kv"
)

// Memstore is a struct which satisfies the Store interface
// and works entirely in memory. It's useful for testing the behavior of the system.
func NewMemstore() *Memstore {
	return &Memstore{
		List: NewSkiplist[string, []byte](),
	}
}

// Memstore is an in-memory key-value store designed to satisfy the Store interface.
type Memstore struct {
	List *Skiplist[string, []byte]
}

// Get retrieves the value associated with the given key.
// If the key does not exist in the Memstore, it returns nil and no error.
//
// Parameters:
//
//	key - The key whose associated value is to be returned.
//
// Returns:
//
//	[]byte - The value associated with the specified key, or nil if the key does not exist.
//	error - An error if there is an issue retrieving the value, otherwise nil.
func (m *Memstore) Get(key string) ([]byte, error) {
	node := m.List.Get(key)
	if node == nil {
		return nil, nil
	}
	return node.Val, nil
}

// GetRange retrieves a range of elements from the Memstore starting from the
// specified 'start' key up to, but not including, the 'end' key.
//
// Parameters:
// - start: The starting key of the range to retrieve.
// - end: The ending key of the range to retrieve.
//
// Returns:
// - kv.Cursor: A cursor pointing to the start of the range within the Memstore.
// - error: An error if the range cannot be retrieved.
//
// The function searches for the node corresponding to the 'start' key. If the node
// is not found, it attempts to use the previous node's next pointer. If still not
// found, it defaults to the head of the list. The returned cursor will iterate
// from the found node up to the 'end' key.
func (m *Memstore) GetRange(start, end string) (kv.Cursor, error) {
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

// Put stores the given key-value pair in the Memstore.
// If the key already exists, its value is updated with the new value.
//
// Parameters:
//
//	key - The key to be stored.
//	value - The value to be associated with the specified key.
//
// Returns:
//
//	error - An error if there is an issue storing the key-value pair, otherwise nil.
func (m *Memstore) Put(key string, value []byte) error {
	fmt.Println("putting", key, string(value))
	_, err := m.List.Put(key, value)
	return err
}
