package catalog

import (
	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/meta"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/schema"
)

/*
Bootstrap is a utility function for populating a fresh database.
It adds new meta objects to the schema and the store. It should only be
called if the database is empty, unexpected behavior may occur if executed
otherwise.

Because it depends on the logic required to create system objects from the
manager, it initalizes a temporary, throwaway manager instance to handle the
storage of the meta objects.
*/
func Bootstrap(st kv.Store, sc *schema.Schema, meta *meta.Meta) error {
	tmp := &Manager{
		Schema: sc,
		Store:  st,
		Meta:   meta,
	}
	for _, obj := range tmp.Meta.Objects() {
		var err error
		// hmm..my generic fu isn't strong enough to escape this switch.
		switch v := obj.(type) {
		case *desc.Table:
			err = Create(tmp, v)
		case *desc.Sequence:
			err = Create(tmp, v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
