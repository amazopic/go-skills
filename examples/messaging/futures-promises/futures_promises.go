// Package futures implements a generic Future[T] / Promise[T] pair.
//
// A Promise is the write side: Resolve or Reject is called exactly once.
// A Future is the read side: Get blocks until a result is available or ctx expires.
// Multiple callers may share a Future; each call to Get retrieves the value via
// caching so all callers see the same result without draining the channel.
package futures

import (
	"context"
	"sync"
)

// result holds either a value or an error.
type result[T any] struct {
	val T
	err error
}

// Future is the consumer handle for an async computation.
// Safe to share across goroutines; Get may be called any number of times.
type Future[T any] struct {
	once  sync.Once
	ch    <-chan result[T]
	cache result[T]
}

// Get blocks until the promise is fulfilled or ctx is cancelled.
// If ctx expires before the promise is fulfilled, (zero, ctx.Err()) is returned.
// Subsequent calls always return the same result as the first successful Get.
func (f *Future[T]) Get(ctx context.Context) (T, error) {
	// Check if already resolved (fast path after the first Get).
	done := make(chan struct{})
	go func() {
		// Exactly one goroutine wins the Once and populates the cache.
		f.once.Do(func() {
			select {
			case r := <-f.ch:
				f.cache = r
			case <-ctx.Done():
				f.cache = result[T]{err: ctx.Err()}
			}
		})
		close(done)
	}()
	<-done
	return f.cache.val, f.cache.err
}

// Promise is the producer handle. Resolve or Reject must be called exactly once.
type Promise[T any] struct {
	ch   chan result[T]
	once sync.Once
}

// Resolve fulfils the future with a value.
func (p *Promise[T]) Resolve(v T) {
	p.once.Do(func() { p.ch <- result[T]{val: v} })
}

// Reject fulfils the future with an error.
func (p *Promise[T]) Reject(err error) {
	p.once.Do(func() { p.ch <- result[T]{err: err} })
}

// New returns a linked Promise and Future.
// The channel is buffered so Resolve/Reject never block the producer.
func New[T any]() (*Promise[T], *Future[T]) {
	ch := make(chan result[T], 1)
	promise := &Promise[T]{ch: ch}
	future := &Future[T]{ch: ch}
	return promise, future
}

// Go launches fn in a goroutine and returns a Future for its result.
// ctx is forwarded to fn; if ctx is cancelled, Get returns ctx.Err() even if
// fn eventually produces a value.
func Go[T any](ctx context.Context, fn func(context.Context) (T, error)) *Future[T] {
	promise, future := New[T]()
	go func() {
		v, err := fn(ctx)
		if err != nil {
			promise.Reject(err)
		} else {
			promise.Resolve(v)
		}
	}()
	return future
}
