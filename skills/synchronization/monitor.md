---
name: synchronization-monitor
description: Use when shared mutable state and its synchronisation logic should be encapsulated in a single type — exported methods are the only safe entry points, and the mutex is private.
category: synchronization
go-version-min: "1.22"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/synchronization/monitor/
---

# Monitor

## Intent

A monitor bundles shared state with the mutex that protects it into one type, then exposes only synchronised public methods. Callers never touch the lock — they call methods and the monitor guarantees thread safety internally. This is the object-oriented formalisation of "lock discipline": because the mutex is unexported and the methods control all access, it is impossible for a caller to forget to lock or to accidentally bypass synchronisation.

## Context

Hoare (1974) defined a monitor as a programming construct in which data and their synchronisation are encapsulated together. In Go, this maps naturally to a struct with unexported fields (data + mutex) and exported methods (operations). Every Go type with a `sync.Mutex` or `sync.RWMutex` field and consistent lock discipline is a monitor, whether the author uses that name or not.

The monitor is the idiomatic Go answer to the question "how do I share state safely?" — more so than raw locks scattered across the codebase, because the type enforces the discipline at compile time (callers cannot see the mutex). A `sync.Map` is a standard-library monitor. A DB connection pool is a monitor. An HTTP server's connection table is a monitor.

## Implementation in Go

```go
// Pool is a monitor: private state + private mutex + public methods.
type Pool struct {
    mu       sync.Mutex
    cond     *sync.Cond
    items    []Item
    maxItems int
}

func NewPool(max int) *Pool {
    p := &Pool{maxItems: max}
    p.cond = sync.NewCond(&p.mu)
    return p
}

// Acquire blocks until an item is available and returns it.
func (p *Pool) Acquire() Item {
    p.mu.Lock()
    defer p.mu.Unlock()
    for len(p.items) == 0 {
        p.cond.Wait()
    }
    item := p.items[len(p.items)-1]
    p.items = p.items[:len(p.items)-1]
    return item
}

// Release returns an item to the pool.
func (p *Pool) Release(item Item) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.items = append(p.items, item)
    p.cond.Signal()
}
```

The caller sees `Acquire` and `Release`. The mutex and condition variable are invisible. Adding a new method that forgets to lock would be a bug visible at code review — the pattern doesn't eliminate that risk, but it localises it to one file.

## When to use

- Any type that wraps shared mutable state: caches, pools, registries, counters with compound operations.
- When the synchronisation logic is non-trivial (condition variables, invariants) and must be hidden from callers.
- When building library code that must be goroutine-safe without burdening callers with lock management.
- Thread-safe data structures that expose high-level operations (Get-or-set, acquire-and-modify).

## When NOT to use

- When there is only one goroutine accessing the data — no synchronisation needed.
- When the state transitions can be expressed cleanly as channel messages — a goroutine owning the data is simpler.
- When the mutex grain is too coarse and contention is measured — consider sharding the monitor (N shards, each its own mutex) or lock-free structures.

## Gotchas

- **Method calls other methods on the same monitor**: if `MethodA` locks and then calls `MethodB` which also locks — deadlock. Pattern: unexported "inner" methods that assume the lock is held, called by exported methods that lock/unlock.
- **Returning references to internal state**: returning a pointer to a field (e.g., `*p.items[0]`) allows the caller to mutate it without holding the lock — a hidden race. Return copies or use read-only view types.
- **Copy of the monitor struct**: copying a struct with a `sync.Mutex` copies the lock state — `go vet` catches this. Always pass monitors by pointer.
- **Large critical sections**: putting complex computation inside the monitor's lock degrades concurrency. Compute outside, lock only to update state.

## See also

- skills/synchronization/mutex.md
- skills/synchronization/condition-variable.md
- skills/synchronization/read-write-lock.md
- examples/synchronization/monitor/
