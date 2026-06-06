// Package cascading illustrates the Cascading Failures anti-pattern and the
// mitigations that prevent downstream failure from spreading upstream.
//
// Scenario: three services — Frontend → OrderService → PaymentService.
// PaymentService becomes slow or down. Without mitigations, the slowness
// propagates: threads/goroutines pile up, memory grows, the whole stack falls.
//
// This file shows both the broken pattern (NaiveClient) and the mitigated
// pattern (SafeClient) using timeouts, bulkheads, and a circuit-breaker token.
package cascading

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// --- Broken pattern ---------------------------------------------------------

// ErrTimeout is returned when an operation exceeds its deadline.
var ErrTimeout = errors.New("cascading: timeout")

// ErrCircuitOpen is returned when the circuit breaker refuses the call.
var ErrCircuitOpen = errors.New("cascading: circuit open")

// ErrPoolExhausted is returned when the bulkhead pool is full.
var ErrPoolExhausted = errors.New("cascading: pool exhausted")

// SlowDep simulates a dependency that hangs for hang duration before
// returning an error. It counts how many calls are currently in-flight.
type SlowDep struct {
	hang     time.Duration
	inFlight atomic.Int64
}

// Call blocks for hang, then returns an error. Callers without a deadline will
// block for the full duration, piling up goroutines.
func (d *SlowDep) Call(ctx context.Context) error {
	d.inFlight.Add(1)
	defer d.inFlight.Add(-1)

	select {
	case <-time.After(d.hang):
		return errors.New("dep: unavailable")
	case <-ctx.Done():
		return ctx.Err()
	}
}

// InFlight returns the number of goroutines currently blocked inside Call.
func (d *SlowDep) InFlight() int64 { return d.inFlight.Load() }

// NaiveClient calls the dependency with no timeout and no concurrency cap.
// Under load, every goroutine blocks for dep.hang, memory grows linearly, and
// the caller's goroutine pool exhausts — a textbook cascading failure.
type NaiveClient struct {
	dep *SlowDep
}

// NewNaiveClient creates a client with no protections.
func NewNaiveClient(dep *SlowDep) *NaiveClient { return &NaiveClient{dep: dep} }

// Call passes the context unmodified — if ctx has no deadline, all
// goroutines pile up for the full dep.hang duration.
func (c *NaiveClient) Call(ctx context.Context) error {
	return c.dep.Call(ctx)
}

// --- Mitigation 1: per-call timeout ----------------------------------------

// TimeoutClient wraps calls with a per-call timeout, bounding how long a
// single goroutine can block waiting for a slow dependency.
type TimeoutClient struct {
	dep     *SlowDep
	timeout time.Duration
}

// NewTimeoutClient creates a client that caps each call to timeout.
func NewTimeoutClient(dep *SlowDep, timeout time.Duration) *TimeoutClient {
	return &TimeoutClient{dep: dep, timeout: timeout}
}

// Call runs dep.Call with a deadline bounded to timeout or the parent
// context's remaining budget, whichever is shorter.
func (c *TimeoutClient) Call(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	err := c.dep.Call(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	return err
}

// --- Mitigation 2: bulkhead (concurrency cap) -------------------------------

// BulkheadClient wraps a TimeoutClient with a semaphore that caps concurrent
// calls. Excess callers are shed immediately with ErrPoolExhausted rather than
// blocking and accumulating.
type BulkheadClient struct {
	inner *TimeoutClient
	sem   chan struct{}
}

// NewBulkheadClient creates a client with both timeout and concurrency cap.
func NewBulkheadClient(dep *SlowDep, timeout time.Duration, maxConcurrent int) *BulkheadClient {
	return &BulkheadClient{
		inner: NewTimeoutClient(dep, timeout),
		sem:   make(chan struct{}, maxConcurrent),
	}
}

// Call acquires a slot (non-blocking) then delegates to the timeout client.
func (c *BulkheadClient) Call(ctx context.Context) error {
	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	default:
		return ErrPoolExhausted
	}
	return c.inner.Call(ctx)
}

// --- Mitigation 3: circuit breaker token ------------------------------------

// circuitState tracks the circuit breaker state.
type circuitState int32

const (
	stateClosed   circuitState = iota // normal operation
	stateOpen                         // fast-fail
	stateHalfOpen                     // probing
)

// CircuitBreaker is a simple three-state circuit breaker. In production use
// sony/gobreaker or a similar battle-tested library.
type CircuitBreaker struct {
	mu              sync.Mutex
	state           circuitState
	failures        int
	threshold       int
	resetAfter      time.Duration
	lastFailureTime time.Time
}

// NewCircuitBreaker creates a breaker that opens after threshold consecutive
// failures and probes again after resetAfter.
func NewCircuitBreaker(threshold int, resetAfter time.Duration) *CircuitBreaker {
	return &CircuitBreaker{threshold: threshold, resetAfter: resetAfter}
}

// Allow reports whether the caller may proceed. It transitions states
// as needed and returns ErrCircuitOpen when the circuit is open.
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case stateClosed:
		return nil
	case stateOpen:
		if time.Since(cb.lastFailureTime) >= cb.resetAfter {
			cb.state = stateHalfOpen
			return nil
		}
		return ErrCircuitOpen
	case stateHalfOpen:
		return nil
	}
	return nil
}

// Record records the outcome of a call and updates state.
func (cb *CircuitBreaker) Record(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if err != nil && !errors.Is(err, context.Canceled) {
		cb.failures++
		cb.lastFailureTime = time.Now()
		if cb.failures >= cb.threshold {
			cb.state = stateOpen
		}
	} else {
		cb.failures = 0
		cb.state = stateClosed
	}
}

// State returns the current circuit state string for diagnostics.
func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case stateClosed:
		return "closed"
	case stateOpen:
		return "open"
	case stateHalfOpen:
		return "half-open"
	}
	return "unknown"
}

// FullyProtectedClient combines timeout + bulkhead + circuit breaker.
// This is the production-safe composition that prevents cascading failures.
type FullyProtectedClient struct {
	dep *SlowDep
	cb  *CircuitBreaker
	bh  *BulkheadClient
}

// NewFullyProtectedClient wires all three mitigations together.
func NewFullyProtectedClient(dep *SlowDep, timeout time.Duration, maxConcurrent, cbThreshold int, cbReset time.Duration) *FullyProtectedClient {
	return &FullyProtectedClient{
		dep: dep,
		cb:  NewCircuitBreaker(cbThreshold, cbReset),
		bh:  NewBulkheadClient(dep, timeout, maxConcurrent),
	}
}

// Call enforces: circuit open? fail fast. Pool full? shed. Timeout? return.
// All three layers compose cleanly because each is a simple function wrapper.
func (c *FullyProtectedClient) Call(ctx context.Context) error {
	if err := c.cb.Allow(); err != nil {
		return err
	}
	err := c.bh.Call(ctx)
	c.cb.Record(err)
	return err
}

// CBState returns the circuit breaker state for tests/diagnostics.
func (c *FullyProtectedClient) CBState() string { return c.cb.State() }
