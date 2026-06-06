---
name: messaging-publish-subscribe
description: Publish/Subscribe — broker-mediated event distribution: publishers emit topics, subscribers receive matching messages. Use for decoupled cross-component communication.
category: messaging
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/messaging/publish_subscribe.md
example: examples/messaging/publish-subscribe/
---

Publish & Subscribe Messaging Pattern
============
Publish-Subscribe is a messaging pattern used to communicate messages between 
different components without these components knowing anything about each other's identity.

It is similar to the Observer behavioral design pattern. 
The fundamental design principals of both Observer and Publish-Subscribe is the decoupling of
those interested in being informed about `Event Messages` from the informer (Observers or Publishers).
Meaning that you don't have to program the messages to be sent directly to specific receivers.

To accomplish this, an intermediary, called a "message broker" or "event bus", 
receives published messages, and then routes them on to subscribers.


The three roles are **publishers**, the **broker** (event bus), and **subscribers**.
In Go the idiomatic shape is a generic broker that maps each topic to a set of
per-subscriber channels. Publishing fans a copy of the message out to every
current subscriber of the topic (broadcast, not work-sharing).

```go
// Broker is a topic-based pub/sub bus, safe for concurrent use.
type Broker[T any] struct {
	mu     sync.RWMutex
	topics map[string]map[*subscription[T]]struct{}
	closed bool
}

type subscription[T any] struct {
	ch   chan T
	done chan struct{} // closed once on teardown
}

// Subscribe registers a subscriber and returns a receive-only stream.
func (b *Broker[T]) Subscribe(ctx context.Context, topic string, buf int) (*Subscription[T], error) {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return nil, ErrClosed
	}
	sub := &subscription[T]{ch: make(chan T, buf), done: make(chan struct{})}
	if b.topics[topic] == nil {
		b.topics[topic] = make(map[*subscription[T]]struct{})
	}
	b.topics[topic][sub] = struct{}{}
	b.mu.Unlock()

	s := &Subscription[T]{C: sub.ch, broker: b, topic: topic, sub: sub}
	// Tie lifetime to ctx without leaking a goroutine and without stealing
	// messages: wait on done (closed by teardown), never on sub.ch.
	go func() {
		select {
		case <-ctx.Done():
			s.Unsubscribe()
		case <-sub.done:
		}
	}()
	return s, nil
}
```

```go
// Publish broadcasts msg to every current subscriber. It never blocks: a
// subscriber whose buffer is full has its copy dropped (straggler isolation).
func (b *Broker[T]) Publish(ctx context.Context, topic string, msg T) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return 0, ErrClosed
	}
	// Snapshot subscribers, then send outside the lock.
	targets := make([]*subscription[T], 0, len(b.topics[topic]))
	for sub := range b.topics[topic] {
		targets = append(targets, sub)
	}
	b.mu.RUnlock()

	delivered := 0
	for _, sub := range targets {
		select {
		case <-ctx.Done():
			return delivered, ctx.Err()
		case sub.ch <- msg:
			delivered++
		default: // full buffer: drop for this subscriber, keep the bus flowing
		}
	}
	return delivered, nil
}
```

Gotchas
============
- **Do not block publishers on slow subscribers.** Give each subscriber its own
  buffered channel and use a `default` (drop) or a deadline; never let one
  straggler stall the whole bus.
- **Close every subscriber channel exactly once.** Guard teardown with a
  membership check (and `sync.Once` on the handle) so `Unsubscribe`, ctx
  cancellation, and `Broker.Close` cannot double-close.
- **Send outside the lock.** Snapshot the subscriber set under the lock, then do
  channel sends without holding it, so a slow send never blocks other topics.
- **The ctx lifetime goroutine must not receive on the message channel** — that
  would steal messages from the real consumer. Signal teardown on a separate
  `done` channel.

Improvements
============
Events can be published in parallel by fanning out to subscribers in separate
goroutines. For broadcast (every subscriber sees every message) keep a
per-subscriber channel; for work-sharing (each message handled once) use the
push-pull pattern instead.
