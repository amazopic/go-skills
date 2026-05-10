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
  <a href="README.ar.md">العربية</a>
</p>

---

## Inhalt

- **Methodik** — ein 18-Kapitel-Leitfaden für den Aufbau von Go-Services in Produktion (Verzeichnisstruktur, Schichtenarchitektur, manuelles DI, Konfiguration, Wiederholungsversuche, Speicher, Transport, Jobs, Logging, Validierung, Fehler, Tests, Build, Deploy). Lesen Sie das [kanonische Volldokument](skills/methodology/00-canonical-full.md) oder wählen Sie einen [chapter-skill](METHODOLOGY.md).
- **Muster** — 50+ Einträge in 9 Kategorien: erzeugende, strukturelle, verhaltensorientierte, Nebenläufigkeits-, Synchronisierungs-, Messaging-, Stabilitäts-, Profiling-, Idiom- und Anti-Muster. Siehe [PATTERNS.md](PATTERNS.md).
- **Beispiele** — 23 ausführbare Go-Musterpakete mit `_test.go`. Ein gemeinsames Modul unter [`examples/`](examples/).
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
