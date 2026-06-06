// Package observer is an example of the Observer Pattern.
// Push model.
package observer

import "sync"

// Publisher interface.
type Publisher interface {
	Attach(observer Observer)
	Detach(observer Observer)
	SetState(state string)
	Notify()
}

// Observer provides a subscriber interface.
type Observer interface {
	Update(state string)
}

// ConcretePublisher implements the Publisher interface.
//
// The observer list is mutated by Attach/Detach and iterated by Notify, all of
// which may run concurrently from different goroutines. A sync.Mutex guards the
// shared state so those operations are safe under the race detector.
type ConcretePublisher struct {
	mu        sync.Mutex
	observers []Observer
	state     string
}

// NewPublisher is the Publisher constructor.
func NewPublisher() Publisher {
	return &ConcretePublisher{}
}

// Attach a Observer
func (s *ConcretePublisher) Attach(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

// Detach removes a previously attached Observer. Detaching an Observer that was
// never attached is a no-op.
func (s *ConcretePublisher) Detach(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, o := range s.observers {
		if o == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			return
		}
	}
}

// SetState sets new state
func (s *ConcretePublisher) SetState(state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = state
}

// Notify sends notifications to subscribers.
// Push model.
//
// The observer list and state are snapshotted under the lock, then the
// observers are updated after releasing it. This avoids holding the mutex while
// running arbitrary observer code (which could call back into the publisher and
// deadlock).
func (s *ConcretePublisher) Notify() {
	s.mu.Lock()
	state := s.state
	observers := make([]Observer, len(s.observers))
	copy(observers, s.observers)
	s.mu.Unlock()

	for _, observer := range observers {
		observer.Update(state)
	}
}

// ConcreteObserver implements the Observer interface.
type ConcreteObserver struct {
	state string
}

// Update set new state
func (s *ConcreteObserver) Update(state string) {
	s.state = state
}
