# Publish / Subscribe

`Broker[T]` — a topic-based message bus that decouples publishers from
subscribers. Publishers send a message to a topic; the broker **broadcasts** a
copy to every current subscriber of that topic. Neither side knows the other's
identity.

This is broadcast (fan-out), not work-sharing: contrast with
[`push-pull`](../push-pull/), where each item goes to exactly *one* consumer.

## When to use

- Decoupled, event-driven communication between components.
- One event must reach many independent consumers (cache invalidation, UI
  updates, audit/metrics taps).
- You want slow consumers isolated so they cannot stall publishers or peers.

## API

```go
b := pubsub.New[Event]()
defer b.Close()

// Subscriber: buffered channel; ranges until the subscription ends.
sub, _ := b.Subscribe(ctx, "orders", 64)
go func() {
    for ev := range sub.C { handle(ev) }
}()

// Publisher: non-blocking fan-out; returns how many subscribers got it.
n, _ := b.Publish(ctx, "orders", ev)

// Teardown: Unsubscribe, cancel ctx, or Broker.Close (closes every sub).
sub.Unsubscribe()
```

## Key properties

- **Broadcast:** every subscriber of a topic receives every message.
- **Slow-consumer isolation:** `Publish` never blocks. If a subscriber's buffer
  is full its copy is dropped (counted by `Dropped()`); other subscribers and
  the publisher are unaffected. Size the buffer for your burst tolerance.
- **No goroutine leaks:** context cancellation tears the subscription down and
  closes `sub.C`; the lifetime goroutine waits on a private `done` channel so it
  never steals messages from the consumer.
- **Concurrency-safe:** any number of publishers and subscribers. `Publish`
  fans out under the broker's *read* lock and teardown (`Unsubscribe`/`Close`)
  closes channels under the *write* lock, so the two are mutually exclusive — a
  channel can never be closed mid-send (no "send on closed channel" panic). The
  send is non-blocking (`select`/`default`), so holding the read lock during
  fan-out never stalls and still admits concurrent publishers.

## Run

```bash
cd examples && go test -race ./messaging/publish-subscribe/
```
