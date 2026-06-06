<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Patterns d'architecture et méthodologie pour les services Go — packagés comme Claude Code skills.</p>

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

> **Nouveau sur go-skills ? Commencez par [Project Assessment](skills/workflow/project-assessment.md).** Pointez-le sur votre projet Go existant et obtenez un score de maturité ainsi qu'une feuille de route priorisée où chaque élément renvoie au pattern ou au chapitre go-skills exact à appliquer — en lecture seule, il ne modifie jamais votre code.

## Contenu

- **Méthodologie** — un guide pratique en 18 chapitres pour construire des services Go en production (organisation des répertoires, architecture en couches, DI manuel, configuration, tentatives, stockage, transport, tâches, journalisation, validation, erreurs, tests, compilation, déploiement). Lisez le [document canonique complet](skills/methodology/00-canonical-full.md) ou choisissez un [chapter-skill](METHODOLOGY.md).
- **Patterns** — plus de 52 entrées réparties en 10 catégories : créationnels, structurels, comportementaux, concurrence, synchronisation, messagerie, stabilité, profilage, idiomes, anti-patterns. Voir [PATTERNS.md](PATTERNS.md).
- **Exemples** — 52 packages Go de patterns exécutables avec `_test.go`. Un module partagé sous [`examples/`](examples/).
- **Workflows** — skills exécutables basés sur le travail d'équipe : [Project Assessment](skills/workflow/project-assessment.md) (évaluer un projet existant), [Feature Development](skills/workflow/feature-development.md) (construire une fonctionnalité), [Security Code Review](skills/workflow/security-review.md) (audit + remédiation).
- **Site** — GitHub Pages au style Linear : <https://amazopic.github.io/go-skills/>
- **MCP server** — ébauche pour une itération future. Voir [`mcp/`](mcp/).

## Installation comme Claude Code skills

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Ajoutez à `~/.claude/settings.json` :

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Redémarrez Claude Code.

## Lancer les exemples

```bash
cd examples
go test ./...
```

## Licence

MIT — voir [LICENSE](LICENSE).

## Auteur

Yevgeniy Achin · <https://github.com/amazopic>
