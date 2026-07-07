package imports

import (
	"sync"

	"github.com/lsdch/biome/lib/progress"
	"github.com/lsdch/biome/models"
)

type ImportEvent struct {
	Workflow models.ImportWorkflow
	Status   RunnerStatus
	GBIF     progress.ProgressSnapshot
	Error    error
}

type ImportEventSink[T any] interface {
	Publish(T)
}

type EventBroker[T any] struct {
	mu          sync.RWMutex
	subscribers map[chan T]struct{}
}

func NewEventBroker[T any]() *EventBroker[T] {
	return &EventBroker[T]{
		subscribers: make(map[chan T]struct{}),
	}
}

func (b *EventBroker[T]) Subscribe() (<-chan T, func()) {
	ch := make(chan T, 16)

	b.mu.Lock()
	b.subscribers[ch] = struct{}{}
	b.mu.Unlock()

	unsubscribe := func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		if _, ok := b.subscribers[ch]; ok {
			delete(b.subscribers, ch)
			close(ch)
		}
	}

	return ch, unsubscribe
}

func (b *EventBroker[T]) Publish(event T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.subscribers {
		select {
		case ch <- event:
		default:
			// If the channel is full, we can choose to drop the event or block.
			// Here, we choose to drop the event to avoid blocking.
		}
	}
}
