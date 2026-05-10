---
name: structural-decorator
description: Decorator — wrap a value in a function or type that adds behavior while preserving the same interface. Use for cross-cutting concerns: logging, caching, retry.
category: structural
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/structural/decorator.md
example: examples/structural/decorator/
---

# Decorator Pattern
Decorator structural pattern allows extending the function of an existing object dynamically without altering its internals.

Decorators provide a flexible method to extend functionality of objects.

## Implementation
`LogDecorate` decorates a function with the signature `func(int) int` that 
manipulates integers and adds input/output logging capabilities.

```go
type Object func(int) int

func LogDecorate(fn Object) Object {
	return func(n int) int {
		log.Println("Starting the execution with the integer", n)

		result := fn(n)

		log.Println("Execution is completed with the result", result)

        return result
	}
}
```

## Usage
```go
func Double(n int) int {
    return n * 2
}

f := LogDecorate(Double)

f(5)
// Starting execution with the integer 5
// Execution is completed with the result 10
```

## Rules of Thumb
- Unlike Adapter pattern, the object to be decorated is obtained by **injection**.
- Decorators should not alter the interface of an object.
