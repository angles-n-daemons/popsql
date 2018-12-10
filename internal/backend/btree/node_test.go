package btree

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/matryer/is"
)

func TestSerialization(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name     string
		node     *Node
		notEqual bool
		errors   bool
	}{
		{
			name: "simple test",
			node: &Node{
				nodeType: interior,
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
				nodeType: leaf,
				pageSize: 32,
				order:    4,
				cells:    []cell{},
			},
		},
		{
			name: "obscure page size",
			node: &Node{
				nodeType: leaf,
				pageSize: 37,
				order:    4,
				cells:    []cell{},
			},
		},

		// not equal
		{
			name: "unrealistic order",
			node: &Node{
				nodeType: interior,
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
				nodeType: nodeType(0),
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

			is.NoErr(err)

			newNode := &Node{}
			err = newNode.UnmarshalBinary(b)

			if test.errors {
				is.True(err != nil)
			} else {
				is.NoErr(err)
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
