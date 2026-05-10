// Package steadystate demonstrates the Steady-State stability pattern.
//
// A system reaches steady state when resource accumulation is bounded: caches
// are evicted, logs are rotated, stale entries are compacted. Without active
// bounding, memory and disk grow without limit under sustained load, eventually
// causing OOM kills or disk exhaustion. The pattern is implemented as a janitor
// goroutine that runs a bounded data structure's eviction on a regular schedule.
package steadystate

import (
	"context"
	"sync"
	"time"
)

// entry is a cache item with an expiry timestamp.
type entry struct {
	value     any
	expiresAt time.Time
}

// Cache is a TTL-bounded, capacity-bounded in-memory cache. It never grows
// beyond maxSize entries. A background janitor goroutine evicts expired entries
// on a configurable interval, preventing unbounded accumulation.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]entry
	maxSize int

	// janitor coordination
	janitorInterval time.Duration
	stop            chan struct{}
	stopped         chan struct{}
}

// NewCache creates a Cache with the given capacity cap and eviction interval.
// Call Start to launch the background janitor. Call Stop to shut it down.
func NewCache(maxSize int, janitorInterval time.Duration) *Cache {
	if maxSize <= 0 {
		panic("steadystate: maxSize must be > 0")
	}
	if janitorInterval <= 0 {
		panic("steadystate: janitorInterval must be > 0")
	}
	return &Cache{
		items:           make(map[string]entry),
		maxSize:         maxSize,
		janitorInterval: janitorInterval,
		stop:            make(chan struct{}),
		stopped:         make(chan struct{}),
	}
}

// Start launches the background janitor goroutine. The goroutine exits when
// ctx is cancelled or Stop is called, whichever comes first.
func (c *Cache) Start(ctx context.Context) {
	go func() {
		defer close(c.stopped)
		ticker := time.NewTicker(c.janitorInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.evictExpired()
			case <-ctx.Done():
				return
			case <-c.stop:
				return
			}
		}
	}()
}

// Stop signals the janitor to exit and waits for it to finish.
func (c *Cache) Stop() {
	close(c.stop)
	<-c.stopped
}

// Set inserts or updates a key with the given TTL. If the cache is full and
// the new key is not already present, Set evicts the oldest expired entry
// before inserting. If there are no expired entries to evict, the insert is
// dropped (capacity-capped write shedding).
func (c *Cache) Set(key string, value any, ttl time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update existing — always allowed.
	if _, exists := c.items[key]; exists {
		c.items[key] = entry{value: value, expiresAt: time.Now().Add(ttl)}
		return true
	}

	// At capacity — try to make room by evicting one expired entry.
	if len(c.items) >= c.maxSize {
		now := time.Now()
		for k, e := range c.items {
			if now.After(e.expiresAt) {
				delete(c.items, k)
				break
			}
		}
		// Still full after eviction attempt — shed the write.
		if len(c.items) >= c.maxSize {
			return false
		}
	}

	c.items[key] = entry{value: value, expiresAt: time.Now().Add(ttl)}
	return true
}

// Get returns the value for key if it exists and has not expired.
func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

// Len returns the number of entries currently in the cache (including expired
// entries not yet swept by the janitor).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// evictExpired removes all entries whose TTL has elapsed. Called by janitor.
func (c *Cache) evictExpired() {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, e := range c.items {
		if now.After(e.expiresAt) {
			delete(c.items, k)
		}
	}
}

// EvictExpired is the exported version for testing without waiting for the timer.
func (c *Cache) EvictExpired() {
	c.evictExpired()
}

// RingBuffer is a fixed-size FIFO buffer that overwrites the oldest entry on
// overflow, bounding memory to exactly cap elements — the archetypal
// steady-state data structure for log rings, metrics windows, etc.
type RingBuffer[T any] struct {
	mu   sync.Mutex
	buf  []T
	head int // write position
	size int // number of valid entries
	cap  int
}

// NewRingBuffer creates a RingBuffer with the given fixed capacity.
func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	if capacity <= 0 {
		panic("steadystate: ring buffer capacity must be > 0")
	}
	return &RingBuffer[T]{buf: make([]T, capacity), cap: capacity}
}

// Push appends value, overwriting the oldest entry if the buffer is full.
func (r *RingBuffer[T]) Push(v T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buf[r.head] = v
	r.head = (r.head + 1) % r.cap
	if r.size < r.cap {
		r.size++
	}
}

// Slice returns a copy of all valid entries in oldest-first order.
func (r *RingBuffer[T]) Slice() []T {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]T, r.size)
	start := (r.head - r.size + r.cap) % r.cap
	for i := range out {
		out[i] = r.buf[(start+i)%r.cap]
	}
	return out
}

// Len returns the number of valid entries.
func (r *RingBuffer[T]) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.size
}
