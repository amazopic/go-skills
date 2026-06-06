package bridge

import (
	"sync"
	"testing"
)

func TestBridge(t *testing.T) {

	expect := "SssuuuuZzzuuuuKkiiiii"

	car := NewCar(&EngineSuzuki{})

	sound := car.Race()

	if sound != expect {
		t.Errorf("Expect sound to %s, but %s", expect, sound)
	}
}

// TestBridge_ConcurrentRace exercises the renamed Race() method from many
// goroutines so the shared Car (and its injected Enginer) is hammered under
// the -race detector. A correct, immutable bridge field must stay race-free.
func TestBridge_ConcurrentRace(t *testing.T) {
	expect := "HhoooNnnnnnnnnDddaaaaaaa"

	car := NewCar(&EngineHonda{})

	const goroutines = 64
	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				if sound := car.Race(); sound != expect {
					t.Errorf("Expect sound to %s, but %s", expect, sound)
					return
				}
			}
		}()
	}

	wg.Wait()
}
