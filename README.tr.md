<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Go servisleri için mimari pattern'ler ve metodoloji — Claude Code skills olarak paketlenmiş.</p>

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

> **go-skills'te yeni misiniz? [Project Assessment](skills/workflow/project-assessment.md) ile başlayın.** Onu mevcut Go projenize yönlendirin; her maddenin uygulanacak tam go-skills pattern'ine veya bölümüne bağlandığı bir olgunluk puanı ve önceliklendirilmiş bir yol haritası elde edin — yalnızca okuma yapar, kodunuzu asla değiştirmez.

## İçindekiler

- **Metodoloji** — production Go servisleri kurmak için 18 bölümlük bir el kitabı (dizin düzeni, katmanlı mimari, manuel DI, yapılandırma, yeniden denemeler, depolama, taşıma, işler, loglama, doğrulama, hatalar, test, derleme, dağıtım). [Kanonik tam belgeyi](skills/methodology/00-canonical-full.md) okuyun ya da bir [bölüm-skill](METHODOLOGY.md) seçin.
- **Pattern'ler** — 10 kategoride 52 madde: yaratımsal, yapısal, davranışsal, eşzamanlılık, senkronizasyon, mesajlaşma, kararlılık, profilleme, deyim, anti-pattern'ler. Bkz. [PATTERNS.md](PATTERNS.md).
- **Örnekler** — `_test.go` içeren 52 çalıştırılabilir Go pattern paketi. [`examples/`](examples/) altında tek bir ortak modül.
- **Workflows** — takım tabanlı çalıştırılabilir skill'ler: [Project Assessment](skills/workflow/project-assessment.md) (mevcut bir projeyi değerlendir), [Feature Development](skills/workflow/feature-development.md) (bir özellik geliştir), [Security Code Review](skills/workflow/security-review.md) (denetle + düzelt).
- **Site** — Linear tarzı GitHub Pages: <https://amazopic.github.io/go-skills/>
- **MCP sunucusu** — gelecek bir yineleme için yer tutucu. Bkz. [`mcp/`](mcp/).

## Claude Code skills olarak kurulum

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

`~/.claude/settings.json` dosyasına ekleyin:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Claude Code'u yeniden başlatın.

## Örnekleri çalıştırma

```bash
cd examples
go test ./...
```

## Lisans

MIT — bkz. [LICENSE](LICENSE).

## Yazar

Yevgeniy Achin · <https://github.com/amazopic>
