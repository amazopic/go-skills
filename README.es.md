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
  <a href="README.pl.md">Polski</a> ·
  <a href="README.th.md">ไทย</a> ·
  <a href="README.he.md">עברית</a> ·
  <a href="README.bn.md">বাংলা</a> ·
  <a href="README.ur.md">اردو</a>
</p>

---

> **¿Nuevo en go-skills? Empieza con [Project Assessment](skills/workflow/project-assessment.md).** Apúntalo a tu proyecto Go existente y obtén una puntuación de madurez más una hoja de ruta priorizada donde cada elemento enlaza al patrón o capítulo exacto de go-skills que debes aplicar — solo lectura, nunca cambia tu código.

## Contenido

- **Metodología** — una guía práctica de 18 capítulos para construir servicios Go en producción (estructura de directorios, arquitectura por capas, DI manual, configuración, reintentos, almacenamiento, transport, tareas, logging, validación, errores, testing, build, deploy). Lee el [documento canónico completo](skills/methodology/00-canonical-full.md) o elige un [chapter-skill](METHODOLOGY.md).
- **Patrones** — más de 52 entradas en 10 categorías: creacionales, estructurales, de comportamiento, concurrencia, sincronización, mensajería, estabilidad, profiling, idiomas, anti-patrones. Ver [PATTERNS.md](PATTERNS.md).
- **Ejemplos** — 52 paquetes Go de patrones ejecutables con `_test.go`. Un módulo compartido en [`examples/`](examples/).
- **Flujos de trabajo** — skills ejecutables basados en equipos: [Project Assessment](skills/workflow/project-assessment.md) (evalúa un proyecto existente), [Feature Development](skills/workflow/feature-development.md) (construye una funcionalidad), [Security Code Review](skills/workflow/security-review.md) (auditoría + corrección).
- **Sitio** — GitHub Pages con estilo Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — marcador de posición para una iteración futura. Ver [`mcp/`](mcp/).

## Instalación como Claude Code skills

Lo más fácil — pega esto en Claude Code y lo configura todo:

```text
Set up go-skills for me — handle everything, I'll just vibe.

1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills
2. Add "~/.claude/plugins/go-skills/skills" to skillSources in ~/.claude/settings.json
3. Reload settings and confirm the skills are live
4. Tell me in one line what you can now help with — patterns, methodology, workflows

Then ask what I'm building and start coding Go at senior level — idiomatic, race-safe, no rookie mistakes.
```

O configúralo a mano:

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
