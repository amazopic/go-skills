# Cascading Failures Example

Demonstrates how a slow downstream dependency spreads upstream when callers
lack timeouts and concurrency caps — and shows how each mitigation layer
(timeout, bulkhead, circuit breaker) prevents the cascade.

## Structure

| File | Purpose |
|---|---|
| `cascading.go` | `SlowDep`, `NaiveClient` (broken), `TimeoutClient`, `BulkheadClient`, `CircuitBreaker`, `FullyProtectedClient` |
| `cascading_test.go` | Tests showing goroutine accumulation in the naive case, then proving each mitigation works independently and together |

## Run

```bash
go test -race ./anti-patterns/cascading-failures/
```

## Client progression

| Client | Protection |
|---|---|
| `NaiveClient` | None — goroutines pile up |
| `TimeoutClient` | Per-call timeout caps goroutine lifetime |
| `BulkheadClient` | Timeout + semaphore sheds excess callers immediately |
| `FullyProtectedClient` | Timeout + bulkhead + circuit breaker fast-fails when error rate is high |

## Key points

- The root cause of cascading failures is always the same: unbounded waiting multiplied by unbounded concurrency.
- Apply mitigations in composition order: timeout → bulkhead → circuit breaker. Each layer is independently meaningful.
- `FullyProtectedClient` demonstrates that 20 concurrent callers against a 500ms-hang dep all return in under 300ms total when timeout=30ms and bulkhead=3.
