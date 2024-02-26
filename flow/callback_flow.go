package flow

import "context"

type ProducerScope[T any] interface {
	send(value T)
	sendBlocking(value T)
	awaitClose(cleanup func())
	close()
}

func callbackFlow[T any](ctx context.Context, block func(send ProducerScope[T])) Flow[T] {
	return newCallbackFlowBuilder(ctx, block)
}
