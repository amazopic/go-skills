// Package objectpool implements the Object Pool creational pattern: a bounded,
// concurrency-safe set of pre-allocated, reusable instances.
//
// Use it when constructing an object is markedly more expensive than keeping it
// alive between uses (database connections, large scratch buffers, parsers).
// Callers Get an object, use it, then Put it back. The pool caps the number of
// live objects, providing back-pressure instead of unbounded allocation.
//
// For pure zero-allocation, GC-managed reuse where a strict cap is NOT required,
// prefer the standard library's sync.Pool. Reach for this Pool when you need a
// hard upper bound on live resources and context-aware blocking acquisition.
package objectpool

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Sentinel errors returned by Pool, suitable for errors.Is checks.
var (
	// ErrPoolClosed is returned by Get and Put after Close has been called.
	ErrPoolClosed = errors.New("objectpool: pool is closed")

	// ErrNoFactory is returned by New when the factory function is nil.
	ErrNoFactory = errors.New("objectpool: factory must not be nil")

	// ErrInvalidSize is returned by New when size < 1.
	ErrInvalidSize = errors.New("objectpool: size must be >= 1")
)

// Pool is a bounded pool of reusable values of type T.
//
// The zero value is not usable; construct a Pool with New. A Pool is safe for
// concurrent use by multiple goroutines. Objects are created lazily by the
// factory up to the configured size, so a large pool costs nothing until it is
// actually exercised.
type Pool[T any] struct {
	// free carries idle objects. Its capacity is the pool size and also acts as
	// the semaphore that bounds the number of live objects.
	free chan T

	// factory builds a fresh object. It may return an error (e.g. dialing a DB).
	factory func(context.Context) (T, error)

	// reset, if non-nil, is invoked on every Put to scrub per-use state before
	// the object returns to the free list (e.g. buf.Reset()).
	reset func(T)

	mu     sync.Mutex // guards closed and created
	closed bool
	// created counts objects the factory has produced. It never exceeds the
	// pool size, which is how lazy creation stays bounded.
	created int
	size    int
}

// New creates a Pool that holds at most size objects, built on demand by
// factory. The optional reset hook is called on each Put to clear per-use state
// before reuse; pass nil if objects need no scrubbing.
//
// New returns ErrInvalidSize if size < 1 and ErrNoFactory if factory is nil.
func New[T any](size int, factory func(context.Context) (T, error), reset func(T)) (*Pool[T], error) {
	if size < 1 {
		return nil, ErrInvalidSize
	}
	if factory == nil {
		return nil, ErrNoFactory
	}
	return &Pool[T]{
		free:    make(chan T, size),
		factory: factory,
		reset:   reset,
		size:    size,
	}, nil
}

// Get returns an object from the pool, blocking until one is available or ctx is
// done. If an idle object exists it is returned immediately. Otherwise, if the
// pool has not yet created its maximum number of objects, a new one is built via
// the factory. If the pool is at capacity with all objects in use, Get blocks
// until another goroutine Puts one back or ctx is cancelled.
//
// On success the caller owns the returned object and must eventually Put it back
// (typically via defer) to avoid starving the pool.
//
// Get returns ErrPoolClosed if the pool is closed, ctx.Err() (wrapped) if the
// context is done while waiting, or a wrapped factory error if construction
// fails. When Get returns a non-nil error, the zero value of T is returned and
// no pool slot is consumed.
func (p *Pool[T]) Get(ctx context.Context) (T, error) {
	var zero T

	// Fast path: hand back an idle object without touching the lock.
	select {
	case obj, ok := <-p.free:
		if !ok {
			return zero, ErrPoolClosed
		}
		return obj, nil
	default:
	}

	// No idle object. Try to create one if we are below the cap.
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return zero, ErrPoolClosed
	}
	if p.created < p.size {
		p.created++
		p.mu.Unlock()

		obj, err := p.factory(ctx)
		if err != nil {
			// Construction failed: release the reserved slot so the pool does
			// not leak capacity, then surface the wrapped error.
			p.mu.Lock()
			p.created--
			p.mu.Unlock()
			return zero, fmt.Errorf("objectpool: factory failed: %w", err)
		}
		return obj, nil
	}
	p.mu.Unlock()

	// At capacity: block for a returned object or context cancellation.
	select {
	case obj, ok := <-p.free:
		if !ok {
			return zero, ErrPoolClosed
		}
		return obj, nil
	case <-ctx.Done():
		return zero, fmt.Errorf("objectpool: get cancelled: %w", ctx.Err())
	}
}

// Put returns obj to the pool for reuse. If a reset hook was configured it runs
// first to scrub per-use state. Put never blocks: the free channel has exactly
// enough capacity for every object the pool can create, so a correctly used
// pool (one Put per successful Get) always has room.
//
// Put returns ErrPoolClosed if the pool has been closed; in that case obj is
// dropped (and eligible for GC) rather than retained.
func (p *Pool[T]) Put(obj T) error {
	if p.reset != nil {
		p.reset(obj)
	}

	// Hold the lock across the send. Close also closes p.free under this same
	// lock, so checking closed and sending in one critical section closes the
	// TOCTOU window: the channel cannot be closed between the check and the send.
	// A naive "check closed, unlock, then send" leaves a gap where a concurrent
	// Close turns the send into a panic ("send on closed channel").
	//
	// The send is non-blocking (select/default), so holding the lock here cannot
	// deadlock: it returns immediately whether the free list has room or not.
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return ErrPoolClosed
	}

	// Non-blocking by construction: cap(free) == size >= live objects.
	select {
	case p.free <- obj:
		return nil
	default:
		// Defensive: a full channel means Put was called more times than Get
		// (a caller bug). Drop the surplus object rather than block forever.
		return nil
	}
}

// Available reports how many idle objects are currently in the pool and ready to
// be handed out without blocking or invoking the factory.
func (p *Pool[T]) Available() int {
	return len(p.free)
}

// Size returns the maximum number of objects the pool can hold.
func (p *Pool[T]) Size() int {
	return p.size
}

// Close shuts the pool down. After Close, Get and Put return ErrPoolClosed.
// Idle objects are drained from the free list and, if dispose is non-nil, passed
// to it so callers can release underlying resources (close connections, return
// buffers). Objects currently checked out are not affected; their later Put will
// return ErrPoolClosed and drop them.
//
// Close is idempotent: calling it more than once is a no-op and returns nil.
func (p *Pool[T]) Close(dispose func(T)) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	close(p.free)
	p.mu.Unlock()

	// Drain idle objects after closing the channel; the range stops once the
	// buffered objects are consumed.
	for obj := range p.free {
		if dispose != nil {
			dispose(obj)
		}
	}
	return nil
}
