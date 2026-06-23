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

> **Novi v go-skills? Začnite s [Project Assessment](skills/workflow/project-assessment.md).** Usmerite ga na svoj obstoječi Go projekt in pridobite oceno zrelosti ter prednostno razvrščen načrt, kjer je vsaka postavka povezana z natančno go-skills vzorcem ali poglavjem, ki ga je treba uporabiti — samo za branje, nikoli ne spremeni vaše kode.

## Vsebina

- **Metodologija** — priročnik v 18 poglavjih za gradnjo production Go storitev (struktura imenikov, slojevita arhitektura, ročni DI, konfiguracija, ponovni poskusi, shranjevanje, transport, opravila, beleženje, validacija, napake, testiranje, gradnja, namestitev). Preberite [kanonični celoten dokument](skills/methodology/00-canonical-full.md) ali izberite posamezen [chapter-skill](METHODOLOGY.md).
- **Vzorci** — več kot 52 vnosov v 10 kategorijah: ustvarjalni, strukturni, vedenjski, vzporednost, sinhronizacija, sporočanje, stabilnost, profiliranje, idiomi, anti-vzorci. Glejte [PATTERNS.md](PATTERNS.md).
- **Primeri** — 52 izvršljivih paketov Go vzorcev z `_test.go`. En skupni modul v [`examples/`](examples/).
- **Workflows** — izvršljive skills, ki temeljijo na ekipi: [Project Assessment](skills/workflow/project-assessment.md) (ocenite obstoječ projekt), [Feature Development](skills/workflow/feature-development.md) (zgradite funkcionalnost), [Security Code Review](skills/workflow/security-review.md) (revizija + odprava).
- **Spletna stran** — GitHub Pages v slogu Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — ogrodišče za prihodnjo iteracijo. Glejte [`mcp/`](mcp/).

## Namestitev kot Claude Code skills

Najlažje — prilepi to v Claude Code in vse se uredi samo:

```text
Set up go-skills for me — handle everything, I'll just vibe.

1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills
2. Add "~/.claude/plugins/go-skills/skills" to skillSources in ~/.claude/settings.json
3. Reload settings and confirm the skills are live
4. Tell me in one line what you can now help with — patterns, methodology, workflows

Then ask what I'm building and start coding Go at senior level — idiomatic, race-safe, no rookie mistakes.
```

Ali pa nastavi ročno:

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
