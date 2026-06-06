// Package generator implements the Generator concurrency pattern: a goroutine
// lazily produces a sequence of values and emits them on a channel, returned to
// the caller as a receive-only <-chan T.
//
// The receive-only return type hides the channel's send/close internals so
// callers can only range over the stream. Every generator is context-aware: if
// the caller abandons the stream early it must cancel the context (or the
// producer would block forever on send and leak its goroutine). Each generator
// owns its channel and closes it exactly once when the sequence is exhausted or
// the context is cancelled, so `for v := range gen(...)` always terminates.
package generator

import "context"

// Generate runs produce in a new goroutine and returns the receive-only stream
// it writes to. produce is invoked with a yield function; calling yield(v)
// emits v on the stream and blocks until either the consumer receives it or ctx
// is cancelled. yield reports false once ctx is cancelled, at which point
// produce should return promptly to release the goroutine.
//
// The returned channel is closed when produce returns (normally or because
// yield reported cancellation), so ranging over it always terminates. Callers
// that stop ranging early MUST cancel ctx to avoid leaking the producer.
func Generate[T any](ctx context.Context, produce func(yield func(T) bool)) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		produce(func(v T) bool {
			select {
			case out <- v:
				return true
			case <-ctx.Done():
				return false
			}
		})
	}()
	return out
}

// Count emits the integers start, start+1, ..., end-1 in order (a half-open
// range [start, end)). If start >= end the stream is empty. The producer
// goroutine exits as soon as ctx is cancelled, even mid-range.
func Count(ctx context.Context, start, end int) <-chan int {
	return Generate(ctx, func(yield func(int) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	})
}

// FromSlice emits each element of items in order. The slice is not retained
// beyond the stream's lifetime. The producer exits promptly on ctx cancellation.
func FromSlice[T any](ctx context.Context, items []T) <-chan T {
	return Generate(ctx, func(yield func(T) bool) {
		for _, v := range items {
			if !yield(v) {
				return
			}
		}
	})
}

// Map returns a new stream that applies fn to every value received from in.
// The returned channel closes when in is drained or ctx is cancelled. Map does
// not drain in on cancellation; cancel the context of the upstream generator so
// it too can exit.
func Map[T, U any](ctx context.Context, in <-chan T, fn func(T) U) <-chan U {
	return Generate(ctx, func(yield func(U) bool) {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				if !yield(fn(v)) {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	})
}

// Take returns a stream that emits at most the first n values from in, then
// closes. If n <= 0 the result is empty. Take stops receiving from in once n
// values have been forwarded; cancel the upstream context to let the source
// generator exit cleanly.
func Take[T any](ctx context.Context, in <-chan T, n int) <-chan T {
	return Generate(ctx, func(yield func(T) bool) {
		for i := 0; i < n; i++ {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				if !yield(v) {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	})
}
