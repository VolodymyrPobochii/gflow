package flow

import "context"

type safeCollector[T any] struct {
	ctx  context.Context
	emit Collectable[T]
}

func (s *safeCollector[T]) safeEmit(value T) error {
	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	default:
		return s.emit(value)
	}
}
