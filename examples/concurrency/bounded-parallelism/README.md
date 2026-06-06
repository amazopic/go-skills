# Bounded Parallelism

Process an arbitrarily long input slice with **at most K** worker goroutines
running concurrently. A fixed pool of K workers competes to pull indexed jobs
from a shared channel; each worker writes its result back into a pre-sized
output slice at the input index, so output order matches input order with no
result channel or post-sort.

## When to use

- The input list is large or unbounded and "one goroutine per item" would
  exhaust file descriptors, sockets, DB connections, or memory.
- You need to cap concurrency (e.g. respect an API rate limit or a connection
  pool size) while still parallelizing.

Prefer this over [`parallelism`](../parallelism/) (fan-out everything) whenever
concurrency must be capped.

## API

```go
// Parallel map, results in input order:
out, err := boundedparallelism.Map(ctx, workers, in,
    func(ctx context.Context, item T) (R, error) { ... })

// Side-effecting work, no results collected:
err := boundedparallelism.ForEach(ctx, workers, in,
    func(ctx context.Context, item T) error { ... })
```

## Key properties

- **Bounded:** never more than `workers` invocations of `fn` run at once.
- **Order-preserving:** `out[i]` corresponds to `in[i]`.
- **Context-first:** returns `ctx.Err()` (wrapped) on cancellation; cancellation
  stops dispatch and lets in-flight workers drain.
- **Fail-fast:** the first `fn` error cancels remaining work and is reported
  wrapped with `%w` (so `errors.Is` works); `Map` returns `nil` results.
- **Leak-free:** all goroutines are joined before returning; verified under
  `-race`.
- `workers < 1` returns `ErrInvalidWorkers`.

## Run

```bash
cd examples && go test -race ./concurrency/bounded-parallelism/
```
