---
name: messaging-push-pull
description: Use when a set of pushers fan-out messages to a pool of pullers that compete to consume each message — load distribution without a central dispatcher.
category: messaging
go-version-min: "1.22"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/messaging/push-pull/
---

# Push & Pull

## Intent

In Push-Pull topology (ZeroMQ calls it PUSH/PULL, sometimes called "pipeline" or "task distribution"), pushers send work items and pullers compete to receive them. Each item is delivered to exactly one puller — they share the queue. This is load-balanced distribution without a central scheduler: the fastest puller naturally picks up more items.

The key distinction from broadcast (fan-out to all) is exclusivity: each message goes to exactly one consumer.

## Context

The pattern appears wherever work must be distributed across a dynamic pool of workers without a routing key or topic: task queues, image/video processing pipelines, log ingestion fans, map phase of map-reduce. The implementation in Go is a single shared buffered channel that all pullers receive from — the Go scheduler distributes fairly across blocked receivers.

Comparison:
- **Broadcast**: one sender, all receivers get every message. Push-Pull: one sender, exactly one receiver per message.
- **Producer-Consumer**: conceptually the same at the channel level. Push-Pull emphasises the competitive pulling semantic and often has many pushers and many pullers.
- **Fan-out**: one producer, one channel per receiver (no competition). Push-Pull: one shared channel, receivers compete.

## Implementation in Go

```go
// A single buffered channel is the entire Push-Pull topology.
func NewPipeline[T any](buf int) (push func(T), pull func() (T, bool), close_ func()) {
    ch := make(chan T, buf)
    return func(v T) { ch <- v },
           func() (T, bool) { v, ok := <-ch; return v, ok },
           func() { close(ch) }
}

// Pushers: goroutines that call push(item).
// Pullers: goroutines that loop: for { v, ok := pull(); if !ok { return } ... }
```

A more structured API separates pusher and puller roles into a typed `Pipeline[T]` struct:

```go
type Pipeline[T any] struct{ ch chan T }

func NewPipeline[T any](buf int) *Pipeline[T] { return &Pipeline[T]{make(chan T, buf)} }
func (p *Pipeline[T]) Push(ctx context.Context, v T) error {
    select { case p.ch <- v: return nil; case <-ctx.Done(): return ctx.Err() }
}
func (p *Pipeline[T]) Pull() <-chan T { return p.ch }  // pullers range over this
func (p *Pipeline[T]) Close()          { close(p.ch) }
```

## When to use

- Distributing tasks across a dynamic pool of workers where the fastest available worker should get the next item.
- Task queues: N ingestion goroutines push items; M worker goroutines pull and process.
- Pipeline stages: output of one stage is the input queue for the next — each item flows through exactly once.
- Throttling: the channel's capacity provides backpressure; pushers block when the queue is full.

## When NOT to use

- When every receiver must see every message — use broadcast (fan-out channels or `sync.Cond.Broadcast`).
- When messages must be routed based on content — use a topic-based pub-sub or a dispatcher goroutine.
- When strict ordering per producer must be preserved end-to-end — a single puller (no competition) or a keyed dispatcher.
- When the number of pullers is exactly one — plain producer-consumer; the "competing" aspect adds no value.

## Gotchas

- **Close from multiple pushers**: only close the channel once. Use a `sync.WaitGroup` sentinel goroutine (same as producer-consumer).
- **Puller goroutine leak**: if pushers stop but the channel is never closed, pullers block forever. Always close on shutdown.
- **Back-pressure sizing**: if pushers are faster than pullers and the buffer fills, pushers block. Monitor channel length (`len(ch)`) in production; alert when consistently at capacity.
- **Context not propagated to Push**: pushers that ignore ctx will keep enqueuing after cancellation, potentially overwhelming the queue.

## See also

- skills/concurrency/producer-consumer.md
- skills/concurrency/broadcast.md
- skills/messaging/futures-promises.md
- examples/messaging/push-pull/
