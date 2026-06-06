package fanout

import (
	"context"
	"errors"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
)

// feed returns a channel pre-loaded with 0..n-1 and already closed, so workers
// drain it deterministically without any producer goroutine racing the test.
func feed(n int) <-chan int {
	in := make(chan int, n)
	for i := range n {
		in <- i
	}
	close(in)
	return in
}

func TestFanOut_AllItemsProcessedExactlyOnce(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		items   int
	}{
		{"1worker", 1, 100},
		{"4workers", 4, 1000},
		{"more_workers_than_items", 16, 5},
		{"single_item", 8, 1},
		{"zero_items", 4, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			in := feed(tt.items)

			// process doubles the input so we can verify the mapping, not just
			// the count.
			out := FanOut(ctx, tt.workers, in, func(v int) int { return v * 2 })

			got, err := Collect(ctx, out)
			if err != nil {
				t.Fatalf("Collect: %v", err)
			}
			if len(got) != tt.items {
				t.Fatalf("got %d results, want %d", len(got), tt.items)
			}

			sort.Ints(got)
			for i, v := range got {
				if v != i*2 {
					t.Errorf("got[%d] = %d, want %d", i, v, i*2)
				}
			}
		})
	}
}

// TestFanOut_EachItemHandledOnce verifies the work-distribution property: every
// source item is consumed by exactly one worker (no drops, no duplicates),
// regardless of how the runtime schedules the workers.
func TestFanOut_EachItemHandledOnce(t *testing.T) {
	const (
		items   = 2000
		workers = 12
	)
	ctx := context.Background()
	in := feed(items)

	var counts [items]atomic.Int32
	out := FanOut(ctx, workers, in, func(v int) int {
		counts[v].Add(1)
		return v
	})

	if _, err := Collect(ctx, out); err != nil {
		t.Fatalf("Collect: %v", err)
	}
	for i := range items {
		if got := counts[i].Load(); got != 1 {
			t.Errorf("item %d processed %d times, want 1", i, got)
		}
	}
}

// TestFanOut_WorkParallelises confirms multiple workers run concurrently: each
// invocation of process gates on a barrier that only releases once `workers`
// callers have arrived, so the test deadlocks (and fails) unless at least that
// many workers are truly in-flight at the same time. Synchronisation is via a
// WaitGroup barrier — no sleeps, no timing assumptions.
func TestFanOut_WorkParallelises(t *testing.T) {
	const (
		workers = 8
		items   = workers * 50
	)
	ctx := context.Background()
	in := feed(items)

	// barrier releases after `workers` simultaneous arrivals, proving that many
	// workers are active at once. We only need the barrier to fire once.
	var (
		mu       sync.Mutex
		waiting  int
		released bool
		cond     = sync.NewCond(&mu)
		maxSeen  atomic.Int64
	)

	out := FanOut(ctx, workers, in, func(v int) int {
		mu.Lock()
		waiting++
		if int64(waiting) > maxSeen.Load() {
			maxSeen.Store(int64(waiting))
		}
		if waiting >= workers {
			released = true
			cond.Broadcast()
		}
		// Block until the barrier has fired at least once; thereafter pass
		// straight through so the remaining items drain.
		for !released {
			cond.Wait()
		}
		waiting--
		mu.Unlock()
		return v
	})

	got, err := Collect(ctx, out)
	if err != nil {
		t.Fatalf("Collect: %v", err)
	}
	if len(got) != items {
		t.Fatalf("got %d results, want %d", len(got), items)
	}
	if maxSeen.Load() < workers {
		t.Errorf("max concurrent workers = %d, want %d", maxSeen.Load(), workers)
	}
}

// TestFanOut_ContextCancellationStopsWorkers ensures that cancelling the context
// unblocks workers stuck on send and closes the output channel — proving there
// is no goroutine leak when the consumer abandons the stream.
func TestFanOut_ContextCancellationStopsWorkers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Unbuffered, never-closed source: workers would block forever on receive
	// without context support.
	in := make(chan int)

	out := FanOut(ctx, 4, in, func(v int) int { return v })

	cancel()

	// out must become closed; if any worker leaked it would never close.
	for range out {
		// Drain any in-flight value (there should be none, but be tolerant).
	}
	// Reaching here means the closer goroutine ran => all workers exited.
}

// TestFanOut_CancelMidStreamUnblocksSend cancels while a worker is blocked
// trying to send because the consumer has stopped reading. The output must
// still close.
func TestFanOut_CancelMidStreamUnblocksSend(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	in := make(chan int)
	// Producer pushes a few items then leaves the rest; we never close in.
	go func() {
		for i := 0; i < 3; i++ {
			select {
			case in <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	out := FanOut(ctx, 2, in, func(v int) int { return v })

	// Read exactly one value, then cancel and stop reading. A worker may be
	// parked on `out <- ...`; cancellation must release it.
	<-out
	cancel()

	// Drain to closure; must terminate.
	for range out {
	}
}

func TestFanOut_PanicsOnZeroWorkers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for n=0")
		}
	}()
	FanOut(context.Background(), 0, feed(1), func(v int) int { return v })
}

func TestFanOutErr_InvalidWorkers(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		wantErr bool
	}{
		{"zero", 0, true},
		{"negative", -3, true},
		{"valid", 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch, err := FanOutErr(context.Background(), tt.workers, feed(4), func(v int) int { return v })
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidWorkers) {
					t.Fatalf("err = %v, want ErrInvalidWorkers", err)
				}
				if ch != nil {
					t.Error("expected nil channel on error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if _, err := Collect(context.Background(), ch); err != nil {
				t.Fatalf("Collect: %v", err)
			}
		})
	}
}

func TestCollect_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// A live, never-closed channel: Collect must return promptly via ctx.
	ch := make(chan int)
	_, err := Collect(ctx, ch)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled (wrapped)", err)
	}
}

func TestCollect_ReturnsPartialOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan int)

	// Deliver two items, then cancel; Collect should return what it has.
	done := make(chan struct{})
	go func() {
		defer close(done)
		ch <- 1
		ch <- 2
		cancel()
	}()

	got, err := Collect(ctx, ch)
	<-done
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Errorf("partial = %v, want [1 2]", got)
	}
}
