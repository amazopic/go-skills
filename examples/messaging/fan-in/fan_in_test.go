package fanin

import (
	"context"
	"errors"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
)

// makeChan returns a channel that emits vs in order and is then closed.
func makeChan(vs ...int) <-chan int {
	c := make(chan int, len(vs))
	for _, v := range vs {
		c <- v
	}
	close(c)
	return c
}

func TestMerge_AllValuesDelivered(t *testing.T) {
	tests := []struct {
		name   string
		inputs [][]int
		want   []int // sorted; Merge does not preserve cross-input order
	}{
		{
			name:   "single input",
			inputs: [][]int{{1, 2, 3}},
			want:   []int{1, 2, 3},
		},
		{
			name:   "two inputs",
			inputs: [][]int{{1, 2, 3}, {4, 5, 6}},
			want:   []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:   "many inputs uneven lengths",
			inputs: [][]int{{1}, {2, 3, 4, 5}, {}, {6, 7}},
			want:   []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:   "all empty inputs",
			inputs: [][]int{{}, {}, {}},
			want:   nil,
		},
		{
			name:   "no inputs",
			inputs: nil,
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := make([]<-chan int, 0, len(tt.inputs))
			for _, in := range tt.inputs {
				cs = append(cs, makeChan(in...))
			}

			out := Merge(context.Background(), cs...)

			var got []int
			for v := range out {
				got = append(got, v)
			}
			sort.Ints(got)

			if len(got) != len(tt.want) {
				t.Fatalf("got %v (len %d), want %v (len %d)", got, len(got), tt.want, len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("got[%d] = %d, want %d (full got=%v want=%v)", i, got[i], tt.want[i], got, tt.want)
				}
			}
		})
	}
}

// TestMerge_OutputClosed verifies the merged channel is closed (range exits)
// for both the no-input and drained-input cases.
func TestMerge_OutputClosed(t *testing.T) {
	t.Run("no inputs closes immediately", func(t *testing.T) {
		out := Merge[int](context.Background())
		v, ok := <-out
		if ok {
			t.Fatalf("expected closed channel, got value %d", v)
		}
	})

	t.Run("drained inputs close output", func(t *testing.T) {
		out := Merge(context.Background(), makeChan(1, 2), makeChan(3))
		count := 0
		for range out {
			count++
		}
		if count != 3 {
			t.Fatalf("received %d values, want 3", count)
		}
		// Second receive on a drained+closed channel must be the zero/closed.
		if _, ok := <-out; ok {
			t.Fatal("channel not closed after draining")
		}
	})
}

// TestMerge_ConcurrentSourcesExactlyOnce sends a large number of distinct
// values across many producers and asserts every value is delivered exactly
// once. Synchronisation is via channel close + range, never time.Sleep.
func TestMerge_ConcurrentSourcesExactlyOnce(t *testing.T) {
	const (
		producers   = 16
		perProducer = 500
		total       = producers * perProducer
	)

	cs := make([]<-chan int, producers)
	var startWg sync.WaitGroup
	startWg.Add(producers)
	for p := range producers {
		c := make(chan int)
		cs[p] = c
		go func(c chan<- int, base int) {
			defer close(c)
			startWg.Done()
			for j := range perProducer {
				c <- base + j
			}
		}(c, p*perProducer)
	}
	startWg.Wait()

	out := Merge(context.Background(), cs...)

	var counts [total]atomic.Int32
	received := 0
	for v := range out {
		if v < 0 || v >= total {
			t.Fatalf("value out of range: %d", v)
		}
		counts[v].Add(1)
		received++
	}

	if received != total {
		t.Fatalf("received %d values, want %d", received, total)
	}
	for i := range counts {
		if got := counts[i].Load(); got != 1 {
			t.Fatalf("value %d delivered %d times, want 1", i, got)
		}
	}
}

// TestMerge_ContextCancelStopsForwarders verifies that cancelling the context
// releases forwarders even when the consumer stops reading, so the output
// channel still closes (no goroutine leak). A blocked unbuffered source that
// never closes would hang forever without ctx handling.
func TestMerge_ContextCancelStopsForwarders(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// A source that produces forever and is never closed by us.
	blocking := make(chan int)
	producerStopped := make(chan struct{})
	go func() {
		defer close(producerStopped)
		for {
			select {
			case blocking <- 42:
			case <-ctx.Done():
				return
			}
		}
	}()

	out := Merge(ctx, blocking)

	// Read a couple of values, then abandon the stream and cancel.
	<-out
	<-out
	cancel()

	// Draining out must terminate: Merge closes out after forwarders exit.
	for range out {
	}

	// Producer goroutine must also unwind via ctx (proves no dangling sender).
	<-producerStopped
}

func TestDrain(t *testing.T) {
	t.Run("collects all values", func(t *testing.T) {
		got, err := Drain(context.Background(), makeChan(1, 2), makeChan(3, 4))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		sort.Ints(got)
		want := []int{1, 2, 3, 4}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("got %v, want %v", got, want)
			}
		}
	})

	t.Run("no inputs returns ErrNoInputs", func(t *testing.T) {
		got, err := Drain[int](context.Background())
		if !errors.Is(err, ErrNoInputs) {
			t.Fatalf("err = %v, want ErrNoInputs", err)
		}
		if got != nil {
			t.Fatalf("got %v, want nil", got)
		}
	})

	t.Run("cancelled context wraps cause", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		blocking := make(chan int) // never closed, never sends
		cancel()

		got, err := Drain(ctx, blocking)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v, want wrapped context.Canceled", err)
		}
		if got != nil {
			t.Fatalf("got %v, want nil collected values", got)
		}
	})
}
