---
name: go-skills — knowledge base, skills repository, future MCP
description: Design for go-skills — a Go-language knowledge base (architectural patterns + service methodology) packaged as Claude Code skills, exposed via a static GitHub Pages site, with a future MCP server. Mirrors the structure and SEO/i18n stack of claude-customize/status-line; uses the Linear design system from awesome-design-md/design-md/linear.app; sources content from go-skills-md/go-old-pattern (golang-service-methodology.md + go-patterns-1 + go-patterns-2).
status: draft
date: 2026-05-10
---

# go-skills — Design Document

## 0. Goal

Build a self-contained repository at `/Users/mylive/project/go-skills-md/go-skills/` that serves three audiences:

1. **A Go developer reading the GitHub README** — finds a concise, multilingual entry point that links to the static site and the skills.
2. **A Claude Code user** — drops the repo's `skills/` into their config and gets the full Go pattern + methodology library as triggerable skills.
3. **Search engines and AI crawlers** — index a Linear-styled GitHub Pages site with full SEO (JSON-LD, hreflang, llms.txt, sitemap, robots).

A future MCP server will expose the same content programmatically; this iteration leaves a placeholder folder and a roadmap section on the site, no implementation.

## 1. Sources of Truth (read-only — never modified)

| Source | Role | Volume |
|---|---|---|
| `go-skills-md/go-old-pattern/golang-service-methodology.md` | Canonical service-architecture methodology, 18 chapters | 930 lines, 29 KB |
| `go-skills-md/go-old-pattern/go-patterns-1/` | Pattern theory, SUMMARY taxonomy (9 categories) | ~24 .md, partially stubbed |
| `go-skills-md/go-old-pattern/go-patterns-2/` | Runnable Go pattern code (README + .go + _test.go per pattern) | ~24 patterns |
| `claude-customize/status-line/` | Repository structure template (README × 13 langs, docs/, llms.txt, JSON-LD stack) | reference |
| `md-data/deign/awesome-design-md/design-md/linear.app/` | Visual design system tokens | DESIGN.md tokens |

**Merge rule.** Where a pattern exists in both sources, take theory from `go-patterns-1` and runnable code from `go-patterns-2`. Where only one side has it, use what's there.

## 2. Repository Layout

