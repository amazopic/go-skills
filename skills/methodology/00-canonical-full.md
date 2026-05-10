---
name: go-service-architecture-canonical
description: Canonical full methodology for building backend services in Go — directory layout, layer separation, DI, configuration, transport, storage, jobs, logging, build, deploy. Use when reviewing or designing an entire Go service end-to-end. For chapter-scoped skills see 01–13.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md
---

# Go Service Architecture Methodology

A reference for building backend services in Go: directory layout, layer separation, dependency injection, configuration, transport, storage, background jobs, logging, and build/deploy.

This document is anonymized — it captures only the methodology, with no project-specific names. Use it as a template for new services and as a baseline for code reviews.

---

## 1. Principles

1. **Strict layer separation**: transport → business logic → data access. Upper layers depend only on lower layers.
2. **Explicit dependency injection** via constructors. No global singletons, service locators, or reflection-based wiring.
3. **Minimal frameworks**: standard library plus small focused libraries (router, env parser, DB drivers). No heavyweight DI/ORM frameworks.
4. **One service = one binary = one folder under `cmd/`**. A monorepo can host several services that share `pkg/common`.
5. **`context.Context` is propagated everywhere** — first argument of any method that performs I/O.
6. **Errors are explicit**: `error` is the last return value; wrap with `fmt.Errorf("...: %w", err)`.
7. **Graceful shutdown is mandatory** for every service.

---

## 2. Directory Layout

```
.
├── cmd/                                  # Entry points (one folder per binary)
│   └── <service-name>/
│       ├── main.go                       # Bootstrap, http.Server, signal handling
│       └── app/
│           ├── app.go                    # App struct, resource init, Close()
│           └── di.go                     # Dependency graph assembled in stages
│
├── internal/                             # Private code (not importable from outside the module)
│   └── <service-name>/
│       ├── handler/                      # Transport layer (HTTP)
│       │   ├── http.go                   # Root Handler struct
│       │   ├── router.go                 # Route + middleware registration
│       │   ├── middleware/               # Service-specific middleware
│       │   ├── dto/                      # Request/response models
│       │   └── <operation>/
│       │       ├── handler.go            # Per-operation constructor + deps
│       │       └── <operation>.go        # HTTP method implementation
│       │
│       ├── service/                      # Business logic
│       │   └── <domain>/
│       │       ├── service.go            # Service struct + constructor
│       │       └── <operation>.go        # Domain methods
│       │
│       ├── fetcher/                      # Read-side data accessors
│       │   └── <source>/
│       │       ├── fetcher.go
│       │       └── <operation>.go
│       │
│       ├── crons/                        # Scheduled background jobs
│       │   ├── cron.go                   # Cron struct + constructor
│       │   └── <task>.go                 # Job implementation
│       │
│       ├── validator/                    # Input validation
│       └── docs/                         # Generated Swagger output
│
├── pkg/                                  # Reusable utilities (importable)
│   └── common/
│       ├── config/                       # Config loading (ENV + .env)
│       ├── logger/                       # Simple coloured logger
│       ├── slogger/                      # Structured logger (slog wrapper)
│       ├── database/<driver>/            # DB connectors (postgres, clickhouse, ...)
│       ├── cache/<driver>/               # Cache connectors (redis, local, ...)
│       ├── lib/cache/                    # Generic in-memory TTL cache
│       ├── httpKit/                      # HTTP helpers (decoders, renderers)
│       ├── middleware/                   # Shared middleware + context-key constants
│       ├── tracer/                       # Tracing helpers
│       └── error.go                      # Unified ErrorResponse
│
├── go.mod
├── go.sum
├── makefile                              # Build/test/lint targets
├── .gitlab-ci.yml                        # CI/CD pipeline
└── test.sh                               # Optional smoke check
```

### Placement rules

- **`cmd/`** holds bootstrap only. No business logic.
- **`internal/`** is the domain layer. Go enforces that it is not importable from outside the module.
- **`pkg/common/`** contains only utilities that are reused by ≥ 2 services and carry no domain knowledge.
- An operation folder (`handler/<operation>/`, `service/<domain>/`) contains **at least two files**: `handler.go`/`service.go` (struct + constructor) and one or more files with the actual methods. This keeps reads small and minimizes diff churn.

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

## 4. Entry point and bootstrap

### 4.1 `cmd/<service>/main.go`

```go
func main() {
    application := app.Init()

    server := http.Server{
        Addr:              fmt.Sprintf(":%d", application.Cfg.HTTPPort),
        Handler:           application.HttpHandler,
        ReadTimeout:       10 * time.Second,
        ReadHeaderTimeout: 5 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
        MaxHeaderBytes:    1 << 20,
    }

    go func() {
        if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            logger.Error.Fatalln(err)
        }
    }()

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    <-sigs

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _ = server.Shutdown(ctx)
    application.Close()
}
```

