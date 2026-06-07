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
  <a href="README.ar.md">العربية</a> ·
  <a href="README.pt-BR.md">Português (BR)</a> ·
  <a href="README.tr.md">Türkçe</a> ·
  <a href="README.id.md">Bahasa Indonesia</a> ·
  <a href="README.vi.md">Tiếng Việt</a> ·
  <a href="README.hi.md">हिन्दी</a> ·
  <a href="README.zh-TW.md">繁體中文</a> ·
  <a href="README.pl.md">Polski</a> ·
  <a href="README.th.md">ไทย</a> ·
  <a href="README.he.md">עברית</a> ·
  <a href="README.bn.md">বাংলা</a> ·
  <a href="README.ur.md">اردو</a>
</p>

---

> **New to go-skills? Start with [Project Assessment](skills/workflow/project-assessment.md).** Point it at your existing Go project and get a maturity score plus a prioritized roadmap where every item links to the exact go-skills pattern or chapter to apply — read-only, it never changes your code.

## What's inside

- **Methodology** — an 18-chapter playbook for building production Go services (directory layout, layered architecture, manual DI, configuration, retries, storage, transport, jobs, logging, validation, errors, testing, build, deploy). Read the [canonical full doc](skills/methodology/00-canonical-full.md) or pick a [chapter-skill](METHODOLOGY.md).
- **Patterns** — 52 entries across 10 categories: creational, structural, behavioral, concurrency, synchronization, messaging, stability, profiling, idiom, anti-patterns. See [PATTERNS.md](PATTERNS.md).
- **Examples** — 52 runnable Go pattern packages with `_test.go`. One shared module under [`examples/`](examples/).
- **Workflows** — runnable team-based skills: [Project Assessment](skills/workflow/project-assessment.md) (assess an existing project), [Feature Development](skills/workflow/feature-development.md) (build a feature), [Security Code Review](skills/workflow/security-review.md) (audit + remediate).
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

## Versioning

This repo follows [Semantic Versioning](https://semver.org/). Current release: **v1.0.0**.
See [CHANGELOG.md](CHANGELOG.md) for the release history and what MAJOR/MINOR/PATCH mean
for the content surface (skills, examples, locales, site).

## License

MIT — see [LICENSE](LICENSE).

## Author

Yevgeniy Achin · <https://github.com/amazopic>
