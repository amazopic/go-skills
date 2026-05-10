---
name: go-testing
description: Use when writing unit/integration/example/race tests: table tests + t.Run, Example*, fakes over mocks, integration build tag, makefile targets test/test-short/test-cover/test-race/test-bench.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§13)
related:
  - skills/methodology/00-canonical-full.md
---

## 13. Testing

### makefile targets

```makefile
test:
	go test ./... -v

test-short:
	go test ./... -short

test-cover:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-race:
	go test ./... -race -short

test-bench:
	go test ./... -bench=. -benchmem
```

### Approach

- **Unit tests** live next to the code (`*_test.go`).
- **Table tests + `t.Run`** — the standard form.
- **`Example*`** — both documentation and quick output checks.
- **No heavy mocking frameworks** — interfaces plus hand-written fakes are enough.
- **Integration tests** — separate build tag (`//go:build integration`) and real dependencies via docker-compose.
