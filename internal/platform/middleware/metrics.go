package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// httpRequestsTotal and httpRequestDuration are registered via promauto on
// the default Prometheus registry, so observability.MetricsHandler's plain
// promhttp.Handler() (internal/platform/observability/health.go) picks them
// up automatically -- no wiring needed on that end.
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests, labeled by method, route pattern and status code.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request latency in seconds, labeled by method, route pattern and status code.",
			// 1ms .. 1s: this is a fast in-process API backed by local
			// Postgres/Redis, not an external-call-heavy service.
			Buckets: []float64{0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"method", "path", "status"},
	)
)

// Metrics records http_requests_total/http_request_duration_seconds for
// every request. The path label is chi's resolved route pattern (e.g.
// "/api/requests/{id}"), not the raw URL -- read via chi.RouteContext AFTER
// next.ServeHTTP returns, not before: chi only finishes filling in
// RoutePattern() once routing has fully matched, and router.go mounts
// several sub-routers under /api, each contributing its own segment of the
// pattern as matching descends into it. Reading it before next.ServeHTTP
// would see an empty or partial pattern instead of the final one.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		// "unmatched" (not the raw path) for anything chi never resolved to
		// a registered route -- otherwise a flood of 404s against arbitrary
		// paths would each mint their own label combination.
		pattern := "unmatched"
		if rctx := chi.RouteContext(r.Context()); rctx != nil {
			if p := rctx.RoutePattern(); p != "" {
				pattern = p
			}
		}

		status := strconv.Itoa(ww.Status())
		httpRequestsTotal.WithLabelValues(r.Method, pattern, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, pattern, status).Observe(time.Since(start).Seconds())
	})
}
