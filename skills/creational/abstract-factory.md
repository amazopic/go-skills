---
name: creational-abstract-factory
description: Abstract Factory — produce families of related objects through a single factory interface, keeping the consumer agnostic of which concrete family is in use. Use when multiple related types must vary together as a matched set.
category: creational
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Creational/AbstractFactory/
example: examples/creational/abstract-factory/
---

# Abstract Factory

## Intent

Provide an interface for creating families of related or dependent objects without specifying their concrete types. The factory itself is the abstraction; swapping the factory changes every product it produces in a coordinated, consistent way.

## Context

Abstract Factory is the natural step up from Factory Method when you have multiple related types that must be kept in sync. A UI toolkit that must render either native macOS widgets or web widgets is the textbook case: you swap one `WidgetFactory` and every button, dialog, and scrollbar it creates belongs to the same family.

In Go, the full GoF pattern is rarely the right first tool. Before reaching for an interface with multiple `Create*` methods, ask:

1. **Does only one type vary?** Use a simple factory function or [Factory Method](skills/creational/factory-method.md) instead.
2. **Do the products share state only through construction?** A plain struct with configuration fields and multiple `New*()` methods on it is easier to read and test than a factory interface.
3. **Must the factory itself be swappable at runtime by external code (plugin, DI container)?** Only then does the full interface form earn its keep.

## Implementation in Go

### Simple form — struct factory with configuration (recommended)

When the factory is not plug-in swappable, a concrete struct with `New*()` methods is cleaner and needs no interface:

```go
// DrinkFactory creates a matched set of drink products for a given brand.
type DrinkFactory struct {
    brand string
}

func NewDrinkFactory(brand string) *DrinkFactory {
    return &DrinkFactory{brand: brand}
}

func (f *DrinkFactory) NewWater(volumeL float64) *Water {
    return &Water{brand: f.brand, volumeL: volumeL}
}

func (f *DrinkFactory) NewBottle(capacityL float64) *Bottle {
    return &Bottle{brand: f.brand, capacityL: capacityL}
}
```

The consumer receives a `*DrinkFactory` and calls `f.NewWater(...)` and `f.NewBottle(...)`. Changing the brand is a one-line constructor swap.

### Full interface form — plug-in swappable factory

Use this when the factory must be injected and replaced (tests, plugins, multi-tenant systems):

```go
// Factory creates a matched family of drink products.
type Factory interface {
    CreateWater(volumeL float64) Water
    CreateBottle(capacityL float64) Bottle
}

// Water and Bottle are the product abstractions.
type Water interface {
    Volume() float64
}

type Bottle interface {
    Fill(Water)
    WaterVolume() float64
    Capacity() float64
}

// CocaColaFactory is one concrete family.
type CocaColaFactory struct{}

func (CocaColaFactory) CreateWater(v float64) Water   { return &colaWater{v} }
func (CocaColaFactory) CreateBottle(c float64) Bottle { return &colaBottle{capacity: c} }

// PepsiFactory is the other concrete family.
type PepsiFactory struct{}

func (PepsiFactory) CreateWater(v float64) Water   { return &pepsiWater{v} }
func (PepsiFactory) CreateBottle(c float64) Bottle { return &pepsiBottle{capacity: c} }
```

Consumer code works with `Factory` only:

```go
func FillBottle(f Factory, vol float64) Bottle {
    water  := f.CreateWater(vol)
    bottle := f.CreateBottle(vol)
    bottle.Fill(water)
    return bottle
}
```

Swapping `CocaColaFactory{}` for `PepsiFactory{}` changes the entire product family without touching `FillBottle`.

## When to use

- You need to produce multiple related types (a "family") and the whole family must change together — e.g., cloud provider clients, themed UI components, test doubles for a subsystem.
- The factory must be injected and swapped at runtime without the consumer knowing the concrete family.
- You want to enforce that products from different families are not accidentally mixed (a `CocaColaBottle` filled with `PepsiWater` is a logical error the type system can prevent when both are created through the same factory).

## When NOT to use

- Only one product type varies — use a factory function or [Factory Method](skills/creational/factory-method.md).
- The factory is not swappable — a concrete struct with `New*()` methods is simpler and equally testable.
- The family has only a single product — Abstract Factory is overkill; reach for [Functional Options](skills/idiom/functional-options.md) or a plain constructor.
- Products in the family do not interact or share invariants — there is no coupling to protect; independent factory functions are cleaner.

## Gotchas

- **Interface explosion**: every new product type forces a new method on the factory interface, recompiling every implementation. If the family grows frequently, prefer the struct form so adding a method does not break existing factories.
- **Product interface coupling**: when product interfaces are too wide, they leak implementation details. Keep product interfaces narrow — only what the consumer actually needs.
- **Value vs pointer receivers on factory structs**: if the factory holds no mutable state (most do not), define factory methods on value receivers (`func (CocaColaFactory) CreateWater(...)`) so the factory can be stored by value, avoiding an unnecessary heap allocation.
- **Nil product interfaces**: an uninitialized `Factory` variable is a nil interface; calling `Create*` on it panics. Fail fast at construction time with a `if f == nil { panic("factory required") }` guard rather than discovering the nil deep in a call stack.

## See also

- skills/creational/factory-method.md
- skills/creational/builder.md
- skills/idiom/functional-options.md
- examples/creational/abstract-factory/
