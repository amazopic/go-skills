---
name: synchronization-mutex
description: Use when multiple goroutines share mutable state that cannot be expressed as a channel message — protect reads and writes with sync.Mutex or sync.RWMutex.
category: synchronization
go-version-min: "1.22"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/synchronization/mutex/
---

# Mutex

## Intent

A mutex (mutual exclusion lock) ensures that only one goroutine at a time executes a critical section — a block of code that reads or writes shared state. Every other goroutine that tries to lock the mutex blocks until the current holder calls `Unlock`. In Go, `sync.Mutex` is the standard implementation; `sync.RWMutex` is its read-optimised variant.

## Context

Go's preferred model is "communicate by sharing" — pass data through channels instead of sharing memory directly. But shared memory is sometimes unavoidable: a cache, an in-process counter, a map that multiple goroutines must update. For these, a mutex is the right tool. Keep the critical section as short as possible: only the lines that touch shared state, not I/O or computation.

`sync.RWMutex` allows any number of concurrent readers (`RLock`/`RUnlock`) but only one writer (`Lock`/`Unlock`). It wins when reads are overwhelmingly more frequent than writes. Under write-heavy workloads, `RWMutex` can be slower than `Mutex` due to the additional bookkeeping.

## Implementation in Go

```go
type SafeCounter struct {
    mu sync.Mutex
    v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock() // ALWAYS defer to avoid forget-to-unlock bugs
    c.v[key]++
}

func (c *SafeCounter) Value(key string) int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.v[key]
}
```

For read-heavy counters, swap to `sync.RWMutex`:

```go
type ReadHeavyCounter struct {
    mu sync.RWMutex
    v  int
}

func (c *ReadHeavyCounter) Inc()      { c.mu.Lock(); defer c.mu.Unlock(); c.v++ }
func (c *ReadHeavyCounter) Value() int { c.mu.RLock(); defer c.mu.RUnlock(); return c.v }
```

## When to use

- Multiple goroutines read and/or write the same in-memory data structure (map, slice, struct fields).
- The critical section is too complex or stateful to express as a channel message.
- You need a simple atomic read-modify-write on a value that cannot be served by `sync/atomic`.
- Read-heavy shared data (cache, config snapshot) — use `sync.RWMutex`.

## When NOT to use

- When data ownership can be transferred via channels — prefer channels.
- For a single integer counter — `sync/atomic.Int64` is faster and lock-free.
- For "write once, read many" configuration — `sync.Once` + immutable value, or `atomic.Value`.
- When goroutines do not actually share data — no mutex needed if each goroutine owns its data.

## Gotchas

- **Copy trap**: `sync.Mutex` must not be copied after first use (the lock state is embedded). Pass structs that contain a mutex by pointer, never by value. `go vet` catches this with `copylocks`.
- **Forget to unlock**: a panic inside the critical section leaves the mutex locked forever. Always `defer mu.Unlock()` immediately after `mu.Lock()`.
- **Lock ordering deadlock**: if goroutine A holds mu1 and tries to acquire mu2, while goroutine B holds mu2 and tries mu1 — deadlock. Establish a consistent global lock order and document it.
- **Holding a lock across I/O**: network calls inside a critical section can take seconds, blocking all other goroutines. Copy the data you need, unlock, then do I/O.
- **RLock-then-upgrade**: you cannot promote an `RLock` to a `Lock` without first calling `RUnlock`. Attempting to do so deadlocks (`RLock` prevents new writers; the upgrading goroutine is itself a reader).
- **Starvation with `RWMutex`**: a stream of continuous readers can starve writers indefinitely. Go's implementation mitigates this (writers eventually block new readers), but profile under load.

## See also

- skills/synchronization/read-write-lock.md
- skills/synchronization/monitor.md
- skills/synchronization/condition-variable.md
- examples/synchronization/mutex/
