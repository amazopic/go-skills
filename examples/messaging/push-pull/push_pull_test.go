package pushpull

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
)

func TestPipeline_AllItemsDelivered(t *testing.T) {
	tests := []struct {
		name      string
		pushers   int
		pullers   int
		perPusher int
		buf       int
	}{
		{"1P1C", 1, 1, 100, 16},
		{"4P4C", 4, 4, 50, 32},
		{"1P8C", 1, 8, 200, 64},
		{"8P1C", 8, 1, 25, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPipeline[int](tt.buf)
			total := tt.pushers * tt.perPusher
			received := make(chan int, total)

			// Start pullers.
			var pullWg sync.WaitGroup
			pullWg.Add(tt.pullers)
			for range tt.pullers {
				go func() {
					defer pullWg.Done()
					for item := range p.Pull() {
						received <- item
					}
				}()
			}

			// Push all items.
			ctx := context.Background()
			var pushWg sync.WaitGroup
			pushWg.Add(tt.pushers)
			for i := range tt.pushers {
				go func(id int) {
					defer pushWg.Done()
					for j := range tt.perPusher {
						if err := p.Push(ctx, id*tt.perPusher+j); err != nil {
							t.Errorf("Push error: %v", err)
							return
						}
					}
				}(i)
			}
			pushWg.Wait()
			p.Close()
			pullWg.Wait()
			close(received)

			var got []int
			for v := range received {
				got = append(got, v)
			}
			sort.Ints(got)
			if len(got) != total {
				t.Errorf("received %d items, want %d", len(got), total)
			}
			for i, v := range got {
				if v != i {
					t.Errorf("got[%d] = %d, want %d", i, v, i)
				}
			}
		})
	}
}

func TestPipeline_EachItemDeliveredExactlyOnce(t *testing.T) {
	const (
		items   = 500
		pullers = 10
		buf     = 32
	)
	p := NewPipeline[int](buf)
	var counts [items]atomic.Int32

	var wg sync.WaitGroup
	wg.Add(pullers)
	for range pullers {
		go func() {
			defer wg.Done()
			for v := range p.Pull() {
				counts[v].Add(1)
			}
		}()
	}

	ctx := context.Background()
	for i := range items {
		if err := p.Push(ctx, i); err != nil {
			t.Fatal(err)
		}
	}
	p.Close()
	wg.Wait()

	for i := range items {
		if got := counts[i].Load(); got != 1 {
			t.Errorf("item %d received %d times, want 1", i, got)
		}
	}
}

func TestPipeline_RunPullers(t *testing.T) {
	p := NewPipeline[int](32)
	var sum atomic.Int64
	ctx := context.Background()

	go func() {
		for i := range 100 {
			_ = p.Push(ctx, i)
		}
		p.Close()
	}()

	p.RunPullers(4, func(v int) { sum.Add(int64(v)) })

	want := int64(100 * 99 / 2)
	if got := sum.Load(); got != want {
		t.Errorf("sum = %d, want %d", got, want)
	}
}

func TestPipeline_ContextCancelsPush(t *testing.T) {
	p := NewPipeline[int](1) // tiny buffer
	ctx, cancel := context.WithCancel(context.Background())

	// Fill the buffer.
	_ = p.Push(ctx, 0)

	// Cancel then try to push — should get ctx.Err().
	cancel()
	err := p.Push(ctx, 1)
	if err == nil {
		t.Error("expected error on cancelled push")
	}
}

func TestPipeline_PanicOnZeroBuf(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for buf=0")
		}
	}()
	NewPipeline[int](0)
}
