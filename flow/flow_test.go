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

	getCBFlow(ctx).Collect(func(v *string) error {
		println("got value = ", *v)
		return nil
	})
}

func getFlow(ctx context.Context) Flow[string] {
	return New[string](ctx, func(emit Collectable[string]) {
		println("sending value...")
		for i := 0; i < 10; i++ {
			s := string(rune(i+35)) + ": hey there!"
			if err := emit(&s); err != nil {
				println(err)
			}
			time.Sleep(time.Second)
		}
		println("done sending")
	})
}

func getCBFlow(ctx context.Context) Flow[string] {
	return NewCallback[string](ctx, func(ps ProducerScope[string]) {
		for i := 0; i < 100; i++ {
			go func(i int) {
				v := string(rune(i+35)) + ": hey there!"
				ps.Send(&v)
				println("value sent:", v)
			}(i)
		}

		ps.Close()

		ps.AwaitClose(func(_ error) {
			println("cleaned up")
		})
	})
}
