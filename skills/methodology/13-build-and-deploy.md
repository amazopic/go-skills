---
name: go-build-and-deploy
description: Use when setting up build/CI: makefile (run/build/build-risc/fmt/vet/lint/check/clean), static linux build flags, GitLab CI three-stage pipeline (build → deploy → smoke), default approved stack.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§14, §15)
related:
  - skills/methodology/00-canonical-full.md
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
