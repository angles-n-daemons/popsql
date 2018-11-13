package pager

import "fmt"

const (
	RESERVED_HEADER_PAGE_COUNT = 1
)

type file interface {
	Read(int) ([]byte, error)
	Write(int, []byte) error
	Size() (int, error)
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
	Get retrieves the data page at offset i after the header
*/
func (p *Pager) Get(int i) ([]byte, error) {
	if i < 1 || i > p.freelistIndex {
		return nil, fmt.Errorf(
			"cannot read page %d, index must be between 1 and %d",
			i,
			p.freelistIndex,
		)
	}

	offset = RESERVED_HEADER_PAGE_COUNT * i * pageSize
	return p.fs.Read(offset)
}

/*
	Set puts a data page at offset i with the bytes in content
*/
func (p *Pager) Set(pageNumber int, content []byte) error {
	if pageNumber < 1 || pageNumber > p.freelistIndex {
		return nil, fmt.Errorf(
			"cannot write page %d, index must be between 1 and %d",
			pageNumber,
			p.freelistIndex,
		)
	}

	offset = RESERVED_HEADER_PAGE_COUNT * pageNumber * pageSize
	return p.fs.Write(offset, content)
}
