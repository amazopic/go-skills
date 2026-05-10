---
name: go-validation
description: Use when writing input validators: pure functions returning error, table-driven tests, Example* for documentation; handler runs validation before service.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§11)
related:
  - skills/methodology/00-canonical-full.md
---

## 11. Validation

A dedicated `internal/<service>/validator/` package. Pure functions returning `error`. Tests use table-driven style and `Example*` for documentation.

```go
func ValidateEvent(e *dto.EventDTO) error {
    if e.EventID == "" { return errors.New("event_id is required") }
    if len(e.EventName) > 256 { return errors.New("event_name too long") }
    return nil
}
```

The handler runs validation **before** invoking the service. The service trusts its input.
