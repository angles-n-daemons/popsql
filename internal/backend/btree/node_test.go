package btree

import "testing"

func TestSerialization(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name string
		node *Node
	}{
		name: "simple test",
		node: &Node{
			n.nodeType: interior,
			n.pageSize: 32,
			n.order:    3,
			n.cells: []cell{
				{1, 2},
				{3, 4},
				{5, 6},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := test.node.MarshalBinary()

			is.NoErr(err)

			newNode := &Node{}
			newNode.UnmarshalBinary(b)

			is.Equal(test.node, newNode)
		})

	}
}
