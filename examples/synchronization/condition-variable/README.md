# Condition Variable

`BoundedQueue` is a capacity-limited FIFO backed by `sync.Cond`:

- `Put` blocks when full (producer waits on `cond`).
- `Get` blocks when empty (consumer waits on `cond`).
- `Signal` wakes exactly one waiter after each state change.
- `Drain` removes all items and calls `Broadcast` to wake all waiters.

## Key rules

1. **Always loop** on the predicate: `for !ready { cond.Wait() }` — Mesa semantics.
2. **Hold the lock** when calling `Wait`, `Signal`, and `Broadcast`.
3. **Never copy** a `sync.Cond`; always store as a pointer.

## Run

```bash
go test -race -v ./synchronization/condition-variable/
```
