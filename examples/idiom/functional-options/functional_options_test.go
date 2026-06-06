package options

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestNew_Defaults(t *testing.T) {
	s, err := New("localhost:8080")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if got, want := s.Addr(), "localhost:8080"; got != want {
		t.Errorf("Addr() = %q, want %q", got, want)
	}
	if got, want := s.Timeout(), 30*time.Second; got != want {
		t.Errorf("Timeout() = %v, want %v", got, want)
	}
	if got, want := s.MaxConns(), 100; got != want {
		t.Errorf("MaxConns() = %d, want %d", got, want)
	}
	if got, want := s.Retries(), 3; got != want {
		t.Errorf("Retries() = %d, want %d", got, want)
	}
	if !s.TLS() {
		t.Error("TLS() = false, want true (default)")
	}
	if got, want := s.LoggerName(), "default"; got != want {
		t.Errorf("LoggerName() = %q, want %q", got, want)
	}
}

func TestNew_AppliesAndOverrides(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
		check   func(t *testing.T, s *Server)
	}{
		{
			name: "single option",
			opts: []Option{WithTimeout(5 * time.Second)},
			check: func(t *testing.T, s *Server) {
				if got := s.Timeout(); got != 5*time.Second {
					t.Errorf("Timeout() = %v, want 5s", got)
				}
			},
		},
		{
			name: "multiple options",
			opts: []Option{
				WithMaxConns(10),
				WithRetries(0),
				WithTLS(false),
				WithLoggerName("api"),
			},
			check: func(t *testing.T, s *Server) {
				if s.MaxConns() != 10 {
					t.Errorf("MaxConns() = %d, want 10", s.MaxConns())
				}
				if s.Retries() != 0 {
					t.Errorf("Retries() = %d, want 0", s.Retries())
				}
				if s.TLS() {
					t.Error("TLS() = true, want false")
				}
				if s.LoggerName() != "api" {
					t.Errorf("LoggerName() = %q, want %q", s.LoggerName(), "api")
				}
			},
		},
		{
			name: "later option wins (last-write semantics)",
			opts: []Option{WithTimeout(1 * time.Second), WithTimeout(9 * time.Second)},
			check: func(t *testing.T, s *Server) {
				if got := s.Timeout(); got != 9*time.Second {
					t.Errorf("Timeout() = %v, want 9s (last option should win)", got)
				}
			},
		},
		{
			name: "nil option in slice is skipped",
			opts: []Option{nil, WithRetries(7), nil},
			check: func(t *testing.T, s *Server) {
				if got := s.Retries(); got != 7 {
					t.Errorf("Retries() = %d, want 7", got)
				}
			},
		},
		{
			name: "zero retries allowed",
			opts: []Option{WithRetries(0)},
			check: func(t *testing.T, s *Server) {
				if got := s.Retries(); got != 0 {
					t.Errorf("Retries() = %d, want 0", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New("addr", tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			tt.check(t, s)
		})
	}
}

func TestNew_ValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		addr string
		opts []Option
	}{
		{"empty addr", "", nil},
		{"zero timeout", "addr", []Option{WithTimeout(0)}},
		{"negative timeout", "addr", []Option{WithTimeout(-time.Second)}},
		{"maxConns below 1", "addr", []Option{WithMaxConns(0)}},
		{"negative retries", "addr", []Option{WithRetries(-1)}},
		{"empty logger name", "addr", []Option{WithLoggerName("")}},
		{"valid then invalid", "addr", []Option{WithMaxConns(5), WithTimeout(0)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(tt.addr, tt.opts...)
			if err == nil {
				t.Fatalf("New(%q) = %+v, nil; want error", tt.addr, s)
			}
			if s != nil {
				t.Errorf("New returned non-nil Server alongside error: %+v", s)
			}
			if !errors.Is(err, ErrInvalidOption) {
				t.Errorf("errors.Is(err, ErrInvalidOption) = false; err = %v", err)
			}
		})
	}
}

// TestNew_OptionsAreReusableAndConcurrencySafe verifies that an Option value
// is just a pure function of a config pointer: the same option can be reused
// across many concurrent New calls without sharing state or racing. Each New
// owns its own config, so there is nothing to synchronise.
func TestNew_OptionsAreReusableAndConcurrencySafe(t *testing.T) {
	// Build a shared option set once; reuse it from every goroutine.
	shared := []Option{
		WithTimeout(2 * time.Second),
		WithMaxConns(42),
		WithTLS(false),
		WithLoggerName("worker"),
	}

	const goroutines = 64
	var wg sync.WaitGroup
	wg.Add(goroutines)
	errs := make(chan error, goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			s, err := New("addr", shared...)
			if err != nil {
				errs <- err
				return
			}
			if s.Timeout() != 2*time.Second || s.MaxConns() != 42 ||
				s.TLS() || s.LoggerName() != "worker" {
				errs <- errors.New("unexpected resolved config in concurrent build")
			}
		}()
	}

	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}

// TestNew_DefaultsAreIsolatedPerCall guards against the classic bug of sharing a
// single mutable config across constructions: building with options must not
// leak into a later default-only construction.
func TestNew_DefaultsAreIsolatedPerCall(t *testing.T) {
	if _, err := New("a", WithMaxConns(1), WithTimeout(time.Millisecond)); err != nil {
		t.Fatalf("first New failed: %v", err)
	}
	s, err := New("b")
	if err != nil {
		t.Fatalf("second New failed: %v", err)
	}
	if s.MaxConns() != 100 || s.Timeout() != 30*time.Second {
		t.Errorf("defaults leaked: MaxConns=%d Timeout=%v, want 100 and 30s",
			s.MaxConns(), s.Timeout())
	}
}
