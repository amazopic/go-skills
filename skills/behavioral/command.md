---
name: behavioral-command
description: Use when you need to encapsulate an operation as a first-class value — for queuing, undo/redo, logging, or deferred execution.
category: behavioral
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Behavioral/Command/
example: examples/behavioral/command/
---

# Command

## Intent

Encapsulate a request as an object (or function value) so that callers and executors are decoupled. The invoker knows only how to execute a command; it does not know what the command does or who performs it.

## Context

The GoF Command pattern encodes an operation as a struct implementing a single `Execute() T` method. In Go this collapses neatly to a **function value** in the common case. Use the full struct form only when you need:

1. **Undo/Redo** — the struct can carry an inverse operation or the pre-change state.
2. **Serialization / queuing** — commands must be encoded (JSON, protobuf) and dispatched across a queue or stored in a log.
3. **Rich metadata** — timestamps, retry counts, correlation IDs that travel with the operation.
4. **Multiple methods on one command** — `Execute`, `Undo`, `Validate`, `Describe`.

For everything else, `type Command func() error` is the idiomatic Go answer and needs no pattern at all. Don't reach for an interface with a single method when a function type already is that interface.

The example in `examples/behavioral/command/` shows the GoF struct approach with `Command` interface, `Invoker`, and `Receiver` — a good teaching example for cases where struct richness is warranted (the Invoker queues and can UnStore commands before execution).

## Implementation in Go

```go
// Command provides the command interface.
type Command interface {
    Execute() string
}

// Receiver knows how to perform the actual work.
type Receiver struct{}

func (r *Receiver) ToggleOn() string  { return "Toggle On" }
func (r *Receiver) ToggleOff() string { return "Toggle Off" }

// Invoker queues and fires commands.
type Invoker struct {
    commands []Command
}

func (i *Invoker) StoreCommand(c Command)   { i.commands = append(i.commands, c) }
func (i *Invoker) UnStoreCommand()          {
    if len(i.commands) > 0 {
        i.commands = i.commands[:len(i.commands)-1]
    }
}
func (i *Invoker) Execute() string {
    var out string
    for _, c := range i.commands {
        out += c.Execute() + "\n"
    }
    return out
}
```

The lightweight Go-native alternative for a simple fire-and-forget queue:

```go
type Job func() error

type Queue struct{ jobs []Job }

func (q *Queue) Add(j Job)     { q.jobs = append(q.jobs, j) }
func (q *Queue) Flush() error  {
    for _, j := range q.jobs {
        if err := j(); err != nil {
            return err
        }
    }
    q.jobs = q.jobs[:0]
    return nil
}
```

## When to use

- Undo/Redo stacks (text editors, drawing tools, database transaction logs).
- Work queues or job schedulers where operations are produced and consumed separately.
- Macro recording — capture a sequence of user actions for replay.
- Transactional batching — collect commands, validate the batch, then commit or abort.
- Audit logging — the command object carries context (who, when, what) alongside the action.

## When NOT to use

- When you just need to call a function later — use `func() error` directly; no interface needed.
- When there is no invoker / no deferred execution — call the receiver directly.
- When undo is not required and the operation is stateless — a plain function is clearer and leaner.
- Avoid wrapping every method call in a Command "just in case" — that is over-engineering.

## Gotchas

- **Captured mutable state**: command closures or structs that capture a pointer to shared state can race if the invoker is concurrent. Make commands capture values, not pointers, or synchronize.
- **Nil receiver**: a `Command` interface value holding a nil `*ConcreteCommand` is non-nil; calling `Execute()` on it panics. Ensure `StoreCommand` rejects nil values.
- **Growing queue**: an `Invoker` that only appends and never flushes leaks memory. Decide on clear ownership of when the queue drains.
- **Undo ordering**: undo must replay in reverse order. Off-by-one bugs in stack management are the most common defect in undo implementations.
- **Serialization coupling**: if commands must be serialized, add a `type` discriminator field and a registry; without it, JSON unmarshalling cannot reconstruct the concrete type.

## See also

- `skills/behavioral/chain-of-responsibility.md` — CoR passes a request through handlers; Command encapsulates the request itself.
- `skills/behavioral/memento.md` — Memento stores the state snapshot that Command's Undo restores.
- `examples/behavioral/command/`
