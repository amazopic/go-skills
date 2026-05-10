# Broadcast

Two broadcast mechanisms in idiomatic Go:

- **Signal** — one-shot broadcast via `close(ch)`. Zero allocation per waiter.
  Used for shutdown, "ready", config-loaded signals.
- **Broadcaster[T]** — repeatable typed broadcast to a dynamic subscriber set.
  Fan-out with drop-on-full policy; call `Unsubscribe` when a receiver exits.

## Run

```bash
go test -race -v ./concurrency/broadcast/
```

## Key points

- Closing a channel wakes **all** receivers simultaneously — the canonical Go
  one-shot broadcast. Protect with `sync.Once` to prevent double-close panics.
- For repeatable events, maintain a slice of buffered subscriber channels.
  Choose between drop (non-blocking send) or block based on delivery guarantees.
- `context.WithCancel` already provides broadcast cancellation — reach for it
  before writing your own `Signal`.
