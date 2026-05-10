<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Arhitekturni vzorci in metodologija za Go storitve — pakirani kot Claude Code skills.</p>

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

## Vsebina

- **Metodologija** — priročnik v 18 poglavjih za gradnjo production Go storitev (struktura imenikov, slojevita arhitektura, ročni DI, konfiguracija, ponovni poskusi, shranjevanje, transport, opravila, beleženje, validacija, napake, testiranje, gradnja, namestitev). Preberite [kanonični celoten dokument](skills/methodology/00-canonical-full.md) ali izberite posamezen [chapter-skill](METHODOLOGY.md).
- **Vzorci** — več kot 50 vnosov v 9 kategorijah: ustvarjalni, strukturni, vedenjski, vzporednost, sinhronizacija, sporočanje, stabilnost, profiliranje, idiomi, anti-vzorci. Glejte [PATTERNS.md](PATTERNS.md).
- **Primeri** — 23 izvršljivih paketov Go vzorcev z `_test.go`. En skupni modul v [`examples/`](examples/).
- **Spletna stran** — GitHub Pages v slogu Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — ogrodišče za prihodnjo iteracijo. Glejte [`mcp/`](mcp/).

## Namestitev kot Claude Code skills

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Dodajte v `~/.claude/settings.json`:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Znova zaženite Claude Code.

## Zagon primerov

```bash
cd examples
go test ./...
```

## Licenca

MIT — glejte [LICENSE](LICENSE).

## Avtor

Yevgeniy Achin · <https://github.com/amazopic>
