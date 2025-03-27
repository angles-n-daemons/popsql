package catalog_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/kv/memtable"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
	"github.com/angles-n-daemons/popsql/pkg/test/catalogT"
)

func TestBootstrap(t *testing.T) {
	// Create a new Memstore and Manager
	st := memtable.NewMemstore()
	// create manager by hand.
	m := &catalog.Manager{Store: st, Schema: schema.New()}

	m.Bootstrap()

	checkManagerTable(t, m, catalog.InitMetaTable())
	checkManagerSequence(t, m, catalog.InitMetaTableSequence())
	checkManagerTable(t, m, catalog.InitSequencesTable())
	checkManagerSequence(t, m, catalog.InitSequencesTableSequence())
}

func TestBootstrapIdempotence(t *testing.T) {
	// Create a new Memstore and Manager
	m := catalogT.Manager(t)
	m.Schema = schema.New()

	// Bootstrap twice - should not error
	err := m.Bootstrap()
	assert.NoError(t, err)

	err = m.Bootstrap()
	assert.IsError(t, err, "table '__tables__' already exists")
}
