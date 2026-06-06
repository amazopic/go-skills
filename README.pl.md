<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Wzorce architektoniczne i metodologia dla usług w Go — spakowane jako skills Claude Code.</p>

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

## Co znajdziesz w środku

- **Metodologia** — 18-rozdziałowy podręcznik budowy produkcyjnych usług w Go (układ katalogów, architektura warstwowa, ręczne DI, konfiguracja, ponawianie, przechowywanie, transport, zadania, logowanie, walidacja, błędy, testowanie, budowanie, wdrażanie). Przeczytaj [kanoniczny pełny dokument](skills/methodology/00-canonical-full.md) lub wybierz [rozdział-skill](METHODOLOGY.md).
- **Wzorce** — ponad 52 wpisów w 10 kategoriach: kreacyjne, strukturalne, behawioralne, współbieżność, synchronizacja, komunikaty, stabilność, profilowanie, idiomy, antywzorce. Zobacz [PATTERNS.md](PATTERNS.md).
- **Przykłady** — 52 uruchamialne pakiety wzorców w Go z `_test.go`. Jeden współdzielony moduł w [`examples/`](examples/).
- **Strona** — GitHub Pages w stylu Linear: <https://amazopic.github.io/go-skills/>
- **Serwer MCP** — placeholder dla przyszłej iteracji. Zobacz [`mcp/`](mcp/).

## Instalacja jako skills Claude Code

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Dodaj do `~/.claude/settings.json`:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Uruchom ponownie Claude Code.

## Uruchamianie przykładów

```bash
cd examples
go test ./...
```

## Licencja

MIT — zobacz [LICENSE](LICENSE).

## Autor

Yevgeniy Achin · <https://github.com/amazopic>
