# Service Architecture Methodology

A reference for building backend services in Go. Distilled from `golang-service-methodology.md` (canonical, 930 lines) into 13 chapter-grouped skills plus the original full document.

## Read

- **Full canonical** (read end-to-end): [skills/methodology/00-canonical-full.md](skills/methodology/00-canonical-full.md)

## Skills (chapter-grouped)

| # | Skill | Covers |
|---|---|---|
| 01 | [Principles & Layout](skills/methodology/01-principles-and-layout.md) | §1 Principles · §2 Directory Layout · §16 Naming · §18 New-service checklist |
| 02 | [Layered Architecture](skills/methodology/02-layered-architecture.md) | §3 Transport / Service / Fetcher / DTO · §17 Cross-layer flow |
| 03 | [Bootstrap & DI](skills/methodology/03-bootstrap-and-di.md) | §4 main.go · app.go · di.go (four-stage manual DI) |
| 04 | [Configuration](skills/methodology/04-configuration.md) | §5 caarlos0/env, godotenv, per-subsystem Config |
| 05 | [External Connections](skills/methodology/05-external-connections.md) | §6 Connect() factories with retry + timeout |
| 06 | [Storage](skills/methodology/06-storage.md) | §7 PostgreSQL · ClickHouse batched · Redis Streams · in-memory TTL cache |
| 07 | [HTTP Transport](skills/methodology/07-http-transport.md) | §8 chi · middleware · rate-limit · context_key · ErrorResponse |
| 08 | [Background Jobs](skills/methodology/08-background-jobs.md) | §9 robfig/cron/v3 |
| 09 | [Logging](skills/methodology/09-logging.md) | §10 simple + structured (slog) two-tier |
| 10 | [Validation](skills/methodology/10-validation.md) | §11 Pure validator package, table tests |
| 11 | [Error Handling](skills/methodology/11-error-handling.md) | §12 wrap-with-%w, log once, map to HTTP |
| 12 | [Testing](skills/methodology/12-testing.md) | §13 unit / integration / example / race; makefile targets |
| 13 | [Build & Deploy](skills/methodology/13-build-and-deploy.md) | §14 makefile · §15 default approved stack |

## Default stack

The methodology converges on a deliberate stack — see [skill 13](skills/methodology/13-build-and-deploy.md) §15.

Deliberately avoided: ORMs, DI frameworks, heavy loggers, HTTP frameworks beyond chi.
