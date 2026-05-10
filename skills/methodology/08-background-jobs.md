---
name: go-background-jobs
description: Use when scheduling background tasks with robfig/cron/v3 (second precision); register via di.go after dependency graph.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§9)
related:
  - skills/methodology/00-canonical-full.md
---

## 9. Background jobs (cron)

`robfig/cron/v3` with second-level precision.

```go
// internal/<service>/crons/cron.go
type Cron struct {
    Cache   *cache.Cache
    Service *domain.Service
    Logger  *slogger.Logger
}

func NewCron(c *cache.Cache, s *domain.Service, l *slogger.Logger) *Cron {
    return &Cron{Cache: c, Service: s, Logger: l}
}

// internal/<service>/crons/sync.go
func (c *Cron) Sync(ctx context.Context) error {
    items := c.Cache.GetCount(1000)
    for _, it := range items {
        if err := c.Service.Process(ctx, it); err != nil {
            return fmt.Errorf("process: %w", err)
        }
    }
    return nil
}
```

Register jobs in `di.go` via `app.cron.AddFunc("...", func)`. Start the scheduler with `app.cron.Start()` after `InitializeDependencies`.
