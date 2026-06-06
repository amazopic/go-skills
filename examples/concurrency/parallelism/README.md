# Parallelism (Scatter-Gather)

`Map[T, R]` fans a slice of inputs across a bounded pool of worker goroutines,
runs a per-item function on each, and gathers the results back **in input
order**. This is the parallel-map / scatter-gather shape for embarrassingly
parallel, CPU-bound work.

## When to use

- Each item can be processed independently (no cross-item shared state).
- The work is CPU-bound and you want to use multiple cores.
- You want bounded concurrency: a 1M-element input must not spawn 1M goroutines.

If items depend on each other, or you need streaming back-pressure between
stages, reach for the pipeline / producer-consumer patterns instead.

## API

```go
results, err := parallelism.Map(ctx, workers, inputs,
    func(ctx context.Context, item In) (Out, error) {
        return transform(ctx, item)
    })

// Convenience: workers = runtime.NumCPU()
results, err := parallelism.MapCPU(ctx, inputs, fn)
```

## Key properties

- **Order-preserving:** `results[i]` corresponds to `inputs[i]`; no post-sort.
- **Bounded:** at most `workers` invocations of `fn` run concurrently.
- **Fail-fast:** the first error is returned (wrapped with `%w`) and the shared
  context is cancelled so the remaining workers stop pulling new work.
- **Context-first:** caller cancellation unblocks in-flight workers; `Map`
  always waits for every goroutine to exit, so it never leaks.
- `workers < 1` returns `parallelism.ErrNoWorkers`.

## Run

```bash
cd examples && go test -race ./concurrency/parallelism/
```
