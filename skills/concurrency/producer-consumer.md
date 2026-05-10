---
name: concurrency-producer-consumer
description: Use when one or more goroutines generate work items and one or more goroutines process them — decouple production from consumption with backpressure and graceful shutdown.
category: concurrency
go-version-min: "1.22"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/concurrency/producer-consumer/
---

# Producer Consumer

## Intent

Producer Consumer decouples work generation from work processing. Producers push items into a bounded buffer; consumers pull and process them independently. The buffer absorbs bursts, and its bounded capacity creates natural backpressure: when consumers are slow, producers block rather than building up unbounded memory. Graceful shutdown propagates from producers (close the channel) to consumers (range exits when channel is empty and closed).

## Context

Wherever a program reads from a slow source (network, disk, scraper) and processes items with CPU-intensive logic, or vice versa, the producer-consumer split allows both sides to run at their natural pace without one starving the other. Common instances: HTTP request ingestion → background job processing; file crawler → indexer; event stream → analytics pipeline.

In Go, a **buffered channel** is the queue. It provides thread safety, backpressure, and close-signalling for free. The idiom is: producers send to `jobs`; consumers `range jobs`; producers close `jobs` when done. A `context.Context` provides cancellation for both sides.

## Implementation in Go

```go
func Run(ctx context.Context, nProducers, nConsumers, bufSize int) {
    jobs := make(chan int, bufSize) // bounded buffer = backpressure

    // Producers: close jobs when all are done.
    var prodWg sync.WaitGroup
    prodWg.Add(nProducers)
    for i := range nProducers {
        go func(id int) {
            defer prodWg.Done()
            for j := range 100 {
                select {
                case jobs <- id*100 + j:
                case <-ctx.Done():
                    return
                }
            }
        }(i)
    }
    go func() { prodWg.Wait(); close(jobs) }()

    // Consumers: range exits when jobs is closed and drained.
    var consWg sync.WaitGroup
    consWg.Add(nConsumers)
    for range nConsumers {
        go func() {
            defer consWg.Done()
            for job := range jobs {
                process(job)
            }
        }()
    }
    consWg.Wait()
}
```

Key points: `close(jobs)` is called exactly once, from a dedicated goroutine that waits for all producers. Consumers use `range` which unblocks on close. The `select` on `ctx.Done()` lets producers respect cancellation without closing the channel themselves.

## When to use

- Decoupling a slow source from a CPU-bound processor (or vice versa).
- When you need explicit backpressure: buffer full → producers stall → upstream pressure.
- When the number of producers and consumers should scale independently.
- Any "work queue" scenario: job processing, event pipelines, ETL.

## When NOT to use

- When items must be processed in strict FIFO order by a single worker — a simple goroutine + channel is enough.
- When you need priority queuing — a channel is FIFO; use `container/heap` behind a mutex.
- When production and consumption are tightly coupled and always run at the same rate — the indirection adds latency without benefit.
- When items carry errors that must propagate back to producers — consider `errgroup` or a result channel alongside the jobs channel.

## Gotchas

- **Closing from multiple producers**: only one goroutine may close a channel. Use a `sync.WaitGroup` + sentinel goroutine. Never close from inside a producer.
- **Panic on send to closed channel**: if a consumer closes the channel (it shouldn't) or a producer escapes the WaitGroup, the send panics. Keep close ownership clear.
- **Context not threaded through**: if producers do not select on `ctx.Done()`, a cancelled context does not stop them — they keep filling the buffer.
- **Buffer size tuning**: too small → producers starve; too large → memory pressure and slow cancellation. Profile with `go tool pprof` and tune to your burst characteristics.
- **Consumer errors lost**: `range jobs` discards the job on error. Always collect errors via an `errs chan error` or `errgroup`.

## See also

- skills/concurrency/broadcast.md
- skills/concurrency/n-barrier.md
- skills/messaging/push-pull.md
- examples/concurrency/producer-consumer/