### 4.2 `cmd/<service>/app/app.go` — struct and initialization

```go
type App struct {
    Cfg            config.Config
    HttpHandler    http.Handler
    Logger         *slogger.Logger
    ctx            context.Context

    db             *sqlx.DB
    redisClient    *redis.Client
    localCache     *cache.Cache
    clickhouseConn clickhouse.Conn
    cron           *cron.Cron
}

func Init() *App {
    app := &App{ctx: context.Background()}

    logger.Init(serviceName)
    app.Logger = slogger.New(serviceName)

    cfg, err := config.LoadConfig(serviceName)
    if err != nil {
        logger.Error.Fatalln(err)
    }
    app.Cfg = cfg

    app.redisClient, _ = cacheRedis.Connect(cfg.Redis)
    app.clickhouseConn, _ = clickhouseDB.Connect(cfg.Clickhouse)
    app.localCache = local.Connect(cfg.LocalCache)

    app.cron = cron.New(cron.WithLocation(time.UTC), cron.WithSeconds())

    if err := app.InitializeDependencies(); err != nil {
        logger.Error.Fatalln(err)
    }

    return app
}

func (app *App) Close() {
    if app.db != nil { _ = app.db.Close() }
    if app.redisClient != nil { _ = app.redisClient.Close() }
    if app.clickhouseConn != nil { _ = app.clickhouseConn.Close() }
}
```

### 4.3 `cmd/<service>/app/di.go` — dependency graph

DI is **manual and staged**. No containers (wire/fx). Stages:

1. **Services** (business logic, fetchers).
2. **Handlers** (depend on services).
3. **Cron jobs** (depend on services).
4. **HTTP router** (depends on handlers).

```go
type Dependencies struct {
    // services
    DomainService *domain.Service
    Fetcher       *src.Fetcher

    // handlers
    OperationAHandler *opA.Handler
    OperationBHandler *opB.Handler

    // crons
    SyncCron func(context.Context) error
}

func (app *App) InitializeDependencies() error {
    deps := app.initServices()
    deps = app.initHandlers(deps)
    if err := app.initCronJobs(deps); err != nil {
        return err
    }
    app.HttpHandler = app.initHTTPHandler(deps)
    app.cron.Start()
    return nil
}

func (app *App) initServices() *Dependencies {
    return &Dependencies{
        DomainService: domain.NewService(app.Logger, app.redisClient, app.localCache),
        Fetcher:       src.NewFetcher(app.Logger, app.localCache),
    }
}

func (app *App) initHandlers(d *Dependencies) *Dependencies {
    d.OperationAHandler = opA.NewHandler(app.Logger, d.DomainService)
    d.OperationBHandler = opB.NewHandler(app.Logger, d.Fetcher)
    return d
}

func (app *App) initCronJobs(d *Dependencies) error {
    c := crons.NewCron(app.localCache, d.DomainService, app.Logger)
    d.SyncCron = c.Sync

    _, err := app.cron.AddFunc("*/10 * * * * *", func() {
        if err := d.SyncCron(app.ctx); err != nil {
            logger.Error.Println("sync cron failed: " + err.Error())
        }
    })
    return err
}

func (app *App) initHTTPHandler(d *Dependencies) http.Handler {
    h := handler.NewHandler(app.Logger, app.Cfg.HTTPPort)
    h.OperationAHandler = d.OperationAHandler
    h.OperationBHandler = d.OperationBHandler
    return h.NewRouter(app.Cfg.AllowedCORSOrigins(), app.localCache)
}
```

---

## 5. Configuration

### Approach

- **Source of truth**: environment variables.
- In dev mode, fall back to `.env` via `godotenv`.
- Parsing: `caarlos0/env` with struct tags (`env:"..."`, `envDefault:"..."`).
- Each subsystem owns its `Config` struct under `pkg/common/<subsystem>/`. The root `Config` aggregates them.

### Root config

```go
type Config struct {
    HTTPPort           int    `env:"HTTP_PORT"`
    DevMode            string `env:"DEV_MODE"`
    AllowedCorsOrigins string `env:"ALLOWED_CORS_ORIGINS"`

    Postgres   postgres.Config
    Clickhouse clickhouse.Config
    Redis      redis.Config
    LocalCache local.Config
}

func LoadConfig(serviceID string) (Config, error) {
    envPath := filepath.Join(".", ".env")
    if _, err := os.Stat(envPath); errors.Is(err, os.ErrNotExist) {
        envPath = filepath.Join("/etc/<app>/", serviceID, ".env")
    }
    _ = godotenv.Load(envPath)

    var cfg Config
    if err := env.Parse(&cfg); err != nil {
        return Config{}, fmt.Errorf("parse env: %w", err)
    }
    common.DevMode = cfg.DevMode
    return cfg, nil
}
```

