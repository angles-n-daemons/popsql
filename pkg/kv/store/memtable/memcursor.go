package memtable

import "math"

// Memcursor is a simple Cursor implementation using a memory
// store.
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

func (m *Memcursor) Next() ([]byte, error) {
	if m.IsAtEnd() {
		return nil, nil
	}
	val := m.Node.Val
	m.Node = m.Node.Next()
	return val, nil
}

func (m *Memcursor) IsAtEnd() bool {
	return m.Node == nil || m.End <= m.Node.Key
}
