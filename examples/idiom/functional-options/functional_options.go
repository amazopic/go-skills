// Package options demonstrates the Functional Options idiom: a constructor
// accepts a variadic list of func(*config) values, each of which mutates a
// private configuration struct. It is the idiomatic Go alternative to the
// Builder pattern and to telescoping constructors for optional configuration.
//
// Why it wins over the alternatives:
//
//   - Backward compatible: adding a new option never changes existing call sites.
//   - Self-documenting: New(addr, WithTimeout(2*time.Second)) reads at the call site.
//   - Defaults live in one place (the constructor), not scattered across overloads.
//   - Options can validate and fail; the constructor surfaces the error with %w.
//
// This example shows the validating variant: an Option may return an error so a
// bad value is rejected at construction time rather than producing a half-built
// object. Options that cannot fail simply return nil.
package options

import (
	"errors"
	"fmt"
	"time"
)

// ErrInvalidOption is the sentinel wrapped by every option-validation failure.
// Callers can test for it with errors.Is(err, options.ErrInvalidOption).
var ErrInvalidOption = errors.New("options: invalid option")

// config holds the resolved configuration for a Server. It is unexported: the
// only way to populate it is through Option functions, which keeps the surface
// minimal and lets us evolve fields without breaking callers.
type config struct {
	timeout    time.Duration
	maxConns   int
	retries    int
	tls        bool
	loggerName string
}

// Option mutates a config and reports whether the requested value is valid.
// Returning a non-nil error aborts construction in New. Options that never fail
// return nil.
type Option func(*config) error

// defaults returns the baseline configuration applied before any options run.
// Centralising defaults here is the whole point of the pattern.
func defaults() config {
	return config{
		timeout:    30 * time.Second,
		maxConns:   100,
		retries:    3,
		tls:        true,
		loggerName: "default",
	}
}

// WithTimeout sets the request timeout. A non-positive duration is rejected.
func WithTimeout(d time.Duration) Option {
	return func(c *config) error {
		if d <= 0 {
			return fmt.Errorf("%w: timeout must be > 0, got %s", ErrInvalidOption, d)
		}
		c.timeout = d
		return nil
	}
}

// WithMaxConns sets the maximum number of concurrent connections.
// n must be >= 1.
func WithMaxConns(n int) Option {
	return func(c *config) error {
		if n < 1 {
			return fmt.Errorf("%w: maxConns must be >= 1, got %d", ErrInvalidOption, n)
		}
		c.maxConns = n
		return nil
	}
}

// WithRetries sets the number of retry attempts. n must be >= 0.
func WithRetries(n int) Option {
	return func(c *config) error {
		if n < 0 {
			return fmt.Errorf("%w: retries must be >= 0, got %d", ErrInvalidOption, n)
		}
		c.retries = n
		return nil
	}
}

// WithTLS toggles TLS. This option cannot fail, so it always returns nil — the
// canonical shape for a non-validating option.
func WithTLS(enabled bool) Option {
	return func(c *config) error {
		c.tls = enabled
		return nil
	}
}

// WithLoggerName sets a label for the server's logger. An empty name is rejected.
func WithLoggerName(name string) Option {
	return func(c *config) error {
		if name == "" {
			return fmt.Errorf("%w: logger name must not be empty", ErrInvalidOption)
		}
		c.loggerName = name
		return nil
	}
}

// Server is the configured object the constructor produces. Its fields are
// unexported and exposed through accessors so the resolved config is immutable
// after New returns.
type Server struct {
	addr string
	cfg  config
}

// New builds a Server for addr, applying opts in order over the defaults.
// Later options override earlier ones for the same field. If any option fails,
// New returns the wrapped error and a nil Server, leaving no partially built
// object behind. New also rejects an empty addr.
func New(addr string, opts ...Option) (*Server, error) {
	if addr == "" {
		return nil, fmt.Errorf("%w: addr must not be empty", ErrInvalidOption)
	}

	cfg := defaults()
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

// Addr returns the address the server was built for.
func (s *Server) Addr() string { return s.addr }

// Timeout returns the resolved request timeout.
func (s *Server) Timeout() time.Duration { return s.cfg.timeout }

// MaxConns returns the resolved maximum connection count.
func (s *Server) MaxConns() int { return s.cfg.maxConns }

// Retries returns the resolved retry count.
func (s *Server) Retries() int { return s.cfg.retries }

// TLS reports whether TLS is enabled.
func (s *Server) TLS() bool { return s.cfg.tls }

// LoggerName returns the resolved logger label.
func (s *Server) LoggerName() string { return s.cfg.loggerName }
