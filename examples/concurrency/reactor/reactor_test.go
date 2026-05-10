package reactor

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestReactor_DispatchToHandler(t *testing.T) {
	r := New(16)
	var count atomic.Int32
	r.Register("inc", func(e Event) { count.Add(1) })

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); r.Run(ctx) }()

	const n = 10
	for range n {
		r.Send(Event{Kind: "inc"})
	}

	// Give the loop time to process before cancelling.
	for count.Load() < n {
		time.Sleep(time.Millisecond)
	}
	cancel()
	wg.Wait()

	if got := count.Load(); got != n {
		t.Errorf("count = %d, want %d", got, n)
	}
}

func TestReactor_UnknownKindNoOp(t *testing.T) {
	r := New(4)
	// No handler registered — sending unknown kind must not panic.
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); r.Run(ctx) }()

	r.Send(Event{Kind: "unknown"})
	cancel()
	wg.Wait()
}

func TestReactor_MultipleKinds(t *testing.T) {
	r := New(32)
	var aCount, bCount atomic.Int32
	r.Register("a", func(e Event) { aCount.Add(1) })
	r.Register("b", func(e Event) { bCount.Add(1) })

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); r.Run(ctx) }()

	for range 5 {
		r.Send(Event{Kind: "a"})
		r.Send(Event{Kind: "b"})
	}
	for aCount.Load() < 5 || bCount.Load() < 5 {
		time.Sleep(time.Millisecond)
	}
	cancel()
	wg.Wait()

	if aCount.Load() != 5 {
		t.Errorf("aCount = %d, want 5", aCount.Load())
	}
	if bCount.Load() != 5 {
		t.Errorf("bCount = %d, want 5", bCount.Load())
	}
}

func TestReactor_PayloadDelivered(t *testing.T) {
	r := New(4)
	got := make(chan any, 1)
	r.Register("data", func(e Event) { got <- e.Payload })

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); r.Run(ctx) }()

	r.Send(Event{Kind: "data", Payload: 42})
	v := <-got
	if v != 42 {
		t.Errorf("payload = %v, want 42", v)
	}
	cancel()
	wg.Wait()
}

func TestReactor_TrySendDropsWhenFull(t *testing.T) {
	r := New(1) // tiny buffer
	// Fill it.
	r.events <- Event{Kind: "x"}
	if r.TrySend(Event{Kind: "x"}) {
		t.Error("TrySend should return false when buffer is full")
	}
}

func TestReactor_DrainOnCancel(t *testing.T) {
	// Events queued before cancel should still be dispatched.
	r := New(32)
	var count atomic.Int32
	r.Register("inc", func(e Event) { count.Add(1) })

	ctx, cancel := context.WithCancel(context.Background())
	// Queue events before starting the loop.
	const n = 20
	for range n {
		r.Send(Event{Kind: "inc"})
	}
	cancel() // cancel immediately

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); r.Run(ctx) }()
	wg.Wait()

	// All queued events must have been dispatched (drain on cancel).
	if got := count.Load(); got != n {
		t.Errorf("count after drain = %d, want %d", got, n)
	}
}
