package chat

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	platformauth "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/auth"
	platformcache "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/cache"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/middleware"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = pongWait * 9 / 10
)

// upgrader.CheckOrigin defaults to same-origin only, which would reject the
// browser client here: the SPA is served by its own nginx (frontend/nginx.conf,
// port 5173) while this WebSocket terminates on the gateway (:8000) -- two
// different origins by design (see docker-compose.yml), same as every plain
// HTTP API call already crossing that boundary via CORS. Allowing all
// origins here is the WebSocket-handshake equivalent of the CORS
// middleware every other route already goes through.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WSHandler holds what the WebSocket handshake needs that the rest of
// Handler doesn't: token auth happens here via a query param (browsers
// can't set a custom Authorization header on a WebSocket upgrade request),
// not the JWTAuth middleware, so this route is registered outside it
// (see server/router.go) and authenticates itself.
type WSHandler struct {
	service     *Service
	hub         *Hub
	tokens      *platformauth.Service
	cacheClient *platformcache.Client
}

func NewWSHandler(service *Service, hub *Hub, tokens *platformauth.Service, cacheClient *platformcache.Client) *WSHandler {
	return &WSHandler{service: service, hub: hub, tokens: tokens, cacheClient: cacheClient}
}

func (h *WSHandler) Serve(w http.ResponseWriter, r *http.Request) {
	requestID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}

	auth, err := middleware.AuthenticateToken(r.Context(), h.tokens, h.cacheClient, r.URL.Query().Get("token"))
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Same participant + assigned-only rule as the REST endpoints -- a
	// WebSocket is just another way to reach the same authorized thread,
	// not a separate access path with its own rules.
	if _, err := h.service.ListMessages(r.Context(), requestID, auth.UserID, 0); err != nil {
		status := http.StatusForbidden
		if errors.Is(err, ErrRequestNotFound) {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	c := &client{conn: conn, send: make(chan []byte, 16)}
	h.hub.register(requestID, c)

	go writePump(c)
	readPump(h.hub, requestID, c)
}

// readPump's only real job is noticing disconnects: this hub is push-only
// (the frontend still POSTs new messages over plain REST, see
// handler.go's SendMessage), it doesn't accept chat text over the socket.
// Blocking on ReadMessage is what makes a closed connection observable at
// all, so this loop has to keep calling it even though it discards
// whatever comes back.
func readPump(hub *Hub, requestID int64, c *client) {
	defer func() {
		hub.unregister(requestID, c)
		_ = c.conn.Close()
	}()
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

// writePump is the single writer for this connection (see the client
// struct's doc comment for why) -- both hub broadcasts and the periodic
// ping go through here, never a direct WriteMessage call from elsewhere.
func writePump(c *client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case payload, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
