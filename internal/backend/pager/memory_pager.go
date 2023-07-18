package pager

import "errors"

// MemoryPager implements the Pager interface with an in memory map
// it's mostly around for testing purposes
type MemoryPager struct {
	pageSize uint16
	data     map[uint32][]byte
}

func NewMemoryPager(pageSize uint16) *MemoryPager {
	return &MemoryPager{
		pageSize: pageSize,
		data:     make(map[uint32][]byte),
	}
}

func (m *MemoryPager) ReadPage(pageNum uint32) ([]byte, error) {
	if data, ok := m.data[pageNum]; ok {
		return data, nil
	}
	return nil, errors.New("page not found")
}

func (m *MemoryPager) WritePage(pageNum uint32, data []byte) error {
	m.data[pageNum] = data
	return nil
}

func (m *MemoryPager) PageSize() uint16 {
	return m.pageSize
}
