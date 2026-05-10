package rwlock

import (
	"fmt"
	"sync"
	"testing"
)

func TestConfigStore_SetGet(t *testing.T) {
	c := NewConfigStore()
	c.Set("host", "localhost")
	c.Set("port", "8080")

	tests := []struct {
		key  string
		want string
		ok   bool
	}{
		{"host", "localhost", true},
		{"port", "8080", true},
		{"missing", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, ok := c.Get(tt.key)
			if ok != tt.ok {
				t.Errorf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Errorf("value = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConfigStore_Delete(t *testing.T) {
	c := NewConfigStore()
	c.Set("key", "val")
	c.Delete("key")
	if _, ok := c.Get("key"); ok {
		t.Error("key should be absent after Delete")
	}
}

func TestConfigStore_Snapshot(t *testing.T) {
	c := NewConfigStore()
	for i := range 5 {
		c.Set(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i))
	}
	snap := c.Snapshot()
	if len(snap) != 5 {
		t.Errorf("snapshot len = %d, want 5", len(snap))
	}
	// Modifying snapshot must not affect the store.
	snap["extra"] = "x"
	if c.Len() != 5 {
		t.Error("store mutated via snapshot")
	}
}

func TestConfigStore_ConcurrentReads(t *testing.T) {
	c := NewConfigStore()
	c.Set("key", "value")

	var wg sync.WaitGroup
	const readers = 100
	wg.Add(readers)
	for range readers {
		go func() {
			defer wg.Done()
			v, ok := c.Get("key")
			if !ok || v != "value" {
				t.Errorf("Get returned (%q, %v)", v, ok)
			}
		}()
	}
	wg.Wait()
}

func TestConfigStore_ConcurrentWriters(t *testing.T) {
	c := NewConfigStore()
	const writers = 50
	var wg sync.WaitGroup
	wg.Add(writers)
	for i := range writers {
		go func(i int) {
			defer wg.Done()
			c.Set(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i))
		}(i)
	}
	wg.Wait()
	if c.Len() != writers {
		t.Errorf("Len = %d, want %d", c.Len(), writers)
	}
}

func TestConfigStore_MixedReadWrite_Race(t *testing.T) {
	c := NewConfigStore()
	c.Set("shared", "initial")

	var wg sync.WaitGroup
	const (
		readers = 80
		writers = 20
	)
	wg.Add(readers + writers)

	for range readers {
		go func() {
			defer wg.Done()
			c.Get("shared")
		}()
	}
	for i := range writers {
		go func(i int) {
			defer wg.Done()
			c.Set("shared", fmt.Sprintf("update-%d", i))
		}(i)
	}
	wg.Wait()
}
