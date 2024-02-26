package flow

type Collectable[T any] func(value T) error

type Collector[T any] interface {
	Emit(value T) error
}
