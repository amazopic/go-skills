---
name: behavioral-registry
description: Use when multiple components need to locate named singletons or plugins by string key — codecs, drivers, handlers — without creating direct import dependencies between them.
category: behavioral
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/behavioral/registry/
---

# Registry

## Intent

Maintain a goroutine-safe, typed lookup table that maps string names to shared values — singletons, plugin implementations, codecs, drivers, or factories. Consumers retrieve values by name without importing the concrete packages that provide them; producers register at init time without knowing who will consume.

## Context

The Registry pattern arises whenever a system is extended by plugins or named implementations that must be discovered at runtime. Examples: `database/sql` driver registration (`sql.Register`), `image` format decoders (`image.RegisterFormat`), HTTP handler routing, codec maps in serialization libraries.

**Registry vs. Dependency Injection**: a Registry is a global lookup table; DI is an explicit constructor argument. They solve opposite problems:

- **Registry** wins when the number of registered types is open (plugins written by third parties) and the lookup key is determined at runtime (e.g., a config file says `"driver": "postgres"`).
- **DI wins in almost every other case.** Hard-coding a global registry into business logic makes it impossible to swap implementations in tests without mutating global state. If you know the dependency at construction time, inject it — don't register it.

The Go standard library uses the Registry pattern sparingly and deliberately. Before adding one to your codebase, ask: "can I pass this as a constructor argument instead?"

## Implementation in Go

```go
// Registry[T] is a goroutine-safe map from string names to values of type T.
type Registry[T any] struct {
    mu    sync.RWMutex
    items map[string]T
}

// Register adds name → value. Returns an error on duplicate.
func (r *Registry[T]) Register(name string, value T) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if r.items == nil {
        r.items = make(map[string]T)
    }
    if _, exists := r.items[name]; exists {
        return fmt.Errorf("registry: %q is already registered", name)
    }
    r.items[name] = value
    return nil
}

// Get retrieves a value by name. Returns (zero, false) when not found.
func (r *Registry[T]) Get(name string) (T, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    v, ok := r.items[name]
    return v, ok
}

// MustRegister panics on duplicate — for package-level init registration.
func (r *Registry[T]) MustRegister(name string, value T) {
    if err := r.Register(name, value); err != nil {
        panic(err)
    }
}
```

Package-level registration pattern (same as `database/sql`):

```go
// codecs/registry.go
var Codecs Registry[Codec]

// codecs/json/json.go
func init() {
    codecs.Codecs.MustRegister("json", JSONCodec{})
}

// Usage: import _ "codecs/json" to trigger init registration.
```

## When to use

- Plugin architectures where third parties provide implementations at compile time (blank-import `init`).
- Named driver / codec lookup determined by configuration: `"driver": "postgres"` → look up `pg.Driver`.
- Protocol multiplexers: receive a message type tag, look up the handler by tag name.
- Feature flags backed by registered handlers.
- Any place the standard library uses `init`-based registration (`database/sql`, `image`, `net/http`).

## When NOT to use

- When the dependency is known at construction time — inject it directly; no registry needed.
- When there is only one implementation — a registry of one is pointless ceremony.
- In unit-tested business logic — global registry state leaks between tests; use interface injection instead.
- When the registry would grow without bound at runtime — a registry is for a fixed set of known implementations, not a cache.
- Avoid a mutable registry after the init phase — write-after-startup races and makes the system non-deterministic.

## Gotchas

- **Global mutable state**: a package-level `var Registry[T]` is mutated during `init`. Tests that run in parallel (or in the same binary with different expected contents) will race. Make registries test-injectable by passing them as constructor arguments, or reset them between tests with a `Clear()` method gated behind a build tag.
- **Init order**: Go's init order within a package is deterministic but across packages depends on import order. If two `init` functions register the same name, the second one panics — but which `init` runs first is not guaranteed by the language spec. Never rely on registry insertion order.
- **`MustRegister` panic in init**: a panic in `init` produces a cryptic stack trace. Log a human-readable error before panicking, or use `Register` and fail the startup explicitly in `main`.
- **Generics require Go 1.18+**: the `Registry[T]` implementation uses generics. For Go < 1.18, use `map[string]interface{}` with type assertions, or a concrete type per registry (`map[string]Codec`).
- **Nil map**: the zero value of `Registry[T]` is safe (`items` is lazily initialized in `Register`), but `Get` on a zero-value registry returns `(zero, false)` — correct, not a panic.
- **Lock granularity**: `sync.RWMutex` allows concurrent reads. If `T` itself is a mutable object shared across goroutines, the registry only protects the map — not the values inside it. Document thread-safety expectations for `T`.

## See also

- `skills/behavioral/strategy.md` — Strategy selects an algorithm; a Registry often stores the available strategies.
- `skills/creational/singleton.md` — Singleton is one instance; a Registry is many named instances with the same guarantees.
- `examples/behavioral/registry/`
