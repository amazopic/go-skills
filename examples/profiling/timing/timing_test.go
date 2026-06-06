package timing

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// fakeClock is a deterministic Clock. Each call to Now returns the next value
// from a pre-programmed sequence; this lets tests assert exact durations
// without any wall-clock sleeps and keeps results stable under -race.
type fakeClock struct {
	mu    sync.Mutex
	times []time.Time
	i     int
}

func newFakeClock(ts ...time.Time) *fakeClock {
	return &fakeClock{times: ts}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.i >= len(c.times) {
		// Exhausting the sequence is a test bug; return the last value rather
		// than panicking inside a deferred Stop.
		return c.times[len(c.times)-1]
	}
	t := c.times[c.i]
	c.i++
	return t
}

func at(ms int) time.Time {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	return base.Add(time.Duration(ms) * time.Millisecond)
}

func TestRecorder_Stats(t *testing.T) {
	tests := []struct {
		name     string
		observe  []time.Duration
		want     Stats
		wantErr  error
		statsKey string
	}{
		{
			name:     "no samples",
			observe:  nil,
			wantErr:  ErrNoSamples,
			statsKey: "missing",
		},
		{
			name:     "single sample",
			observe:  []time.Duration{5 * time.Millisecond},
			statsKey: "x",
			want: Stats{
				Count: 1,
				Sum:   5 * time.Millisecond,
				Min:   5 * time.Millisecond,
				Max:   5 * time.Millisecond,
				Mean:  5 * time.Millisecond,
			},
		},
		{
			name: "multiple samples min max mean",
			observe: []time.Duration{
				10 * time.Millisecond,
				30 * time.Millisecond,
				20 * time.Millisecond,
			},
			statsKey: "x",
			want: Stats{
				Count: 3,
				Sum:   60 * time.Millisecond,
				Min:   10 * time.Millisecond,
				Max:   30 * time.Millisecond,
				Mean:  20 * time.Millisecond,
			},
		},
		{
			name:     "negative clamped to zero",
			observe:  []time.Duration{-7 * time.Millisecond, 4 * time.Millisecond},
			statsKey: "x",
			want: Stats{
				Count: 2,
				Sum:   4 * time.Millisecond,
				Min:   0,
				Max:   4 * time.Millisecond,
				Mean:  2 * time.Millisecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRecorder(nil) // SystemClock; Observe does not use the clock.
			for _, d := range tt.observe {
				r.Observe("x", d)
			}

			got, err := r.Stats(tt.statsKey)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Stats err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Stats = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestRecorder_StartStop_DeterministicElapsed(t *testing.T) {
	// Start consumes one Now (begin), Stop consumes one Now (end).
	clk := newFakeClock(at(0), at(15), at(100), at(130))
	r := NewRecorder(clk)

	func() {
		defer r.Start("op")() // begin=0, end=15 -> 15ms
	}()
	func() {
		defer r.Start("op")() // begin=100, end=130 -> 30ms
	}()

	got, err := r.Stats("op")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := Stats{
		Count: 2,
		Sum:   45 * time.Millisecond,
		Min:   15 * time.Millisecond,
		Max:   30 * time.Millisecond,
		Mean:  22500 * time.Microsecond,
	}
	if got != want {
		t.Errorf("Stats = %+v, want %+v", got, want)
	}
}

func TestMeasure_PreservesResultAndError(t *testing.T) {
	clk := newFakeClock(at(0), at(50))
	r := NewRecorder(clk)

	wantErr := errors.New("boom")
	got, err := Measure(r, "m", func() (int, error) {
		return 42, wantErr
	})
	if got != 42 {
		t.Errorf("value = %d, want 42", got)
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}

	st, sErr := r.Stats("m")
	if sErr != nil {
		t.Fatalf("Stats error: %v", sErr)
	}
	if st.Count != 1 || st.Sum != 50*time.Millisecond {
		t.Errorf("Stats = %+v, want Count=1 Sum=50ms", st)
	}
}

func TestRecorder_LabelsAndReset(t *testing.T) {
	r := NewRecorder(nil)
	r.Observe("beta", time.Millisecond)
	r.Observe("alpha", time.Millisecond)
	r.Observe("alpha", time.Millisecond)

	labels := r.Labels()
	want := []string{"alpha", "beta"}
	if len(labels) != len(want) {
		t.Fatalf("Labels = %v, want %v", labels, want)
	}
	for i := range want {
		if labels[i] != want[i] {
			t.Errorf("Labels[%d] = %q, want %q", i, labels[i], want[i])
		}
	}

	r.Reset()
	if got := r.Labels(); len(got) != 0 {
		t.Errorf("after Reset Labels = %v, want empty", got)
	}
	if _, err := r.Stats("alpha"); !errors.Is(err, ErrNoSamples) {
		t.Errorf("after Reset Stats err = %v, want ErrNoSamples", err)
	}
}

// TestRecorder_ConcurrentObserve verifies the Recorder is race-free and counts
// every sample when many goroutines record simultaneously. Synchronisation is
// done with a WaitGroup; the SystemClock is irrelevant because Observe records
// fixed durations.
func TestRecorder_ConcurrentObserve(t *testing.T) {
	const (
		goroutines = 16
		perG       = 500
	)
	r := NewRecorder(nil)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < perG; i++ {
				r.Observe("hot", time.Millisecond)
			}
		}()
	}
	wg.Wait()

	st, err := r.Stats("hot")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := goroutines * perG; st.Count != want {
		t.Errorf("Count = %d, want %d", st.Count, want)
	}
	if want := time.Duration(goroutines*perG) * time.Millisecond; st.Sum != want {
		t.Errorf("Sum = %s, want %s", st.Sum, want)
	}
	if st.Min != time.Millisecond || st.Max != time.Millisecond {
		t.Errorf("Min/Max = %s/%s, want 1ms/1ms", st.Min, st.Max)
	}
}

// TestRecorder_ConcurrentStartAndRead exercises Start/Stop and Stats running
// concurrently to prove the read path copies under the lock and is race-clean.
func TestRecorder_ConcurrentStartAndRead(t *testing.T) {
	r := NewRecorder(SystemClock{})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			stop := r.Start("w")
			stop()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_, _ = r.Stats("w") // may legitimately return ErrNoSamples early
		}
	}()

	wg.Wait()

	st, err := r.Stats("w")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st.Count != 1000 {
		t.Errorf("Count = %d, want 1000", st.Count)
	}
}

func TestTrack_DoesNotPanic(t *testing.T) {
	// Track logs through the standard logger; we only assert it runs cleanly in
	// the deferred position it is designed for.
	func() {
		defer Track(time.Now(), "unit")
	}()
}
