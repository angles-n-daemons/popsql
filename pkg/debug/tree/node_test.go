package tree_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/debug/tree"
)

func n8(c string) *tree.Node {
	return &tree.Node{Content: []string{strings.Repeat(c, 8)}}
}

func n8Deep(c string) *tree.Node {
	content := []string{}
	for range 4 {
		content = append(content, strings.Repeat(c, 8))
	}
	return &tree.Node{Content: content}
}

func basicTree() *tree.Node {
	a := n8Deep("A")
	b := n8Deep("B")
	c := n8Deep("C")
	d := n8Deep("D")
	e := n8Deep("E")
	f := n8Deep("F")
	g := n8Deep("G")
	h := n8Deep("H")

	a.AddChild(b)
	a.AddChild(h)
	b.AddChild(c)
	c.AddChild(d)
	c.AddChild(g)
	d.AddChild(e)
	d.AddChild(f)
	return a
}

func TestTreeBasic(t *testing.T) {
	fmt.Println(tree.Visualize(basicTree()))
}
