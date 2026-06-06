# Object Pool

`Pool[T]` — a bounded, concurrency-safe set of pre-allocated, reusable objects.
Callers `Get` an object, use it, then `Put` it back. The pool caps live objects
and provides back-pressure instead of unbounded allocation.

## When to use

- Construction is much more expensive than keeping the object alive (DB
  connections, large scratch buffers, parsers, compiled regexes).
- You need a **hard upper bound** on simultaneously live resources.
- You want **context-aware** blocking acquisition under contention.

For pure zero-allocation reuse without a strict cap, prefer the standard
library's `sync.Pool` instead.

## API

```go
p, err := objectpool.New[*bytes.Buffer](
    poolSize,
    func(ctx context.Context) (*bytes.Buffer, error) { return &bytes.Buffer{}, nil }, // factory
    func(b *bytes.Buffer) { b.Reset() },                                              // reset (optional, may be nil)
)

obj, err := p.Get(ctx) // immediate, lazily-created, or blocks until one frees up
defer p.Put(obj)       // returns it for reuse (runs reset first)

p.Close(func(b *bytes.Buffer) { /* dispose idle objects */ })
```

## Key properties

- **Lazy creation:** objects are built on demand up to the cap, so a large pool
  costs nothing until it is exercised. The factory runs at most `size` times.
- **Bounded:** the number of simultaneously checked-out objects never exceeds
  `size`; excess `Get` calls block until a `Put` (or `ctx` is cancelled).
- **Context-first:** `Get` returns a wrapped `ctx.Err()` if cancelled while
  waiting; no goroutine leaks.
- **Safe lifecycle:** factory errors release the reserved slot (no capacity
  leak); `Close` is idempotent and disposes idle objects.

## Run

```bash
cd examples && go test -race ./creational/object-pool/
```
