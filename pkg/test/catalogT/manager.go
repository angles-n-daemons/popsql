package catalogT

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/db/kv/store"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
)

func Manager(t *testing.T) *catalog.Manager {
	st := store.NewMemStore()
	m, err := catalog.NewManager(st)
	assert.NoError(t, err)
	return m
}
