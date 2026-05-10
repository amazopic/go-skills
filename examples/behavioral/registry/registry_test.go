package registry

import (
	"fmt"
	"sync"
	"testing"
)

// Codec is a sample type used to demonstrate a typed registry.
type Codec struct {
	Name string
}

func TestRegisterAndGet(t *testing.T) {
	var r Registry[*Codec]

	c := &Codec{Name: "json"}
	if err := r.Register("json", c); err != nil {
		t.Fatalf("Register: unexpected error: %v", err)
	}

	got, ok := r.Get("json")
	if !ok {
		t.Fatal("Get: expected to find \"json\", got not-found")
	}
	if got != c {
		t.Errorf("Get: got %v, want %v", got, c)
	}
}

func TestRegisterDuplicate(t *testing.T) {
	var r Registry[int]

	if err := r.Register("key", 1); err != nil {
		t.Fatalf("first Register: unexpected error: %v", err)
	}
	if err := r.Register("key", 2); err == nil {
		t.Fatal("second Register: expected an error for duplicate key, got nil")
	}
}

func TestGetMissing(t *testing.T) {
	var r Registry[string]

	_, ok := r.Get("nonexistent")
	if ok {
		t.Error("Get: expected not-found for unregistered key")
	}
}

func TestMustRegisterPanicsOnDuplicate(t *testing.T) {
	var r Registry[int]
	r.MustRegister("x", 42)

	defer func() {
		if rec := recover(); rec == nil {
			t.Error("MustRegister: expected panic on duplicate, got none")
		}
	}()
	r.MustRegister("x", 99)
}

func TestMustGetPanicsOnMissing(t *testing.T) {
	var r Registry[int]

	defer func() {
		if rec := recover(); rec == nil {
			t.Error("MustGet: expected panic on missing key, got none")
		}
	}()
	r.MustGet("missing")
}

func TestMustGet(t *testing.T) {
	var r Registry[string]
	r.MustRegister("lang", "go")

	if got := r.MustGet("lang"); got != "go" {
		t.Errorf("MustGet: got %q, want %q", got, "go")
	}
}

func TestNamesAndLen(t *testing.T) {
	var r Registry[int]
	keys := []string{"a", "b", "c"}
	for i, k := range keys {
		r.MustRegister(k, i)
	}

	if r.Len() != len(keys) {
		t.Errorf("Len: got %d, want %d", r.Len(), len(keys))
	}

	names := r.Names()
	if len(names) != len(keys) {
		t.Errorf("Names: got %d names, want %d", len(names), len(keys))
	}
	nameSet := make(map[string]bool, len(names))
	for _, n := range names {
		nameSet[n] = true
	}
	for _, k := range keys {
		if !nameSet[k] {
			t.Errorf("Names: missing expected key %q", k)
		}
	}
}

// TestConcurrentRegisterAndGet exercises the registry under concurrent access
// to confirm it is race-safe (run with -race).
func TestConcurrentRegisterAndGet(t *testing.T) {
	const workers = 20
	const entriesPerWorker = 50

	var r Registry[int]
	var wg sync.WaitGroup

	// Phase 1: concurrent Registers (each worker uses a unique key namespace).
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < entriesPerWorker; i++ {
				key := fmt.Sprintf("w%d-k%d", workerID, i)
				_ = r.Register(key, workerID*entriesPerWorker+i)
			}
		}(w)
	}
	wg.Wait()

	// Phase 2: concurrent Gets from all workers simultaneously.
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < entriesPerWorker; i++ {
				key := fmt.Sprintf("w%d-k%d", workerID, i)
				val, ok := r.Get(key)
				if !ok {
					t.Errorf("Get(%q): key not found after concurrent Register", key)
					return
				}
				want := workerID*entriesPerWorker + i
				if val != want {
					t.Errorf("Get(%q): got %d, want %d", key, val, want)
				}
			}
		}(w)
	}
	wg.Wait()

	want := workers * entriesPerWorker
	if r.Len() != want {
		t.Errorf("Len: got %d, want %d", r.Len(), want)
	}
}
