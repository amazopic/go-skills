package monitor

import (
	"sync"
	"testing"
)

func TestResourcePool_AcquireRelease(t *testing.T) {
	p := NewResourcePool(3)
	if got := p.Available(); got != 3 {
		t.Fatalf("Available = %d, want 3", got)
	}
	t1 := p.Acquire()
	t2 := p.Acquire()
	if p.Available() != 1 {
		t.Errorf("Available = %d, want 1", p.Available())
	}
	p.Release(t1)
	p.Release(t2)
	if p.Available() != 3 {
		t.Errorf("Available = %d, want 3", p.Available())
	}
}

func TestResourcePool_BlocksWhenEmpty(t *testing.T) {
	p := NewResourcePool(1)
	token := p.Acquire() // exhausts the pool

	released := make(chan struct{})
	acquired := make(chan int, 1)

	go func() {
		<-released              // wait until Release is called
		acquired <- p.Acquire() // should unblock
	}()

	// Give the goroutine time to call Acquire and block.
	// Use a channel handshake: release after goroutine is waiting.
	// Since we can't observe goroutine state directly, we release after
	// a goroutine signals readiness.
	ready := make(chan struct{})
	go func() {
		close(ready)
		acquired <- p.Acquire()
	}()
	<-ready
	p.Release(token)
	got := <-acquired
	_ = got
	close(released) // signal the first goroutine (won't block — closed)
}

func TestResourcePool_ConcurrentUse(t *testing.T) {
	const (
		poolSize  = 4
		workers   = 20
		workItems = 10
	)
	p := NewResourcePool(poolSize)
	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for range workItems {
				tok := p.Acquire()
				// Simulate work — nothing here; just return the token.
				p.Release(tok)
			}
		}()
	}
	wg.Wait()
	if got := p.Available(); got != poolSize {
		t.Errorf("Available after all work = %d, want %d", got, poolSize)
	}
}

func TestResourcePool_PanicOnOverflow(t *testing.T) {
	p := NewResourcePool(1)
	p.Acquire()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on overflow")
		}
	}()
	p.Release(0)
	p.Release(0) // one extra → panic
}

func TestAtomicRegistry_GetOrSet(t *testing.T) {
	r := NewAtomicRegistry()

	v, created := r.GetOrSet("key", "first")
	if !created || v != "first" {
		t.Errorf("GetOrSet = (%q, %v), want (\"first\", true)", v, created)
	}

	v2, created2 := r.GetOrSet("key", "second")
	if created2 || v2 != "first" {
		t.Errorf("GetOrSet = (%q, %v), want (\"first\", false)", v2, created2)
	}
}

func TestAtomicRegistry_ConcurrentGetOrSet(t *testing.T) {
	r := NewAtomicRegistry()
	const goroutines = 50
	results := make(chan bool, goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			_, created := r.GetOrSet("once", "value")
			results <- created
		}()
	}
	wg.Wait()
	close(results)

	createdCount := 0
	for c := range results {
		if c {
			createdCount++
		}
	}
	if createdCount != 1 {
		t.Errorf("GetOrSet created %d times, want exactly 1", createdCount)
	}
}
