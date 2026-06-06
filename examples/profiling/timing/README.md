# Timing

Lightweight, stdlib-only latency measurement for hot paths where reaching for
pprof or a full metrics framework is overkill.

The idiomatic one-liner uses a deferred call whose arguments are evaluated at
function entry:

```go
func work() {
	defer timing.Track(time.Now(), "work") // logs "work took 1.2ms"
	// ...
}
```

For aggregating call latencies (and for deterministic tests), use a `Recorder`
driven by an injectable `Clock`:

```go
rec := timing.NewRecorder(nil) // nil -> SystemClock

func handle() {
	defer rec.Start("handle")() // records elapsed under "handle"
	// ...
}

n, err := timing.Measure(rec, "parse", func() (int, error) { return parse() })

stats, _ := rec.Stats("handle") // Count, Sum, Min, Max, Mean
```

## When to use

- Quick "where is the time going?" checks during optimization.
- Per-call latency aggregation without external dependencies.
- Anywhere you want timing under test control: inject a fake `Clock` so
  measurements are exact and never wall-clock-flaky.

## Structure

| File | Purpose |
|---|---|
| `timing.go` | `Track` defer helper, `Recorder` (concurrent-safe), `Clock`/`SystemClock`, generic `Measure` |
| `timing_test.go` | Table-driven, race-safe tests using a deterministic fake clock |

## Key points

- `Recorder` is safe for concurrent use; `Stats` copies samples under the lock.
- `Clock` injection makes elapsed times deterministic — no `time.Sleep` in tests.
- Negative durations are clamped to zero so a misbehaving clock cannot corrupt `Min`.
- `Stats` returns `timing.ErrNoSamples` for an unrecorded label.

## Run

```bash
cd examples && go test -race ./profiling/timing/
```
