package observer

import (
	"sync"
	"sync/atomic"
	"testing"
)

// recordingObserver is a test double that records the last pushed state and how
// many times it was notified. It is safe for concurrent Update calls so it can
// be used under the race detector.
type recordingObserver struct {
	mu    sync.Mutex
	last  string
	calls int
}

func (r *recordingObserver) Update(state string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.last = state
	r.calls++
}

func (r *recordingObserver) snapshot() (string, int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.last, r.calls
}

// TestNotify drives attach/detach/notify correctness through table-driven
// scenarios. Each step mutates the publisher, then Notify is fired and every
// observer's recorded state is asserted.
func TestNotify(t *testing.T) {
	type step struct {
		attach []*recordingObserver
		detach []*recordingObserver
		state  string
	}

	a := &recordingObserver{}
	b := &recordingObserver{}
	c := &recordingObserver{}

	tests := []struct {
		name string
		step step
		// want maps an observer to the state it must hold after Notify, and the
		// number of times it must have been called so far.
		want map[*recordingObserver]struct {
			state string
			calls int
		}
	}{
		{
			name: "single observer receives state",
			step: step{attach: []*recordingObserver{a}, state: "one"},
			want: map[*recordingObserver]struct {
				state string
				calls int
			}{
				a: {"one", 1},
			},
		},
		{
			name: "multiple observers all receive state",
			step: step{attach: []*recordingObserver{b, c}, state: "two"},
			want: map[*recordingObserver]struct {
				state string
				calls int
			}{
				a: {"two", 2},
				b: {"two", 1},
				c: {"two", 1},
			},
		},
		{
			name: "detached observer stops receiving state",
			step: step{detach: []*recordingObserver{b}, state: "three"},
			want: map[*recordingObserver]struct {
				state string
				calls int
			}{
				a: {"three", 3},
				b: {"two", 1}, // unchanged: detached before this Notify
				c: {"three", 2},
			},
		},
		{
			name: "detach of unattached observer is a no-op",
			step: step{detach: []*recordingObserver{&recordingObserver{}}, state: "four"},
			want: map[*recordingObserver]struct {
				state string
				calls int
			}{
				a: {"four", 4},
				b: {"two", 1},
				c: {"four", 3},
			},
		},
	}

	pub := NewPublisher()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, o := range tt.step.attach {
				pub.Attach(o)
			}
			for _, o := range tt.step.detach {
				pub.Detach(o)
			}
			pub.SetState(tt.step.state)
			pub.Notify()

			for obs, want := range tt.want {
				gotState, gotCalls := obs.snapshot()
				if gotState != want.state {
					t.Errorf("observer state = %q, want %q", gotState, want.state)
				}
				if gotCalls != want.calls {
					t.Errorf("observer calls = %d, want %d", gotCalls, want.calls)
				}
			}
		})
	}
}

// TestNotifyNoObservers verifies Notify on an empty publisher does not panic.
func TestNotifyNoObservers(t *testing.T) {
	pub := NewPublisher()
	pub.SetState("ignored")
	pub.Notify() // must not panic
}

// TestConcurrentAttachNotify exercises the shared observer list from many
// goroutines simultaneously: attachers append while notifiers iterate and
// state setters mutate. Before the mutex fix, append-during-iterate triggers
// the race detector (concurrent read/write of s.observers and s.state).
//
// We assert every attached observer is eventually notified at least once after
// the storm settles, proving the snapshot-under-lock logic is also correct.
func TestConcurrentAttachNotify(t *testing.T) {
	pub := NewPublisher()

	const (
		attachers          = 16
		observersPerWorker = 16
		notifiers          = 16
		notifyRounds       = 200
	)

	all := make([]*recordingObserver, 0, attachers*observersPerWorker)
	var allMu sync.Mutex

	var wg sync.WaitGroup

	// Attachers: concurrently grow the observer list.
	for i := 0; i < attachers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < observersPerWorker; j++ {
				o := &recordingObserver{}
				allMu.Lock()
				all = append(all, o)
				allMu.Unlock()
				pub.Attach(o)
			}
		}()
	}

	// State setters + notifiers: concurrently read the list and the state.
	var counter int64
	for i := 0; i < notifiers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := 0; r < notifyRounds; r++ {
				pub.SetState("s" + string(rune('A'+int(atomic.AddInt64(&counter, 1))%26)))
				pub.Notify()
			}
		}()
	}

	wg.Wait()

	// All observers must have been attached.
	allMu.Lock()
	got := len(all)
	allMu.Unlock()
	if want := attachers * observersPerWorker; got != want {
		t.Fatalf("attached observers = %d, want %d", got, want)
	}

	// A final Notify after the storm must reach every observer (proving none
	// were lost by a racy append) and leave them all on the same final state.
	pub.SetState("final")
	pub.Notify()

	allMu.Lock()
	defer allMu.Unlock()
	for i, o := range all {
		state, calls := o.snapshot()
		if calls == 0 {
			t.Errorf("observer %d was never notified", i)
		}
		if state != "final" {
			t.Errorf("observer %d final state = %q, want %q", i, state, "final")
		}
	}
}

// ExampleObserver demonstrates the push-model flow end to end.
func ExampleObserver() {

	publisher := NewPublisher()

	publisher.Attach(&ConcreteObserver{})
	publisher.Attach(&ConcreteObserver{})
	publisher.Attach(&ConcreteObserver{})

	publisher.SetState("New State...")

	publisher.Notify()
}
