package data

type Store interface {
	Get(string) ([][]byte, error)
	GetRange(string, string) Cursor
	Put(string, []byte) error
	PutRange(string, string, []byte) error
}
