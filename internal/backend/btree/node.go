package btree

import "fmt"

type pager interface {
	PageSize() int
}

/*
	Node is the BTree logical implementation

	order: btree order which is the same number of pointers
	       in the node and the number of keys plus one
	cells: key pointer pairs ordered by the key value
*/
type Node struct {
	NodeType
	Order int
	Cells []cell
	Pager pager
}

/*
	UnmarshalBinary reads in a page and attempts to return a typed
	Node object. It accomplishes by doing the following things.

	1. Read the node header (8 or 12 bytes).
	2. Set the node type
	4. Set the order based on page size
	3. Read the cells
	4. Set the page size for when the node needs to be serialized
*/
func (n *Node) UnmarshalBinary(data []byte) error {
	if len(data) < 6 {
		return fmt.Errorf(
			"Unexpected length of Node page: %d",
			len(data),
		)
	}

	n.pageSize = len(data)

	/*
		nodeHeader
		0:   node type
		1-2: numCells
		3-6: rightmost pointer
		7:   unused
	*/
	nodeHeader := data[:8]
	n.nodeType = nodeType(nodeHeader[0])
	numCells := int16(nodeHeader[1:3])
	n.rightPointer = int32(nodeHeader[3:7])

	if n.nodeType != interior && n.nodeType != leaf {
		return fmt.Errorf("Unknown node type: %d", n.nodeType)
	}

	n.cells = make([]cell, 0, numCells)

	for i := 0; i < numCells; i++ {
		// Read cells in 8 byte increments
		offset := (i * 8) + 8
		cellBytes := data[offset : offset+8]

		n.cells[i] = cell{
			key:     int32(cellBytes[:4]),
			pointer: int32(cellBytes[4:8]),
		}
	}

	// order is the number of pointers which is equal to
	// the number of cells plus 1
	n.order = ((data - 8) / 8) + 1
}

/*
	MarshalBinary performs the opposite action of UnmarshalBinary
*/
func (n *Node) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0, n.pageSize)
	data = append(data, n.nodeType)

	numCells := int16(len(n.cells))
	// append the numCells bytes
	data = append(data, []byte(numCells)...)
	// append the right pointer
	data = append(data, []byte(n.rightPointer)...)
	// append the empty header byte
	data = append(data, 0)

	for _, cell := range n.cells {
		data = append(data, []byte(cell.key)...)
		data = append(data, []byte(cell.pointer)...)
	}

	// fill the remaining space
	data = append(
		data,
		make([]byte, cap(data)-len(data)),
	)
}

/*
	nodeType describes whether a node is a leaf or not
*/
type nodeType byte

const (
	interior nodeType = 1
	leaf
)

/*
	cell is a key pointer pair in a btree node

	key: corresponds to a value for comparison
	pointer: left pointer for the key. all children's keys will
	         be less than or equal to the key value


	serialized it is laid out on disk [key, pointer]
*/
type cell struct {
	key     int32
	pointer int32
}
