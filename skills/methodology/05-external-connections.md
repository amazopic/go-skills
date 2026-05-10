---
name: go-external-connections
description: Use when writing a Connect() factory for redis/clickhouse/postgres: retry loop, timeouts, ping-after-dial verification.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§6)
related:
  - skills/methodology/00-canonical-full.md
---

## 6. External connections

**All connectors live in `pkg/common/<subsystem>/connect.go`** with mandatory retry logic and timeouts.

```go
const (
    MaxAttempts  = 5
    RetryDelay   = 2 * time.Second
    DialTimeout  = 5 * time.Second
    ReadTimeout  = 3 * time.Second
    WriteTimeout = 3 * time.Second
)

func Connect(cfg Config) (*redis.Client, error) {
    var client *redis.Client
    var err error
    ctx := context.Background()

    for attempt := 1; attempt <= MaxAttempts; attempt++ {
        client = redis.NewClient(&redis.Options{
            Addr:         cfg.Host + ":" + cfg.Port,
            DB:           cfg.DB,
            DialTimeout:  DialTimeout,
            ReadTimeout:  ReadTimeout,
            WriteTimeout: WriteTimeout,
        })

        if _, err = client.Ping(ctx).Result(); err == nil {
            logger.Info.Printf("connected to redis after %d attempts", attempt)
            return client, nil
        }
        time.Sleep(RetryDelay)
    }
    return nil, fmt.Errorf("connect redis after %d attempts: %w", MaxAttempts, err)
}
```
