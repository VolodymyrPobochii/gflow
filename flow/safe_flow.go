package flow

import "context"

type safeFlow[T any] struct {
	ctx   context.Context
	block func(emit Collectable[T])
}

func newSafeFlow[T any](ctx context.Context, block func(emit Collectable[T])) *safeFlow[T] {
	return &safeFlow[T]{ctx: ctx, block: block}
}

func (f *safeFlow[T]) Collect(emit Collectable[T]) {
	safeCollector := &safeCollector[T]{f.ctx, emit}
	f.block(safeCollector.safeEmit)
}