```
go-skills/
├── README.md                       # English (default)
├── README.{ru,fr,de,uk,sl,it,es,zh,ja,ko,ar}.md
├── PATTERNS.md                     # Catalog index (analog of status-line BLOCKS.md)
├── METHODOLOGY.md                  # Methodology summary + links to skills/methodology/*
├── LICENSE                         # MIT or Source-Available (decide before publish)
├── llms.txt                        # AI-crawler index, hand-curated
├── llms-full.txt                   # Concatenated corpus (script-built)
├── .gitignore
├── .nojekyll                       # GitHub Pages directive
├── make-llms.sh                    # builds llms-full.txt from skills/
│
├── skills/                         # Claude Code-style skills (frontmatter + body)
│   ├── README.md                   # how-to-use
│   ├── methodology/                # 13 chapter-grouped skills + 1 canonical doc
│   │   ├── 00-canonical-full.md
│   │   ├── 01-principles-and-layout.md         # §1 §2 §16 §18
│   │   ├── 02-layered-architecture.md          # §3 §17
│   │   ├── 03-bootstrap-and-di.md              # §4
│   │   ├── 04-configuration.md                 # §5
│   │   ├── 05-external-connections.md          # §6
│   │   ├── 06-storage.md                       # §7
│   │   ├── 07-http-transport.md                # §8
│   │   ├── 08-background-jobs.md               # §9
│   │   ├── 09-logging.md                       # §10
│   │   ├── 10-validation.md                    # §11
│   │   ├── 11-error-handling.md                # §12
│   │   ├── 12-testing.md                       # §13
│   │   └── 13-build-and-deploy.md              # §14 §15
│   │
│   ├── creational/   abstract-factory · builder · factory-method · object-pool · prototype · singleton
│   ├── structural/   adapter · bridge · composite · decorator · facade · flyweight · proxy
│   ├── behavioral/   chain-of-responsibility · command · iterator · mediator · memento · observer · state · strategy · template-method · visitor
│   ├── concurrency/  bounded-parallelism · generator · parallelism · (+ stubs for barrier · broadcast · coroutine · reactor · producer-consumer per SUMMARY)
│   ├── synchronization/ semaphore · (+ stubs for condition_variable · mutex · monitor · read_write_lock)
│   ├── messaging/    fan-in · fan-out · publish-subscribe · (+ stubs for futures-promises · push-pull)
│   ├── stability/    circuit-breaker · (+ stubs for bulkhead · deadline · fail-fast · handshaking · steady-state)
│   ├── profiling/    timing
│   ├── idiom/        functional-options · specification
│   └── anti-patterns/ cascading-failures (stub)
│
├── examples/                       # runnable .go from go-patterns-2
│   ├── creational/<pattern>/{pattern.go,pattern_test.go,README.md}
│   ├── structural/...
│   └── ...
│
├── commands/                       # placeholder for future Claude Code slash commands
│   └── README.md
│
├── mcp/                            # placeholder for future MCP server
│   └── README.md                   # design notes only
│
└── docs/                           # GitHub Pages site (static, Linear design)
    ├── .nojekyll
    ├── index.html
    ├── README.md
    ├── favicon.svg
    ├── apple-touch-icon.png
    ├── og-image.png
    ├── og-image.svg
    ├── robots.txt
    ├── sitemap.xml
    ├── llms.txt
    ├── llms-full.txt
    ├── css/
    │   ├── tokens.css              # Linear tokens (#010102, #5e6ad2, ...)
    │   ├── base.css                # reset, typography, layout primitives
    │   └── sections.css            # per-section styles (s-hero, s-catalog, ...)
    ├── js/
    │   ├── main.js                 # entrypoint: initialize sub-modules
    │   ├── i18n.js                 # 13-language UI string switching
    │   └── themes.js               # dark/light toggle (default dark)
    └── assets/
        └── specimens/              # code-specimen screenshots / SVGs
```

### Skills frontmatter contract

Every file under `skills/**` (except READMEs) starts with:

```yaml
---
name: <hyphenated-skill-name>
description: <one-line, used by Claude Code matcher — be specific about when to invoke>
category: <methodology | creational | structural | behavioral | concurrency | synchronization | messaging | stability | profiling | idiom | anti-patterns>
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/<file>
  - go-old-pattern/go-patterns-2/<folder>
example: examples/<category>/<pattern>/
---
```

Body is the original markdown content (kept verbatim where it stands; light edits where stubs exist).

## 3. The Site (docs/) — Linear Design

### 3.1 Visual tokens

```
canvas:           #010102
surface-1:        #0f1011    # cards
surface-2:        #141516    # nested
hairline:         #23252a
hairline-strong:  #34343a
ink:              #f7f8f8
ink-muted:        #d0d6e0
ink-subtle:       #8a8f98
accent:           #5e6ad2    # exactly one chromatic accent (focus, brand mark, key CTAs)
accent-hover:     #828fff
```

Typography:
- **Display / headings**: Editorial New + Neue Montreal (already loaded by status-line; closest fallback to Linear Display)
- **Body UI**: Neue Montreal 400/500
- **Code**: Geist Mono 300/400/500/600

Spacing scale: 4 / 8 / 12 / 16 / 24 / 32 / 48 / 64 / 96 (px). 8-pt grid.

### 3.2 Page sections

Section IDs mirror status-line's `s-*` naming so any shared CSS conventions transfer.

