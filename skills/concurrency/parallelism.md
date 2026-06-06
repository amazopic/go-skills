---
name: concurrency-parallelism
description: Parallel scatter-gather — fan a workload across goroutines and collect results via WaitGroup or channel. Use for embarrassingly parallel CPU-bound work.
category: concurrency
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/concurrency/parallelism.md
example: examples/concurrency/parallelism/
---

# Parallelism Pattern

[Parallelism](https://blog.golang.org/pipelines#TOC_8.) allows multiple "jobs" or tasks to be run concurrently and asynchronously.

# Implementation and Example

The idiomatic Go shape is a **bounded parallel map**: fan a slice of inputs
across at most `workers` goroutines, run a per-item function on each, and
gather results back in input order. Each worker writes a distinct result index,
so no lock is needed. The shared context is cancelled on the first error so the
remaining workers stop pulling new work — no goroutine leaks.

```go
// Map applies fn to every element of in concurrently, using at most workers
// goroutines, and returns the results in the same order as in. The first
// non-nil error cancels the shared context and is returned (wrapped with %w).
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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]R, len(in))
	type job struct {
		idx int
		val T
	}
	jobs := make(chan job)

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
				if err := ctx.Err(); err != nil {
					return
				}
				r, err := fn(ctx, j.val)
				if err != nil {
					fail(fmt.Errorf("parallelism: item %d: %w", j.idx, err))
					return
				}
				results[j.idx] = r // distinct index per worker: race-free
			}
		}()
	}

	// Dispatch; stop early if the context is cancelled.
	go func() {
		defer close(jobs)
		for i, v := range in {
			select {
			case jobs <- job{idx: i, val: v}:
			case <-ctx.Done():
				return
			}
		}
	}()
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("parallelism: cancelled: %w", err)
	}
	return results, nil
}
```

A complete, runnable example (with `MapCPU`, sentinel errors, and race-safe
tests) can be found in [parallelism.go](parallelism.go).
