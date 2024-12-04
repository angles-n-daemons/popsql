package catalog

import (
	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
)

// System Table Keys
var CATALOG_KEYS_PREFIX = keys.New("__tables")
var CATALOG_KEYS_END = CATALOG_KEYS_PREFIX.Next()

// catalog.Manager is responsible for persisting and loading the database
// schema to and from disk.
type Manager struct {
	Schema *Schema
	Store  *kv.Store
}

func (m *Manager) NewManager(store *kv.Store) *Manager {
	// check to see if the system table is in the cache

	return &Manager{
		Store: store,
	}
}
