---
name: stability-fail-fast
description: Use when you want to reject work at the earliest detectable failure point — invalid config at startup, bad request at ingress, unready dependency at handler boundary — rather than allowing bad state to propagate.
category: stability
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/stability/fail-fast/
---

# Fail-Fast

## Intent

Surface errors at the earliest detectable point in the system's lifecycle:
at startup (configuration), at request ingress (input validation), or at the
handler boundary (readiness probe). A system that fails fast gives operators
and callers an immediate, precise signal rather than a slow, confusing cascade
of downstream errors.

## Context

The opposite default is retry-with-backoff: the system optimistically attempts
work and handles failure later. Fail-Fast is appropriate when the system can
determine locally, without attempting the work, that success is impossible —
a required field is missing, a dependency is down, or a precondition is
obviously violated. Retrying in that state wastes resources and delays the
useful error signal.

Fail-Fast appears in three distinct phases:

1. **Startup**: Validate all configuration and exit immediately if anything is
   wrong. Operators see the full list of problems, not one per restart.
2. **Ingress**: Validate every inbound request at the handler boundary before
   any downstream work begins. Reject early with a typed error.
3. **Readiness**: Register probes (DB ping, dependency health check) that
   re-validate the system's ability to serve. Reject requests when probes fail.

## Implementation in Go

```go
// Startup: collect all config errors before returning.
func (c Config) Validate() error {
    var errs []string
    if c.Addr == "" {
        errs = append(errs, "Addr: must not be empty")
    }
    if c.MaxConns <= 0 {
        errs = append(errs, "MaxConns: must be > 0")
    }
    if len(errs) > 0 {
        return errors.New(strings.Join(errs, "; "))
    }
    return nil
}

// Ingress: typed error so callers can errors.As to HTTP 400.
type ErrBadRequest struct{ Reason string }
func (e *ErrBadRequest) Error() string { return "bad request: " + e.Reason }

func (r Request) Validate() error {
    if r.UserID == "" {
        return &ErrBadRequest{"UserID must not be empty"}
    }
    return nil
}

// Readiness: atomic bool, safe under concurrent HTTP requests.
type Handler struct {
    ready  atomic.Bool
    probes []func(context.Context) error
}

func (h *Handler) Handle(ctx context.Context, fn func(context.Context) error) error {
    if !h.ready.Load() {
        return ErrNotReady
    }
    for _, p := range h.probes {
        if err := p(ctx); err != nil {
            return fmt.Errorf("%w: probe: %w", ErrNotReady, err)
        }
    }
    return fn(ctx)
}
```

## When to use

- Config-heavy services where a mis-typed env var should crash the process, not degrade silently.
- HTTP/gRPC handlers where malformed requests are common and expensive to propagate.
- Services with hard external dependencies that must be healthy before serving traffic.
- Kubernetes deployments — pair with a readiness probe endpoint so the load balancer never routes to an unready pod.

## When NOT to use

- When the error is transient and retrying after a backoff is likely to succeed (network blip, lock contention).
- When partial success is acceptable — e.g. a batch job that should continue processing other items despite one failure.
- When the "ready" check itself is expensive and would create a hot-loop under high request rates (cache the probe result with a short TTL instead).

## Gotchas

- **Collecting vs. returning the first error.** Return all config errors at once; operators should not have to restart the process three times to find three problems.
- **Probe caching vs. freshness.** Probes run on every request if uncached. Under high load this can become the bottleneck. Cache with a short TTL (1–5s) and return the cached result instead.
- **Race on `ready` flag.** Use `atomic.Bool` (Go 1.19+), not a bare `bool` guarded by no lock. Multiple goroutines read `ready` on every request.
- **Forgetting the startup check exits non-zero.** `os.Exit(1)` (or `log.Fatal`) after a failed `Validate()` is intentional — a process that silently continues with bad config is worse than one that crashes clearly.
- **Conflating fail-fast with panic.** Reserve `panic` for truly unrecoverable programmer errors (nil dereference, invariant violation). Fail-fast for operator-fixable problems uses `error` returns.

## See also

- `skills/stability/deadline.md`
- `skills/stability/circuit-breaker.md`
- `skills/stability/handshaking.md`
- `examples/stability/fail-fast/`
