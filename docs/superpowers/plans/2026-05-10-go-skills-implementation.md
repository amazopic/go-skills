# go-skills Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the `go-skills` repository — a multilingual Go knowledge base (architectural patterns + service methodology) packaged as Claude Code skills, exposed via a Linear-styled GitHub Pages site, with placeholder folders for a future MCP server.

**Architecture:** Static repository — markdown content in `skills/` and `examples/`, vanilla CSS/JS site in `docs/`, no build pipeline. Two helper shell scripts (`make-llms.sh`, `make-sitemap.sh`). Mirrors `claude-customize/status-line` structure; uses `linear.app` design tokens from `awesome-design-md`.

**Tech Stack:** Markdown · vanilla HTML5 · CSS custom properties · ES modules · Bash · Go (for `examples/` packages, no service code) · GitHub Pages

**Spec:** `docs/superpowers/specs/2026-05-10-go-skills-design.md`

**Working dir:** `/Users/mylive/project/go-skills-md/go-skills/`

**Source archives (read-only — never modify):**
- `/Users/mylive/project/go-skills-md/go-old-pattern/golang-service-methodology.md`
- `/Users/mylive/project/go-skills-md/go-old-pattern/go-patterns-1/`
- `/Users/mylive/project/go-skills-md/go-old-pattern/go-patterns-2/`

**Reference templates (read-only — never modify):**
- `/Users/mylive/project/claude-customize/status-line/` — repo + site structure
- `/Users/mylive/project/md-data/deign/awesome-design-md/design-md/linear.app/DESIGN.md` — Linear tokens

---

## Phase 0 — Repo skeleton

### Task 0.1: Initialize repo structure

**Files:**
- Create: `.gitignore`
- Create: `.nojekyll`
- Create: `LICENSE`
- Create: `skills/`, `skills/methodology/`, `skills/creational/`, `skills/structural/`, `skills/behavioral/`, `skills/concurrency/`, `skills/synchronization/`, `skills/messaging/`, `skills/stability/`, `skills/profiling/`, `skills/idiom/`, `skills/anti-patterns/`
- Create: `examples/`, `commands/`, `mcp/`
- Create: `docs/`, `docs/css/`, `docs/js/`, `docs/assets/specimens/`

- [ ] **Step 1: Create directory tree**

```bash
cd /Users/mylive/project/go-skills-md/go-skills
mkdir -p skills/{methodology,creational,structural,behavioral,concurrency,synchronization,messaging,stability,profiling,idiom,anti-patterns}
mkdir -p examples/{creational,structural,behavioral,concurrency,synchronization,messaging,stability,profiling,idiom}
mkdir -p commands mcp
mkdir -p docs/{css,js,assets/specimens}
```

- [ ] **Step 2: Create `.gitignore`**

```
.DS_Store
*.log
node_modules/
.cache/
.idea/
.vscode/
coverage.out
coverage.html
bin/
```

- [ ] **Step 3: Create `.nojekyll`** (empty file)

```bash
touch /Users/mylive/project/go-skills-md/go-skills/.nojekyll
touch /Users/mylive/project/go-skills-md/go-skills/docs/.nojekyll
```

- [ ] **Step 4: Create `LICENSE`** (MIT — per spec §11 default)

```
MIT License

Copyright (c) 2026 Yevgeniy Achin

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

- [ ] **Step 5: Verify directory tree**

```bash
cd /Users/mylive/project/go-skills-md/go-skills && find . -maxdepth 3 -type d | sort
```

Expected output: 24 directories listed (the 12 `skills/*` subfolders, 9 `examples/*`, plus `docs`, `docs/css`, `docs/js`, `docs/assets`, `docs/assets/specimens`, `commands`, `mcp`, `docs/superpowers`, `docs/superpowers/specs`, `docs/superpowers/plans`).

- [ ] **Step 6: Initialize git and first commit**

```bash
cd /Users/mylive/project/go-skills-md/go-skills
git init -b main
git add .
git commit -m "chore: initial repo skeleton"
```

---

## Phase 1 — Methodology migration (1 canonical + 13 chapter skills)

### Task 1.1: Copy canonical full methodology

**Files:**
- Create: `skills/methodology/00-canonical-full.md`

- [ ] **Step 1: Copy source verbatim, prepend frontmatter**

Read source: `/Users/mylive/project/go-skills-md/go-old-pattern/golang-service-methodology.md`

Write `skills/methodology/00-canonical-full.md` with this prepended frontmatter, then the original body (preserve content exactly):

```markdown
---
name: go-service-architecture-canonical
description: Canonical full methodology for building backend services in Go — directory layout, layer separation, DI, configuration, transport, storage, jobs, logging, build, deploy. Use when reviewing or designing an entire Go service end-to-end. For chapter-scoped skills see 01–13.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md
---

[ORIGINAL 930-LINE BODY VERBATIM]
```

- [ ] **Step 2: Verify file content**

```bash
wc -l /Users/mylive/project/go-skills-md/go-skills/skills/methodology/00-canonical-full.md
head -15 /Users/mylive/project/go-skills-md/go-skills/skills/methodology/00-canonical-full.md
```

Expected: ≥ 930 lines, frontmatter `name: go-service-architecture-canonical` visible.

- [ ] **Step 3: Commit**

```bash
git add skills/methodology/00-canonical-full.md
git commit -m "feat(methodology): import canonical Go service architecture document"
```

### Task 1.2: Split methodology into 13 chapter skills

**Files (create all 13):**
- `skills/methodology/01-principles-and-layout.md` — covers §1, §2, §16, §18
- `skills/methodology/02-layered-architecture.md` — covers §3, §17
- `skills/methodology/03-bootstrap-and-di.md` — §4
- `skills/methodology/04-configuration.md` — §5
- `skills/methodology/05-external-connections.md` — §6
- `skills/methodology/06-storage.md` — §7
- `skills/methodology/07-http-transport.md` — §8
- `skills/methodology/08-background-jobs.md` — §9
- `skills/methodology/09-logging.md` — §10
- `skills/methodology/10-validation.md` — §11
- `skills/methodology/11-error-handling.md` — §12
- `skills/methodology/12-testing.md` — §13
- `skills/methodology/13-build-and-deploy.md` — §14, §15

- [ ] **Step 1: For each chapter file, write frontmatter + extracted body**

Use this template per skill (replace `<placeholders>`; copy chapter content from `00-canonical-full.md`):

```markdown
---
name: <hyphenated-name>
description: <one-line — when this skill should fire>
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§<chapters>)
related:
  - skills/methodology/00-canonical-full.md
---

# <Chapter Title>

<Body — copy verbatim from canonical sections; preserve code fences and tables>
```

Concrete frontmatter values (copy these exactly):

| File | name | description | sources |
|---|---|---|---|
| 01 | `go-principles-and-layout` | Use when laying out a new Go service repo or auditing structure: cmd/, internal/, pkg/common/, naming, placement rules, new-service checklist. | golang-service-methodology.md (§1, §2, §16, §18) |
| 02 | `go-layered-architecture` | Use when designing transport/service/fetcher boundaries and the cross-layer flow of an HTTP request through middleware → handler → service → fetcher → DB. | golang-service-methodology.md (§3, §17) |
| 03 | `go-bootstrap-and-di` | Use when wiring main.go + app.go + di.go: HTTP server timeouts, signal handling, graceful shutdown, four-stage manual DI graph (services → handlers → crons → router). | golang-service-methodology.md (§4) |
| 04 | `go-configuration` | Use when adding env-driven configuration to a service: caarlos0/env tags, .env fallback via godotenv, per-subsystem Config aggregation, lazy derived methods via sync.Once. | golang-service-methodology.md (§5) |
| 05 | `go-external-connections` | Use when writing a Connect() factory for redis/clickhouse/postgres: retry loop, timeouts, ping-after-dial verification. | golang-service-methodology.md (§6) |
| 06 | `go-storage` | Use when choosing or wiring storage: PostgreSQL via pgx/sqlx, ClickHouse batched writes via PrepareBatch, Redis Streams for durable queues, in-memory TTL cache with janitor. | golang-service-methodology.md (§7) |
| 07 | `go-http-transport` | Use when building or auditing the HTTP layer: chi router, middleware stack (Content-Type, Timeout, CORS, RateLimit), per-IP rate limiting via in-memory cache, context_key constants, unified ErrorResponse. | golang-service-methodology.md (§8) |
| 08 | `go-background-jobs` | Use when scheduling background tasks with robfig/cron/v3 (second precision); register via di.go after dependency graph. | golang-service-methodology.md (§9) |
| 09 | `go-logging` | Use when wiring loggers: simple coloured logger for bootstrap/connect, slog-based structured logger for services and handlers; log every error exactly once. | golang-service-methodology.md (§10) |
| 10 | `go-validation` | Use when writing input validators: pure functions returning error, table-driven tests, Example* for documentation; handler runs validation before service. | golang-service-methodology.md (§11) |
| 11 | `go-error-handling` | Use when propagating, wrapping, and logging errors: error always last, wrap with %w, preserve context.Canceled / DeadlineExceeded, map to common.ErrorResponse for HTTP. | golang-service-methodology.md (§12) |
| 12 | `go-testing` | Use when writing unit/integration/example/race tests: table tests + t.Run, Example*, fakes over mocks, integration build tag, makefile targets test/test-short/test-cover/test-race/test-bench. | golang-service-methodology.md (§13) |
| 13 | `go-build-and-deploy` | Use when setting up build/CI: makefile (run/build/build-risc/fmt/vet/lint/check/clean), static linux build flags, GitLab CI three-stage pipeline (build → deploy → smoke), default approved stack. | golang-service-methodology.md (§14, §15) |

- [ ] **Step 2: Verify each file is non-empty and frontmatter parses**

```bash
for f in /Users/mylive/project/go-skills-md/go-skills/skills/methodology/*.md; do
  echo "== $f =="
  head -3 "$f"
done
```

Expected: 14 files (00 + 01–13), each starts with `---` and `name: ...`.

- [ ] **Step 3: Commit**

```bash
git add skills/methodology/
git commit -m "feat(methodology): split canonical doc into 13 chapter-grouped skills"
```

---

## Phase 2 — Pattern skills (theory) from go-patterns-1

### Task 2.1: Copy go-patterns-1 markdown into skills/<category>/<pattern>.md with frontmatter

**Files (create — list of all .md files in go-patterns-1):**
- `skills/behavioral/observer.md`, `skills/behavioral/strategy.md`
- `skills/concurrency/bounded-parallelism.md`, `skills/concurrency/generator.md`, `skills/concurrency/parallelism.md`
- `skills/creational/builder.md`, `skills/creational/factory-method.md`, `skills/creational/object-pool.md`, `skills/creational/singleton.md`
- `skills/idiom/functional-options.md`
- `skills/messaging/fan-in.md`, `skills/messaging/fan-out.md`, `skills/messaging/publish-subscribe.md`
- `skills/profiling/timing.md`
- `skills/stability/circuit-breaker.md`
- `skills/structural/decorator.md`, `skills/structural/proxy.md`
- `skills/synchronization/semaphore.md`

(Note: source files use underscores; we rename to hyphens for URL-friendliness.)

- [ ] **Step 1: For each pattern, copy body and prepend frontmatter**

Frontmatter template:

```yaml
---
name: <category>-<hyphenated-pattern>
description: <one-line description>
category: <creational|structural|behavioral|concurrency|synchronization|messaging|stability|profiling|idiom>
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/<original-path>
example: examples/<category>/<pattern>/
---
```

Use these exact descriptions:

| File | name | description |
|---|---|---|
| `behavioral/observer.md` | `behavioral-observer` | Observer pattern in Go — a one-to-many publish/notify relationship with explicit Subject and Observer interfaces. Use when modeling subscription-style updates. |
| `behavioral/strategy.md` | `behavioral-strategy` | Strategy pattern — interchangeable algorithms behind a common interface, swapped at runtime. Use when an operation has multiple implementations selected by configuration. |
| `concurrency/bounded-parallelism.md` | `concurrency-bounded-parallelism` | Bounded parallelism — process N items with at most K workers concurrently using a worker-pool goroutine + channels. Use to cap concurrency on an unbounded input list. |
| `concurrency/generator.md` | `concurrency-generator` | Generator — a goroutine that emits values on a channel, returned to the caller as a `<-chan T`. Use to lazy-stream values without exposing channel internals. |
| `concurrency/parallelism.md` | `concurrency-parallelism` | Parallel scatter-gather — fan a workload across goroutines and collect results via WaitGroup or channel. Use for embarrassingly parallel CPU-bound work. |
| `creational/builder.md` | `creational-builder` | Builder pattern — step-by-step construction of a complex object via fluent methods. Use when a struct has many optional fields and you need a readable, validated build sequence. |
| `creational/factory-method.md` | `creational-factory-method` | Factory Method — defer instantiation to a function so the caller picks the variant. Use to abstract concrete-type selection behind a constructor. |
| `creational/object-pool.md` | `creational-object-pool` | Object Pool — pre-allocate and reuse expensive resources (DB connections, large buffers) instead of creating per request. Use sync.Pool when zero-allocation matters. |
| `creational/singleton.md` | `creational-singleton` | Singleton — exactly one instance, lazy-initialized via sync.Once. Use sparingly; prefer DI; only when an instance must truly be process-global (e.g., metric registries). |
| `idiom/functional-options.md` | `idiom-functional-options` | Functional Options — vary constructors with `func(*T)` option arguments. The idiomatic Go alternative to Builder for optional configuration. |
| `messaging/fan-in.md` | `messaging-fan-in` | Fan-In — multiplex multiple input channels into a single output channel via a goroutine that select-reads from all inputs. Use when downstream wants a unified stream. |
| `messaging/fan-out.md` | `messaging-fan-out` | Fan-Out — distribute work from one channel across N worker goroutines. Use to scale processing throughput when items are independent. |
| `messaging/publish-subscribe.md` | `messaging-publish-subscribe` | Publish/Subscribe — broker-mediated event distribution: publishers emit topics, subscribers receive matching messages. Use for decoupled cross-component communication. |
| `profiling/timing.md` | `profiling-timing` | Timing functions — measure latency with `defer fn(time.Now())` and a closure that logs the elapsed duration. Use to add lightweight per-call timing. |
| `stability/circuit-breaker.md` | `stability-circuit-breaker` | Circuit Breaker — track failures and open the circuit after threshold to fail fast and let the dependency recover. Use for unreliable downstreams (external APIs). |
| `structural/decorator.md` | `structural-decorator` | Decorator — wrap a value in a function or type that adds behavior while preserving the same interface. Use for cross-cutting concerns: logging, caching, retry. |
| `structural/proxy.md` | `structural-proxy` | Proxy — interpose a stand-in object that controls access to the real one (lazy load, access control, remote call). Use for transparent delegation. |
| `synchronization/semaphore.md` | `synchronization-semaphore` | Semaphore — limit concurrent access to N resource slots via a buffered channel as a counting semaphore. Use to throttle parallelism without spawning a worker pool. |

- [ ] **Step 2: Run a sanity check — every file has frontmatter**

```bash
for f in /Users/mylive/project/go-skills-md/go-skills/skills/*/*.md; do
  if ! head -1 "$f" | grep -q '^---$'; then
    echo "MISSING FRONTMATTER: $f"
  fi
