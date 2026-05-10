---
name: behavioral-visitor
description: Use when you need to add operations to a type hierarchy without modifying the types — especially when the hierarchy is stable but operations grow frequently.
category: behavioral
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Behavioral/Visitor/
example: examples/behavioral/visitor/
---

# Visitor

## Intent

Represent an operation to be performed on elements of an object structure. Visitor lets you define a new operation without changing the classes of the elements on which it operates. It achieves double dispatch: the operation varies by both the visitor type and the element type.

## Context

In Go, the honest answer is: **a type switch is usually the right tool**, not the GoF Visitor. Reach for the full Visitor structure only when:

1. The **element hierarchy is closed and stable** (you own it) but **operations are added frequently** by external consumers.
2. You need **true extensibility** — external packages can define new Visitor implementations without modifying the element types.
3. The double-dispatch ceremony is worth the compile-time guarantee that every operation handles every element type (the `Visitor` interface's method set enforces exhaustiveness).

When the element hierarchy is in your package and operations are few, a plain `switch v := node.(type)` is more readable, incurs no interface overhead, and avoids the `Accept/Visit` boilerplate.

The example in `examples/behavioral/visitor/` shows the full GoF shape: `Visitor` interface with per-element methods, `Place` interface with `Accept(Visitor)`, and `City` as the object structure. This is the correct shape when the set of places (SushiBar, Pizzeria, BurgerBar) is fixed but the operations on them (People visiting, Health inspector visiting, etc.) grow over time.

## Implementation in Go

**Full GoF Visitor (element hierarchy stable, operations grow):**

```go
// Visitor declares a Visit method for each concrete element type.
type Visitor interface {
    VisitSushiBar(p *SushiBar) string
    VisitPizzeria(p *Pizzeria) string
    VisitBurgerBar(p *BurgerBar) string
}

// Place is the element interface.
type Place interface {
    Accept(v Visitor) string
}

// SushiBar is a concrete element.
type SushiBar struct{}

func (s *SushiBar) Accept(v Visitor) string { return v.VisitSushiBar(s) }
func (s *SushiBar) BuySushi() string        { return "Buy sushi..." }

// People is a concrete visitor.
type People struct{}

func (p *People) VisitSushiBar(s *SushiBar) string  { return s.BuySushi() }
func (p *People) VisitPizzeria(p2 *Pizzeria) string  { return p2.BuyPizza() }
func (p *People) VisitBurgerBar(b *BurgerBar) string { return b.BuyBurger() }
```

**Type-switch Visitor (idiomatic Go for internal hierarchies):**

```go
// Node is the element interface — no Accept() needed.
type Node interface{ node() }

type Num struct{ Val int }
func (Num) node() {}

type Add struct{ Left, Right Node }
func (Add) node() {}

// Eval is the "visitor" — a plain function with a type switch.
func Eval(n Node) int {
    switch v := n.(type) {
    case Num:
        return v.Val
    case Add:
        return Eval(v.Left) + Eval(v.Right)
    default:
        panic(fmt.Sprintf("unknown node %T", v))
    }
}
```

The type-switch form is not exhaustively checked at compile time (the `default` panic is the workaround), but it is dramatically simpler and is the idiomatic Go choice for internal hierarchies.

## When to use

- AST traversal in compilers and linters where the AST node types are fixed but analysis passes grow.
- Document rendering: a fixed set of document nodes (Paragraph, Heading, Table) with multiple renderers (HTML, Markdown, PDF).
- Serialization: visit each node type with a serialization visitor without adding `Marshal()` to every element.
- When an external package must add new operations to your element types without you modifying them (true plugin extensibility).

## When NOT to use

- When the element hierarchy is internal and operations are few — use a type switch.
- When the hierarchy changes frequently (new element types added often) — every new type forces all existing Visitors to add a method; use a map of handlers or the type-switch approach instead.
- When elements are simple value types — just write functions.
- When you control both elements and operations — you can add methods to the elements directly.

## Gotchas

- **Exhaustiveness gap in type switches**: Go does not warn if you miss a `case` in a type switch. The compiler catches missing methods in a GoF Visitor interface. If exhaustiveness matters, the GoF form provides it; otherwise add a `default: panic(...)` in the type switch.
- **Type assertion panic**: `Accept` passes `s *SushiBar` to `VisitSushiBar`; the concrete type must match exactly. Accidental nil concrete pointers inside a non-nil `Place` interface will panic.
- **Adding element types**: every new concrete element type requires adding a method to the `Visitor` interface, which breaks all existing visitor implementations. This is the fundamental tension of the pattern.
- **Cyclic visits**: a `City.Accept` iterates over `places` calling `Accept` on each. If a `Place` implementation's `Accept` itself calls `city.Accept`, you recurse infinitely. Visitor trees must be acyclic.
- **State in visitors**: a stateful visitor (e.g., accumulating a result) is not goroutine-safe if called concurrently on different elements. Either make the visitor immutable or use per-goroutine visitor instances.

## See also

- `skills/behavioral/iterator.md` — Iterator traverses a structure; Visitor applies an operation to each element during traversal.
- `skills/behavioral/strategy.md` — Strategy swaps the algorithm for a single operation; Visitor dispatches multiple operations across a type hierarchy.
- `examples/behavioral/visitor/`
