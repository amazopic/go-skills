# Generator

A generator is a goroutine that lazily produces a sequence of values and emits
them on a channel, handed back to the caller as a receive-only `<-chan T`. The
caller ranges over the stream; it can never send or close the channel — those
internals stay with the producer.

## When to use

- Stream values one at a time instead of materializing a whole slice.
- Decouple production from consumption (the producer blocks until the consumer
  is ready — natural back-pressure).
- Build pipelines: `Count -> Map -> Take`, each stage a generator.

## API

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // releases producers if you stop early

for v := range generator.Count(ctx, 0, 10) {
    fmt.Println(v) // 0..9
}

// Compose stages:
nums    := generator.Count(ctx, 1, 1_000)
doubled := generator.Map(ctx, nums, func(n int) int { return n * 2 })
first5  := generator.Take(ctx, doubled, 5) // 2,4,6,8,10

// Custom producer via yield:
gen := generator.Generate(ctx, func(yield func(string) bool) {
    for _, s := range []string{"a", "b", "c"} {
        if !yield(s) { return } // yield reports false on cancellation
    }
})
```

## Key properties

- **Receive-only return**: callers get `<-chan T`; they cannot misuse send/close.
- **Context-first, leak-free**: `yield` (and every stage) selects on `ctx.Done()`,
  so a blocked-on-send producer unblocks the instant you cancel. Always cancel
  the context if you stop ranging before the stream is exhausted.
- **Always terminates**: each generator closes its own channel exactly once when
  the sequence ends or the context is cancelled, so `range` is guaranteed to end.
- **Generic**: `Generate`, `FromSlice`, `Map`, and `Take` work for any element
  type; `Map` can change the element type.

## Run

```bash
cd examples && go test -race ./concurrency/generator/
```
