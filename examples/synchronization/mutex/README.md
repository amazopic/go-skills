# Mutex

Two goroutine-safe data structures demonstrating `sync.Mutex` and
`sync.RWMutex`:

- **SafeMap** — `sync.Mutex` protecting a `map[string]int`. All reads and
  writes are exclusive. Use when writes are as frequent as reads.
- **ReadHeavyCache** — `sync.RWMutex` protecting a `map[string]string`. Many
  goroutines may call `Lookup` concurrently; `Put` is exclusive. Use when reads
  dominate.

## Rules of thumb

- `defer mu.Unlock()` immediately after `mu.Lock()` — never forget to unlock.
- Never copy a struct containing a mutex — pass by pointer.
- Keep critical sections short: copy what you need, unlock, then compute.

## Run

```bash
go test -race -v ./synchronization/mutex/
```
