package nbarrier

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestBarrier_SingleRound(t *testing.T) {
	const n = 8
	b := NewBarrier(n)
	var arrived atomic.Int32

	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			arrived.Add(1)
			b.Wait()
			// After Wait, all n goroutines must have arrived.
			if got := arrived.Load(); got != n {
				t.Errorf("after Wait: arrived = %d, want %d", got, n)
			}
		}()
	}
	wg.Wait()
}

func TestBarrier_MultipleRounds(t *testing.T) {
	const (
		n      = 5
		rounds = 10
	)
	b := NewBarrier(n)

	// Each goroutine increments a per-round counter, barriers, then reads back.
	// If any goroutine sees a partial count the test fails.
	var counter atomic.Int32
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			for r := range rounds {
				counter.Add(1)
				b.Wait()
				// All n increments for this round must be visible.
				want := int32((r + 1) * n)
				if got := counter.Load(); got != want {
					t.Errorf("round %d: counter = %d, want %d", r, got, want)
				}
				b.Wait() // second barrier so reads don't race with next-round increments
			}
		}()
	}
	wg.Wait()
}

func TestBarrier_N1(t *testing.T) {
	// A barrier of 1 should return immediately.
	b := NewBarrier(1)
	done := make(chan struct{})
	go func() {
		b.Wait()
		close(done)
	}()
	<-done
}

func TestBarrier_PanicOnZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for n=0")
		}
	}()
	NewBarrier(0)
}

func TestBarrier_Race(t *testing.T) {
	// Stress test: many goroutines, many rounds — run with -race.
	const (
		n      = 16
		rounds = 50
	)
	b := NewBarrier(n)
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			for range rounds {
				b.Wait()
			}
		}()
	}
	wg.Wait()
}
