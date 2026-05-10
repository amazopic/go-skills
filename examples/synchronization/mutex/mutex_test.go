package mutex

import (
	"sync"
	"testing"
)

func TestSafeMap_SetGet(t *testing.T) {
	m := NewSafeMap()
	m.Set("a", 1)
	m.Set("b", 2)

	tests := []struct {
		key  string
		want int
		ok   bool
	}{
		{"a", 1, true},
		{"b", 2, true},
		{"c", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, ok := m.Get(tt.key)
			if ok != tt.ok {
				t.Errorf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Errorf("value = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSafeMap_Inc_Race(t *testing.T) {
	m := NewSafeMap()
	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			m.Inc("counter")
		}()
	}
	wg.Wait()

	got, ok := m.Get("counter")
	if !ok {
		t.Fatal("key 'counter' missing")
	}
	if got != goroutines {
		t.Errorf("counter = %d, want %d", got, goroutines)
	}
}

func TestSafeMap_Len(t *testing.T) {
	m := NewSafeMap()
	for _, k := range []string{"x", "y", "z"} {
		m.Set(k, 1)
	}
	if got := m.Len(); got != 3 {
		t.Errorf("Len = %d, want 3", got)
	}
}

func TestReadHeavyCache_PutLookup(t *testing.T) {
	c := NewReadHeavyCache()
	c.Put("lang", "Go")

	tests := []struct {
		key  string
		want string
		ok   bool
	}{
		{"lang", "Go", true},
		{"missing", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, ok := c.Lookup(tt.key)
			if ok != tt.ok {
				t.Errorf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Errorf("value = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadHeavyCache_ConcurrentReads(t *testing.T) {
	c := NewReadHeavyCache()
	c.Put("k", "v")

	var wg sync.WaitGroup
	const readers = 50
	wg.Add(readers)
	for range readers {
		go func() {
			defer wg.Done()
			v, ok := c.Lookup("k")
			if !ok || v != "v" {
				t.Errorf("Lookup returned (%q, %v)", v, ok)
			}
		}()
	}
	wg.Wait()
}

func TestReadHeavyCache_WriterExclusive(t *testing.T) {
	c := NewReadHeavyCache()
	var wg sync.WaitGroup
	const (
		writers = 10
		readers = 40
	)
	wg.Add(writers + readers)
	for i := range writers {
		go func(i int) {
			defer wg.Done()
			_ = i
			c.Put("shared", "value")
		}(i)
	}
	for range readers {
		go func() {
			defer wg.Done()
			c.Lookup("shared")
		}()
	}
	wg.Wait()
}
