---
name: messaging-fan-out
description: Fan-Out — distribute work from one channel across N worker goroutines. Use to scale processing throughput when items are independent.
category: messaging
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/messaging/fan_out.md
example: examples/messaging/fan-out/
---

Fan-Out Messaging Pattern
=========================
Fan-Out is a messaging pattern used for distributing work amongst workers (producer: source, consumers: destination).

We can model fan-out with Go channels: a single source channel feeding `n`
worker goroutines that compete to receive each item. Because each value is
received by exactly one worker, the runtime balances load for free — fast
workers naturally pull more items than slow ones (work-stealing semantics).

```go
// FanOut starts n workers, each draining the shared in channel and emitting
// process(item) onto the returned channel. The result channel is closed once
// every worker has finished, i.e. after in is drained (or ctx is cancelled).
func FanOut[In, Out any](ctx context.Context, n int, in <-chan In, process func(In) Out) <-chan Out {
	out := make(chan Out)

	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-in:
					if !ok {
						return // source drained
					}
					select {
					case out <- process(item):
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	// Close out exactly once, after all workers exit — no goroutine leak.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
```

The `FanOut` function spreads work from one source channel across `n` worker
goroutines. Each worker ranges the shared input until it is drained or the
context is cancelled, so there is no fixed round-robin assignment — workers
self-balance. A single closer goroutine waits on the `WaitGroup` and closes the
output channel, guaranteeing every send completes before the channel closes and
that no goroutine is left blocked.

Gotchas:

- **Close the output exactly once.** Let a dedicated goroutine `wg.Wait()` then
  `close(out)`; never close from inside a worker.
- **Honor the context on both send and receive.** A worker can block sending to
  `out` if the consumer goes away; the inner `select` on `ctx.Done()` prevents a
  leak.
- **The caller owns `in`.** Workers stop when `in` is closed; the producer must
  close it. Closing the result channel does not close the source.
