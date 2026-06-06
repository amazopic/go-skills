// Package fanout implements the Fan-Out messaging pattern: distribute work
// from a single source channel across N worker goroutines.
//
// Items on the source channel are received by exactly one worker, so the
// runtime balances load automatically — fast workers pull more items than slow
// ones (work-stealing semantics) rather than following a fixed round-robin
// assignment. Use Fan-Out to scale throughput when work items are independent
// and a single consumer is the bottleneck.
//
// The package is generic over the input and output element types and is
// context-first: every blocking operation is cancellable, and no goroutine is
// left running after the source drains or the context is cancelled.
package fanout

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrInvalidWorkers is returned by FanOutErr when n < 1.
var ErrInvalidWorkers = errors.New("fanout: workers must be >= 1")

// FanOut starts n worker goroutines, each draining the shared in channel and
// emitting process(item) onto the returned channel. Results from different
// workers are interleaved in arrival order; no ordering relative to in is
// guaranteed.
//
// The returned channel is closed exactly once, after every worker has exited —
// which happens when in is drained (closed by the producer) or ctx is
// cancelled. Callers must range over the result to completion (or cancel ctx)
// to avoid leaking the workers blocked on send.
//
// FanOut panics if n < 1; use FanOutErr for a non-panicking variant.
func FanOut[In, Out any](ctx context.Context, n int, in <-chan In, process func(In) Out) <-chan Out {
	if n < 1 {
		panic("fanout: workers must be >= 1")
	}

	out := make(chan Out)

	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-in:
					if !ok {
						return // source drained
					}
					select {
					case out <- process(item):
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	// A single closer goroutine owns closing out: it runs only after every
	// worker has returned, so all sends have completed and the channel is
	// closed exactly once. No worker ever closes out.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// FanOutErr is FanOut without the panic: it returns ErrInvalidWorkers when
// n < 1 and otherwise behaves identically.
func FanOutErr[In, Out any](ctx context.Context, n int, in <-chan In, process func(In) Out) (<-chan Out, error) {
	if n < 1 {
		return nil, ErrInvalidWorkers
	}
	return FanOut(ctx, n, in, process), nil
}

// Collect drains ch into a slice, returning early with ctx.Err() (wrapped) if
// the context is cancelled before ch is closed. It is a convenience for callers
// that want all results materialised rather than streamed.
func Collect[T any](ctx context.Context, ch <-chan T) ([]T, error) {
	var out []T
	for {
		select {
		case <-ctx.Done():
			return out, fmt.Errorf("fanout: collect cancelled: %w", ctx.Err())
		case v, ok := <-ch:
			if !ok {
				return out, nil
			}
			out = append(out, v)
		}
	}
}
