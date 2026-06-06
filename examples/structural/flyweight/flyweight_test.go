package flyweight

import (
	"strconv"
	"sync"
	"testing"
)

func TestFlyweight(t *testing.T) {
	var testCases = []struct {
		filename string
		width    int
		height   int
		opacity  float64
		expect   string
	}{
		{"cat.jpg", 100, 100, 0.95, "draw image: cat.jpg, width: 100, height: 100, opacity: 0.95"},
		{"cat.jpg", 200, 200, 0.75, "draw image: cat.jpg, width: 200, height: 200, opacity: 0.75"},
		{"dog.jpg", 300, 300, 0.50, "draw image: dog.jpg, width: 300, height: 300, opacity: 0.50"},
	}

	var factory = new(FlyweightFactory)

	for i, tt := range testCases {
		t.Run("case "+strconv.Itoa(i), func(t *testing.T) {
			var flyweight = factory.GetFlyweight(tt.filename)
			var result = flyweight.Draw(tt.width, tt.height, tt.opacity)
			if result != tt.expect {
				t.Errorf("Expect result to equal %s, but %s.\n", tt.expect, result)
			}
		})
	}
}

// TestFlyweightConcurrent hammers GetFlyweight from many goroutines on both the
// same and different keys. It asserts the factory's interning guarantee — every
// request for a given key must return the *same* pointer — while exercising the
// shared cache hard enough to trip the race detector if the map were unguarded.
func TestFlyweightConcurrent(t *testing.T) {
	const (
		goroutines = 64
		iterations = 200
	)

	// A small key set so many goroutines collide on the same entries (forcing
	// concurrent read/write of the same map slots), plus distinct keys to drive
	// concurrent inserts of brand-new entries.
	keys := []string{"cat.jpg", "dog.jpg", "fox.jpg", "owl.jpg"}

	// Factory is intentionally zero-valued so the lazy pool init also races
	// across goroutines on the very first calls.
	var factory FlyweightFactory

	// Collect, per key, the set of distinct pointers observed. Interning means
	// each key must map to exactly one pointer no matter the concurrency.
	var mu sync.Mutex
	seen := make(map[string]map[Flyweighter]struct{}, len(keys))
	for _, k := range keys {
		seen[k] = make(map[Flyweighter]struct{})
	}

	var start sync.WaitGroup
	start.Add(1)
	var done sync.WaitGroup
	done.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(g int) {
			defer done.Done()
			start.Wait() // release all goroutines at once to maximize contention
			for i := 0; i < iterations; i++ {
				key := keys[(g+i)%len(keys)]
				fw := factory.GetFlyweight(key)
				if fw == nil {
					t.Errorf("GetFlyweight(%q) returned nil", key)
					return
				}
				mu.Lock()
				seen[key][fw] = struct{}{}
				mu.Unlock()
			}
		}(g)
	}

	start.Done()
	done.Wait()

	for _, k := range keys {
		if n := len(seen[k]); n != 1 {
			t.Errorf("interning violated for key %q: observed %d distinct flyweights, want 1", k, n)
		}
	}
}
