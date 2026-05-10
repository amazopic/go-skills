---
name: go-bootstrap-and-di
description: Use when wiring main.go + app.go + di.go: HTTP server timeouts, signal handling, graceful shutdown, four-stage manual DI graph (services → handlers → crons → router).
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§4)
related:
  - skills/methodology/00-canonical-full.md
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
