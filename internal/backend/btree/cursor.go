package btree

type pager interface {
	GetPage(offset int) ([]byte, error)
	SetPage(offset int, content []byte) error
}

/*
	NewCursor returns a btree cursor implementation by reading in
	the root page and returning a cursor with the root node and
	the pager.
*/
func NewCursor(rootPage int, pager pager) (*Cursor, error) {
	page, err := pager.GetPage(rootPage)
	if err != nil {
		return nil, err
	}

	var n Node
	err = n.UnmarshalBinary(page)
	if err != nil {
		return nil, err
	}

	return &Node{
		node:  node,
		pager: pager,
	}, nil
}

/*
	Cursor is the handle to the BTree and implements its interface.

	It functions by reading pages from the pager, and making decisions
	based on the properties of each individual node.

*/
type Cursor struct {
	node Node
	pager
}

/*
	Search is a tree traversal implementation which iteratively goes
	down indexed btree pages for the referenced key.

	If the key is found, the stored value is returned.

	If the key is not found, a nil value with no error is returned.
*/
func (c *Cursor) Search(key uint32) (interface{}, error) {

	return nil, nil
}

func (c *Cursor) Insert(key uint32, value interface{}) error {
	return nil
}

func (c *Cursor) Delete(key uint32) error {
	return nil
}
