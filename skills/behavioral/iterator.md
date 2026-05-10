---
name: behavioral-iterator
description: Use when you need sequential access to a collection's elements without exposing its internal structure — especially for lazy sequences, external cursors, or tree traversal.
category: behavioral
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Behavioral/Iterator/
example: examples/behavioral/iterator/
---

# Iterator

## Intent

Provide a uniform way to traverse elements of a collection without exposing its underlying representation. Separate the traversal logic from the collection so clients can iterate without knowing whether they are walking a slice, a tree, a database cursor, or a network stream.

## Context

Go largely obsoletes the classic GoF Iterator for in-memory slices and maps — `range` does everything the pattern promises, with zero boilerplate. The pattern earns its place in three real Go scenarios:

1. **Lazy / infinite sequences**: a generator that computes the next value on demand without materializing the full collection (pagination, fibonacci stream, file lines).
2. **Tree and graph traversal**: a stateful cursor that tracks a position in a non-linear structure.
3. **External resource cursors**: `sql.Rows`, `bufio.Scanner`, gRPC streaming responses — these are iterators over I/O that cannot be held in memory all at once.

The example in `examples/behavioral/iterator/` implements the full GoF shape (`Iterator` + `Aggregate` interfaces, `BookIterator`, `BookShelf`) — useful for teaching and for cases where you need bidirectional traversal (`Next` / `Prev`) or arbitrary repositioning (`End`, `Reset`).

Go 1.23 added `iter.Seq[V]` (push-style iterators over `range`) as the stdlib-blessed form for library iterators. For Go 1.21 target, channel-based generators or the classic cursor struct are the go-to approaches.

## Implementation in Go

**GoF cursor (from the example — bidirectional, stateful):**

```go
type Iterator interface {
    Index() int
    Value() interface{}
    Has() bool
    Next()
    Prev()
    Reset()
    End()
}

type Aggregate interface {
    Iterator() Iterator
}
```

**Channel-based generator (lazy, unidirectional — idiomatic Go for pipelines):**

```go
// Integers generates an infinite sequence of integers starting at start.
func Integers(start int) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := start; ; i++ {
            ch <- i
        }
    }()
    return ch
}

// Usage: for n := range Integers(1) { if n > 5 { break } }
```

**Closure-based pull iterator (no goroutine, no allocation on every call):**

```go
// Lines returns a function that yields the next line on each call.
// Returns ("", false) when exhausted.
func Lines(s *bufio.Scanner) func() (string, bool) {
    return func() (string, bool) {
        if s.Scan() {
            return s.Text(), true
        }
        return "", false
    }
}
```

## When to use

- Traversing a tree, graph, or other non-linear structure where `range` does not apply.
- Generating lazy / infinite sequences (pagination tokens, time ticks, prime numbers).
- Wrapping an external resource (database rows, file scanner, network stream) with a uniform `Has() / Next() / Value()` API.
- Bidirectional or random-access traversal that `range` cannot express.
- Providing a stable traversal API to callers while you refactor the underlying collection type.

## When NOT to use

- Iterating over a slice, map, or channel — use `range` directly.
- When the full collection fits in memory and ordering is simple — `sort.Slice` + `range` is cleaner.
- When you only need to filter or transform — use a plain function, `slices.DeleteFunc`, or a pipeline of channels.
- Avoid the GoF cursor struct when a closure or channel generator expresses the same intent with less ceremony.

## Gotchas

- **Goroutine leak in channel generators**: if the consumer breaks early, the goroutine blocks forever sending to an unread channel. Use a context with cancellation or pass a done channel.
- **Concurrent modification**: the GoF cursor holds an index into a mutable collection. Modifying the collection during iteration (append, delete) invalidates the iterator silently.
- **`interface{}` Value()**: the example uses `interface{}` for broad compatibility. In Go 1.18+ generics, prefer `Iterator[T]` to eliminate type assertions and panics.
- **`Has()` / `Next()` ordering**: the `Has()` → read → `Next()` loop idiom differs from Go's `for scanner.Scan()` idiom. Be consistent within your API so callers don't off-by-one.
- **Resource cleanup**: external iterators (DB rows, HTTP streams) must be closed even if iteration is interrupted. Implement `io.Closer` alongside the iterator interface.

## See also

- `skills/behavioral/visitor.md` — Visitor traverses a structure too, but applies an operation to each element rather than just yielding values.
- `examples/behavioral/iterator/`
