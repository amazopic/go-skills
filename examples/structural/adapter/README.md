# Adapter — Example

## Pattern summary

Adapter converts the interface of an existing type (`Adaptee`) into the interface a client expects (`Target`), without modifying the original type.

## Structure

| Type | Role |
|---|---|
| `Target` | Interface the client expects |
| `Adaptee` | Existing type with an incompatible method (`SpecificRequest`) |
| `Adapter` | Wraps `Adaptee` via embedding; satisfies `Target` by translating `Request()` → `SpecificRequest()` |

## Key design choices

- `Adapter` embeds `*Adaptee` so it inherits all of `Adaptee`'s methods automatically and only overrides the one that needs translation.
- `NewAdapter` returns `Target` (not `*Adapter`), keeping the concrete type opaque to callers.
- No external dependencies — stdlib only.

## Run the tests

```bash
go test -race ./...
```

## Related skill

`skills/structural/adapter.md`
