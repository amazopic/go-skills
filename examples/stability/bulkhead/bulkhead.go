// Package bulkhead demonstrates the Bulkhead stability pattern.
//
// A bulkhead isolates resource pools so that exhaustion or failure in one
// consumer (tenant, call type, service) cannot drain shared capacity and
// bring down unrelated consumers. Each pool owns its own semaphore; acquiring
// one pool's slot never blocks another pool.
package bulkhead

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// ErrPoolExhausted is returned when all slots in a pool are occupied.
var ErrPoolExhausted = errors.New("bulkhead: pool exhausted")

// Pool is a concurrency-limited execution compartment. Callers acquire a slot
// before doing work and release it afterwards. The pool is safe for concurrent
// use by multiple goroutines.
type Pool struct {
	name string
	sem  chan struct{}

	// metrics — updated atomically so tests can inspect without a lock.
	acquired atomic.Int64
	rejected atomic.Int64
	current  atomic.Int64
}

// NewPool creates a Pool with the given name and maximum concurrency limit.
// limit must be > 0.
func NewPool(name string, limit int) *Pool {
	if limit <= 0 {
		panic("bulkhead: limit must be > 0")
	}
	return &Pool{
		name: name,
		sem:  make(chan struct{}, limit),
	}
}

// Do runs fn inside the pool. It returns ErrPoolExhausted immediately if no
// slot is available. Context cancellation is respected while waiting is not
// attempted here — this is a non-blocking bulkhead (fail-fast on overflow).
func (p *Pool) Do(ctx context.Context, fn func(context.Context) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	select {
	case p.sem <- struct{}{}:
		// slot acquired
	default:
		p.rejected.Add(1)
		return ErrPoolExhausted
	}

	p.acquired.Add(1)
	p.current.Add(1)
	defer func() {
		<-p.sem
		p.current.Add(-1)
	}()

	return fn(ctx)
}

// Stats returns a snapshot of (acquired, rejected, current) counters.
func (p *Pool) Stats() (acquired, rejected, current int64) {
	return p.acquired.Load(), p.rejected.Load(), p.current.Load()
}

// Bulkhead groups named pools. Callers reference pools by tenant or call type.
type Bulkhead struct {
	mu    sync.RWMutex
	pools map[string]*Pool
}

// New creates an empty Bulkhead.
func New() *Bulkhead {
	return &Bulkhead{pools: make(map[string]*Pool)}
}

// AddPool registers a named pool with the given concurrency limit.
// Panics if the name is already registered.
func (b *Bulkhead) AddPool(name string, limit int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.pools[name]; exists {
		panic("bulkhead: duplicate pool name: " + name)
	}
	b.pools[name] = NewPool(name, limit)
}

// Do executes fn in the named pool. Returns ErrPoolExhausted if the pool is
// full, or an error if the pool name is unknown.
func (b *Bulkhead) Do(ctx context.Context, pool string, fn func(context.Context) error) error {
	b.mu.RLock()
	p, ok := b.pools[pool]
	b.mu.RUnlock()
	if !ok {
		return errors.New("bulkhead: unknown pool: " + pool)
	}
	return p.Do(ctx, fn)
}

// Pool returns the named Pool or nil.
func (b *Bulkhead) Pool(name string) *Pool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.pools[name]
}
