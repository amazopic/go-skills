# Coroutines

Go's coroutine idiom: a goroutine paired with two unbuffered channels so that
exactly one side runs at a time — the Go equivalent of cooperative yield.

## Provided types

- `Coroutine[T]` — generic generator. Call `Next()` to advance; body calls
  `yield(v)` to suspend and deliver a value.
- `Fibonacci(ctx)` — example infinite generator using `Coroutine[int]`.

## Run

```bash
go test -race -v ./concurrency/coroutines/
```

## Key design choices

- Context-aware: body honours `ctx.Done()` so callers can abort without leaks.
- `Stop()` lets callers abandon early without exhausting the sequence.
- Unbuffered `yield` channel enforces strict cooperative alternation — no
  buffering means the body cannot run ahead of the caller.
