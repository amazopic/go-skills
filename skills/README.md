# Skills

This folder packages the knowledge base as Claude Code-compatible skills. Every file (except this README) starts with YAML frontmatter declaring `name`, `description`, `category`, `sources`, and (where applicable) `example`.

## Install as Claude Code skills

```bash
# 1. Clone the repo
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills

# 2. Add to your Claude Code skills path (~/.claude/settings.json):
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}

# 3. Restart Claude Code.
```

Now Claude Code will surface skills from this repo when their `description` matches your task.

## Layout

- `workflow/` — team-based development process skills (e.g. Feature Development Flow), with a runnable orchestration under `../workflows/`.
- `methodology/` — 13 chapter-grouped skills + 1 canonical full doc.
- `creational/`, `structural/`, `behavioral/` — GoF + Go-flavored patterns.
- `concurrency/`, `synchronization/`, `messaging/` — Go's strengths.
- `stability/`, `profiling/`, `idiom/`, `anti-patterns/` — operational and Go-specific.

See `../PATTERNS.md` for the full catalog with status (filled / partial / stub).
