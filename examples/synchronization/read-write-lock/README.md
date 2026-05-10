# Read-Write Lock

`ConfigStore` is a goroutine-safe key-value store backed by `sync.RWMutex`:

- `Get` / `Snapshot` / `Len` — shared read lock; many callers run concurrently.
- `Set` / `Delete` — exclusive write lock; blocks all readers and other writers.

## When it wins

Read-to-write ratio ≥ 10:1. For write-heavy or contention-free workloads,
`sync.Mutex` is often faster due to lower overhead.

## Gotchas

- Never upgrade RLock → Lock in the same goroutine (deadlock).
- `sync.RWMutex` must not be copied — pass the struct by pointer.
- Keep RLock sections short; release before any I/O.

## Run

```bash
go test -race -v ./synchronization/read-write-lock/
```
