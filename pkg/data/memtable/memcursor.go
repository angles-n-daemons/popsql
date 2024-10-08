package memtable

import "math"

type Memcursor struct {
	Node *SkiplistNode[string, []byte]
	End  string
}

func (m *Memcursor) ReadAll(num int) ([][]byte, error) {
	return m.ReadAll(math.MaxInt)
}

func (m *Memcursor) Read(num int) ([][]byte, error) {
	vals := [][]byte{}
	for i := 0; i < num; i++ {
		if m.IsAtEnd() {
			return vals, nil
		}
		vals = append(vals, m.Node.Val)
		m.Node = m.Node.Next()

	}
	return vals, nil
}

func (m *Memcursor) IsAtEnd() bool {
	return m.Node == nil || m.End <= m.Node.Key
}
