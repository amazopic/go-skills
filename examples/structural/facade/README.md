# Facade — Example

## Pattern summary

Facade provides a simplified interface to a set of cooperating subsystems. Here, `Man` is the facade that orchestrates three independent subsystems (`House`, `Tree`, `Child`) behind a single `Todo()` call.

## Structure

| Type | Role |
|---|---|
| `Man` | Facade — wires up subsystems in `NewMan`; exposes `Todo()` |
| `House` | Subsystem A — `Build()` |
| `Tree` | Subsystem B — `Grow()` |
| `Child` | Subsystem C — `Born()` |

## What the facade hides

- Construction and wiring of three subsystem objects.
- The sequence in which subsystem methods must be called.
- The fact that three separate types are involved at all.

Callers see only `NewMan()` and `.Todo()`.

## Note

Subsystem types remain exported. Advanced callers can use `House`, `Tree`, and `Child` directly — the facade does not restrict access, it only provides a convenient default path.

## Run the tests

```bash
go test -race ./...
```

## Related skill

`skills/structural/facade.md`
