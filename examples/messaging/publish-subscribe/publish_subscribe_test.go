package pubsub

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

// drainExpect receives exactly want messages from s.C, returning them. It fails
// the test if the channel closes early.
func drainExpect[T any](t *testing.T, s *Subscription[T], want int) []T {
	t.Helper()
	got := make([]T, 0, want)
	for range want {
		v, ok := <-s.C
		if !ok {
			t.Fatalf("subscription channel closed after %d/%d messages", len(got), want)
		}
		got = append(got, v)
	}
	return got
}

func TestBroker_BroadcastToAllSubscribers(t *testing.T) {
	tests := []struct {
		name        string
		subscribers int
		messages    int
	}{
		{"single subscriber", 1, 5},
		{"three subscribers", 3, 10},
		{"many subscribers", 16, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New[int]()
			defer b.Close()
			ctx := context.Background()

			subs := make([]*Subscription[int], tt.subscribers)
			for i := range subs {
				s, err := b.Subscribe(ctx, "topic", tt.messages)
				if err != nil {
					t.Fatalf("Subscribe: %v", err)
				}
				subs[i] = s
			}

			// Every subscriber must receive every message, in order.
			var wg sync.WaitGroup
			wg.Add(tt.subscribers)
			for _, s := range subs {
				go func(s *Subscription[int]) {
					defer wg.Done()
					got := drainExpect(t, s, tt.messages)
					for j, v := range got {
						if v != j {
							t.Errorf("message %d = %d, want %d", j, v, j)
						}
					}
				}(s)
			}

			for j := range tt.messages {
				n, err := b.Publish(ctx, "topic", j)
				if err != nil {
					t.Fatalf("Publish: %v", err)
				}
				if n != tt.subscribers {
					t.Errorf("Publish delivered to %d subscribers, want %d", n, tt.subscribers)
				}
			}
			wg.Wait()

			if d := b.Dropped(); d != 0 {
				t.Errorf("Dropped = %d, want 0 (buffers were large enough)", d)
			}
		})
	}
}

func TestBroker_TopicIsolation(t *testing.T) {
	b := New[string]()
	defer b.Close()
	ctx := context.Background()

	a, err := b.Subscribe(ctx, "a", 1)
	if err != nil {
		t.Fatal(err)
	}
	c, err := b.Subscribe(ctx, "c", 1)
	if err != nil {
		t.Fatal(err)
	}

	if n, _ := b.Publish(ctx, "a", "hello"); n != 1 {
		t.Errorf("publish to a delivered to %d, want 1", n)
	}
	if n, _ := b.Publish(ctx, "c", "world"); n != 1 {
		t.Errorf("publish to c delivered to %d, want 1", n)
	}

	if got := <-a.C; got != "hello" {
		t.Errorf("a received %q, want %q", got, "hello")
	}
	if got := <-c.C; got != "world" {
		t.Errorf("c received %q, want %q", got, "world")
	}

	// Publishing to an unknown topic delivers to nobody.
	if n, _ := b.Publish(ctx, "nope", "x"); n != 0 {
		t.Errorf("publish to empty topic delivered to %d, want 0", n)
	}
}

func TestBroker_Unsubscribe(t *testing.T) {
	b := New[int]()
	defer b.Close()
	ctx := context.Background()

	s, err := b.Subscribe(ctx, "t", 1)
	if err != nil {
		t.Fatal(err)
	}
	if got := b.NumSubscribers("t"); got != 1 {
		t.Fatalf("NumSubscribers = %d, want 1", got)
	}

	s.Unsubscribe()
	// Idempotent: a second call must not panic or double-close.
	s.Unsubscribe()

	if got := b.NumSubscribers("t"); got != 0 {
		t.Errorf("NumSubscribers after Unsubscribe = %d, want 0", got)
	}
	// Channel must be closed.
	if _, ok := <-s.C; ok {
		t.Error("expected closed channel after Unsubscribe")
	}
	// Publishing now reaches nobody.
	if n, _ := b.Publish(ctx, "t", 1); n != 0 {
		t.Errorf("publish after unsubscribe delivered to %d, want 0", n)
	}
}

func TestBroker_ContextCancelUnsubscribes(t *testing.T) {
	b := New[int]()
	defer b.Close()

	ctx, cancel := context.WithCancel(context.Background())
	s, err := b.Subscribe(ctx, "t", 1)
	if err != nil {
		t.Fatal(err)
	}

	cancel()

	// The subscriber channel must eventually close as a result of cancellation.
	// Ranging blocks until close — deterministic, no sleep needed.
	for range s.C {
		// drain any in-flight (none expected)
	}
	if got := b.NumSubscribers("t"); got != 0 {
		t.Errorf("NumSubscribers after cancel = %d, want 0", got)
	}
}

