package desc

// Object is the interface that all descriptors must implement.
type Object[V any] interface {
	WithID(id uint64)
	ID() uint64
	Key() string
	Name() string
	Equal(o V) bool
}
