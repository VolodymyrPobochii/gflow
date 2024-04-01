package flow

import "context"

type Flow[T any] interface {
	Collect(emit Collectable[T])
}

func New[T any](ctx context.Context, block func(emit Collectable[T])) Flow[T] {
	return newSafe[T](ctx, block)
}
