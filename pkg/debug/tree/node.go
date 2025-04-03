package tree

import "strings"

type Node struct {
	Content  []string
	Children []*Node
}

func NewNode(content []string) *Node {
	return &Node{
		Content:  content,
		Children: []*Node{},
	}
}

func (n *Node) AddChild(c *Node) {
	n.Children = append(n.Children, c)
}

type Direction int

const (
	UP Direction = iota
	DOWN
)

const max_depth = 30

func Visualize(n *Node) string {
	lines := (&Visualizer{
		depth: 0,
		open:  make([]bool, max_depth),
		pad:   8,
		box:   true,
	}).Render(n)
	return strings.Join(lines, "\n")
}

type Visualizer struct {
	depth int

	// below settings aren't used currently
	pad       int
	box       bool
	open      []bool
	direction Direction
}

/*
reference characters
┌─┬─┐  ╔═╦═╗  ╭─┬─╮
│ │ │  ║ ║ ║  │ │ │
├─┼─┤  ╠═╬═╣  ├─┼─┤
└─┴─┘  ╚═╩═╝  ╰─┴─╯

→ ← ↑ ↓ ↔ ↕
⇒ ⇐ ⇑ ⇓ ⇄ ⇅

├── (branch)
└── (last branch)
│   (vertical)
*/

func (v *Visualizer) Render(n *Node) []string {
	result := []string{}
	content := prepContent(n.Content, v.box)
	if n == nil {
		return []string{"tree node should not be nil"}
	}

	if v.depth >= max_depth {
		return []string{"max depth exceeded"}
	}

	height := len(content)
	for row, line := range content {
		result = append(result, v.Indent(row, height)+line)
	}

	for i, child := range n.Children {
		// if i < len(n.Children)-1, set open to true for this depth
		// otherwise set to false
		// open is guaranteed to be closed if opened
		if i < len(n.Children)-1 {
			v.open[v.depth] = true
		} else {
			v.open[v.depth] = false
		}
		v.depth++
		cLines := v.Render(child)
		result = append(result, cLines...)
		v.depth--
	}
	return result
}

// r, h stand for row and height respectively
func (v *Visualizer) Indent(row int, height int) string {
	pad := v.pad
	edgeHeight := (height - 1) / 2
	indent := ""
	// pre-current intent
	for i := range v.depth {
		// either before the current depth, or after the edge.
		if i < v.depth-1 || row > edgeHeight {
			if v.open[i] {
				indent += " │" + strings.Repeat(" ", pad-2)
			} else {
				indent += strings.Repeat(" ", pad)
			}
		} else if row < edgeHeight {
			// at where the edge should intersect
			indent += " │" + strings.Repeat(" ", pad-2)
		} else {
			// edge condition
			edge := strings.Repeat("─", pad-3) + "▶"
			if v.open[v.depth-1] {
				indent += " ├" + edge
			} else {
				indent += " ╰" + edge

			}
		}
	}
	return indent
}

func prepContent(lines []string, box bool) []string {
	if !box {
		lines = append([]string{""}, lines...)
		return lines
	}
	longest := 0
	for _, line := range lines {
		if len(line) > longest {
			longest = len(line)
		}
	}
	x := longest
	top := "┌" + strings.Repeat("─", x) + "┐"
	bottom := "└" + strings.Repeat("─", x) + "┘"
	result := []string{top}
	for _, line := range lines {
		result = append(result, "│"+line+"│")
	}
	result = append(result, bottom)
	return result
}
