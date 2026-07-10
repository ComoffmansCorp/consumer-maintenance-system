package observability

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Health struct {
	pool *db.Pool
}

func NewHealth(pool *db.Pool) *Health {
	return &Health{pool: pool}
}

func (h *Health) Liveness(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Health) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.pool.Ping(ctx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "not_ready",
			"detail": "database unavailable",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
