<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Padrões arquiteturais e metodologia para serviços Go — empacotados como skills do Claude Code.</p>

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

> **Novo no go-skills? Comece pelo [Project Assessment](skills/workflow/project-assessment.md).** Aponte-o para seu projeto Go existente e obtenha uma pontuação de maturidade mais um roteiro priorizado em que cada item vincula ao padrão ou capítulo exato do go-skills a aplicar — somente leitura, nunca altera seu código.

## O que tem dentro

- **Metodologia** — um manual de 18 capítulos para construir serviços Go em produção (estrutura de diretórios, arquitetura em camadas, DI manual, configuração, retentativas, armazenamento, transporte, jobs, logging, validação, erros, testes, build, deploy). Leia o [documento canônico completo](skills/methodology/00-canonical-full.md) ou escolha um [capítulo-skill](METHODOLOGY.md).
- **Padrões** — mais de 52 itens em 10 categorias: criacionais, estruturais, comportamentais, concorrência, sincronização, mensageria, estabilidade, profiling, idiomas, antipadrões. Veja [PATTERNS.md](PATTERNS.md).
- **Exemplos** — 52 pacotes Go de padrões executáveis com `_test.go`. Um único módulo compartilhado em [`examples/`](examples/).
- **Workflows** — skills executáveis baseadas em equipe: [Project Assessment](skills/workflow/project-assessment.md) (avalie um projeto existente), [Feature Development](skills/workflow/feature-development.md) (construa uma funcionalidade), [Security Code Review](skills/workflow/security-review.md) (audite + corrija).
- **Site** — GitHub Pages com estilo Linear: <https://amazopic.github.io/go-skills/>
- **Servidor MCP** — placeholder para uma iteração futura. Veja [`mcp/`](mcp/).

## Instalar como skills do Claude Code

O jeito mais fácil — cole isto no Claude Code e ele configura tudo:

```text
Set up go-skills for me — handle everything, I'll just vibe.

1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills
2. Add "~/.claude/plugins/go-skills/skills" to skillSources in ~/.claude/settings.json
3. Reload settings and confirm the skills are live
4. Tell me in one line what you can now help with — patterns, methodology, workflows

Then ask what I'm building and start coding Go at senior level — idiomatic, race-safe, no rookie mistakes.
```

Ou configure manualmente:

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Adicione ao `~/.claude/settings.json`:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Reinicie o Claude Code.

## Executar os exemplos

```bash
cd examples
go test ./...
```

## Licença

MIT — veja [LICENSE](LICENSE).

## Autor

Yevgeniy Achin · <https://github.com/amazopic>
