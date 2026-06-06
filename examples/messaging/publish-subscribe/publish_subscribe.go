// Package pubsub implements the Publish/Subscribe messaging pattern: a broker
// mediates between publishers and subscribers so neither knows the other's
// identity. Publishers send messages to a topic; the broker fans each message
// out to every current subscriber of that topic (broadcast, not work-sharing).
//
// The Broker is generic over the message type and safe for concurrent use by
// any number of publishers and subscribers.
//
// Delivery semantics:
//   - Each subscriber owns an independent buffered channel.
//   - Publish is non-blocking per subscriber: if a subscriber's buffer is full,
//     that subscriber's message is dropped (slow-consumer isolation) and counted,
//     rather than blocking the publisher or other subscribers. This prevents one
//     straggler from stalling the whole bus, as the skill recommends.
//   - Subscribers must Unsubscribe (or cancel their context) to release resources.
package pubsub

import (
	"context"
	"errors"
	"sync"
)

// ErrClosed is returned by Publish and Subscribe after the broker is closed.
var ErrClosed = errors.New("pubsub: broker closed")

// subscription is the broker's internal view of one subscriber on one topic.
type subscription[T any] struct {
	ch   chan T        // message stream delivered to the subscriber
	done chan struct{} // closed exactly once on teardown
}

// Broker is a topic-based publish/subscribe message bus. The zero value is not
// usable; create one with New. A Broker is safe for concurrent use.
type Broker[T any] struct {
	mu     sync.RWMutex
	topics map[string]map[*subscription[T]]struct{}
	closed bool

	// droppedMu guards dropped independently of mu so concurrent Publishers,
	// which fan out under mu's read lock, can record drops without escalating
	// to mu's write lock (which would deadlock against their own RLock).
	droppedMu sync.Mutex
	dropped   int // messages dropped due to full subscriber buffers
}

// New creates an empty Broker.
func New[T any]() *Broker[T] {
	return &Broker[T]{
		topics: make(map[string]map[*subscription[T]]struct{}),
	}
}

// Subscription is the consumer-facing handle returned by Subscribe. Receive
// messages by ranging over C; the channel is closed when the subscription ends
// (via Unsubscribe, context cancellation, or Broker.Close).
type Subscription[T any] struct {
	C <-chan T // receive-only stream of messages for the subscribed topic

	broker *Broker[T]
	topic  string
	sub    *subscription[T]

	once sync.Once // guards exactly-once teardown
}

// Subscribe registers a new subscriber on topic and returns its handle. The
// underlying channel is buffered with capacity buf; buf may be 0, in which case
// every Publish to this subscriber is dropped unless a receiver is mid-send-ready
// — for broadcast use a buf >= 1.
//
// If ctx is cancelled, the subscription is automatically torn down and C is
// closed. Subscribe returns ErrClosed if the broker is already closed.
func (b *Broker[T]) Subscribe(ctx context.Context, topic string, buf int) (*Subscription[T], error) {
	if buf < 0 {
		buf = 0
	}

	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return nil, ErrClosed
	}
	sub := &subscription[T]{ch: make(chan T, buf), done: make(chan struct{})}
	subs := b.topics[topic]
	if subs == nil {
		subs = make(map[*subscription[T]]struct{})
		b.topics[topic] = subs
	}
	subs[sub] = struct{}{}
	b.mu.Unlock()

	s := &Subscription[T]{
		C:      sub.ch,
		broker: b,
		topic:  topic,
		sub:    sub,
	}

	// Tie the subscription lifetime to ctx without leaking a goroutine and
	// without stealing messages from the subscriber: the goroutine waits on
	// either ctx cancellation or teardown (sub.done, closed by remove/Close).
	go func() {
		select {
		case <-ctx.Done():
			s.Unsubscribe()
		case <-sub.done:
			// Torn down by Unsubscribe or Close; nothing more to do.
		}
	}()

	return s, nil
}

// Unsubscribe removes the subscription from its topic and closes C. It is
// idempotent and safe to call concurrently. Calling it more than once is a
// no-op.
func (s *Subscription[T]) Unsubscribe() {
	s.once.Do(func() {
		s.broker.remove(s.topic, s.sub)
	})
}

// remove detaches sub from topic and closes its channel exactly once.
func (b *Broker[T]) remove(topic string, sub *subscription[T]) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subs := b.topics[topic]
	if subs == nil {
		return
	}
	if _, ok := subs[sub]; !ok {
		return
	}
	delete(subs, sub)
	if len(subs) == 0 {
		delete(b.topics, topic)
	}
	close(sub.done)
	close(sub.ch)
}

// Publish delivers msg to every current subscriber of topic. It never blocks:
// a subscriber whose buffer is full has this message dropped (and counted via
// Dropped). It returns the number of subscribers that received the message and
// ErrClosed if the broker is closed.
//
// ctx is honored for early cancellation while fanning out; on cancellation it
// returns the subscribers reached so far and ctx.Err().
//
// Concurrency note: teardown (remove/Close) closes a subscriber's channel only
// while holding b.mu for writing, and Publish performs every send while holding
// b.mu for reading. Because RLock and Lock are mutually exclusive, a channel can
// never be closed mid-send: this closes the send-on-closed-channel TOCTOU that
// would otherwise let a concurrent Unsubscribe/Close panic an in-flight Publish.
// The send stays non-blocking (select with default), so holding the read lock
// during fan-out never stalls — it only excludes teardown, not other Publishers.
func (b *Broker[T]) Publish(ctx context.Context, topic string, msg T) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return 0, ErrClosed
	}

	delivered := 0
	dropped := 0
	for sub := range b.topics[topic] {
		select {
		case <-ctx.Done():
			b.addDropped(dropped)
			return delivered, ctx.Err()
		case sub.ch <- msg:
			// Safe: the channel cannot be closed here because teardown needs
			// the write lock, which this read lock excludes.
			delivered++
		default:
			// Full buffer: drop for this subscriber, keep the bus flowing.
			dropped++
		}
	}
	b.addDropped(dropped)
	return delivered, nil
}

// addDropped accumulates the per-Publish drop count into the broker total via a
// dedicated mutex, so it is safe to call while Publish still holds b.mu's read
// lock during fan-out (it never touches b.mu).
func (b *Broker[T]) addDropped(n int) {
	if n == 0 {
		return
	}
	b.droppedMu.Lock()
	b.dropped += n
	b.droppedMu.Unlock()
}

// Dropped returns the cumulative number of messages dropped because a
// subscriber's buffer was full at publish time.
func (b *Broker[T]) Dropped() int {
	b.droppedMu.Lock()
	defer b.droppedMu.Unlock()
	return b.dropped
}

// NumSubscribers returns the current subscriber count for topic.
func (b *Broker[T]) NumSubscribers(topic string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.topics[topic])
}

// Close shuts the broker down: it closes every subscriber channel on every
// topic and rejects subsequent Subscribe/Publish calls with ErrClosed. Close is
// idempotent.
func (b *Broker[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for topic, subs := range b.topics {
		for sub := range subs {
			close(sub.done)
			close(sub.ch)
		}
		delete(b.topics, topic)
	}
}
