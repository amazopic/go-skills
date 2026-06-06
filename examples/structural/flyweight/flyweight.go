// Package flyweight is an example of the Flyweight Pattern.
package flyweight

import (
	"fmt"
	"sync"
)

// Flyweighter interface
type Flyweighter interface {
	Draw(width, height int, opacity float64) string
}

// FlyweightFactory implements a factory.
// If a suitable flyweighter is in pool, then returns it.
//
// The pool is a shared interning cache, so concurrent GetFlyweight calls must
// not race on it. An RWMutex guards the map: lookups take the cheap read lock,
// and only first-time inserts upgrade to the write lock.
type FlyweightFactory struct {
	mu   sync.RWMutex
	pool map[string]Flyweighter
}

// GetFlyweight creates or returns a suitable Flyweighter by state.
// It is safe for concurrent use and guarantees interning: repeated calls with
// the same filename always return the same *ConcreteFlyweight.
func (f *FlyweightFactory) GetFlyweight(filename string) Flyweighter {
	// Fast path: an existing entry only needs a shared read lock.
	f.mu.RLock()
	fw, ok := f.pool[filename]
	f.mu.RUnlock()
	if ok {
		return fw
	}

	// Slow path: create the entry under the write lock. Re-check after locking
	// (double-checked locking) since another goroutine may have inserted it
	// between releasing the read lock and acquiring the write lock.
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.pool == nil {
		f.pool = make(map[string]Flyweighter)
	}
	if fw, ok = f.pool[filename]; !ok {
		fw = &ConcreteFlyweight{filename: filename}
		f.pool[filename] = fw
	}
	return fw
}

// ConcreteFlyweight implements a Flyweighter interface.
type ConcreteFlyweight struct {
	filename string // internal state
}

// Draw draws image. Args width, height and opacity is external state.
func (f *ConcreteFlyweight) Draw(width, height int, opacity float64) string {
	return fmt.Sprintf("draw image: %s, width: %d, height: %d, opacity: %.2f", f.filename, width, height, opacity)
}
