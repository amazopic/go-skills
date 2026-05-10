---
name: concurrency-broadcast
description: Use when one sender must deliver the same message or signal to an unbounded number of receivers simultaneously — e.g., config reload, shutdown signal, cache invalidation.
category: concurrency
go-version-min: "1.22"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/concurrency/broadcast/
---

# Broadcast

## Intent

Broadcast delivers one event to every current subscriber atomically — all receivers unblock at the same logical instant. Unlike a regular channel send (which wakes exactly one goroutine), a broadcast wakes all of them. The pattern captures the distinction between "pick one worker to handle this" (fan-out / work distribution) and "notify everyone about this" (broadcast / publish).

## Context

Broadcast arises whenever shared state changes and every interested party must react: configuration reloads, shutdown signals, cache invalidation, leader-election notifications, "ready" signals in test harnesses. Two idiomatic Go mechanisms deliver it:

1. **Closing a channel** (`close(ch)`) — the simplest broadcast. A closed channel is readable by any number of goroutines with immediate, non-blocking receipt of the zero value. Works for one-shot signals (e.g., `context.Done()`).
2. **`sync.Cond.Broadcast()`** — for repeated or data-carrying events where the set of waiters changes dynamically. Waiters register under a mutex, wait on the condition, and the broadcaster wakes all of them in one call.

A third approach — maintaining a slice of per-subscriber channels — gives more control (typed payloads, back-pressure) but adds bookkeeping.

## Implementation in Go

**One-shot broadcast via channel close** (recommended for signals):

```go
type Signal struct{ once sync.Once; ch chan struct{} }

func NewSignal() *Signal { return &Signal{ch: make(chan struct{})} }

// Broadcast fires the signal; safe to call multiple times.
func (s *Signal) Broadcast() { s.once.Do(func() { close(s.ch) }) }

// Wait blocks until Broadcast has been called.
func (s *Signal) Wait() <-chan struct{} { return s.ch }
```

Callers do `<-sig.Wait()` or `select { case <-sig.Wait(): ... case <-ctx.Done(): ... }`.

**Repeatable broadcast via fan-out channels** (for typed payloads):

```go
type Broadcaster[T any] struct {
    mu   sync.Mutex
    subs []chan T
}

func (b *Broadcaster[T]) Subscribe(buf int) <-chan T {
    ch := make(chan T, buf)
    b.mu.Lock()
    b.subs = append(b.subs, ch)
    b.mu.Unlock()
    return ch
}

func (b *Broadcaster[T]) Send(v T) {
    b.mu.Lock()
    defer b.mu.Unlock()
    for _, ch := range b.subs {
        select {
        case ch <- v:
        default: // drop if subscriber is full; adjust policy as needed
        }
    }
}
```

## When to use

- One-time signal to all goroutines: shutdown, "ready", config loaded — use `close(ch)` wrapped in `sync.Once`.
- Repeatable typed notifications to a known-at-startup set of listeners — fan-out channels.
- Dynamic subscription list with complex wait predicates — `sync.Cond.Broadcast`.
- `context.WithCancel` already provides broadcast semantics for cancellation; reuse it before rolling your own.

## When NOT to use

- When only one goroutine should handle the event — use a plain channel send.
- When you need guaranteed delivery with back-pressure — add buffered subscriber channels and handle the full case explicitly instead of dropping.
- High-frequency data streams to many subscribers — consider a ring-buffer or event-bus library to avoid per-subscriber heap allocation per message.

## Gotchas

- **Closing a closed channel panics.** Always guard with `sync.Once` or an explicit closed-flag under a mutex.
- **`sync.Cond` spurious wakeups**: always re-check the predicate in a `for` loop after `Wait()` returns — it may have been woken by `Signal` (not `Broadcast`) or a spurious wake.
- **Subscriber leak**: if a goroutine subscribes but exits without draining its channel, the `Send` loop may block or drop messages indefinitely. Provide an `Unsubscribe` or use a `select/default` drop policy.
- **Lock ordering**: do not call `Broadcast()` while holding a lock that subscribers also acquire — classic deadlock.

## See also

- skills/concurrency/producer-consumer.md
- skills/messaging/push-pull.md
- skills/synchronization/condition-variable.md
- examples/concurrency/broadcast/
