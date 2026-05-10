# Flyweight — Example

## Pattern summary

Flyweight shares intrinsic (immutable, shared) state across many objects and passes extrinsic (context-specific) state per call. A factory ensures each unique intrinsic state is allocated only once.

## Structure

| Type | Role |
|---|---|
| `Flyweighter` | Interface exposing `Draw(width, height int, opacity float64)` |
| `ConcreteFlyweight` | Holds intrinsic state (`filename`); extrinsic state arrives via `Draw` |
| `FlyweightFactory` | Caches `Flyweighter` instances by filename; creates on first request |

## State split

| State | Location | Example value |
|---|---|---|
| Intrinsic | `ConcreteFlyweight.filename` | `"cat.jpg"` |
| Extrinsic | `Draw` parameters | `width=100, height=100, opacity=0.95` |

Two calls with `"cat.jpg"` return the **same** `*ConcreteFlyweight`. Only `Draw` is invoked again with different extrinsic arguments — no second allocation.

## Concurrency note

The factory's `pool` map is **not** goroutine-safe as written. For concurrent use, protect it with a `sync.RWMutex` (see `skills/structural/flyweight.md` for the double-checked locking snippet).

## Run the tests

```bash
go test -race ./...
```

## Related skill

`skills/structural/flyweight.md`
