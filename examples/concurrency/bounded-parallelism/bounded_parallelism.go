// Package boundedparallelism implements bounded parallelism: process an
// arbitrarily long input slice with at most K worker goroutines running
// concurrently, regardless of the input size.
//
// The naive "one goroutine per item" approach scales goroutine count with the
// input and can exhaust file descriptors, sockets, or memory on large inputs.
// Bounded parallelism caps concurrency at a fixed K by launching exactly K
// workers that compete to pull indexed jobs from a shared channel. Results are
// written back to a pre-sized output slice indexed by position, so no result
// channel or post-sort is needed and output order matches input order.
//
// The pattern is context-first (cancellation stops dispatch and is reported),
// leak-free (every goroutine is joined via sync.WaitGroup before returning),
// and generic over the input and output element types.
package boundedparallelism

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrInvalidWorkers is returned when the requested worker count is < 1.
var ErrInvalidWorkers = errors.New("boundedparallelism: workers must be >= 1")

// Map applies fn to every element of in using at most workers concurrent
// goroutines and returns the results in input order.
//
// fn receives the per-call context and a single input element. If fn returns an
// error for any element, Map cancels the remaining work, waits for in-flight
// workers to return, and reports the first error encountered (wrapped with %w).
// On error the returned slice is nil.
//
// Map returns ErrInvalidWorkers if workers < 1, and ctx.Err() (wrapped) if ctx
// is already cancelled. For an empty input, Map returns an empty, non-nil slice
// without launching any workers.
//
// Map blocks until all work is complete or cancelled; it never leaks goroutines.
func Map[T, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]R, error) {
	if workers < 1 {
		return nil, ErrInvalidWorkers
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("boundedparallelism: %w", err)
	}
	if len(in) == 0 {
		return []R{}, nil
	}

	// Cap workers at the number of items: more workers than items is wasteful
	// (the extra goroutines would start and immediately exit).
	if workers > len(in) {
		workers = len(in)
	}

	// Derive a cancellable context so the first error (or external cancel)
	// stops the dispatcher and lets workers drain quickly.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	out := make([]R, len(in))

	// jobs carries the index of the next item to process. The dispatcher sends
	// indices; workers receive them. A buffered result/error path is not needed
	// because each worker writes directly into out[i] (disjoint indices ⇒ no
	// data race) and reports its first error once via firstErr.
	jobs := make(chan int)

	var (
		errOnce  sync.Once
		firstErr error
	)
	fail := func(err error) {
		errOnce.Do(func() {
			firstErr = err
			cancel() // stop the dispatcher and remaining workers
		})
	}

	// Dispatcher: feed indices until done or cancelled, then close jobs so
	// workers' range loops terminate.
	var dispatchWg sync.WaitGroup
	dispatchWg.Add(1)
	go func() {
		defer dispatchWg.Done()
		defer close(jobs)
		for i := range in {
			select {
			case jobs <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Workers: pull indices, run fn, write results in place.
	var workerWg sync.WaitGroup
	workerWg.Add(workers)
	for range workers {
		go func() {
			defer workerWg.Done()
			for i := range jobs {
				// Bail out promptly if another worker already failed or ctx
				// was cancelled, rather than starting more work.
				if ctx.Err() != nil {
					return
				}
				r, err := fn(ctx, in[i])
				if err != nil {
					fail(fmt.Errorf("boundedparallelism: item %d: %w", i, err))
					return
				}
				out[i] = r
			}
		}()
	}

	workerWg.Wait()
	dispatchWg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	// Our own cancel() runs only via fail (which sets firstErr). So if firstErr
	// is nil yet ctx is done, the parent cancelled us mid-flight — surface it.
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("boundedparallelism: %w", err)
	}
	return out, nil
}

// ForEach applies fn to every element of in using at most workers concurrent
// goroutines, discarding per-item results. It has the same error, cancellation,
// and leak-safety semantics as Map. Use it for side-effecting work (writes,
// network calls) where no return value is collected.
func ForEach[T any](ctx context.Context, workers int, in []T, fn func(context.Context, T) error) error {
	_, err := Map(ctx, workers, in, func(ctx context.Context, v T) (struct{}, error) {
		return struct{}{}, fn(ctx, v)
	})
	return err
}
