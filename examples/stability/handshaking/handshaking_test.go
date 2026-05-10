package handshaking

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestServer_Capacity_FullAdvertised(t *testing.T) {
	s := NewServer(3)
	resp := s.Capacity()
	if resp.Slots != 3 {
		t.Errorf("want Slots=3, got %d", resp.Slots)
	}
}

func TestServer_Capacity_DecreasesUnderLoad(t *testing.T) {
	s := NewServer(3)
	ctx := context.Background()

	started := make(chan struct{})
	release := make(chan struct{})

	go func() {
		_ = s.Handle(ctx, func(ctx context.Context) error {
			close(started)
			<-release
			return nil
		})
	}()
	<-started

	resp := s.Capacity()
	if resp.Slots != 2 {
		t.Errorf("want Slots=2 with 1 in-flight, got %d", resp.Slots)
	}
	close(release)
}

func TestServer_Handle_RejectsOverCapacity(t *testing.T) {
	s := NewServer(1)
	ctx := context.Background()

	started := make(chan struct{})
	release := make(chan struct{})

	go func() {
		_ = s.Handle(ctx, func(ctx context.Context) error {
			close(started)
			<-release
			return nil
		})
	}()
	<-started

	err := s.Handle(ctx, func(ctx context.Context) error { return nil })
	if !errors.Is(err, ErrCapacityExceeded) {
		t.Errorf("want ErrCapacityExceeded, got %v", err)
	}
	close(release)
}

func TestServer_InFlight_SettlesAfterDone(t *testing.T) {
	s := NewServer(5)
	ctx := context.Background()

	_ = s.Handle(ctx, func(ctx context.Context) error { return nil })
	if s.InFlight() != 0 {
		t.Errorf("want InFlight=0 after handle, got %d", s.InFlight())
	}
}

func TestClient_Send_RejectsWhenNoSlots(t *testing.T) {
	fetcher := func(ctx context.Context) (CapacityResponse, error) {
		return CapacityResponse{Slots: 0, LeaseTTL: 0}, nil
	}
	c := NewClient(fetcher)

	err := c.Send(context.Background(), func(ctx context.Context) error { return nil })
	if !errors.Is(err, ErrCapacityExceeded) {
		t.Errorf("want ErrCapacityExceeded, got %v", err)
	}

	_, rej := c.Stats()
	if rej != 1 {
		t.Errorf("want rejected=1, got %d", rej)
	}
}

func TestClient_Send_SucceedsWithSlots(t *testing.T) {
	fetcher := func(ctx context.Context) (CapacityResponse, error) {
		return CapacityResponse{Slots: 5, LeaseTTL: 0}, nil
	}
	c := NewClient(fetcher)

	err := c.Send(context.Background(), func(ctx context.Context) error { return nil })
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	sent, _ := c.Stats()
	if sent != 1 {
		t.Errorf("want sent=1, got %d", sent)
	}
}

func TestClient_CachesLease(t *testing.T) {
	fetchCount := 0
	fetcher := func(ctx context.Context) (CapacityResponse, error) {
		fetchCount++
		return CapacityResponse{Slots: 10, LeaseTTL: 5 * time.Second}, nil
	}
	c := NewClient(fetcher)

	for i := 0; i < 5; i++ {
		_ = c.Send(context.Background(), func(ctx context.Context) error { return nil })
	}

	// Should have fetched only once because LeaseTTL hasn't expired.
	if fetchCount != 1 {
		t.Errorf("want 1 fetch within lease, got %d", fetchCount)
	}
}

func TestClient_RefetchesAfterLeaseExpiry(t *testing.T) {
	fetchCount := 0
	fetcher := func(ctx context.Context) (CapacityResponse, error) {
		fetchCount++
		return CapacityResponse{Slots: 10, LeaseTTL: 10 * time.Millisecond}, nil
	}
	c := NewClient(fetcher)

	_ = c.Send(context.Background(), func(ctx context.Context) error { return nil })
	time.Sleep(20 * time.Millisecond) // let lease expire
	_ = c.Send(context.Background(), func(ctx context.Context) error { return nil })

	if fetchCount < 2 {
		t.Errorf("want >= 2 fetches after lease expiry, got %d", fetchCount)
	}
}

func TestClient_FetcherError_Propagated(t *testing.T) {
	sentinel := errors.New("network error")
	fetcher := func(ctx context.Context) (CapacityResponse, error) {
		return CapacityResponse{}, sentinel
	}
	c := NewClient(fetcher)

	err := c.Send(context.Background(), func(ctx context.Context) error { return nil })
	if !errors.Is(err, sentinel) {
		t.Errorf("want sentinel fetcher error, got %v", err)
	}
}

func TestServerClient_Integration(t *testing.T) {
	s := NewServer(2)
	fetcher := func(ctx context.Context) (CapacityResponse, error) {
		return s.Capacity(), nil
	}
	c := NewClient(fetcher)
	ctx := context.Background()

	var wg sync.WaitGroup
	errs := make([]error, 4)
	started := make(chan struct{}, 2)
	release := make(chan struct{})

	// Hold 2 slots on the server.
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Handle(ctx, func(ctx context.Context) error {
				started <- struct{}{}
				<-release
				return nil
			})
		}()
	}
	// Wait until both are in-flight.
	<-started
	<-started

	// Now send 4 requests via client — all should be rejected (Slots=0).
	for i := 0; i < 4; i++ {
		errs[i] = c.Send(ctx, func(ctx context.Context) error { return nil })
	}

	close(release)
	wg.Wait()

	for i, err := range errs {
		if !errors.Is(err, ErrCapacityExceeded) {
			t.Errorf("request %d: want ErrCapacityExceeded, got %v", i, err)
		}
	}
}
