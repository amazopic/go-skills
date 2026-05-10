# Reactor

A single-goroutine event loop that dispatches events to registered handler
functions. Because all handlers execute sequentially inside the loop goroutine,
they may share state without mutexes — but must not block.

## API

```go
r := reactor.New(bufSize)
r.Register("eventKind", func(e reactor.Event) { /* handle */ })

ctx, cancel := context.WithCancel(context.Background())
go r.Run(ctx)   // starts the loop

r.Send(reactor.Event{Kind: "eventKind", Payload: data})
cancel()        // stops the loop; queued events are drained first
```

## Run

```bash
go test -race -v ./concurrency/reactor/
```

## Trade-offs

- Serial dispatch = no lock needed inside handlers.
- Blocking handlers stall all other events — dispatch long work to goroutines.
- Buffer size controls back-pressure; `TrySend` gives non-blocking sends.
