// Package prodcons demonstrates the Producer-Consumer pattern:
// N producers send work items into a bounded channel; M consumers process them.
// Context cancellation and graceful shutdown are fully supported.
package prodcons

import (
	"context"
	"sync"
)

// Item is the unit of work passed from producers to consumers.
type Item struct {
	ProducerID int
	Value      int
}

// Result holds the outcome of processing one item.
type Result struct {
	Item  Item
	Value int // transformed value
}

// Pipeline runs nProducers and nConsumers goroutines connected by a bounded
// channel of size bufSize.
//
// produce(id, yield) is called in each producer goroutine; it should call
// yield(item) for each work item. yield returns false when ctx is cancelled.
//
// consume(item) is called for each item by consumer goroutines. Its return
// value is collected and returned.
//
// Pipeline blocks until all items are processed or ctx is cancelled.
func Pipeline(
	ctx context.Context,
	nProducers, nConsumers, bufSize int,
	produce func(id int, yield func(Item) bool),
	consume func(Item) Result,
) []Result {
	jobs := make(chan Item, bufSize)
	results := make(chan Result, bufSize)

	// --- Producers ---
	var prodWg sync.WaitGroup
	prodWg.Add(nProducers)
	for i := range nProducers {
		go func(id int) {
			defer prodWg.Done()
			produce(id, func(item Item) bool {
				select {
				case jobs <- item:
					return true
				case <-ctx.Done():
					return false
				}
			})
		}(i)
	}
	// Close jobs exactly once when all producers finish.
	go func() {
		prodWg.Wait()
		close(jobs)
	}()

	// --- Consumers ---
	var consWg sync.WaitGroup
	consWg.Add(nConsumers)
	for range nConsumers {
		go func() {
			defer consWg.Done()
			for item := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}
				r := consume(item)
				select {
				case results <- r:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Close results when all consumers finish.
	go func() {
		consWg.Wait()
		close(results)
	}()

	// Collect all results.
	var out []Result
	for r := range results {
		out = append(out, r)
	}
	return out
}
