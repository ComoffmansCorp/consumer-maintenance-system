// Package broker is a minimal in-process publish/subscribe bus used for the
// cross-domain effects described in the architecture doc (e.g. "task
// completed" -> "notification created"). It stands in for a real message
// broker; domains only depend on the Bus type, so swapping it for an actual
// queue later does not touch domain code.
package broker

import (
	"context"
	"log/slog"
	"sync"
)

type Event struct {
	Type     string
	TenantID int64
	Payload  map[string]any
}

type Handler func(ctx context.Context, event Event) error

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
	logger   *slog.Logger
}

func NewBus(logger *slog.Logger) *Bus {
	return &Bus{handlers: make(map[string][]Handler), logger: logger}
}

func (b *Bus) Subscribe(eventType string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], h)
}

// Publish delivers the event synchronously to every subscriber. A failing
// handler is logged, never propagated: notification delivery must not roll
// back the domain operation that raised the event.
func (b *Bus) Publish(ctx context.Context, event Event) {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers[event.Type]...)
	b.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, event); err != nil {
			b.logger.Error("event handler failed", "event_type", event.Type, "error", err)
		}
	}
}
