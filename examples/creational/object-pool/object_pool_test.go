package objectpool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// newCountingFactory returns a factory that produces sequentially numbered
// *int objects and an atomic counter of how many were ever built.
func newCountingFactory() (func(context.Context) (*int, error), *atomic.Int64) {
	var calls atomic.Int64
	factory := func(context.Context) (*int, error) {
		n := int(calls.Add(1))
		return &n, nil
	}
	return factory, &calls
}

func TestNew_Validation(t *testing.T) {
	okFactory := func(context.Context) (int, error) { return 0, nil }

	tests := []struct {
		name    string
		size    int
		factory func(context.Context) (int, error)
		wantErr error
	}{
		{"valid", 4, okFactory, nil},
		{"size zero", 0, okFactory, ErrInvalidSize},
		{"size negative", -3, okFactory, ErrInvalidSize},
		{"nil factory", 2, nil, ErrNoFactory},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New[int](tt.size, tt.factory, nil)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("New error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if p == nil {
					t.Fatal("New returned nil pool with nil error")
				}
				if got := p.Size(); got != tt.size {
					t.Errorf("Size() = %d, want %d", got, tt.size)
				}
				if got := p.Available(); got != 0 {
					t.Errorf("Available() = %d, want 0 (lazy creation)", got)
				}
			}
		})
	}
}

func TestGetPut_LazyCreationAndReuse(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"size 1", 1},
		{"size 3", 3},
		{"size 8", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, calls := newCountingFactory()
			p, err := New[*int](tt.size, factory, nil)
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			ctx := context.Background()

			// Get one object, return it, get again: the factory must run only
			// once because the second Get reuses the idle object.
			obj1, err := p.Get(ctx)
			if err != nil {
				t.Fatalf("Get #1: %v", err)
			}
			if got := calls.Load(); got != 1 {
				t.Fatalf("factory calls after first Get = %d, want 1", got)
			}
			if err := p.Put(obj1); err != nil {
				t.Fatalf("Put: %v", err)
			}
			if got := p.Available(); got != 1 {
				t.Fatalf("Available after Put = %d, want 1", got)
			}

			obj2, err := p.Get(ctx)
			if err != nil {
				t.Fatalf("Get #2: %v", err)
			}
			if got := calls.Load(); got != 1 {
				t.Errorf("factory calls after reuse = %d, want 1 (no new object)", got)
			}
			if obj2 != obj1 {
				t.Errorf("Get #2 returned a different object; expected reuse")
			}
		})
	}
}

