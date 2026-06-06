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
package timing

import (
    "log"
    "time"
)

// Track logs the time elapsed since start under name. Call it deferred so that
// start (the argument) is evaluated at function entry.
func Track(start time.Time, name string) {
    log.Printf("%s took %s", name, time.Since(start))
}
```

## Usage

```go
func BigIntFactorial(x big.Int) *big.Int {
    // Arguments to a defer statement are evaluated immediately and stored.
    // The deferred function receives the pre-evaluated values when it is invoked.
    defer timing.Track(time.Now(), "BigIntFactorial")

    y := big.NewInt(1)
    for one := big.NewInt(1); x.Sign() > 0; x.Sub(x, one) {
        y.Mul(y, x)
    }

    return x.Set(y)
}
```

## Aggregating & testing

The bare `Track` helper logs wall-clock elapsed time, which is fine for a quick
look but impossible to assert on in a test. For aggregating per-call latencies —
and for deterministic tests — use a `Recorder` driven by an injectable `Clock`:

```go
rec := timing.NewRecorder(nil) // nil -> SystemClock

func handle() {
    defer rec.Start("handle")() // records elapsed under "handle"
    // ...
}

stats, _ := rec.Stats("handle") // Count, Sum, Min, Max, Mean
```

Injecting a fake `Clock` in tests yields exact, non-flaky durations with no
`time.Sleep`. See `examples/profiling/timing/` for the full, race-tested code.
