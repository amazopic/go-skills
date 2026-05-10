# Deadline Example

Bounds the maximum wall-clock time a request may occupy across all layers.
Demonstrates end-to-end budget propagation via `context.WithDeadline`,
per-attempt timeout capping, and early budget checks at service boundaries.

## Structure

| File | Purpose |
|---|---|
| `deadline.go` | `WithBudget`, `CheckBudget`, `Call.Do`, `Chain`, `Attempt`, `IsDeadlineError` |
| `deadline_test.go` | Table-driven tests covering budget shrinkage, short-circuit, and parent/child deadline interaction |

## Run

```bash
go test -race ./stability/deadline/
```

## Key points

- `context.WithDeadline` automatically picks the shorter of parent vs. requested deadline — never set a longer deadline than the parent.
- `CheckBudget` at each layer boundary enables fail-fast before expensive work when budget is already insufficient.
- `Attempt` caps per-attempt timeouts to the remaining parent budget — critical for retry loops that must not exceed the end-to-end SLA.
- Always `defer cancel()` immediately after `WithDeadline`/`WithTimeout` to release the timer.
