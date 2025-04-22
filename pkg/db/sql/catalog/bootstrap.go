package catalog

import (
	"github.com/angles-n-daemons/popsql/pkg/db/kv"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/catalog/schema"
)

/*
Bootstrap is a utility function for populating a fresh database.
It adds new system objects to the store. It should only be called
if the database is empty, unexpected behavior may occur if
executed otherwise.

Because it depends on the logic required to create system objects
from the manager, it initalizes a temporary, throwaway manager
instance to handle the storage of the meta objects.
*/
func Bootstrap(st kv.Store, sc *schema.Schema) error {
	tmp := &Manager{
		Schema: sc,
		Store:  st,
	}
	for _, t := range sc.Tables.All() {
		err := save(tmp, t)
		if err != nil {
			return err
		}
	}
	for _, s := range sc.Sequences.All() {
		err := save(tmp, s)
		if err != nil {
			return err
		}
	}
	return nil
}
