package condvar

import (
	"sort"
	"sync"
	"testing"
)

func TestBoundedQueue_PutGet(t *testing.T) {
	q := NewBoundedQueue(4)
	for i := range 4 {
		q.Put(i)
	}
	for i := range 4 {
		if got := q.Get(); got != i {
			t.Errorf("Get() = %d, want %d", got, i)
		}
	}
}

func TestBoundedQueue_Len(t *testing.T) {
	q := NewBoundedQueue(10)
	if q.Len() != 0 {
		t.Errorf("initial Len = %d, want 0", q.Len())
	}
	q.Put(1)
	q.Put(2)
	if q.Len() != 2 {
		t.Errorf("Len = %d, want 2", q.Len())
	}
}

func TestBoundedQueue_ConcurrentProducersConsumers(t *testing.T) {
	const (
		cap      = 8
		items    = 200
		nProd    = 4
		nCons    = 4
	)
	q := NewBoundedQueue(cap)
	results := make(chan int, items)

	var consWg sync.WaitGroup
	consWg.Add(nCons)
	for range nCons {
		go func() {
			defer consWg.Done()
			for range items / nCons {
				results <- q.Get()
			}
		}()
	}

	var prodWg sync.WaitGroup
	prodWg.Add(nProd)
	for p := range nProd {
		go func(id int) {
			defer prodWg.Done()
			for i := range items / nProd {
				q.Put(id*(items/nProd) + i)
			}
		}(p)
	}

	prodWg.Wait()
	consWg.Wait()
	close(results)

	var got []int
	for v := range results {
		got = append(got, v)
	}
	sort.Ints(got)
	if len(got) != items {
		t.Errorf("received %d items, want %d", len(got), items)
	}
	// Verify each value is in expected range.
	for _, v := range got {
		if v < 0 || v >= items {
			t.Errorf("unexpected value %d", v)
		}
	}
}

func TestBoundedQueue_Drain(t *testing.T) {
	q := NewBoundedQueue(10)
	for i := range 5 {
		q.Put(i)
	}
	drained := q.Drain()
	if len(drained) != 5 {
		t.Errorf("Drain returned %d items, want 5", len(drained))
	}
	if q.Len() != 0 {
		t.Errorf("Len after Drain = %d, want 0", q.Len())
	}
}

func TestBoundedQueue_PanicOnZeroCap(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for cap=0")
		}
	}()
	NewBoundedQueue(0)
}