| # | id | Purpose | Notes |
|---|---|---|---|
| 1 | `s-hero` | Brand wordmark, tagline, CTAs (Catalog / Methodology / GitHub) | Single-column, generous whitespace |
| 2 | `s-vibe` | Manifesto: three pillars Methodology · Patterns · Skills (+ MCP soon) | One paragraph + 3 hairlined cards |
| 3 | `s-numbers` | Headline metrics: 9 categories · 50+ patterns · 18 methodology chapters · 13 languages · 24 runnable examples | Big-number row |
| 4 | `s-methodology` | Methodology preview: 18-chapter TOC + lead paragraph; CTA to full doc | Anchors to skill files |
| 5 | `s-catalog` | Pattern catalog grouped by 9 categories. Each card: name, 1-line summary, 3–6-line code excerpt | Dark hairline cards, monospace code |
| 6 | `s-specimens` | 2–3 featured patterns shown in full: Functional Options, Circuit Breaker, Fan-Out | "Specimens" lifted from status-line — same name, different content |
| 7 | `s-install` | How to install as Claude Code skills (clone + settings.json snippet) + future MCP teaser | Two terminal-styled blocks |
| 8 | `s-roadmap` | Public roadmap: MCP server, additional translations, anti-patterns expansion | Bullet list with statuses |
| 9 | `s-faq` | 8–10 FAQ Q/A; mirrored as JSON-LD `FAQPage` | Structured semantic |
| 10 | `s-colophon` | License, author, llms.txt links, GitHub link | Status-line uses same name |

### 3.3 i18n strategy

- README files: 13 hand-localized files at repo root (concise: title · pitch · features · install · links). Heavy content stays English-only in this iteration.
- Site UI strings: switchable via `?lang=xx` and `js/i18n.js`. Same approach status-line uses.
- Pattern/methodology bodies remain English-only for now; localizing them is on the roadmap.

### 3.4 SEO stack

- 13 `<link rel="alternate" hreflang>` entries + canonical
- OpenGraph + Twitter Card tags
- JSON-LD: `SoftwareApplication`, `WebSite`, `FAQPage`, `BreadcrumbList`, `HowTo`, plus `TechArticle` for the methodology
- `robots.txt` (allow all, link to sitemap), `sitemap.xml`, `og-image.png` (1200×630)
- AI-crawler hints: `<link rel="alternate" type="text/plain" href="/llms.txt">`

### 3.5 Build (none)

The site is fully static. No SSG. Two helper scripts only:

- `make-llms.sh` — concatenates `skills/**/*.md` (skipping frontmatter) into `llms-full.txt`. Runs on demand; output is committed.
- `make-sitemap.sh` — emits `sitemap.xml` from a hard-coded list of in-page anchors. Single source of truth.

## 4. Information Architecture — Mapping Sources to Skills

### 4.1 Methodology (18 → 13 skills)

| Skill | Methodology chapters | Trigger description (frontmatter) |
|---|---|---|
| `01-principles-and-layout` | §1 Principles · §2 Directory Layout · §16 Naming · §18 Checklist | When laying out a new Go service or auditing an existing one's structure |
| `02-layered-architecture` | §3 Layers (transport · service · fetcher · DTO) · §17 Cross-layer flow | When designing transport/service/data boundaries |
| `03-bootstrap-and-di` | §4 main.go · app.go · di.go | When wiring dependencies, signal handling, graceful shutdown |
| `04-configuration` | §5 Configuration | When adding env-driven config to a service |
| `05-external-connections` | §6 External connections (retry/timeout) | When writing a Connect() factory for redis/clickhouse/postgres |
| `06-storage` | §7 OLTP/OLAP/Cache/in-memory | When choosing a storage layer or batching writes |
| `07-http-transport` | §8 chi · middleware · rate-limit · errors · context | When building or auditing the HTTP layer |
| `08-background-jobs` | §9 Cron | When scheduling background tasks |
| `09-logging` | §10 Two-tier logging | When wiring loggers into services and bootstrap |
| `10-validation` | §11 Validation | When writing input validators |
| `11-error-handling` | §12 Error handling | When propagating, wrapping, and logging errors |
| `12-testing` | §13 Testing | When writing unit/integration/example/race tests |
| `13-build-and-deploy` | §14 makefile · §15 Default stack | When setting up build targets, CI/CD, and dependency choices |

`00-canonical-full.md` keeps the original 930-line document intact for end-to-end reading.

### 4.2 Patterns (≈ 51 skills, including stubs)

Categories follow `go-patterns-1/SUMMARY.md`:

