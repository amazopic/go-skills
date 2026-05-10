package broadcast

import (
	"sync"
	"sync/atomic"
	"testing"
)

// --------------------------------------------------------------------------
// Signal tests
// --------------------------------------------------------------------------

func TestSignal_BroadcastWakesAll(t *testing.T) {
	const n = 20
	sig := NewSignal()
	var ready sync.WaitGroup
	var woken atomic.Int32

	ready.Add(n)
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			ready.Done()
			<-sig.Wait()
			woken.Add(1)
		}()
	}
	ready.Wait() // ensure all goroutines are blocked before firing
	sig.Broadcast()
	wg.Wait()

	if got := woken.Load(); got != n {
		t.Errorf("woken = %d, want %d", got, n)
	}
}

func TestSignal_MultipleCallsAreSafe(t *testing.T) {
	sig := NewSignal()
	// Calling Broadcast multiple times must not panic.
	for range 10 {
		sig.Broadcast()
	}
	<-sig.Wait() // must be immediately readable
}

func TestSignal_WaitAfterBroadcast(t *testing.T) {
	sig := NewSignal()
	sig.Broadcast()
	// Late subscriber should not block.
	select {
	case <-sig.Wait():
	default:
		t.Error("Wait() blocked after Broadcast")
	}
}

// --------------------------------------------------------------------------
// Broadcaster tests
// --------------------------------------------------------------------------

func TestBroadcaster_DeliverToAll(t *testing.T) {
	var b Broadcaster[int]
	const (
		subscribers = 5
		messages    = 4
	)

	subs := make([]<-chan int, subscribers)
	for i := range subscribers {
		subs[i] = b.Subscribe(messages) // buffer = number of messages so no drops
	}

	for i := range messages {
		b.Send(i)
	}
	b.Close()

	for s, ch := range subs {
		got := make([]int, 0, messages)
		for v := range ch {
			got = append(got, v)
		}
		if len(got) != messages {
			t.Errorf("subscriber %d: received %d messages, want %d", s, len(got), messages)
		}
		for i, v := range got {
			if v != i {
				t.Errorf("subscriber %d: msg[%d] = %d, want %d", s, i, v, i)
			}
		}
	}
}

func TestBroadcaster_Unsubscribe(t *testing.T) {
	var b Broadcaster[string]
	ch1 := b.Subscribe(10)
	ch2 := b.Subscribe(10)

	b.Send("hello")
	b.Unsubscribe(ch2) // ch2 is closed here
	b.Send("world")

	// ch1 should have both messages
	if v := <-ch1; v != "hello" {
		t.Errorf("ch1[0] = %q, want %q", v, "hello")
	}
	if v := <-ch1; v != "world" {
		t.Errorf("ch1[1] = %q, want %q", v, "world")
	}

	// ch2 should have exactly one message then be closed
	if v := <-ch2; v != "hello" {
		t.Errorf("ch2[0] = %q, want %q", v, "hello")
	}
	if _, ok := <-ch2; ok {
		t.Error("ch2 should be closed after Unsubscribe")
	}
	b.Close()
}

func TestBroadcaster_Race(t *testing.T) {
	var b Broadcaster[int]
	const subs = 8
	channels := make([]<-chan int, subs)
	for i := range subs {
		channels[i] = b.Subscribe(100)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 50 {
			b.Send(i)
		}
		b.Close()
	}()

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for range c {
			}
		}(ch)
	}
	wg.Wait()
}
