# Bridge — Example

## Pattern summary

Bridge decouples an abstraction (`Car`/`Carer`) from its implementation (`Enginer`) so that both can vary independently. New car types and new engine types can be added without modifying each other.

## Structure

| Type | Role |
|---|---|
| `Carer` | Abstraction interface |
| `Car` | Refined abstraction — holds an `Enginer` field (the bridge) |
| `Enginer` | Implementation interface |
| `EngineSuzuki`, `EngineHonda`, `EngineLada` | Concrete implementations |

## How the bridge works

`Car.engine` is the bridge field. Any `Enginer` can be injected at construction time via `NewCar(engine Enginer)`. The `Rase()` method delegates entirely to `engine.GetSound()`. Three engines × any number of car types = no combinatorial explosion.

## Run the tests

```bash
go test -race ./...
```

## Related skill

`skills/structural/bridge.md`
