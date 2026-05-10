---
name: synchronization-condition-variable
description: Use when goroutines must wait for a shared predicate that can only become true as a side-effect of another goroutine's action — e.g., bounded queues, thread pools, event notification with shared state.
category: synchronization
go-version-min: "1.22"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/synchronization/condition-variable/
---

# Condition Variable

## Intent

A condition variable allows goroutines to block efficiently until a predicate on shared state becomes true, without spinning. The goroutine atomically releases the mutex and suspends; when another goroutine signals the condition, the waiter reacquires the mutex and re-checks the predicate. This three-step dance (lock → check predicate → wait/proceed) is the foundation of all wait/notify synchronisation.

## Context

Go's `sync.Cond` implements the Mesa-style condition variable (the same semantics as `pthread_cond_t`). It is lower-level than channels but more efficient when many goroutines share a single predicate on a large piece of state — because a channel-based solution would need a goroutine per waiter or a broadcast channel that carries no payload.

The canonical use case is a **bounded queue** (producer blocks when full, consumer blocks when empty) or a **fixed-size resource pool** (goroutine blocks until a resource is available). In both cases, state is a shared data structure — not a simple token — so `sync.Cond` is the right fit.

Modern Go often prefers channels for wait/notify. Use `sync.Cond` when:
- You need `Broadcast` (wake all waiters) with shared complex state, or
- `Signal` (wake exactly one waiter) and you want O(1) rather than per-waiter channels, or
- You are wrapping a large existing struct and do not want to change its storage shape.

## Implementation in Go

```go
type BoundedQueue struct {
    mu    sync.Mutex
    cond  *sync.Cond
    buf   []int
    cap   int
}

func NewBoundedQueue(cap int) *BoundedQueue {
    q := &BoundedQueue{cap: cap}
    q.cond = sync.NewCond(&q.mu)
    return q
}

func (q *BoundedQueue) Put(v int) {
    q.mu.Lock()
    defer q.mu.Unlock()
    for len(q.buf) == q.cap { // ALWAYS loop — Mesa semantics
        q.cond.Wait()         // releases mu, suspends, reacquires mu on wake
    }
    q.buf = append(q.buf, v)
    q.cond.Signal() // wake one waiter (a consumer)
}

func (q *BoundedQueue) Get() int {
    q.mu.Lock()
    defer q.mu.Unlock()
    for len(q.buf) == 0 { // ALWAYS loop
        q.cond.Wait()
    }
    v := q.buf[0]
    q.buf = q.buf[1:]
    q.cond.Signal() // wake one waiter (a producer)
    return v
}
```

## When to use

- **Bounded buffer / pool**: producers block when full, consumers block when empty — one `sync.Cond` shared by both sides with `Signal`.
- **Wake-all on state change**: e.g., all workers waiting for "leader elected" — use `Broadcast`.
- **Replacing a per-goroutine channel for scalability**: N goroutines waiting on one condition is O(1) memory vs O(N) channels.

## When NOT to use

- When the predicate reduces to a simple on/off signal — a closed channel (`close(ch)`) is cleaner and context-cancellable.
- When there is exactly one waiter — a plain channel send/receive is simpler.
- When the predicate depends on data flowing through the signal — embed the data in the channel message.
- For timeouts/context: `sync.Cond.Wait` is not interruptible. Implement a timed-wait with a helper goroutine and `Broadcast`, or redesign with a channel.

## Gotchas

- **If block instead of for loop**: Mesa semantics guarantee the waiter holds the lock when `Wait` returns, but do NOT guarantee the predicate is still true. Always use `for pred { cond.Wait() }`.
- **Signal outside the lock**: calling `Signal` or `Broadcast` without holding the mutex is technically allowed but creates a race where the signal fires before the waiter enters `Wait`. Lock before signalling when in doubt.
- **`sync.Cond` is not copyable**: embedding `sync.Cond` in a struct and then copying the struct loses the internal state. Always store as a pointer.
- **Deadlock on panic**: if a goroutine panics inside the critical section it holds the mutex forever. Use `defer mu.Unlock()`.
- **No context support**: `Wait` blocks indefinitely. For cancellable waits, run a goroutine that sleeps and then calls `Broadcast` after the deadline.

## See also

- skills/synchronization/mutex.md
- skills/synchronization/monitor.md
- skills/concurrency/broadcast.md
- examples/synchronization/condition-variable/
