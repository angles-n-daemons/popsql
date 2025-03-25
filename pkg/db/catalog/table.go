package catalog

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sys/schema/desc"
)

func (c *Catalog) CreateTable(t *desc.Table) error {
	if len(t.PrimaryKey) == 0 {
		s, err := t.AddInternalPrimaryKey()
		if err != nil {
			return err
		}
		s, err = c.addSequence(s)
		if err != nil {
			return err
		}
	}

	return c.addTable(t)
}

func (c *Catalog) addTable(t *desc.Table) error {
	id, err := c.SequenceNext(c.metaTableSequence)
	if err != nil {
		return err
	}
	t.ID = id

	err = c.Schema.AddTable(t)
	if err != nil {
		return err
	}

	err = c.storeTable(t)
	if err != nil {
		err = c.Schema.DropTable(t.Key())
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("could not put table definition in store %w", err)
	}
	return nil
}

func (c *Catalog) storeTable(t *desc.Table) error {
	key := c.metaTable.Prefix().WithID(t.Key())
	tableBytes, err := t.Value()
	if err != nil {
		return fmt.Errorf("failed encoding table while saving to store %w", err)
	}
	err = c.Store.Put(key.Encode(), tableBytes)
	if err != nil {
		return fmt.Errorf("could not put table definition in store %w", err)
	}
	return nil
}
