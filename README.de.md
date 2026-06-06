<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Architekturmuster und Methodik für Go-Services — als Claude Code Skills paketiert.</p>

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
  <a href="README.pl.md">Polski</a>
</p>

---

> **Neu bei go-skills? Beginnen Sie mit [Project Assessment](skills/workflow/project-assessment.md).** Richten Sie es auf Ihr bestehendes Go-Projekt aus und erhalten Sie eine Reifegrad-Bewertung sowie eine priorisierte Roadmap, bei der jeder Punkt mit dem genauen go-skills-Muster oder -Kapitel verknüpft ist, das angewendet werden soll — schreibgeschützt, es ändert Ihren Code nie.

## Inhalt

- **Methodik** — ein 18-Kapitel-Leitfaden für den Aufbau von Go-Services in Produktion (Verzeichnisstruktur, Schichtenarchitektur, manuelles DI, Konfiguration, Wiederholungsversuche, Speicher, Transport, Jobs, Logging, Validierung, Fehler, Tests, Build, Deploy). Lesen Sie das [kanonische Volldokument](skills/methodology/00-canonical-full.md) oder wählen Sie einen [chapter-skill](METHODOLOGY.md).
- **Muster** — 52 Einträge in 10 Kategorien: erzeugende, strukturelle, verhaltensorientierte, Nebenläufigkeits-, Synchronisierungs-, Messaging-, Stabilitäts-, Profiling-, Idiom- und Anti-Muster. Siehe [PATTERNS.md](PATTERNS.md).
- **Beispiele** — 52 ausführbare Go-Musterpakete mit `_test.go`. Ein gemeinsames Modul unter [`examples/`](examples/).
- **Workflows** — ausführbare teambasierte Skills: [Project Assessment](skills/workflow/project-assessment.md) (ein bestehendes Projekt bewerten), [Feature Development](skills/workflow/feature-development.md) (ein Feature entwickeln), [Security Code Review](skills/workflow/security-review.md) (auditieren + beheben).
- **Website** — GitHub Pages im Linear-Stil: <https://amazopic.github.io/go-skills/>
- **MCP server** — Platzhalter für eine künftige Iteration. Siehe [`mcp/`](mcp/).

## Als Claude Code Skills installieren

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Fügen Sie zu `~/.claude/settings.json` hinzu:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Starten Sie Claude Code neu.

## Beispiele ausführen

```bash
cd examples
go test ./...
```

## Lizenz

MIT — siehe [LICENSE](LICENSE).

## Autor

Yevgeniy Achin · <https://github.com/amazopic>
