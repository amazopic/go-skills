// Package registry is an example of the Registry Pattern.
// A Registry is a typed, goroutine-safe lookup table for named values.
// It is commonly used to register plugins, codecs, drivers, or shared
// singletons by name without creating a direct import dependency between
// producer and consumer.
package registry

import (
	"fmt"
	"sync"
)

// Registry is a generic, goroutine-safe map from string keys to values of type T.
// The zero value is ready to use.
type Registry[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

// Register adds name → value to the registry.
// Returns an error if name is already registered.
func (r *Registry[T]) Register(name string, value T) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.items == nil {
		r.items = make(map[string]T)
	}
	if _, exists := r.items[name]; exists {
		return fmt.Errorf("registry: %q is already registered", name)
	}
	r.items[name] = value
	return nil
}

// MustRegister is like Register but panics on duplicate names.
// Intended for use in package-level init functions.
func (r *Registry[T]) MustRegister(name string, value T) {
	if err := r.Register(name, value); err != nil {
		panic(err)
	}
}

// Get retrieves the value for name.
// The second return value is false when name has not been registered.
func (r *Registry[T]) Get(name string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.items[name]
	return v, ok
}

// MustGet is like Get but panics when name is not found.
func (r *Registry[T]) MustGet(name string) T {
	v, ok := r.Get(name)
	if !ok {
		panic(fmt.Sprintf("registry: %q not found", name))
	}
	return v
}

// Names returns a snapshot of all registered names in unspecified order.
func (r *Registry[T]) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.items))
	for k := range r.items {
		names = append(names, k)
	}
	return names
}

// Len returns the number of registered entries.
func (r *Registry[T]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.items)
}