func TestGet_BlocksAtCapacityThenUnblocksOnPut(t *testing.T) {
	factory, calls := newCountingFactory()
	p, err := New[*int](2, factory, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ctx := context.Background()

	// Exhaust the pool.
	a, err := p.Get(ctx)
	if err != nil {
		t.Fatalf("Get a: %v", err)
	}
	b, err := p.Get(ctx)
	if err != nil {
		t.Fatalf("Get b: %v", err)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("factory calls = %d, want 2", got)
	}

	// A third Get must block until something is returned. Use a channel to
	// observe completion deterministically (no time.Sleep).
	gotCh := make(chan *int, 1)
	go func() {
		obj, gErr := p.Get(ctx)
		if gErr != nil {
			t.Errorf("blocked Get: %v", gErr)
		}
		gotCh <- obj
	}()

	// The goroutine must still be blocked: nothing on gotCh yet.
	select {
	case <-gotCh:
		t.Fatal("Get returned while pool was exhausted")
	default:
	}

	// Returning an object unblocks the waiter and hands it the same object.
	if err := p.Put(a); err != nil {
		t.Fatalf("Put a: %v", err)
	}
	got := <-gotCh
	if got != a {
		t.Errorf("blocked Get returned %v, want returned object %v", got, a)
	}
	if calls.Load() != 2 {
		t.Errorf("factory ran again while unblocking; calls = %d, want 2", calls.Load())
	}

	_ = b
}

func TestGet_ContextCancellationWhileBlocked(t *testing.T) {
	factory, _ := newCountingFactory()
	p, err := New[*int](1, factory, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Take the only object so the next Get blocks.
	held, err := p.Get(context.Background())
	if err != nil {
		t.Fatalf("Get held: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		_, gErr := p.Get(ctx)
		errCh <- gErr
	}()

	// Cancel and verify the blocked Get returns a wrapped context error.
	cancel()
	gErr := <-errCh
	if !errors.Is(gErr, context.Canceled) {
		t.Fatalf("Get error = %v, want context.Canceled", gErr)
	}

	// The held object is still valid and returnable.
	if err := p.Put(held); err != nil {
		t.Errorf("Put after cancel: %v", err)
	}
}

func TestGet_FactoryErrorDoesNotLeakCapacity(t *testing.T) {
	errBoom := errors.New("boom")
	var attempts atomic.Int64
	factory := func(context.Context) (*int, error) {
		// First call fails, later calls succeed.
		if attempts.Add(1) == 1 {
			return nil, errBoom
		}
		n := 42
		return &n, nil
	}

	p, err := New[*int](1, factory, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ctx := context.Background()

	_, err = p.Get(ctx)
	if !errors.Is(err, errBoom) {
		t.Fatalf("Get error = %v, want wrapped %v", err, errBoom)
	}

	// The failed slot must have been released: a subsequent Get can still
	// create the one allowed object. If capacity leaked, this would block
	// forever / require the pool to think it is full.
	got, err := p.Get(ctx)
	if err != nil {
		t.Fatalf("Get after factory error: %v", err)
	}
	if *got != 42 {
		t.Errorf("got *%d, want 42", *got)
	}
}

func TestPut_ResetHookRuns(t *testing.T) {
	resetCalls := 0
	factory := func(context.Context) (*bytes.Buffer, error) {
		return &bytes.Buffer{}, nil
	}
	reset := func(b *bytes.Buffer) {
		resetCalls++
		b.Reset()
	}

	p, err := New[*bytes.Buffer](1, factory, reset)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ctx := context.Background()

	buf, err := p.Get(ctx)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	buf.WriteString("dirty state")
	if err := p.Put(buf); err != nil {
		t.Fatalf("Put: %v", err)
	}
	if resetCalls != 1 {
		t.Fatalf("reset called %d times, want 1", resetCalls)
	}

	reused, err := p.Get(ctx)
	if err != nil {
		t.Fatalf("Get #2: %v", err)
	}
	if reused.Len() != 0 {
		t.Errorf("reused buffer not scrubbed: len = %d, want 0", reused.Len())
	}
}

func TestClose(t *testing.T) {
	t.Run("get and put after close", func(t *testing.T) {
		factory, _ := newCountingFactory()
		p, err := New[*int](2, factory, nil)
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		if err := p.Close(nil); err != nil {
			t.Fatalf("Close: %v", err)
		}
		if _, err := p.Get(context.Background()); !errors.Is(err, ErrPoolClosed) {
			t.Errorf("Get after Close = %v, want ErrPoolClosed", err)
		}
		v := 1
		if err := p.Put(&v); !errors.Is(err, ErrPoolClosed) {
			t.Errorf("Put after Close = %v, want ErrPoolClosed", err)
		}
	})

	t.Run("disposes idle objects", func(t *testing.T) {
		factory, _ := newCountingFactory()
		p, err := New[*int](3, factory, nil)
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		ctx := context.Background()

		// Create and return two objects so they sit idle in the free list.
		a, _ := p.Get(ctx)
		b, _ := p.Get(ctx)
		_ = p.Put(a)
		_ = p.Put(b)

		var disposed atomic.Int64
		if err := p.Close(func(*int) { disposed.Add(1) }); err != nil {
			t.Fatalf("Close: %v", err)
		}
		if got := disposed.Load(); got != 2 {
			t.Errorf("disposed %d objects, want 2", got)
		}
	})

	t.Run("idempotent", func(t *testing.T) {
		factory, _ := newCountingFactory()
		p, err := New[*int](1, factory, nil)
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		if err := p.Close(nil); err != nil {
			t.Fatalf("first Close: %v", err)
		}
		if err := p.Close(nil); err != nil {
			t.Errorf("second Close = %v, want nil (idempotent)", err)
		}
	})
}

// TestConcurrent_NeverExceedsCapacity is the core race-safety check: many
// goroutines hammer Get/Put concurrently, and at no point may the number of
// simultaneously checked-out objects exceed the pool size. It also confirms the
// factory is never invoked more than `size` times.
func TestConcurrent_NeverExceedsCapacity(t *testing.T) {
	tests := []struct {
		name       string
		size       int
		goroutines int
		opsEach    int
	}{
		{"small pool high contention", 2, 32, 200},
		{"medium pool", 8, 16, 300},
		{"single slot", 1, 20, 150},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, calls := newCountingFactory()
			p, err := New[*int](tt.size, factory, nil)
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			ctx := context.Background()

			var inUse atomic.Int64
			var maxInUse atomic.Int64

			var wg sync.WaitGroup
			wg.Add(tt.goroutines)
			for g := 0; g < tt.goroutines; g++ {
				go func() {
					defer wg.Done()
					for i := 0; i < tt.opsEach; i++ {
						obj, gErr := p.Get(ctx)
						if gErr != nil {
							t.Errorf("Get: %v", gErr)
							return
						}
						cur := inUse.Add(1)
						// Track the high-water mark of concurrent checkouts.
						for {
							old := maxInUse.Load()
							if cur <= old || maxInUse.CompareAndSwap(old, cur) {
								break
							}
						}
						if cur > int64(tt.size) {
							t.Errorf("in-use %d exceeded pool size %d", cur, tt.size)
						}
						inUse.Add(-1)
						if pErr := p.Put(obj); pErr != nil {
							t.Errorf("Put: %v", pErr)
							return
						}
					}
				}()
			}
			wg.Wait()

			if got := calls.Load(); got > int64(tt.size) {
				t.Errorf("factory invoked %d times, must not exceed size %d", got, tt.size)
			}
			if got := maxInUse.Load(); got > int64(tt.size) {
				t.Errorf("max concurrent in-use = %d, want <= %d", got, tt.size)
			}
			if got := p.Available(); got > tt.size {
				t.Errorf("Available = %d, want <= %d", got, tt.size)
			}
		})
	}
}

// Example demonstrates the typical Get/defer-Put lifecycle.
func Example() {
	p, err := New[*bytes.Buffer](
		4,
		func(context.Context) (*bytes.Buffer, error) { return &bytes.Buffer{}, nil },
		func(b *bytes.Buffer) { b.Reset() },
	)
	if err != nil {
		panic(err)
	}
	defer p.Close(nil)

	buf, err := p.Get(context.Background())
	if err != nil {
		panic(err)
	}
	defer p.Put(buf)

	buf.WriteString("hello pool")
	fmt.Println(buf.String())
	// Output: hello pool
}
