package prototype

import (
	"sync"
	"testing"
)

func TestPrototype(t *testing.T) {

	product := NewConcreteProduct("A")
	cloneProduct := product.Clone()

	if cloneProduct.GetName() != product.GetName() {
		t.Error("Expect name \"A\" to equal, but not equal.")
	}
}

// TestPrototypeConcurrentClone exercises the clone path from many goroutines
// against a single shared prototype. Clone must be safe for concurrent use
// (it only reads the prototype's state), so this must stay green under -race.
func TestPrototypeConcurrentClone(t *testing.T) {
	prototype := NewConcreteProduct("A")

	const goroutines = 64
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				clone := prototype.Clone()
				if clone.GetName() != "A" {
					t.Errorf("clone name = %q, want %q", clone.GetName(), "A")
					return
				}
			}
		}()
	}
	wg.Wait()
}
