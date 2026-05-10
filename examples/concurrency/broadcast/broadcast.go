// Package broadcast demonstrates two idiomatic Go broadcast mechanisms:
//  1. Signal — a one-shot broadcast via channel close (zero allocation per waiter).
//  2. Broadcaster — a repeatable typed broadcast to a dynamic set of subscribers.
package broadcast

import "sync"

// --------------------------------------------------------------------------
// Signal: one-shot broadcast (e.g., shutdown, ready)
// --------------------------------------------------------------------------

// Signal is a one-shot broadcast primitive. Closing the underlying channel
// wakes every goroutine blocked on Wait() simultaneously.
type Signal struct {
	once sync.Once
	ch   chan struct{}
}

// NewSignal returns an armed Signal.
func NewSignal() *Signal {
	return &Signal{ch: make(chan struct{})}
}

// Broadcast fires the signal. Subsequent calls are no-ops.
func (s *Signal) Broadcast() {
	s.once.Do(func() { close(s.ch) })
}

// Wait returns a channel that is closed when Broadcast is called.
// Use in a select to support context cancellation:
//
//	select { case <-sig.Wait(): ... case <-ctx.Done(): ... }
func (s *Signal) Wait() <-chan struct{} {
	return s.ch
}

// --------------------------------------------------------------------------
// Broadcaster: repeatable typed broadcast to dynamic subscribers
// --------------------------------------------------------------------------

// Broadcaster fans a value out to every current subscriber.
// Subscribe before the first Send; late subscribers miss earlier messages.
type Broadcaster[T any] struct {
	mu   sync.Mutex
	subs []chan T
}

// Subscribe registers a new receiver with the given buffer size and returns
// its receive channel. Call Unsubscribe when the receiver exits.
func (b *Broadcaster[T]) Subscribe(buf int) <-chan T {
	ch := make(chan T, buf)
	b.mu.Lock()
	b.subs = append(b.subs, ch)
	b.mu.Unlock()
	return ch
}

// Unsubscribe removes the channel and closes it so the receiver goroutine
// can detect the end-of-stream.
func (b *Broadcaster[T]) Unsubscribe(sub <-chan T) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, ch := range b.subs {
		if ch == sub {
			b.subs = append(b.subs[:i], b.subs[i+1:]...)
			close(ch)
			return
		}
	}
}

// Send delivers v to every subscriber. If a subscriber's buffer is full the
// message is dropped for that subscriber (non-blocking policy). Adjust the
// select to block if guaranteed delivery is required.
func (b *Broadcaster[T]) Send(v T) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.subs {
		select {
		case ch <- v:
		default: // subscriber too slow; drop rather than block the broadcaster
		}
	}
}

// Close shuts down the broadcaster by closing all subscriber channels.
func (b *Broadcaster[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.subs {
		close(ch)
	}
	b.subs = nil
}
