---
name: creational-object-pool
description: Object Pool — pre-allocate and reuse expensive resources (DB connections, large buffers) instead of creating per request. Use sync.Pool when zero-allocation matters.
category: creational
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/creational/object-pool.md
example: examples/creational/object-pool/
---

# Object Pool Pattern

The object pool creational design pattern is used to prepare and keep multiple
instances according to the demand expectation.

## Implementation

A production-grade pool is generic, bounded, and context-aware. A buffered
channel doubles as both the free list and the semaphore that caps live objects;
a factory builds them lazily up to the cap. See
`examples/creational/object-pool/` for the full, race-tested version.

```go
package objectpool

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrPoolClosed = errors.New("objectpool: pool is closed")

type Pool[T any] struct {
	free    chan T
	factory func(context.Context) (T, error)

	mu      sync.Mutex
	closed  bool
	created int
	size    int
}

func New[T any](size int, factory func(context.Context) (T, error)) (*Pool[T], error) {
	if size < 1 {
		return nil, errors.New("objectpool: size must be >= 1")
	}
	if factory == nil {
		return nil, errors.New("objectpool: factory must not be nil")
	}
	return &Pool[T]{free: make(chan T, size), factory: factory, size: size}, nil
}

// Get returns a pooled object, creating one lazily if below the cap, or
// blocking until a Put frees one. It honours context cancellation.
func (p *Pool[T]) Get(ctx context.Context) (T, error) {
	var zero T
	select {
	case obj, ok := <-p.free:
		if !ok {
			return zero, ErrPoolClosed
		}
		return obj, nil
	default:
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return zero, ErrPoolClosed
	}
	if p.created < p.size {
		p.created++
		p.mu.Unlock()
		obj, err := p.factory(ctx)
		if err != nil {
			p.mu.Lock()
			p.created-- // release the reserved slot so capacity does not leak
			p.mu.Unlock()
			return zero, fmt.Errorf("objectpool: factory failed: %w", err)
		}
		return obj, nil
	}
	p.mu.Unlock()

	select {
	case obj, ok := <-p.free:
		if !ok {
			return zero, ErrPoolClosed
		}
		return obj, nil
	case <-ctx.Done():
		return zero, fmt.Errorf("objectpool: get cancelled: %w", ctx.Err())
	}
}
```

## Usage

```go
p, err := objectpool.New(2, func(ctx context.Context) (*Conn, error) {
	return dial(ctx)
})
if err != nil { /* ... */ }

obj, err := p.Get(ctx) // immediate, lazily-created, or blocks until one frees up
if err != nil { /* handle ErrPoolClosed or ctx cancellation */ }
defer p.Put(obj)       // return for reuse

obj.Do( /* ... */ )
```

## Rules of Thumb

- Object pool pattern is useful in cases where object initialization is more
  expensive than the object maintenance.
- If there are spikes in demand as opposed to a steady demand, the maintenance
  overhead might overweigh the benefits of an object pool.
- It has positive effects on performance due to objects being initialized beforehand.
