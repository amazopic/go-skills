# MCP Server (roadmap)

This folder is reserved for an MCP (Model Context Protocol) server that will expose the `skills/` content as MCP tools and resources.

Planned shape:
- `tools.list_patterns` — returns the catalog with filters (category, status).
- `tools.get_pattern` — returns body + example code for a given pattern slug.
- `resources://methodology/<n>` — pinned methodology chapters as resources.

Status: not yet implemented. See spec `docs/superpowers/specs/2026-05-10-go-skills-design.md` §0 and §12 for context.
