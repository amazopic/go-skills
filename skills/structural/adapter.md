---
name: structural-adapter
description: Use when you need to make an existing type fit an interface it was not designed for, without modifying the original type. Thin wrapper — translates one API surface to another.
category: structural
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Structural/Adapter/
example: examples/structural/adapter/
---

# Adapter

## Intent

Adapter converts the interface of an existing type into the interface a client expects, enabling types that would otherwise be incompatible to cooperate. It is the "plug adapter" of software design: no internals change, only the shape of the socket.

## Context

Adapter appears whenever you must integrate a third-party library, a legacy package, or a stdlib primitive into a codebase that has defined its own interface contract. It is also the mechanism behind every `io.Reader`/`io.Writer` adapter in the standard library (`strings.NewReader`, `bytes.NewBuffer`, `gzip.NewWriter`, etc.) — the canonical Go examples of the pattern in the wild.

Reach for Adapter when:
- You cannot (or should not) modify the source type.
- The mismatch is purely structural: the methods exist but have different names, signatures, or calling conventions.
- A functional option or middleware chain would be overkill for a one-off translation.

If you control both sides, prefer redesigning the interface rather than layering an adapter.

## Implementation in Go

Go's small, implicit interfaces make Adapter especially lightweight. A struct that embeds the adaptee gets all its methods for free and only needs to override the ones that need translation. The adapter is rarely more than 15 lines.

```go
// Target is the interface the rest of the codebase expects.
type Target interface {
    Request() string
}

// Adaptee is a third-party type with an incompatible method name.
type Adaptee struct{}

func (a *Adaptee) SpecificRequest() string { return "Request" }

// Adapter wraps Adaptee and satisfies Target.
// Embedding gives us Adaptee's methods; we override only what differs.
type Adapter struct {
    *Adaptee
}

func (a *Adapter) Request() string {
    return a.SpecificRequest() // translate the call
}

// NewAdapter is the constructor; callers only see Target.
func NewAdapter(adaptee *Adaptee) Target {
    return &Adapter{adaptee}
}
```

This is the shape used in `examples/structural/adapter/`. The embedding means `Adapter` inherits any future methods added to `Adaptee` at no cost, while still presenting the `Target` interface to callers.

### Functional alternative

When the interface has a single method (very common in Go), a plain function type adapter is even simpler:

```go
type TargetFunc func() string

func (f TargetFunc) Request() string { return f() }

// Wrap an incompatible function in one line — no struct needed.
var t Target = TargetFunc(existingFunc)
```

## When to use

- Integrating a third-party package whose types do not match your domain interface.
- Wrapping stdlib primitives (`http.ResponseWriter`, `io.Reader`) to add domain-specific methods.
- Testing: adapting a real dependency to a mock-friendly interface without touching the dependency.
- Making two packages developed independently interoperate.

## When NOT to use

- You own both sides — just change one of them.
- The mismatch is in behaviour, not interface shape; Adapter only translates signatures.
- You need to add cross-cutting concerns (logging, retries) to every call — use Decorator instead.
- The adaptee has dozens of methods and you only care about three — consider a narrow interface and a thin delegation struct without embedding.

## Gotchas

- **Embedding leaks the adaptee's full API.** If `Adaptee` has 20 methods and `Target` has 1, embedding exposes all 20. Use a named field (`adaptee *Adaptee`) and delegate explicitly when the public surface must stay tight.
- **Pointer vs value receiver confusion.** Embedding `*Adaptee` satisfies interfaces that require pointer receivers on `Adaptee`. Embedding `Adaptee` (value) does not — a common compile error when wrapping third-party types.
- **Nil adaptee.** Constructor functions should guard against a nil adaptee if the adapter will be stored and called later.
- **Interface widening.** Returning `Target` from the constructor (not `*Adapter`) prevents callers from casting back to the concrete type and bypassing the adapter — maintain the abstraction boundary.

## See also

- skills/structural/decorator.md — adds behaviour; Adapter only translates the interface
- skills/structural/proxy.md — controls access; Adapter changes the API shape
- examples/structural/adapter/
