---
name: go-configuration
description: Use when adding env-driven configuration to a service: caarlos0/env tags, .env fallback via godotenv, per-subsystem Config aggregation, lazy derived methods via sync.Once.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§5)
related:
  - skills/methodology/00-canonical-full.md
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
