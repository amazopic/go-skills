---
name: structural-facade
description: Use when you want to offer a simple, high-level API over a complex or multi-step subsystem. Reduces coupling between callers and subsystem internals.
category: structural
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Structural/Facade/
example: examples/structural/facade/
---

# Facade

## Intent

Facade provides a single, simplified entry point to a subsystem of cooperating types. It does not add new behaviour — it orchestrates existing behaviour and hides the sequence, dependencies, and internal types from callers. The subsystem remains accessible directly when needed; the facade is a convenience, not a wall.

## Context

As a subsystem grows, callers accumulate knowledge of its internals: which types to instantiate, in what order to call methods, and which error conditions to handle. That coupling makes callers fragile to subsystem changes and hard to test. A Facade absorbs that knowledge into one place.

In Go, the idiom often manifests as a thin top-level package or a single exported struct that delegates to several internal (unexported) types. The `net/http` package is the canonical stdlib example: `http.ListenAndServe` is a facade over `net.Listen`, `http.Server`, and the underlying `net.Conn` machinery. You can bypass it and use those primitives directly, but 95 % of HTTP servers never need to.

## Implementation in Go

The example in `examples/structural/facade/` models the old saying "a real man must build a house, grow a tree, and raise a child." `Man` is the facade that sequences three subsystems (`House`, `Tree`, `Child`) behind a single `Todo()` call.

```go
// NewMan wires up the subsystem — callers do not construct subsystem types.
func NewMan() *Man {
    return &Man{
        house: &House{},
        tree:  &Tree{},
        child: &Child{},
    }
}

type Man struct {
    house *House
    tree  *Tree
    child *Child
}

// Todo is the single façade method; it orchestrates three subsystems.
func (m *Man) Todo() string {
    return strings.Join([]string{
        m.house.Build(),
        m.tree.Grow(),
        m.child.Born(),
    }, "\n")
}
```

Subsystem types (`House`, `Tree`, `Child`) remain exported so "advanced" callers can use them directly — consistent with the GoF guideline that a facade does not restrict access.

### Package-level facade

In idiomatic Go, the facade is often a package rather than a struct. A thin `api` or `service` package exposes a handful of functions, while the heavy lifting happens in unexported sub-packages. This is the shape of `database/sql` (facade) over the internal `driver` interfaces.

```go
// api/api.go — package-level facade
func ProcessOrder(id string) error {
    inv := inventory.Reserve(id)
    pay := payment.Charge(inv.Price())
    return shipping.Dispatch(inv, pay)
}
```

## When to use

- Onboarding developers: a facade reduces the learning surface of a complex subsystem.
- You have a stable high-level workflow and want to decouple it from subsystem evolution.
- Testing: mock the facade rather than all its internal dependencies.
- Providing a "batteries included" default while still leaving the internals available for power users.
- Reducing import graphs — the facade re-exports only what callers need; subsystem packages stay internal.

## When NOT to use

- The subsystem is simple — a facade would just be a passthrough with no value.
- Callers genuinely need fine-grained control; hiding that behind a facade forces them to work around it.
- The facade becomes a "god object" that grows unbounded — if every feature request ends up in the facade, break it into multiple smaller facades.
- You want the facade to enforce access control — use Proxy instead, which is designed for controlled delegation.

## Gotchas

- **Facade creep.** It is tempting to add methods to the facade forever. Keep it focused on the most common workflows; unusual use cases should go to the subsystem directly.
- **Hidden errors.** Orchestrating multiple subsystem calls means multiple potential error sites. Do not swallow errors silently. Propagate them clearly, or define a domain-specific error type.
- **Constructor complexity.** If `NewMan` needs to wire dozens of dependencies, consider dependency injection or a builder. A constructor that does a lot of work is itself a smell.
- **Testing the facade vs the subsystems.** Test the facade at integration level and each subsystem in isolation. Avoid re-testing subsystem logic through the facade — that duplicates tests and obscures which component failed.
- **Concurrency.** If the facade is shared across goroutines and delegates to stateful subsystems, ensure each subsystem is safe for concurrent use or serialize access in the facade.

## Facade vs other patterns

| Pattern | Key question | Typical shape |
|---|---|---|
| Facade | "How do I simplify a subsystem?" | One struct delegates to N internal types |
| Adapter | "How do I make an incompatible type fit an interface?" | One struct wraps one foreign type |
| Proxy | "How do I control access to one object?" | One struct wraps one same-interface object |
| Mediator | "How do I decouple N objects from each other?" | Central hub objects talk to; they don't know each other |

The distinction between Facade and Mediator is subtle: the facade speaks to the subsystem but the subsystem does not speak back through the facade. In a Mediator the communication is bidirectional.

## See also

- skills/structural/adapter.md — makes incompatible interfaces fit; Facade simplifies a compatible but complex interface
- skills/structural/proxy.md — controls access to one object; Facade orchestrates many objects
- examples/structural/facade/