- **Creational** (6): abstract-factory · builder · factory-method · object-pool · prototype · singleton
- **Structural** (7): adapter · bridge · composite · decorator · facade · flyweight · proxy
- **Behavioral** (10): chain-of-responsibility · command · iterator · mediator · memento · observer · state · strategy · template-method · visitor
- **Concurrency** (8): n-barrier · bounded-parallelism · broadcast · coroutines · generator · reactor · parallelism · producer-consumer
- **Synchronization** (5): condition-variable · mutex · monitor · read-write-lock · semaphore
- **Messaging** (5): fan-in · fan-out · futures-promises · publish-subscribe · push-pull
- **Stability** (6): bulkhead · circuit-breaker · deadline · fail-fast · handshaking · steady-state
- **Profiling** (1): timing
- **Idiom** (2): functional-options · specification
- **Anti-patterns** (1): cascading-failures (stub)

Where source content is missing, the skill file is a stub frontmatter + a "Status: stub — see SUMMARY taxonomy" body. Stubs are openly listed; not every skill needs to be filled in v1. Stretch: fill all stubs by porting from upstream go-patterns archives (out of scope for this iteration).

## 5. Data Flow

```
go-old-pattern/                                 (read-only source archive)
        │
        ├─ golang-service-methodology.md ───► skills/methodology/00-canonical-full.md (verbatim)
        │                                      │
        │                                      └─► skills/methodology/01..13-*.md (chapter-split, frontmatter added)
        │
        ├─ go-patterns-1/<cat>/<pattern>.md ─► skills/<cat>/<pattern>.md (theory + frontmatter)
        │
        └─ go-patterns-2/<Cat>/<Pattern>/ ────► examples/<cat>/<pattern>/ (.go + _test.go + README)

skills/**/*.md ──┬──► docs/index.html#s-catalog (curated excerpts)
                 ├──► PATTERNS.md (full catalog index)
                 ├──► llms-full.txt (concatenated by make-llms.sh)
                 └──► future: mcp/ (programmatic exposure)
```

## 6. Components and Their Boundaries

| Unit | What it does | Inputs | Outputs |
|---|---|---|---|
| `skills/` | Authoritative knowledge content | Hand-edited markdown with frontmatter | Read by Claude Code, by docs/ build, by llms-full builder |
| `docs/index.html` | Marketing + navigation surface | Reads css/, js/, references skills/ via links only (not embedded) | Rendered HTML |
| `docs/css/tokens.css` | Linear design tokens | — | CSS custom properties consumed by base.css/sections.css |
| `docs/js/i18n.js` | UI-string language switching | `?lang=xx` query param | Updates `data-i18n` attributes in DOM |
| `make-llms.sh` | Build llms-full.txt corpus | `skills/**/*.md` | `llms-full.txt` (committed) |
| `examples/` | Runnable Go code | imported from go-patterns-2 | `go test` passes per package |
| `README.{lang}.md` × 13 | Multilingual entry point | — | GitHub renders directly |

Each unit can be understood, edited, and tested in isolation. CSS layers don't know about JS. JS doesn't read skills/. Site doesn't depend on examples/ being compiled. Adding a new pattern requires touching exactly: `skills/<cat>/<name>.md`, optional `examples/<cat>/<name>/`, and the catalog list in `PATTERNS.md` + `docs/index.html#s-catalog`.

## 7. Error Handling

This is a static-content repo. Failure surfaces and their handling:

| Surface | Failure mode | Handling |
|---|---|---|
| `make-llms.sh` | A skill file lacks valid frontmatter | Script fails fast, prints offending file path |
| `docs/index.html` | i18n key missing for current language | `js/i18n.js` falls back to English (no UI breakage) |
| `examples/<pattern>/` | `go test ./...` fails | CI marks failed, blocks merge to main |
| Internal links in markdown | Pattern referenced but file absent | Pre-commit hook (or CI step) using `markdown-link-check` flags broken links |
| README × 13 inconsistency | One language drifts in feature list | Linter checks structural parity (same headers, same link count) — implementation deferred; manual audit in v1 |

## 8. Testing & Acceptance

### v1 acceptance criteria

