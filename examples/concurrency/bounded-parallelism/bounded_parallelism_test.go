package boundedparallelism

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestMap_Results(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		in      []int
		want    []int
	}{
		{"empty", 4, []int{}, []int{}},
		{"single", 4, []int{5}, []int{25}},
		{"one_worker", 1, []int{1, 2, 3, 4}, []int{1, 4, 9, 16}},
		{"workers_equal_items", 4, []int{1, 2, 3, 4}, []int{1, 4, 9, 16}},
		{"more_workers_than_items", 16, []int{2, 3}, []int{4, 9}},
		{"large_input", 8, seq(1000), squares(1000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Map(context.Background(), tt.workers, tt.in,
				func(_ context.Context, v int) (int, error) { return v * v, nil })
			if err != nil {
				t.Fatalf("Map returned error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("len(got) = %d, want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("got[%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestMap_ConcurrencyIsBounded(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		items   int
	}{
		{"k2", 2, 200},
		{"k4", 4, 200},
		{"k8", 8, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// release gates every worker call so they all pile up to the limit
			// at once; we then assert the observed peak never exceeds workers.
			release := make(chan struct{})
			var inFlight, peak atomic.Int64

			// allArrived fires once we've seen `workers` concurrent calls, so
			// the test makes progress without time.Sleep.
			allArrived := make(chan struct{})
			var arrivedOnce sync.Once

			fn := func(_ context.Context, v int) (int, error) {
				cur := inFlight.Add(1)
				for {
					p := peak.Load()
					if cur <= p || peak.CompareAndSwap(p, cur) {
						break
					}
				}
				if cur == int64(tt.workers) {
					arrivedOnce.Do(func() { close(allArrived) })
				}
				<-release // block until the test opens the gate
				inFlight.Add(-1)
				return v * v, nil
			}

			done := make(chan struct{})
			go func() {
				defer close(done)
				if _, err := Map(context.Background(), tt.workers, seq(tt.items), fn); err != nil {
					t.Errorf("Map error: %v", err)
				}
			}()

			<-allArrived // exactly `workers` calls are now blocked
			close(release)
			<-done

			if got := peak.Load(); got != int64(tt.workers) {
				t.Errorf("peak concurrency = %d, want %d", got, tt.workers)
			}
		})
	}
}

func TestMap_AllItemsProcessedExactlyOnce(t *testing.T) {
	const (
		items   = 1000
		workers = 7
	)
	var counts [items]atomic.Int32

	got, err := Map(context.Background(), workers, seq(items),
		func(_ context.Context, v int) (int, error) {
			counts[v].Add(1)
			return v, nil
		})
	if err != nil {
		t.Fatalf("Map error: %v", err)
	}
	if len(got) != items {
		t.Fatalf("len(got) = %d, want %d", len(got), items)
	}
	for i := range items {
		if c := counts[i].Load(); c != 1 {
			t.Errorf("item %d processed %d times, want 1", i, c)
		}
		if got[i] != i {
			t.Errorf("got[%d] = %d, want %d (order not preserved)", i, got[i], i)
		}
	}
}

func TestMap_FirstErrorReported(t *testing.T) {
	sentinel := errors.New("boom")
	const failAt = 42

	_, err := Map(context.Background(), 4, seq(1000),
		func(_ context.Context, v int) (int, error) {
			if v == failAt {
				return 0, sentinel
			}
			return v, nil
		})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("error %v does not wrap sentinel", err)
	}
}

func TestMap_ErrorStopsRemainingWork(t *testing.T) {
	// Worker 0's item fails immediately; gate the rest so they would block
	// forever if they ran. The error must cancel them so Map returns.
	gate := make(chan struct{})
	defer close(gate)
	sentinel := errors.New("stop")
	var started atomic.Int64

	_, err := Map(context.Background(), 4, seq(1000),
		func(ctx context.Context, v int) (int, error) {
			started.Add(1)
			if v == 0 {
				return 0, sentinel
			}
			// Block until either the test ends or ctx is cancelled by the
			// failing item. A leak/hang here means the test times out.
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-gate:
				return v, nil
			}
		})
	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want wrap of sentinel", err)
	}
	if s := started.Load(); s > 1000 {
		t.Errorf("started %d calls, want <= input length", s)
	}
}

func TestMap_ContextCancellation(t *testing.T) {
	tests := []struct {
		name      string
		cancelPre bool // cancel before calling Map
	}{
		{"already_cancelled", true},
		{"cancelled_mid_flight", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			if tt.cancelPre {
				cancel()
				defer cancel()
				_, err := Map(ctx, 4, seq(100),
					func(_ context.Context, v int) (int, error) { return v, nil })
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("err = %v, want context.Canceled", err)
				}
				return
			}

			// Mid-flight: block the first worker calls, cancel, then assert.
			entered := make(chan struct{}, 1)
			_, err := Map(ctx, 1, seq(100),
				func(ctx context.Context, v int) (int, error) {
					select {
					case entered <- struct{}{}:
						cancel() // cancel from inside the very first call
					default:
					}
					<-ctx.Done()
					return 0, ctx.Err()
				})
			defer cancel()
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("err = %v, want context.Canceled", err)
			}
		})
	}
}

func TestMap_InvalidWorkers(t *testing.T) {
	for _, w := range []int{0, -1, -100} {
		_, err := Map(context.Background(), w, seq(10),
			func(_ context.Context, v int) (int, error) { return v, nil })
		if !errors.Is(err, ErrInvalidWorkers) {
			t.Errorf("workers=%d: err = %v, want ErrInvalidWorkers", w, err)
		}
	}
}

func TestForEach(t *testing.T) {
	t.Run("side_effects_applied_once", func(t *testing.T) {
		const items = 500
		var sum atomic.Int64
		err := ForEach(context.Background(), 8, seq(items),
			func(_ context.Context, v int) error {
				sum.Add(int64(v))
				return nil
			})
		if err != nil {
			t.Fatalf("ForEach error: %v", err)
		}
		want := int64(items * (items - 1) / 2)
		if got := sum.Load(); got != want {
			t.Errorf("sum = %d, want %d", got, want)
		}
	})

	t.Run("propagates_error", func(t *testing.T) {
		sentinel := errors.New("nope")
		err := ForEach(context.Background(), 4, seq(50),
			func(_ context.Context, v int) error {
				if v == 7 {
					return sentinel
				}
				return nil
			})
		if !errors.Is(err, sentinel) {
			t.Errorf("err = %v, want wrap of sentinel", err)
		}
	})

	t.Run("invalid_workers", func(t *testing.T) {
		err := ForEach(context.Background(), 0, seq(5),
			func(_ context.Context, _ int) error { return nil })
		if !errors.Is(err, ErrInvalidWorkers) {
			t.Errorf("err = %v, want ErrInvalidWorkers", err)
		}
	})
}

// --- helpers ---

func seq(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

func squares(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i * i
	}
	return s
}

// Sanity-check that helper output is what tests assume.
func TestHelpers(t *testing.T) {
	if got := fmt.Sprint(seq(3)); got != "[0 1 2]" {
		t.Fatalf("seq(3) = %s", got)
	}
}
