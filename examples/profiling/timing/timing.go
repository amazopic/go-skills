// Package timing provides lightweight, dependency-free latency measurement for
// hot paths where reaching for pprof or a metrics framework would be overkill.
//
// The idiomatic shape is a one-line defer at the top of a function:
//
//	func work() {
//		defer timing.Track(time.Now(), "work")
//		// ... do the work ...
//	}
//
// Arguments to a deferred call are evaluated immediately, so time.Now() captures
// the entry timestamp and the elapsed duration is computed when the function
// returns. Track logs through the standard logger.
//
// For tests and for aggregating call latencies, use a Recorder. A Recorder is
// driven by an injectable Clock so its behaviour is fully deterministic — no
// reliance on wall-clock sleeps — and is safe for concurrent use.
package timing

import (
	"errors"
	"log"
	"sort"
	"sync"
	"time"
)

// ErrNoSamples is returned by Stats when a label has not recorded any samples.
var ErrNoSamples = errors.New("timing: no samples recorded")

// Track logs the time elapsed since start under name. Call it as a deferred
// statement so start is captured at function entry:
//
//	defer timing.Track(time.Now(), "BigComputation")
func Track(start time.Time, name string) {
	log.Printf("%s took %s", name, time.Since(start))
}

// Clock abstracts the source of time so callers (and tests) can control it.
// The standard library's time.Now satisfies the intent of Now; SystemClock
// adapts it to this interface.
type Clock interface {
	Now() time.Time
}

// SystemClock is the production Clock backed by time.Now.
type SystemClock struct{}

// Now returns the current wall-clock time.
func (SystemClock) Now() time.Time { return time.Now() }

// Stats is an immutable summary of the durations recorded for one label.
type Stats struct {
	Count int
	Sum   time.Duration
	Min   time.Duration
	Max   time.Duration
	Mean  time.Duration
}

// Recorder accumulates timing samples grouped by label. The zero value is not
// usable; construct one with NewRecorder. A Recorder is safe for concurrent use
// by multiple goroutines.
type Recorder struct {
	clk Clock

	mu      sync.Mutex
	samples map[string][]time.Duration
}

// NewRecorder returns a Recorder driven by clk. If clk is nil, SystemClock is
// used. Pass a controllable Clock in tests to keep measurements deterministic.
func NewRecorder(clk Clock) *Recorder {
	if clk == nil {
		clk = SystemClock{}
	}
	return &Recorder{
		clk:     clk,
		samples: make(map[string][]time.Duration),
	}
}

// Start records the current time and returns a Stop function. Defer the
// returned function to record the elapsed duration under label:
//
//	defer rec.Start("query")()
//
// The returned Stop is idempotent in the sense that it captures the elapsed
// time when invoked; calling it more than once records additional samples, so
// invoke it exactly once (which the defer form guarantees).
func (r *Recorder) Start(label string) (stop func()) {
	begin := r.clk.Now()
	return func() {
		r.record(label, r.clk.Now().Sub(begin))
	}
}

// Observe records an already-measured duration d under label. Negative
// durations are clamped to zero so a misbehaving clock cannot corrupt Min.
func (r *Recorder) Observe(label string, d time.Duration) {
	if d < 0 {
		d = 0
	}
	r.record(label, d)
}

func (r *Recorder) record(label string, d time.Duration) {
	r.mu.Lock()
	r.samples[label] = append(r.samples[label], d)
	r.mu.Unlock()
}

// Stats returns the aggregated statistics for label. It returns ErrNoSamples
// if nothing has been recorded under that label.
func (r *Recorder) Stats(label string) (Stats, error) {
	r.mu.Lock()
	ds := r.samples[label]
	// Copy under the lock so we never read a slice another goroutine appends to.
	cp := make([]time.Duration, len(ds))
	copy(cp, ds)
	r.mu.Unlock()

	if len(cp) == 0 {
		return Stats{}, ErrNoSamples
	}

	s := Stats{
		Count: len(cp),
		Min:   cp[0],
		Max:   cp[0],
	}
	for _, d := range cp {
		s.Sum += d
		if d < s.Min {
			s.Min = d
		}
		if d > s.Max {
			s.Max = d
		}
	}
	s.Mean = s.Sum / time.Duration(s.Count)
	return s, nil
}

// Labels returns the recorded labels in sorted order. The result is a fresh
// slice owned by the caller.
func (r *Recorder) Labels() []string {
	r.mu.Lock()
	out := make([]string, 0, len(r.samples))
	for label := range r.samples {
		out = append(out, label)
	}
	r.mu.Unlock()
	sort.Strings(out)
	return out
}

// Reset discards all recorded samples, returning the Recorder to an empty state.
func (r *Recorder) Reset() {
	r.mu.Lock()
	r.samples = make(map[string][]time.Duration)
	r.mu.Unlock()
}

// Measure times the execution of fn, records the elapsed duration under label,
// and returns fn's value. It is the generic, return-value-preserving companion
// to the defer form:
//
//	n, err := timing.Measure(rec, "parse", func() (int, error) { ... })
func Measure[T any](r *Recorder, label string, fn func() (T, error)) (T, error) {
	stop := r.Start(label)
	defer stop()
	return fn()
}
