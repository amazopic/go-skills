---
name: creational-singleton
description: Singleton — exactly one instance, lazy-initialized via sync.Once. Use sparingly; prefer DI; only when an instance must truly be process-global (e.g., metric registries).
category: creational
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/creational/singleton.md
example: examples/creational/singleton/
---

# Singleton Pattern

Singleton creational design pattern restricts the instantiation of a type to a single object.

## Implementation

```go
package singleton

type singleton map[string]string

var (
    once sync.Once

    instance singleton
)

func New() singleton {
	once.Do(func() {
		instance = make(singleton)
	})

	return instance
}
```

## Usage

```go
s := singleton.New()

s["this"] = "that"

s2 := singleton.New()

fmt.Println("This is ", s2["this"])
// This is that
```

## Rules of Thumb

- Singleton pattern represents a global state and most of the time reduces testability.
