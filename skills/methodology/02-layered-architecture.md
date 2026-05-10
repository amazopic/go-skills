---
name: go-layered-architecture
description: Use when designing transport/service/fetcher boundaries and the cross-layer flow of an HTTP request through middleware → handler → service → fetcher → DB.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§3, §17)
related:
  - skills/methodology/00-canonical-full.md
---

## 3. Layers

### 3.1 Transport (handler)

**Responsibility**: parse and validate input, call a service, render the response. No business logic here.

**Root struct** (`handler/http.go`):

```go
type Handler struct {
    port              int
    Logger            *slogger.Logger
    OperationAHandler *opA.Handler
    OperationBHandler *opB.Handler
}
```

**Route registration** (`handler/router.go`):

```go
func (h *Handler) NewRouter(allowedCORS []string, guard *cache.Cache) http.Handler {
    r := chi.NewRouter()

    r.Use(
        middleware.JSON(),
        chiMiddleware.Timeout(20*time.Second),
        middleware.Cors(allowedCORS),
        middleware.RateLimit(guard),
    )

    r.Get("/health", healthHandler)

    r.Route("/v1", func(r chi.Router) {
        r.Post("/operation-a", h.OperationAHandler.Handle)
        r.Get("/operation-b", h.OperationBHandler.Handle)
    })

    return r
}
```

**Per-operation handler** (`handler/<operation>/handler.go`):

```go
type Handler struct {
    Logger  *slogger.Logger
    Service *domainService.Service
}

func NewHandler(logger *slogger.Logger, svc *domainService.Service) *Handler {
    return &Handler{Logger: logger, Service: svc}
}
```

**HTTP method** (`handler/<operation>/<operation>.go`):

```go
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
    var req dto.Request
    if err := httpKit.Decode(r, &req); err != nil {
        render.Render(w, r, common.NewErrorResponse("decode error", http.StatusBadRequest))
        return
    }

    if err := validator.Validate(&req); err != nil {
        render.Render(w, r, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
        return
    }

    if err := h.Service.Do(r.Context(), req.ToCommand()); err != nil {
        h.Logger.Error("operation failed", "err", err)
        render.Render(w, r, common.NewErrorResponse("internal", http.StatusInternalServerError))
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
```

### 3.2 Business logic (service)

**Responsibility**: orchestration — calling fetchers, repositories, external APIs; applying business rules.

**Struct + constructor** (`service/<domain>/service.go`):

```go
type Service struct {
    Logger      *slogger.Logger
    RedisClient *redis.Client
    LocalCache  *cache.Cache
}

func NewService(logger *slogger.Logger, rc *redis.Client, lc *cache.Cache) *Service {
    return &Service{Logger: logger, RedisClient: rc, LocalCache: lc}
}
```

**Domain method** (`service/<domain>/<operation>.go`):

```go
type DoCommand struct {
    ID      string
    Payload string
}

func (s *Service) Do(ctx context.Context, cmd DoCommand) error {
    if err := s.RedisClient.XAdd(ctx, &redis.XAddArgs{
        Stream: streamName,
        Values: map[string]any{"id": cmd.ID, "payload": cmd.Payload},
    }).Err(); err != nil {
        return fmt.Errorf("xadd: %w", err)
    }
    return nil
}
```

**Rules**:
- Constructor `NewService(...)` accepts every dependency. Returns `*Service`.
- Methods take `context.Context` as the first argument when they do I/O.
- Methods accept primitives or `Command` structs — never HTTP types.
- Methods return `error` last.
- No global state.

### 3.3 Fetcher — specialized read layer

When a read path is non-trivial or shared across services, factor it into a Fetcher. This makes reads easy to mock in tests and keeps services lean.

```go
type Fetcher struct {
    Logger      *slogger.Logger
    RedisClient *redis.Client
    LocalCache  *cache.Cache
}

func NewFetcher(logger *slogger.Logger, rc *redis.Client, lc *cache.Cache) *Fetcher {
    return &Fetcher{Logger: logger, RedisClient: rc, LocalCache: lc}
}

func (f *Fetcher) GetSize(ctx context.Context) (int, error) {
    return f.LocalCache.ItemCount(), nil
}
```

### 3.4 DTOs

Flat structs with json tags. No methods, except simple converters such as `ToCommand()` / `FromEntity()`.

```go
type EventDTO struct {
    EventID   string            `json:"event_id"`
    EventName string            `json:"event_name"`
    Payload   map[string]string `json:"payload"`
}
```

If a request arrives base64-encoded or in any other non-standard form, decode it in the handler — never in the service.

---

## 17. Cross-layer flow

```
HTTP request
    ↓
[middleware] → context enriched with request data (IP, request-id)
    ↓
[handler.<operation>.Handle]
    ↓ decode → validate → build Command
[service.<domain>.Method(ctx, cmd)]
    ↓
[fetcher / db client / redis client / external API]
    ↓
result bubbles back up: error or value
    ↓
handler renders the HTTP response
```

Rules:
- A handler **never** calls a DB client or fetcher directly — only services. Exception: pure read-only endpoints may go straight to a fetcher.
- A service **knows nothing about HTTP** — no `http.Request`, no `http.ResponseWriter`.
- A fetcher **knows nothing about business rules** — only reads and shapes data.
