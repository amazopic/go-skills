# Bulkhead Example

Isolates resource pools per tenant or call type using per-pool semaphores.
Failure or exhaustion in one pool cannot drain capacity from another.

## Structure

| File | Purpose |
|---|---|
| `bulkhead.go` | `Pool` (single-pool semaphore) and `Bulkhead` (named pool registry) |
| `bulkhead_test.go` | Table-driven, race-safe tests covering isolation, rejection, and metrics |

## Run

```bash
go test -race ./stability/bulkhead/
```

## Key points

- Each `Pool` owns a buffered channel as a semaphore; acquiring it never blocks other pools.
- Overflow is fail-fast (`ErrPoolExhausted`), not queued — prevents unbounded goroutine accumulation.
- `atomic.Int64` counters for acquired/rejected/current are safe under `-race`.
- `Bulkhead` is a thin named registry; add pools at startup for each resource class (payment, search, notifications).