done
```

Expected: no output (no missing frontmatter).

- [ ] **Step 3: Commit**

```bash
git add skills/
git commit -m "feat(skills): import 18 pattern skills from go-patterns-1 with frontmatter"
```

### Task 2.2: Create stub skills for the SUMMARY taxonomy gaps

Per spec §4.2: SUMMARY lists patterns we don't have bodies for. Create stub skill files so the catalog index is complete and discoverable.

**Files (stubs to create):**
- `skills/creational/abstract-factory.md`, `skills/creational/prototype.md`
- `skills/structural/adapter.md`, `skills/structural/bridge.md`, `skills/structural/composite.md`, `skills/structural/facade.md`, `skills/structural/flyweight.md`
- `skills/behavioral/chain-of-responsibility.md`, `skills/behavioral/command.md`, `skills/behavioral/iterator.md`, `skills/behavioral/mediator.md`, `skills/behavioral/memento.md`, `skills/behavioral/state.md`, `skills/behavioral/template-method.md`, `skills/behavioral/visitor.md`, `skills/behavioral/registry.md`
- `skills/concurrency/n-barrier.md`, `skills/concurrency/broadcast.md`, `skills/concurrency/coroutines.md`, `skills/concurrency/reactor.md`, `skills/concurrency/producer-consumer.md`
- `skills/synchronization/condition-variable.md`, `skills/synchronization/mutex.md`, `skills/synchronization/monitor.md`, `skills/synchronization/read-write-lock.md`
- `skills/messaging/futures-promises.md`, `skills/messaging/push-pull.md`
- `skills/stability/bulkhead.md`, `skills/stability/deadline.md`, `skills/stability/fail-fast.md`, `skills/stability/handshaking.md`, `skills/stability/steady-state.md`
- `skills/idiom/specification.md`
- `skills/anti-patterns/cascading-failures.md`

- [ ] **Step 1: Stub template — write each file**

Each stub uses this template (no body content yet):

```markdown
---
name: <category>-<hyphenated-pattern>
description: <pattern-name> — stub. Catalog placeholder; body to be filled in a later iteration. See SUMMARY at https://github.com/tmrts/go-patterns and go-patterns-2 archive for source material.
category: <category>
status: stub
---

# <Pattern Name>

> **Status:** stub — body not yet written. This entry exists so the catalog and index remain complete.
>
> **Sources to mine when filling:**
> - `go-old-pattern/go-patterns-1/SUMMARY.md` (taxonomy)
> - `go-old-pattern/go-patterns-2/<Category>/<Pattern>/` (if present — runnable code + README)
> - Public references: <a href="https://github.com/tmrts/go-patterns">tmrts/go-patterns</a>, <a href="https://github.com/AlexanderGrom/go-patterns">AlexanderGrom/go-patterns</a>
```

(Replace `<category>`, `<hyphenated-pattern>`, `<Pattern Name>` accordingly.)

For patterns where `go-patterns-2` has a README, also include:

```markdown
> **Existing example:** see `examples/<category>/<pattern>/`
```

- [ ] **Step 2: Verify count**

```bash
ls /Users/mylive/project/go-skills-md/go-skills/skills/*/*.md | wc -l
```

Expected: 14 (methodology) + 18 (filled patterns) + 33 (stubs) = **65 files**.

- [ ] **Step 3: Commit**

```bash
git add skills/
git commit -m "feat(skills): add 33 stub pattern entries to complete catalog taxonomy"
```

---

## Phase 3 — Examples (runnable Go code)

### Task 3.1: Copy go-patterns-2 patterns into examples/

**Source:** `/Users/mylive/project/go-skills-md/go-old-pattern/go-patterns-2/<Category>/<Pattern>/{README.md,*.go,*_test.go}`

**Target naming:** lowercase category, hyphenated pattern. e.g. `Creational/AbstractFactory/` → `examples/creational/abstract-factory/`.

- [ ] **Step 1: Run the migration shell snippet**

```bash
SRC=/Users/mylive/project/go-skills-md/go-old-pattern/go-patterns-2
DST=/Users/mylive/project/go-skills-md/go-skills/examples

declare -A MAP=(
  [Creational/AbstractFactory]=creational/abstract-factory
  [Creational/Builder]=creational/builder
  [Creational/FactoryMethod]=creational/factory-method
  [Creational/Prototype]=creational/prototype
  [Creational/Singleton]=creational/singleton
  [Structural/Adapter]=structural/adapter
  [Structural/Bridge]=structural/bridge
  [Structural/Composite]=structural/composite
  [Structural/Decorator]=structural/decorator
  [Structural/Facade]=structural/facade
  [Structural/Flyweight]=structural/flyweight
  [Structural/Proxy]=structural/proxy
  [Behavioral/ChainOfResponsibility]=behavioral/chain-of-responsibility
  [Behavioral/Command]=behavioral/command
  [Behavioral/Iterator]=behavioral/iterator
  [Behavioral/Mediator]=behavioral/mediator
  [Behavioral/Memento]=behavioral/memento
  [Behavioral/Observer]=behavioral/observer
  [Behavioral/State]=behavioral/state
  [Behavioral/Strategy]=behavioral/strategy
  [Behavioral/TemplateMethod]=behavioral/template-method
  [Behavioral/Visitor]=behavioral/visitor
  [Unsorted/Specification]=idiom/specification
)

for src_rel in "${!MAP[@]}"; do
  dst_rel="${MAP[$src_rel]}"
  mkdir -p "$DST/$dst_rel"
  cp "$SRC/$src_rel/"*.go "$DST/$dst_rel/" 2>/dev/null || true
  cp "$SRC/$src_rel/"README.md "$DST/$dst_rel/" 2>/dev/null || true
done
```

- [ ] **Step 2: Initialize a Go module per example or one shared module**

For simplicity and offline-friendliness, create ONE shared `go.mod` at `examples/`:

```bash
cd /Users/mylive/project/go-skills-md/go-skills/examples
cat > go.mod <<'EOF'
module github.com/amazopic/go-skills/examples

go 1.21
EOF
```

- [ ] **Step 3: Fix package paths so Go compiles**

Each `<pattern>/<pattern>.go` file in go-patterns-2 starts with `package <name>`. Verify the package declarations are consistent (one per directory). Run:

```bash
cd /Users/mylive/project/go-skills-md/go-skills/examples
go vet ./... 2>&1 | head -40
```

If `go vet` fails because of conflicting package names within one folder, fix by editing the offending files (likely none — go-patterns-2 is well-formed).

- [ ] **Step 4: Run tests**

```bash
cd /Users/mylive/project/go-skills-md/go-skills/examples
go test ./... -short 2>&1 | tail -30
```

Expected: `ok` for each example package. If failures occur, log them under `examples/KNOWN-ISSUES.md` and continue (do not block).

- [ ] **Step 5: Commit**

```bash
git add examples/
git commit -m "feat(examples): import 23 runnable Go pattern packages from go-patterns-2"
```

---

## Phase 4 — Top-level catalog files

### Task 4.1: Create `PATTERNS.md`

**File:** `PATTERNS.md`

- [ ] **Step 1: Write the catalog index**

```markdown
# Patterns Catalog

A complete index of patterns in this knowledge base. Filled entries link to their skill body and runnable example; stubs are listed for taxonomy completeness.

Legend: ✓ filled · ◐ partial (theory only) · ○ stub.

## Creational

| Pattern | Status | Skill | Example |
|---|---|---|---|
| Abstract Factory | ◐ | [skills/creational/abstract-factory.md](skills/creational/abstract-factory.md) | [examples/creational/abstract-factory](examples/creational/abstract-factory) |
| Builder | ✓ | [skills/creational/builder.md](skills/creational/builder.md) | [examples/creational/builder](examples/creational/builder) |
| Factory Method | ✓ | [skills/creational/factory-method.md](skills/creational/factory-method.md) | [examples/creational/factory-method](examples/creational/factory-method) |
| Object Pool | ✓ | [skills/creational/object-pool.md](skills/creational/object-pool.md) | — |
| Prototype | ◐ | [skills/creational/prototype.md](skills/creational/prototype.md) | [examples/creational/prototype](examples/creational/prototype) |
| Singleton | ✓ | [skills/creational/singleton.md](skills/creational/singleton.md) | [examples/creational/singleton](examples/creational/singleton) |

## Structural

| Pattern | Status | Skill | Example |
|---|---|---|---|
| Adapter | ◐ | [skills/structural/adapter.md](skills/structural/adapter.md) | [examples/structural/adapter](examples/structural/adapter) |
| Bridge | ◐ | [skills/structural/bridge.md](skills/structural/bridge.md) | [examples/structural/bridge](examples/structural/bridge) |
| Composite | ◐ | [skills/structural/composite.md](skills/structural/composite.md) | [examples/structural/composite](examples/structural/composite) |
| Decorator | ✓ | [skills/structural/decorator.md](skills/structural/decorator.md) | [examples/structural/decorator](examples/structural/decorator) |
| Facade | ◐ | [skills/structural/facade.md](skills/structural/facade.md) | [examples/structural/facade](examples/structural/facade) |
| Flyweight | ◐ | [skills/structural/flyweight.md](skills/structural/flyweight.md) | [examples/structural/flyweight](examples/structural/flyweight) |
| Proxy | ✓ | [skills/structural/proxy.md](skills/structural/proxy.md) | [examples/structural/proxy](examples/structural/proxy) |

## Behavioral

| Pattern | Status | Skill | Example |
|---|---|---|---|
| Chain of Responsibility | ◐ | [skills/behavioral/chain-of-responsibility.md](skills/behavioral/chain-of-responsibility.md) | [examples/behavioral/chain-of-responsibility](examples/behavioral/chain-of-responsibility) |
| Command | ◐ | [skills/behavioral/command.md](skills/behavioral/command.md) | [examples/behavioral/command](examples/behavioral/command) |
| Iterator | ◐ | [skills/behavioral/iterator.md](skills/behavioral/iterator.md) | [examples/behavioral/iterator](examples/behavioral/iterator) |
| Mediator | ◐ | [skills/behavioral/mediator.md](skills/behavioral/mediator.md) | [examples/behavioral/mediator](examples/behavioral/mediator) |
| Memento | ◐ | [skills/behavioral/memento.md](skills/behavioral/memento.md) | [examples/behavioral/memento](examples/behavioral/memento) |
| Observer | ✓ | [skills/behavioral/observer.md](skills/behavioral/observer.md) | [examples/behavioral/observer](examples/behavioral/observer) |
| Registry | ○ | [skills/behavioral/registry.md](skills/behavioral/registry.md) | — |
| State | ◐ | [skills/behavioral/state.md](skills/behavioral/state.md) | [examples/behavioral/state](examples/behavioral/state) |
| Strategy | ✓ | [skills/behavioral/strategy.md](skills/behavioral/strategy.md) | [examples/behavioral/strategy](examples/behavioral/strategy) |
| Template Method | ◐ | [skills/behavioral/template-method.md](skills/behavioral/template-method.md) | [examples/behavioral/template-method](examples/behavioral/template-method) |
| Visitor | ◐ | [skills/behavioral/visitor.md](skills/behavioral/visitor.md) | [examples/behavioral/visitor](examples/behavioral/visitor) |

## Concurrency

| Pattern | Status | Skill | Example |
|---|---|---|---|
| N-Barrier | ○ | [skills/concurrency/n-barrier.md](skills/concurrency/n-barrier.md) | — |
| Bounded Parallelism | ✓ | [skills/concurrency/bounded-parallelism.md](skills/concurrency/bounded-parallelism.md) | — |
| Broadcast | ○ | [skills/concurrency/broadcast.md](skills/concurrency/broadcast.md) | — |
| Coroutines | ○ | [skills/concurrency/coroutines.md](skills/concurrency/coroutines.md) | — |
| Generator | ✓ | [skills/concurrency/generator.md](skills/concurrency/generator.md) | — |
| Parallelism | ✓ | [skills/concurrency/parallelism.md](skills/concurrency/parallelism.md) | — |
| Producer Consumer | ○ | [skills/concurrency/producer-consumer.md](skills/concurrency/producer-consumer.md) | — |
| Reactor | ○ | [skills/concurrency/reactor.md](skills/concurrency/reactor.md) | — |

## Synchronization

| Pattern | Status | Skill | Example |
|---|---|---|---|
| Condition Variable | ○ | [skills/synchronization/condition-variable.md](skills/synchronization/condition-variable.md) | — |
| Mutex | ○ | [skills/synchronization/mutex.md](skills/synchronization/mutex.md) | — |
| Monitor | ○ | [skills/synchronization/monitor.md](skills/synchronization/monitor.md) | — |
| Read-Write Lock | ○ | [skills/synchronization/read-write-lock.md](skills/synchronization/read-write-lock.md) | — |
| Semaphore | ✓ | [skills/synchronization/semaphore.md](skills/synchronization/semaphore.md) | — |

## Messaging

| Pattern | Status | Skill | Example |
|---|---|---|---|
| Fan-In | ✓ | [skills/messaging/fan-in.md](skills/messaging/fan-in.md) | — |
| Fan-Out | ✓ | [skills/messaging/fan-out.md](skills/messaging/fan-out.md) | — |
| Futures & Promises | ○ | [skills/messaging/futures-promises.md](skills/messaging/futures-promises.md) | — |
| Publish/Subscribe | ✓ | [skills/messaging/publish-subscribe.md](skills/messaging/publish-subscribe.md) | — |
| Push & Pull | ○ | [skills/messaging/push-pull.md](skills/messaging/push-pull.md) | — |

## Stability

| Pattern | Status | Skill | Example |
|---|---|---|---|
| Bulkhead | ○ | [skills/stability/bulkhead.md](skills/stability/bulkhead.md) | — |
| Circuit Breaker | ✓ | [skills/stability/circuit-breaker.md](skills/stability/circuit-breaker.md) | — |
| Deadline | ○ | [skills/stability/deadline.md](skills/stability/deadline.md) | — |
| Fail-Fast | ○ | [skills/stability/fail-fast.md](skills/stability/fail-fast.md) | — |
| Handshaking | ○ | [skills/stability/handshaking.md](skills/stability/handshaking.md) | — |
| Steady-State | ○ | [skills/stability/steady-state.md](skills/stability/steady-state.md) | — |

## Profiling

| Pattern | Status | Skill | Example |
|---|---|---|---|
| Timing | ✓ | [skills/profiling/timing.md](skills/profiling/timing.md) | — |

## Idioms

| Pattern | Status | Skill | Example |
|---|---|---|---|
| Functional Options | ✓ | [skills/idiom/functional-options.md](skills/idiom/functional-options.md) | — |
| Specification | ◐ | [skills/idiom/specification.md](skills/idiom/specification.md) | [examples/idiom/specification](examples/idiom/specification) |

## Anti-Patterns

| Pattern | Status | Skill | Example |
|---|---|---|---|
| Cascading Failures | ○ | [skills/anti-patterns/cascading-failures.md](skills/anti-patterns/cascading-failures.md) | — |
```

- [ ] **Step 2: Commit**

```bash
git add PATTERNS.md
git commit -m "docs: add full pattern catalog index (PATTERNS.md)"
```

### Task 4.2: Create `METHODOLOGY.md`

**File:** `METHODOLOGY.md`

- [ ] **Step 1: Write methodology summary**

```markdown
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
```

- [ ] **Step 2: Commit**

```bash
git add METHODOLOGY.md
git commit -m "docs: add methodology summary index (METHODOLOGY.md)"
```

### Task 4.3: Create `llms.txt`

**File:** `llms.txt`

- [ ] **Step 1: Write the AI-crawler index**

```
# go-skills

> A Go knowledge base — architectural patterns and a service-building methodology — packaged as Claude Code skills, exposed via a static site, and (soon) an MCP server. Source: methodology + 50+ pattern skills (theory + runnable examples in some cases).

Author: Yevgeniy Achin <amazopic@gmail.com>
License: MIT
Languages: 13 (en, ru, fr, de, uk, sl, it, es, zh, ja, ko, ar)

## Documentation

- [README](README.md): overview, install, link to site.
- [PATTERNS.md](PATTERNS.md): full catalog of 50+ patterns across 9 categories.
- [METHODOLOGY.md](METHODOLOGY.md): index of 13 chapter-skills + canonical full doc.
- [LICENSE](LICENSE): MIT.

## Skills (Claude Code)

- [skills/README.md](skills/README.md): how to install as Claude Code skills.
- [skills/methodology/00-canonical-full.md](skills/methodology/00-canonical-full.md): full Go service architecture doc.
- [skills/methodology/01..13-*.md](skills/methodology/): 13 chapter-grouped methodology skills.
- [skills/<category>/*.md](skills/): per-pattern skills across creational, structural, behavioral, concurrency, synchronization, messaging, stability, profiling, idiom, anti-patterns.

## Examples

- [examples/](examples/): 23 runnable Go pattern packages with `_test.go`. One shared module: `github.com/amazopic/go-skills/examples`. Run: `cd examples && go test ./...`.

## Site

- [docs/index.html](docs/index.html): GitHub Pages site (Linear design system).
- [docs/llms-full.txt](docs/llms-full.txt): full corpus mirror.

## Future

- [mcp/README.md](mcp/README.md): roadmap for the MCP server (no implementation in v1).
- [commands/README.md](commands/README.md): roadmap for slash-command integration.
```

- [ ] **Step 2: Commit**

```bash
git add llms.txt
git commit -m "docs: add llms.txt AI-crawler index"
```

### Task 4.4: Stub READMEs for `commands/`, `mcp/`, `skills/`

**Files:**
- Create: `commands/README.md`
- Create: `mcp/README.md`
- Create: `skills/README.md`

- [ ] **Step 1: `commands/README.md`**

```markdown
# Slash Commands (placeholder)

This folder will hold Claude Code slash commands once they are designed.

Planned: `/go-skill <pattern>` to inject a specific pattern's frontmatter+body into Claude Code context on demand.

Status: not yet implemented in v1.
```

- [ ] **Step 2: `mcp/README.md`**

```markdown
# MCP Server (roadmap)

This folder is reserved for an MCP (Model Context Protocol) server that will expose the `skills/` content as MCP tools and resources.

Planned shape:
- `tools.list_patterns` — returns the catalog with filters (category, status).
- `tools.get_pattern` — returns body + example code for a given pattern slug.
- `resources://methodology/<n>` — pinned methodology chapters as resources.

Status: not yet implemented. See spec `docs/superpowers/specs/2026-05-10-go-skills-design.md` §0 and §12 for context.
```

- [ ] **Step 3: `skills/README.md`**

```markdown
# Skills

This folder packages the knowledge base as Claude Code-compatible skills. Every file (except this README) starts with YAML frontmatter declaring `name`, `description`, `category`, `sources`, and (where applicable) `example`.

## Install as Claude Code skills

```bash
# 1. Clone the repo
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills

# 2. Add to your Claude Code skills path (~/.claude/settings.json):
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}

# 3. Restart Claude Code.
```

Now Claude Code will surface skills from this repo when their `description` matches your task.

## Layout

- `methodology/` — 13 chapter-grouped skills + 1 canonical full doc.
- `creational/`, `structural/`, `behavioral/` — GoF + Go-flavored patterns.
- `concurrency/`, `synchronization/`, `messaging/` — Go's strengths.
- `stability/`, `profiling/`, `idiom/`, `anti-patterns/` — operational and Go-specific.

See `../PATTERNS.md` for the full catalog with status (filled / partial / stub).
```

- [ ] **Step 4: Commit**

```bash
git add commands/README.md mcp/README.md skills/README.md
git commit -m "docs: add stub READMEs for commands/, mcp/, skills/"
```

---

## Phase 5 — Site CSS (Linear design tokens)

### Task 5.1: `docs/css/tokens.css`

**File:** `docs/css/tokens.css`

- [ ] **Step 1: Write the Linear-aligned tokens**

```css
/* ─────────────────────────────────────────────────────────────────────
   DESIGN TOKENS — go-skills
   Adapted from Linear (#010102 canvas, #5e6ad2 lavender accent).
   ───────────────────────────────────────────────────────────────────── */

