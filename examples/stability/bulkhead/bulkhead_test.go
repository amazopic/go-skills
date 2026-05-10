package bulkhead

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestPool_Do_AcquireAndRelease(t *testing.T) {
	p := NewPool("test", 2)
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := p.Do(ctx, func(ctx context.Context) error {
				time.Sleep(20 * time.Millisecond)
				return nil
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()

	acq, rej, cur := p.Stats()
	if acq != 2 {
		t.Errorf("want acquired=2, got %d", acq)
	}
	if rej != 0 {
		t.Errorf("want rejected=0, got %d", rej)
	}
	if cur != 0 {
		t.Errorf("want current=0, got %d", cur)
	}
}

func TestPool_Do_RejectsWhenFull(t *testing.T) {
	p := NewPool("limited", 1)
	ctx := context.Background()

	// Hold the single slot.
	started := make(chan struct{})
	release := make(chan struct{})
	go func() {
		_ = p.Do(ctx, func(ctx context.Context) error {
			close(started)
			<-release
			return nil
		})
	}()
	<-started

	// Second call must be rejected immediately.
	err := p.Do(ctx, func(ctx context.Context) error { return nil })
	if !errors.Is(err, ErrPoolExhausted) {
		t.Fatalf("want ErrPoolExhausted, got %v", err)
	}
	close(release)
}

func TestPool_Do_ContextAlreadyCancelled(t *testing.T) {
	p := NewPool("ctx", 2)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := p.Do(ctx, func(ctx context.Context) error { return nil })
	if err == nil {
		t.Fatal("want error for cancelled context, got nil")
	}
}

func TestPool_Do_PropagatesError(t *testing.T) {
	p := NewPool("err", 2)
	sentinel := errors.New("work failed")

	err := p.Do(context.Background(), func(ctx context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("want sentinel error, got %v", err)
	}
}

func TestBulkhead_IsolatesPools(t *testing.T) {
	b := New()
	b.AddPool("payment", 1)
	b.AddPool("search", 5)

	ctx := context.Background()

	// Saturate payment pool.
	started := make(chan struct{})
	release := make(chan struct{})
	go func() {
		_ = b.Do(ctx, "payment", func(ctx context.Context) error {
			close(started)
			<-release
			return nil
		})
	}()
	<-started

	// Payment is full — search must still work.
	err := b.Do(ctx, "search", func(ctx context.Context) error { return nil })
	if err != nil {
		t.Errorf("search pool blocked by payment pool exhaustion: %v", err)
	}

	// Payment must be rejected.
	err = b.Do(ctx, "payment", func(ctx context.Context) error { return nil })
	if !errors.Is(err, ErrPoolExhausted) {
		t.Errorf("want ErrPoolExhausted for payment pool, got %v", err)
	}

	close(release)
}

func TestBulkhead_UnknownPool(t *testing.T) {
	b := New()
	err := b.Do(context.Background(), "ghost", func(ctx context.Context) error { return nil })
	if err == nil {
		t.Fatal("want error for unknown pool")
	}
}

func TestPool_ConcurrentRace(t *testing.T) {
	// Run with -race to exercise all atomic paths.
	p := NewPool("race", 4)
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = p.Do(ctx, func(ctx context.Context) error {
				time.Sleep(time.Millisecond)
				return nil
			})
		}()
	}
	wg.Wait()

	acq, _, cur := p.Stats()
	if cur != 0 {
		t.Errorf("want current=0 after all goroutines done, got %d", cur)
	}
	if acq < 4 {
		t.Errorf("want at least 4 acquired, got %d", acq)
	}
}
