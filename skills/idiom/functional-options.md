---
name: idiom-functional-options
description: Functional Options — vary constructors with `func(*T)` option arguments. The idiomatic Go alternative to Builder for optional configuration.
category: idiom
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/idiom/functional-options.md
example: examples/idiom/functional-options/
---

# Functional Options

Functional options are a method of implementing clean/eloquent APIs in Go.
Options implemented as a function set the state of that option.

## Idiomatic shape (validating, error-returning)

Prefer an `Option` that can fail: it lets options validate their input and lets
the constructor reject a bad value at build time instead of producing a
half-built object. Defaults live in one place; later options override earlier
ones; failures wrap a package sentinel with `%w`. See the runnable example in
`examples/idiom/functional-options/`.

```go
package options

import (
	"errors"
	"fmt"
	"time"
)

// ErrInvalidOption is the sentinel wrapped by every validation failure.
var ErrInvalidOption = errors.New("options: invalid option")

// config is unexported; the only way to populate it is through Option funcs.
type config struct {
	timeout  time.Duration
	maxConns int
}

// Option mutates a config and reports whether the value is valid.
// Options that cannot fail return nil.
type Option func(*config) error

func WithTimeout(d time.Duration) Option {
	return func(c *config) error {
		if d <= 0 {
			return fmt.Errorf("%w: timeout must be > 0, got %s", ErrInvalidOption, d)
		}
		c.timeout = d
		return nil
	}
}

func WithMaxConns(n int) Option {
	return func(c *config) error {
		if n < 1 {
			return fmt.Errorf("%w: maxConns must be >= 1, got %d", ErrInvalidOption, n)
		}
		c.maxConns = n
		return nil
	}
}

type Server struct {
	addr string
	cfg  config
}

func New(addr string, opts ...Option) (*Server, error) {
	if addr == "" {
		return nil, fmt.Errorf("%w: addr must not be empty", ErrInvalidOption)
	}
	cfg := config{timeout: 30 * time.Second, maxConns: 100} // defaults in one place
	for i, opt := range opts {
		if opt == nil { // tolerate a nil option in the variadic slice
			continue
		}
		if err := opt(&cfg); err != nil {
			return nil, fmt.Errorf("options.New: applying option %d: %w", i, err)
		}
	}
	return &Server{addr: addr, cfg: cfg}, nil
}
```

```go
s, err := options.New("localhost:8080",
	options.WithTimeout(5*time.Second),
	options.WithMaxConns(50),
)
if errors.Is(err, options.ErrInvalidOption) { /* reject config */ }
```

## Gotchas

- Apply options in order over a single `config`; the last option for a field
  wins. Build the `config` per call — never share one mutable struct across
  constructions, or defaults leak between callers.
- Keep `config` unexported and expose accessors, so the resolved config is
  immutable after `New` returns.
- A non-validating option still returns `nil` (uniform `func(*config) error`
  signature) rather than mixing `func(*config)` and `func(*config) error`.
- Don't write `} else { ... }` after a branch that `return`s — drop the `else`.

## Implementation (non-validating variant)

The original side-effecting example below uses an exported `Options` struct and
options that cannot fail; it predates the validating shape above and is kept for
reference.

### Options

```go
package file

type Options struct {
	UID         int
	GID         int
	Flags       int
	Contents    string
	Permissions os.FileMode
}

type Option func(*Options)

func UID(userID int) Option {
	return func(args *Options) {
		args.UID = userID
	}
}

func GID(groupID int) Option {
	return func(args *Options) {
		args.GID = groupID
	}
}

func Contents(c string) Option {
	return func(args *Options) {
		args.Contents = c
	}
}

func Permissions(perms os.FileMode) Option {
	return func(args *Options) {
		args.Permissions = perms
	}
}
```

### Constructor

```go
package file

func New(filepath string, setters ...Option) error {
	// Default Options
	args := &Options{
		UID:         os.Getuid(),
		GID:         os.Getgid(),
		Contents:    "",
		Permissions: 0666,
		Flags:       os.O_CREATE | os.O_EXCL | os.O_WRONLY,
	}

	for _, setter := range setters {
		setter(args)
	}

	f, err := os.OpenFile(filepath, args.Flags, args.Permissions)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(args.Contents); err != nil {
		return err
	}

	return f.Chown(args.UID, args.GID)
}
```

## Usage

```go
emptyFile, err := file.New("/tmp/empty.txt")
if err != nil {
    panic(err)
}

fillerFile, err := file.New("/tmp/file.txt", file.UID(1000), file.Contents("Lorem Ipsum Dolor Amet"))
if err != nil {
    panic(err)
}
```
