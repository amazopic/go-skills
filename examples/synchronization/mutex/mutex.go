// Package mutex demonstrates idiomatic use of sync.Mutex and sync.RWMutex
// for protecting shared mutable state in Go.
package mutex

import "sync"

// --------------------------------------------------------------------------
// SafeMap: a goroutine-safe string→int map using sync.Mutex
// --------------------------------------------------------------------------

// SafeMap wraps a map with a mutex. All exported methods are goroutine-safe.
// The zero value is not ready to use; use NewSafeMap.
type SafeMap struct {
	mu sync.Mutex
	m  map[string]int
}

// NewSafeMap returns an initialised SafeMap.
func NewSafeMap() *SafeMap {
	return &SafeMap{m: make(map[string]int)}
}

// Set stores key → value.
func (s *SafeMap) Set(key string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

// Get returns the value for key and whether it was found.
func (s *SafeMap) Get(key string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.m[key]
	return v, ok
}

// Inc atomically increments the value associated with key.
func (s *SafeMap) Inc(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key]++
}

// Len returns the number of entries.
func (s *SafeMap) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.m)
}

// --------------------------------------------------------------------------
// ReadHeavyCache: goroutine-safe cache using sync.RWMutex
// Optimised for far more reads than writes.
// --------------------------------------------------------------------------

// ReadHeavyCache stores string→string entries with an RWMutex so many
// goroutines can read concurrently while writes remain exclusive.
type ReadHeavyCache struct {
	mu sync.RWMutex
	m  map[string]string
}

// NewReadHeavyCache returns an initialised ReadHeavyCache.
func NewReadHeavyCache() *ReadHeavyCache {
	return &ReadHeavyCache{m: make(map[string]string)}
}

// Put stores key → value (exclusive write lock).
func (c *ReadHeavyCache) Put(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = value
}

// Lookup returns (value, true) or ("", false). Many goroutines may call
// Lookup concurrently without blocking each other.
func (c *ReadHeavyCache) Lookup(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[key]
	return v, ok
}
