<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Pattern architetturali e metodologia per servizi Go — confezionati come Claude Code skills.</p>

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

## Contenuto

- **Metodologia** — una guida pratica in 18 capitoli per costruire servizi Go in produzione (struttura delle directory, architettura a strati, DI manuale, configurazione, retry, storage, transport, job, logging, validazione, errori, testing, build, deploy). Leggi il [documento canonico completo](skills/methodology/00-canonical-full.md) o scegli un [chapter-skill](METHODOLOGY.md).
- **Pattern** — oltre 50 voci in 9 categorie: creazionali, strutturali, comportamentali, concorrenza, sincronizzazione, messaggistica, stabilità, profiling, idiomi, anti-pattern. Vedi [PATTERNS.md](PATTERNS.md).
- **Esempi** — 23 pacchetti Go di pattern eseguibili con `_test.go`. Un modulo condiviso in [`examples/`](examples/).
- **Sito** — GitHub Pages in stile Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — segnaposto per una futura iterazione. Vedi [`mcp/`](mcp/).

## Installazione come Claude Code skills

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Aggiungi a `~/.claude/settings.json`:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Riavvia Claude Code.

## Eseguire gli esempi

```bash
cd examples
go test ./...
```

## Licenza

MIT — vedi [LICENSE](LICENSE).

## Autore

Yevgeniy Achin · <https://github.com/amazopic>
