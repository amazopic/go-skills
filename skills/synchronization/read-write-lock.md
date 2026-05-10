---
name: synchronization-read-write-lock
description: Use when reads vastly outnumber writes on shared mutable state — sync.RWMutex lets many goroutines read concurrently while writes remain exclusive.
category: synchronization
go-version-min: "1.22"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/synchronization/read-write-lock/
---

# Read-Write Lock

## Intent

A read-write lock allows any number of goroutines to read a shared resource simultaneously, but grants write access exclusively to one goroutine at a time, blocking all readers while the write is in progress. The benefit over a plain mutex is throughput: concurrent reads do not block each other, so read-heavy workloads scale with the number of goroutines rather than serialising at the lock.

## Context

`sync.RWMutex` is Go's standard read-write lock. It uses the same interface as `sync.Mutex` for writes (`Lock`/`Unlock`) and adds `RLock`/`RUnlock` for reads. The Go runtime implements a priority scheme that prevents writer starvation: once a writer is waiting, new readers are queued behind it.

Alternatives for read-heavy data:
- `sync/atomic.Value` — lock-free copy-on-write for whole-value snapshots. Zero reader contention.
- `sync.Map` — built-in concurrent map optimised for stable key sets with infrequent writes.
- Copy-on-write with `atomic.Pointer[T]` — writers copy, modify the copy, swap the pointer.

`RWMutex` wins when: the guarded data is too large or structured for copy-on-write; write frequency is low (< ~10% of operations); and multiple readers must see the latest committed state without a copy.

## Implementation in Go

```go
type Config struct {
    mu   sync.RWMutex
    data map[string]string
}

// Get reads a key — many goroutines may call concurrently.
func (c *Config) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.data[key]
    return v, ok
}

// Set writes a key — exclusive; blocks all readers and other writers.
func (c *Config) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

## When to use

- Read-to-write ratio is high (e.g., ≥ 10:1). Profile first — under low contention, `Mutex` is often faster due to simpler implementation.
- Guarded state is large or structurally complex, making copy-on-write expensive.
- In-process caches, configuration snapshots, routing tables, connection registries.

## When NOT to use

- Write-heavy workloads: `RWMutex` overhead from tracking concurrent readers makes it slower than `Mutex`.
- Simple scalar values — use `sync/atomic` types (`atomic.Int64`, `atomic.Value`) which are lock-free.
- When you need atomic swap of an entire value — `atomic.Pointer[T]` with copy-on-write is zero-reader-contention.
- When the data fits in a stable key set with rare writes — `sync.Map` specialises for that case.

## Gotchas

- **Copy trap**: `sync.RWMutex` must not be copied after first use. Pass the enclosing struct by pointer; `go vet` enforces this.
- **RLock-then-Lock upgrade is a deadlock**: you cannot promote a read lock to a write lock. Call `RUnlock` first, then re-check the predicate after acquiring `Lock` (the state may have changed).
- **Starvation of readers by writers**: Go's implementation gives preference to waiting writers, so a flood of writers can starve readers. Under extreme write load, test your workload.
- **Holding RLock across a slow operation**: a goroutine holding `RLock` prevents all writers. Keep `RLock` sections as short as reading the value; release before any I/O.
- **Never call `mu.Lock()` while holding `mu.RLock()`** in the same goroutine — deadlocks because the write lock waits for all read locks to release, including the current goroutine's own.

## See also

- skills/synchronization/mutex.md
- skills/synchronization/monitor.md
- examples/synchronization/read-write-lock/
