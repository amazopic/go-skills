// Package nbarrier implements a cyclic N-Barrier synchronisation primitive.
// All N goroutines block at Wait until the last one arrives, then all proceed.
// The barrier re-arms automatically, supporting multiple rounds.
package nbarrier

import "sync"

// Barrier is a reusable rendezvous point for exactly n goroutines.
// The zero value is not usable; use NewBarrier.
type Barrier struct {
	n       int
	mu      sync.Mutex
	waiting int
	done    chan struct{} // closed to broadcast "proceed"
}

// NewBarrier creates a Barrier that releases when n goroutines have called Wait.
// Panics if n < 1.
func NewBarrier(n int) *Barrier {
	if n < 1 {
		panic("nbarrier: n must be >= 1")
	}
	return &Barrier{n: n, done: make(chan struct{})}
}

// Wait blocks until n goroutines (including the caller) have called Wait.
// The call is safe to make from multiple goroutines concurrently.
// After Wait returns, the barrier is armed for the next round.
func (b *Barrier) Wait() {
	b.mu.Lock()
	b.waiting++
	if b.waiting == b.n {
		// Last arrival: broadcast by closing the current gate and arm a new one.
		old := b.done
		b.done = make(chan struct{})
		b.waiting = 0
		b.mu.Unlock()
		close(old) // unblocks all waiters simultaneously
		return
	}
	// Not the last — capture the gate under the lock, then block outside it.
	gate := b.done
	b.mu.Unlock()
	<-gate
}
