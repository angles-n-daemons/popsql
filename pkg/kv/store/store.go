package store

import (
	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/kv/store/debug"
	"github.com/angles-n-daemons/popsql/pkg/kv/store/memtable"
)

func NewDebugStore(store kv.Store) *debug.DebugStore {
	return debug.NewStore(store)
}

func NewMemStore() *memtable.Memstore {
	return memtable.NewStore()
}
