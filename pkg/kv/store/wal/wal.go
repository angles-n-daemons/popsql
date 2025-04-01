package wal

import "io"

func NewStore() *WALStore {
	return &WALStore{}
}

type WALStore struct {
	f *io.ReadWriter
}

func (w *WALStore) Put(string, []byte) {

}
