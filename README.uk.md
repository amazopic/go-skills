<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Архітектурні патерни та методологія для Go-сервісів — упаковані як Claude Code skills.</p>

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

> **Новачок у go-skills? Почніть із [Project Assessment](skills/workflow/project-assessment.md).** Спрямуйте її на ваш наявний Go-проєкт і отримайте оцінку зрілості разом із пріоритезованою дорожньою картою, де кожен пункт посилається на конкретний патерн або розділ go-skills, який слід застосувати — лише для читання, вона ніколи не змінює ваш код.

## Що всередині

- **Методологія** — керівництво з 18 розділів для побудови production Go-сервісів (структура директорій, шарова архітектура, ручний DI, конфігурація, повторні спроби, сховища, transport, завдання, логування, валідація, помилки, тестування, збірка, деплой). Читайте [повний канонічний документ](skills/methodology/00-canonical-full.md) або оберіть окремий [chapter-skill](METHODOLOGY.md).
- **Патерни** — 52 записів у 10 категоріях: породжуючі, структурні, поведінкові, конкурентні, синхронізації, обміну повідомленнями, відмовостійкості, профілювання, ідіоми, анти-патерни. Див. [PATTERNS.md](PATTERNS.md).
- **Приклади** — 52 запускувані пакети Go-патернів із `_test.go`. Один спільний модуль у [`examples/`](examples/).
- **Робочі процеси** — запускувані skills на основі команд: [Project Assessment](skills/workflow/project-assessment.md) (оцінити наявний проєкт), [Feature Development](skills/workflow/feature-development.md) (побудувати функцію), [Security Code Review](skills/workflow/security-review.md) (аудит + усунення вад).
- **Сайт** — GitHub Pages у стилі Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — заготовка для майбутньої ітерації. Див. [`mcp/`](mcp/).

## Встановлення як Claude Code skills

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Додайте до `~/.claude/settings.json`:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Перезапустіть Claude Code.

## Запуск прикладів

```bash
cd examples
go test ./...
```

## Ліцензія

MIT — див. [LICENSE](LICENSE).

## Автор

Yevgeniy Achin · <https://github.com/amazopic>
