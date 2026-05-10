package failfast

import (
	"context"
	"errors"
	"testing"
)

// --- Config.Validate tests --------------------------------------------------

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name:    "valid config",
			cfg:     Config{Addr: ":8080", MaxConns: 100, ServiceName: "svc"},
			wantErr: false,
		},
		{
			name:    "missing Addr",
			cfg:     Config{MaxConns: 10, ServiceName: "svc"},
			wantErr: true,
		},
		{
			name:    "zero MaxConns",
			cfg:     Config{Addr: ":8080", MaxConns: 0, ServiceName: "svc"},
			wantErr: true,
		},
		{
			name:    "negative MaxConns",
			cfg:     Config{Addr: ":8080", MaxConns: -1, ServiceName: "svc"},
			wantErr: true,
		},
		{
			name:    "missing ServiceName",
			cfg:     Config{Addr: ":8080", MaxConns: 10},
			wantErr: true,
		},
		{
			name:    "all fields missing",
			cfg:     Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// --- Handler tests ----------------------------------------------------------

func TestHandler_RejectsWhenNotReady(t *testing.T) {
	h := &Handler{}
	err := h.Handle(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if !errors.Is(err, ErrNotReady) {
		t.Errorf("want ErrNotReady, got %v", err)
	}
}

func TestHandler_AllowsWhenReady(t *testing.T) {
	h := &Handler{}
	h.SetReady(true)
	err := h.Handle(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandler_ProbeFailureMakesNotReady(t *testing.T) {
	h := &Handler{}
	h.SetReady(true)
	failingProbe := ReadinessProbe(func(ctx context.Context) error {
		return errors.New("db down")
	})
	h.AddProbe(failingProbe)

	err := h.Handle(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if !errors.Is(err, ErrNotReady) {
		t.Errorf("want ErrNotReady wrapping probe error, got %v", err)
	}
}

func TestHandler_SetReady_Transitions(t *testing.T) {
	h := &Handler{}
	if h.IsReady() {
		t.Error("should not be ready initially")
	}
	h.SetReady(true)
	if !h.IsReady() {
		t.Error("should be ready after SetReady(true)")
	}
	h.SetReady(false)
	if h.IsReady() {
		t.Error("should not be ready after SetReady(false)")
	}
}

// --- Request.Validate tests -------------------------------------------------

func TestRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     Request
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     Request{UserID: "u1", Payload: []byte("hello"), Limit: 10},
			wantErr: false,
		},
		{
			name:    "empty UserID",
			req:     Request{Payload: []byte("x"), Limit: 1},
			wantErr: true,
		},
		{
			name:    "empty Payload",
			req:     Request{UserID: "u1", Limit: 1},
			wantErr: true,
		},
		{
			name:    "zero Limit",
			req:     Request{UserID: "u1", Payload: []byte("x"), Limit: 0},
			wantErr: true,
		},
		{
			name:    "Limit too high",
			req:     Request{UserID: "u1", Payload: []byte("x"), Limit: 1001},
			wantErr: true,
		},
		{
			name:    "Limit at max boundary",
			req:     Request{UserID: "u1", Payload: []byte("x"), Limit: 1000},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			var badReq *ErrBadRequest
			if tt.wantErr && !errors.As(err, &badReq) {
				t.Errorf("want ErrBadRequest, got %T: %v", err, err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestHandler_ConcurrentReadiness verifies SetReady / IsReady are race-safe.
func TestHandler_ConcurrentReadiness(t *testing.T) {
	h := &Handler{}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 1000; i++ {
			h.SetReady(i%2 == 0)
		}
	}()
	for {
		select {
		case <-done:
			return
		default:
			_ = h.IsReady()
		}
	}
}
