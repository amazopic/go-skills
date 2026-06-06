// Package fanin implements the Fan-In messaging pattern: multiplexing several
// input channels into a single output channel.
//
// A forwarding goroutine is started per input; each copies values to a shared
// output channel until its input is drained or the context is cancelled. A
// final goroutine closes the output once every forwarder has finished, so the
// consumer can simply range over the merged stream.
//
// All forwarders are context-aware on both receive and send, so the consumer
// can abandon the merged stream early without leaking goroutines.
package fanin

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrNoInputs is returned by [Drain] when no input channels are supplied.
var ErrNoInputs = errors.New("fanin: no input channels")

// Merge fans the input channels into a single output channel.
//
// The output channel is closed once every input is drained, or once ctx is
// cancelled — whichever happens first. Each forwarding goroutine returns on
// ctx cancellation even if the consumer has stopped reading, so no goroutine
// leaks regardless of how the consumer behaves.
//
// Calling Merge with no inputs returns an already-closed channel.
func Merge[T any](ctx context.Context, cs ...<-chan T) <-chan T {
	out := make(chan T)

	var wg sync.WaitGroup
	wg.Add(len(cs))

	// forward copies values from c to out until c is closed or ctx is done.
	forward := func(c <-chan T) {
		defer wg.Done()
		for {
			select {
			case v, ok := <-c:
				if !ok {
					return
				}
				// Guard the send too: if the consumer has stopped
				// reading, ctx cancellation must still release us.
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}

	for _, c := range cs {
		go forward(c)
	}

	// Close out once all forwarders are done. Must start after wg.Add so the
	// closer can never observe a prematurely-zero counter.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Drain merges the inputs and collects every value into a slice, preserving no
// particular ordering across inputs (values from a single input keep their
// relative order; values from different inputs interleave).
//
// Drain blocks until all inputs are drained or ctx is cancelled. If ctx is
// cancelled before every input is exhausted it returns the values collected so
// far wrapped together with ctx.Err(). With no inputs it returns ErrNoInputs.
func Drain[T any](ctx context.Context, cs ...<-chan T) ([]T, error) {
	if len(cs) == 0 {
		return nil, ErrNoInputs
	}

	out := Merge(ctx, cs...)

	var collected []T
	for v := range out {
		collected = append(collected, v)
	}

	if err := context.Cause(ctx); err != nil {
		return collected, fmt.Errorf("fanin: drain interrupted: %w", err)
	}
	return collected, nil
}
