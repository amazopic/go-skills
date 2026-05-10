---
name: behavioral-chain-of-responsibility
description: Use when a request must pass through a sequence of handlers, each deciding to process or forward it — middleware, validation pipelines, log routing.
category: behavioral
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Behavioral/ChainOfResponsibility/
example: examples/behavioral/chain-of-responsibility/
---

# Chain of Responsibility

## Intent

Decouple a sender from its receivers by giving multiple objects a chance to handle a request. Handlers are linked in a chain; each either handles the request or passes it to its successor, avoiding a hard-coded if-else ladder spread across the caller.

## Context

Chain of Responsibility appears whenever you have an ordered set of processing stages where (a) the right handler is not known at compile time, or (b) the set of handlers changes at runtime. Classic Go examples include HTTP middleware stacks, validation pipelines, log-level routing, and permission checks.

Go's most idiomatic CoR form is the **middleware function**: `func(http.Handler) http.Handler`. The standard library's `net/http` package is built on exactly this shape. The GoF struct-with-a-next-pointer approach still makes sense when handlers need to carry significant per-instance state or when you need to introspect the chain.

Two shapes to know:

1. **Linked-node (GoF)**: each handler holds a `next Handler` pointer. Explicit, easy to debug, slightly verbose.
2. **Slice of handlers + loop**: no explicit chaining; a coordinator iterates handlers until one returns. Simpler when all handlers live in one package and the chain is fixed.

The example in `examples/behavioral/chain-of-responsibility/` uses the linked-node shape, which maps cleanly to the GoF description.

## Implementation in Go

```go
// Handler is the chain link.
type Handler interface {
    SendRequest(message int) string
}

// ConcreteHandlerA handles message 1; forwards the rest.
type ConcreteHandlerA struct {
    next Handler
}

func (h *ConcreteHandlerA) SendRequest(message int) string {
    if message == 1 {
        return "Im handler 1"
    }
    if h.next != nil {
        return h.next.SendRequest(message)
    }
    return ""
}

// Build the chain: A → B → C
handlers := &ConcreteHandlerA{
    next: &ConcreteHandlerB{
        next: &ConcreteHandlerC{},
    },
}
result := handlers.SendRequest(2) // "Im handler 2"
```

For HTTP-style middleware the idiomatic Go form is even leaner:

```go
type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, mw ...Middleware) http.Handler {
    for i := len(mw) - 1; i >= 0; i-- {
        h = mw[i](h)
    }
    return h
}
```

No structs, no interface boilerplate — just functions composed in order.

## When to use

- HTTP middleware (auth, logging, rate-limiting, tracing) — the canonical Go case.
- Validation pipelines where early-exit on first failure is correct.
- Log routing: DEBUG → INFO → WARN → ERROR with different sinks per level.
- Plugin or handler registries where the chain is assembled at startup from configuration.
- Any place you would otherwise write a long `switch` or `if-else` in the *caller*.

## When NOT to use

- When every handler must always run (use a slice + range loop instead).
- When only one handler will ever exist — just call it directly.
- When the chain never changes and all handlers are in one function — a plain `switch` is clearer.
- When you need parallel fan-out — use goroutines + channels, not a chain.

## Gotchas

- **Nil next pointer**: always guard with `if h.next != nil` before forwarding, or use a sentinel no-op handler as the tail.
- **Infinite loops**: if a handler accidentally adds itself as its own `next`, you get a stack overflow. Keep chain construction separate from handler logic.
- **Silent drops**: a request that no handler claims returns silently. Return a sentinel error or log a warning at the tail.
- **Ordering matters**: inserting a handler in the wrong position (e.g., auth after business logic) is a runtime bug the compiler cannot catch.
- **Shared state**: if handlers share mutable state (a counter, a buffer), protect it with a mutex or move state into the request context.

## See also

- `skills/behavioral/command.md` — encapsulates a request as an object; often used *inside* CoR handlers.
- `skills/behavioral/decorator.md` — similar wrapping structure but adds behavior rather than routing a request.
- `examples/behavioral/chain-of-responsibility/`
