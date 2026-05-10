---
name: messaging-futures-promises
description: Use when a computation runs asynchronously and callers need a handle to retrieve the result later — decouples launch from retrieval, enables concurrent fan-out with result collection.
category: messaging
go-version-min: "1.22"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/messaging/futures-promises/
---

# Futures & Promises

## Intent

A Future is a read-only placeholder for a value that will be computed asynchronously. A Promise is the write side — the entity that fulfils or rejects the future. Together they decouple "start the computation" from "wait for its result", allowing callers to kick off multiple concurrent operations and collect results at their convenience.

## Context

Go does not have built-in futures. The idiomatic equivalent is: goroutine + channel. `go func() { ch <- compute() }()` is a one-liner future. The caller blocks on `<-ch` when it needs the result. This is so natural in Go that dedicated "Future" types are rare — most code just passes channels. But a typed `Future[T]` wrapper is useful when:

- The result must be retrievable multiple times by multiple callers (fan-out read).
- The computation must be cancellable via context.
- You want to compose futures (wait-all, wait-any).
- You want to enforce a clean API that hides the channel from callers.

The `Promise[T]`/`Future[T]` split mirrors Scala/Java futures: `Promise` is the producer side (`Resolve`/`Reject`), `Future` is the consumer side (`Get`).

## Implementation in Go

```go
type result[T any] struct {
    val T
    err error
}

type Future[T any] struct{ ch <-chan result[T] }

type Promise[T any] struct{ ch chan result[T] }

func NewPromise[T any]() (*Promise[T], Future[T]) {
    ch := make(chan result[T], 1) // buffered: Resolve never blocks
    return &Promise[T]{ch}, Future[T]{ch}
}

func (p *Promise[T]) Resolve(v T)    { p.ch <- result[T]{val: v} }
func (p *Promise[T]) Reject(err error) { p.ch <- result[T]{err: err} }

func (f Future[T]) Get(ctx context.Context) (T, error) {
    select {
    case r := <-f.ch:
        return r.val, r.err
    case <-ctx.Done():
        var zero T
        return zero, ctx.Err()
    }
}
```

A buffered channel of size 1 means `Resolve`/`Reject` never block the producer, and the result is available to any number of callers — though only one gets it (use `sync.Once` + cached result if N callers need the same value).

## When to use

- Kick off N independent async computations and collect results as each finishes.
- API boundaries where the caller should not know about goroutines or channels.
- Timeout/cancellation on long-running computations: `Get(ctx)` returns when ctx expires.
- "Memoised async lookup": compute once, distribute the future, all callers wait for the same result.

## When NOT to use

- Simple one-off async work with no result needed — just `go func()`.
- Streaming results — a channel iterator or `iter.Seq` (Go 1.23) is cleaner.
- When multiple callers must each receive the value — add a fan-out layer or use a shared cached result with `sync.Once`.
- When you already use `errgroup` — `Group.Go` + closures capturing a result slice is more Go-idiomatic for concurrent fan-out.

## Gotchas

- **Channel size 0 (unbuffered)**: if no one calls `Get` before `Resolve`, the producer goroutine blocks forever — goroutine leak. Always buffer the result channel (size 1 for single-result futures).
- **Multiple Resolve calls**: a second `Resolve` on a buffered-1 channel blocks the caller. Guard with `sync.Once`.
- **Result channel read by multiple goroutines**: only one goroutine receives from an unbuffered channel. Fan-out via `sync.Once` + caching if N callers need the same result.
- **Context not cancelling the computation**: `Get` can time out, but the underlying goroutine keeps running. Pass `ctx` into the computation itself.

## See also

- skills/messaging/push-pull.md
- skills/concurrency/producer-consumer.md
- examples/messaging/futures-promises/
