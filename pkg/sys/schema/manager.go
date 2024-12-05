package schema

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
)

// System Table Keys
var CATALOG_KEYS_PREFIX = keys.New("__tables")
var CATALOG_KEYS_END = CATALOG_KEYS_PREFIX.Next()

// Meta Table Name
const MetaTableName = "__schema__"

// Manager is responsible for holding the entire schema as well as keeping it
// in sync with the underlying data store.
type Manager struct {
	Schema *Schema
	Store  kv.Store
}

func (m *Manager) NewManager(store kv.Store) (*Manager, error) {
	return &Manager{
		Store: store,
	}, nil
}

func (m *Manager) LoadSchema() error {
	cur, err := m.Store.GetRange(MetaTable.Prefix().String(), MetaTable.PrefixEnd().String())
	if err != nil {
		return fmt.Errorf("failed to read the table catalog from the store %w", err)
	}

	tablesBytes, err := cur.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read the table catalog from a cursor %w", err)
	}

	schema, err := SchemaFromBytes(tablesBytes)
	if err != nil {
		return err
	}

	m.Schema = schema
	return nil
}

func (m *Manager) AddTable(t *Table) error {
	err := m.Schema.AddTable(t)
	if err != nil {
		return err
	}

	tableBytes, err := t.Value()
	if err != nil {
		return fmt.Errorf("failed encoding table while saving to store %w", err)
	}
	err = m.Store.Put(t.Key(), tableBytes)
	if err != nil {
		err = m.Schema.DropTable(t.Key())
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("could not put table definition in store %w", err)
	}
	return nil
}

var MetaTable = &Table{
	Name: MetaTableName,
	Columns: []*Column{
		{
			Name:     "name",
			DataType: STRING,
		},
	},
	PrimaryKey: []string{"name"},
}
