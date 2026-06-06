// Package semaphore implements a counting semaphore that limits concurrent
// access to N resource slots ("tickets").
//
// The implementation is a buffered channel used as a counting semaphore:
// sending to the channel claims a slot, receiving from it returns one. This
// bounds parallelism without spawning a worker pool — callers run their own
// goroutines and gate the hot section with Acquire/Release.
//
// Blocking acquisition is context-first: Acquire honours cancellation and
// deadlines and never leaks a timer. The zero value of Semaphore is not
// usable; create one with New.
package semaphore

import (
	"context"
	"errors"
	"fmt"
)

// ErrNoTickets is returned by TryAcquire when no slot is free.
var ErrNoTickets = errors.New("semaphore: no tickets available")

// Semaphore limits the number of goroutines that may hold a ticket at once.
// It is safe for concurrent use by multiple goroutines.
type Semaphore struct {
	tickets chan struct{}
}

// New creates a Semaphore that allows at most n concurrent holders.
// It panics if n <= 0.
func New(n int) *Semaphore {
	if n <= 0 {
		panic("semaphore: n must be > 0")
	}
	return &Semaphore{tickets: make(chan struct{}, n)}
}

// Acquire blocks until a ticket is free or ctx is done.
//
// On success it returns nil and the caller must call Release exactly once
// (use defer s.Release()). If ctx is cancelled or its deadline is exceeded
// before a ticket is obtained, Acquire returns the wrapped context error and
// no ticket is held.
func (s *Semaphore) Acquire(ctx context.Context) error {
	// Fast path: respect an already-cancelled context before blocking, and
	// avoid the non-deterministic select choice when a ticket happens to be
	// free at the same time the context is done.
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("semaphore: acquire: %w", err)
	}
	select {
	case s.tickets <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("semaphore: acquire: %w", ctx.Err())
	}
}

// TryAcquire grabs a ticket without blocking. It returns ErrNoTickets if every
// slot is currently held. On success the caller must call Release exactly once.
func (s *Semaphore) TryAcquire() error {
	select {
	case s.tickets <- struct{}{}:
		return nil
	default:
		return ErrNoTickets
	}
}

// Release returns one ticket to the semaphore. It must be called exactly once
// for each successful Acquire or TryAcquire. Calling Release without holding a
// ticket is a programming error and panics, because silently dropping the call
// would let more than N holders through.
func (s *Semaphore) Release() {
	select {
	case <-s.tickets:
	default:
		panic("semaphore: Release called without a held ticket")
	}
}

// Len reports how many tickets are currently held.
func (s *Semaphore) Len() int {
	return len(s.tickets)
}

// Cap reports the maximum number of concurrent holders (the value passed to New).
func (s *Semaphore) Cap() int {
	return cap(s.tickets)
}

// Do runs fn while holding a ticket, acquiring (blocking, context-aware) before
// and releasing after — even if fn panics. If the context is done before a
// ticket can be acquired, fn is not run and the context error is returned.
func (s *Semaphore) Do(ctx context.Context, fn func(context.Context) error) error {
	if err := s.Acquire(ctx); err != nil {
		return err
	}
	defer s.Release()
	return fn(ctx)
}
