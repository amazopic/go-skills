package deadline

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithBudget_UsesShortestDeadline(t *testing.T) {
	// Parent already has a tight deadline.
	parent, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Request a much longer budget — parent deadline should win.
	ctx, done := WithBudget(parent, time.Hour)
	defer done()

	dl, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	// Should be within ~10ms, not 1h.
	if time.Until(dl) > 20*time.Millisecond {
		t.Errorf("expected parent deadline to win, remaining=%v", time.Until(dl))
	}
}

func TestRemainingBudget_NoDeadline(t *testing.T) {
	_, ok := context.Background().Deadline()
	if ok {
		t.Skip("background ctx has a deadline in this env")
	}
	rem, set := RemainingBudget(context.Background())
	if set {
		t.Error("expected set=false for background context")
	}
	if rem < time.Minute {
		t.Errorf("sentinel should be large, got %v", rem)
	}
}

func TestCheckBudget_PassesWithSufficientBudget(t *testing.T) {
	ctx, cancel := WithBudget(context.Background(), time.Second)
	defer cancel()
	if err := CheckBudget(ctx, MinBudget); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCheckBudget_RejectsExpiredBudget(t *testing.T) {
	ctx, cancel := WithBudget(context.Background(), time.Millisecond)
	defer cancel()
	time.Sleep(5 * time.Millisecond) // let it expire
	err := CheckBudget(ctx, MinBudget)
	if err == nil {
		t.Fatal("expected error for expired budget")
	}
}

func TestCheckBudget_RejectsTooSmallBudget(t *testing.T) {
	ctx, cancel := WithBudget(context.Background(), MinBudget/2)
	defer cancel()
	err := CheckBudget(ctx, MinBudget)
	if !errors.Is(err, ErrBudgetExceeded) {
		t.Errorf("want ErrBudgetExceeded, got %v", err)
	}
}

func TestCall_Do_SucceedsWithinBudget(t *testing.T) {
	ctx, cancel := WithBudget(context.Background(), 200*time.Millisecond)
	defer cancel()

	call := &Call{Work: func(ctx context.Context) error {
		return nil
	}}
	if err := call.Do(ctx); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCall_Do_RejectsExhaustedBudget(t *testing.T) {
	ctx, cancel := WithBudget(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(5 * time.Millisecond)

	call := &Call{Work: func(ctx context.Context) error { return nil }}
	err := call.Do(ctx)
	if err == nil {
		t.Fatal("expected deadline error")
	}
}

func TestChain_ShortCircuitsOnError(t *testing.T) {
	ctx, cancel := WithBudget(context.Background(), time.Second)
	defer cancel()

	sentinel := errors.New("step2 failed")
	called := 0

	calls := []*Call{
		{Work: func(ctx context.Context) error { called++; return nil }},
		{Work: func(ctx context.Context) error { called++; return sentinel }},
		{Work: func(ctx context.Context) error { called++; return nil }},
	}

	err := Chain(ctx, calls...)
	if !errors.Is(err, sentinel) {
		t.Errorf("want sentinel, got %v", err)
	}
	if called != 2 {
		t.Errorf("want 2 calls before short-circuit, got %d", called)
	}
}

func TestAttempt_CapsToRemainingBudget(t *testing.T) {
	// Parent has 50ms left; attemptTimeout=1s → should be capped to ~50ms.
	parent, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var gotDeadline time.Time
	err := Attempt(parent, time.Second, func(ctx context.Context) error {
		dl, _ := ctx.Deadline()
		gotDeadline = dl
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The deadline set inside Attempt must be no more than ~50ms from now.
	if time.Until(gotDeadline) > 60*time.Millisecond {
		t.Errorf("attempt deadline not capped: %v remaining", time.Until(gotDeadline))
	}
}

func TestAttempt_RejectsExhaustedParent(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(5 * time.Millisecond)

	err := Attempt(parent, time.Second, func(ctx context.Context) error { return nil })
	if !IsDeadlineError(err) {
		t.Errorf("want deadline error, got %v", err)
	}
}

func TestIsDeadlineError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"DeadlineExceeded", context.DeadlineExceeded, true},
		{"Canceled", context.Canceled, true},
		{"ErrBudgetExceeded", ErrBudgetExceeded, true},
		{"nil", nil, false},
		{"other", errors.New("random"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDeadlineError(tt.err); got != tt.want {
				t.Errorf("IsDeadlineError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
