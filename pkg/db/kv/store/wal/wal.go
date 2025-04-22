package wal

import (
	"io"

	"github.com/angles-n-daemons/popsql/pkg/db/kv"
	"github.com/angles-n-daemons/popsql/pkg/db/kv/store/memtable"
)

func NewStore() *WALStore {
	return &WALStore{}
}

func NewStoreFromBackup() {

}

type WALStore struct {
	f *io.ReadWriter
	m *memtable.Memstore
}

func (w *WALStore) Put(k string, v []byte) error {
	return w.m.Put(k, v)
}

func (w *WALStore) Get(k string) ([]byte, error) {
	return w.m.Get(k)
}

func (w *WALStore) Scan(start, end string) (kv.Cursor, error) {
	return w.m.Scan(start, end)

}
