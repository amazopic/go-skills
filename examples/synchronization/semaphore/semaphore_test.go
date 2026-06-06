package semaphore

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNew_PanicsOnNonPositive(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"zero", 0},
		{"negative", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("New(%d) did not panic", tt.n)
				}
			}()
			New(tt.n)
		})
	}
}

func TestNew_CapAndLen(t *testing.T) {
	s := New(3)
	if got := s.Cap(); got != 3 {
		t.Errorf("Cap() = %d, want 3", got)
	}
	if got := s.Len(); got != 0 {
		t.Errorf("fresh Len() = %d, want 0", got)
	}
}

func TestTryAcquire_ExhaustsAndRecovers(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"one slot", 1},
		{"three slots", 3},
		{"ten slots", 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.n)

			// Fill every slot.
			for i := 0; i < tt.n; i++ {
				if err := s.TryAcquire(); err != nil {
					t.Fatalf("TryAcquire #%d: unexpected error %v", i, err)
				}
			}
			if got := s.Len(); got != tt.n {
				t.Errorf("Len() = %d, want %d", got, tt.n)
			}

			// One more must fail fast.
			if err := s.TryAcquire(); !errors.Is(err, ErrNoTickets) {
				t.Errorf("over-acquire error = %v, want ErrNoTickets", err)
			}

			// Release one, then a single TryAcquire must succeed again.
			s.Release()
			if got := s.Len(); got != tt.n-1 {
				t.Errorf("after Release Len() = %d, want %d", got, tt.n-1)
			}
			if err := s.TryAcquire(); err != nil {
				t.Errorf("re-acquire after Release: %v", err)
			}
		})
	}
}

func TestAcquire_RespectsCancelledContext(t *testing.T) {
	t.Run("already cancelled", func(t *testing.T) {
		s := New(1)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := s.Acquire(ctx)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Acquire error = %v, want context.Canceled", err)
		}
		if got := s.Len(); got != 0 {
			t.Errorf("Len() = %d after failed Acquire, want 0", got)
		}
	})

	t.Run("blocks until full then cancels", func(t *testing.T) {
		s := New(1)
		// Take the only ticket so the next Acquire must block.
		if err := s.TryAcquire(); err != nil {
			t.Fatalf("setup TryAcquire: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		errc := make(chan error, 1)
		started := make(chan struct{})
		go func() {
			close(started)
			errc <- s.Acquire(ctx)
		}()

		<-started // ensure the goroutine is running before we cancel
		cancel()

		err := <-errc
		if !errors.Is(err, context.Canceled) {
			t.Errorf("blocked Acquire error = %v, want context.Canceled", err)
		}
		if got := s.Len(); got != 1 {
			t.Errorf("Len() = %d, want 1 (only the original holder)", got)
		}
	})
}

func TestAcquire_BlocksUntilReleased(t *testing.T) {
	s := New(1)
	if err := s.Acquire(context.Background()); err != nil {
		t.Fatalf("first Acquire: %v", err)
	}

	acquired := make(chan struct{})
	go func() {
		// This Acquire blocks until the main goroutine releases.
		if err := s.Acquire(context.Background()); err != nil {
			t.Errorf("second Acquire: %v", err)
		}
		close(acquired)
	}()

	select {
	case <-acquired:
		t.Fatal("second Acquire returned before Release")
	default:
	}

	s.Release() // hand the ticket over to the waiting goroutine
	<-acquired  // must now proceed
}

func TestRelease_PanicsWithoutTicket(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Release without held ticket did not panic")
		}
	}()
	New(2).Release()
}

func TestDo_RunsAndReleases(t *testing.T) {
	s := New(1)
	ran := false
	err := s.Do(context.Background(), func(context.Context) error {
		ran = true
		if got := s.Len(); got != 1 {
			t.Errorf("inside Do Len() = %d, want 1", got)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if !ran {
		t.Error("fn did not run")
	}
	if got := s.Len(); got != 0 {
		t.Errorf("after Do Len() = %d, want 0", got)
	}
}

func TestDo_ReleasesOnPanic(t *testing.T) {
	s := New(1)
	func() {
		defer func() { _ = recover() }()
		_ = s.Do(context.Background(), func(context.Context) error {
			panic("boom")
		})
	}()
	if got := s.Len(); got != 0 {
		t.Errorf("after panicking Do Len() = %d, want 0 (ticket leaked)", got)
	}
}

func TestDo_PropagatesError(t *testing.T) {
	s := New(1)
	sentinel := errors.New("work failed")
	err := s.Do(context.Background(), func(context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Errorf("Do error = %v, want %v", err, sentinel)
	}
}

// TestConcurrentLimitNeverExceeded is the core correctness property: under heavy
// contention the number of goroutines simultaneously holding a ticket never
// exceeds the configured limit.
func TestConcurrentLimitNeverExceeded(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		goroutines int
	}{
		{"limit1", 1, 50},
		{"limit4", 4, 100},
		{"limit8", 8, 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.limit)
			var inside atomic.Int64
			var maxSeen atomic.Int64
			var completed atomic.Int64

			var wg sync.WaitGroup
			wg.Add(tt.goroutines)
			for i := 0; i < tt.goroutines; i++ {
				go func() {
					defer wg.Done()
					err := s.Do(context.Background(), func(context.Context) error {
						n := inside.Add(1)
						// Track the high-water mark of concurrent holders.
						for {
							m := maxSeen.Load()
							if n <= m || maxSeen.CompareAndSwap(m, n) {
								break
							}
						}
						inside.Add(-1)
						completed.Add(1)
						return nil
					})
					if err != nil {
						t.Errorf("Do: %v", err)
					}
				}()
			}
			wg.Wait()

			if got := completed.Load(); got != int64(tt.goroutines) {
				t.Errorf("completed = %d, want %d", got, tt.goroutines)
			}
			if got := maxSeen.Load(); got > int64(tt.limit) {
				t.Errorf("max concurrent holders = %d, exceeds limit %d", got, tt.limit)
			}
			if got := s.Len(); got != 0 {
				t.Errorf("after all work Len() = %d, want 0", got)
			}
		})
	}
}
