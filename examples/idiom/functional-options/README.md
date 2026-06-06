# Functional Options

`New(addr, opts ...Option)` configures a `Server` from a variadic list of
`func(*config) error` options applied over a single set of defaults. It is the
idiomatic Go alternative to the Builder pattern and telescoping constructors.

## When to use

- A constructor has several optional parameters and you want call sites to read
  clearly (`New(addr, WithTimeout(2*time.Second), WithTLS(false))`).
- You want to add new options later without breaking existing callers.
- Options may need to validate their input and fail construction.

Skip it when there are only one or two required parameters — plain arguments are
simpler.

## API

```go
s, err := options.New("localhost:8080",
    options.WithTimeout(5*time.Second),
    options.WithMaxConns(50),
    options.WithTLS(false),
)
```

## Key properties

- Defaults live in one place (`defaults()`); options override them in order, so
  the last option for a field wins.
- An `Option` may return an error; `New` wraps it with `%w` and the failing
  option index, returns a nil `Server`, and never yields a half-built object.
- Validation failures wrap the package sentinel `ErrInvalidOption`, so callers
  can branch with `errors.Is(err, options.ErrInvalidOption)`.
- Options are pure functions over a per-call `config`, so the same option slice
  is reusable and safe to share across concurrent `New` calls.
- A `nil` option in the variadic slice is tolerated and skipped.

## Run

```bash
cd examples && go test -race ./idiom/functional-options/
```
