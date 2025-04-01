package catalog

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/sys"
)

// NewManager is a constructor for the Manager struct.
//
// The general sequence of events is as follows:
//
//  1. Create a new Schema object, consisting of only system
//     descriptors.
//  2. Using that schema, attempt to load the schema from the
//     store if one exists.
//  3. If there is no data in the store, "bootstrap" or store
//     the system schema from step 1. Return the manager created
//     in bootstrap.
//  4. Otherwise, return a new manager created from the loaded
//     (2) schema and passed in store.
func NewManager(st kv.Store) (*Manager, error) {
	initSchema, err := sys.InitSchema()
	if err != nil {
		return nil, err
	}

	sc, err := LoadSchema(initSchema, st)
	if err != nil {
		return nil, err
	}

	if schema.Empty(sc) {
		// Set the schema to be the init schema and persist it.
		sc = initSchema
		err = Bootstrap(st, sc)
		if err != nil {
			return nil, err
		}
	}

	return &Manager{
		Schema: sc,
		Store:  st,
	}, nil
}

func LoadSchema(sc *schema.Schema, st kv.Store) (*schema.Schema, error) {
	tables, err := LoadCollection[*desc.Table](sc, st)
	if err != nil {
		return nil, err
	}

	sequences, err := LoadCollection[*desc.Sequence](sc, st)
	if err != nil {
		return nil, err
	}
	return schema.SchemaFromCollections(tables, sequences), nil
}

func LoadCollection[V desc.Any[V]](sc *schema.Schema, st kv.Store) (*schema.Collection[V], error) {
	var zero V
	span := getSystemTable[V](sc).Span()
	cur, err := st.GetRange(span.Start.Encode(), span.End.Encode())
	if err != nil {
		return nil, err
	}

	bytesArr, err := cur.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to '%T' descriptors from a cursor %w", zero, err)
	}

	return schema.CollectionFromBytes[V](bytesArr)
}
