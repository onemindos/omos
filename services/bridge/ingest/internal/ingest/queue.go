package ingest

import (
	"context"
	"sync"
)

// DropQueue is a bounded FIFO. When full, Push drops the oldest item; with Block set,
// PushCtx instead waits for room.
type DropQueue[T any] struct {
	// Block makes PushCtx wait for room instead of dropping the oldest item.
	Block  bool
	mu     sync.Mutex
	items  []T
	head   int
	n      int
	notify chan struct{}
	space  chan struct{}
}

func NewDropQueue[T any](capacity int) *DropQueue[T] {
	return &DropQueue[T]{items: make([]T, max(capacity, 1)), notify: make(chan struct{}, 1), space: make(chan struct{}, 1)}
}

// PushCtx adds v. With Block it waits for room (or ctx) and never drops; otherwise it
// behaves like Push. It reports whether the oldest item was dropped.
func (q *DropQueue[T]) PushCtx(ctx context.Context, v T) (bool, error) {
	if !q.Block {
		return q.Push(v), nil
	}
	for {
		q.mu.Lock()
		if q.n < len(q.items) {
			q.mu.Unlock()
			return q.Push(v), nil // single producer: room cannot vanish in between
		}
		q.mu.Unlock()
		select {
		case <-q.space:
		case <-ctx.Done():
			return false, ctx.Err()
		}
	}
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
			select {
			case q.space <- struct{}{}:
			default:
			}
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
