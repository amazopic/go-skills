---
name: concurrency-n-barrier
description: Use when N goroutines must all reach a rendezvous point before any of them proceeds — e.g., parallel test phases, simulation rounds, multi-stage pipelines.
category: concurrency
go-version-min: "1.22"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
example: examples/concurrency/n-barrier/
---

# N-Barrier

## Intent

An N-Barrier is a synchronisation point at which exactly N goroutines must arrive before any of them is allowed to continue. It models the concept of a "rendezvous": every participant waits until the last one checks in, then all are released together. This is the concurrency analogue of a starting pistol — nobody runs until everyone is on the blocks.

## Context

The pattern appears whenever a parallel workload is split into discrete phases and correctness requires that phase K is fully complete before phase K+1 begins. Classic examples: parallel numerical solvers (all goroutines finish one Jacobi sweep before the next), distributed simulation (all agents step forward one tick in lock-step), integration-test harnesses (all workers reach a known state before assertions fire). The alternative — ad-hoc `time.Sleep` or hoping that goroutines finish in order — is both fragile and racy.

Go's standard library does not ship a reusable cyclic barrier, but one can be built trivially from channels or `sync.WaitGroup`. When you only need a single rendezvous, `WaitGroup` is idiomatic. When you need repeated cycles (re-arming), a channel-based barrier is cleaner.

## Implementation in Go

**Single-shot** — `sync.WaitGroup` suffices:

```go
var wg sync.WaitGroup
wg.Add(n)
for range n {
    go func() {
        defer wg.Done()
        doPhaseWork()
    }()
}
wg.Wait() // barrier point — all n goroutines finished
```

**Cyclic (re-arming)** — replace the WaitGroup with a generation-based barrier. Each generation tracks a countdown; when it hits zero a new channel is closed to broadcast "go", and a fresh generation begins.

```go
type Barrier struct {
    n       int
    mu      sync.Mutex
    waiting int
    done    chan struct{}
}

func NewBarrier(n int) *Barrier {
    return &Barrier{n: n, done: make(chan struct{})}
}

// Wait blocks until n goroutines have called Wait, then all return.
func (b *Barrier) Wait() {
    b.mu.Lock()
    b.waiting++
    if b.waiting == b.n {
        // Last arrival: open the gate and arm a new one.
        close(b.done)
        b.done = make(chan struct{})
        b.waiting = 0
        b.mu.Unlock()
        return
    }
    gate := b.done
    b.mu.Unlock()
    <-gate
}
```

The `close(ch)` broadcast is idiomatic Go: closing a channel unblocks all receivers simultaneously, which is exactly what a barrier needs.

## When to use

- Multi-phase parallel algorithms (numerical solvers, game simulations) where each phase depends on all results from the previous.
- Integration tests that spin up N workers and need a "steady state reached" guarantee before asserting.
- Pipeline stages where the next stage must not start until all producers of the current stage have finished.
- Parallel benchmarks where you want all goroutines to start at the same instant (use barrier before work begins).

## When NOT to use

- When you only need to wait for goroutines to finish without re-use — `sync.WaitGroup` already is a single-shot barrier; no custom type needed.
- When N is not known at construction time or varies dynamically — prefer a fan-in channel pattern or `errgroup`.
- When participants can proceed as soon as their own prerequisite is satisfied rather than everyone's — that is a pipeline, not a barrier.
- High-frequency tight loops: the mutex inside the cyclic barrier adds contention; consider lock-free designs or restructuring the algorithm.

## Gotchas

- **Goroutine leak on panic**: if any participant panics before calling `Wait`, the remaining goroutines block forever. Always `defer` calls that release the barrier or use `errgroup` for panic-safe coordination.
- **Wrong N**: passing a count that never matches the number of actual goroutines deadlocks silently. Prefer computing N from `len(workers)` rather than a magic constant.
- **Re-arming race**: in the cyclic version, reading `b.done` and then blocking on it must be atomic with respect to the counter update — the mutex covers both; do not cache `b.done` outside the lock then release the lock before blocking.
- **Context cancellation**: the bare `<-gate` in `Wait` is not cancellable. In production, select on both `gate` and `ctx.Done()` and propagate `ctx.Err()`.

## See also

- skills/concurrency/producer-consumer.md
- skills/synchronization/condition-variable.md
- examples/concurrency/n-barrier/