- All 13 README files render in GitHub without errors and link to one another via the language-flag block (status-line convention).
- `docs/index.html` opens locally and on `https://amazopic.github.io/go-skills/` with no console errors.
- Lighthouse audit (mobile): SEO ≥ 95, Performance ≥ 90, Accessibility ≥ 95, Best Practices ≥ 95.
- Every skill file under `skills/**` has valid frontmatter; `make-llms.sh` produces a non-empty `llms-full.txt`.
- All `examples/**` packages: `go vet` and `go test ./...` pass.
- All in-page anchors in `index.html` resolve to existing section IDs.
- All hreflang URLs return 200 (just `?lang=xx` query strings on the same page).
- Sitemap entries match actual anchors and language URLs.
- JSON-LD validates via Google's Rich Results Test.
- No broken links in any markdown file (markdown-link-check pass).

### Tooling

- HTML: W3C validator (manual run)
- Lighthouse: `lighthouse https://localhost:8080 --view`
- Markdown links: `markdown-link-check '**/*.md'`
- Go: `make check` (vet + short tests + race) per `examples/`

## 9. Out of Scope (v1)

YAGNI line, explicitly excluded:

- MCP server implementation (placeholder folder + roadmap text only)
- Translating pattern bodies and methodology bodies into the 12 non-English languages
- Filling every stubbed pattern (anti-patterns, some concurrency/messaging/stability entries) — flag as "stub" and continue
- Build pipelines / SSG / Tailwind / React — vanilla CSS+JS only
- Pre-commit hooks (we run linters manually in v1; CI added later)
- Custom slash commands under `commands/` (placeholder folder only)
- Rich pattern visualizations (sequence diagrams etc.) beyond text + code

## 10. Risks and Mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| 13-language READMEs drift over time | Inconsistent metadata across languages | Keep README short and structurally identical; add a line-count parity check to v2 |
| Stubbed pattern files mislead users | Looks like noise | Each stub has explicit "Status: stub" banner and links to upstream sources |
| Linear's true display font is unavailable | Visual feels off | Editorial New + Neue Montreal as a near-equivalent (already proven on status-line) |
| GitHub Pages indexing of examples/ Go files | Unwanted in search | `robots.txt` allows only top-level + docs/; site links only to relevant anchors |
| llms-full.txt drifts from skills | Stale corpus served to AI crawlers | `make-llms.sh` runs as part of release checklist; v2 adds CI step |

## 11. Open Decisions (defaults set; flag if you disagree)

| Decision | Default |
|---|---|
| LICENSE | MIT (open, permissive). Status-line is Source-Available; for a knowledge base MIT is friendlier. |
| GitHub owner/repo | `amazopic/go-skills` |
| Pages URL | `https://amazopic.github.io/go-skills/` |
| Default branch | `main` |
| OG image style | Linear-styled dark hero with brand mark + tagline; reuses the same SVG-to-PNG flow status-line uses |

## 12. Future Iterations (post-v1)

1. **Fill stubs** — port remaining pattern bodies from public Go patterns archives.
2. **MCP server** — expose skills as MCP tools/resources. Schema TBD.
3. **Translate content** — paginate methodology chapters and pattern bodies into the 12 secondary languages.
4. **CI** — GitHub Actions: `markdown-link-check`, `go test ./examples/...`, Lighthouse-CI, JSON-LD validation, README parity check.
5. **Slash commands** — `/go-skill <pattern>` to inject a specific pattern's frontmatter+body into Claude Code context on demand.
6. **Diagrams** — text-based UML/sequence diagrams (mermaid, mscgen) for select stability/concurrency patterns.

## 13. Glossary

- **Skill** — a markdown file under `skills/` with frontmatter; consumable by Claude Code as a triggerable knowledge unit.
- **Pattern** — a specific Go-language design or concurrency pattern (e.g., Functional Options).
- **Methodology** — the 18-chapter document on Go service architecture; split into 13 chapter-skills + 1 canonical full doc.
- **Specimen** — a featured pattern shown in full on the site (copying status-line's `s-specimens` naming).
- **Canonical** — the original, unmodified content from `go-old-pattern`, preserved for fidelity.