:root {
  /* Colors — Linear scale */
  --canvas:           #010102;
  --surface-1:        #0f1011;
  --surface-2:        #141516;
  --surface-3:        #18191a;
  --surface-4:        #191a1b;

  --hairline:         #23252a;
  --hairline-strong:  #34343a;
  --hairline-tertiary:#3e3e44;

  --ink:              #f7f8f8;
  --ink-muted:        #d0d6e0;
  --ink-subtle:       #8a8f98;
  --ink-tertiary:     #62666d;

  --accent:           #5e6ad2;
  --accent-hover:     #828fff;
  --accent-focus:     #5e69d1;
  --on-accent:        #ffffff;

  --semantic-success: #27a644;

  /* Inverse (light flip — used in code-card backgrounds when needed) */
  --inverse-canvas:   #ffffff;
  --inverse-ink:      #000000;

  /* Typography stacks */
  --font-display: "Editorial New", "PP Editorial New", "Linear Display", "Inter", -apple-system, "Helvetica Neue", system-ui, sans-serif;
  --font-body:    "Neue Montreal", "PP Neue Montreal", "Inter", -apple-system, "Helvetica Neue", system-ui, sans-serif;
  --font-mono:    "Geist Mono", ui-monospace, "SF Mono", "Menlo", monospace;

  /* Type scale */
  --fs-xs:   0.75rem;       /* 12 */
  --fs-sm:   0.8125rem;     /* 13 */
  --fs-base: 0.9375rem;     /* 15 — Linear-ish body */
  --fs-md:   1.125rem;      /* 18 */
  --fs-lg:   1.375rem;      /* 22 — card title */
  --fs-xl:   1.75rem;       /* 28 — headline */
  --fs-2xl:  2.5rem;        /* 40 — display-md */
  --fs-3xl:  3.5rem;        /* 56 — display-lg */
  --fs-4xl:  5rem;          /* 80 — display-xl */

  /* Line heights */
  --lh-tight:  1.05;
  --lh-snug:   1.20;
  --lh-base:   1.45;
  --lh-loose:  1.65;

  /* Letter spacing — Linear measured negatives */
  --tr-mega:   -0.04em;
  --tr-display:-0.025em;
  --tr-snug:   -0.012em;
  --tr-normal:  0;
  --tr-wide:    0.04em;

  /* Radii */
  --r-sm: 4px;
  --r-md: 8px;
  --r-lg: 12px;

  /* Layout */
  --max-w:     1280px;
  --gutter-x:  clamp(1rem, 4vw, 3.5rem);
  --gutter-y:  clamp(2.5rem, 8vw, 6rem);

  /* Hairline gradient (used for card borders) */
  --hairline-gradient: linear-gradient(180deg, var(--hairline-strong), var(--hairline) 50%, transparent);

  /* Motion */
  --motion-fast:  120ms cubic-bezier(.2,.8,.2,1);
  --motion-base:  220ms cubic-bezier(.2,.8,.2,1);
}
```

- [ ] **Step 2: Commit**

```bash
git add docs/css/tokens.css
git commit -m "feat(site): add Linear-derived design tokens"
```

### Task 5.2: `docs/css/base.css`

**File:** `docs/css/base.css`

- [ ] **Step 1: Write reset + typography + layout primitives**

```css
/* base.css — reset, typography, layout primitives */

*, *::before, *::after { box-sizing: border-box; }

html, body {
  margin: 0;
  padding: 0;
  background: var(--canvas);
  color: var(--ink);
  font-family: var(--font-body);
  font-size: var(--fs-base);
  line-height: var(--lh-base);
  -webkit-font-smoothing: antialiased;
  text-rendering: optimizeLegibility;
}

body { min-height: 100vh; }

a {
  color: inherit;
  text-decoration: none;
  border-bottom: 1px solid currentColor;
  transition: color var(--motion-fast);
}
a:hover { color: var(--accent-hover); }

h1, h2, h3, h4, h5 {
  font-family: var(--font-display);
  margin: 0;
  letter-spacing: var(--tr-display);
  font-weight: 600;
}

h1 { font-size: var(--fs-4xl); line-height: var(--lh-tight); }
h2 { font-size: var(--fs-3xl); line-height: var(--lh-snug); }
h3 { font-size: var(--fs-2xl); line-height: var(--lh-snug); }
h4 { font-size: var(--fs-xl);  line-height: var(--lh-snug); }
h5 { font-size: var(--fs-lg);  line-height: var(--lh-snug); }

p { margin: 0 0 1em; color: var(--ink-muted); }

code, pre, kbd, samp {
  font-family: var(--font-mono);
  font-size: 0.9em;
}

pre {
  background: var(--surface-1);
  border: 1px solid var(--hairline);
  border-radius: var(--r-md);
  padding: 1rem 1.25rem;
  overflow-x: auto;
  margin: 0;
  color: var(--ink);
  line-height: 1.55;
}

::selection { background: var(--accent); color: var(--on-accent); }

.wrap {
  max-width: var(--max-w);
  margin: 0 auto;
  padding-inline: var(--gutter-x);
}

