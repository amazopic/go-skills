# Monitor

The Monitor pattern bundles shared state with its synchronisation primitives
into one type. Callers use public methods; the mutex is invisible.

Two examples:

- **ResourcePool** — fixed-capacity token pool using `sync.Mutex` + `sync.Cond`.
  `Acquire` blocks when empty; `Release` wakes a waiter.
- **AtomicRegistry** — map with atomic `GetOrSet` using `sync.RWMutex`.
  The compound check-and-set is safe because the lock spans both steps.

## Design rules

1. Mutex is unexported — callers cannot lock it directly.
2. Every exported method that touches shared state must lock before reading
   or writing and unlock (via `defer`) before returning.
3. Unexported helper methods that assume the lock is already held avoid
   double-lock deadlocks.

## Run

```bash
go test -race -v ./synchronization/monitor/
```
