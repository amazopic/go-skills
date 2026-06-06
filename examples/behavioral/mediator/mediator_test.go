package mediator

import (
	"sync"
	"testing"
)

func TestMediator(t *testing.T) {

	farmer := new(Farmer)
	cannery := new(Cannery)
	shop := new(Shop)

	farmer.AddMoney(7500.00)
	cannery.AddMoney(15000.00)
	shop.AddMoney(30000.00)

	ConnectColleagues(farmer, cannery, shop)

	// A farmer grows a 1000kg tomato
	// and informs the mediator about the completion of his work.
	// Next, the mediator sends the tomatoes to the cannery.
	// After the cannery produces 1000 packs of ketchup,
	// he informs the mediator about his delivery to the store.
	farmer.GrowTomato(1000)

	expect := float64(54750)
	result := shop.GetMoney()

	if result != expect {
		t.Errorf("Expect result to equal %f, but %f.\n", expect, result)
	}
}

// TestConnectColleaguesConcurrent calls the public ConnectColleagues by its
// Latin name (the bug was a Cyrillic homoglyph "С" U+0421 that made the
// documented Latin name uncallable; this test source would not even compile
// against the pre-fix code). Each goroutine wires and drives a completely
// independent set of colleagues, so the Mediator wiring is exercised many
// times in parallel under -race.
func TestConnectColleaguesConcurrent(t *testing.T) {
	const goroutines = 64

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			farmer := new(Farmer)
			cannery := new(Cannery)
			shop := new(Shop)

			farmer.AddMoney(7500.00)
			cannery.AddMoney(15000.00)
			shop.AddMoney(30000.00)

			ConnectColleagues(farmer, cannery, shop)

			farmer.GrowTomato(1000)

			if got, want := shop.GetMoney(), float64(54750); got != want {
				t.Errorf("Expect result to equal %f, but %f.\n", want, got)
			}
		}()
	}
	wg.Wait()
}
