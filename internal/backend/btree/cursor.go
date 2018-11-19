package btree

type Cursor struct {
	Node
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
