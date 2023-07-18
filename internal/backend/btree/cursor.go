package btree

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/internal/backend/pager"
)

type Cursor struct {
	pager pager.Pager
}

func (c *Cursor) Create(pageNum uint16) error {
	_, ok := c.pager.ReadPage(pageNum)
	if ok == nil {
		return fmt.Errorf("cannot create b-tree at page %d, page is in use", pageNum)
	}

	// TODO: finish implementation

	return nil
}

func (c *Cursor) Search(key uint32) (interface{}, error) {
	return nil, nil
}

func (c *Cursor) Insert(key uint32, value interface{}) error {
	return nil
}

func (c *Cursor) Delete(key uint32) error {
	return nil
}
