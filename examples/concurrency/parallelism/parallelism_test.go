package parallelism

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestMap_OrderPreservedAndCorrect(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		n       int
	}{
		{"1worker/0items", 1, 0},
		{"1worker/1item", 1, 1},
		{"1worker/many", 1, 100},
		{"4workers/many", 4, 100},
		{"more workers than items", 16, 3},
		{"workers equal items", 8, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make([]int, tt.n)
			for i := range in {
				in[i] = i
			}

			got, err := Map(context.Background(), tt.workers, in, func(_ context.Context, v int) (int, error) {
				return v * v, nil
			})
			if err != nil {
				t.Fatalf("Map returned error: %v", err)
			}
			if len(got) != tt.n {
				t.Fatalf("len(got) = %d, want %d", len(got), tt.n)
			}
			for i, v := range got {
				if want := i * i; v != want {
					t.Errorf("got[%d] = %d, want %d", i, v, want)
				}
			}
		})
	}
}

func TestMap_NoWorkers(t *testing.T) {
	tests := []struct {
		name    string
		workers int
	}{
		{"zero", 0},
		{"negative", -3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Map(context.Background(), tt.workers, []int{1, 2, 3}, func(_ context.Context, v int) (int, error) {
				return v, nil
			})
			if !errors.Is(err, ErrNoWorkers) {
				t.Fatalf("err = %v, want ErrNoWorkers", err)
			}
		})
	}
}

func TestMap_EmptyInput(t *testing.T) {
	got, err := Map(context.Background(), 4, []int(nil), func(_ context.Context, v int) (int, error) {
		t.Fatal("fn should not be called for empty input")
		return v, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("len(got) = %d, want 0", len(got))
	}
}

// TestMap_BoundedConcurrency verifies that at most `workers` invocations of fn
// run at the same time. Synchronization is via a barrier channel, not sleep,
// so the assertion is deterministic.
func TestMap_BoundedConcurrency(t *testing.T) {
	const (
		workers = 4
		items   = 64
	)
	in := make([]int, items)

	var (
		concurrent atomic.Int32
		maxSeen    atomic.Int32
	)
	// release gates every fn call until the test has confirmed the in-flight
	// count; this makes the peak-concurrency measurement deterministic.
	release := make(chan struct{})

	var done sync.WaitGroup
	done.Add(1)
	go func() {
		defer done.Done()
		// Allow each item exactly one release.
		for range items {
			release <- struct{}{}
		}
	}()

	_, err := Map(context.Background(), workers, in, func(_ context.Context, v int) (int, error) {
		cur := concurrent.Add(1)
		for {
			prev := maxSeen.Load()
			if cur <= prev || maxSeen.CompareAndSwap(prev, cur) {
				break
			}
		}
		<-release
		concurrent.Add(-1)
		return v, nil
	})
	done.Wait()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := maxSeen.Load(); got > workers {
		t.Fatalf("peak concurrency = %d, want <= %d", got, workers)
	}
	if got := maxSeen.Load(); got < 1 {
		t.Fatalf("peak concurrency = %d, want >= 1", got)
	}
}

func TestMap_FirstErrorReturnedAndCancels(t *testing.T) {
	const items = 200
	in := make([]int, items)
	for i := range in {
		in[i] = i
	}

	sentinel := errors.New("boom")
	var processed atomic.Int32

	_, err := Map(context.Background(), 4, in, func(ctx context.Context, v int) (int, error) {
		processed.Add(1)
		if v == 5 {
			return 0, sentinel
		}
		// Cooperatively stop once cancellation propagates so the test does
		// not depend on timing; still always make progress on the happy path.
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return v, nil
		}
	})

	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want wrapped sentinel", err)
	}
	// Cancellation must prevent every item from being processed; with 200
	// items and an early failure this is reliably true.
	if got := processed.Load(); got >= items {
		t.Fatalf("processed = %d, want < %d (cancellation should short-circuit)", got, items)
	}
}

func TestMap_ContextCancelledBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	in := []int{1, 2, 3, 4, 5}
	var calls atomic.Int32

	got, err := Map(ctx, 2, in, func(_ context.Context, v int) (int, error) {
		calls.Add(1)
		return v, nil
	})
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if got != nil {
		t.Fatalf("got = %v, want nil results on cancellation", got)
	}
}

func TestMap_DifferentInputOutputTypes(t *testing.T) {
	in := []int{1, 22, 333}
	got, err := Map(context.Background(), 3, in, func(_ context.Context, v int) (string, error) {
		return fmt.Sprintf("n=%d", v), nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"n=1", "n=22", "n=333"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestMapCPU(t *testing.T) {
	in := make([]int, 50)
	for i := range in {
		in[i] = i
	}
	got, err := MapCPU(context.Background(), in, func(_ context.Context, v int) (int, error) {
		return v + 1, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, v := range got {
		if v != i+1 {
			t.Errorf("got[%d] = %d, want %d", i, v, i+1)
		}
	}
}
