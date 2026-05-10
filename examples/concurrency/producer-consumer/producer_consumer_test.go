package prodcons

import (
	"context"
	"sort"
	"testing"
)

func TestPipeline_AllItemsProcessed(t *testing.T) {
	tests := []struct {
		name      string
		producers int
		consumers int
		bufSize   int
		itemsEach int
	}{
		{"1P1C", 1, 1, 4, 10},
		{"3P2C", 3, 2, 8, 20},
		{"1P8C", 1, 8, 16, 50},
		{"8P1C", 8, 1, 16, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			totalItems := tt.producers * tt.itemsEach

			results := Pipeline(
				ctx,
				tt.producers,
				tt.consumers,
				tt.bufSize,
				func(id int, yield func(Item) bool) {
					for i := range tt.itemsEach {
						if !yield(Item{ProducerID: id, Value: i}) {
							return
						}
					}
				},
				func(item Item) Result {
					return Result{Item: item, Value: item.Value * 2}
				},
			)

			if len(results) != totalItems {
				t.Errorf("got %d results, want %d", len(results), totalItems)
			}
			// Each result.Value should equal item.Value * 2.
			for _, r := range results {
				if r.Value != r.Item.Value*2 {
					t.Errorf("result.Value = %d, want %d", r.Value, r.Item.Value*2)
				}
			}
		})
	}
}

func TestPipeline_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	resultsCh := make(chan []Result, 1)

	go func() {
		r := Pipeline(
			ctx,
			2, 2, 4,
			func(id int, yield func(Item) bool) {
				// Infinite producer; stops when yield returns false (ctx cancelled).
				for i := 0; ; i++ {
					if !yield(Item{ProducerID: id, Value: i}) {
						return
					}
				}
			},
			func(item Item) Result {
				return Result{Item: item, Value: item.Value}
			},
		)
		resultsCh <- r
	}()

	// Cancel immediately; pipeline must terminate without deadlock.
	cancel()
	r := <-resultsCh
	// Result count is indeterminate — we only verify termination.
	_ = r
}

func TestPipeline_OrderIndependent(t *testing.T) {
	ctx := context.Background()
	const n = 100

	results := Pipeline(
		ctx,
		1, 4, 10,
		func(id int, yield func(Item) bool) {
			for i := range n {
				if !yield(Item{ProducerID: id, Value: i}) {
					return
				}
			}
		},
		func(item Item) Result {
			return Result{Item: item, Value: item.Value}
		},
	)

	if len(results) != n {
		t.Fatalf("got %d results, want %d", len(results), n)
	}
	// Sort and verify all values 0..n-1 appear.
	sort.Slice(results, func(i, j int) bool { return results[i].Value < results[j].Value })
	for i, r := range results {
		if r.Value != i {
			t.Errorf("results[%d].Value = %d, want %d", i, r.Value, i)
		}
	}
}
