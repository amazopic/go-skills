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
  <a href="README.pl.md">Polski</a> ·
  <a href="README.th.md">ไทย</a> ·
  <a href="README.he.md">עברית</a> ·
  <a href="README.bn.md">বাংলা</a> ·
  <a href="README.ur.md">اردو</a>
</p>

---

> **Sei nuovo di go-skills? Inizia da [Project Assessment](skills/workflow/project-assessment.md).** Puntalo sul tuo progetto Go esistente e ottieni un punteggio di maturità più una roadmap prioritizzata in cui ogni voce rimanda all'esatto pattern o capitolo go-skills da applicare — sola lettura, non modifica mai il tuo codice.

## Contenuto

- **Metodologia** — una guida pratica in 18 capitoli per costruire servizi Go in produzione (struttura delle directory, architettura a strati, DI manuale, configurazione, retry, storage, transport, job, logging, validazione, errori, testing, build, deploy). Leggi il [documento canonico completo](skills/methodology/00-canonical-full.md) o scegli un [chapter-skill](METHODOLOGY.md).
- **Pattern** — oltre 52 voci in 10 categorie: creazionali, strutturali, comportamentali, concorrenza, sincronizzazione, messaggistica, stabilità, profiling, idiomi, anti-pattern. Vedi [PATTERNS.md](PATTERNS.md).
- **Esempi** — 52 pacchetti Go di pattern eseguibili con `_test.go`. Un modulo condiviso in [`examples/`](examples/).
- **Workflows** — skill eseguibili basate su team: [Project Assessment](skills/workflow/project-assessment.md) (valuta un progetto esistente), [Feature Development](skills/workflow/feature-development.md) (costruisci una feature), [Security Code Review](skills/workflow/security-review.md) (audit + correzione).
- **Sito** — GitHub Pages in stile Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — segnaposto per una futura iterazione. Vedi [`mcp/`](mcp/).

## Installazione come Claude Code skills

Il modo più semplice — incolla questo in Claude Code e configura tutto:

```text
Set up go-skills for me — handle everything, I'll just vibe.

1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills
2. Add "~/.claude/plugins/go-skills/skills" to skillSources in ~/.claude/settings.json
3. Reload settings and confirm the skills are live
4. Tell me in one line what you can now help with — patterns, methodology, workflows

Then ask what I'm building and start coding Go at senior level — idiomatic, race-safe, no rookie mistakes.
```

Oppure configuralo a mano:

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
