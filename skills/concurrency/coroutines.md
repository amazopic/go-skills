---
name: concurrency-coroutines
description: Use when you need cooperative, interleaved execution between two or more logical threads of control — e.g., generators, cooperative schedulers, state machines with yield points.
category: concurrency
go-version-min: "1.22"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/concurrency/coroutines/
---

# Coroutines

## Intent

A coroutine is a function that can suspend its execution at a yield point and transfer control to another coroutine, resuming later from exactly where it left off. Unlike threads, coroutines are cooperative: they explicitly yield rather than being preempted. The pattern enables clean expression of state machines, generators, and pipelines where two logical agents take turns producing and consuming.

## Context

Go does not have first-class coroutine syntax. However, goroutines + channels compose into the same cooperative semantics: a "yield" is a channel send (or receive) that blocks the sender until the partner is ready. This is not a limitation — it is idiomatic Go. The send/receive handshake is explicit, deterministic, and race-free by construction.

When the use case is pure generation (one direction of yield), a goroutine closing over a channel is the standard "generator" pattern. When two coroutines truly take turns (ping-pong), two goroutines communicate over two channels — one for each direction of control transfer.

`golang.org/x/exp/coroutine` exists but brings an external dependency; avoid it for production code when stdlib suffices.

## Implementation in Go

```go
// Coroutine encapsulates a goroutine that yields values of type T
// and resumes when the caller calls Next.
type Coroutine[T any] struct {
    yield  chan T
    resume chan struct{}
    done   chan struct{}
}

// Start launches the coroutine body. body receives a yield func that returns
// true to continue or false when the coroutine has been stopped.
// The body should return when yield returns false.
func Start[T any](ctx context.Context, body func(yield func(T) bool)) *Coroutine[T] {
    c := &Coroutine[T]{
        yield:  make(chan T),
        resume: make(chan struct{}),
        done:   make(chan struct{}),
    }
    go func() {
        defer close(c.yield)
        body(func(v T) bool {
            select {
            case c.yield <- v: // suspend: deliver value
            case <-ctx.Done():
                return false
            case <-c.done:
                return false
            }
            select {
            case <-c.resume:
                return true
            case <-ctx.Done():
                return false
            case <-c.done:
                return false
            }
        })
    }()
    return c
}

// Next advances the coroutine and returns (value, true) or (zero, false)
// when the body has returned or the coroutine has been stopped.
func (c *Coroutine[T]) Next() (T, bool) {
    v, ok := <-c.yield
    if ok {
        select {
        case c.resume <- struct{}{}:
        case <-c.done:
        }
    }
    return v, ok
}

// Stop signals the coroutine to cease yielding. Idempotent.
func (c *Coroutine[T]) Stop() { /* closes c.done via sync.Once */ }
```

The send/receive pattern ensures that exactly one goroutine is running at any time — true cooperative scheduling.

## When to use

- **Generators**: produce an infinite or lazy sequence (Fibonacci, prime sieve, paginated API results) without materialising the whole slice.
- **State machines** with complex transitions that are cleaner expressed as sequential code than as explicit state enum switches.
- **Cooperative ping-pong**: two logical agents alternating turns (parser ↔ lexer, test driver ↔ system-under-test).
- Protocol implementations where client and server speak in alternating turns.

## When NOT to use

- When the data flows in one direction at high throughput — a buffered channel pipeline is more efficient.
- When you need parallelism rather than interleaving — use goroutines without the cooperative handshake.
- When the "coroutine" is just an iterator over a slice — use a plain `for` loop or closure returning `iter.Seq` (Go 1.23+).
- When the body can panic and must be recovered safely — the paired channels leak if the goroutine panics without draining; add a recover + close.

## Gotchas

- **Goroutine leak on early exit**: if the caller stops calling `Next` before the body returns, the goroutine blocks forever on `c.yield <- v`. Add a `done` channel or context to the yield function so the goroutine can detect abandonment.
- **Not using generics (pre-1.18)**: coroutines over `any` require type assertions at every `Next` call — messy. Prefer the generic form shown above.
- **Deadlock with buffered channels**: using a buffered `yield` channel breaks the strict cooperative semantics — the body can run ahead of the caller.
- **Shared state between coroutine and caller**: even though only one runs at a time, the Go race detector does not know that. Use channel values to pass data; do not share pointers without explicit synchronisation.

## See also

- skills/concurrency/producer-consumer.md
- skills/concurrency/reactor.md
- examples/concurrency/coroutines/
