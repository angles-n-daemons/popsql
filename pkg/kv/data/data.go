package data

type Store interface {
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
	GetRange(start, end string) (Cursor, error)
	PutRange(start, end string, value []byte) error
}

type Cursor interface {
	ReadAll() ([][]byte, error)
	Read(num int) ([][]byte, error)
	IsAtEnd() bool
}
