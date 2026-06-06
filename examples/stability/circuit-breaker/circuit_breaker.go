// Package circuitbreaker implements the Circuit Breaker stability pattern.
//
// A circuit breaker wraps a call to an unreliable dependency (an external API,
// a database, a downstream service). It tracks consecutive failures and, once a
// threshold is crossed, "opens" the circuit: subsequent calls fail fast with
// [ErrOpenCircuit] instead of piling onto a struggling dependency. After a
// cooldown the breaker moves to half-open and lets a single trial call through;
// success closes the circuit, failure re-opens it with an exponentially longer
// cooldown.
//
// The breaker is safe for concurrent use by multiple goroutines.
package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrOpenCircuit is returned by [Breaker.Do] when the circuit is open and the
// call is rejected without invoking the wrapped function. Callers can branch on
// it with errors.Is.
var ErrOpenCircuit = errors.New("circuitbreaker: circuit is open")

// State is the breaker's lifecycle state.
type State int

const (
	// StateClosed lets all calls through; failures are counted.
	StateClosed State = iota
	// StateOpen rejects all calls until the cooldown elapses.
	StateOpen
	// StateHalfOpen allows a single trial call to probe the dependency.
	StateHalfOpen
)

// String implements fmt.Stringer.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Operation is the unit of work guarded by a breaker: a context-aware call that
// returns an error. The type parameter R is the result the operation produces.
type Operation[R any] func(context.Context) (R, error)

// Breaker guards an Operation[R] against a failing dependency. The zero value
// is not usable; construct one with [New].
type Breaker[R any] struct {
	failureThreshold uint32        // consecutive failures that trip the circuit
	baseCooldown     time.Duration // cooldown after the first trip
	now              func() time.Time

	mu             sync.Mutex
	state          State
	consecFailures uint32
	openCount      uint32    // number of times we have entered open since closed
	retryAt        time.Time // when a half-open trial is permitted
}

// Option configures a Breaker.
type Option[R any] func(*Breaker[R])

// WithClock overrides the time source. It exists for deterministic tests; in
// production the default (time.Now) is correct.
func WithClock[R any](now func() time.Time) Option[R] {
	return func(b *Breaker[R]) { b.now = now }
}

// New builds a Breaker that trips after failureThreshold consecutive failures
// and stays open for baseCooldown (doubling on each successive trip until the
// dependency recovers).
//
// failureThreshold must be > 0 and baseCooldown must be > 0; otherwise New
// panics, as both signal a programming error at wiring time.
func New[R any](failureThreshold uint32, baseCooldown time.Duration, opts ...Option[R]) *Breaker[R] {
	if failureThreshold == 0 {
		panic("circuitbreaker: failureThreshold must be > 0")
	}
	if baseCooldown <= 0 {
		panic("circuitbreaker: baseCooldown must be > 0")
	}
	b := &Breaker[R]{
		failureThreshold: failureThreshold,
		baseCooldown:     baseCooldown,
		now:              time.Now,
		state:            StateClosed,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Do executes op under the breaker's protection.
//
//   - Closed: op runs; a failure increments the counter and may trip the circuit.
//   - Open: op is not run; Do returns ErrOpenCircuit until the cooldown elapses.
//   - Half-open: a single trial op runs; success closes the circuit, failure
//     re-opens it with a longer cooldown.
//
// If ctx is already cancelled Do returns ctx.Err() without touching the breaker
// state. When op returns an error it is returned unchanged (so callers can
// inspect it with errors.Is/As); only the breaker's own rejection surfaces as
// ErrOpenCircuit.
func (b *Breaker[R]) Do(ctx context.Context, op Operation[R]) (R, error) {
	var zero R
	if err := ctx.Err(); err != nil {
		return zero, err
	}

	if !b.allow() {
		return zero, ErrOpenCircuit
	}

	res, err := op(ctx)
	if err != nil {
		b.onResult(false)
		return zero, err
	}
	b.onResult(true)
	return res, nil
}

// allow reports whether a call may proceed, transitioning Open -> HalfOpen when
// the cooldown has elapsed. It holds the lock for the whole decision so state
// transitions are atomic with respect to concurrent callers.
func (b *Breaker[R]) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if b.now().Before(b.retryAt) {
			return false // still cooling down: fail fast
		}
		// Cooldown elapsed: promote to half-open and admit one trial call.
		b.state = StateHalfOpen
		return true
	case StateHalfOpen:
		// A trial is already in flight; reject everyone else until it resolves.
		return false
	default:
		return false
	}
}

// onResult records the outcome of a completed call and advances the state.
func (b *Breaker[R]) onResult(success bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if success {
		// Any success closes the circuit and clears the failure history.
		b.state = StateClosed
		b.consecFailures = 0
		b.openCount = 0
		b.retryAt = time.Time{}
		return
	}

	b.consecFailures++

	switch b.state {
	case StateHalfOpen:
		// Trial failed: re-open immediately with a longer cooldown.
		b.trip()
	case StateClosed:
		if b.consecFailures >= b.failureThreshold {
			b.trip()
		}
	case StateOpen:
		// Should not happen (open rejects calls), but stay safe and idempotent.
	}
}

// trip opens the circuit and schedules the next retry using exponential backoff
// on baseCooldown. Caller must hold b.mu.
func (b *Breaker[R]) trip() {
	b.state = StateOpen
	// backoff: baseCooldown * 2^openCount, capped to avoid overflow / absurd waits.
	shift := b.openCount
	if shift > 16 {
		shift = 16
	}
	b.retryAt = b.now().Add(b.baseCooldown << shift)
	b.openCount++
}

// State returns the current breaker state. Useful for metrics and tests.
func (b *Breaker[R]) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// ConsecutiveFailures returns the current consecutive-failure count.
func (b *Breaker[R]) ConsecutiveFailures() uint32 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.consecFailures
}
