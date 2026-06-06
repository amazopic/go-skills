package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeClock is a deterministic, race-safe time source. Tests advance it
// explicitly instead of sleeping, so behavior is reproducible under -race.
type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func newFakeClock() *fakeClock { return &fakeClock{t: time.Unix(0, 0)} }

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

var errDownstream = errors.New("downstream boom")

// ok and fail are canned operations returning a string result.
func ok(_ context.Context) (string, error)   { return "ok", nil }
func fail(_ context.Context) (string, error) { return "", errDownstream }

func TestBreaker_Transitions(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T, b *Breaker[string], clk *fakeClock)
	}{
		{
			name: "happy path stays closed",
			run: func(t *testing.T, b *Breaker[string], _ *fakeClock) {
				for range 5 {
					res, err := b.Do(context.Background(), ok)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					if res != "ok" {
						t.Fatalf("res = %q, want ok", res)
					}
				}
				if got := b.State(); got != StateClosed {
					t.Fatalf("state = %v, want closed", got)
				}
			},
		},
		{
			name: "trips after threshold consecutive failures",
			run: func(t *testing.T, b *Breaker[string], _ *fakeClock) {
				// threshold is 3 (see makeBreaker). Two failures: still closed.
				for i := range 2 {
					if _, err := b.Do(context.Background(), fail); !errors.Is(err, errDownstream) {
						t.Fatalf("call %d: err = %v, want errDownstream", i, err)
					}
					if got := b.State(); got != StateClosed {
						t.Fatalf("after %d failures state = %v, want closed", i+1, got)
					}
				}
				// Third failure trips the breaker.
				if _, err := b.Do(context.Background(), fail); !errors.Is(err, errDownstream) {
					t.Fatalf("err = %v, want errDownstream", err)
				}
				if got := b.State(); got != StateOpen {
					t.Fatalf("state = %v, want open", got)
				}
			},
		},
		{
			name: "open fails fast without invoking op",
			run: func(t *testing.T, b *Breaker[string], _ *fakeClock) {
				trip(t, b)

				var calls atomic.Int32
				probe := func(_ context.Context) (string, error) {
					calls.Add(1)
					return "ok", nil
				}
				_, err := b.Do(context.Background(), probe)
				if !errors.Is(err, ErrOpenCircuit) {
					t.Fatalf("err = %v, want ErrOpenCircuit", err)
				}
				if n := calls.Load(); n != 0 {
					t.Fatalf("op invoked %d times while open, want 0", n)
				}
			},
		},
		{
			name: "half-open recovery closes on success",
			run: func(t *testing.T, b *Breaker[string], clk *fakeClock) {
				trip(t, b)
				// Before cooldown: still rejected.
				if _, err := b.Do(context.Background(), ok); !errors.Is(err, ErrOpenCircuit) {
					t.Fatalf("err = %v, want ErrOpenCircuit before cooldown", err)
				}
				// Advance past the cooldown; one trial is admitted and succeeds.
				clk.Advance(2 * time.Second)
				res, err := b.Do(context.Background(), ok)
				if err != nil {
					t.Fatalf("half-open trial err = %v, want nil", err)
				}
				if res != "ok" {
					t.Fatalf("res = %q, want ok", res)
				}
				if got := b.State(); got != StateClosed {
					t.Fatalf("state = %v, want closed after recovery", got)
				}
				if got := b.ConsecutiveFailures(); got != 0 {
					t.Fatalf("consecutive failures = %d, want 0", got)
				}
			},
		},
		{
			name: "half-open failure re-opens with longer backoff",
			run: func(t *testing.T, b *Breaker[string], clk *fakeClock) {
				trip(t, b) // openCount now 1, cooldown was base (1s)

				clk.Advance(time.Second) // reach first retry window
				// Trial fails -> re-open with base<<1 = 2s.
				if _, err := b.Do(context.Background(), fail); !errors.Is(err, errDownstream) {
					t.Fatalf("err = %v, want errDownstream", err)
				}
				if got := b.State(); got != StateOpen {
					t.Fatalf("state = %v, want open", got)
				}
				// After 1s (less than 2s) still open.
				clk.Advance(time.Second)
				if _, err := b.Do(context.Background(), ok); !errors.Is(err, ErrOpenCircuit) {
					t.Fatalf("err = %v, want ErrOpenCircuit during longer backoff", err)
				}
				// After another 1s (total 2s) the trial is admitted again.
				clk.Advance(time.Second)
				if _, err := b.Do(context.Background(), ok); err != nil {
					t.Fatalf("err = %v, want nil after longer backoff elapsed", err)
				}
				if got := b.State(); got != StateClosed {
					t.Fatalf("state = %v, want closed", got)
				}
			},
		},
		{
			name: "cancelled context short-circuits without changing state",
			run: func(t *testing.T, b *Breaker[string], _ *fakeClock) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				var calls atomic.Int32
				probe := func(_ context.Context) (string, error) {
					calls.Add(1)
					return "ok", nil
				}
				_, err := b.Do(ctx, probe)
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("err = %v, want context.Canceled", err)
				}
				if n := calls.Load(); n != 0 {
					t.Fatalf("op invoked %d times with cancelled ctx, want 0", n)
				}
				if got := b.State(); got != StateClosed {
					t.Fatalf("state = %v, want closed", got)
				}
			},
		},
		{
			name: "success resets a partial failure streak",
			run: func(t *testing.T, b *Breaker[string], _ *fakeClock) {
				_, _ = b.Do(context.Background(), fail)
				_, _ = b.Do(context.Background(), fail) // 2 failures, threshold 3
				if got := b.ConsecutiveFailures(); got != 2 {
					t.Fatalf("consecutive failures = %d, want 2", got)
				}
				if _, err := b.Do(context.Background(), ok); err != nil {
					t.Fatalf("err = %v, want nil", err)
				}
				if got := b.ConsecutiveFailures(); got != 0 {
					t.Fatalf("consecutive failures = %d, want 0 after success", got)
				}
				if got := b.State(); got != StateClosed {
					t.Fatalf("state = %v, want closed", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clk := newFakeClock()
			b := makeBreaker(clk)
			tt.run(t, b, clk)
		})
	}
}

