# N-Barrier

A reusable synchronisation barrier: all N goroutines block at `Wait` until the
last one arrives, then all proceed together. The barrier re-arms automatically
for the next round.

## Use case

Multi-phase parallel algorithms where phase K must finish completely before
phase K+1 begins (numerical solvers, game simulations, staged pipelines).

## API

```go
b := nbarrier.NewBarrier(n) // arm for n participants
b.Wait()                    // block until all n have called Wait
```

## Run

```bash
go test -race -v ./concurrency/n-barrier/
```

## Trade-offs

A `sync.WaitGroup` covers the single-shot case with less ceremony.  Use this
cyclic barrier when the same group of goroutines repeats many rounds.
