# Prototype — Example

Demonstrates the **Prototype** pattern in Go: new objects are created by cloning an existing instance via the `Clone()` method rather than constructing from scratch.

## Structure

| Type | Role |
|---|---|
| `Prototyper` | Interface declaring `Clone() Prototyper` and `GetName() string` |
| `ConcreteProduct` | Concrete type; `Clone()` returns an independent copy |
| `NewConcreteProduct` | Constructor that creates the initial prototype |

## Run

```bash
go test -race ./creational/prototype/
```

## Key points

- `Clone()` returns a new `*ConcreteProduct` with the same `name` field — an independent copy; modifying the clone does not affect the original.
- The example is intentionally minimal (a single `string` field). In production, `Clone()` must deep-copy any slice, map, or pointer fields to avoid sharing mutable state between original and clone.
- For deeply nested structures without a hand-written `Clone()`, an `encoding/gob` or `encoding/json` round-trip achieves the same result at the cost of some performance.

See [skills/creational/prototype.md](../../../skills/creational/prototype.md) for full pattern documentation, including deep-copy semantics and when to prefer a plain constructor instead.
