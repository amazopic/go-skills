---
name: synchronization-semaphore
description: Semaphore — limit concurrent access to N resource slots via a buffered channel as a counting semaphore. Use to throttle parallelism without spawning a worker pool.
category: synchronization
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/synchronization/semaphore.md
example: examples/synchronization/semaphore/
---

# Semaphore Pattern
A semaphore is a synchronization pattern/primitive that imposes mutual exclusion on a limited number of resources. 

## Implementation

```go
package semaphore

import (
	"context"
	"errors"
)

// ErrNoTickets is returned by TryAcquire when no slot is free.
var ErrNoTickets = errors.New("semaphore: no tickets available")

// Semaphore is a counting semaphore backed by a buffered channel. The buffer
// capacity is the number of tickets (concurrent holders) allowed. The zero
// value is not usable; create one with New.
type Semaphore struct {
	tickets chan struct{}
}

// New creates a Semaphore that allows n concurrent holders. n must be > 0.
func New(n int) *Semaphore {
	if n <= 0 {
		panic("semaphore: n must be > 0")
	}
	return &Semaphore{tickets: make(chan struct{}, n)}
}

// Acquire blocks until a ticket is free or ctx is done. On success it returns
// nil; the caller must call Release exactly once. On cancellation it returns
// ctx.Err() wrapped, and no ticket is held.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.tickets <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire grabs a ticket without blocking. It returns ErrNoTickets if none
// is free. On success the caller must call Release exactly once.
func (s *Semaphore) TryAcquire() error {
	select {
	case s.tickets <- struct{}{}:
		return nil
	default:
		return ErrNoTickets
	}
}

// Release returns one ticket. It must be called exactly once per successful
// Acquire/TryAcquire; releasing without holding a ticket panics.
func (s *Semaphore) Release() {
	select {
	case <-s.tickets:
	default:
		panic("semaphore: Release called without a held ticket")
	}
}
```

## Usage
### Blocking acquire (with cancellation/deadline via context)

```go
s := semaphore.New(3) // at most 3 concurrent holders

ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
defer cancel()

if err := s.Acquire(ctx); err != nil {
    return err // ctx cancelled or deadline exceeded; no ticket held
}
defer s.Release()

// Do important work — guaranteed at most 3 goroutines reach here at once.
```

### Non-blocking acquire (fail fast)

```go
s := semaphore.New(1)

if err := s.TryAcquire(); err != nil {
    if errors.Is(err, semaphore.ErrNoTickets) {
        return // no slot free right now; skip / shed load
    }
    return err
}
defer s.Release()
```

## Gotchas

- **Always pair Acquire with Release** — use `defer s.Release()` right after a
  successful acquire so an early return or panic still returns the ticket.
- **Never Release without a held ticket.** Doing so would over-fill the channel
  and let more than N holders through; this implementation panics instead.
- **Prefer `context` over timers.** A `time.After`-based timeout leaks the timer
  until it fires and is harder to compose; `Acquire(ctx)` cancels cleanly and
  reports `ctx.Err()`.
