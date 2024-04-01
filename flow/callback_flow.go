package flow

import "context"

type ProducerScope[T any] interface {
	Send(value T)
	SendBlocking(value T)
	AwaitClose(cleanup func())
	Close()
}

func NewCallback[T any](ctx context.Context, block func(send ProducerScope[T])) Flow[T] {
	return newCallbackBuilder(ctx, block)
}
