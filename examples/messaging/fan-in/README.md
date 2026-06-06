# Fan-In

`Merge[T]` multiplexes several input channels into a single output channel: one
forwarding goroutine per input copies values to a shared output, and a final
goroutine closes the output once every input is drained. The consumer just
ranges over the merged stream.

## When to use

- A downstream stage wants one unified stream from many producers (workers,
  shards, sources) without caring which producer a value came from.
- The dual of fan-out: many-to-one instead of one-to-many.

## API

```go
out := fanin.Merge(ctx, c1, c2, c3) // <-chan T
for v := range out {                // exits when all inputs drained or ctx done
    process(v)
}

// Or collect everything into a slice:
vals, err := fanin.Drain(ctx, c1, c2, c3)
```

## Key properties

- **Generic** over the element type.
- **No goroutine leaks**: forwarders select on `ctx.Done()` for both receive
  and send, so cancelling the context releases them even if the consumer stops
  reading.
- **Output closed exactly once**, from a single closer goroutine that waits on
  a `sync.WaitGroup` — never from inside a forwarder.
- Cross-input ordering is not preserved (values interleave); per-input order is.

## Run

```bash
cd examples && go test -race ./messaging/fan-in/
```
