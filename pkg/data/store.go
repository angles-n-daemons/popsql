package data

type Store interface {
	Put(key string, value []byte, overwrite bool) error
	Get(start string) ([]byte, error)
	Scan(start string, end string) ([][]byte, error)
}
