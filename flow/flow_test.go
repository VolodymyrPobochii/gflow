package flow

import (
	"context"
	"testing"
	"time"
)

func TestNewFlow(t *testing.T) {
	ctx := context.Background()
	//getFlow(ctx).Collect(func(value string) error {
	//	println("got value = ", value)
	//	return nil
	//})

	ctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
	defer cancel()

	getCBFlow(ctx).Collect(func(v string) error {
		println("got value = ", v)
		return nil
	})
}

func getFlow(ctx context.Context) Flow[string] {
	return flow[string](ctx, func(emit Collectable[string]) {
		println("sending value...")
		for i := 0; i < 10; i++ {
			if err := emit(string(rune(i+35)) + ": hey there!"); err != nil {
				println(err)
			}
			time.Sleep(time.Second)
		}
		println("done sending")
	})
}

func getCBFlow(ctx context.Context) Flow[string] {
	return callbackFlow[string](ctx, func(ps ProducerScope[string]) {
		for i := 0; i < 100; i++ {
			go func(i int) {
				value := string(rune(i+35)) + ": hey there!"
				ps.send(value)
				println("value sent:", value)
			}(i)
		}

		ps.close()

		ps.awaitClose(func() {
			println("cleaned up")
		})
	})
}
