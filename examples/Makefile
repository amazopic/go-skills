# Examples module — the test/quality targets the methodology (chapter 13) prescribes.
# Run from examples/. `make check` is what CI gates on.

.PHONY: test test-race test-cover test-bench vet fmt fmt-check vuln check

test: ## run tests
	go test ./...

test-race: ## run tests with the race detector
	go test -race ./...

test-cover: ## run tests with coverage
	go test -cover ./...

test-bench: ## run benchmarks only
	go test -bench=. -benchmem -run='^$$' ./...

vet: ## go vet
	go vet ./...

fmt: ## format in place
	gofmt -w .

fmt-check: ## fail if anything is unformatted
	@drift="$$(gofmt -l .)"; if [ -n "$$drift" ]; then echo "gofmt drift:"; echo "$$drift"; exit 1; fi

vuln: ## scan for known vulnerabilities (needs govulncheck on PATH)
	govulncheck ./...

check: fmt-check vet test-race ## the full local gate (matches CI)