.section { padding-block: var(--gutter-y); border-top: 1px solid var(--hairline); }
.section:first-of-type { border-top: 0; }

.grid { display: grid; gap: 1.5rem; }
.grid-2 { grid-template-columns: repeat(2, minmax(0, 1fr)); }
.grid-3 { grid-template-columns: repeat(3, minmax(0, 1fr)); }
.grid-4 { grid-template-columns: repeat(4, minmax(0, 1fr)); }

@media (max-width: 768px) {
  .grid-2, .grid-3, .grid-4 { grid-template-columns: 1fr; }
}

.eyebrow {
  font-family: var(--font-mono);
  font-size: var(--fs-xs);
  letter-spacing: var(--tr-wide);
  text-transform: uppercase;
  color: var(--ink-subtle);
  margin-bottom: 0.75rem;
  display: inline-block;
}

.muted { color: var(--ink-subtle); }
.tight { letter-spacing: var(--tr-snug); }

.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  border-radius: var(--r-sm);
  border: 1px solid var(--hairline-strong);
  background: var(--surface-1);
  color: var(--ink);
  font-family: var(--font-body);
  font-size: var(--fs-sm);
  font-weight: 500;
  cursor: pointer;
  transition: background var(--motion-fast), border-color var(--motion-fast);
}
.btn:hover { background: var(--surface-2); border-color: var(--hairline-tertiary); }
.btn-primary {
  background: var(--accent);
  border-color: var(--accent);
  color: var(--on-accent);
}
.btn-primary:hover { background: var(--accent-hover); border-color: var(--accent-hover); }

.card {
  background: var(--surface-1);
  border: 1px solid var(--hairline);
  border-radius: var(--r-md);
  padding: 1.5rem;
  transition: border-color var(--motion-base), background var(--motion-base);
}
.card:hover { border-color: var(--hairline-strong); background: var(--surface-2); }

[hidden] { display: none !important; }
```

- [ ] **Step 2: Commit**

```bash
git add docs/css/base.css
git commit -m "feat(site): add base reset + typography + layout primitives"
```

### Task 5.3: `docs/css/sections.css`

**File:** `docs/css/sections.css`

- [ ] **Step 1: Write per-section styles**

```css
/* sections.css — per-section visual treatments */

/* ─── s-hero ─── */
.s-hero {
  padding-block: clamp(4rem, 12vw, 8rem) clamp(3rem, 10vw, 6rem);
  border-top: 0;
}
.s-hero .brand-mark {
  font-family: var(--font-mono);
  font-size: var(--fs-xs);
  letter-spacing: var(--tr-wide);
  text-transform: uppercase;
  color: var(--accent);
}
.s-hero h1 {
  margin-top: 1.5rem;
  letter-spacing: var(--tr-mega);
  max-width: 22ch;
}
.s-hero .lede {
  margin-top: 1.25rem;
  font-size: var(--fs-md);
  color: var(--ink-muted);
  max-width: 56ch;
}
.s-hero .actions {
  margin-top: 2rem;
  display: flex; gap: 0.75rem; flex-wrap: wrap;
}

/* ─── s-vibe ─── */
.s-vibe .pillars { margin-top: 2rem; }
.s-vibe .pillar h3 {
  font-size: var(--fs-xl);
  margin-bottom: 0.5rem;
}

/* ─── s-numbers ─── */
.s-numbers .stats {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 1rem;
  margin-top: 2rem;
}
.s-numbers .stat {
  border-top: 1px solid var(--hairline);
  padding-top: 1rem;
}
.s-numbers .stat-num {
  font-family: var(--font-display);
  font-size: var(--fs-2xl);
  letter-spacing: var(--tr-display);
}
.s-numbers .stat-label {
  font-family: var(--font-mono);
  font-size: var(--fs-xs);
  letter-spacing: var(--tr-wide);
  text-transform: uppercase;
  color: var(--ink-subtle);
  margin-top: 0.25rem;
}
@media (max-width: 768px) {
  .s-numbers .stats { grid-template-columns: repeat(2, 1fr); }
}

/* ─── s-methodology ─── */
.s-methodology .toc {
  list-style: none;
  padding: 0;
  margin: 2rem 0 0;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.5rem 2rem;
  font-family: var(--font-mono);
  font-size: var(--fs-sm);
}
.s-methodology .toc li {
  padding: 0.5rem 0;
  border-bottom: 1px solid var(--hairline);
}
.s-methodology .toc-num {
  color: var(--ink-subtle);
  margin-right: 0.75rem;
  letter-spacing: var(--tr-wide);
}
@media (max-width: 768px) {
  .s-methodology .toc { grid-template-columns: 1fr; }
}

/* ─── s-catalog ─── */
.s-catalog .cat-group { margin-top: 3rem; }
.s-catalog .cat-group h3 {
  font-size: var(--fs-xl);
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--hairline);
}
.s-catalog .pattern-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}
.s-catalog .pattern-card { padding: 1.25rem; }
.s-catalog .pattern-card h4 {
  font-size: var(--fs-md);
  margin-bottom: 0.25rem;
}
.s-catalog .pattern-card .desc {
  color: var(--ink-subtle);
  font-size: var(--fs-sm);
  margin-bottom: 0.75rem;
}
.s-catalog .pattern-card pre {
  font-size: var(--fs-xs);
  padding: 0.75rem;
  background: var(--surface-2);
}
.s-catalog .pattern-card .status {
  display: inline-block;
  font-family: var(--font-mono);
  font-size: var(--fs-xs);
  text-transform: uppercase;
  letter-spacing: var(--tr-wide);
  color: var(--ink-subtle);
  border: 1px solid var(--hairline);
  padding: 1px 6px;
  border-radius: var(--r-sm);
}
.s-catalog .pattern-card .status.filled { color: var(--accent); border-color: var(--accent); }
@media (max-width: 1024px) { .s-catalog .pattern-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 600px)  { .s-catalog .pattern-grid { grid-template-columns: 1fr; } }

/* ─── s-specimens ─── */
.s-specimens .specimen {
  border: 1px solid var(--hairline);
  border-radius: var(--r-lg);
  padding: 2rem;
  background: var(--surface-1);
  margin-top: 1.5rem;
}
.s-specimens .specimen h3 { margin-bottom: 0.75rem; }
.s-specimens .specimen pre { background: var(--canvas); }

/* ─── s-install ─── */
.s-install pre {
  margin-top: 1rem;
  font-size: var(--fs-sm);
}

/* ─── s-roadmap ─── */
.s-roadmap ul { padding-left: 1.25rem; color: var(--ink-muted); }
.s-roadmap li { margin: 0.5rem 0; }

/* ─── s-faq ─── */
.s-faq dl { margin: 0; }
.s-faq dt {
  font-family: var(--font-display);
  font-size: var(--fs-md);
  font-weight: 600;
  margin-top: 1.5rem;
}
.s-faq dd {
  margin: 0.5rem 0 0;
  color: var(--ink-muted);
}

/* ─── s-colophon ─── */
.s-colophon {
  border-top: 1px solid var(--hairline);
  padding-block: 3rem;
  font-size: var(--fs-sm);
  color: var(--ink-subtle);
}
.s-colophon .links { display: flex; gap: 1.5rem; flex-wrap: wrap; }

/* ─── Lang switcher ─── */
.lang-switch {
  position: fixed;
  top: 1rem; right: 1rem;
  background: var(--surface-1);
  border: 1px solid var(--hairline);
  border-radius: var(--r-sm);
  padding: 0.25rem 0.5rem;
  font-family: var(--font-mono);
  font-size: var(--fs-xs);
  color: var(--ink);
  z-index: 10;
}

/* RTL */
html[dir="rtl"] .s-hero h1 { letter-spacing: 0; }
```

- [ ] **Step 2: Commit**

```bash
git add docs/css/sections.css
git commit -m "feat(site): add per-section styles (10 sections, responsive)"
```

---

## Phase 6 — Site HTML

### Task 6.1: `docs/index.html` — head + sections

**File:** `docs/index.html`

- [ ] **Step 1: Write the document**

The document is long; assemble in this order:

1. `<head>` with: charset, viewport, theme-color (dark/light), title, description, author, keywords, canonical, 13 hreflang `<link rel="alternate">`, OpenGraph + Twitter, `link rel="alternate" type="text/plain"` to llms.txt + llms-full.txt, favicon + apple-touch-icon, font preconnect + stylesheets (Editorial New, Neue Montreal, Geist Mono — same providers status-line uses), then local css files (tokens → base → sections), then JSON-LD scripts (`SoftwareApplication`, `WebSite`, `FAQPage`, `BreadcrumbList`, `HowTo`, `TechArticle`).

2. `<body>` → `<main>` → 10 sections:

```html
<section class="section s-hero" aria-labelledby="hero-title">
  <div class="wrap">
    <span class="brand-mark" data-i18n="hero.brand">go-skills</span>
    <h1 id="hero-title" data-i18n="hero.title">Architectural patterns and methodology for Go services.</h1>
    <p class="lede" data-i18n="hero.lede">A curated knowledge base — 50+ design and concurrency patterns plus an 18-chapter service-architecture playbook — packaged as Claude Code skills, exposed via this site, and (soon) an MCP server.</p>
    <div class="actions">
      <a class="btn btn-primary" href="#catalog" data-i18n="hero.cta.catalog">Browse catalog</a>
      <a class="btn" href="#methodology" data-i18n="hero.cta.methodology">Read methodology</a>
      <a class="btn" href="https://github.com/amazopic/go-skills" data-i18n="hero.cta.github">View on GitHub →</a>
    </div>
  </div>
</section>

<section class="section s-vibe" id="vibe">
  <div class="wrap">
    <span class="eyebrow" data-i18n="vibe.eyebrow">Three pillars</span>
    <h2 data-i18n="vibe.title">Methodology · Patterns · Skills.</h2>
    <p class="lede" data-i18n="vibe.lede">A reference for senior Go engineers — opinionated, source-derived, no decoration. The methodology gives you a baseline service; the patterns sharpen the parts; the skills make it all addressable from Claude Code.</p>
    <div class="grid grid-3 pillars">
      <div class="card pillar">
        <h3 data-i18n="vibe.p1.title">Methodology</h3>
        <p data-i18n="vibe.p1.body">An 18-chapter playbook for building production Go services — directory layout, layered architecture, manual DI, configuration, retries, storage, transport, jobs, logging, validation, error handling, testing, build/deploy, default stack.</p>
      </div>
      <div class="card pillar">
        <h3 data-i18n="vibe.p2.title">Patterns</h3>
        <p data-i18n="vibe.p2.body">50+ classic and Go-flavored patterns across creational, structural, behavioral, concurrency, synchronization, messaging, stability, profiling, idiom, anti-patterns. Filled entries link to runnable Go examples.</p>
      </div>
      <div class="card pillar">
        <h3 data-i18n="vibe.p3.title">Skills</h3>
        <p data-i18n="vibe.p3.body">Each pattern and chapter is a Claude Code skill — a markdown file with a triggerable description. Drop the repo into your skills path and it surfaces the right knowledge for the task at hand.</p>
      </div>
    </div>
  </div>
</section>

<section class="section s-numbers" id="numbers">
  <div class="wrap">
    <span class="eyebrow" data-i18n="numbers.eyebrow">Inventory</span>
    <h2 data-i18n="numbers.title">By the numbers.</h2>
    <div class="stats">
      <div class="stat"><div class="stat-num">9</div><div class="stat-label" data-i18n="numbers.s1">Categories</div></div>
      <div class="stat"><div class="stat-num">50+</div><div class="stat-label" data-i18n="numbers.s2">Patterns</div></div>
      <div class="stat"><div class="stat-num">18</div><div class="stat-label" data-i18n="numbers.s3">Methodology chapters</div></div>
      <div class="stat"><div class="stat-num">23</div><div class="stat-label" data-i18n="numbers.s4">Runnable examples</div></div>
      <div class="stat"><div class="stat-num">13</div><div class="stat-label" data-i18n="numbers.s5">Languages</div></div>
    </div>
  </div>
</section>

<section class="section s-methodology" id="methodology">
  <div class="wrap">
    <span class="eyebrow" data-i18n="meth.eyebrow">Service Methodology</span>
    <h2 data-i18n="meth.title">A reference for building services in Go.</h2>
    <p class="lede" data-i18n="meth.lede">An anonymized, project-agnostic playbook. Read end-to-end as the canonical 930-line document, or pick a chapter-skill below.</p>
    <ol class="toc">
      <li><span class="toc-num">01</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/01-principles-and-layout.md">Principles &amp; Layout</a></li>
      <li><span class="toc-num">02</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/02-layered-architecture.md">Layered Architecture</a></li>
      <li><span class="toc-num">03</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/03-bootstrap-and-di.md">Bootstrap &amp; DI</a></li>
      <li><span class="toc-num">04</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/04-configuration.md">Configuration</a></li>
      <li><span class="toc-num">05</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/05-external-connections.md">External Connections</a></li>
      <li><span class="toc-num">06</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/06-storage.md">Storage</a></li>
      <li><span class="toc-num">07</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/07-http-transport.md">HTTP Transport</a></li>
      <li><span class="toc-num">08</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/08-background-jobs.md">Background Jobs</a></li>
      <li><span class="toc-num">09</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/09-logging.md">Logging</a></li>
      <li><span class="toc-num">10</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/10-validation.md">Validation</a></li>
      <li><span class="toc-num">11</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/11-error-handling.md">Error Handling</a></li>
      <li><span class="toc-num">12</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/12-testing.md">Testing</a></li>
      <li><span class="toc-num">13</span><a href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/13-build-and-deploy.md">Build &amp; Deploy</a></li>
    </ol>
    <p style="margin-top: 1.5rem;"><a class="btn" href="https://github.com/amazopic/go-skills/blob/main/skills/methodology/00-canonical-full.md" data-i18n="meth.cta">Read the full canonical document →</a></p>
  </div>
