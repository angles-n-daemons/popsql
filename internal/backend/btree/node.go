package btree

import (
	"encoding/binary"
	"fmt"
)

func NewNode(nodeType NodeType) {
	return
}

/*
	Node is the BTree logical implementation

	order: btree order which is the same number of pointers
	       in the node and the number of keys plus one
	cells: key pointer pairs ordered by the key value
	rightPointer: the page number of the node to the right
*/
type Node struct {
	NodeType
	order        int
	cells        []cell
	rightPointer uint32
	pageSize     int
}

/*
	UnmarshalBinary reads in a page and attempts to return a typed
	Node object. It accomplishes by executing the following steps.

	1. Read the node header (8 or 12 bytes).
	2. Set the node type
	4. Set the order based on page size
	3. Read the cells
	4. Set the page size for when the node needs to be serialized
*/
func (n *Node) UnmarshalBinary(data []byte) error {
	if len(data) < 16 {
		return fmt.Errorf(
			"Node page too small: %d",
			len(data),
		)
	}

	/*
		nodeHeader
		0:   node type
		1-2: numCells
		3-6: rightmost pointer
		7:   unused
	*/
	nodeHeader := data[:8]
	n.NodeType = NodeType(nodeHeader[0])
	numCells := binary.LittleEndian.Uint16(nodeHeader[1:3])
	n.rightPointer = binary.LittleEndian.Uint32(nodeHeader[3:7])

	if n.NodeType != table_interior && n.NodeType != table_leaf {
		return fmt.Errorf("Unknown node type: %d", n.NodeType)
	}

	n.cells = make([]cell, numCells)

	for i := uint16(0); i < numCells; i++ {
		// Read cells in 8 byte increments
		offset := (i * 8) + 8
		cellBytes := data[offset : offset+8]

		c := cell{
			key:  binary.LittleEndian.Uint32(cellBytes[:4]),
			page: binary.LittleEndian.Uint32(cellBytes[4:8]),
		}
		n.cells[i] = c
	}

	// order is the number of pointers which is equal to
	// the number of cells plus 1
	n.order = ((len(data) - 8) / 8) + 1

	// pagesize stored for when the node is reserialized
	n.pageSize = len(data)

	return nil
}

/*
	MarshalBinary performs the opposite action of UnmarshalBinary
*/
func (n *Node) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0, n.pageSize)
	data = append(data, byte(n.NodeType))

	// append the numCells bytes
	numCells16 := uint16(len(n.cells))
	numCellsBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(numCellsBytes, numCells16)
	data = append(data, numCellsBytes...)

	// append the right pointer
	rightPointerBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(rightPointerBytes, n.rightPointer)
	data = append(data, rightPointerBytes...)

	// append the empty header byte
	data = append(data, 0)

	// add the cells
	for _, cell := range n.cells {
		keyBytes := make([]byte, 4)
		pointerBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(keyBytes, cell.key)
		binary.LittleEndian.PutUint32(pointerBytes, cell.page)
		data = append(data, keyBytes...)
		data = append(data, pointerBytes...)
	}

	// fill the remaining space
	data = append(
		data,
		make([]byte, cap(data)-len(data))...,
	)

	return data, nil
}

/*
	NodeType describes whether a node is a leaf or not
*/
type NodeType byte

const (
	table_interior NodeType = 1
	table_leaf
)

/*
	cell is a key pointer pair in a btree node

	key: in a leaf node it's the exact value of the corresponding data. in an interior node, the key represents the left-most (lowest) key of the child that it points to.
	page: page number for the key. all children's keys will
	         be less than or equal to the key value


	serialized it is laid out on disk [key, pointer]
*/
type cell struct {
	key  uint32
	page uint32
}
