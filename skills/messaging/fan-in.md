---
name: messaging-fan-in
description: Fan-In — multiplex multiple input channels into a single output channel via a goroutine that select-reads from all inputs. Use when downstream wants a unified stream.
category: messaging
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/messaging/fan_in.md
example: examples/messaging/fan-in/
---

Fan-In Messaging Patterns
===================================
Fan-In is a messaging pattern used to create a funnel for work amongst workers (clients: source, server: destination).

We can model fan-in using Go channels. The idiomatic shape is generic over the
element type and context-aware so that the consumer can abandon the merged
stream without leaking the forwarding goroutines.

```go
// Merge fans the input channels into a single output channel. The output is
// closed once every input is drained, or once ctx is cancelled — whichever
// happens first. Forwarding goroutines never leak: each one returns on ctx
// cancellation even if the consumer has stopped reading.
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

	// Close out once all forwarders are done. Started after wg.Add.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
```

`Merge` converts a list of channels into a single channel by starting a
goroutine for each inbound channel that copies its values to the sole outbound
channel. Once all forwarders return, a final goroutine closes `out`.

Gotchas:

- **Don't leak forwarders.** A naive `for n := range c { out <- n }` blocks
  forever if the consumer stops reading. Select on `ctx.Done()` for both the
  receive and the send so every forwarder can exit.
- **Close `out` exactly once**, from a single goroutine that waits on the
  `WaitGroup` — never from inside a forwarder.
- **`wg.Add(len(cs))` before launching** the forwarders, so the closer
  goroutine can't observe a zero count prematurely.