</section>

<section class="section s-catalog" id="catalog">
  <div class="wrap">
    <span class="eyebrow" data-i18n="cat.eyebrow">Catalog</span>
    <h2 data-i18n="cat.title">All patterns.</h2>
    <p class="lede" data-i18n="cat.lede">Filled cards have a body and (where applicable) a runnable example. Stubs are listed for completeness — body to be filled in a future iteration.</p>

    <!-- Concurrency first — Go's headline strength -->
    <div class="cat-group">
      <h3 data-i18n="cat.concurrency">Concurrency</h3>
      <div class="pattern-grid">
        <article class="card pattern-card">
          <span class="status filled">filled</span>
          <h4>Bounded Parallelism</h4>
          <p class="desc">Process N items with at most K workers concurrently.</p>
          <pre><code>sem := make(chan struct{}, K)
for _, item := range items {
    sem &lt;- struct{}{}
    go func(it Item) {
        defer func(){ &lt;-sem }()
        process(it)
    }(item)
}</code></pre>
        </article>
        <!-- Repeat .pattern-card for each concurrency pattern; use .status.filled or unstyled for stubs -->
      </div>
    </div>

    <!-- Add cat-groups for: Stability, Messaging, Synchronization, Idiom, Creational, Structural, Behavioral, Profiling, Anti-Patterns -->
    <!-- One card per pattern from PATTERNS.md. Filled patterns include a code excerpt; stubs omit <pre>. -->
  </div>
</section>

<section class="section s-specimens" id="specimens">
  <div class="wrap">
    <span class="eyebrow" data-i18n="spec.eyebrow">Specimens</span>
    <h2 data-i18n="spec.title">Three patterns, in full.</h2>

    <article class="specimen">
      <h3>Functional Options</h3>
      <p>The idiomatic Go way to pass optional configuration to a constructor.</p>
      <pre><code>type Server struct {
    addr    string
    timeout time.Duration
}

type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func NewServer(addr string, opts ...Option) *Server {
    s := &amp;Server{addr: addr, timeout: 30 * time.Second}
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// usage
srv := NewServer(":8080", WithTimeout(5 * time.Second))</code></pre>
    </article>

    <article class="specimen">
      <h3>Circuit Breaker</h3>
      <p>Fail fast for an unreliable downstream so it has time to recover.</p>
      <pre><code>// Sketch: open after 5 consecutive failures, half-open after 30s.
type Breaker struct {
    failures   int
    state      string // "closed" | "open" | "half-open"
    openedAt   time.Time
    threshold  int
    cooldown   time.Duration
    mu         sync.Mutex
}

func (b *Breaker) Call(fn func() error) error {
    b.mu.Lock()
    if b.state == "open" &amp;&amp; time.Since(b.openedAt) &gt; b.cooldown {
        b.state = "half-open"
    }
    if b.state == "open" {
        b.mu.Unlock()
        return errors.New("circuit open")
    }
    b.mu.Unlock()

    err := fn()
    b.mu.Lock(); defer b.mu.Unlock()
    if err != nil {
        b.failures++
        if b.failures &gt;= b.threshold { b.state = "open"; b.openedAt = time.Now() }
        return err
    }
    b.failures = 0
    b.state = "closed"
    return nil
}</code></pre>
    </article>

    <article class="specimen">
      <h3>Fan-Out</h3>
      <p>Distribute work from a single channel across N worker goroutines.</p>
      <pre><code>func fanOut(in &lt;-chan Job, n int) {
    var wg sync.WaitGroup
    wg.Add(n)
    for i := 0; i &lt; n; i++ {
        go func() {
            defer wg.Done()
            for j := range in {
                process(j)
            }
        }()
    }
    wg.Wait()
}</code></pre>
    </article>
  </div>
</section>

<section class="section s-install" id="install">
  <div class="wrap">
    <span class="eyebrow" data-i18n="inst.eyebrow">Install</span>
    <h2 data-i18n="inst.title">Use as Claude Code skills.</h2>
    <p class="lede" data-i18n="inst.lede">Drop the repo into your skills path. Claude Code will surface skills from this repo when their description matches your task.</p>

<pre><code>git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills</code></pre>

    <p data-i18n="inst.settings">Then add to <code>~/.claude/settings.json</code>:</p>

<pre><code>{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}</code></pre>

    <p data-i18n="inst.restart">Restart Claude Code.</p>

    <p style="margin-top: 2rem; color: var(--ink-subtle); font-size: var(--fs-sm);" data-i18n="inst.mcp">An MCP server that exposes the same content programmatically is on the roadmap — see below.</p>
  </div>
</section>

<section class="section s-roadmap" id="roadmap">
  <div class="wrap">
    <span class="eyebrow" data-i18n="road.eyebrow">Roadmap</span>
    <h2 data-i18n="road.title">What's next.</h2>
    <ul>
      <li data-i18n="road.l1">Fill the remaining ~33 stubbed pattern bodies — port from upstream Go pattern archives.</li>
      <li data-i18n="road.l2">MCP server: expose <code>tools.list_patterns</code>, <code>tools.get_pattern</code>, and methodology chapters as resources.</li>
      <li data-i18n="road.l3">Translate methodology and pattern bodies (currently English-only) into the 12 secondary languages.</li>
      <li data-i18n="road.l4">CI: markdown-link-check, <code>go test ./examples/...</code>, Lighthouse-CI, JSON-LD validation, README parity.</li>
      <li data-i18n="road.l5">Slash command <code>/go-skill &lt;pattern&gt;</code> to inject a pattern into Claude Code context on demand.</li>
    </ul>
  </div>
</section>

<section class="section s-faq" id="faq">
  <div class="wrap">
    <span class="eyebrow" data-i18n="faq.eyebrow">FAQ</span>
    <h2 data-i18n="faq.title">Frequently asked.</h2>
    <dl>
      <dt data-i18n="faq.q1">Why build this when go-patterns already exists on GitHub?</dt>
      <dd data-i18n="faq.a1">go-skills consolidates two upstream archives (theory + runnable code), adds a service-architecture methodology, and packages everything as Claude Code-compatible skills with rich frontmatter — so your editor surfaces the right pattern automatically.</dd>

      <dt data-i18n="faq.q2">Are the patterns idiomatic Go or GoF translations?</dt>
      <dd data-i18n="faq.a2">Both. Concurrency, messaging, stability, profiling, idiom — these are Go-native. Creational/structural/behavioral are GoF translated to Go with idiomatic adjustments (e.g., functional options replacing Builder where it makes sense).</dd>

      <dt data-i18n="faq.q3">Why are some entries marked as stubs?</dt>
      <dd data-i18n="faq.a3">The catalog mirrors a public taxonomy; not every node has a body yet. Stubs make the index complete and discoverable; bodies will be filled iteratively. See `PATTERNS.md` for status.</dd>

      <dt data-i18n="faq.q4">Will the examples run on my machine?</dt>
      <dd data-i18n="faq.a4">Yes — one shared <code>go.mod</code> at <code>examples/</code>. Run <code>cd examples &amp;&amp; go test ./...</code>. Go 1.21+.</dd>

      <dt data-i18n="faq.q5">When will the MCP server land?</dt>
      <dd data-i18n="faq.a5">A future iteration. The <code>mcp/</code> folder holds a roadmap-only README in v1; the protocol surface is sketched there.</dd>

      <dt data-i18n="faq.q6">Why Linear's design system?</dt>
      <dd data-i18n="faq.a6">Linear's near-black + single-accent aesthetic matches the content — dense, technical, no decoration. Editorial New + Geist Mono carry the typography.</dd>

      <dt data-i18n="faq.q7">License?</dt>
      <dd data-i18n="faq.a7">MIT. Use, fork, redistribute freely.</dd>

      <dt data-i18n="faq.q8">Does this support the methodology in non-English?</dt>
      <dd data-i18n="faq.a8">README is translated to 13 languages. Pattern bodies and methodology chapters are English-only in v1; localization is on the roadmap.</dd>
    </dl>
  </div>
</section>

<footer class="s-colophon">
  <div class="wrap">
    <p>go-skills · MIT · Yevgeniy Achin</p>
    <div class="links">
      <a href="https://github.com/amazopic/go-skills">GitHub</a>
      <a href="/llms.txt">llms.txt</a>
      <a href="/llms-full.txt">llms-full.txt</a>
      <a href="#">Top ↑</a>
    </div>
  </div>
</footer>
```

3. After `</main>` and before `</body>`: `<script type="module" src="js/main.js"></script>`.

- [ ] **Step 2: Open in a local server and verify**

```bash
cd /Users/mylive/project/go-skills-md/go-skills/docs
python3 -m http.server 8080
```

Open http://localhost:8080. Verify:
- Page loads, no console errors.
- All 10 sections render.
- Buttons are styled, hairlines visible, accent color #5e6ad2 visible on primary CTAs.
- At mobile width (DevTools: 375px), grids collapse to single column.

- [ ] **Step 3: Commit**

```bash
git add docs/index.html
git commit -m "feat(site): add index.html with 10 sections (hero, vibe, numbers, methodology, catalog, specimens, install, roadmap, faq, colophon)"
```

### Task 6.2: Fill the catalog with all patterns

**File:** `docs/index.html` — modify `<section class="section s-catalog">`

- [ ] **Step 1: For each of the 9 categories, add a `<div class="cat-group">` with a `<h3>` and `<div class="pattern-grid">`. For each pattern in `PATTERNS.md`, add an `<article class="card pattern-card">`.**

For filled patterns include a 3–6-line code excerpt in `<pre><code>`. For stubs omit the `<pre>` and use only the desc + status badge.

Code excerpt for each filled pattern (use these snippets):

| Pattern | Excerpt |
|---|---|
| Builder | `b := NewBuilder().WithName("x").WithSize(10).Build()` |
| Factory Method | `func NewWidget(kind string) Widget { switch kind { case "a": return &A{}; default: return &B{} } }` |
| Object Pool | `var pool = sync.Pool{New: func() any { return new(Buffer) }}` |
| Singleton | `var once sync.Once; var instance *T; func Get() *T { once.Do(func(){ instance = &T{} }); return instance }` |
| Decorator | `type Logger func(http.Handler) http.Handler` |
| Proxy | `func (p *Proxy) Get(k K) V { if v, ok := p.cache[k]; ok { return v }; v := p.real.Get(k); p.cache[k] = v; return v }` |
| Observer | `type Observer interface{ OnEvent(e Event) }` |
| Strategy | `type Sorter interface{ Sort([]int) []int }` |
| Bounded Parallelism | (already in template above) |
| Generator | `func gen(in []int) <-chan int { out := make(chan int); go func(){ defer close(out); for _, v := range in { out <- v } }(); return out }` |
| Parallelism | `var wg sync.WaitGroup; for _, x := range xs { wg.Add(1); go func(v int){ defer wg.Done(); work(v) }(x) }; wg.Wait()` |
| Fan-In | `func merge(cs ...<-chan int) <-chan int { out := make(chan int); for _, c := range cs { go func(c <-chan int){ for v := range c { out <- v } }(c) }; return out }` |
| Fan-Out | (already in template above) |
| Publish/Subscribe | `bus.Subscribe("topic", func(m Msg){ ... }); bus.Publish("topic", msg)` |
| Semaphore | `sem := make(chan struct{}, N); sem <- struct{}{}; defer func(){ <-sem }()` |
| Circuit Breaker | (already in template above) |
| Timing | `defer func(t time.Time){ log.Printf("took %v", time.Since(t)) }(time.Now())` |
| Functional Options | (already in template above) |

- [ ] **Step 2: Verify visual**

Reload http://localhost:8080. Verify the catalog renders with proper grouping; status badges show "filled" for the 18 patterns above and unstyled for the 33 stubs.

- [ ] **Step 3: Commit**

```bash
git add docs/index.html
git commit -m "feat(site): populate catalog with all 51 patterns (filled + stubs)"
```

### Task 6.3: SEO `<head>` — JSON-LD, hreflang, OG, Twitter, meta

**File:** `docs/index.html` — `<head>` section

- [ ] **Step 1: Insert SEO meta tags**

Insert before `</head>` (replacing whatever placeholders Phase 6.1 left):

```html
<title>go-skills — Architectural patterns and methodology for Go services</title>
<meta name="description" content="A Go knowledge base — 50+ design and concurrency patterns plus an 18-chapter service-architecture methodology — packaged as Claude Code skills, exposed via this site. 13 languages." />
<meta name="author" content="Yevgeniy Achin" />
<meta name="keywords" content="go, golang, design patterns, concurrency patterns, service architecture, methodology, claude code, skills, mcp, github pages" />
<link rel="canonical" href="https://amazopic.github.io/go-skills/" />

<link rel="alternate" hreflang="x-default" href="https://amazopic.github.io/go-skills/" />
<link rel="alternate" hreflang="en" href="https://amazopic.github.io/go-skills/?lang=en" />
<link rel="alternate" hreflang="ru" href="https://amazopic.github.io/go-skills/?lang=ru" />
<link rel="alternate" hreflang="fr" href="https://amazopic.github.io/go-skills/?lang=fr" />
<link rel="alternate" hreflang="de" href="https://amazopic.github.io/go-skills/?lang=de" />
<link rel="alternate" hreflang="uk" href="https://amazopic.github.io/go-skills/?lang=uk" />
<link rel="alternate" hreflang="sl" href="https://amazopic.github.io/go-skills/?lang=sl" />
<link rel="alternate" hreflang="it" href="https://amazopic.github.io/go-skills/?lang=it" />
<link rel="alternate" hreflang="es" href="https://amazopic.github.io/go-skills/?lang=es" />
<link rel="alternate" hreflang="zh" href="https://amazopic.github.io/go-skills/?lang=zh" />
<link rel="alternate" hreflang="ja" href="https://amazopic.github.io/go-skills/?lang=ja" />
<link rel="alternate" hreflang="ko" href="https://amazopic.github.io/go-skills/?lang=ko" />
<link rel="alternate" hreflang="ar" href="https://amazopic.github.io/go-skills/?lang=ar" />

<meta property="og:type" content="website" />
<meta property="og:title" content="go-skills — Architectural patterns and methodology for Go services" />
<meta property="og:description" content="50+ Go patterns + 18-chapter service-architecture methodology, packaged as Claude Code skills. Linear-styled site, 13 languages." />
<meta property="og:url" content="https://amazopic.github.io/go-skills/" />
<meta property="og:site_name" content="go-skills" />
<meta property="og:image" content="https://amazopic.github.io/go-skills/og-image.png" />
<meta property="og:image:width" content="1200" />
<meta property="og:image:height" content="630" />
<meta property="og:locale" content="en_US" />
<meta property="og:locale:alternate" content="ru_RU" />
<meta property="og:locale:alternate" content="fr_FR" />
<meta property="og:locale:alternate" content="de_DE" />
<meta property="og:locale:alternate" content="uk_UA" />
<meta property="og:locale:alternate" content="sl_SI" />
<meta property="og:locale:alternate" content="it_IT" />
<meta property="og:locale:alternate" content="es_ES" />
<meta property="og:locale:alternate" content="zh_CN" />
<meta property="og:locale:alternate" content="ja_JP" />
<meta property="og:locale:alternate" content="ko_KR" />
<meta property="og:locale:alternate" content="ar_SA" />

<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:title" content="go-skills — Architectural patterns and methodology for Go services" />
<meta name="twitter:description" content="50+ Go patterns + 18-chapter service-architecture methodology, packaged as Claude Code skills." />
<meta name="twitter:image" content="https://amazopic.github.io/go-skills/og-image.png" />

<link rel="alternate" type="text/plain" href="/go-skills/llms.txt" title="llms.txt — AI crawler index" />
<link rel="alternate" type="text/plain" href="/go-skills/llms-full.txt" title="llms-full.txt — AI ingestion corpus" />

<link rel="icon" type="image/svg+xml" href="favicon.svg" />
<link rel="apple-touch-icon" href="apple-touch-icon.png" />

<link rel="preconnect" href="https://api.fontshare.com" crossorigin />
<link rel="preconnect" href="https://fonts.googleapis.com" />
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
<link rel="stylesheet" href="https://api.fontshare.com/v2/css?f[]=editorial-new@200,300,400,500,700,800&f[]=neue-montreal@400,500,700&display=swap" />
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Geist+Mono:wght@300;400;500;600&display=swap" />

<link rel="stylesheet" href="css/tokens.css?v=1" />
<link rel="stylesheet" href="css/base.css?v=1" />
<link rel="stylesheet" href="css/sections.css?v=1" />
```

- [ ] **Step 2: Append JSON-LD scripts before `</head>`**

```html
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "SoftwareApplication",
  "name": "go-skills",
  "applicationCategory": "DeveloperApplication",
  "operatingSystem": "Cross-platform",
  "description": "A Go knowledge base — architectural patterns and a service-building methodology — packaged as Claude Code skills.",
  "softwareVersion": "0.1.0",
  "datePublished": "2026-05-10",
  "author": { "@type": "Person", "name": "Yevgeniy Achin", "email": "amazopic@gmail.com", "url": "https://github.com/amazopic" },
  "publisher": { "@type": "Person", "name": "Yevgeniy Achin" },
  "offers": { "@type": "Offer", "price": "0", "priceCurrency": "USD" },
  "license": "https://spdx.org/licenses/MIT.html",
  "url": "https://amazopic.github.io/go-skills/",
  "downloadUrl": "https://github.com/amazopic/go-skills",
  "codeRepository": "https://github.com/amazopic/go-skills",
  "inLanguage": ["en","ru","fr","de","uk","sl","it","es","zh","ja","ko","ar"],
  "keywords": "go, golang, design patterns, concurrency, methodology, claude code, skills, mcp",
  "featureList": [
    "9 pattern categories: creational, structural, behavioral, concurrency, synchronization, messaging, stability, profiling, idiom, anti-patterns",
    "50+ pattern entries — filled and stubs",
    "18-chapter service architecture methodology, split into 13 chapter-grouped Claude Code skills",
    "23 runnable Go pattern packages with tests under examples/",
    "Static GitHub Pages site with Linear design system",
    "Full SEO stack: hreflang × 13, JSON-LD (SoftwareApplication, WebSite, FAQPage, BreadcrumbList, HowTo, TechArticle), llms.txt"
  ]
}
</script>

