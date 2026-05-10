---
name: behavioral-memento
description: Use when you need to snapshot and restore an object's internal state — undo/redo, checkpoints, rollback — without violating its encapsulation.
category: behavioral
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Behavioral/Memento/
example: examples/behavioral/memento/
---

# Memento

## Intent

Without violating encapsulation, capture and externalize an object's internal state so it can be restored to that state later. The owner of the state (Originator) creates a snapshot (Memento); the snapshot is held by a Caretaker that treats it as opaque.

## Context

Memento appears wherever "undo" or "time travel" is needed: text editors, transaction rollback, game save points, configuration checkpoints. The key constraint is that the Caretaker must not be able to read or modify the snapshot content — it just holds it.

In Go, encapsulation is enforced by package boundaries. The classic Go Memento approach:

- Define `Memento` with **unexported fields** in the same package as the Originator.
- Expose only `Snapshot() Memento` and `Restore(m Memento)` on the Originator.
- The Caretaker (possibly in another package) holds `[]Memento` but cannot touch the fields.

The example in `examples/behavioral/memento/` follows this shape: `Memento.state` is unexported, `Originator.CreateMemento()` captures it, `Originator.SetMemento()` restores it, and `Caretaker` holds a pointer to the snapshot.

**Serialization consideration**: if snapshots must survive process restarts (game saves, database WAL), add a `MarshalBinary() / UnmarshalBinary()` pair to `Memento` keeping the data layout internal to the package. Avoid directly marshalling the Originator — that leaks its structure.

## Implementation in Go

```go
// Originator owns the mutable state.
type Originator struct {
    State string
}

// CreateMemento captures current state as an opaque snapshot.
func (o *Originator) CreateMemento() *Memento {
    return &Memento{state: o.State}
}

// SetMemento restores state from a snapshot.
func (o *Originator) SetMemento(m *Memento) {
    o.State = m.GetState()
}

// Memento is the opaque snapshot — state is unexported.
type Memento struct {
    state string
}

func (m *Memento) GetState() string { return m.state }

// Caretaker holds the snapshot; cannot read its contents.
type Caretaker struct {
    Memento *Memento
}
```

Usage:

```go
o := &Originator{State: "On"}
ct := &Caretaker{Memento: o.CreateMemento()} // save
o.State = "Off"
o.SetMemento(ct.Memento)                      // restore
// o.State == "On"
```

For undo stacks, the Caretaker holds a slice:

```go
type History struct{ stack []*Memento }

func (h *History) Push(m *Memento)       { h.stack = append(h.stack, m) }
func (h *History) Pop() (*Memento, bool) {
    if len(h.stack) == 0 {
        return nil, false
    }
    n := len(h.stack) - 1
    m := h.stack[n]
    h.stack = h.stack[:n]
    return m, true
}
```

## When to use

- Undo/Redo in editors, drawing apps, form wizards.
- Database-style savepoints: begin, do work, rollback to savepoint.
- Game checkpoints: serialize player state before entering a boss battle.
- Configuration management: snapshot current config before applying a potentially broken change.
- Speculative execution: try a transformation, revert if it fails validation.

## When NOT to use

- When the state is already stored externally (database row, file) — just reload from the source.
- When snapshots are large and frequent — copying multi-megabyte state per keystroke is impractical; use a diff/patch approach instead.
- When the Originator's state is already immutable — immutable value types are their own memento; just keep a copy.
- Don't use Memento as a general-purpose serialization mechanism — use encoding/json or encoding/gob for that and accept the coupling.

## Gotchas

- **Shallow copy**: if `Originator.State` contains pointers or slices, `CreateMemento` must deep-copy them. A naive assignment copies the pointer, not the data, so modifying the live object also mutates the "snapshot".
- **Package boundary**: Go unexported fields only protect within a package. If Originator and Caretaker are in separate packages, `Memento` fields must be exported (or use an opaque `[]byte` snapshot). Decide your encapsulation boundary up front.
- **Memory growth**: an undo stack with no cap leaks memory proportional to the number of edits. Implement a maximum depth with LRU eviction.
- **Nil Memento**: `SetMemento(nil)` panics. Guard: `if m == nil { return }`.
- **Concurrent snapshots**: taking a snapshot while another goroutine mutates the Originator races. Lock before `CreateMemento` and release after. Alternatively, make `State` a value type and rely on the copy being race-free.

## See also

- `skills/behavioral/command.md` — Command's Undo operation typically restores a Memento.
- `skills/behavioral/state.md` — State machines can snapshot their current state node as a Memento.
- `examples/behavioral/memento/`
