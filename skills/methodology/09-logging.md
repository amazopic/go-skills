---
name: go-logging
description: Use when wiring loggers: simple coloured logger for bootstrap/connect, slog-based structured logger for services and handlers; log every error exactly once.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§10)
related:
  - skills/methodology/00-canonical-full.md
---

## 10. Logging

### Two tiers

1. **Simple logger** (`pkg/common/logger`) — for bootstrap, fatal errors, and low-level debugging. Coloured, built on top of `log.Logger`.

```go
logger.Init(serviceName)
logger.Info.Printf("started on :%d", port)
logger.Error.Println("connection failed: " + err.Error())
```

2. **Structured logger** (`pkg/common/slogger`) — for everything inside services and handlers. A wrapper over `log/slog` with JSON output, source location, and a shared `app_info` group.

```go
logger := slogger.New(serviceName)
logger.Info("event accepted", "event_id", id, "ip", ip)
logger.Error("xadd failed", "err", err, "stream", streamName)
```

### Rules

- The **structured** logger is injected through DI into every service and handler.
- The **simple** logger is used only in `cmd/.../app/*.go` and `pkg/common/*/connect.go`.
- Never log the same error twice — whoever decides not to propagate it is responsible for logging.
