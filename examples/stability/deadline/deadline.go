// Package deadline demonstrates the Deadline stability pattern.
//
// A deadline bounds the maximum wall-clock time a request may occupy,
// end-to-end. Every layer propagates the same context so that sub-calls
// automatically inherit the remaining budget. This is distinct from a
// per-attempt timeout: the budget shrinks as it passes through layers and
// is never reset on retry.
package deadline

import (
	"context"
	"errors"
	"time"
)

// ErrBudgetExceeded is the sentinel returned when a layer detects that the
// remaining context budget is too small to be useful before even attempting
// work. This enables fail-fast inside the call graph, complementing the
// automatic cancellation that ctx.Done() provides.
var ErrBudgetExceeded = errors.New("deadline: insufficient budget remaining")

// MinBudget is the minimum remaining duration considered worth attempting
// an RPC or sub-call. Callers may override this in their own validation.
const MinBudget = 5 * time.Millisecond

// WithBudget creates a child context whose deadline is set to now+budget.
// If a deadline is already set on the parent and it expires sooner, the
// parent deadline wins — context.WithDeadline guarantees this.
func WithBudget(parent context.Context, budget time.Duration) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, time.Now().Add(budget))
}

// RemainingBudget returns the time left before ctx's deadline. If no deadline
// is set it returns a large sentinel value (1 hour) and ok=false. Negative
// values mean the deadline has already passed.
func RemainingBudget(ctx context.Context) (remaining time.Duration, ok bool) {
	dl, set := ctx.Deadline()
	if !set {
		return time.Hour, false
	}
	return time.Until(dl), true
}

// CheckBudget returns ErrBudgetExceeded if the remaining context budget is
// smaller than min. Embed this at service-layer boundaries before expensive
// work to propagate refusals eagerly.
func CheckBudget(ctx context.Context, min time.Duration) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	rem, set := RemainingBudget(ctx)
	if set && rem < min {
		return ErrBudgetExceeded
	}
	return nil
}

// Call represents a downstream dependency call. Work simulates latency.
type Call struct {
	// Work is called with the propagated context. Implementations must
	// respect ctx.Done() and return promptly when it fires.
	Work func(ctx context.Context) error
}

// Do runs the Call, propagating ctx unchanged. It first checks that the
// remaining budget exceeds MinBudget and returns ErrBudgetExceeded if not.
func (c *Call) Do(ctx context.Context) error {
	if err := CheckBudget(ctx, MinBudget); err != nil {
		return err
	}
	return c.Work(ctx)
}

// Chain runs a sequence of Calls with the same context, short-circuiting on
// the first error. Each call inherits the shrinking remaining budget
// automatically via context propagation.
func Chain(ctx context.Context, calls ...*Call) error {
	for _, call := range calls {
		if err := call.Do(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Attempt runs fn with a per-attempt timeout that is capped to the
// remaining parent budget. This is useful for retries where each attempt
// should not consume more than attemptTimeout, but can never exceed what is
// left on the end-to-end deadline.
//
// The pattern is: budget = min(attemptTimeout, remainingParentBudget).
func Attempt(parent context.Context, attemptTimeout time.Duration, fn func(context.Context) error) error {
	remaining, set := RemainingBudget(parent)
	if set && remaining <= MinBudget {
		return ErrBudgetExceeded
	}

	// Cap the per-attempt timeout to the remaining parent budget.
	budget := attemptTimeout
	if set && remaining < attemptTimeout {
		budget = remaining
	}

	ctx, cancel := context.WithTimeout(parent, budget)
	defer cancel()
	return fn(ctx)
}

// IsDeadlineError reports whether err is a context deadline or budget error.
func IsDeadlineError(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, ErrBudgetExceeded) ||
		errors.Is(err, context.Canceled)
}
