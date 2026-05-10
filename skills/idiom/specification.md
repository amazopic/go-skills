---
name: idiom-specification
description: Specification — encapsulate a named boolean predicate over a domain entity for reuse and And/Or/Not composition. Use when filter logic is complex, must be named, tested alone, or combined at runtime from independently reusable rules.
category: idiom
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Unsorted/Specification/
example: examples/idiom/specification/
---

# Specification

## Intent

Encapsulate a single business rule as a named, composable predicate object. Rules written as `Specification` values can be combined with `And`, `Or`, and `Not` combinators without scattering conditional logic across the call sites that need them.

## Context

The pattern originates from Domain-Driven Design (Evans, Fowler). It shines when:

- the same predicate appears in multiple places (query, validation, UI display)
- rules must be composed at runtime based on user configuration
- individual rules must be unit-tested in isolation with a meaningful name

In straightforward Go code, a plain `func(T) bool` or a method on the domain type is usually sufficient. Reserve the heavier named-type Specification only when you genuinely need serialization, reflection, rule-engine integration, or a large catalog of named, independently testable rules.

## Implementation in Go

### Idiomatic form: functional combinators

Prefer this for most Go code. No interface overhead, no struct embedding, just closures.

```go
// Predicate is a typed boolean function over any domain entity.
type Predicate[T any] func(T) bool

// And returns a predicate satisfied only when all given predicates are satisfied.
func And[T any](specs ...Predicate[T]) Predicate[T] {
    return func(v T) bool {
        for _, s := range specs {
            if !s(v) {
                return false
            }
        }
        return true
    }
}

// Or returns a predicate satisfied when at least one given predicate is satisfied.
func Or[T any](specs ...Predicate[T]) Predicate[T] {
    return func(v T) bool {
        for _, s := range specs {
            if s(v) {
                return true
            }
        }
        return false
    }
}

// Not negates a predicate.
func Not[T any](spec Predicate[T]) Predicate[T] {
    return func(v T) bool { return !spec(v) }
}
```

Usage:

```go
overDue     := Predicate[Invoice](func(inv Invoice) bool { return inv.DaysPastDue >= 30 })
noticeSent  := Predicate[Invoice](func(inv Invoice) bool { return inv.NoticeCount >= 3 })
inCollection := Predicate[Invoice](func(inv Invoice) bool { return inv.SentToCollection })

sendToCollection := And(overDue, noticeSent, Not(inCollection))

if sendToCollection(invoice) {
    collector.Send(invoice)
}
```

### Named-type form (GoF / DDD style)

Use this when rules must be registered by name, persisted, or composed via a rule engine.

```go
type Specification[T any] interface {
    IsSatisfiedBy(T) bool
}

type andSpec[T any] struct{ left, right Specification[T] }

func (s andSpec[T]) IsSatisfiedBy(v T) bool {
    return s.left.IsSatisfiedBy(v) && s.right.IsSatisfiedBy(v)
}

func And[T any](left, right Specification[T]) Specification[T] {
    return andSpec[T]{left, right}
}
```

The example in `examples/idiom/specification/` uses the named-type form without generics (Go 1.17-era style). Its `BaseSpecification` embedding pattern wires `And`/`Or`/`Not` factory methods onto every concrete specification so they can be chained: `overDue.And(noticeSent).And(inCollection.Not())`.

## When to use

- Business rules need individual names, documentation, and unit tests.
- Rules are assembled dynamically (e.g., from user configuration or a rule table).
- The same predicate is used across queries, validation, and display layers and duplication is becoming a maintenance problem.
- You need serialization or storage of rule combinations (persistence of a composed specification as data).

## When NOT to use

- A plain `if` statement or a `func(T) bool` closure is readable and not repeated — do not add the indirection.
- The predicate is used in exactly one place; name it with a well-chosen function name instead.
- Performance is critical in a hot path — a closure chain adds allocations and indirect calls; a hand-written boolean expression is faster.
- The domain is simple enough that the pattern adds ceremony without benefit.

## Gotchas

- **`InCollectionSpecification` fix**: the upstream example had inverted semantics — `IsSatisfiedBy` returned `!elm.IsSent`, treating `IsSent == false` as "already in collection". The correct semantics: `IsSatisfiedBy` returns `true` when the invoice *is* in the collection (`IsSent == true`). Compose with `Not()` to express "not yet sent to collection". The test error message was also inverted and has been corrected.
- **Embedding pitfalls**: the example's `BaseSpecification` approach embeds an interface in a struct to inherit `And`/`Or`/`Not`. The embedded field must be wired back to the concrete type via `Relate()` so that `And(s, other).IsSatisfiedBy(x)` dispatches to `s`'s concrete implementation, not to `BaseSpecification.IsSatisfiedBy` (which always returns `false`). Forget `a.Relate(a)` in a constructor and you get silent wrong answers.
- **Interface nil trap**: a `Specification` variable holding a typed nil pointer (`var s *OverDueSpecification`) is not `nil` when compared as an interface; `s.IsSatisfiedBy(...)` will still dispatch, reaching the concrete method — but the zero `Specification` field in `BaseSpecification` will be `nil`, causing a panic at `And`/`Or`/`Not` call time.
- **Allocation per composition**: every `And`/`Or`/`Not` call allocates a new struct or closure on the heap. Build the composed specification once and reuse it; do not rebuild inside a tight loop.
- **Generic vs non-generic**: the functional-combinator form with generics (`Predicate[T]`) is type-safe and reusable across entity types. Pre-generics code (Go < 1.18) must either use `interface{}` with type assertions or repeat the pattern per entity type.

## See also

- skills/idiom/functional-options.md
- skills/behavioral/strategy.md
- examples/idiom/specification/
