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
  <a href="README.ar.md">العربية</a>
</p>

---

## Что внутри

- **Методология** — руководство из 18 глав по созданию production Go-сервисов (структура директорий, слоистая архитектура, ручной DI, конфигурация, повторные попытки, хранилища, transport, задачи, логирование, валидация, ошибки, тестирование, сборка, деплой). Читайте [полный канонический документ](skills/methodology/00-canonical-full.md) или выберите отдельный [chapter-skill](METHODOLOGY.md).
- **Паттерны** — 50+ записей в 9 категориях: порождающие, структурные, поведенческие, конкурентные, синхронизации, обмена сообщениями, отказоустойчивости, профилирования, идиомы, анти-паттерны. См. [PATTERNS.md](PATTERNS.md).
- **Примеры** — 23 запускаемых пакета Go-паттернов с `_test.go`. Один общий модуль в [`examples/`](examples/).
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
