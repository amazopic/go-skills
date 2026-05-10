// Package rwlock demonstrates sync.RWMutex: concurrent reads, exclusive writes.
// ConfigStore caches string key-value pairs. Many goroutines may Get concurrently;
// Set is exclusive.
package rwlock

import "sync"

// ConfigStore is a goroutine-safe key-value store optimised for reads.
// Any number of goroutines may call Get concurrently without blocking each
// other. Set acquires an exclusive lock, blocking all readers and writers.
type ConfigStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewConfigStore returns an empty ConfigStore.
func NewConfigStore() *ConfigStore {
	return &ConfigStore{data: make(map[string]string)}
}

// Get returns the value for key and whether it was found.
// Many goroutines may call Get concurrently.
func (c *ConfigStore) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

// Set stores key → value exclusively.
func (c *ConfigStore) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// Delete removes a key.
func (c *ConfigStore) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Snapshot returns a copy of all current entries.
// Acquires a read lock for the duration of the copy.
func (c *ConfigStore) Snapshot() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]string, len(c.data))
	for k, v := range c.data {
		out[k] = v
	}
	return out
}

// Len returns the number of entries.
func (c *ConfigStore) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}
