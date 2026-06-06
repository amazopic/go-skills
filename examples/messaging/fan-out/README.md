# Fan-Out

`FanOut[In, Out]` — distribute work from a single source channel across `N`
worker goroutines. Each item is received by exactly one worker, so the runtime
balances load automatically (work-stealing, not fixed round-robin): fast workers
pull more items than slow ones.

## When to use

- Work items are independent and order does not matter.
- A single consumer is the throughput bottleneck and the work is CPU- or
  IO-bound enough to benefit from parallelism.
- You want bounded concurrency (`N` workers) over an unbounded stream.

## API

```go
ctx := context.Background()

in := make(chan int)
go func() {
    defer close(in) // producer owns closing the source
    for _, item := range items {
        in <- item
    }
}()

// Stream results as they arrive:
out := fanout.FanOut(ctx, 8, in, func(v int) int { return v * v })
for r := range out {
    use(r)
}

// Or materialise everything (returns wrapped ctx.Err() on cancellation):
results, err := fanout.Collect(ctx, out)
```

## Key properties

- Each source item processed exactly once — no drops, no duplicates.
- The result channel is closed exactly once, after all workers exit; a single
  closer goroutine owns the `close`.
- Context-aware on both receive and send: cancelling `ctx` unblocks workers and
  guarantees no goroutine leak even if the consumer abandons the stream.
- `FanOut` panics on `n < 1`; `FanOutErr` returns `ErrInvalidWorkers` instead.

## Run

```bash
cd examples && go test -race ./messaging/fan-out/
```
