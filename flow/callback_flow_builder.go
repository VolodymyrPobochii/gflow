package flow

import (
	"context"
	"errors"
	"sync"
)

const UNBUFFERED = 0
const BUFFERED = -1

type ChannelFlow[T any] struct {
	ctx      context.Context
	capacity int
}

type ChannelFlowBuilder[T any] struct {
	ChannelFlow[T]
	ctx      context.Context
	capacity int
	block    func(send ProducerScope[T])
}

type CallbackFlowBuilder[T any] struct {
	//ChannelFlowBuilder[T]
	ctx      context.Context
	capacity int
	block    func(send ProducerScope[T])
}

func (c *CallbackFlowBuilder[T]) collectTo(ps ProducerScope[T]) {
	c.block(ps)
}

func (c *CallbackFlowBuilder[T]) Collect(emit Collectable[T]) {
	cp := newChannelProducer[T](c.ctx, c.capacity)

	go func(emit Collectable[T]) {
		for {
			select {
			case <-c.ctx.Done():
				println("context done")
				close(cp.ch)
				cp.exit <- c.ctx.Err()
				return
			case v, ok := <-cp.ch:
				if !ok {
					println("channel closed")
					cp.exit <- errors.New("channel closed")
					return
				}
				_ = emit(v)
				cp.wg.Done()
			}
		}
	}(emit)

	c.collectTo(cp)
}

func newCallbackBuilder[T any](ctx context.Context, block func(send ProducerScope[T])) *CallbackFlowBuilder[T] {
	return &CallbackFlowBuilder[T]{
		ctx:      ctx,
		block:    block,
		capacity: BUFFERED,
	}
}

type channelProducer[T any] struct {
	ctx  context.Context
	ch   chan T
	exit chan error
	wg   sync.WaitGroup
}

func (c *channelProducer[T]) send(value T) {
	c.wg.Add(1)
	c.ch <- value
}

func (c *channelProducer[T]) sendBlocking(value T) {
	c.wg.Add(1)
	c.ch <- value
}

func (c *channelProducer[T]) awaitClose(cleanup func()) {
	err := <-c.exit
	println(err)
	cleanup()
}

func (c *channelProducer[T]) close() {
	c.wg.Wait()
	close(c.ch)
}

func newChannelProducer[T any](ctx context.Context, capacity int) *channelProducer[T] {
	var ch chan T
	if capacity == BUFFERED {
		ch = make(chan T, 64)
	} else {
		ch = make(chan T)
	}
	return &channelProducer[T]{
		ctx:  ctx,
		ch:   ch,
		exit: make(chan error),
	}
}
