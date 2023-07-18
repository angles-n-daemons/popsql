package pager

import (
	"io"
	"os"
)

type File interface {
	io.ReaderAt
	io.WriterAt
}

// FilePager struct implementing Pager using a file handle
type FilePager struct {
	file     File
	pageSize uint16
}

func NewFilePager(filePath string, pageSize uint16) (*FilePager, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &FilePager{
		file:     file,
		pageSize: pageSize,
	}, nil
}

func (f *FilePager) ReadPage(pageNum uint32) ([]byte, error) {
	pageOffset := int64(pageNum) * int64(f.pageSize)
	pageData := make([]byte, f.pageSize)

	_, err := f.file.ReadAt(pageData, pageOffset)
	if err != nil {
		return nil, err
	}

	return pageData, nil
}

func (f *FilePager) WritePage(pageNum uint32, data []byte) error {
	pageOffset := int64(pageNum) * int64(f.pageSize)
	_, err := f.file.WriteAt(data, pageOffset)
	if err != nil {
		return err
	}

	return nil
}

func (f *FilePager) PageSize() uint16 {
	return f.pageSize
}
