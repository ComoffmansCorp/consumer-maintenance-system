package chat

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/broker"
)

// client owns exclusive write access to its own connection -- gorilla's
// websocket.Conn permits at most one concurrent writer, so every outbound
// message (whether a hub broadcast or the connection's own control frames)
// goes through this one channel and the single writer goroutine draining it
// (see writePump in ws_handler.go), rather than calling WriteMessage
// directly from whichever goroutine has a message to send.
type client struct {
	conn *websocket.Conn
	send chan []byte
}

// Hub fans out newly created chat messages to every WebSocket connection
// currently watching a given request's thread -- the real-time replacement
// for the frontend's previous 5s polling loop. It subscribes to the same
// broker.Bus the rest of the domain layer already uses for cross-domain
// effects, so Service.SendMessage doesn't need to know Hub exists at all.
//
// In-process only: on a horizontally scaled `app` (docker compose --scale
// app=N), two clients on the same thread connected to different replicas
// won't see each other's messages pushed live -- each replica only knows
// about the connections it's holding itself. A real deployment would swap
// this for Redis pub/sub (already in the stack for the cache/session-
// revocation roles) as the fan-out layer instead. Documented limitation,
// not a bug: this demo stack runs a single `app` instance.
type Hub struct {
	mu    sync.Mutex
	conns map[int64]map[*client]struct{}
}

func NewHub() *Hub {
	return &Hub{conns: make(map[int64]map[*client]struct{})}
}

// Subscribe wires the hub into the event bus. Called once at startup
// (cmd/api/main.go), alongside payment's RegisterHandlers for the same bus.
func (h *Hub) Subscribe(bus *broker.Bus) {
	bus.Subscribe(broker.EventChatMessageCreated, h.onMessageCreated)
}

func (h *Hub) onMessageCreated(_ context.Context, event broker.Event) error {
	requestID, ok := event.Payload["request_id"].(int64)
	if !ok {
		return nil
	}
	dto, ok := event.Payload["message"].(MessageDTO)
	if !ok {
		return nil
	}
	payload, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	h.broadcast(requestID, payload)
	return nil
}

func (h *Hub) register(requestID int64, c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.conns[requestID] == nil {
		h.conns[requestID] = make(map[*client]struct{})
	}
	h.conns[requestID][c] = struct{}{}
}

func (h *Hub) unregister(requestID int64, c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conns[requestID], c)
	if len(h.conns[requestID]) == 0 {
		delete(h.conns, requestID)
	}
	close(c.send)
}

// broadcast is best-effort per connection: a client whose send buffer is
// already full (a stalled reader on the other end) gets this message
// dropped rather than blocking every other connection on the same thread,
// or blocking SendMessage's own request/response cycle (Bus.Publish runs
// synchronously on the goroutine that created the message).
func (h *Hub) broadcast(requestID int64, payload []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.conns[requestID] {
		select {
		case c.send <- payload:
		default:
		}
	}
}
