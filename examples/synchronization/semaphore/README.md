# Semaphore

`Semaphore` — a counting semaphore backed by a buffered channel. It limits how
many goroutines may hold a "ticket" at once, throttling parallelism without
spinning up a worker pool. Callers run their own goroutines and gate the hot
section with `Acquire`/`Release`.

## When to use

- Cap concurrent calls to a downstream dependency (DB connections, an API quota).
- Bound fan-out so a burst of work does not overwhelm a shared resource.
- Shed load fast (`TryAcquire`) when every slot is busy, instead of queueing.

Prefer this over a worker pool when the goroutines already exist and you only
need to limit how many of them are active simultaneously.

## API

```go
s := semaphore.New(3) // at most 3 concurrent holders

// Blocking, context-aware (cancellation / deadline):
if err := s.Acquire(ctx); err != nil {
    return err // no ticket held
}
defer s.Release()

// Non-blocking, fail fast:
if err := s.TryAcquire(); errors.Is(err, semaphore.ErrNoTickets) {
    return // skip / shed load
}
defer s.Release()

// Convenience: acquire, run, release (even on panic):
err := s.Do(ctx, func(ctx context.Context) error { return work(ctx) })
```

## Key properties

- At most N goroutines hold a ticket at once (verified under `-race`).
- `Acquire` is context-first: returns `ctx.Err()` (wrapped with `%w`) on
  cancellation and leaks no timer — no `time.After`.
- `Release` must pair with a successful acquire; calling it without a held
  ticket panics rather than silently letting an extra holder through.

## Run

```bash
cd examples && go test -race ./synchronization/semaphore/
```
