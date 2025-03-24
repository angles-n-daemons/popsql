package catalog_test

import (
	"encoding/json"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/db/catalog"
	"github.com/angles-n-daemons/popsql/pkg/kv/memtable"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema/desc"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
)

func verifyManagerTable(t *testing.T, cat *catalog.Catalog, table *desc.Table, name string) {
	// Verify that the table was added to the desc.
	schemaTable, ok := cat.Schema.GetTable(name)
	if !ok {
		t.Fatalf("table '%s' was not in the desc", name)
	}
	assert.Equal(t, schemaTable, table)

	// Check that the table was written to the store.
	tableBytes, err := cat.Store.Get(catalog.MetaTableKey(table))
	assert.NoError(t, err)

	storeTable := &desc.Table{}
	err = json.Unmarshal(tableBytes, storeTable)
	assert.NoError(t, err)
	assert.Equal(t, storeTable, table)
}

func verifyManagerSequence(
	t *testing.T, cat *catalog.Catalog, sequence *desc.Sequence, name string,
) {
	// Verify that the sequence was added to the desc.
	schemaSequence, ok := cat.Schema.GetSequence(name)
	if !ok {
		t.Fatalf("sequence '%s' was not in the desc", name)
	}
	assert.Equal(t, schemaSequence, sequence)

	// Check that the sequence was written to the store.
	sequenceBytes, err := cat.Store.Get(catalog.SequenceKey(sequence))
	assert.NoError(t, err)

	storeSequence := &desc.Sequence{}
	err = json.Unmarshal(sequenceBytes, storeSequence)
	assert.NoError(t, err)
	assert.Equal(t, storeSequence, sequence)
}

func TestBootstrap(t *testing.T) {
	// Create a new Memstore and Manager
	store := memtable.NewMemstore()
	cat, err := catalog.NewCatalog(store)
	assert.NoError(t, err)

	verifyManagerTable(t, cat, catalog.InitMetaTable, catalog.MetaTableName)
	verifyManagerTable(t, cat, catalog.InitSequencesTable, catalog.SequencesTableName)

	verifyManagerSequence(t, cat, catalog.InitMetaTableSequence, catalog.MetaTableSequenceName)
	verifyManagerSequence(t, cat, catalog.InitSequencesTableSequence, catalog.SequencesTableSequenceName)

	// Test that we can load the desc from the store without a Bootstrap call.
	// Create a new cat to verify persistence
	newManager, err := catalog.NewCatalog(store)
	assert.NoError(t, err)

	verifyManagerTable(t, newManager, catalog.InitMetaTable, catalog.MetaTableName)
	verifyManagerTable(t, newManager, catalog.InitSequencesTable, catalog.SequencesTableName)

	// need to also load sequences
	verifyManagerSequence(t, newManager, catalog.InitMetaTableSequence, catalog.MetaTableSequenceName)
	verifyManagerSequence(t, newManager, catalog.InitSequencesTableSequence, catalog.SequencesTableSequenceName)
}

func TestBootstrapIdempotence(t *testing.T) {
	// Create a new Memstore and Manager
	store := memtable.NewMemstore()
	cat, err := catalog.NewCatalog(store)
	assert.NoError(t, err)
	cat.Schema = schema.NewSchema()

	// Bootstrap twice - should not error
	err = cat.Bootstrap()
	assert.NoError(t, err)

	err = cat.Bootstrap()
	assert.IsError(t, err, "table '__tables___sequence' already exists")
}