func TestBroker_SlowConsumerDrops(t *testing.T) {
	b := New[int]()
	defer b.Close()
	ctx := context.Background()

	// Buffer of 1, never drained: the 2nd+ publishes drop for this subscriber.
	s, err := b.Subscribe(ctx, "t", 1)
	if err != nil {
		t.Fatal(err)
	}
	_ = s

	n0, _ := b.Publish(ctx, "t", 0) // fills buffer
	n1, _ := b.Publish(ctx, "t", 1) // dropped
	n2, _ := b.Publish(ctx, "t", 2) // dropped

	if n0 != 1 {
		t.Errorf("first publish delivered to %d, want 1", n0)
	}
	if n1 != 0 || n2 != 0 {
		t.Errorf("subsequent publishes delivered to %d/%d, want 0/0", n1, n2)
	}
	if d := b.Dropped(); d != 2 {
		t.Errorf("Dropped = %d, want 2", d)
	}
}

func TestBroker_ClosedRejects(t *testing.T) {
	b := New[int]()
	ctx := context.Background()

	s, err := b.Subscribe(ctx, "t", 1)
	if err != nil {
		t.Fatal(err)
	}

	b.Close()
	// Idempotent.
	b.Close()

	// Existing subscription channel is closed.
	if _, ok := <-s.C; ok {
		t.Error("expected closed channel after Broker.Close")
	}

	if _, err := b.Subscribe(ctx, "t", 1); !errors.Is(err, ErrClosed) {
		t.Errorf("Subscribe after Close: err = %v, want ErrClosed", err)
	}
	if _, err := b.Publish(ctx, "t", 1); !errors.Is(err, ErrClosed) {
		t.Errorf("Publish after Close: err = %v, want ErrClosed", err)
	}
}

func TestBroker_PublishRespectsContext(t *testing.T) {
	b := New[int]()
	defer b.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := b.Publish(ctx, "t", 1); !errors.Is(err, context.Canceled) {
		t.Errorf("Publish with cancelled ctx: err = %v, want context.Canceled", err)
	}
}

// TestBroker_ConcurrentPublishersAndSubscribers exercises the broker under load
// with many concurrent publishers and subscribers. It asserts each subscriber
// receives every message exactly once (buffers are sized to never drop).
func TestBroker_ConcurrentPublishersAndSubscribers(t *testing.T) {
	const (
		publishers     = 8
		perPublisher   = 200
		subscribers    = 6
		totalPublished = publishers * perPublisher
	)

	b := New[int]()
	defer b.Close()
	ctx := context.Background()

	// Counters: total messages each subscriber received.
	var received [subscribers]atomic.Int64

	// Sum check: each subscriber should see the sum of all published values.
	var subSums [subscribers]atomic.Int64

	subs := make([]*Subscription[int], subscribers)
	for i := range subs {
		s, err := b.Subscribe(ctx, "load", totalPublished)
		if err != nil {
			t.Fatal(err)
		}
		subs[i] = s
	}

	var consumeWg sync.WaitGroup
	consumeWg.Add(subscribers)
	for i, s := range subs {
		go func(i int, s *Subscription[int]) {
			defer consumeWg.Done()
			for v := range s.C {
				received[i].Add(1)
				subSums[i].Add(int64(v))
			}
		}(i, s)
	}

	// Publishers each send a disjoint contiguous range of values.
	var pubWg sync.WaitGroup
	pubWg.Add(publishers)
	var wantSum int64
	for v := 0; v < totalPublished; v++ {
		wantSum += int64(v)
	}
	for p := range publishers {
		go func(p int) {
			defer pubWg.Done()
			for j := range perPublisher {
				val := p*perPublisher + j
				if _, err := b.Publish(ctx, "load", val); err != nil {
					t.Errorf("Publish: %v", err)
					return
				}
			}
		}(p)
	}
	pubWg.Wait()

	// All published; tear down subscriptions so consumer goroutines exit.
	for _, s := range subs {
		s.Unsubscribe()
	}
	consumeWg.Wait()

	for i := range subscribers {
		if got := received[i].Load(); got != totalPublished {
			t.Errorf("subscriber %d received %d, want %d", i, got, totalPublished)
		}
		if got := subSums[i].Load(); got != wantSum {
			t.Errorf("subscriber %d sum = %d, want %d", i, got, wantSum)
		}
	}
	if d := b.Dropped(); d != 0 {
		t.Errorf("Dropped = %d, want 0", d)
	}
}
