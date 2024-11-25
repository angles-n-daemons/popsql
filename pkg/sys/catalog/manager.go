package catalog

import "github.com/angles-n-daemons/popsql/pkg/kv/store"

type Manager struct {
	Schema Schema
	Store  store.Store
}

func NewManager(store store.Store)