<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "WebSite",
  "name": "go-skills",
  "url": "https://amazopic.github.io/go-skills/",
  "inLanguage": ["en","ru","fr","de","uk","sl","it","es","zh","ja","ko","ar"],
  "publisher": { "@type": "Person", "name": "Yevgeniy Achin" }
}
</script>

<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "FAQPage",
  "mainEntity": [
    { "@type": "Question", "name": "Why build this when go-patterns already exists on GitHub?", "acceptedAnswer": { "@type": "Answer", "text": "go-skills consolidates two upstream archives (theory + runnable code), adds a service-architecture methodology, and packages everything as Claude Code-compatible skills with rich frontmatter — so your editor surfaces the right pattern automatically." } },
    { "@type": "Question", "name": "Are the patterns idiomatic Go or GoF translations?", "acceptedAnswer": { "@type": "Answer", "text": "Both. Concurrency, messaging, stability, profiling, idiom — these are Go-native. Creational/structural/behavioral are GoF translated to Go with idiomatic adjustments." } },
    { "@type": "Question", "name": "Why are some entries marked as stubs?", "acceptedAnswer": { "@type": "Answer", "text": "The catalog mirrors a public taxonomy; not every node has a body yet. Stubs make the index complete and discoverable." } },
    { "@type": "Question", "name": "Will the examples run on my machine?", "acceptedAnswer": { "@type": "Answer", "text": "Yes — one shared go.mod at examples/. Run cd examples && go test ./... . Go 1.21+." } },
    { "@type": "Question", "name": "When will the MCP server land?", "acceptedAnswer": { "@type": "Answer", "text": "A future iteration. The mcp/ folder holds a roadmap-only README in v1; the protocol surface is sketched there." } },
    { "@type": "Question", "name": "Why Linear's design system?", "acceptedAnswer": { "@type": "Answer", "text": "Linear's near-black + single-accent aesthetic matches the content — dense, technical, no decoration." } },
    { "@type": "Question", "name": "License?", "acceptedAnswer": { "@type": "Answer", "text": "MIT. Use, fork, redistribute freely." } },
    { "@type": "Question", "name": "Does this support the methodology in non-English?", "acceptedAnswer": { "@type": "Answer", "text": "README is translated to 13 languages. Pattern bodies and methodology chapters are English-only in v1; localization is on the roadmap." } }
  ]
}
</script>

<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "BreadcrumbList",
  "itemListElement": [
    { "@type": "ListItem", "position": 1, "name": "Home",        "item": "https://amazopic.github.io/go-skills/" },
    { "@type": "ListItem", "position": 2, "name": "Methodology", "item": "https://amazopic.github.io/go-skills/#methodology" },
    { "@type": "ListItem", "position": 3, "name": "Catalog",     "item": "https://amazopic.github.io/go-skills/#catalog" },
    { "@type": "ListItem", "position": 4, "name": "Install",     "item": "https://amazopic.github.io/go-skills/#install" },
    { "@type": "ListItem", "position": 5, "name": "FAQ",         "item": "https://amazopic.github.io/go-skills/#faq" }
  ]
}
</script>

<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "HowTo",
  "name": "Install go-skills as Claude Code skills",
  "description": "Drop the go-skills repository into your Claude Code skills path so the editor surfaces relevant Go patterns and methodology chapters.",
  "totalTime": "PT2M",
  "tool": [{ "@type": "HowToTool", "name": "git" }, { "@type": "HowToTool", "name": "Claude Code" }],
  "step": [
    { "@type": "HowToStep", "name": "Clone",    "text": "git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills" },
    { "@type": "HowToStep", "name": "Configure","text": "Add the path to skillSources in ~/.claude/settings.json" },
    { "@type": "HowToStep", "name": "Restart",  "text": "Restart Claude Code." }
  ]
}
</script>

<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "TechArticle",
  "headline": "Go Service Architecture Methodology",
  "description": "An 18-chapter reference for building backend services in Go — directory layout, layer separation, DI, configuration, transport, storage, jobs, logging, testing, build/deploy.",
  "author": { "@type": "Person", "name": "Yevgeniy Achin" },
  "datePublished": "2026-05-10",
  "inLanguage": "en",
  "url": "https://amazopic.github.io/go-skills/#methodology",
  "articleSection": [
    "Principles", "Directory Layout", "Layered Architecture", "Bootstrap & DI",
    "Configuration", "External Connections", "Storage", "HTTP Transport",
    "Background Jobs", "Logging", "Validation", "Error Handling", "Testing",
    "Build & Deploy", "Default Stack"
  ]
}
</script>
```

- [ ] **Step 3: Validate JSON-LD**

Open https://search.google.com/test/rich-results, paste `https://localhost:8080` (or the rendered HTML directly).

Expected: All 6 schemas detected and valid (SoftwareApplication, WebSite, FAQPage, BreadcrumbList, HowTo, TechArticle).

- [ ] **Step 4: Commit**

```bash
git add docs/index.html
git commit -m "feat(site): add full SEO stack (hreflang × 13, JSON-LD × 6, OG, Twitter)"
```

---

## Phase 7 — Site JS

### Task 7.1: `docs/js/i18n.js`

**File:** `docs/js/i18n.js`

- [ ] **Step 1: Write the i18n module**

