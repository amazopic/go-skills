---
name: go-storage
description: Use when choosing or wiring storage: PostgreSQL via pgx/sqlx, ClickHouse batched writes via PrepareBatch, Redis Streams for durable queues, in-memory TTL cache with janitor.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§7)
related:
  - skills/methodology/00-canonical-full.md
---

## 7. Storage

### Defaults

- **OLTP** → PostgreSQL via `jackc/pgx` or `jmoiron/sqlx`.
- **Analytics / events** → ClickHouse via `ClickHouse/clickhouse-go/v2`. Always batch via `PrepareBatch`.
- **Cache / streams / queues** → Redis (`go-redis/v9`). Use Streams (`XAdd`/`XReadGroup`) for durable queues.
- **In-memory TTL cache** → in-house generic implementation with a janitor goroutine.

### ClickHouse batch insert

```go
func (s *Service) InsertBatch(ctx context.Context, items []dto.Event) error {
    batch, err := s.ClickHouseConn.PrepareBatch(ctx, "INSERT INTO events")
    if err != nil {
        return fmt.Errorf("prepare batch: %w", err)
    }
    for i, it := range items {
        if err := batch.AppendStruct(&it); err != nil {
            return fmt.Errorf("append %d: %w", i, err)
        }
    }
    if err := batch.Send(); err != nil {
        return fmt.Errorf("send batch: %w", err)
    }
    return nil
}
```

### In-memory TTL cache

API:
- `Set(key, value, ttl)`
- `Add(key, value, ttl)` — fails if the key exists
- `Get(key) (any, bool)`
- `GetCount(n int) []Item` — drain N keys for batch processing
- `Delete(key)`
- `ItemCount() int`
- janitor goroutine — periodically evicts expired entries

Implementation: `map[string]Item` guarded by `sync.RWMutex`; item = `{Object any, Expiration int64}`.
