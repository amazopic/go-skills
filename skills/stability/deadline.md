---
name: stability-deadline
description: Use when you need to bound the maximum wall-clock time a request may occupy end-to-end, propagating a shrinking time budget through every layer of the call graph.
category: stability
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/stability/deadline/
---

# Deadline

## Intent

Assign a single end-to-end time budget to a request at ingress and propagate
it unchanged through every downstream call. Each layer automatically inherits
the shrinking remaining budget via `context.Context`. No layer can silently
extend the budget — the deadline is monotonically non-increasing.

## Context

Without deadlines, a slow dependency causes goroutines to block indefinitely.
Retries without a shrinking budget can fire after the client has already given
up, wasting resources. Timeouts per-attempt and end-to-end deadlines solve
different problems: a per-attempt timeout resets on retry; an end-to-end
deadline never resets. Most production systems need both: a short per-attempt
timeout (prevents a single slow RPC from blocking too long) plus an overall
deadline (prevents retries from running after the client SLA expires).

## Implementation in Go

`context.WithDeadline` and `context.WithTimeout` are the primitives.
The key guarantee: if the parent context already has an earlier deadline,
the child inherits it — you cannot accidentally extend a deadline by passing
a later one.

```go
// WithBudget creates a child context capped at budget from now.
// If parent's deadline is sooner, parent's deadline wins.
func WithBudget(parent context.Context, budget time.Duration) (context.Context, context.CancelFunc) {
    return context.WithDeadline(parent, time.Now().Add(budget))
}

// CheckBudget fails fast if the remaining context budget is less than min.
// Embed at layer boundaries before expensive sub-calls.
func CheckBudget(ctx context.Context, min time.Duration) error {
    if ctx.Err() != nil {
        return ctx.Err()
    }
    if dl, ok := ctx.Deadline(); ok && time.Until(dl) < min {
        return ErrBudgetExceeded
    }
    return nil
}
```

For retries, cap the per-attempt timeout to the remaining parent budget so
retries cannot outlive the end-to-end SLA:

```go
func Attempt(parent context.Context, perAttempt time.Duration, fn func(context.Context) error) error {
    rem, _ := RemainingBudget(parent)
    budget := perAttempt
    if rem < perAttempt {
        budget = rem
    }
    ctx, cancel := context.WithTimeout(parent, budget)
    defer cancel()
    return fn(ctx)
}
```

## When to use

- Any RPC or database call that must complete within a customer-facing SLA.
- Retry loops where individual attempts must not collectively exceed the outer deadline.
- Fan-out calls where all branches share the same originating request deadline.
- gRPC services (gRPC propagates deadlines in metadata automatically; use this to enforce them server-side too).

## When NOT to use

- Background jobs with no latency contract — use a cancelable context without a deadline.
- Batch processing pipelines where work should complete regardless of wall clock.
- Long-running streaming RPCs where a hard deadline would cut the stream mid-transfer.

## Gotchas

- **Forgetting `defer cancel()`** after `WithDeadline`/`WithTimeout` leaks the timer goroutine until the parent is cancelled. Always defer immediately after creation.
- **Using `time.Sleep` instead of select.** `time.Sleep` ignores context cancellation. Use `select { case <-time.After(d): ...; case <-ctx.Done(): ... }` or `context.WithTimeout`.
- **Not checking `errors.Is(err, context.DeadlineExceeded)`** — wrapping libraries often return `context.DeadlineExceeded` unwrapped; use `errors.Is`, not `err == context.DeadlineExceeded`.
- **Budget not propagated.** If you create a new `context.Background()` mid-call instead of propagating the parent, the deadline is silently lost. Always thread the context through every function boundary.
- **Clock skew in distributed systems.** The local deadline is meaningless to a remote server. Use gRPC deadlines or pass a `Deadline-Unix-Millis` header so the server can check budget server-side.

## See also

- `skills/stability/fail-fast.md`
- `skills/stability/circuit-breaker.md`
- `skills/anti-patterns/cascading-failures.md`
- `examples/stability/deadline/`