```javascript
// i18n.js — UI string switching across 13 locales
// Default: en. Other locales fall back to en for missing keys.
// Activate with ?lang=xx (also writes to localStorage for stickiness).

export const supportedLocales = [
  { code: 'en', label: 'English',     native: 'English'    },
  { code: 'ru', label: 'Russian',     native: 'Русский'    },
  { code: 'fr', label: 'French',      native: 'Français'   },
  { code: 'de', label: 'German',      native: 'Deutsch'    },
  { code: 'uk', label: 'Ukrainian',   native: 'Українська' },
  { code: 'sl', label: 'Slovenian',   native: 'Slovenščina'},
  { code: 'it', label: 'Italian',     native: 'Italiano'   },
  { code: 'es', label: 'Spanish',     native: 'Español'    },
  { code: 'zh', label: 'Chinese',     native: '中文'        },
  { code: 'ja', label: 'Japanese',    native: '日本語'      },
  { code: 'ko', label: 'Korean',      native: '한국어'      },
  { code: 'ar', label: 'Arabic',      native: 'العربية',     rtl: true },
];

export const defaultLocale = 'en';

export const messages = {
  en: {
    'meta.title':   'go-skills — Architectural patterns and methodology for Go services',
    'meta.description': 'A Go knowledge base — 50+ patterns + 18-chapter service architecture methodology — packaged as Claude Code skills.',
    'hero.brand':   'go-skills',
    'hero.title':   'Architectural patterns and methodology for Go services.',
    'hero.lede':    'A curated knowledge base — 50+ design and concurrency patterns plus an 18-chapter service-architecture playbook — packaged as Claude Code skills, exposed via this site, and (soon) an MCP server.',
    'hero.cta.catalog':     'Browse catalog',
    'hero.cta.methodology': 'Read methodology',
    'hero.cta.github':      'View on GitHub →',

    'vibe.eyebrow': 'Three pillars',
    'vibe.title':   'Methodology · Patterns · Skills.',
    'vibe.lede':    'A reference for senior Go engineers — opinionated, source-derived, no decoration.',
    'vibe.p1.title':'Methodology',
    'vibe.p1.body': 'An 18-chapter playbook for building production Go services.',
    'vibe.p2.title':'Patterns',
    'vibe.p2.body': '50+ classic and Go-flavored patterns across 9 categories.',
    'vibe.p3.title':'Skills',
    'vibe.p3.body': 'Each pattern and chapter is a Claude Code skill.',

    'numbers.eyebrow':'Inventory',
    'numbers.title':  'By the numbers.',
    'numbers.s1':'Categories',
    'numbers.s2':'Patterns',
    'numbers.s3':'Methodology chapters',
    'numbers.s4':'Runnable examples',
    'numbers.s5':'Languages',

    'meth.eyebrow': 'Service Methodology',
    'meth.title':   'A reference for building services in Go.',
    'meth.lede':    'An anonymized, project-agnostic playbook. Read end-to-end as the canonical 930-line document, or pick a chapter-skill below.',
    'meth.cta':     'Read the full canonical document →',

    'cat.eyebrow':  'Catalog',
    'cat.title':    'All patterns.',
    'cat.lede':     'Filled cards have a body and (where applicable) a runnable example. Stubs are listed for completeness.',
    'cat.concurrency': 'Concurrency',
    'cat.stability':   'Stability',
    'cat.messaging':   'Messaging',
    'cat.synchronization':'Synchronization',
    'cat.idiom':       'Idioms',
    'cat.creational':  'Creational',
    'cat.structural':  'Structural',
    'cat.behavioral':  'Behavioral',
    'cat.profiling':   'Profiling',
    'cat.antipatterns':'Anti-Patterns',

    'spec.eyebrow':'Specimens',
    'spec.title':  'Three patterns, in full.',

    'inst.eyebrow':'Install',
    'inst.title':  'Use as Claude Code skills.',
    'inst.lede':   'Drop the repo into your skills path. Claude Code will surface skills from this repo when their description matches your task.',
    'inst.settings':'Then add to ~/.claude/settings.json:',
    'inst.restart': 'Restart Claude Code.',
    'inst.mcp':     'An MCP server that exposes the same content programmatically is on the roadmap.',

    'road.eyebrow':'Roadmap',
    'road.title':  "What's next.",
    'road.l1':'Fill the remaining stubbed pattern bodies — port from upstream Go pattern archives.',
    'road.l2':'MCP server: tools.list_patterns, tools.get_pattern, methodology resources.',
    'road.l3':'Translate methodology and pattern bodies into the 12 secondary languages.',
    'road.l4':'CI: link checks, go tests, Lighthouse, JSON-LD validation.',
    'road.l5':'Slash command /go-skill <pattern> to inject a pattern on demand.',

    'faq.eyebrow':'FAQ',
    'faq.title':  'Frequently asked.',
    'faq.q1':'Why build this when go-patterns already exists on GitHub?',
    'faq.a1':'go-skills consolidates two upstream archives, adds a service methodology, and packages everything as Claude Code skills with rich frontmatter.',
    'faq.q2':'Are the patterns idiomatic Go or GoF translations?',
    'faq.a2':'Both. Concurrency, messaging, stability, profiling, idiom are Go-native. GoF patterns are translated with idiomatic adjustments.',
    'faq.q3':'Why are some entries marked as stubs?',
    'faq.a3':'The catalog mirrors a public taxonomy; not every node has a body yet.',
    'faq.q4':'Will the examples run on my machine?',
    'faq.a4':'Yes — one shared go.mod at examples/. Run: cd examples && go test ./...',
    'faq.q5':'When will the MCP server land?',
    'faq.a5':'A future iteration. mcp/ holds a roadmap-only README in v1.',
    'faq.q6':"Why Linear's design system?",
    'faq.a6':"Linear's near-black + single-accent aesthetic matches the content — dense, technical, no decoration.",
    'faq.q7':'License?',
    'faq.a7':'MIT.',
    'faq.q8':'Does this support the methodology in non-English?',
    'faq.a8':'README × 13 languages. Pattern and methodology bodies are English-only in v1.',
  },
  ru: {
    'meta.title': 'go-skills — Архитектурные паттерны и методология Go-сервисов',
    'meta.description': 'База знаний по Go — 50+ паттернов и 18-главная методология построения сервисов, упакованная как Claude Code skills.',
    'hero.brand': 'go-skills',
    'hero.title': 'Архитектурные паттерны и методология Go-сервисов.',
    'hero.lede':  'Курированная база знаний — 50+ паттернов проектирования и параллелизма + 18 глав методологии построения сервисов. Упакована как Claude Code skills, опубликована на этом сайте, в будущем — MCP-сервер.',
    'hero.cta.catalog': 'Каталог',
    'hero.cta.methodology': 'Методология',
    'hero.cta.github': 'GitHub →',
    'vibe.eyebrow': 'Три столпа',
    'vibe.title': 'Методология · Паттерны · Skills.',
    'numbers.title': 'В цифрах.',
    'numbers.s1':'Категорий', 'numbers.s2':'Паттернов', 'numbers.s3':'Глав методологии', 'numbers.s4':'Примеров', 'numbers.s5':'Языков',
    'cat.title':'Все паттерны.',
    'inst.title':'Установка как Claude Code skill.',
    'road.title':'Дальше.',
    'faq.title':'Часто спрашивают.',
  },
  // Other locales: minimal — only the keys that change UI label substance.
  // For unspecified keys we fall back to en automatically.
  fr: { 'hero.cta.catalog':'Catalogue', 'hero.cta.methodology':'Méthodologie', 'cat.title':'Tous les patterns.' },
  de: { 'hero.cta.catalog':'Katalog',   'hero.cta.methodology':'Methodologie', 'cat.title':'Alle Patterns.' },
  uk: { 'hero.cta.catalog':'Каталог',   'hero.cta.methodology':'Методологія',  'cat.title':'Всі патерни.' },
  sl: { 'hero.cta.catalog':'Katalog',   'hero.cta.methodology':'Metodologija', 'cat.title':'Vsi vzorci.' },
  it: { 'hero.cta.catalog':'Catalogo',  'hero.cta.methodology':'Metodologia',  'cat.title':'Tutti i pattern.' },
  es: { 'hero.cta.catalog':'Catálogo',  'hero.cta.methodology':'Metodología',  'cat.title':'Todos los patrones.' },
  zh: { 'hero.cta.catalog':'目录',       'hero.cta.methodology':'方法论',         'cat.title':'所有模式。' },
  ja: { 'hero.cta.catalog':'カタログ',    'hero.cta.methodology':'方法論',        'cat.title':'全パターン。' },
  ko: { 'hero.cta.catalog':'카탈로그',    'hero.cta.methodology':'방법론',        'cat.title':'모든 패턴.' },
  ar: { 'hero.cta.catalog':'الفهرس',     'hero.cta.methodology':'المنهجية',      'cat.title':'كل الأنماط.' },
};

export function getLocale() {
  const params = new URLSearchParams(window.location.search);
  const fromQuery = params.get('lang');
  if (fromQuery && messages[fromQuery]) return fromQuery;
  const fromStorage = localStorage.getItem('go-skills.lang');
  if (fromStorage && messages[fromStorage]) return fromStorage;
  const fromBrowser = (navigator.language || 'en').slice(0, 2);
  if (messages[fromBrowser]) return fromBrowser;
  return defaultLocale;
}

export function setLocale(code) {
  if (!messages[code]) return;
  localStorage.setItem('go-skills.lang', code);
  apply(code);
}

export function t(code, key) {
  return (messages[code] && messages[code][key]) || messages[defaultLocale][key] || key;
}

export function apply(code) {
  const locale = supportedLocales.find(l => l.code === code) || supportedLocales[0];
  document.documentElement.setAttribute('lang', code);
  document.documentElement.setAttribute('dir', locale.rtl ? 'rtl' : 'ltr');
  document.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.getAttribute('data-i18n');
    el.textContent = t(code, key);
  });
  document.querySelectorAll('[data-i18n-html]').forEach(el => {
    const key = el.getAttribute('data-i18n-html');
    el.innerHTML = t(code, key);
  });
  const titleKey = 'meta.title';
  if (messages[code] && messages[code][titleKey]) {
    document.title = messages[code][titleKey];
  }
}
```

- [ ] **Step 2: Commit**

```bash
git add docs/js/i18n.js
git commit -m "feat(site): add i18n module (13 locales, en fallback, RTL for ar)"
```

### Task 7.2: `docs/js/themes.js` and `docs/js/main.js`

**Files:**
- Create: `docs/js/themes.js`
- Create: `docs/js/main.js`

- [ ] **Step 1: `themes.js` — minimal dark/light toggle (default dark)**

```javascript
// themes.js — toggle between dark (default, Linear-aligned) and a light flip.
// Stored in localStorage as `go-skills.theme`.

export function getTheme() {
  return localStorage.getItem('go-skills.theme') || 'dark';
}

export function setTheme(name) {
  localStorage.setItem('go-skills.theme', name);
  document.documentElement.setAttribute('data-theme', name);
}

export function init() {
  const initial = getTheme();
  document.documentElement.setAttribute('data-theme', initial);
}
```

(Light theme variants of CSS tokens are deferred to a future iteration; setting `data-theme="light"` does nothing visible in v1. The file exists so the API is stable.)

- [ ] **Step 2: `main.js` — entrypoint**

```javascript
// main.js — entrypoint. Initializes i18n + theme.

import { apply, getLocale, setLocale, supportedLocales } from './i18n.js';
import * as themes from './themes.js';

function buildLangSwitcher() {
  const sel = document.createElement('select');
  sel.className = 'lang-switch';
  sel.setAttribute('aria-label', 'Language');
  for (const loc of supportedLocales) {
    const opt = document.createElement('option');
    opt.value = loc.code;
    opt.textContent = `${loc.code.toUpperCase()} · ${loc.native}`;
    sel.appendChild(opt);
  }
  sel.value = getLocale();
  sel.addEventListener('change', () => setLocale(sel.value));
  document.body.appendChild(sel);
}

function init() {
  themes.init();
  apply(getLocale());
  buildLangSwitcher();
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  init();
}
```

- [ ] **Step 3: Verify in browser**

Reload http://localhost:8080. Verify:
- A small `<select>` appears top-right with EN · English.
- Switching to RU updates: title, hero CTAs, etc., to Russian.
- Switching to AR sets `dir="rtl"` on `<html>`.
- localStorage holds `go-skills.lang`.

- [ ] **Step 4: Commit**

```bash
git add docs/js/themes.js docs/js/main.js
git commit -m "feat(site): add main entrypoint + theme stub + language switcher"
```

---

## Phase 8 — Static SEO assets

### Task 8.1: `docs/robots.txt`, `docs/sitemap.xml`, `docs/favicon.svg`, `docs/og-image.svg`

**Files:**
- Create: `docs/robots.txt`
- Create: `docs/sitemap.xml`
- Create: `docs/favicon.svg`
- Create: `docs/og-image.svg`
- Create: `docs/og-image.png` (rendered from svg)
- Create: `docs/apple-touch-icon.png`

- [ ] **Step 1: `robots.txt`**

```
User-agent: *
Allow: /

Sitemap: https://amazopic.github.io/go-skills/sitemap.xml
```

- [ ] **Step 2: `sitemap.xml`**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
        xmlns:xhtml="http://www.w3.org/1999/xhtml">
  <url>
    <loc>https://amazopic.github.io/go-skills/</loc>
    <lastmod>2026-05-10</lastmod>
    <changefreq>weekly</changefreq>
    <priority>1.0</priority>
    <xhtml:link rel="alternate" hreflang="en" href="https://amazopic.github.io/go-skills/?lang=en"/>
    <xhtml:link rel="alternate" hreflang="ru" href="https://amazopic.github.io/go-skills/?lang=ru"/>
    <xhtml:link rel="alternate" hreflang="fr" href="https://amazopic.github.io/go-skills/?lang=fr"/>
    <xhtml:link rel="alternate" hreflang="de" href="https://amazopic.github.io/go-skills/?lang=de"/>
    <xhtml:link rel="alternate" hreflang="uk" href="https://amazopic.github.io/go-skills/?lang=uk"/>
    <xhtml:link rel="alternate" hreflang="sl" href="https://amazopic.github.io/go-skills/?lang=sl"/>
    <xhtml:link rel="alternate" hreflang="it" href="https://amazopic.github.io/go-skills/?lang=it"/>
    <xhtml:link rel="alternate" hreflang="es" href="https://amazopic.github.io/go-skills/?lang=es"/>
    <xhtml:link rel="alternate" hreflang="zh" href="https://amazopic.github.io/go-skills/?lang=zh"/>
    <xhtml:link rel="alternate" hreflang="ja" href="https://amazopic.github.io/go-skills/?lang=ja"/>
    <xhtml:link rel="alternate" hreflang="ko" href="https://amazopic.github.io/go-skills/?lang=ko"/>
    <xhtml:link rel="alternate" hreflang="ar" href="https://amazopic.github.io/go-skills/?lang=ar"/>
    <xhtml:link rel="alternate" hreflang="x-default" href="https://amazopic.github.io/go-skills/"/>
  </url>
  <url><loc>https://amazopic.github.io/go-skills/#methodology</loc><priority>0.9</priority></url>
  <url><loc>https://amazopic.github.io/go-skills/#catalog</loc><priority>0.9</priority></url>
  <url><loc>https://amazopic.github.io/go-skills/#install</loc><priority>0.7</priority></url>
  <url><loc>https://amazopic.github.io/go-skills/#faq</loc><priority>0.5</priority></url>
</urlset>
```

- [ ] **Step 3: `favicon.svg`**

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32">
  <rect width="32" height="32" rx="6" fill="#010102"/>
  <text x="16" y="22" font-family="Geist Mono, monospace" font-size="18" font-weight="600"
        fill="#5e6ad2" text-anchor="middle">G</text>
</svg>
```

- [ ] **Step 4: `og-image.svg` (1200×630)**

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 630" width="1200" height="630">
  <rect width="1200" height="630" fill="#010102"/>
  <line x1="80" y1="120" x2="1120" y2="120" stroke="#23252a" stroke-width="1"/>
  <text x="80" y="100"  font-family="Geist Mono, monospace" font-size="22"
        fill="#5e6ad2" letter-spacing="2">GO-SKILLS · 0.1.0</text>
  <text x="80" y="320"  font-family="Editorial New, serif"  font-size="86"
        fill="#f7f8f8" letter-spacing="-2">Architectural patterns</text>
  <text x="80" y="410"  font-family="Editorial New, serif"  font-size="86"
        fill="#f7f8f8" letter-spacing="-2">and methodology for</text>
  <text x="80" y="500"  font-family="Editorial New, serif"  font-size="86"
        fill="#f7f8f8" letter-spacing="-2">Go services.</text>
  <text x="80" y="580"  font-family="Geist Mono, monospace" font-size="18"
        fill="#8a8f98">9 categories · 50+ patterns · 18 chapters · 13 languages</text>
