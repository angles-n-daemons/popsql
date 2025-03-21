package catalog_test

import (
	"encoding/json"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/db/catalog"
	"github.com/angles-n-daemons/popsql/pkg/kv/memtable"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
	"github.com/angles-n-daemons/popsql/pkg/testing/assert"
)

func verifyManagerTable(t *testing.T, manager *catalog.Manager, table *schema.Table, name string) {
	// Verify that the table was added to the schema.
	schemaTable, ok := manager.Schema.GetTable(name)
	if !ok {
		t.Fatalf("table '%s' was not in the schema", name)
	}
	assert.Equal(t, schemaTable, table)

	// Check that the table was written to the store.
	tableBytes, err := manager.Store.Get(catalog.MetaTableKey(table))
	assert.NoError(t, err)

	storeTable := &schema.Table{}
	err = json.Unmarshal(tableBytes, storeTable)
	assert.NoError(t, err)
	assert.Equal(t, storeTable, table)
}

func verifyManagerSequence(
	t *testing.T, manager *catalog.Manager, sequence *schema.Sequence, name string,
) {
	// Verify that the sequence was added to the schema.
	schemaSequence, ok := manager.Schema.GetSequence(name)
	if !ok {
		t.Fatalf("sequence '%s' was not in the schema", name)
	}
	assert.Equal(t, schemaSequence, sequence)

	// Check that the sequence was written to the store.
	sequenceBytes, err := manager.Store.Get(catalog.SequenceKey(sequence))
	assert.NoError(t, err)

	storeSequence := &schema.Sequence{}
	err = json.Unmarshal(sequenceBytes, storeSequence)
	assert.NoError(t, err)
	assert.Equal(t, storeSequence, sequence)
}

func TestBootstrap(t *testing.T) {
	// Create a new Memstore and Manager
	store := memtable.NewMemstore()
	manager := catalog.NewManager(store)

	// Initialize the schema
	manager.Schema = schema.NewSchema()

	// Call Bootstrap
	err := manager.Bootstrap()
	if err != nil {
		t.Fatalf("Failed to bootstrap: %v", err)
	}

	verifyManagerTable(t, manager, catalog.InitMetaTable, catalog.MetaTableName)
	verifyManagerTable(t, manager, catalog.InitSequencesTable, catalog.SequencesTableName)

	verifyManagerSequence(t, manager, catalog.InitMetaTableSequence, catalog.MetaTableSequenceName)
	verifyManagerSequence(t, manager, catalog.InitSequencesTableSequence, catalog.SequencesTableSequenceName)

	// Test that we can load the schema from the store without a Bootstrap call.
	// Create a new manager to verify persistence
	newManager := catalog.NewManager(store)
	newManager.Schema = schema.NewSchema()
	err = newManager.LoadSchema()
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
	manager := catalog.NewManager(store)
	manager.Schema = schema.NewSchema()

	// Bootstrap twice - should not error
	err := manager.Bootstrap()
	assert.NoError(t, err)

	err = manager.Bootstrap()
	assert.IsError(t, err, "table '__tables___sequence' already exists")
}
