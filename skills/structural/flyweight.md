---
name: structural-flyweight
description: Use when you create large numbers of fine-grained objects that share most of their state. Cache the shared (intrinsic) state; pass the varying (extrinsic) state per call.
category: structural
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Structural/Flyweight/
example: examples/structural/flyweight/
---

# Flyweight

## Intent

Flyweight reduces memory consumption by sharing common state across many fine-grained objects. The pattern splits object state into two parts: **intrinsic** (immutable, shared, cached) and **extrinsic** (context-specific, passed in at call time). A factory ensures that each unique intrinsic state is instantiated only once.

## Context

Flyweight is a memory-optimisation pattern. Apply it when profiling reveals that a large number of nearly-identical objects is consuming significant heap. The trigger is usually: N objects, each holding the same M bytes of immutable data — multiplied blindly instead of shared.

Classic examples: glyphs in a text editor (character bitmap shared, position extrinsic), particles in a game engine (sprite image shared, position/velocity extrinsic), icons in a file manager (image data shared, filename/path extrinsic).

In Go the pattern often converges on two idioms:
1. **A keyed factory map** (as in the example) — best when object identity matters and you query by a string key.
2. **`sync.Pool`** — best when objects are short-lived (request-scoped buffers), creation is cheap, and you don't need per-key identity.

Know which one fits before reaching for the pattern.

## Implementation in Go

The example in `examples/structural/flyweight/` models image rendering. The filename is intrinsic state (loaded once, shared). Width, height, and opacity are extrinsic state (passed to `Draw` each call).

```go
// Flyweighter defines the shared object's interface.
// Extrinsic state (width, height, opacity) is passed per call.
type Flyweighter interface {
    Draw(width, height int, opacity float64) string
}

// ConcreteFlyweight holds only intrinsic state.
type ConcreteFlyweight struct {
    filename string // intrinsic — immutable, shared
}

func (f *ConcreteFlyweight) Draw(width, height int, opacity float64) string {
    return fmt.Sprintf("draw image: %s, width: %d, height: %d, opacity: %.2f",
        f.filename, width, height, opacity)
}

// FlyweightFactory caches instances by intrinsic key.
type FlyweightFactory struct {
    pool map[string]Flyweighter
}

func (f *FlyweightFactory) GetFlyweight(filename string) Flyweighter {
    if f.pool == nil {
        f.pool = make(map[string]Flyweighter)
    }
    if _, ok := f.pool[filename]; !ok {
        f.pool[filename] = &ConcreteFlyweight{filename: filename}
    }
    return f.pool[filename]
}
```

Two calls with `"cat.jpg"` return the **same** `*ConcreteFlyweight`; only `Draw` is called again with different extrinsic state.

### sync.Pool alternative

When flyweights are short-lived (e.g., request buffers), `sync.Pool` is idiomatic and GC-friendly:

```go
var bufPool = sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}

func handle(data []byte) {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufPool.Put(buf)
    buf.Write(data)
    // use buf ...
}
```

`sync.Pool` does not guarantee identity across GC cycles — use the keyed factory map when you need stable, long-lived shared objects.

## When to use

- Thousands or millions of fine-grained objects are eating measurable memory.
- Objects share a large chunk of immutable state (intrinsic) and differ only in small per-call context (extrinsic).
- Object creation is expensive (IO, decoding) and creation frequency is high enough that caching pays off.
- Short-lived, poolable allocations with `sync.Pool` (buffers, decoders, serialisers).

## When NOT to use

- The object count is small — premature optimisation adds complexity with no measurable benefit.
- Intrinsic state is mutable — sharing mutable state is a data race waiting to happen.
- The extrinsic state is large enough that passing it per call costs more than the allocation saved.
- Object identity semantics matter to callers (pointer equality) — callers may be surprised to receive the same pointer from two "different" requests.

## Gotchas

- **Concurrent factory access.** The keyed map is not safe for concurrent use. Protect with a `sync.RWMutex` or use `sync.Map` if reads vastly outnumber writes.

```go
type FlyweightFactory struct {
    mu   sync.RWMutex
    pool map[string]Flyweighter
}

func (f *FlyweightFactory) GetFlyweight(filename string) Flyweighter {
    f.mu.RLock()
    fw, ok := f.pool[filename]
    f.mu.RUnlock()
    if ok {
        return fw
    }
    f.mu.Lock()
    defer f.mu.Unlock()
    if fw, ok = f.pool[filename]; !ok { // double-check after acquiring write lock
        fw = &ConcreteFlyweight{filename: filename}
        f.pool[filename] = fw
    }
    return fw
}
```

- **Intrinsic state mutation.** If a flyweight is shared and its fields are mutated, all sharers see the change — and the race detector will fire. Make intrinsic state read-only by convention or use unexported fields with no setters.
- **Factory memory growth.** The pool map grows without bound if the key space is unbounded. Add an eviction policy (LRU, TTL) for long-running services with many unique keys.
- **sync.Pool eviction.** `sync.Pool` objects can be reclaimed at any GC cycle. Never store state in a pooled object that must survive beyond a single operation.
- **Measuring first.** Flyweight is a response to a measured memory problem. Applying it speculatively produces complex code for no benefit. Profile before you pool.

## See also

- skills/structural/proxy.md — controls access; Flyweight controls allocation
- skills/structural/facade.md — simplifies an API; Flyweight optimises memory layout
- examples/structural/flyweight/
