<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Patrones de arquitectura y metodología para servicios Go — empaquetados como Claude Code skills.</p>

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

## Contenido

- **Metodología** — una guía práctica de 18 capítulos para construir servicios Go en producción (estructura de directorios, arquitectura por capas, DI manual, configuración, reintentos, almacenamiento, transport, tareas, logging, validación, errores, testing, build, deploy). Lee el [documento canónico completo](skills/methodology/00-canonical-full.md) o elige un [chapter-skill](METHODOLOGY.md).
- **Patrones** — más de 52 entradas en 10 categorías: creacionales, estructurales, de comportamiento, concurrencia, sincronización, mensajería, estabilidad, profiling, idiomas, anti-patrones. Ver [PATTERNS.md](PATTERNS.md).
- **Ejemplos** — 52 paquetes Go de patrones ejecutables con `_test.go`. Un módulo compartido en [`examples/`](examples/).
- **Sitio** — GitHub Pages con estilo Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — marcador de posición para una iteración futura. Ver [`mcp/`](mcp/).

## Instalación como Claude Code skills

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Añade a `~/.claude/settings.json`:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Reinicia Claude Code.

## Ejecutar ejemplos

```bash
cd examples
go test ./...
```

## Licencia

MIT — ver [LICENSE](LICENSE).

## Autor

Yevgeniy Achin · <https://github.com/amazopic>
