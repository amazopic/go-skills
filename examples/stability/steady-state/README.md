# Steady-State Example

Bounds resource accumulation so the system reaches stable equilibrium under
sustained load. Two complementary data structures: a TTL+capacity-bounded
`Cache` with a background janitor goroutine, and a fixed-size `RingBuffer`
that overwrites the oldest entry on overflow.

## Structure

| File | Purpose |
|---|---|
| `steadystate.go` | `Cache` (TTL eviction + capacity cap + janitor), `RingBuffer[T]` (generic, fixed-size FIFO) |
| `steadystate_test.go` | Tests for eviction, capacity shedding, janitor lifecycle, ring overflow, and concurrent access |

## Run

```bash
go test -race ./stability/steady-state/
```

## Key points

- The janitor goroutine exits cleanly via `context.Done()` or `Stop()` — no goroutine leak.
- `Cache.Set` makes a best-effort eviction attempt when at capacity before shedding the write — this prevents a thundering-herd of "cache full" errors when entries expire naturally.
- `RingBuffer` uses a single mutex; the modular arithmetic `(head - size + cap) % cap` is the idiom for ordered-oldest-first slicing without copying twice.
- Both structures are bounded at construction time — no dynamic growth, predictable memory footprint.
