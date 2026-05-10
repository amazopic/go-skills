---
name: behavioral-mediator
description: Use when multiple objects communicate in complex ways and direct peer-to-peer coupling is getting tangled — route all interaction through a central coordinator instead.
category: behavioral
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Behavioral/Mediator/
example: examples/behavioral/mediator/
---

# Mediator

## Intent

Define an object that encapsulates how a set of objects interact. The mediator promotes loose coupling by keeping objects from referring to each other explicitly, and lets you vary the interaction independently from the participants.

## Context

Without a mediator, N peers that talk to each other need up to N² connections. The mediator collapses that to N connections (each peer ↔ mediator), and centralizes the coordination logic.

**Mediator vs. Observer**: these are easy to conflate.

- **Observer** is one-to-many: a subject broadcasts a change; observers react independently. No coordination logic.
- **Mediator** is many-to-many with rules: peers notify the mediator, which applies domain logic and may trigger *different* peers in a specific sequence (e.g., farmer sells tomatoes → cannery processes them → shop gets stock). The mediator *orchestrates*.

The example in `examples/behavioral/mediator/` models this distinction well: `Farmer.GrowTomato` notifies the mediator, which in turn calls `Cannery.MakeKetchup`, then the shop. The peers never hold references to each other — only to the mediator.

In Go, a simple event bus (channel + goroutines) can replace Mediator when the coordination rules are trivial. Use the Mediator struct when the rules between participants are non-trivial and need to live in one place.

## Implementation in Go

```go
// Mediator defines the notification contract.
type Mediator interface {
    Notify(msg string)
}

// ConcreteMediator embeds all colleagues and orchestrates their interactions.
type ConcreteMediator struct {
    *Farmer
    *Cannery
    *Shop
}

func (m *ConcreteMediator) Notify(msg string) {
    switch msg {
    case "Farmer: Tomato complete...":
        m.Cannery.AddMoney(-15000)
        m.Farmer.AddMoney(15000)
        m.Cannery.MakeKetchup(m.Farmer.GetTomato())
    case "Cannery: Ketchup complete...":
        m.Shop.AddMoney(-30000)
        m.Cannery.AddMoney(30000)
        m.Shop.SellKetchup(m.Cannery.GetKetchup())
    }
}

// ConnectColleagues wires peers to a single mediator instance.
func ConnectColleagues(f *Farmer, c *Cannery, s *Shop) {
    m := &ConcreteMediator{Farmer: f, Cannery: c, Shop: s}
    f.SetMediator(m)
    c.SetMediator(m)
    s.SetMediator(m)
}
```

For a lightweight event-bus Mediator (no domain rules, just routing):

```go
type Bus struct {
    mu   sync.RWMutex
    subs map[string][]func(payload any)
}

func (b *Bus) Subscribe(event string, fn func(any)) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.subs[event] = append(b.subs[event], fn)
}

func (b *Bus) Publish(event string, payload any) {
    b.mu.RLock()
    fns := b.subs[event]
    b.mu.RUnlock()
    for _, fn := range fns {
        fn(payload)
    }
}
```

Use the Bus variant when you need loose coupling without encoded orchestration rules. Switch to the struct Mediator when the coordination logic has domain significance.

## When to use

- Air-traffic control, chat rooms, workflow engines — any domain where peers must coordinate through a central authority.
- GUI components (form fields, buttons, dropdowns) that must stay in sync without direct field-to-field references.
- Microservice orchestration: a saga coordinator is a Mediator — it knows the sequence of calls; individual services do not.
- When refactoring a rat's nest of peer-to-peer calls — extract the coordination into a mediator to centralize it.

## When NOT to use

- When there are only two peers — just let them talk directly or use a simple callback.
- When coordination rules are trivial (broadcast only) — use Observer or a channel fan-out.
- When peers are in different services and the mediator becomes a distributed monolith — prefer an event bus with schema contracts instead.
- Avoid a mediator that grows into a "god object" knowing every detail of every peer — split by domain boundary when that happens.

## Gotchas

- **God-object creep**: the mediator's `Notify` switch-case grows without bound as new message types are added. Refactor into a method-dispatch table or sub-mediators when it exceeds ~5 cases.
- **Circular notification**: a colleague notifies the mediator, which calls back into the same colleague, which notifies the mediator again. Guard with a "processing" flag or event deduplication.
- **Concurrency**: if colleagues run concurrently, `Notify` must be goroutine-safe. The struct Mediator in the example is not — add a mutex or serialize via a channel.
- **Testing**: the mediator is the hardest part to unit test because it wires everything together. Inject a mock mediator into each colleague to test them in isolation.
- **Cyrillic function names**: the example has `СonnectСolleagues` with Cyrillic `С` characters — this compiles fine in Go but can surprise code reviewers. Rename to ASCII in production code.

## See also

- `skills/behavioral/observer.md` — Observer for simple 1→N broadcast; Mediator for N→N coordination with rules.
- `skills/behavioral/command.md` — Commands can be the payload the mediator routes between peers.
- `examples/behavioral/mediator/`
