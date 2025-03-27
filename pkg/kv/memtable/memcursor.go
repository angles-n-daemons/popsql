package memtable

import "math"

type Memcursor struct {
	Node *SkiplistNode[string, []byte]
	End  string
}

func (m *Memcursor) ReadAll() ([][]byte, error) {
	return m.Read(math.MaxInt)
}

func (m *Memcursor) Read(num int) ([][]byte, error) {
	vals := [][]byte{}
	for i := 0; i < num && !m.IsAtEnd(); i++ {
		vals = append(vals, m.Node.Val)
		m.Node = m.Node.Next()

	}
	return vals, nil
}

func (m *Memcursor) IsAtEnd() bool {
	return m.Node == nil || m.End <= m.Node.Key
}
