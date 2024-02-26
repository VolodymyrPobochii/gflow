package flow

import "context"

type Flow[T any] interface {
	Collect(emit Collectable[T])
}

func flow[T any](ctx context.Context, block func(emit Collectable[T])) Flow[T] {
	return newSafeFlow[T](ctx, block)
}
