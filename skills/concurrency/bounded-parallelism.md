---
name: concurrency-bounded-parallelism
description: Bounded parallelism — process N items with at most K workers concurrently using a worker-pool goroutine + channels. Use to cap concurrency on an unbounded input list.
category: concurrency
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/concurrency/bounded_parallelism.md
example: examples/concurrency/bounded-parallelism/
---

# Bounded Parallelism Pattern

[Bounded parallelism](https://blog.golang.org/pipelines#TOC_9.) is similar to [parallelism](parallelism.md), but caps how many goroutines run at once. The naive "one goroutine per item" approach scales goroutine count with the input and can exhaust file descriptors, sockets, or memory on large inputs. Bounded parallelism launches exactly K workers that compete to pull indexed jobs from a shared channel, so at most K invocations run concurrently regardless of input size.

# Idiomatic Go

Launch K workers over a shared `jobs` channel of indices; each worker writes its result into a pre-sized output slice at its own index (disjoint indices ⇒ no data race), so output order matches input and no result channel or sort is needed. Be context-first, fail-fast on the first error (wrapped with `%w`), and join every goroutine before returning so nothing leaks.

```go
// Map applies fn to every element of in using at most workers goroutines,
// returning results in input order. The first fn error cancels the rest.
func Map[T, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]R, error) {
	if workers < 1 {
		return nil, errors.New("workers must be >= 1")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(in) == 0 {
		return []R{}, nil
	}
	if workers > len(in) {
		workers = len(in) // more workers than items is wasteful
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	out := make([]R, len(in))
	jobs := make(chan int)

	var errOnce sync.Once
	var firstErr error
	fail := func(err error) { errOnce.Do(func() { firstErr = err; cancel() }) }

	// Dispatcher: feed indices, then close so workers' range loops end.
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

	// K workers: pull indices, run fn, write results in place.
	var workerWg sync.WaitGroup
	workerWg.Add(workers)
	for range workers {
		go func() {
			defer workerWg.Done()
			for i := range jobs {
				if ctx.Err() != nil {
					return
				}
				r, err := fn(ctx, in[i])
				if err != nil {
					fail(fmt.Errorf("item %d: %w", i, err))
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
	if err := ctx.Err(); err != nil { // parent cancelled mid-flight
		return nil, err
	}
	return out, nil
}
```

## Gotchas

- Don't spawn one goroutine per item — that's [parallelism](parallelism.md), not *bounded* parallelism. The whole point is the fixed worker count.
- Workers writing `out[i]` at disjoint indices is race-free; sharing a map or appending to one slice is not — guard those with a mutex or a result channel.
- Always `close(jobs)` from a single owner (the dispatcher) so worker `range` loops terminate; closing from a worker risks a double close.
- Join both the dispatcher and the workers (`WaitGroup`) before returning, or a cancelled dispatch leaks a goroutine blocked on `jobs <- i`.
- Cap `workers` at `len(in)`; otherwise extra workers spin up only to exit immediately.

# Implementation and Example

A runnable, race-tested version (with `ForEach` and a `Map` variant) is in [examples/concurrency/bounded-parallelism/](../../examples/concurrency/bounded-parallelism/).
