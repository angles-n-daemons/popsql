package catalog

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/meta"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
)

// NewManager is a constructor for the Manager struct.
//
// On its own, it's a fairly gnarly startup sequence, either cold or hot.
// The first tricky bit is that the system has a chicken and egg problem.
// Ideally, the meta tables (see meta.go) are read from the story, but
// in order to read anything from the store, the key spans of the meta
// tables are required.
//
// We get around this by initalizing an new meta object, which is only
// to be used for startup.
func NewManager(st kv.Store) (*Manager, error) {
	initMeta := meta.InitSystemMeta()

	sc, err := LoadSchema(initMeta, st)
	if err != nil {
		return nil, err
	}

	if schema.Empty(sc) {
		err = Bootstrap(st, sc, initMeta)
		if err != nil {
			return nil, err
		}
	}

	// At this point, the meta tables should exist, either
	// from load into the schema, or from explicit bootstrap.
	meta, err := meta.FromSchema(sc)
	if err != nil {
		return nil, err
	}

	return &Manager{
		Schema: sc,
		Store:  st,
		Meta:   meta,
	}, nil
}

func LoadSchema(meta *meta.Meta, st kv.Store) (*schema.Schema, error) {
	tables, err := LoadCollection[*desc.Table](meta.Tables.Table.Span(), st)
	if err != nil {
		return nil, err
	}

	sequences, err := LoadCollection[*desc.Sequence](meta.Sequences.Table.Span(), st)
	if err != nil {
		return nil, err
	}
	return schema.SchemaFromCollections(tables, sequences), nil
}

func LoadCollection[V desc.Object[V]](span *keys.Span, st kv.Store) (*schema.Collection[V], error) {
	var zero V
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
