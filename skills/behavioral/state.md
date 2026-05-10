---
name: behavioral-state
description: Use when an object's behavior changes based on its internal state and the number of states or transitions is large enough to make a flat switch-case unmaintainable.
category: behavioral
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Behavioral/State/
example: examples/behavioral/state/
---

# State

## Intent

Allow an object to alter its behavior when its internal state changes. The object will appear to change its class. Instead of one big `switch` statement checked in every method, each state is an object that implements the full behavior for that state.

## Context

The State pattern is an object-oriented finite state machine. It shines when:

- There are many states (more than ~4).
- Each state has multiple methods with different behavior.
- Transitions between states are numerous and belong to the states themselves rather than to a god-switch in the context.

**State vs. switch-case**: for two or three states and one behavior method, a `switch` is clearer. The pattern earns its keep when you find yourself duplicating `switch state` across many methods, or when adding a new state requires edits in multiple places.

**State vs. Strategy**: they look identical in code — both use an interface to swap behavior. The difference is intent and lifecycle. Strategy is swapped by the *caller* and stays fixed for the operation's duration. State is swapped by the *state itself* or by the context as part of a state machine — transitions are automatic.

The example in `examples/behavioral/state/` models a mobile phone alert with two states (`MobileAlertVibration`, `MobileAlertSong`). The context (`MobileAlert`) delegates `Alert()` to its current state — clean and extensible.

## Implementation in Go

```go
// MobileAlertStater defines the behavior interface for each state.
type MobileAlertStater interface {
    Alert() string
}

// MobileAlert is the context — delegates to its current state.
type MobileAlert struct {
    state MobileAlertStater
}

func NewMobileAlert() *MobileAlert {
    return &MobileAlert{state: &MobileAlertVibration{}}
}

func (a *MobileAlert) Alert() string            { return a.state.Alert() }
func (a *MobileAlert) SetState(s MobileAlertStater) { a.state = s }

// MobileAlertVibration is the "vibrate" state.
type MobileAlertVibration struct{}

func (a *MobileAlertVibration) Alert() string { return "Vrrr... Brrr... Vrrr..." }

// MobileAlertSong is the "ring" state.
type MobileAlertSong struct{}

func (a *MobileAlertSong) Alert() string {
    return "Белые розы, Белые розы. Беззащитны шипы..."
}
```

For a state machine where states control their own transitions:

```go
type State interface {
    Handle(ctx *TrafficLight) State // returns the next state
}

type GreenLight struct{}

func (g *GreenLight) Handle(ctx *TrafficLight) State {
    fmt.Println("Green — go")
    return &YellowLight{}
}
```

The context just calls `current = current.Handle(ctx)` in a loop.

## When to use

- Network connection lifecycle: Disconnected → Connecting → Connected → Closing.
- Order workflow: Pending → Paid → Shipped → Delivered → Refunded.
- Game entity AI: Idle → Patrolling → Chasing → Attacking → Fleeing.
- Document approval: Draft → Review → Approved → Published → Archived.
- Any place you find `switch state { case A: ...; case B: ...; }` scattered across multiple methods.

## When NOT to use

- Two states and one method — use a bool flag and an `if`.
- Transitions are data-driven (loaded from config or database) — use a table-driven state machine (`map[State]map[Event]State`) instead.
- The context has no methods that change behavior by state — just track the state as a value with no need for the pattern.
- Goroutine-per-state machines are sometimes cleaner for concurrent state + event handling than the synchronous pattern.

## Gotchas

- **Nil state**: if `MobileAlert.state` is nil, calling `Alert()` panics. Always initialize to a concrete start state; use `NewMobileAlert()` constructors rather than struct literals.
- **Concurrent state transitions**: `SetState` is not goroutine-safe. If multiple goroutines can trigger transitions, protect `state` with a mutex or use atomic pointer swap (`atomic.Pointer[MobileAlertStater]` in Go 1.19+).
- **Transition explosion**: if every state needs to know about every other state (to create it), you get tight coupling. Use a factory or the context itself to create next-states to break cycles.
- **Lost events during transition**: if an event arrives while a transition is in progress, it may be processed by the old or new state depending on timing. Design your transition protocol carefully in concurrent systems.
- **Interface allocation**: each state transition allocates a new concrete state struct. In hot paths (tight loops), consider a pool of state objects or a value-type state enum with a dispatch table.

## See also

- `skills/behavioral/strategy.md` — Strategy for externally selected algorithms; State for internally driven behavior changes.
- `skills/behavioral/memento.md` — Snapshot the current state node for undo/restore.
- `examples/behavioral/state/`
