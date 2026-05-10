// Package reactor implements a single-goroutine event loop that demultiplexes
// events from multiple sources and dispatches them to registered handlers.
//
// All handlers execute sequentially inside the Run goroutine, so they may
// read/write shared state without additional synchronisation — but must not
// block. Long-running work should be dispatched to a separate goroutine.
package reactor

import "context"

// Event carries a kind identifier and an arbitrary payload.
type Event struct {
	Kind    string
	Payload any
}

// Handler is called synchronously by the reactor loop.
type Handler func(Event)

// Reactor is a single-goroutine event loop.
// Zero value is not usable; use New.
type Reactor struct {
	events   chan Event
	handlers map[string]Handler
}

// New returns a Reactor with a buffered event queue of the given size.
// Panics if buf < 0.
func New(buf int) *Reactor {
	return &Reactor{
		events:   make(chan Event, buf),
		handlers: make(map[string]Handler),
	}
}

// Register associates kind with h. Must be called before Run.
// Registering the same kind twice overwrites the previous handler.
func (r *Reactor) Register(kind string, h Handler) {
	r.handlers[kind] = h
}

// Send enqueues an event for dispatch. Blocks if the internal buffer is full.
// Safe to call from any goroutine.
func (r *Reactor) Send(e Event) {
	r.events <- e
}

// TrySend enqueues an event without blocking. Returns false if the buffer is full.
func (r *Reactor) TrySend(e Event) bool {
	select {
	case r.events <- e:
		return true
	default:
		return false
	}
}

// Run starts the event loop in the calling goroutine and blocks until ctx is
// cancelled. Events in the buffer when ctx cancels are drained and dispatched
// before Run returns, ensuring at-most-once delivery for queued events.
func (r *Reactor) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Drain remaining events before exit.
			r.drain()
			return
		case e := <-r.events:
			r.dispatch(e)
		}
	}
}

func (r *Reactor) drain() {
	for {
		select {
		case e := <-r.events:
			r.dispatch(e)
		default:
			return
		}
	}
}

func (r *Reactor) dispatch(e Event) {
	if h, ok := r.handlers[e.Kind]; ok {
		h(e)
	}
}
