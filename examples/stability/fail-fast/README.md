# Fail-Fast Example

Surfaces errors at the earliest possible point: startup (config), ingress
(request validation), and dependency check (readiness probe). Avoids wasted
work propagating bad state into the call graph.

## Structure

| File | Purpose |
|---|---|
| `failfast.go` | `Config.Validate`, `Handler` with readiness probes, `Request.Validate` |
| `failfast_test.go` | Table-driven tests for all three fail-fast layers + concurrent race test |

## Run

```bash
go test -race ./stability/fail-fast/
```

## Key points

- `Config.Validate` collects all field errors before returning — operators get the full picture rather than fixing one field at a time.
- `Handler` uses `atomic.Bool` for `SetReady`/`IsReady` — safe under concurrent HTTP requests without a mutex.
- Probes are checked even when `ready=true`, guarding against transient dependency outages after startup.
- `Request.Validate` returns typed `*ErrBadRequest` so callers can `errors.As` and map to HTTP 400 without string matching.
