package pager

import (
	"fmt"
	"testing"
)

// MockFile represents a mocked file handle for testing
type MockFile struct {
	data map[int64][]byte
}

func NewMockFile() *MockFile {
	return &MockFile{
		data: make(map[int64][]byte),
	}
}

func (m *MockFile) ReadAt(data []byte, offset int64) (int, error) {
	copy(data, m.data[offset])
	return len(m.data[offset]), nil
}

func (m *MockFile) WriteAt(data []byte, offset int64) (int, error) {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	m.data[offset] = dataCopy
	return len(data), nil
}

func (m *MockFile) Close() error {
	return nil
}

func TestFilePager_ReadPage(t *testing.T) {
	mockFile := NewMockFile()
	filePager := &FilePager{
		file:     mockFile,
		pageSize: 12,
	}

	pageNum := uint32(0)
	expectedData := []byte("Page 0 data ")
	mockFile.data[int64(pageNum)*int64(filePager.pageSize)] = expectedData

	data, err := filePager.ReadPage(pageNum)
	if err != nil {
		t.Errorf("Error reading page from FilePager: %v", err)
	}

	if string(data) != string(expectedData) {
		fmt.Println("expected length ", len(expectedData))
		fmt.Println("actual length ", len(data))
		t.Errorf("ReadPage returned incorrect data. Expected: %s, Got: %s", string(expectedData), string(data))
	}
}

func TestFilePager_WritePage(t *testing.T) {
	mockFile := NewMockFile()
	filePager := &FilePager{
		file:     mockFile,
		pageSize: 4096,
	}

	pageNum := uint32(0)
	data := []byte("Page 0 data")

	err := filePager.WritePage(pageNum, data)
	if err != nil {
		t.Errorf("Error writing page to FilePager: %v", err)
	}

	offset := int64(pageNum) * int64(filePager.pageSize)
	if len(mockFile.data[offset]) != len(data) {
		t.Errorf("WritePage did not write the correct amount of data. Expected: %d bytes, Got: %d bytes", len(data), len(mockFile.data[offset]))
	}

	if string(mockFile.data[offset]) != string(data) {
		t.Errorf("WritePage did not write the correct data. Expected: %s, Got: %s", string(data), string(mockFile.data[offset]))
	}
}

func TestFilePager_PageSize(t *testing.T) {
	mockFile := NewMockFile()
	filePager := &FilePager{
		file:     mockFile,
		pageSize: 4096,
	}

	pageSize := filePager.PageSize()
	if pageSize != 4096 {
		t.Errorf("PageSize returned incorrect value. Expected: 4096, Got: %d", pageSize)
	}
}
