package catalog

import "github.com/angles-n-daemons/popsql/pkg/kv"

// System Table Keys
var CATALOG_KEYS_PREFIX = kv.NewKey("__tables")
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
