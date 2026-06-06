---
name: stability-circuit-breaker
description: Circuit Breaker — track failures and open the circuit after threshold to fail fast and let the dependency recover. Use for unreliable downstreams (external APIs).
category: stability
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/stability/circuit-breaker.md
example: examples/stability/circuit-breaker/
---

# Circuit Breaker Pattern

Similar to electrical fuses that prevent fires when a circuit that is connected
to the electrical grid starts drawing a high amount of power which causes the
wires to heat up and combust, the circuit breaker design pattern is a fail-first
mechanism that shuts down the circuit, request/response relationship or a
service in the case of software development, to prevent bigger failures.

**Note:** The words "circuit" and "service" are used synonymously throught this
document.

## Implementation

Below is the implementation of a very simple circuit breaker to illustrate the purpose
of the circuit breaker design pattern.

### Operation Counter

`circuit.Counter` is a simple counter that records success and failure states of
a circuit along with a timestamp and calculates the consecutive number of
failures.

```go
package circuit

import (
	"time"
)

type State int

const (
	UnknownState State = iota
	FailureState
	SuccessState
)

type Counter interface {
	Count(State)
	ConsecutiveFailures() uint32
	LastActivity() time.Time
	Reset()
}
```

### Circuit Breaker

Circuit is wrapped using the `circuit.Breaker` closure that keeps an internal operation counter.
It returns a fast error if the circuit has failed consecutively more than the specified threshold.
After a while it retries the request and records it.

**Note:** Context type is used here to carry deadlines, cancelation signals, and
other request-scoped values across API boundaries and between processes.

```go
package circuit

import (
	"context"
	"errors"
	"time"
)

// ErrServiceUnavailable is returned when the breaker is open and fails fast.
var ErrServiceUnavailable = errors.New("circuit: service unavailable")

type Circuit func(context.Context) error

// Breaker wraps c and trips after failureThreshold consecutive failures.
// While tripped it fails fast with ErrServiceUnavailable; after an
// exponentially growing backoff it lets a single request through to probe the
// dependency. NewCounter returns a thread-safe Counter implementation.
func Breaker(c Circuit, failureThreshold uint32) Circuit {
	cnt := NewCounter()

	return func(ctx context.Context) error {
		if cnt.ConsecutiveFailures() >= failureThreshold {
			// Bug-prone spots to watch: the parameter is context.Context (not
			// the package "context"), the closure must declare its bool return,
			// the constant is time.Second (singular), and don't shadow cnt.
			canRetry := func() bool {
				backoffLevel := cnt.ConsecutiveFailures() - failureThreshold

				// When should the breaker resume propagating requests.
				shouldRetryAt := cnt.LastActivity().Add(time.Second * 2 << backoffLevel)

				return time.Now().After(shouldRetryAt)
			}

			if !canRetry() {
				// Fail fast: not enough time has passed since the last failure.
				return ErrServiceUnavailable
			}
		}

		// Unless the failure threshold is exceeded the wrapped service mimics
		// the old behavior; the difference appears only after consecutive
		// failures.
		if err := c(ctx); err != nil {
			cnt.Count(FailureState)
			return err
		}

		cnt.Count(SuccessState)
		return nil
	}
}
```

A complete, generic, deterministically tested implementation (with an injectable
clock and explicit open/half-open states) lives in the runnable example at
`examples/stability/circuit-breaker/`.

## Related Works

- [sony/gobreaker](https://github.com/sony/gobreaker) is a well-tested and intuitive circuit breaker implementation for real-world use cases.
