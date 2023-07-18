package btree

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSerialization(t *testing.T) {
	tests := []struct {
		name     string
		node     *Node
		notEqual bool
		errors   bool
	}{
		{
			name: "simple test",
			node: &Node{
				NodeType: table_interior,
				pageSize: 32,
				order:    4,
				cells: []cell{
					{1, 2},
					{3, 4},
					{5, 6},
				},
			},
		},
		{
			name: "empty node",
			node: &Node{
				NodeType: table_leaf,
				pageSize: 32,
				order:    4,
				cells:    []cell{},
			},
		},
		{
			name: "obscure page size",
			node: &Node{
				NodeType: table_leaf,
				pageSize: 37,
				order:    4,
				cells:    []cell{},
			},
		},

		// not equal
		{
			name: "unrealistic order",
			node: &Node{
				NodeType: table_interior,
				pageSize: 32,
				order:    6,
				cells:    []cell{},
			},
			notEqual: true,
		},

		// errors
		{
			name: "invalid node type",
			node: &Node{
				NodeType: NodeType(0),
				pageSize: 32,
				order:    4,
				cells:    []cell{},
			},
			errors:   true,
			notEqual: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := test.node.MarshalBinary()

			if err != nil {
				t.Errorf("Failed to unmarhshal value %v", err)
			}

			newNode := &Node{}
			err = newNode.UnmarshalBinary(b)

			if test.errors && err == nil {
				t.Errorf("Expected error in test '%s' but received none", test.name)
			} else if !test.errors && err != nil {
				t.Errorf("Expected no error in test '%s' but got %v", test.name, err)
			}
			isEqual := reflect.DeepEqual(test.node, newNode)

			if test.notEqual == isEqual {
				b1, _ := json.Marshal(test.node)
				b2, _ := json.Marshal(newNode)

				t.Fatalf(
					"Expected equality: %t\ninput:%s\noutput:%s",
					!test.notEqual,
					b1,
					b2,
				)
			}
		})

	}
}
