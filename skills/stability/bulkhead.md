---
name: stability-bulkhead
description: Use when you need to prevent failure or saturation in one consumer (tenant, call type) from draining shared resources and taking down unrelated consumers.
category: stability
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/stability/bulkhead/
---

# Bulkhead

## Intent

Partition resource pools — goroutines, connections, semaphore slots — into
separate compartments per tenant or call type. Exhaustion or failure inside
one compartment cannot spill into others. The name comes from watertight
bulkheads in ship hulls: flooding one compartment does not sink the vessel.

## Context

A single shared thread pool or connection pool is the default in most services.
Under load, a noisy tenant or a slow call type saturates the pool and starves
all other callers. Circuit breakers detect failure after the fact; bulkheads
prevent the resource drain in the first place. The two patterns are
complementary: bulkheads cap concurrency, circuit breakers cap error propagation.

Use bulkheads when:
- You serve multiple tenants or priority classes from one process.
- Different call types have wildly different latency profiles (e.g. fast reads vs. slow exports).
- You need to guarantee that a degraded dependency can only harm its own callers.

## Implementation in Go

Go's idiomatic bulkhead is a **buffered channel as a semaphore** — one per pool.
A non-blocking send acquires a slot; a receive releases it. Overflow is
rejected immediately, never queued.

```go
type Pool struct {
    sem chan struct{}
}

func NewPool(limit int) *Pool {
    return &Pool{sem: make(chan struct{}, limit)}
}

func (p *Pool) Do(ctx context.Context, fn func(context.Context) error) error {
    if ctx.Err() != nil {
        return ctx.Err()
    }
    select {
    case p.sem <- struct{}{}:
        defer func() { <-p.sem }()
    default:
        return ErrPoolExhausted // fail-fast, never queue
    }
    return fn(ctx)
}
```

Multiple pools are registered by name at startup. The caller names the pool
when dispatching work; unknown or exhausted pools return immediately without
touching any shared state.

A `Bulkhead` registry holds `map[string]*Pool` behind a `sync.RWMutex`.
Reads (dispatching work) take a read lock; writes (adding pools at startup)
take a write lock. After startup, only read locks are needed — contention is
negligible.

Metrics (acquired / rejected / current) are tracked with `atomic.Int64` so
tests and dashboards can inspect them without an additional lock.

## When to use

- Multi-tenant SaaS where one tenant's traffic spike must not affect others.
- Mixed priority workloads (interactive vs. batch) sharing one process.
- Gateway or BFF layers that fan out to multiple downstreams with different SLAs.
- Anywhere a retry storm in one path could consume all goroutines.

## When NOT to use

- When you have a single consumer class — a global semaphore suffices.
- When calls are so short-lived that semaphore overhead matters more than isolation.
- When the downstream can handle unlimited concurrency and your bottleneck is CPU on the calling side.

## Gotchas

- **Forgetting to release.** Always use `defer func() { <-p.sem }()` inside
  the critical section. A panic before the deferred release permanently leaks
  a slot. The pattern above is safe because the defer fires even on panic.
- **Pool-per-goroutine vs. pool-per-tenant.** Pools should map to resource
  classes, not to individual goroutines. A pool with limit=1 per goroutine
  defeats the purpose.
- **Wrong overflow policy.** Queuing behind a full pool hides the problem.
  Fail-fast and surface `ErrPoolExhausted` so callers know to back off.
- **Static pool sizing.** In production, expose the limit via config and emit
  `current`/`limit` as metrics so you can tune without a deploy.
- **Context cancellation before acquire.** Check `ctx.Err()` before the
  channel select; a cancelled context should never wait for a slot.

## See also

- `skills/stability/circuit-breaker.md`
- `skills/stability/fail-fast.md`
- `skills/anti-patterns/cascading-failures.md`
- `examples/stability/bulkhead/`
