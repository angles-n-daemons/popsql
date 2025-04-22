package store

import (
	"github.com/angles-n-daemons/popsql/pkg/db/kv"
	"github.com/angles-n-daemons/popsql/pkg/db/kv/store/debug"
	"github.com/angles-n-daemons/popsql/pkg/db/kv/store/memtable"
)

func NewDebugStore(store kv.Store) *debug.DebugStore {
	return debug.NewStore(store)
}

func NewMemStore() *memtable.Memstore {
	return memtable.NewStore()
}
