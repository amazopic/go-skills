package coroutines

import (
	"context"
	"testing"
)

func TestCoroutine_FiniteSequence(t *testing.T) {
	want := []int{1, 2, 3, 4, 5}
	ctx := context.Background()
	c := Start[int](ctx, func(yield func(int) bool) {
		for _, v := range want {
			if !yield(v) {
				return
			}
		}
	})

	var got []int
	for {
		v, ok := c.Next()
		if !ok {
			break
		}
		got = append(got, v)
	}

	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestCoroutine_Fibonacci(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{name: "first 8 fibonacci", n: 8, want: []int{0, 1, 1, 2, 3, 5, 8, 13}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			fib := Fibonacci(ctx)
			defer fib.Stop()
			got := make([]int, 0, tt.n)
			for range tt.n {
				v, ok := fib.Next()
				if !ok {
					t.Fatal("coroutine ended prematurely")
				}
				got = append(got, v)
			}
			for i, w := range tt.want {
				if got[i] != w {
					t.Errorf("fib[%d] = %d, want %d", i, got[i], w)
				}
			}
		})
	}
}

func TestCoroutine_StopBeforeExhaustion(t *testing.T) {
	// Start an infinite coroutine, consume a few values, stop it.
	// The body checks yield's return value and exits when false.
	// We verify termination by draining until ok=false (yield channel closed).
	ctx := context.Background()
	c := Start[int](ctx, func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return // exit immediately when stopped
			}
		}
	})

	// Consume exactly 3 values.
	for range 3 {
		_, ok := c.Next()
		if !ok {
			t.Fatal("coroutine ended too early")
		}
	}

	// Stop signals the body to return on next yield call.
	c.Stop()

	// Drain until closed — goroutine must exit promptly.
	// If it hangs, the test times out (goroutine leak detection).
	for {
		_, ok := c.Next()
		if !ok {
			break
		}
	}
}

func TestCoroutine_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := Start[int](ctx, func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
		}
	})
	// Consume a couple values then cancel context.
	c.Next()
	c.Next()
	cancel()
	// Drain — body exits when yield returns false due to ctx cancellation.
	for {
		_, ok := c.Next()
		if !ok {
			break
		}
	}
}

func TestCoroutine_EmptyBody(t *testing.T) {
	ctx := context.Background()
	c := Start[string](ctx, func(yield func(string) bool) {
		// yields nothing
	})
	_, ok := c.Next()
	if ok {
		t.Error("expected ok=false for empty body")
	}
}

func TestCoroutine_StopIdempotent(t *testing.T) {
	ctx := context.Background()
	c := Start[int](ctx, func(yield func(int) bool) {
		yield(1) //nolint:errcheck
	})
	c.Stop()
	c.Stop() // must not panic
	c.Stop()
}
