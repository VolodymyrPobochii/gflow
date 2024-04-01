package flow

type Collectable[T any] func(t *T) error

type Collector[T any] interface {
	Emit(t *T) error
}
