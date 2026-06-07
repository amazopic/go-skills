<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Архитектурные паттерны и методология для Go-сервисов — упакованы как Claude Code skills.</p>

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

> **Впервые в go-skills? Начните с [Project Assessment](skills/workflow/project-assessment.md).** Направьте его на свой существующий Go-проект и получите оценку зрелости плюс приоритизированную дорожную карту, где каждый пункт связан с конкретным паттерном или главой go-skills для применения — только чтение, ваш код никогда не изменяется.

## Что внутри

- **Методология** — руководство из 18 глав по созданию production Go-сервисов (структура директорий, слоистая архитектура, ручной DI, конфигурация, повторные попытки, хранилища, transport, задачи, логирование, валидация, ошибки, тестирование, сборка, деплой). Читайте [полный канонический документ](skills/methodology/00-canonical-full.md) или выберите отдельный [chapter-skill](METHODOLOGY.md).
- **Паттерны** — 52 записей в 10 категориях: порождающие, структурные, поведенческие, конкурентные, синхронизации, обмена сообщениями, отказоустойчивости, профилирования, идиомы, анти-паттерны. См. [PATTERNS.md](PATTERNS.md).
- **Примеры** — 52 запускаемых пакета Go-паттернов с `_test.go`. Один общий модуль в [`examples/`](examples/).
- **Workflows** — запускаемые командные skills: [Project Assessment](skills/workflow/project-assessment.md) (оценить существующий проект), [Feature Development](skills/workflow/feature-development.md) (разработать функциональность), [Security Code Review](skills/workflow/security-review.md) (аудит и устранение проблем).
- **Сайт** — GitHub Pages в стиле Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — заглушка для будущей итерации. См. [`mcp/`](mcp/).

## Установка как Claude Code skills

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Добавьте в `~/.claude/settings.json`:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Перезапустите Claude Code.

## Запуск примеров

```bash
cd examples
go test ./...
```

## Лицензия

MIT — см. [LICENSE](LICENSE).

## Автор

Yevgeniy Achin · <https://github.com/amazopic>
