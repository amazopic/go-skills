---
name: creational-prototype
description: Prototype — clone an existing, fully configured instance rather than constructing from scratch. Use when construction is expensive or complex and a good base instance already exists.
category: creational
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Creational/Prototype/
example: examples/creational/prototype/
---

# Prototype

## Intent

Create new objects by copying an existing, fully initialized instance (the prototype). The clone starts with a known-good state and can then be modified independently. Construction cost is paid once; cloning is typically cheaper.

## Context

In most Go code the zero value is useful, constructors are cheap, and `sync.Pool` handles object reuse better than Prototype. The pattern earns its place when:

- **Construction is genuinely expensive** — e.g., a pre-warmed parser, a loaded ML model, a network connection graph — and a baseline configuration is reused with minor tweaks across many calls.
- **The object graph is complex** — the prototype captures a deep hierarchy that would require many constructor calls to recreate each time.
- **Configuration diverges from a common root** — you have a "template" object that many variants start from before customization.

In modern Go, compare against simpler alternatives first:

| Alternative | When it wins |
|---|---|
| Plain constructor | State is cheap to build; no shared baseline |
| `sync.Pool` | Pooling for allocation reduction, not configuration variation |
| Functional Options | Variation through options, not cloning |
| `encoding/json` or `encoding/gob` round-trip | Deep clone of complex graphs without hand-written `Clone()` |

## Implementation in Go

### Explicit `Clone()` with deep-copy semantics (recommended)

Make deep-copy behavior explicit. Never silently share slices or maps between the original and the clone.

```go
// DocumentTemplate is an expensive-to-construct baseline configuration.
type DocumentTemplate struct {
    Title    string
    Tags     []string            // must be deep-copied
    Metadata map[string]string   // must be deep-copied
}

// Clone returns an independent deep copy of the template.
// Modifying the clone's Tags or Metadata does not affect the original.
func (d *DocumentTemplate) Clone() *DocumentTemplate {
    tags := make([]string, len(d.Tags))
    copy(tags, d.Tags)

    meta := make(map[string]string, len(d.Metadata))
    for k, v := range d.Metadata {
        meta[k] = v
    }

    return &DocumentTemplate{
        Title:    d.Title,
        Tags:     tags,
        Metadata: meta,
    }
}
```

Usage:

```go
base := &DocumentTemplate{
    Title:    "Q4 Report",
    Tags:     []string{"finance", "internal"},
    Metadata: map[string]string{"author": "template"},
}

draft := base.Clone()
draft.Title = "Q4 Report — DRAFT"
draft.Metadata["author"] = "alice"
// base is unchanged
```

### Interface-based Prototype (matches existing example)

When multiple types share the cloning contract, define a `Cloner` interface:

```go
// Cloner is satisfied by any type that can produce an independent copy of itself.
type Cloner interface {
    Clone() Cloner
    Name() string
}

type Product struct {
    name string
}

func NewProduct(name string) Cloner { return &Product{name: name} }
func (p *Product) Name() string     { return p.name }
func (p *Product) Clone() Cloner    { return &Product{name: p.name} }
```

### Deep copy via `encoding/gob` (complex graphs, zero boilerplate)

For large, deeply nested structs where hand-rolling `Clone()` is error-prone:

```go
import (
    "bytes"
    "encoding/gob"
)

func DeepClone[T any](src T) (T, error) {
    var buf bytes.Buffer
    if err := gob.NewEncoder(&buf).Encode(src); err != nil {
        var zero T
        return zero, err
    }
    var dst T
    return dst, gob.NewDecoder(&buf).Decode(&dst)
}
```

This works for exported fields only and requires all field types to be gob-encodable. The round-trip is slower than a hand-written `Clone()` but eliminates an entire class of shallow-copy bugs.

## When to use

- Constructing an object requires expensive work (I/O, parsing, computation) and a validated baseline exists.
- Many similar objects are created that share a common starting state and then diverge.
- The exact type of the object to clone is unknown at compile time (polymorphic cloning via interface).
- An object's configuration is determined at runtime and cannot be described statically to a constructor.

## When NOT to use

- Construction is cheap — a plain `New*()` call is clearer and costs nothing.
- The "copy" is just a value type assignment — structs in Go are value types; `b := a` already copies all value fields, making Prototype unnecessary for flat structs.
- You only need object reuse, not configuration variation — use `sync.Pool`.
- Shared state between original and clone is intentional — in that case don't clone, share a pointer.

## Gotchas

- **Shallow vs deep copy**: Go struct assignment is a shallow copy. A naively written `Clone()` that returns `&ConcreteProduct{p.name}` is safe for the string field (`string` is immutable), but any slice or map field will share the underlying array/hash table between original and clone. Always `copy()` slices and range-copy maps.
- **Pointer fields**: a field of type `*Inner` copied by assignment still points to the same `Inner`. Clone the pointed-to value explicitly or use a `Clone()` method on `Inner`.
- **Interface fields**: an interface value contains a pointer to the concrete data. Cloning a struct that holds an interface field requires a type switch or a `Cloner` interface on the concrete type.
- **`sync.Mutex` and similar**: never copy a mutex by value. If the prototype contains a `sync.Mutex`, the clone must initialize a fresh one, not copy the locked/unlocked state.
- **`gob`/`json` round-trip limitations**: unexported fields are silently skipped; channel, function, and unsafe.Pointer fields are not encodable. Verify the round-trip in tests before relying on it.

## See also

- skills/creational/factory-method.md
- skills/creational/builder.md
- skills/idiom/functional-options.md
- examples/creational/prototype/
