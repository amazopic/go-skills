# Circuit Breaker Example

`Breaker[R]` wraps a call to an unreliable dependency (external API, database,
downstream service). It counts consecutive failures and, once a threshold is
crossed, **opens** the circuit so further calls fail fast with `ErrOpenCircuit`
instead of piling onto a struggling dependency. After a cooldown it goes
**half-open** and admits one trial call; success closes the circuit, failure
re-opens it with an exponentially longer cooldown.

## When to use

- Calling a flaky downstream where retries make things worse (cascading load).
- You want to fail fast and shed load while a dependency recovers.
- You need a bounded, self-healing recovery probe rather than blind retries.

## Structure

| File | Purpose |
|---|---|
| `circuit_breaker.go` | `Breaker[R]`, `State`, `ErrOpenCircuit`, closed/open/half-open transitions with exponential backoff |
| `circuit_breaker_test.go` | Table-driven, race-safe tests using an injected fake clock (no `time.Sleep`) plus a concurrent failure-then-recovery stress test |

## API

```go
b := circuitbreaker.New[string](
    3,             // trip after 3 consecutive failures
    time.Second,   // base cooldown (doubles per re-trip)
)

res, err := b.Do(ctx, func(ctx context.Context) (string, error) {
    return callDownstream(ctx)
})
switch {
case errors.Is(err, circuitbreaker.ErrOpenCircuit):
    // failed fast — circuit is open, dependency is presumed down
case err != nil:
    // the downstream returned an error (wrapped/returned as-is)
default:
    // success: res is valid
}
```

## Key points

- **Generic over the result type** (`Breaker[R]`): no `any` boxing at call sites.
- **Context-first**: a cancelled `ctx` short-circuits before touching breaker state.
- **Injectable clock** (`WithClock`) keeps tests deterministic — no real time.
- All state transitions happen under one mutex, so the breaker is safe for
  concurrent use and clean under `-race`.

## Run

```bash
cd examples && go test -race ./stability/circuit-breaker/
```
