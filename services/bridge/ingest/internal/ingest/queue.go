package ingest

import (
	"context"
	"sync"
)

// DropQueue is a bounded FIFO that drops the oldest item when full. Push never blocks.
type DropQueue[T any] struct {
	mu     sync.Mutex
	items  []T
	head   int
	n      int
	notify chan struct{}
}

func NewDropQueue[T any](capacity int) *DropQueue[T] {
	return &DropQueue[T]{items: make([]T, max(capacity, 1)), notify: make(chan struct{}, 1)}
}

// Push adds v and reports whether the oldest item was dropped to make room.
func (q *DropQueue[T]) Push(v T) bool {
	q.mu.Lock()
	dropped := false
	if q.n == len(q.items) {
		var zero T
		q.items[q.head] = zero
		q.head = (q.head + 1) % len(q.items)
		q.n--
		dropped = true
	}
	q.items[(q.head+q.n)%len(q.items)] = v
	q.n++
	q.mu.Unlock()
	select {
	case q.notify <- struct{}{}:
	default:
	}
	return dropped
}

// Pop waits for an item or ctx.
func (q *DropQueue[T]) Pop(ctx context.Context) (T, bool) {
	for {
		q.mu.Lock()
		if q.n > 0 {
			v := q.items[q.head]
			var zero T
			q.items[q.head] = zero
			q.head = (q.head + 1) % len(q.items)
			q.n--
			q.mu.Unlock()
			return v, true
		}
		q.mu.Unlock()
		select {
		case <-q.notify:
		case <-ctx.Done():
			var zero T
			return zero, false
		}
	}
}

func (q *DropQueue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.n
}

func (q *DropQueue[T]) Cap() int { return len(q.items) }
