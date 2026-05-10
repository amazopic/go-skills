---
name: anti-patterns-cascading-failures
description: Avoid when a slow or failed downstream causes goroutines to accumulate upstream because callers lack timeouts, concurrency caps, or circuit breakers — leading to whole-stack collapse.
category: anti-patterns
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/anti-patterns/cascading-failures/
---

# Cascading Failures

## Intent

Cascading failures are an anti-pattern: a downstream outage spreads upstream
because callers block indefinitely, exhaust goroutine/connection pools, and
eventually bring down layers that had no direct dependency on the failed
service. Understanding the mechanism is the first step to preventing it.

## Context

Consider three services: `Frontend → OrderService → PaymentService`.
PaymentService becomes slow — say, 30-second responses. Without mitigations:

1. OrderService goroutines block for 30 seconds waiting for PaymentService.
2. Inbound Frontend requests pile up because OrderService is not responding.
3. Frontend goroutines and connection slots exhaust.
4. The entire stack is unavailable — for a failure that originated in one service.

The root cause is always two factors multiplied: **unbounded wait time** ×
**unbounded concurrency**. Remove either factor and the cascade breaks.

## The Four Mitigations

### 1. Per-call timeout

Cap how long a single goroutine can wait. Use `context.WithTimeout`:

```go
func (c *Client) Call(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
    defer cancel()
    return c.dep.Call(ctx)
}
```

Goroutines now live at most 200ms, not 30 seconds. Memory stays bounded.

### 2. Bulkhead (concurrency cap)

Cap how many goroutines can be waiting simultaneously. Excess callers are
shed immediately rather than queued:

```go
select {
case sem <- struct{}{}:
    defer func() { <-sem }()
default:
    return ErrPoolExhausted // fail-fast, not enqueue
}
```

Even if the timeout fires slowly, at most `cap(sem)` goroutines block.

### 3. Circuit breaker

After N consecutive failures, stop sending altogether and fail fast. This
gives the downstream time to recover without absorbing more load:

```go
if err := cb.Allow(); err != nil {
    return ErrCircuitOpen // instant rejection, no RPC attempted
}
err := doRPC(ctx)
cb.Record(err)
return err
```

### 4. Retry budget

Retries without backoff multiply load on a struggling service. Use exponential
backoff with jitter and a global retry budget (max retries per time window)
rather than per-call unlimited retries.

## Composition

All four mitigations compose cleanly as function wrappers:

```
FullyProtectedClient.Call()
  → cb.Allow()          (circuit: fast-fail if open)
  → sem acquire         (bulkhead: shed if full)
  → WithTimeout(ctx)    (timeout: cap goroutine lifetime)
  → dep.Call(ctx)       (actual RPC)
  → cb.Record(err)      (circuit: record outcome)
```

Apply them in this order. The circuit breaker is outermost so it can prevent
work before touching the semaphore.

## When this anti-pattern occurs

- Any outbound HTTP/gRPC/DB call without a `context.WithTimeout`.
- A `http.Client` with no `Timeout` field set (default: infinite).
- A goroutine pool shared across all tenants (no bulkhead).
- A retry loop with `for { err = call(); if err != nil { continue } }` — infinite retries, no backoff.
- A connection pool with `MaxOpenConns` unset (database/sql defaults to unlimited).

## Gotchas

- **`http.DefaultClient` has no timeout.** Always construct a `http.Client{Timeout: ...}` in production.
- **`database/sql` default pool is unlimited.** Set `db.SetMaxOpenConns`, `db.SetMaxIdleConns`, `db.SetConnMaxLifetime`.
- **Retries without jitter amplify the problem.** N services retrying in lockstep at T+1s create a retry thundering herd. Add `rand.Duration(0, 100ms)` jitter.
- **Circuit breaker open ≠ error.** Callers should distinguish `ErrCircuitOpen` (back off and retry later) from `ErrServiceDown` (alert ops). They have different SRE responses.
- **Timeout too aggressive.** A 50ms timeout on a call that legitimately needs 200ms causes cascade from healthy slowness. Calibrate timeouts from p99 latency in steady state, not best-case.

## See also

- `skills/stability/circuit-breaker.md`
- `skills/stability/bulkhead.md`
- `skills/stability/deadline.md`
- `skills/stability/fail-fast.md`
- `examples/anti-patterns/cascading-failures/`
