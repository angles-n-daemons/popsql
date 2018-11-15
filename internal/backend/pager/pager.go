package pager

import "fmt"

const (
	RESERVED_HEADER_PAGE_COUNT = 1
)

type file interface {
	Read(int, int) ([]byte, error)
	Write(int, []byte) error
}

/*
	Pager manages the database pages for access by the BTree
Properties
	 - pageSize:      the number of bytes in each data page
	 - fs:            file system operator used for reading and writing to the db
	 - freelistIndex: the page number where the freelist starts

	 The database has the following pages:

	 Header:     Metadata about the database
	 TablePages: Data pages for the database
	 Freelist:   Unused pages to be allocated
*/
type Pager struct {
	pageSize      int
	fs            file
	freelistIndex int
}

/*
	Get retrieves the data page at the offset after the header
*/
func (p *Pager) GetPage(offset int) ([]byte, error) {
	if offset < 0 || offset >= p.freelistIndex {
		return nil, fmt.Errorf(
			"cannot read page %d, index must be between 0 and %d",
			offset,
			p.freelistIndex,
		)
	}

	offsetBytes := (offset * p.pageSize) + (RESERVED_HEADER_PAGE_COUNT * p.pageSize)
	return p.fs.Read(offsetBytes, p.pageSize)
}

/*
	Set puts a data page at the offset with the bytes in content
*/
func (p *Pager) SetPage(offset int, content []byte) error {
	if offset < 0 || offset >= p.freelistIndex {
		return fmt.Errorf(
			"cannot write page %d, index must be between 0 and %d",
			offset,
			p.freelistIndex,
		)
	}

	if len(content) != p.pageSize {
		return fmt.Errorf(
			"content length %d larger than pagesize %d",
			len(content),
			p.pageSize,
		)
	}

	offsetBytes := (offset * p.pageSize) + (RESERVED_HEADER_PAGE_COUNT * p.pageSize)
	return p.fs.Write(offsetBytes, content)
}
