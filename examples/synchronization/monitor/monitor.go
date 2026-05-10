// Package monitor demonstrates the Monitor pattern: a type that bundles shared
// state with its synchronisation primitives and exposes only safe public methods.
//
// ResourcePool is a fixed-capacity pool of integer tokens. Callers Acquire a
// token (blocking if none are available) and Release it when done.
// Callers never touch the mutex or condition variable directly.
package monitor

import "sync"

// ResourcePool is a monitor: unexported fields, exported methods.
// Acquire blocks when no token is available; Release returns one.
type ResourcePool struct {
	mu      sync.Mutex
	cond    *sync.Cond
	tokens  []int
	maxSize int
}

// NewResourcePool creates a pool pre-loaded with tokens 0..size-1.
// Panics if size < 1.
func NewResourcePool(size int) *ResourcePool {
	if size < 1 {
		panic("monitor: pool size must be >= 1")
	}
	p := &ResourcePool{maxSize: size}
	p.cond = sync.NewCond(&p.mu)
	p.tokens = make([]int, size)
	for i := range size {
		p.tokens[i] = i
	}
	return p
}

// Acquire removes and returns a token, blocking until one is available.
func (p *ResourcePool) Acquire() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for len(p.tokens) == 0 { // always loop — Mesa semantics
		p.cond.Wait()
	}
	token := p.tokens[len(p.tokens)-1]
	p.tokens = p.tokens[:len(p.tokens)-1]
	return token
}

// Release returns a token to the pool and wakes one waiter.
// Panics if more tokens are released than the pool was initialised with.
func (p *ResourcePool) Release(token int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.tokens) == p.maxSize {
		panic("monitor: pool overflow — Release called more times than Acquire")
	}
	p.tokens = append(p.tokens, token)
	p.cond.Signal()
}

// Available returns the number of tokens currently in the pool.
func (p *ResourcePool) Available() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.tokens)
}

// --------------------------------------------------------------------------
// AtomicRegistry: a monitor demonstrating Get-or-set compound operation.
// --------------------------------------------------------------------------

// AtomicRegistry is a goroutine-safe map that supports atomic Get-or-set:
// if the key is absent, it is set and the new value returned.
type AtomicRegistry struct {
	mu sync.RWMutex
	m  map[string]string
}

// NewAtomicRegistry returns an initialised AtomicRegistry.
func NewAtomicRegistry() *AtomicRegistry {
	return &AtomicRegistry{m: make(map[string]string)}
}

// GetOrSet returns the existing value for key, or stores and returns val if absent.
// The check-and-set is atomic: no other goroutine can interleave.
func (r *AtomicRegistry) GetOrSet(key, val string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.m[key]; ok {
		return existing, false // already existed
	}
	r.m[key] = val
	return val, true // newly set
}

// Get returns (value, true) or ("", false).
func (r *AtomicRegistry) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.m[key]
	return v, ok
}
