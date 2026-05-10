// Package pushpull implements the Push-Pull (pipeline) messaging pattern.
// Pushers fan work items into a shared buffered channel; pullers compete to
// receive each item — exactly one puller processes each item.
package pushpull

import (
	"context"
	"sync"
)

// Pipeline is a bounded work queue that distributes items to competing pullers.
// Zero value is not usable; use NewPipeline.
type Pipeline[T any] struct {
	ch chan T
}

// NewPipeline creates a Pipeline with the given buffer capacity.
// Panics if buf < 1.
func NewPipeline[T any](buf int) *Pipeline[T] {
	if buf < 1 {
		panic("pushpull: buf must be >= 1")
	}
	return &Pipeline[T]{ch: make(chan T, buf)}
}

// Push enqueues v, blocking if the buffer is full.
// Returns ctx.Err() if the context is cancelled before v can be enqueued.
// After Close, Push panics (send on closed channel) — do not call Push after Close.
func (p *Pipeline[T]) Push(ctx context.Context, v T) error {
	select {
	case p.ch <- v:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Pull returns the receive-only end of the pipeline.
// Pullers should range over this channel:
//
//	for item := range pipeline.Pull() { process(item) }
//
// The loop exits when the pipeline is closed and drained.
func (p *Pipeline[T]) Pull() <-chan T {
	return p.ch
}

// Close signals that no more items will be pushed. Pullers draining Pull()
// via range will exit after the last item. Must be called exactly once and
// only after all pushers have finished.
func (p *Pipeline[T]) Close() {
	close(p.ch)
}

// Len returns the number of items currently buffered.
func (p *Pipeline[T]) Len() int {
	return len(p.ch)
}

// RunPullers starts n puller goroutines, each ranging over the pipeline and
// calling process for every item. RunPullers blocks until all pullers finish
// (i.e., the pipeline is closed and drained).
func (p *Pipeline[T]) RunPullers(n int, process func(T)) {
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			for item := range p.Pull() {
				process(item)
			}
		}()
	}
	wg.Wait()
}
