// Package parallelism implements the scatter-gather (parallel map) pattern:
// fan a slice of inputs across a bounded pool of goroutines, run a per-item
// function on each, then gather the results back in input order.
//
// It is intended for embarrassingly parallel, CPU-bound work where each item
// can be processed independently. The pool size bounds concurrency so a large
// input set cannot spawn an unbounded number of goroutines.
//
// All blocking operations are context-aware: cancellation stops the dispatch
// of further work and unblocks in-flight workers without leaking goroutines.
package parallelism

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
)

// ErrNoWorkers is returned by Map when the requested worker count is less than 1.
var ErrNoWorkers = errors.New("parallelism: workers must be >= 1")

// Map applies fn to every element of in concurrently, using at most workers
// goroutines, and returns the results in the same order as in.
//
// If workers <= 0, Map returns ErrNoWorkers. As a convenience, pass
// runtime.NumCPU() (see MapCPU) for CPU-bound workloads.
//
// fn receives the per-call context so it can abort long-running work on
// cancellation. The first non-nil error returned by any fn invocation is
// returned to the caller (wrapped with %w), and the shared context is
// cancelled so the remaining workers stop pulling new items promptly. On
// error the returned slice is nil.
//
// If ctx is cancelled before Map completes, Map returns ctx.Err() (wrapped)
// and nil results.
//
// Map blocks until all workers have exited; it never leaks goroutines.
func Map[T, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]R, error) {
	if workers < 1 {
		return nil, ErrNoWorkers
	}
	if len(in) == 0 {
		return []R{}, nil
	}
	if workers > len(in) {
		workers = len(in)
	}

	// Derive a cancellable context so the first error (or external
	// cancellation) tears down every worker.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]R, len(in))

	// indexed carries each input together with its position so results can be
	// written back in order without a post-sort.
	type job struct {
		idx int
		val T
	}
	jobs := make(chan job)

	// firstErr captures the earliest error; errOnce guards the single write.
	var (
		errOnce  sync.Once
		firstErr error
	)
	fail := func(err error) {
		errOnce.Do(func() {
			firstErr = err
			cancel() // stop dispatch and unblock the other workers
		})
	}

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for j := range jobs {
				// Bail out fast if work has been cancelled.
				if err := ctx.Err(); err != nil {
					return
				}
				r, err := fn(ctx, j.val)
				if err != nil {
					fail(fmt.Errorf("parallelism: item %d: %w", j.idx, err))
					return
				}
				// Each worker writes a distinct index, so no lock is needed
				// for results; this is data-race-free by construction.
				results[j.idx] = r
			}
		}()
	}

	// Dispatch jobs. Stop early if the context is cancelled (either by an
	// errored worker via fail, or by the caller).
	dispatch := func() {
		defer close(jobs)
		for i, v := range in {
			select {
			case jobs <- job{idx: i, val: v}:
			case <-ctx.Done():
				return
			}
		}
	}
	dispatch()
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("parallelism: cancelled: %w", err)
	}
	return results, nil
}

// MapCPU is Map with the worker count set to runtime.NumCPU(), the common
// default for CPU-bound scatter-gather work.
func MapCPU[T, R any](ctx context.Context, in []T, fn func(context.Context, T) (R, error)) ([]R, error) {
	return Map(ctx, runtime.NumCPU(), in, fn)
}
