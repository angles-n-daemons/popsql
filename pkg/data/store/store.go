package data

type Store interface {
	Get(string) ([]byte, error)
	GetRange(string, string) Cursor
	Set(string, []byte) error
	SetRange(string, string, []byte) error
}

type Cursor interface {
	Read(int) ([]byte, bool, error)
	IsAtEnd() bool
	// Close?
}
