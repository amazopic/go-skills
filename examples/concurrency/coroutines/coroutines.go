// Package coroutines demonstrates cooperative, interleaved execution using
// goroutines and channels — the idiomatic Go equivalent of coroutines.
//
// Coroutine[T] is a generic generator: the body calls yield(v) to suspend and
// deliver a value; the caller calls Next() to advance. Exactly one side runs
// at a time — true cooperative scheduling.
package coroutines

import (
	"context"
	"sync"
)

// Coroutine runs a body function that can yield values of type T one at a
// time. The zero value is not usable; use Start.
type Coroutine[T any] struct {
	yield  chan T
	resume chan struct{}
	done   chan struct{} // closed by Stop; signals the body to stop
	once   sync.Once    // guards close(done)
}

// Start launches the coroutine body. body receives a yield function; calling
// yield(v) suspends body and delivers v to the next call of Next. yield returns
// true to continue or false when the coroutine has been stopped (ctx cancelled
// or Stop called). The body should return when yield returns false.
//
// When body returns, the coroutine is exhausted and subsequent Next calls
// return (zero, false).
func Start[T any](ctx context.Context, body func(yield func(T) bool)) *Coroutine[T] {
	c := &Coroutine[T]{
		yield:  make(chan T),
		resume: make(chan struct{}),
		done:   make(chan struct{}),
	}
	go func() {
		defer close(c.yield) // signal exhaustion to Next
		body(func(v T) bool {
			// Suspend: deliver value to caller or abort.
			select {
			case c.yield <- v:
			case <-ctx.Done():
				return false
			case <-c.done:
				return false
			}
			// Wait for Next() to call us back.
			select {
			case <-c.resume:
				return true
			case <-ctx.Done():
				return false
			case <-c.done:
				return false
			}
		})
	}()
	return c
}

// Next advances the coroutine one step.
// Returns (value, true) when a value is available, or (zero, false) when the
// coroutine is exhausted or has been stopped.
func (c *Coroutine[T]) Next() (T, bool) {
	select {
	case v, ok := <-c.yield:
		if ok {
			// Wake the body. If done is closed, the body will see it and exit.
			select {
			case c.resume <- struct{}{}:
			case <-c.done:
			}
		}
		return v, ok
	case <-c.done:
		// Drain the yield channel in case the body sent just before stop.
		select {
		case v, ok := <-c.yield:
			if ok {
				select {
				case c.resume <- struct{}{}:
				case <-c.done:
				}
			}
			return v, ok
		default:
			var zero T
			return zero, false
		}
	}
}

// Stop signals the coroutine to cease yielding. It is idempotent and safe to
// call from any goroutine. After Stop, subsequent Next calls drain any
// in-flight value and then return (zero, false).
func (c *Coroutine[T]) Stop() {
	c.once.Do(func() { close(c.done) })
}

// --------------------------------------------------------------------------
// Fibonacci generator — example usage of Coroutine[int]
// --------------------------------------------------------------------------

// Fibonacci returns a coroutine that yields the infinite Fibonacci sequence.
// The body checks the return value of yield and exits when it returns false.
func Fibonacci(ctx context.Context) *Coroutine[int] {
	return Start[int](ctx, func(yield func(int) bool) {
		a, b := 0, 1
		for {
			if !yield(a) {
				return
			}
			a, b = b, a+b
		}
	})
}