### Subsystem configs

```go
// pkg/common/cache/redis/config.go
type Config struct {
    Host string `env:"REDIS_HOST" envDefault:"127.0.0.1"`
    Port string `env:"REDIS_PORT" envDefault:"6379"`
    DB   int    `env:"REDIS_DB"   envDefault:"0"`
}

// pkg/common/database/clickhouse/config.go
type Config struct {
    Host     string `env:"CLICKHOUSE_HOST"     envDefault:"localhost"`
    Port     int    `env:"CLICKHOUSE_PORT"     envDefault:"9000"`
    User     string `env:"CLICKHOUSE_USER"     envDefault:"default"`
    Password string `env:"CLICKHOUSE_PASSWORD" envDefault:""`
    DBName   string `env:"CLICKHOUSE_DB"       envDefault:"default"`
}
```

### Derived methods

For composite values (lists, parsed strings), use methods on the config with lazy initialization via `sync.Once`:

```go
func (c Config) AllowedCORSOrigins() []string {
    c.corsOnce.Do(func() {
        c.corsCached = strings.Split(c.AllowedCorsOrigins, ";")
    })
    return c.corsCached
}
```

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

---

## 8. HTTP transport

### Router

`go-chi/chi/v5` — lightweight, idiomatic, fully compatible with `net/http`.

### Standard middleware stack

```go
r.Use(
    middleware.JSON(),                       // Content-Type
    chiMiddleware.Timeout(20*time.Second),   // global timeout
    middleware.Cors(allowedCORSOrigins),     // CORS
    middleware.RateLimit(guard),             // per-IP rate limiting
    // tracing.Middleware(),                 // optional
)
```

### Custom middleware example

```go
func JSON() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            next.ServeHTTP(w, r)
        })
    }
}
```

### Rate limiting (in-memory cache)

Approach: per-IP counters in 1s / 1m / 1h windows. When a window exceeds its limit, write a block key with a TTL.

```go
func RateLimit(c *cache.Cache) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := getClientIP(r)
            if _, blocked := c.Get(keyBlocked + ip); blocked {
                http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
                return
            }
            if incr(c, keySec+ip, time.Second) > maxPerSecond {
                c.Set(keyBlocked+ip, true, blockTTL)
                http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

### Passing data via context

Key constants live in `pkg/common/middleware/context_key/`:

```go
type contextKey string
const ClientIP contextKey = "client_ip"

// in middleware
ctx := context.WithValue(r.Context(), context_key.ClientIP, ip)
next.ServeHTTP(w, r.WithContext(ctx))

// in handler
ip, _ := r.Context().Value(context_key.ClientIP).(string)
```

### Unified error response

```go
// pkg/common/error.go
type ErrorResponse struct {
    StatusCode int    `json:"-"`
    StatusText string `json:"status"`
    ErrorText  string `json:"error"`
}

func NewErrorResponse(msg string, code int) render.Renderer {
    return &ErrorResponse{
        StatusCode: code,
        StatusText: http.StatusText(code),
        ErrorText:  msg,
    }
}

func (e *ErrorResponse) Render(_ http.ResponseWriter, r *http.Request) error {
    render.Status(r, e.StatusCode)
    return nil
}
```

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

---

## 11. Validation

A dedicated `internal/<service>/validator/` package. Pure functions returning `error`. Tests use table-driven style and `Example*` for documentation.

```go
func ValidateEvent(e *dto.EventDTO) error {
    if e.EventID == "" { return errors.New("event_id is required") }
    if len(e.EventName) > 256 { return errors.New("event_name too long") }
    return nil
}
```

The handler runs validation **before** invoking the service. The service trusts its input.

---

## 12. Error handling

- Return: `error` is always the last value.
- Wrap: `fmt.Errorf("context: %w", err)` with a meaningful prefix.
- Log: once, close to where the error stops being propagated.
- HTTP: map through `common.NewErrorResponse(msg, code)` + `render.Render(w, r, ...)`.
- Preserve `context.Canceled` / `context.DeadlineExceeded` — check with `errors.Is` where it matters.

---

## 13. Testing

### makefile targets

```makefile
test:
	go test ./... -v

test-short:
	go test ./... -short

test-cover:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-race:
	go test ./... -race -short

test-bench:
	go test ./... -bench=. -benchmem
```

### Approach

- **Unit tests** live next to the code (`*_test.go`).
- **Table tests + `t.Run`** — the standard form.
- **`Example*`** — both documentation and quick output checks.
- **No heavy mocking frameworks** — interfaces plus hand-written fakes are enough.
- **Integration tests** — separate build tag (`//go:build integration`) and real dependencies via docker-compose.

---

## 14. Build and deploy

