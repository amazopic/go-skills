# Push & Pull

`Pipeline[T]` — a shared buffered channel where pushers enqueue work items and
pullers compete to dequeue them. Each item goes to exactly one puller.

## API

```go
p := pushpull.NewPipeline[int](bufSize)

// Pushers (any number of goroutines):
p.Push(ctx, item)

// Pullers (any number of goroutines, range exits when closed):
for item := range p.Pull() { process(item) }

// Or use the convenience helper:
p.RunPullers(nWorkers, func(item int) { process(item) })

// Shutdown (call once, after all pushers are done):
p.Close()
```

## Key properties

- Each item delivered to exactly one puller (work distribution, not broadcast).
- Buffer provides back-pressure: pushers block when full.
- `Push` is context-aware: returns `ctx.Err()` on cancellation.

## Run

```bash
go test -race -v ./messaging/push-pull/
```
