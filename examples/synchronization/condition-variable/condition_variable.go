// Package condvar demonstrates sync.Cond with a classic bounded queue:
// producers block when the queue is full; consumers block when it is empty.
// Signal wakes exactly one waiter; Broadcast is shown via Drain.
package condvar

import "sync"

// BoundedQueue is a FIFO queue with a fixed capacity.
// Producers block when full; consumers block when empty.
type BoundedQueue struct {
	mu   sync.Mutex
	cond *sync.Cond
	buf  []int
	cap  int
}

// NewBoundedQueue creates a BoundedQueue with the given capacity.
// Panics if cap < 1.
func NewBoundedQueue(cap int) *BoundedQueue {
	if cap < 1 {
		panic("condvar: cap must be >= 1")
	}
	q := &BoundedQueue{cap: cap}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Put adds v to the queue, blocking until space is available.
func (q *BoundedQueue) Put(v int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.buf) == q.cap { // Mesa semantics: always loop
		q.cond.Wait()
	}
	q.buf = append(q.buf, v)
	q.cond.Signal() // wake exactly one consumer
}

// Get removes and returns the front item, blocking until one is available.
func (q *BoundedQueue) Get() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.buf) == 0 { // always loop
		q.cond.Wait()
	}
	v := q.buf[0]
	q.buf = q.buf[1:]
	q.cond.Signal() // wake exactly one producer
	return v
}

// Len returns the current number of items in the queue.
func (q *BoundedQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.buf)
}

// Drain removes and returns all current items, waking all waiters.
// Useful for shutdown: producers and consumers are woken so they can exit.
func (q *BoundedQueue) Drain() []int {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]int, len(q.buf))
	copy(out, q.buf)
	q.buf = q.buf[:0]
	q.cond.Broadcast() // wake ALL waiters (e.g., producers that were blocked)
	return out
}