// makeBreaker builds a breaker with threshold 3 and a 1s base cooldown wired to
// the fake clock.
func makeBreaker(clk *fakeClock) *Breaker[string] {
	return New[string](3, time.Second, WithClock[string](clk.Now))
}

// trip drives a fresh breaker from closed to open via threshold failures.
func trip(t *testing.T, b *Breaker[string]) {
	t.Helper()
	for i := range 3 {
		if _, err := b.Do(context.Background(), fail); !errors.Is(err, errDownstream) {
			t.Fatalf("trip call %d: err = %v, want errDownstream", i, err)
		}
	}
	if got := b.State(); got != StateOpen {
		t.Fatalf("after tripping state = %v, want open", got)
	}
}

func TestState_String(t *testing.T) {
	tests := []struct {
		s    State
		want string
	}{
		{StateClosed, "closed"},
		{StateOpen, "open"},
		{StateHalfOpen, "half-open"},
		{State(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.s.String(); got != tt.want {
			t.Errorf("State(%d).String() = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func TestNew_PanicsOnBadConfig(t *testing.T) {
	tests := []struct {
		name      string
		threshold uint32
		cooldown  time.Duration
	}{
		{"zero threshold", 0, time.Second},
		{"zero cooldown", 3, 0},
		{"negative cooldown", 3, -time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for threshold=%d cooldown=%v", tt.threshold, tt.cooldown)
				}
			}()
			New[string](tt.threshold, tt.cooldown)
		})
	}
}

// TestBreaker_ConcurrentSafety hammers the breaker from many goroutines while
// the dependency flips from failing to healthy. It asserts there is no data
// race (run under -race) and that, once the dependency is healthy and the
// cooldown elapses, the breaker converges to closed and serves successes.
func TestBreaker_ConcurrentSafety(t *testing.T) {
	clk := newFakeClock()
	b := New[int](5, time.Second, WithClock[int](clk.Now))

	var healthy atomic.Bool // false: dependency down; true: recovered
	op := func(_ context.Context) (int, error) {
		if healthy.Load() {
			return 42, nil
		}
		return 0, errDownstream
	}

	const (
		goroutines = 32
		iterations = 200
	)

	// Phase 1: dependency is down. Concurrent callers either see the
	// downstream error or a fast ErrOpenCircuit — never a panic or race.
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				_, err := b.Do(context.Background(), op)
				if err != nil && !errors.Is(err, errDownstream) && !errors.Is(err, ErrOpenCircuit) {
					t.Errorf("unexpected error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	// The breaker must have tripped under sustained failure.
	if got := b.State(); got != StateOpen {
		t.Fatalf("state after failure storm = %v, want open", got)
	}

	// Phase 2: dependency recovers and enough time passes for a trial.
	healthy.Store(true)
	clk.Advance(time.Hour) // well past any backoff window

	// A single successful call must close the circuit.
	if _, err := b.Do(context.Background(), op); err != nil {
		t.Fatalf("recovery call err = %v, want nil", err)
	}
	if got := b.State(); got != StateClosed {
		t.Fatalf("state after recovery = %v, want closed", got)
	}

	// Phase 3: now-healthy concurrent traffic all succeeds.
	var successes atomic.Int64
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				res, err := b.Do(context.Background(), op)
				if err != nil {
					t.Errorf("healthy call err = %v, want nil", err)
					return
				}
				if res != 42 {
					t.Errorf("res = %d, want 42", res)
					return
				}
				successes.Add(1)
			}
		}()
	}
	wg.Wait()

	if got, want := successes.Load(), int64(goroutines*iterations); got != want {
		t.Fatalf("successes = %d, want %d", got, want)
	}
	if got := b.State(); got != StateClosed {
		t.Fatalf("final state = %v, want closed", got)
	}
}
