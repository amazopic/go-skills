---
name: structural-bridge
description: Use when an abstraction and its implementation must vary independently. Compose with an interface field instead of locking the two axes together through inheritance.
category: structural
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Structural/Bridge/
example: examples/structural/bridge/
---

# Bridge

## Intent

Bridge decouples an abstraction from its implementation so that the two can vary independently. Instead of a combinatorial explosion of subtypes (M abstractions × N implementations = M×N classes), you compose one axis into the other via an interface field — keeping both axes open to extension without touching existing code.

## Context

The classic symptom is a type hierarchy that grows in two dimensions at once. Imagine shapes and renderers: `CircleVectorRenderer`, `CircleRasterRenderer`, `SquareVectorRenderer`, `SquareRasterRenderer`. Add a third renderer and you need two new types; add a new shape and you need one new type per renderer. Bridge collapses that to `Shapes + Renderers`, combined at runtime.

In Go, the pattern surfaces naturally wherever a struct holds an interface field injected via constructor — which is extremely common. Bridge is less a dramatic refactoring and more a reminder that the injected interface field is doing architectural work: it is the bridge itself.

Compare with plain interface use: if only one axis varies (you swap implementations but the abstraction is fixed), a single interface is sufficient. Bridge is the right name when **both** axes vary and evolve independently.

## Implementation in Go

The example in `examples/structural/bridge/` models cars and engines. `Car` (abstraction) holds an `Enginer` (implementation interface). Any engine snaps into any car; new engines or new car types are independent additions.

```go
// Enginer is the "implementation" side of the bridge.
type Enginer interface {
    GetSound() string
}

// Carer is the "abstraction" side.
type Carer interface {
    Race() string
}

// Car is the refined abstraction — it delegates to the engine.
type Car struct {
    engine Enginer // <-- this field IS the bridge
}

func NewCar(engine Enginer) Carer {
    return &Car{engine: engine}
}

func (c *Car) Race() string {
    return c.engine.GetSound()
}

// Concrete implementations — completely independent of Car.
type EngineSuzuki struct{}
func (e *EngineSuzuki) GetSound() string { return "SssuuuuZzzuuuuKkiiiii" }

type EngineHonda struct{}
func (e *EngineHonda) GetSound() string { return "HhoooNnnnnnnnnDddaaaaaaa" }
```

Usage:

```go
car := NewCar(&EngineSuzuki{})
fmt.Println(car.Race()) // SssuuuuZzzuuuuKkiiiii

car2 := NewCar(&EngineHonda{})
fmt.Println(car2.Race()) // HhoooNnnnnnnnnDddaaaaaaa
```

No new types needed to mix any car with any engine.

### When a plain interface is enough

If your abstraction has only one concrete form, there is no Bridge — there is just good interface usage. Bridge implies at least two concrete `Car` types in addition to multiple `Enginer` implementations. Do not name something a Bridge just because it has an interface field.

## When to use

- Two independent dimensions of variation cause a type-count explosion under single inheritance.
- You want to swap implementations at runtime (e.g., swap a real database engine for an in-memory one in tests) without changing the abstraction.
- Platform portability: the abstraction is stable; multiple platform-specific implementations evolve separately.
- A plugin architecture where the "host" abstraction and "guest" implementations are maintained by different teams.

## When NOT to use

- Only one axis varies — a plain interface is simpler and carries less naming overhead.
- The abstraction and implementation never change after initial design — Bridge's flexibility adds complexity for no gain.
- You find yourself adding a Bridge "just in case" — GoF structural patterns should solve a present problem, not a hypothetical one.

## Gotchas

- **Nil implementation field.** If `engine` is not set before `Race()` is called, you get a nil-pointer panic. Enforce injection in the constructor and return an error or panic early rather than deferring the failure.
- **Over-abstraction.** Two interfaces with one implementation each is not Bridge — it is premature abstraction. Wait until you actually have two concrete implementations on each axis.
- **Interface bloat on the implementation side.** Keep `Enginer` minimal (one or two methods). A fat implementation interface defeats the point: implementations become hard to swap and easy to break.
- **Concurrency.** If the bridge field can be swapped at runtime, guard it with a `sync.RWMutex`. Unsynchronised field replacement is a data race under `-race`.

## See also

- skills/structural/adapter.md — translates an interface; Bridge composes two independent interfaces
- skills/structural/strategy.md — similar composition; Strategy focuses on swappable algorithms, Bridge on decoupling two structural axes
- examples/structural/bridge/
