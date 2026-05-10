# Specification — Example

Demonstrates the **Specification** pattern in Go: business rules are encapsulated as named types that implement `IsSatisfiedBy(Invoice) bool` and composed with `And`, `Or`, and `Not` combinators.

## Domain scenario

An invoice should be sent to a collection agency when ALL of the following are true:

1. It is overdue (`Day >= 30`)
2. Three or more notices have been sent (`Notice >= 3`)
3. It has NOT already been sent to the collection agency (`!IsSent`)

## Structure

| Type | Role |
|---|---|
| `Specification` | Interface: `IsSatisfiedBy`, `And`, `Or`, `Not`, `Relate` |
| `BaseSpecification` | Shared embedding; wires `And`/`Or`/`Not` factory methods |
| `AndSpecification` | Composite: true when both inner specs are satisfied |
| `OrSpecification` | Composite: true when either inner spec is satisfied |
| `NotSpecification` | Negates the wrapped spec |
| `OverDueSpecification` | Leaf: satisfied when `Day >= 30` |
| `NoticeSentSpecification` | Leaf: satisfied when `Notice >= 3` |
| `InCollectionSpecification` | Leaf: satisfied when `IsSent == true` |

## Bug fix (upstream)

The original `InCollectionSpecification.IsSatisfiedBy` returned `!elm.IsSent`, treating
`IsSent == false` (not yet sent) as "already in collection" — semantically inverted.

**Chosen semantics**: `IsSatisfiedBy` returns `true` when the invoice *is* already in the
collection (`IsSent == true`). Compose with `.Not()` to express the rule "not yet in
collection". The test error message was also corrected from inverted `false`/`true` labels
to the correct `true`/`false` expected/actual labels.

## Run

```bash
go test -race ./idiom/specification/
```

## Key points

- The `Relate(self)` call in each constructor is essential: it re-wires the embedded `Specification` interface to the concrete type, so `And`/`Or`/`Not` dispatch to the correct `IsSatisfiedBy` implementation.
- For simpler cases, prefer a `func(Invoice) bool` closure with functional `And`/`Or`/`Not` combinators — no struct embedding, no `Relate` ceremony, and generics make them type-safe.

See [skills/idiom/specification.md](../../../skills/idiom/specification.md) for full pattern documentation.
