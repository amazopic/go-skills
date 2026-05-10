# Producer Consumer

N producers generate `Item` values into a bounded buffered channel; M consumers
read and transform them into `Result` values. The bounded buffer provides
backpressure: when consumers are slow, producers block rather than using
unbounded memory.

## Key design points

- Producers own the close of `jobs`. A sentinel goroutine calls `close(jobs)`
  after all producers finish via `sync.WaitGroup`.
- Consumers use `range jobs` — exits automatically when the channel is closed
  and drained.
- Context cancellation is respected by both producers and consumers.

## Run

```bash
go test -race -v ./concurrency/producer-consumer/
```

## Tuning

`bufSize` controls burst absorption. Too small → producer stalls; too large →
memory pressure and slow cancellation. Profile under realistic load.
