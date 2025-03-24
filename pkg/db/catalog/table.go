package catalog

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema/desc"
)

// TODO: this should take a statement
func (c *Catalog) CreateTable(t *ast.Create) error {
	/*
		In the middle of something not quite working.

		Right now I'm having some trouble conceptualizing how to get
		this to work. On the one hand, I need to create a table from
		a statmenet. This works well if I have an ID, so the created
		table already has an id value. On the other hand, it seems
		better to create the id on the internal step (createTable)
		which is too responsible for:
			- creating an id from the next sequence value
			- adding the table to the schema.
			- adding the table to the store.
	*/
	// fail fast if the table already exists.
	name, err := ast.Identifier(t.Name)
	if err != nil {
		return err
	}
	if _, ok := c.Schema.GetTable(name); ok {
		return fmt.Errorf("table %s already exists", name)
	}

	var dt *desc.Table

	// I need the following for this to work:
	// - if the table doesn't have a primary key, create a sequence and add one.
	dt, err = desc.NewTableFromStmt(t)
	if err != nil {
		return err
	}

	// no primary key created, create a hidden sequence column in its place.
	if len(dt.PrimaryKey) == 0 {
		err := c.createSequence(&desc.Sequence{Name: dt.DefaultSequenceName()})
		if err != nil {
			return err
		}
		// actually use the sequence
	}

	c.createTable(dt)
	return nil
}

func (c *Catalog) createTable(t *desc.Table) error {
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
