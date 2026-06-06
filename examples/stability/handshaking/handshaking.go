// Package handshaking demonstrates the Handshaking stability pattern.
//
// Handshaking is a protocol-level exchange in which a client asks a server
// "how much work can you accept right now?" before submitting that work.
// The server advertises its current capacity (concurrency budget) and an
// optional lease TTL; the client adapts its send rate accordingly.
//
// Canonical real-world example: HTTP/2 SETTINGS frames, where the server
// advertises SETTINGS_MAX_CONCURRENT_STREAMS to rate-limit clients.
//
// This package models a lightweight /capacity endpoint + adaptive client.
package handshaking

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// CapacityResponse is the payload returned by a server's capacity endpoint.
type CapacityResponse struct {
	// Slots is the number of concurrent requests the server will accept.
	Slots int
	// LeaseTTL is how long the client should trust this advertisement before
	// re-querying. Zero means "no caching; query every time".
	LeaseTTL time.Duration
}

// ErrCapacityExceeded is returned when the server's advertised capacity is 0.
var ErrCapacityExceeded = errors.New("handshaking: server has no available capacity")

// Server tracks in-flight requests and advertises its remaining capacity.
// It plays the role of the service behind the /capacity endpoint.
type Server struct {
	mu       sync.Mutex
	limit    int
	inFlight int
}

// NewServer creates a Server with the given concurrency limit.
func NewServer(limit int) *Server {
	if limit <= 0 {
		panic("handshaking: limit must be > 0")
	}
	return &Server{limit: limit}
}

// Capacity returns the server's current concurrency advertisement.
// In a real service this would be the HTTP handler for GET /capacity.
func (s *Server) Capacity() CapacityResponse {
	s.mu.Lock()
	avail := s.limit - s.inFlight
	s.mu.Unlock()
	if avail < 0 {
		avail = 0
	}
	return CapacityResponse{
		Slots:    avail,
		LeaseTTL: 500 * time.Millisecond,
	}
}

// Handle executes fn if the server has capacity. It enforces the limit
// server-side as well; the handshake reduces contention but is not the only
// guard.
func (s *Server) Handle(ctx context.Context, fn func(context.Context) error) error {
	s.mu.Lock()
	if s.inFlight >= s.limit {
		s.mu.Unlock()
		return ErrCapacityExceeded
	}
	s.inFlight++
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.inFlight--
		s.mu.Unlock()
	}()

	return fn(ctx)
}

// InFlight returns the current number of in-flight requests.
func (s *Server) InFlight() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.inFlight
}

// CapacityFetcher is a function that queries the server's capacity endpoint.
// In production this would be an HTTP GET; here it is a plain function for
// testability without networking.
type CapacityFetcher func(ctx context.Context) (CapacityResponse, error)

// Client performs the client-side handshake: it fetches the server's capacity
// advertisement, caches it for LeaseTTL, and refuses to send if Slots == 0.
type Client struct {
	fetch CapacityFetcher

	mu          sync.Mutex
	cached      CapacityResponse
	cacheExpiry time.Time

	// Metrics
	sent     atomic.Int64
	rejected atomic.Int64
}

// NewClient creates a Client backed by the given capacity fetcher.
func NewClient(fetch CapacityFetcher) *Client {
	return &Client{fetch: fetch}
}

// capacity returns the cached advertisement or re-fetches if the lease expired.
func (c *Client) capacity(ctx context.Context) (CapacityResponse, error) {
	c.mu.Lock()
	if time.Now().Before(c.cacheExpiry) && c.cacheExpiry != (time.Time{}) {
		resp := c.cached
		c.mu.Unlock()
		return resp, nil
	}
	c.mu.Unlock()

	resp, err := c.fetch(ctx)
	if err != nil {
		return CapacityResponse{}, err
	}

	c.mu.Lock()
	c.cached = resp
	if resp.LeaseTTL > 0 {
		c.cacheExpiry = time.Now().Add(resp.LeaseTTL)
	}
	c.mu.Unlock()
	return resp, nil
}

// Send performs the client handshake and, if the server has capacity, calls
// fn. It returns ErrCapacityExceeded without calling fn when Slots == 0.
func (c *Client) Send(ctx context.Context, fn func(context.Context) error) error {
	resp, err := c.capacity(ctx)
	if err != nil {
		return fmt.Errorf("handshaking: capacity check: %w", err)
	}
	if resp.Slots <= 0 {
		c.rejected.Add(1)
		return ErrCapacityExceeded
	}

	c.sent.Add(1)
	return fn(ctx)
}

// Stats returns (sent, rejected) counters.
func (c *Client) Stats() (sent, rejected int64) {
	return c.sent.Load(), c.rejected.Load()
}
