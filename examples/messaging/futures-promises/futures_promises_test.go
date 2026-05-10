package futures

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestFuture_Resolve(t *testing.T) {
	p, f := New[int]()
	go p.Resolve(42)
	v, err := f.Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 42 {
		t.Errorf("Get() = %d, want 42", v)
	}
}

func TestFuture_Reject(t *testing.T) {
	want := errors.New("boom")
	p, f := New[string]()
	go p.Reject(want)
	_, err := f.Get(context.Background())
	if !errors.Is(err, want) {
		t.Errorf("err = %v, want %v", err, want)
	}
}

func TestFuture_MultipleGetSameResult(t *testing.T) {
	p, f := New[int]()
	go p.Resolve(7)

	// Wait until resolved.
	v1, err1 := f.Get(context.Background())
	if err1 != nil || v1 != 7 {
		t.Fatalf("first Get = (%d, %v)", v1, err1)
	}

	// Subsequent calls should return the cached result.
	const callers = 20
	var wg sync.WaitGroup
	wg.Add(callers)
	for range callers {
		go func() {
			defer wg.Done()
			v, err := f.Get(context.Background())
			if err != nil || v != 7 {
				t.Errorf("concurrent Get = (%d, %v)", v, err)
			}
		}()
	}
	wg.Wait()
}

func TestFuture_ContextCancellation(t *testing.T) {
	_, f := New[int]() // promise never resolved

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err := f.Get(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want DeadlineExceeded", err)
	}
}

func TestGo_Success(t *testing.T) {
	f := Go[string](context.Background(), func(ctx context.Context) (string, error) {
		return "hello", nil
	})
	v, err := f.Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "hello" {
		t.Errorf("Get() = %q, want %q", v, "hello")
	}
}

func TestGo_Error(t *testing.T) {
	want := errors.New("fail")
	f := Go[int](context.Background(), func(ctx context.Context) (int, error) {
		return 0, want
	})
	_, err := f.Get(context.Background())
	if !errors.Is(err, want) {
		t.Errorf("err = %v, want %v", err, want)
	}
}

func TestFuture_ResolveOnlyOnce(t *testing.T) {
	// Calling Resolve multiple times must not panic or block.
	p, f := New[int]()
	go func() {
		p.Resolve(1)
		p.Resolve(2) // should be a no-op
		p.Resolve(3) // should be a no-op
	}()
	v, err := f.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if v != 1 {
		t.Errorf("v = %d, want 1", v)
	}
}
