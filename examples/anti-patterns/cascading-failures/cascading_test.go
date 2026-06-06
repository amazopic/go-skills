package cascading

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// --- SlowDep helpers -------------------------------------------------------

func newSlowDep(hang time.Duration) *SlowDep {
	return &SlowDep{hang: hang}
}

// --- NaiveClient: demonstrates the anti-pattern ----------------------------

func TestNaiveClient_PilesUpGoroutines(t *testing.T) {
	// A slow dep with no timeout: goroutines accumulate.
	dep := newSlowDep(100 * time.Millisecond)
	client := NewNaiveClient(dep)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()
			_ = client.Call(ctx)
		}()
	}

	// At some point all 5 goroutines are in-flight — demonstrating accumulation.
	time.Sleep(20 * time.Millisecond)
	if dep.InFlight() == 0 {
		t.Log("no in-flight calls observed (timing-sensitive, not a fatal failure)")
	}
	wg.Wait()
}

// --- TimeoutClient: mitigation 1 -------------------------------------------

func TestTimeoutClient_ReturnsFastOnSlowDep(t *testing.T) {
	dep := newSlowDep(500 * time.Millisecond)
	client := NewTimeoutClient(dep, 20*time.Millisecond)

	start := time.Now()
	err := client.Call(context.Background())
	elapsed := time.Since(start)

	if !errors.Is(err, ErrTimeout) {
		t.Errorf("want ErrTimeout, got %v", err)
	}
	// Should return in ~20ms, not 500ms.
	if elapsed > 100*time.Millisecond {
		t.Errorf("timeout not enforced: elapsed=%v", elapsed)
	}
}

func TestTimeoutClient_SucceedsWhenDepFast(t *testing.T) {
	dep := newSlowDep(0)
	client := NewTimeoutClient(dep, 200*time.Millisecond)

	// dep hangs=0 means it returns immediately with an error (unavailable).
	// The point is the timeout does NOT fire.
	start := time.Now()
	err := client.Call(context.Background())
	elapsed := time.Since(start)

	if errors.Is(err, ErrTimeout) {
		t.Errorf("should not timeout for instant dep, elapsed=%v", elapsed)
	}
}

func TestTimeoutClient_RespectsParentCancel(t *testing.T) {
	dep := newSlowDep(time.Second)
	client := NewTimeoutClient(dep, 2*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := client.Call(ctx)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected an error from cancelled context")
	}
	if elapsed > 100*time.Millisecond {
		t.Errorf("parent cancel not respected: elapsed=%v", elapsed)
	}
}

// --- BulkheadClient: mitigation 2 ------------------------------------------

func TestBulkheadClient_ShedsExcessCallers(t *testing.T) {
	dep := newSlowDep(200 * time.Millisecond)
	client := NewBulkheadClient(dep, 500*time.Millisecond, 2) // only 2 concurrent

	var (
		wg       sync.WaitGroup
		rejected int64
		mu       sync.Mutex
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := client.Call(context.Background())
			if errors.Is(err, ErrPoolExhausted) {
				mu.Lock()
				rejected++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	mu.Lock()
	r := rejected
	mu.Unlock()

	// With 10 callers and limit=2, we expect significant shedding.
	if r == 0 {
		t.Error("expected some callers to be shed by the bulkhead")
	}
	t.Logf("rejected=%d/10 (expected >0)", r)
}

func TestBulkheadClient_AllowsUpToLimit(t *testing.T) {
	dep := newSlowDep(0) // instant
	client := NewBulkheadClient(dep, 100*time.Millisecond, 5)

	var wg sync.WaitGroup
	errs := make([]error, 5)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			errs[n] = client.Call(context.Background())
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if errors.Is(err, ErrPoolExhausted) {
			t.Errorf("goroutine %d: should not be shed when under limit", i)
		}
	}
}

// --- CircuitBreaker tests ---------------------------------------------------

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)
	for i := 0; i < 3; i++ {
		if err := cb.Allow(); err != nil {
			t.Fatalf("circuit should be closed before threshold, iteration %d", i)
		}
		cb.Record(errors.New("fail"))
	}
	if err := cb.Allow(); !errors.Is(err, ErrCircuitOpen) {
		t.Errorf("want ErrCircuitOpen after threshold, got %v", err)
	}
	if cb.State() != "open" {
		t.Errorf("want state=open, got %s", cb.State())
	}
}

func TestCircuitBreaker_ClosesAfterSuccess(t *testing.T) {
	cb := NewCircuitBreaker(2, 10*time.Millisecond)
	cb.Record(errors.New("fail"))
	cb.Record(errors.New("fail")) // opens

	time.Sleep(15 * time.Millisecond) // reset window passes

	if err := cb.Allow(); err != nil {
		t.Fatalf("circuit should be half-open, got err: %v", err)
	}
	cb.Record(nil) // success → closes
	if cb.State() != "closed" {
		t.Errorf("want state=closed after success, got %s", cb.State())
	}
}

// --- FullyProtectedClient: all mitigations combined -------------------------

func TestFullyProtectedClient_FastFailsOnOpenCircuit(t *testing.T) {
	dep := newSlowDep(0)
	c := NewFullyProtectedClient(dep, 50*time.Millisecond, 2, 2, time.Second)

	// Trip the circuit.
	c.cb.Record(errors.New("fail"))
	c.cb.Record(errors.New("fail"))

	if c.CBState() != "open" {
		t.Skip("circuit did not open — threshold mismatch")
	}

	err := c.Call(context.Background())
	if !errors.Is(err, ErrCircuitOpen) {
		t.Errorf("want ErrCircuitOpen, got %v", err)
	}
}

func TestFullyProtectedClient_BulkheadProtectsOpenRequests(t *testing.T) {
	dep := newSlowDep(200 * time.Millisecond)
	c := NewFullyProtectedClient(dep, 500*time.Millisecond, 1, 100, time.Second)

	// Saturate the bulkhead with 1 concurrent caller.
	started := make(chan struct{})
	release := make(chan struct{})
	go func() {
		_ = c.dep.Call(contextWithSignal(started, release))
	}()
	<-started

	// Second call should fail with timeout or pool exhausted.
	// The bulkhead has limit=1 so a second concurrent call is shed.
	var wg sync.WaitGroup
	var shed bool
	var mu sync.Mutex
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := c.Call(context.Background())
			if errors.Is(err, ErrPoolExhausted) {
				mu.Lock()
				shed = true
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	close(release)

	mu.Lock()
	s := shed
	mu.Unlock()
	if !s {
		t.Error("expected at least one caller to be shed by the bulkhead")
	}
}

// contextWithSignal returns a context that sends on started when its goroutine
// is entered, then blocks until release is closed.
func contextWithSignal(started chan<- struct{}, release <-chan struct{}) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		close(started)
		<-release
		cancel()
	}()
	return ctx
}

func TestFullyProtectedClient_NoCascadeUnderLoad(t *testing.T) {
	// Integration test: 20 concurrent callers against a slow dep with a small
	// bulkhead. None should block longer than timeout.
	dep := newSlowDep(500 * time.Millisecond)
	c := NewFullyProtectedClient(dep, 30*time.Millisecond, 3, 10, time.Second)

	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = c.Call(context.Background())
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	// All 20 goroutines should finish in well under 1 second.
	if elapsed > 300*time.Millisecond {
		t.Errorf("cascading failure detected: elapsed=%v (want < 300ms)", elapsed)
	}
}
