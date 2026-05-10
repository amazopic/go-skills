# Abstract Factory — Example

Demonstrates the **Abstract Factory** pattern in Go: a single `AbstractFactory` interface produces a matched family of related objects (`AbstractWater` + `AbstractBottle`). Swapping the concrete factory changes every product it creates without modifying consumer code.

## Structure

| Type | Role |
|---|---|
| `AbstractFactory` | Factory interface; declares `CreateWater` and `CreateBottle` |
| `AbstractWater` | Product interface for the liquid |
| `AbstractBottle` | Product interface for the container |
| `CocaColaFactory` | Concrete factory — produces `CocaColaWater` and `CocaColaBottle` |
| `CocaColaWater` | Concrete water product |
| `CocaColaBottle` | Concrete bottle product |

## Run

```bash
go test -race ./creational/abstract-factory/
```

## Key points

- Consumer code (`TestAbstractFactory`) works only with the two abstract interfaces; it never references a concrete type.
- Adding a `PepsiFactory` family requires zero changes to consumer code or product interfaces.
- For a factory that is not plug-in swappable, prefer a concrete struct with `New*()` methods — no interface needed.

See [skills/creational/abstract-factory.md](../../../skills/creational/abstract-factory.md) for full pattern documentation.
