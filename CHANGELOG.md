# Changelog

All notable changes to **go-skills** are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Versioning policy

go-skills is a knowledge base (skills + methodology), a runnable examples module,
and a static site — not a single importable library. SemVer is interpreted for that
content surface:

- **MAJOR** (`x.0.0`) — breaking changes to how the content is consumed: skill
  frontmatter schema changes, renamed or removed patterns/skill files, a methodology
  restructure, the `examples/` module path changing, or site URL/anchor changes that
  break existing links or a `skillSources` setup.
- **MINOR** (`1.x.0`) — backward-compatible additions: new patterns, workflow skills,
  methodology chapters, runnable examples, locales, or site sections.
- **PATCH** (`1.0.x`) — backward-compatible fixes: example bug/race fixes, typo and
  translation corrections, content clarifications, SEO/accessibility/build tweaks.

> The `examples/` Go module is meant to be cloned and run, not imported via `go get`;
> repository tags version the content and site, not a published Go API.

## [1.0.0] - 2026-06-07

First tagged release — marks the complete v1 state of the project.

### Added

- **Methodology** — an 18-chapter playbook for building production Go services, plus
  the canonical full document (`skills/methodology/00-canonical-full.md`).
- **Patterns** — 52 senior-grade patterns across 10 categories (creational, structural,
  behavioral, concurrency, synchronization, messaging, stability, profiling, idiom,
  anti-patterns), each authored as a Claude Code skill with triggerable frontmatter.
  See [PATTERNS.md](PATTERNS.md).
- **Examples** — 52 runnable, race-safe Go example packages under one shared module
  (`examples/`), each with `_test.go`. All pass `gofmt -l`, `go vet ./...`, and
  `go test -race ./...`.
- **Workflow skills** — three runnable team-based skills:
  [Project Assessment](skills/workflow/project-assessment.md) (read-only maturity
  assessment of an existing Go project),
  [Feature Development](skills/workflow/feature-development.md), and
  [Security Code Review](skills/workflow/security-review.md).
- **Site** — a Linear-styled GitHub Pages site (<https://amazopic.github.io/go-skills/>)
  with full SEO (canonical, hreflang, Open Graph/Twitter, six JSON-LD blocks, sitemap,
  robots) and AI-crawler friendliness (`llms.txt`, `llms-full.txt`, an explicit crawler
  allow-list).
- **Internationalization** — 23 locales for the README and site UI, including RTL
  support (ar, he, ur). The site dictionary is lazy-loaded per locale, so a visitor
  downloads only the language they use.
- **CI** — GitHub Actions running gofmt, vet, build, race tests, coverage, and
  govulncheck on the examples module; reproducible locally via `examples/Makefile`.

### Notes

- Pattern and methodology *bodies* are English-only in v1; the README and site UI are
  translated to all 23 languages. Body translation is on the roadmap.
- The MCP server is a roadmap placeholder in v1 (`mcp/`).

[1.0.0]: https://github.com/amazopic/go-skills/releases/tag/v1.0.0
