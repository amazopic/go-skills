---
name: go-principles-and-layout
description: Use when laying out a new Go service repo or auditing structure: cmd/, internal/, pkg/common/, naming, placement rules, new-service checklist.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§1, §2, §16, §18)
related:
  - skills/methodology/00-canonical-full.md
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
