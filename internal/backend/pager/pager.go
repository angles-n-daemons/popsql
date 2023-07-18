package pager

/*
 Pager is the interface with which popsql will access the database files
*/
type Pager interface {
	ReadPage(pageNum uint16) ([]byte, error)
	WritePage(pageNum uint16, data []byte) error
	PageSize() uint16
}
