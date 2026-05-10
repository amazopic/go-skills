package steadystate

import (
	"context"
	"sync"
	"testing"
	"time"
)

// --- Cache tests ------------------------------------------------------------

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache(10, time.Second)
	ok := c.Set("key", "value", time.Minute)
	if !ok {
		t.Fatal("Set should succeed for non-full cache")
	}
	v, found := c.Get("key")
	if !found {
		t.Fatal("Get should find existing key")
	}
	if v.(string) != "value" {
		t.Errorf("want 'value', got %v", v)
	}
}

func TestCache_GetMissing(t *testing.T) {
	c := NewCache(10, time.Second)
	_, found := c.Get("missing")
	if found {
		t.Error("Get should return false for missing key")
	}
}

func TestCache_GetExpired(t *testing.T) {
	c := NewCache(10, time.Second)
	c.Set("key", "val", time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	_, found := c.Get("key")
	if found {
		t.Error("Get should return false for expired key")
	}
}

func TestCache_Update_AlwaysAllowed(t *testing.T) {
	c := NewCache(1, time.Second) // maxSize=1
	c.Set("key", "v1", time.Minute)
	ok := c.Set("key", "v2", time.Minute) // update existing
	if !ok {
		t.Error("update of existing key should always succeed")
	}
	v, _ := c.Get("key")
	if v.(string) != "v2" {
		t.Errorf("want 'v2', got %v", v)
	}
}

func TestCache_CapacityEvictsExpiredOnInsert(t *testing.T) {
	c := NewCache(2, time.Second)
	c.Set("a", 1, time.Millisecond) // will expire
	c.Set("b", 2, time.Minute)
	time.Sleep(5 * time.Millisecond)

	// "a" expired; inserting "c" should evict "a" and succeed.
	ok := c.Set("c", 3, time.Minute)
	if !ok {
		t.Error("insert should succeed after evicting expired entry")
	}
	if c.Len() > 2 {
		t.Errorf("cache exceeds maxSize: len=%d", c.Len())
	}
}

func TestCache_CapacityShedsWhenFull(t *testing.T) {
	c := NewCache(2, time.Second)
	c.Set("a", 1, time.Minute)
	c.Set("b", 2, time.Minute)
	// No expired entries — insert should be shed.
	ok := c.Set("c", 3, time.Minute)
	if ok {
		t.Error("insert should be shed when cache is full and no expired entries")
	}
	if c.Len() != 2 {
		t.Errorf("want len=2, got %d", c.Len())
	}
}

func TestCache_EvictExpired_Sweeps(t *testing.T) {
	c := NewCache(10, time.Second)
	c.Set("short", 1, time.Millisecond)
	c.Set("long", 2, time.Minute)
	time.Sleep(5 * time.Millisecond)

	c.EvictExpired()

	if c.Len() != 1 {
		t.Errorf("want 1 entry after sweep, got %d", c.Len())
	}
	if _, found := c.Get("long"); !found {
		t.Error("long-lived entry should still be present")
	}
}

func TestCache_Janitor_EvictsOverTime(t *testing.T) {
	c := NewCache(10, 10*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c.Start(ctx)

	c.Set("ephemeral", "x", 5*time.Millisecond)
	time.Sleep(30 * time.Millisecond) // let janitor run at least once

	c.Stop()

	if c.Len() != 0 {
		t.Errorf("janitor should have evicted expired entry, len=%d", c.Len())
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := NewCache(50, 10*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c.Start(ctx)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := string(rune('a' + n%26))
			c.Set(key, n, 20*time.Millisecond)
			c.Get(key)
		}(i)
	}
	wg.Wait()
	c.Stop()
}

// --- RingBuffer tests -------------------------------------------------------

func TestRingBuffer_BasicPushSlice(t *testing.T) {
	r := NewRingBuffer[int](3)
	r.Push(1)
	r.Push(2)
	r.Push(3)
	got := r.Slice()
	want := []int{1, 2, 3}
	if len(got) != len(want) {
		t.Fatalf("want len %d, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: want %d, got %d", i, want[i], got[i])
		}
	}
}

func TestRingBuffer_OverflowDropsOldest(t *testing.T) {
	r := NewRingBuffer[int](3)
	for _, v := range []int{1, 2, 3, 4, 5} {
		r.Push(v)
	}
	got := r.Slice()
	want := []int{3, 4, 5}
	if len(got) != len(want) {
		t.Fatalf("want len %d, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: want %d, got %d", i, want[i], got[i])
		}
	}
}

func TestRingBuffer_LenBounded(t *testing.T) {
	r := NewRingBuffer[string](5)
	for i := 0; i < 100; i++ {
		r.Push("x")
	}
	if r.Len() != 5 {
		t.Errorf("want Len=5 (capacity bound), got %d", r.Len())
	}
}

func TestRingBuffer_ConcurrentPush(t *testing.T) {
	r := NewRingBuffer[int](16)
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			r.Push(v)
		}(i)
	}
	wg.Wait()
	if r.Len() != 16 {
		t.Errorf("want Len=16 after concurrent pushes, got %d", r.Len())
	}
}
