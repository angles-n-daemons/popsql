package data

type Cursor interface {
	Read(int) ([][]byte, error)
	ReadAll() ([][]byte, error)
	IsAtEnd() bool
}
