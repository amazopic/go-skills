# Known issues — examples

None. All 52 packages pass `gofmt -l` (clean), `go vet ./...`, and `go test -race ./...` — now
enforced on every push by CI (`.github/workflows/ci.yml`) and reproducible locally with `make check`.

## Resolved

A Project Assessment pass found and fixed several latent defects that the old single-threaded tests
did not surface (each fix ships with a regression/concurrent test):

- `behavioral/mediator` — exported identifier used Cyrillic homoglyphs (`СonnectСolleagues`, U+0421); renamed to ASCII `ConnectColleagues`.
- `structural/bridge` — public method typo `Rase()` → `Race()`.
- `creational/prototype` — package doc named the wrong pattern ("Singleton"); corrected to Prototype.
- `structural/flyweight` — unsynchronized shared interning map; guarded with a mutex + a concurrent test.
- `behavioral/observer` — unsynchronized observer slice and 0% coverage; mutex added + a real test suite (95%+).
- `creational/object-pool` & `messaging/publish-subscribe` — send-on-closed-channel TOCTOU in `Put`/`Publish`; fixed + stress tests that race `Close`.
