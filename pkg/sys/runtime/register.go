package runtime

type Register interface {
	Key() string
	Value() ([]byte, error)
}
