---
name: concurrency-reactor
description: Use when a single goroutine must dispatch multiple asynchronous event streams to registered handlers — e.g., in-memory event loops, timer wheels, protocol state machines.
category: concurrency
go-version-min: "1.22"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/concurrency/reactor/
---

# Reactor

## Intent

The Reactor pattern demultiplexes events arriving on multiple sources and dispatches them synchronously to registered handler functions — all within a single goroutine (the event loop). This eliminates lock contention: because handlers run sequentially in the loop goroutine, they can read and write shared state without mutexes. It mirrors the epoll/kqueue/IOCP model at the OS level and the event-loop model of Node.js/Nginx at the application level.

## Context

Go's `select` statement is the natural implementation of an in-memory reactor: it monitors multiple channels and dispatches to the first ready one. A single goroutine running `select` in a loop is the Go idiom for a single-threaded event loop. This is already what the Go runtime itself does internally for timer and network I/O multiplexing.

Use the pattern when a subsystem must react to several independent event streams (timers, control commands, data events, shutdown) without the overhead of per-stream goroutines or the complexity of shared-state synchronisation. The trade-off is that **handlers must be short** — a long-running handler stalls the entire reactor.

## Implementation in Go

```go
type Event struct{ Kind string; Payload any }
type Handler func(Event)

type Reactor struct {
    events   chan Event
    handlers map[string]Handler
    done     chan struct{}
}

func New(buf int) *Reactor {
    return &Reactor{
        events:   make(chan Event, buf),
        handlers: make(map[string]Handler),
        done:     make(chan struct{}),
    }
}

func (r *Reactor) Register(kind string, h Handler) { r.handlers[kind] = h }
func (r *Reactor) Send(e Event)                    { r.events <- e }

func (r *Reactor) Run(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case e := <-r.events:
            if h, ok := r.handlers[e.Kind]; ok {
                h(e)
            }
        }
    }
}
```

For multiple distinct event channels (timers, data, control), add them as separate cases in the `select`. Each case maps to a dedicated handler.

## When to use

- In-process event loops where several independent event streams converge (timers, commands, data).
- Protocol state machines: events arrive from different sources (timeout, data received, peer closed); a single loop simplifies state transition logic.
- Actor-style isolated components: each actor is a goroutine with a `select` loop; the "reactor" pattern is the loop body.
- Subsystems that must process events in arrival order without reordering across streams.

## When NOT to use

- When handlers are CPU-intensive or block on I/O — offload to worker goroutines instead of blocking the loop.
- When events are rare and goroutine overhead is acceptable — a dedicated goroutine per event source is simpler.
- When you need true parallelism across event types — a reactor is serial by design.
- Network I/O at scale — Go's runtime already multiplexes with `net/http`; building a reactor on raw `net.Conn` re-implements what the runtime gives you.

## Gotchas

- **Blocking handler**: any handler that blocks (waits on a mutex, does network I/O) stalls all other events. Dispatch to a goroutine and return immediately if work is non-trivial.
- **Unbuffered event channel**: senders block when the reactor loop is busy. Use a buffered channel and monitor its length; a growing queue signals handler overload.
- **Handler registration after Run**: if `Register` is called after `Run` starts, there is a data race on the `handlers` map. Fix: register before starting, or use a `sync.Map`, or send a special registration event.
- **Shutdown ordering**: after `ctx` cancels, events still in the buffer are silently dropped. Drain the channel if at-least-once delivery is required.

## See also

- skills/concurrency/producer-consumer.md
- skills/concurrency/broadcast.md
- examples/concurrency/reactor/
