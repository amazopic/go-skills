---
name: profiling-timing
description: Timing functions — measure latency with `defer fn(time.Now())` and a closure that logs the elapsed duration. Use to add lightweight per-call timing.
category: profiling
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/profiling/timing.md
example: examples/profiling/timing/
---

# Timing Functions

When optimizing code, sometimes a quick and dirty time measurement is required
as opposed to utilizing profiler tools/frameworks to validate assumptions.

Time measurements can be performed by utilizing `time` package and `defer` statements.

## Implementation

```go
package profile

import (
    "time"
    "log"
)

func Duration(invocation time.Time, name string) {
    elapsed := time.Since(invocation)

    log.Printf("%s lasted %s", name, elapsed)
}
```

## Usage

```go
func BigIntFactorial(x big.Int) *big.Int {
    // Arguments to a defer statement is immediately evaluated and stored.
    // The deferred function receives the pre-evaluated values when its invoked.
    defer profile.Duration(time.Now(), "IntFactorial")

    y := big.NewInt(1)
    for one := big.NewInt(1); x.Sign() > 0; x.Sub(x, one) {
        y.Mul(y, x)
    }

    return x.Set(y)
}
```
