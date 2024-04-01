package flow

import "context"

type ProducerScope[T any] interface {
	Send(t *T)
	SendBlocking(t *T)
	AwaitClose(cleanup func(error))
	Close()
}

func NewCallback[T any](ctx context.Context, block func(send ProducerScope[T])) Flow[T] {
	return newCallbackBuilder(ctx, block)
}
