---
name: behavioral-template-method
description: Use when multiple algorithms share a fixed skeleton but differ in specific steps — define the skeleton once and delegate the variable steps to concrete implementations.
category: behavioral
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Behavioral/TemplateMethod/
example: examples/behavioral/template-method/
---

# Template Method

## Intent

Define the skeleton of an algorithm in a base type, deferring some steps to concrete implementations. The skeleton stays fixed; the variable steps are overridable without changing the overall structure.

## Context

The GoF Template Method relies on inheritance: a base class implements the invariant steps and abstract methods that subclasses override. Go has no inheritance, so the pattern must be translated.

**Two Go-native translations:**

1. **Interface hook + embedding** (what the example uses): embed a concrete struct that holds the template, and compose it with an interface providing the variable steps. The embedded struct's method calls `q.Open()` and `q.Close()` through the interface — exactly what the example does with `Quotes` embedding `QuotesInterface`.

2. **First-class functions** (often simpler): pass the variable steps as `func` values. This beats the struct form when there are one or two hooks and no shared state.

For most real Go code, **composition + functions wins over Template Method**. The pattern adds value only when the algorithm skeleton is non-trivial (many shared steps) and the variable steps need to be grouped into coherent implementations (e.g., FrenchQuotes bundles `Open` and `Close` as a unit that travels together).

The example in `examples/behavioral/template-method/` cleanly demonstrates the Go embedding approach: `Quotes` embeds `QuotesInterface`, and `Quotes.Quotes(str)` is the template method calling `q.Open()` and `q.Close()`. `FrenchQuotes` and `GermanQuotes` provide the hook implementations.

## Implementation in Go

```go
// QuotesInterface defines the variable steps ("hooks").
type QuotesInterface interface {
    Open() string
    Close() string
}

// Quotes is the template — its Quotes() method is the skeleton.
type Quotes struct {
    QuotesInterface // embedded: provides Open() and Close()
}

// Quotes is the Template Method — invariant skeleton.
func (q *Quotes) Quotes(str string) string {
    return q.Open() + str + q.Close()
}

func NewQuotes(qt QuotesInterface) *Quotes {
    return &Quotes{qt}
}

// FrenchQuotes implements the «French» hook.
type FrenchQuotes struct{}

func (q *FrenchQuotes) Open() string  { return "«" }
func (q *FrenchQuotes) Close() string { return "»" }
```

Usage:

```go
qt := NewQuotes(&FrenchQuotes{})
qt.Quotes("Bonjour") // "«Bonjour»"
```

**Function-based alternative** (no struct needed for simple cases):

```go
func WrapString(s string, open, close func() string) string {
    return open() + s + close()
}

result := WrapString("hello", func() string { return "«" }, func() string { return "»" })
```

## When to use

- Parsing pipelines where the outer loop (read-parse-validate-emit) is fixed but the format-specific steps change.
- Report generation: the report structure (header, body, footer) is invariant; the rendering (HTML, PDF, CSV) varies.
- Test fixtures: a `TestSuite` template runs setup, test, teardown in order; concrete suites provide the test logic.
- Database migrations: the migration runner orchestrates ordering and error handling; each migration implements `Up()` / `Down()`.
- When the algorithm has several meaningful shared steps and the hook implementations belong together as a coherent type.

## When NOT to use

- When there is only one hook — pass a `func` instead; no pattern needed.
- When the skeleton itself varies by context — use Strategy, which swaps the whole algorithm, not just steps.
- When hooks are independent — inject each as a separate `func` rather than grouping them into an interface.
- Avoid Template Method as a default choice just because the code has some shared steps — check if a plain helper function covers it first.

## Gotchas

- **Nil interface panic**: embedding `QuotesInterface` and constructing `&Quotes{}` without setting the interface field results in a nil embedded interface. Calling `q.Open()` panics. Always use the constructor `NewQuotes(qt)`.
- **Infinite recursion via embedding**: if the template method and the hook method share the same name, the embedded call dispatches through the outer struct, not the hook, and recurses. Keep template method names distinct from hook method names.
- **Over-engineering**: Template Method is the most misused GoF pattern. Before reaching for it, ask "does a plain function with function parameters solve this?" — often it does.
- **Thread safety**: if the embedded hook implementations carry mutable state, concurrent calls to the template method race. Make hook implementations stateless or synchronize them.

## See also

- `skills/behavioral/strategy.md` — Strategy swaps the whole algorithm; Template Method keeps the skeleton and varies only the steps.
- `skills/behavioral/chain-of-responsibility.md` — CoR is a runtime chain; Template Method is a compile-time skeleton.
- `examples/behavioral/template-method/`
