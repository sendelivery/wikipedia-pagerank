package queue

import (
	"sync"
)

// SizedQueue is a queue with a fixed size, it can only contain distinct elements.
type SizedQueue struct {
	ch     chan string
	set    sync.Map
	closed bool
}

// Enqueue adds an element to the queue.
// Returns true if the element was successfully enqueued, false otherwise.
// Only distinct elements are enqueued.
func (q *SizedQueue) Enqueue(s string) bool {
	// If the element has previously been enqueued, don't enqueue it again.
	if q.Has(s) {
		return false
	}

	// If there's space in the queue, enqueue the element.
	select {
	case q.ch <- s:
		q.set.Store(s, struct{}{})
		return true
	default:
		// If the queue is full, close the channel and return false.
		if !q.closed {
			close(q.ch)
			q.closed = true
		}
		return false
	}
}

// Dequeue removes and returns an element from the queue.
func (q *SizedQueue) Dequeue() (string, bool) {
	select {
	case s := <-q.ch:
		return s, true
	default:
		return "", false
	}
}

// Has returns true if the queue contains the given element.
func (q *SizedQueue) Has(s string) bool {
	_, ok := q.set.Load(s)
	return ok
}

// Empty returns true if the queue is empty.
func (q *SizedQueue) Empty() bool {
	return len(q.ch) == 0
}

// Full returns true if the queue is full.
func (q *SizedQueue) Full() bool {
	return q.closed
}

// Length returns the number of elements in the queue.
func (q *SizedQueue) Length() int {
	return len(q.ch)
}

// Close closes the queue.
func (q *SizedQueue) Close() {
	if !q.closed {
		close(q.ch)
		q.closed = true
	}
}

// NewSizedQueue returns a new SizedQueue with the given size.
func NewSizedQueue(size int) *SizedQueue {
	return &SizedQueue{
		ch:     make(chan string, size),
		set:    sync.Map{},
		closed: false,
	}
}