### makefile — typical targets

```makefile
# Run locally
run:
	go run ./cmd/<service>/main.go

# Static Linux build (amd64)
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	    go build -a -installsuffix cgo \
	    -ldflags '-extldflags "-static" -w -s' \
	    -o bin/amd64-alpine/<service> ./cmd/<service>/main.go

# Alpine RISC build
build-risc:
	CGO_ENABLED=0 GOOS=linux \
	    go build -ldflags='-s -w' \
	    -o bin/risc-alpine/<service> ./cmd/<service>/main.go

# Code quality
fmt:
	gofmt -s -w .

vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: fmt vet test-short

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
```

### CI/CD (GitLab skeleton)

```yaml
stages: [build, deploy, test]

build:
  stage: build
  script:
    - docker build -t $REGISTRY/<service>:$IMAGE_TAG -f docker/<service>.dockerfile .
    - docker push $REGISTRY/<service>:$IMAGE_TAG

deploy:
  stage: deploy
  script:
    - helm upgrade $RELEASE helm/ -n $NAMESPACE
    - kubectl rollout restart deployment/<service> -n $NAMESPACE

smoke:
  stage: test
  script:
    - ./test.sh <service>
```

---

## 15. Default stack

| Layer             | Dependency                                |
|-------------------|-------------------------------------------|
| HTTP router       | `github.com/go-chi/chi/v5`                |
| CORS              | `github.com/go-chi/cors`                  |
| JSON rendering    | `github.com/go-chi/render`                |
| ENV parsing       | `github.com/caarlos0/env/v11`             |
| .env files        | `github.com/joho/godotenv`                |
| Logging           | `log/slog` (stdlib)                       |
| UUID              | `github.com/google/uuid`                  |
| PostgreSQL        | `github.com/jackc/pgx/v5` + `jmoiron/sqlx`|
| ClickHouse        | `github.com/ClickHouse/clickhouse-go/v2`  |
| Redis             | `github.com/redis/go-redis/v9`            |
| Cron              | `github.com/robfig/cron/v3`               |
| Swagger           | `github.com/swaggo/swag` + `http-swagger/v2` |

**Deliberately avoided**:
- ORMs (`gorm`, `ent`) — prefer explicit SQL via `sqlx` / `pgx`.
- DI frameworks (`uber/fx`, `google/wire`) — manual DI is simpler and more transparent.
- Heavy loggers (`zap`, `logrus`) — `slog` is enough.
- HTTP frameworks (`gin`, `echo`, `fiber`) — `chi` stays compatible with stdlib without imposing abstractions.

---

## 16. Naming conventions

### Directories
- Lowercase, no underscores (`eventreceiver`, not `event_receiver`).
- Package name matches directory name.

### Files
| File                       | Purpose                                                      |
|----------------------------|--------------------------------------------------------------|
| `main.go`                  | Entry point, nothing else                                    |
| `app.go`                   | `App` struct, `Init()`, `Close()`                            |
| `di.go`                    | Dependency graph assembly                                    |
| `router.go`                | Route registration                                           |
| `http.go`                  | Root `Handler` struct                                        |
| `handler.go`               | Per-operation constructor + dependencies                     |
| `<operation>.go`           | HTTP method / domain method implementation                   |
| `service.go`               | `Service` struct + constructor                               |
| `fetcher.go`               | `Fetcher` struct + constructor                               |
| `connect.go`               | External-system connection factory                           |
| `config.go`                | `Config` struct and loader                                   |
| `<entity>DTO.go`           | DTO models                                                   |

### Constructors
- All structs: `func New<Type>(...) *Type`.
- External connections: `func Connect(cfg Config) (Client, error)`.

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

---

## 18. Checklist for a new service

- [ ] `cmd/<service>/` created with `main.go` and `app/{app.go,di.go}`.
- [ ] `internal/<service>/` with `handler/`, `service/`, and (as needed) `fetcher/`, `crons/`, `validator/`.
- [ ] Root `Config` aggregates subsystem configs; loaded from ENV / `.env`.
- [ ] `Init()` sets up loggers, config, connections (with retries), and cron.
- [ ] `InitializeDependencies()` builds the graph in four stages: services → handlers → crons → router.
- [ ] HTTP server has timeouts and graceful shutdown on SIGINT/SIGTERM.
- [ ] `/health` endpoint and a versioned `/v1` prefix.
- [ ] Standard middleware stack: Content-Type, Timeout, CORS, RateLimit.
- [ ] Unified `ErrorResponse` for every error response.
- [ ] Structured logger injected into all services and handlers.
- [ ] makefile with `run`, `build`, `test`, `lint`, `check` targets.
- [ ] CI pipeline: build → push image → deploy → smoke test.
- [ ] `go.mod` uses only the approved stack (see §15).
