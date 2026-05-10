---
name: go-error-handling
description: Use when propagating, wrapping, and logging errors: error always last, wrap with %w, preserve context.Canceled / DeadlineExceeded, map to common.ErrorResponse for HTTP.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§12)
related:
  - skills/methodology/00-canonical-full.md
---

## 12. Error handling

- Return: `error` is always the last value.
- Wrap: `fmt.Errorf("context: %w", err)` with a meaningful prefix.
- Log: once, close to where the error stops being propagated.
- HTTP: map through `common.NewErrorResponse(msg, code)` + `render.Render(w, r, ...)`.
- Preserve `context.Canceled` / `context.DeadlineExceeded` — check with `errors.Is` where it matters.
