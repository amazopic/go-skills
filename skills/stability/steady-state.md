---
name: stability-steady-state
description: Use when you need to bound resource accumulation (memory, disk, connections) so the system reaches stable equilibrium under sustained load rather than growing without limit.
category: stability
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/stability/steady-state/
---

# Steady-State

## Intent

Ensure that every resource the system accumulates is also actively reclaimed
on a schedule, so that memory, disk, and connection usage converge to a stable
ceiling under sustained load rather than growing unboundedly. The system does
not merely survive load — it reaches equilibrium.

## Context

Any cache, log buffer, metrics ring, or connection pool that grows in response
to traffic will exhaust available memory or disk if not actively bounded. The
usual failure mode is an OOM kill hours into a load test, or a disk-full error
at 3 AM. Steady-State addresses this by pairing every accumulation mechanism
with a reclamation mechanism of equal or greater throughput.

Go-specific examples:
- An in-process cache evicting expired entries via a janitor goroutine.
- A metrics ring buffer capped at N entries, overwriting oldest on overflow.
- A background log compactor rotating files when they exceed a size threshold.
- A `sync.Pool` for buffer reuse (Go's GC reclaims pools under pressure automatically).

## Implementation in Go

Two complementary patterns:

**1. TTL cache with a janitor goroutine:**

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]entry  // entry carries expiresAt time.Time
    max   int
    stop  chan struct{}
}

func (c *Cache) Start(ctx context.Context) {
    go func() {
        t := time.NewTicker(c.interval)
        defer t.Stop()
        for {
            select {
            case <-t.C:
                c.evictExpired()
            case <-ctx.Done():
                return
            case <-c.stop:
                return
            }
        }
    }()
}
```

The janitor exits cleanly via `ctx.Done()` or `c.stop` — no goroutine leak.
The cache also enforces a capacity cap at write time, shedding inserts when
full and no expired entry is available to evict.

**2. Fixed-size ring buffer (generic, Go 1.21+):**

```go
type RingBuffer[T any] struct {
    mu   sync.Mutex
    buf  []T
    head int  // next write position
    size int  // number of valid entries
    cap  int
}

func (r *RingBuffer[T]) Push(v T) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.buf[r.head] = v
    r.head = (r.head + 1) % r.cap
    if r.size < r.cap {
        r.size++
    }
    // oldest entry silently overwritten when full — zero allocation
}
```

Memory is bounded at construction time. No dynamic growth; no GC pressure
from append growth; predictable `len(buf) * sizeof(T)` footprint.

## When to use

- In-process caches that must not grow unboundedly under write-heavy workloads.
- Metrics and audit log buffers where the newest N entries are what matter.
- Connection pools that must be periodically pruned of idle connections.
- Any background accumulation that lacks a natural drain (e.g. deduplication sets, bloom filters with a TTL).

## When NOT to use

- When data must be durable — ring buffers and janitor eviction are lossy by design.
- When the reclamation interval is so long that resource usage grows catastrophically between sweeps.
- When `sync.Pool` suffices — the GC already reclaims pooled objects under memory pressure; no janitor needed.

## Gotchas

- **Janitor goroutine leak.** If `Start` is called without a paired `Stop` or a cancellable context, the goroutine runs forever. Always document the lifecycle contract.
- **Eviction under the wrong lock.** `evictExpired` must hold the write lock. Evicting under a read lock causes a data race.
- **Capacity cap that starves long-lived entries.** If the eviction strategy is purely LRU and all entries are non-expired, new writes are always shed. Mix TTL-based and capacity-based eviction.
- **Clock skew in `time.Now()`.** On VMs, wall-clock jumps can cause entries to appear expired that are not. For latency-sensitive eviction, consider monotonic clock (`time.Since`) instead of wall-clock comparison.
- **Ring buffer index arithmetic.** Off-by-one errors in `(head - size + cap) % cap` are the most common bug. Write the formula once, cover it with boundary tests (full, one entry, power-of-two and non-power-of-two capacities).

## See also

- `skills/stability/bulkhead.md`
- `skills/stability/fail-fast.md`
- `examples/stability/steady-state/`
