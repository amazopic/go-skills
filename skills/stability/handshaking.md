---
name: stability-handshaking
description: Use when a client should negotiate with a server about how much work it can accept before submitting that work, adapting send rate to the server's current capacity advertisement.
category: stability
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/stability/handshaking/
---

# Handshaking

## Intent

Before submitting work, the client asks the server "how much can you accept
right now?" and adapts accordingly. The server advertises its available
concurrency budget (and optionally a lease TTL for caching). If the budget
is zero, the client refuses to send — preventing wasted RPCs and reducing
pressure on an already-struggling server.

## Context

Without handshaking, clients discover overload reactively: they send, the
server queues or drops, the client retries, the problem compounds. Handshaking
shifts the conversation earlier: the server signals capacity proactively and
the client adapts before generating load. This is the same mechanism HTTP/2
uses with `SETTINGS_MAX_CONCURRENT_STREAMS` — the server pushes a budget and
the client enforces it locally.

Handshaking complements circuit breakers (which open after observed failures)
and bulkheads (which cap client-side concurrency). It adds a server-driven
signal to the client's admission control, replacing guesswork with negotiation.

## Implementation in Go

The server exposes a capacity endpoint (a plain function in the model below;
an HTTP handler in production). The client fetches and caches it for the
lease TTL, then checks before each send:

```go
type CapacityResponse struct {
    Slots    int           // available concurrent request budget
    LeaseTTL time.Duration // how long to trust this response
}

// Server: enforce limit server-side too — client check is advisory.
func (s *Server) Capacity() CapacityResponse {
    s.mu.Lock()
    avail := s.limit - s.inFlight
    s.mu.Unlock()
    return CapacityResponse{Slots: avail, LeaseTTL: 500 * time.Millisecond}
}

// Client: cache the advertisement for LeaseTTL, then check before each send.
func (c *Client) Send(ctx context.Context, fn func(context.Context) error) error {
    resp, err := c.capacity(ctx) // re-fetches when lease expires
    if err != nil {
        return fmt.Errorf("capacity check: %w", err)
    }
    if resp.Slots <= 0 {
        return ErrCapacityExceeded
    }
    return fn(ctx)
}
```

The server still enforces the limit on its own mutex — the client check is
advisory, not authoritative. This prevents a misbehaving or misconfigured
client from bypassing the guard.

## When to use

- Clients that fan out to a downstream with unpredictable or bursty capacity.
- Rate-limited APIs where you want to avoid wasting quota on requests that will be rejected.
- Systems that need graceful degradation: when capacity drops to zero, callers get a clean error rather than a timeout storm.
- Any protocol that natively supports capacity advertisements (gRPC flow control, HTTP/2 SETTINGS, AMQP prefetch).

## When NOT to use

- When the downstream has effectively unlimited capacity (well-scaled CDN, object storage).
- When the capacity endpoint itself becomes a bottleneck — mitigate with a longer LeaseTTL or a push-based mechanism.
- When latency of the capacity check exceeds the latency of the actual call.

## Gotchas

- **Server-side enforcement is not optional.** A malicious or buggy client can skip the handshake. The server must enforce its limit independently.
- **Lease TTL too long.** A long TTL means the client sends to a server whose capacity has dropped significantly since the last fetch. Tune TTL to be shorter than the time it takes the server to become saturated.
- **Lease TTL too short.** Every request hits the capacity endpoint — doubles your RPC count. Cache aggressively for stable capacity, less so for volatile.
- **Thundering herd on lease expiry.** If 100 clients expire their lease simultaneously, they all hit the capacity endpoint at once. Add jitter: `LeaseTTL + rand.Duration(0, 50ms)`.
- **Zero-slots ≠ error.** `ErrCapacityExceeded` is not a server error; it is a load-shedding signal. The client should backoff and retry, not circuit-break.

## See also

- `skills/stability/bulkhead.md`
- `skills/stability/circuit-breaker.md`
- `skills/stability/fail-fast.md`
- `examples/stability/handshaking/`
