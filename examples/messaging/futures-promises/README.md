# Futures & Promises

Generic `Future[T]` / `Promise[T]` for Go 1.21+.

- `Promise.Resolve(v)` / `Promise.Reject(err)` — fulfil once; subsequent calls
  are no-ops (guarded by `sync.Once`).
- `Future.Get(ctx)` — blocks until resolved or ctx expires. Multiple goroutines
  may call Get on the same Future; all receive the same cached result.
- `Go[T](ctx, fn)` — convenience: launches fn in a goroutine and returns a
  Future.

## Run

```bash
go test -race -v ./messaging/futures-promises/
```

## Why Go rarely needs this

A plain goroutine + buffered channel is a one-liner future. Use the typed
wrapper when: (a) multiple callers need the same result, (b) you want a clean
API that hides channels, or (c) context cancellation must be plumbed through.
