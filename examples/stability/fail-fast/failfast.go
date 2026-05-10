// Package failfast demonstrates the Fail-Fast stability pattern.
//
// Fail-Fast systems surface errors as early as possible — during startup
// validation, at request ingress, or at the first sign a dependency is
// unready — rather than allowing bad state to propagate deep into the call
// graph. This is the opposite default from retry-with-backoff: the system
// refuses work it cannot reasonably complete instead of queuing it.
package failfast

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
)

// --- Configuration validation (startup fail-fast) --------------------------

// Config holds required service configuration.
type Config struct {
	Addr        string
	MaxConns    int
	ServiceName string
}

// ErrInvalidConfig is returned when Validate detects a misconfiguration.
type ErrInvalidConfig struct {
	Field   string
	Message string
}

func (e *ErrInvalidConfig) Error() string {
	return fmt.Sprintf("invalid config: field %q: %s", e.Field, e.Message)
}

// Validate checks all config fields at startup. The caller should call this
// once during program initialisation and exit if it returns a non-nil error.
// Returning all problems at once (not just the first) gives operators a
// complete picture.
func (c Config) Validate() error {
	var errs []string
	if c.Addr == "" {
		errs = append(errs, (&ErrInvalidConfig{"Addr", "must not be empty"}).Error())
	}
	if c.MaxConns <= 0 {
		errs = append(errs, (&ErrInvalidConfig{"MaxConns", "must be > 0"}).Error())
	}
	if c.ServiceName == "" {
		errs = append(errs, (&ErrInvalidConfig{"ServiceName", "must not be empty"}).Error())
	}
	if len(errs) > 0 {
		return errors.New("config validation failed:\n  " + strings.Join(errs, "\n  "))
	}
	return nil
}

// --- Readiness probe (dependency fail-fast) ---------------------------------

// ErrNotReady is returned when the server or a dependency is not yet ready.
var ErrNotReady = errors.New("failfast: service not ready")

// ReadinessProbe is a function that checks whether a dependency is available.
type ReadinessProbe func(ctx context.Context) error

// Handler wraps an HTTP-style handler and rejects requests early when the
// service is not ready. In production this would wrap an http.Handler; here
// it models the pattern with plain functions for testability.
type Handler struct {
	ready  atomic.Bool
	probes []ReadinessProbe
}

// SetReady marks the service as ready or not-ready. Call SetReady(true) after
// all startup checks pass; call SetReady(false) during graceful shutdown.
func (h *Handler) SetReady(ready bool) {
	h.ready.Store(ready)
}

// IsReady returns true when the handler has been marked ready.
func (h *Handler) IsReady() bool {
	return h.ready.Load()
}

// AddProbe registers a readiness probe that will be checked on every Handle
// call while the handler is not in the ready state.
func (h *Handler) AddProbe(p ReadinessProbe) {
	h.probes = append(h.probes, p)
}

// Handle runs fn only if the service is ready and all probes pass. Otherwise
// it returns ErrNotReady immediately — fail-fast at the handler boundary.
func (h *Handler) Handle(ctx context.Context, fn func(context.Context) error) error {
	if !h.ready.Load() {
		return ErrNotReady
	}
	// Even when marked ready, actively verify probes (e.g. DB ping).
	for _, probe := range h.probes {
		if err := probe(ctx); err != nil {
			return fmt.Errorf("%w: probe: %w", ErrNotReady, err)
		}
	}
	return fn(ctx)
}

// --- Request validation (per-request fail-fast) ----------------------------

// Request represents an inbound request with fields to validate.
type Request struct {
	UserID  string
	Payload []byte
	Limit   int
}

// ErrBadRequest signals that a request failed early validation.
type ErrBadRequest struct{ Reason string }

func (e *ErrBadRequest) Error() string { return "bad request: " + e.Reason }

// Validate checks the request at the ingress boundary. Returning early here
// avoids wasted work in downstream layers.
func (r Request) Validate() error {
	if r.UserID == "" {
		return &ErrBadRequest{"UserID must not be empty"}
	}
	if len(r.Payload) == 0 {
		return &ErrBadRequest{"Payload must not be empty"}
	}
	if r.Limit <= 0 || r.Limit > 1000 {
		return &ErrBadRequest{fmt.Sprintf("Limit must be 1-1000, got %d", r.Limit)}
	}
	return nil
}
