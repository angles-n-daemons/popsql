package debug

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/kv"
)

func NewStore(store kv.Store) *DebugStore {
	return &DebugStore{store: store}
}

type DebugStore struct {
	store kv.Store
}

func (d *DebugStore) Get(key string) ([]byte, error) {
	b, err := d.store.Get(key)
	if err != nil {
		fmt.Println("GET ERROR", key, err)
		return nil, err
	}
	fmt.Println("GET", key, string(b))
	return b, nil
}

func (d *DebugStore) Put(key string, value []byte) error {
	err := d.store.Put(key, value)
	if err != nil {
		fmt.Println("PUT ERROR", key, err)
		return err
	}
	fmt.Println("PUT", key, string(value))
	return nil
}

func (d *DebugStore) GetRange(start, end string) (kv.Cursor, error) {
	c, err := d.store.GetRange(start, end)
	if err != nil {
		fmt.Println("GET RANGE ERROR", err)
	}
	return &DebugCursor{cursor: c}, nil
}

type DebugCursor struct {
	start  string
	end    string
	cursor kv.Cursor
}

func (d *DebugCursor) Read(num int) ([][]byte, error) {
	bb, err := d.cursor.Read(num)
	if err != nil {
		fmt.Printf("READ ERROR [%s, %s]: %s\n", d.start, d.end, err)
		return nil, err
	}
	for _, b := range bb {
		fmt.Printf("READ [%s, %s]: %s\n", d.start, d.end, string(b))
	}
	return nil, nil
}

func (d *DebugCursor) ReadAll() ([][]byte, error) {
	bb, err := d.cursor.ReadAll()
	if err != nil {
		fmt.Printf("READ ERROR [%s, %s]: %s\n", d.start, d.end, err)
		return nil, err
	}
	for _, b := range bb {
		fmt.Printf("READ [%s, %s]: %s\n", d.start, d.end, string(b))
	}
	return nil, nil
}

func (d *DebugCursor) IsAtEnd() bool {
	return d.cursor.IsAtEnd()
}
