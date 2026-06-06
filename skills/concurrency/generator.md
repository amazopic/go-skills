---
name: concurrency-generator
description: Generator — a goroutine that emits values on a channel, returned to the caller as a `<-chan T`. Use to lazy-stream values without exposing channel internals.
category: concurrency
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/concurrency/generator.md
example: examples/concurrency/generator/
---

# Generator Pattern

[Generators](https://en.wikipedia.org/wiki/Generator_(computer_programming)) yields a sequence of values one at a time.

## Implementation

Return a **receive-only** `<-chan T` so callers can only range over the stream —
they can never send or close it. Make the producer **context-aware**: a consumer
that abandons the stream early must be able to cancel, or the producer blocks
forever on its send and leaks. `Count` below emits the half-open range
`[start, end)`.

```go
// Count emits start, start+1, ..., end-1 on the returned stream, then closes it.
// The producer goroutine exits as soon as ctx is cancelled, even mid-range.
func Count(ctx context.Context, start, end int) <-chan int {
    ch := make(chan int)

    go func() {
        defer close(ch) // exactly once, so `range` always terminates
        for i := start; i < end; i++ {
            select {
            case ch <- i: // blocks until the consumer receives
            case <-ctx.Done(): // consumer gave up — stop and release the goroutine
                return
            }
        }
    }()

    return ch
}
```

## Usage

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // if we stop ranging early, this lets the producer exit

fmt.Println("No bottles of beer on the wall")

for i := range Count(ctx, 1, 100) {
    fmt.Println("Pass it around, put one up,", i, "bottles of beer on the wall")
    // ... 1 through 99 ...
}

fmt.Println(100, "bottles of beer on the wall")
```