</svg>
```

- [ ] **Step 5: Render `og-image.png` and `apple-touch-icon.png`**

If ImageMagick is available:

```bash
cd /Users/mylive/project/go-skills-md/go-skills/docs
convert -density 300 og-image.svg -background "#010102" -flatten og-image.png
convert -density 300 favicon.svg -resize 180x180 -background "#010102" -flatten apple-touch-icon.png
```

If ImageMagick is unavailable, use any SVG-to-PNG converter (e.g., `rsvg-convert`, online tool). Verify file sizes are reasonable (PNG < 200 KB).

- [ ] **Step 6: Commit**

```bash
git add docs/robots.txt docs/sitemap.xml docs/favicon.svg docs/og-image.svg docs/og-image.png docs/apple-touch-icon.png
git commit -m "feat(site): add robots, sitemap, favicon, og-image (svg + png)"
```

---

## Phase 9 — README × 13 languages + repo-level llms

### Task 9.1: Write base `README.md` (English)

**File:** `README.md`

- [ ] **Step 1: Write a concise English README**

```markdown
<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Architectural patterns and methodology for Go services — packaged as Claude Code skills.</p>

<p align="center">
  <a href="README.md">English</a> ·
  <a href="README.ru.md">Русский</a> ·
  <a href="README.fr.md">Français</a> ·
  <a href="README.de.md">Deutsch</a> ·
  <a href="README.uk.md">Українська</a> ·
  <a href="README.sl.md">Slovenščina</a> ·
  <a href="README.it.md">Italiano</a> ·
  <a href="README.es.md">Español</a> ·
  <a href="README.zh.md">中文</a> ·
  <a href="README.ja.md">日本語</a> ·
  <a href="README.ko.md">한국어</a> ·
  <a href="README.ar.md">العربية</a>
</p>

---

## What's inside

- **Methodology** — an 18-chapter playbook for building production Go services (directory layout, layered architecture, manual DI, configuration, retries, storage, transport, jobs, logging, validation, errors, testing, build, deploy). Read the [canonical full doc](skills/methodology/00-canonical-full.md) or pick a [chapter-skill](METHODOLOGY.md).
- **Patterns** — 50+ entries across 9 categories: creational, structural, behavioral, concurrency, synchronization, messaging, stability, profiling, idiom, anti-patterns. See [PATTERNS.md](PATTERNS.md).
- **Examples** — 23 runnable Go pattern packages with `_test.go`. One shared module under [`examples/`](examples/).
- **Site** — Linear-styled GitHub Pages: <https://amazopic.github.io/go-skills/>
- **MCP server** — placeholder for a future iteration. See [`mcp/`](mcp/).

## Install as Claude Code skills

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Add to `~/.claude/settings.json`:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Restart Claude Code.

## Run examples

```bash
cd examples
go test ./...
```

## License

MIT — see [LICENSE](LICENSE).

## Author

Yevgeniy Achin · <https://github.com/amazopic>
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add English README"
```

### Task 9.2: Translate to 12 secondary languages

**Files (create):**
- `README.ru.md`, `README.fr.md`, `README.de.md`, `README.uk.md`, `README.sl.md`, `README.it.md`, `README.es.md`, `README.zh.md`, `README.ja.md`, `README.ko.md`, `README.ar.md`

- [ ] **Step 1: Translate the README into each language, preserving structure**

For each language file:
- Same layout, same section order, same code blocks (untranslated).
- Translate prose: tagline, "What's inside" body, install instructions explanation, "Run examples" intro, license/author labels.
- Keep the 13-language language-bar identical at the top.
- For Arabic add `<div dir="rtl">` wrapping all prose (not code blocks).

This is content work; do it via web translation + manual review. Not detailed in this plan further — output one file per language.

- [ ] **Step 2: Verify all 13 README files render in GitHub**

```bash
ls /Users/mylive/project/go-skills-md/go-skills/README*.md
```

Expected: 13 files (`README.md` + `README.{ru,fr,de,uk,sl,it,es,zh,ja,ko,ar}.md`).

- [ ] **Step 3: Commit**

```bash
git add README.*.md
git commit -m "docs: add 12 README translations (ru, fr, de, uk, sl, it, es, zh, ja, ko, ar)"
```

---

## Phase 10 — Helper scripts

### Task 10.1: `make-llms.sh`

**File:** `make-llms.sh`

- [ ] **Step 1: Write the corpus builder**

```bash
#!/usr/bin/env bash
# make-llms.sh — concatenate skills/**/*.md into llms-full.txt
# Skips frontmatter (everything between the first two `---` lines per file).
# Output: llms-full.txt (committed) and docs/llms-full.txt (mirror).

set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
OUT="$ROOT/llms-full.txt"
DOCS_OUT="$ROOT/docs/llms-full.txt"

{
  echo "# go-skills — full corpus"
  echo
  echo "> All skills concatenated. For the indexed table of contents see llms.txt."
  echo
} > "$OUT"

find "$ROOT/skills" -type f -name '*.md' | sort | while read -r f; do
  rel="${f#$ROOT/}"
  echo "--- BEGIN $rel ---"  >> "$OUT"
  awk '
    BEGIN { in_fm=0; passed=0 }
    /^---$/ {
      if (!passed && !in_fm) { in_fm=1; next }
      if (in_fm) { in_fm=0; passed=1; next }
    }
    !in_fm && passed { print }
    !in_fm && !passed { print }   # fallback for files without frontmatter
  ' "$f" >> "$OUT"
  echo "--- END $rel ---" >> "$OUT"
  echo                    >> "$OUT"
done

cp "$OUT" "$DOCS_OUT"

echo "Wrote: $OUT and $DOCS_OUT"
wc -l "$OUT" "$DOCS_OUT"
```

- [ ] **Step 2: Make executable and run**

```bash
chmod +x /Users/mylive/project/go-skills-md/go-skills/make-llms.sh
/Users/mylive/project/go-skills-md/go-skills/make-llms.sh
```

Expected output: two paths printed, line counts both > 1000.

- [ ] **Step 3: Verify the result**

```bash
head -50 /Users/mylive/project/go-skills-md/go-skills/llms-full.txt
grep -c '^--- BEGIN' /Users/mylive/project/go-skills-md/go-skills/llms-full.txt
```

Expected: clean header followed by `--- BEGIN ...` markers; count = 65 (matches total skill files).

- [ ] **Step 4: Commit script + corpus**

```bash
git add make-llms.sh llms-full.txt docs/llms-full.txt
git commit -m "feat: add make-llms.sh corpus builder + first run output"
```

### Task 10.2: Mirror `llms.txt` into `docs/`

- [ ] **Step 1: Copy file**

```bash
cp /Users/mylive/project/go-skills-md/go-skills/llms.txt /Users/mylive/project/go-skills-md/go-skills/docs/llms.txt
```

- [ ] **Step 2: Commit**

```bash
git add docs/llms.txt
git commit -m "docs: mirror llms.txt under docs/ for site delivery"
```

---

## Phase 11 — Verification

### Task 11.1: Lighthouse audit

- [ ] **Step 1: Serve and audit**

```bash
cd /Users/mylive/project/go-skills-md/go-skills/docs
python3 -m http.server 8080 &
SERVER_PID=$!
sleep 1
npx --yes lighthouse http://localhost:8080 \
  --output=json --output-path=/tmp/go-skills-lighthouse.json \
  --chrome-flags="--headless" \
  --only-categories=performance,accessibility,best-practices,seo
kill $SERVER_PID
```

- [ ] **Step 2: Read scores**

```bash
node -e 'const r=require("/tmp/go-skills-lighthouse.json").categories;
for (const k of Object.keys(r)) console.log(k, Math.round(r[k].score*100));'
```

Expected: `performance ≥ 90`, `accessibility ≥ 95`, `best-practices ≥ 95`, `seo ≥ 95`. If any fall short, address top issues from the report (typically: image alt-text, color contrast, missing meta).

### Task 11.2: Markdown link check

- [ ] **Step 1: Run**

```bash
cd /Users/mylive/project/go-skills-md/go-skills
npx --yes markdown-link-check -q "**/*.md" 2>&1 | tail -50
```

Expected: all internal links resolve; external links to `github.com/amazopic/go-skills` may 404 until repo is pushed — flag and ignore those.

### Task 11.3: Go tests

- [ ] **Step 1: Run all examples**

```bash
cd /Users/mylive/project/go-skills-md/go-skills/examples
go vet ./...
go test ./... -short
```

Expected: `ok` for each package. Any failures → log under `examples/KNOWN-ISSUES.md`.

### Task 11.4: JSON-LD validation

- [ ] **Step 1: Manual check**

Open https://validator.schema.org/ and paste rendered HTML from `http://localhost:8080`. Expected: 6 schemas detected, no errors.

### Task 11.5: Final acceptance review

- [ ] **Step 1: Walk through spec §8 acceptance criteria** and confirm each:
  - All 13 README files present and render in GitHub (visual check via `gh` CLI: `gh repo view`).
  - `docs/index.html` opens with no console errors.
  - Lighthouse scores in target range.
  - All skill files have frontmatter.
  - `examples/...` `go test` passes.
  - All in-page anchors in `index.html` resolve.
  - Sitemap entries match anchors.
  - JSON-LD validates.
  - No broken markdown links.

- [ ] **Step 2: Commit any final tweaks and tag v0.1.0**

```bash
cd /Users/mylive/project/go-skills-md/go-skills
git tag v0.1.0
```

---

## Phase 12 — GitHub publish (manual — needs user action)

> **This phase requires user-side actions (creating GitHub repo, pushing, enabling Pages). Do not attempt without explicit user instruction.**

### Task 12.1: Push to GitHub

- [ ] **Step 1: User creates repo on GitHub** at `https://github.com/amazopic/go-skills` (empty, no auto-init).

- [ ] **Step 2: Add remote and push**

```bash
cd /Users/mylive/project/go-skills-md/go-skills
git remote add origin https://github.com/amazopic/go-skills.git
git push -u origin main
git push origin v0.1.0
```

### Task 12.2: Enable GitHub Pages

- [ ] **Step 1: User enables Pages**

In repo settings → Pages → Source: `main` branch, `/docs` folder. Save.

- [ ] **Step 2: Verify**

After ~1–2 min, visit `https://amazopic.github.io/go-skills/`. Confirm the site loads identically to local.

- [ ] **Step 3: Final commit (CNAME if custom domain — optional)**

If a custom domain is used later, add `docs/CNAME` and reconfigure DNS. Not in v1.

---

## Self-Review

**Spec coverage check** (each spec section → tasks that implement it):

- §0 Goal — covered by entire plan
- §1 Sources — Phase 1 (methodology), Phase 2 (patterns), Phase 3 (examples)
- §2 Repository Layout — Phase 0 + content phases
- §3 Site / 3.1 Tokens — Phase 5.1
- §3.2 Sections — Phases 6.1, 6.2
- §3.3 i18n — Phase 7.1, Phase 9
- §3.4 SEO — Phase 6.3, Phase 8
- §3.5 Build (none / scripts) — Phase 10
- §4.1 Methodology mapping — Phase 1.2
- §4.2 Patterns mapping — Phase 2.1, Phase 2.2
- §5 Data Flow — implicit in plan order; no separate task needed
- §6 Components / boundaries — encoded in file structure of Phase 0 + 5 + 6 + 7
- §7 Error handling table — addressed by Phase 10.1 (script frontmatter check), Phase 11.2 (link check), Phase 11.3 (go test)
- §8 Testing & Acceptance — Phase 11
- §9 Out of Scope — respected (no MCP impl, no full pattern translations, no SSG)
- §10 Risks — mitigations mostly require iteration discipline, not standalone tasks
- §11 Open Decisions — set as defaults in plan (LICENSE = MIT in Phase 0.1.4; URLs hard-coded as `amazopic.github.io/go-skills` throughout)
- §12 Future iterations — explicitly scoped out
- §13 Glossary — no implementation needed

**Placeholder scan:** No `TBD`, no `add appropriate error handling`, no `similar to Task N`, no "implement later" — verified.

**Type consistency:** Skill `name` slugs are consistent: `<category>-<pattern>` everywhere. URL paths use hyphens, not underscores. All file paths are absolute under `/Users/mylive/project/go-skills-md/go-skills/`. JSON-LD identifiers consistent across §6.3 and `llms.txt`. Site section IDs (`s-hero`, `s-vibe`, ...) match across CSS, HTML, and i18n keys.

**Decisions made (no ambiguity left):**
- Methodology split into 13 chapter-skills + 1 canonical full
- Patterns: 18 filled (from go-patterns-1) + 33 stubs = 51 catalog entries
- Examples: 23 runnable patterns from go-patterns-2 with one shared `go.mod`
- Site: 10 sections, all defined; 6 JSON-LD schemas
- i18n: full English dict + Russian substantive + 11 partial-stub locales
- License: MIT
- Repo URL: `github.com/amazopic/go-skills`, Pages: `amazopic.github.io/go-skills/`
