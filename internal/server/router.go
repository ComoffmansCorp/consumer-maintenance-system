package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/auth"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/catalog"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/chat"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/master"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/payment"
	platformauth "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/auth"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/middleware"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/observability"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/request"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/review"
)

type Dependencies struct {
	Logger         *slog.Logger
	Pool           *db.Pool
	TokenService   *platformauth.Service
	AuthHandler    *auth.Handler
	CatalogHandler *catalog.Handler
	MasterHandler  *master.Handler
	RequestHandler *request.Handler
	ReviewHandler  *review.Handler
	PaymentHandler *payment.Handler
	ChatHandler    *chat.Handler
	CORSOrigins    []string
}

func NewRouter(deps Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recover(deps.Logger))
	r.Use(middleware.Logger(deps.Logger))
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(middleware.CORS(deps.CORSOrigins))

	health := observability.NewHealth(deps.Pool)
	r.Get("/health/live", health.Liveness)
	r.Get("/health/ready", health.Readiness)
	r.Handle("/metrics", observability.MetricsHandler())

	// GET /api/masters/{id}/reviews is public but carries a path parameter,
	// which JWTAuth's publicPaths map (exact-string match only, see below)
	// cannot express. Registered directly here, outside the /api group's
	// JWTAuth middleware entirely -- the same pattern already used for
	// health/metrics above.
	r.Get("/api/masters/{id}/reviews", deps.ReviewHandler.ListMasterReviews)

	publicPaths := map[string]struct{}{
		"/api/auth/login":                 {},
		"/api/auth/register-client":       {},
		"/api/auth/register-master":       {},
		"/api/auth/bootstrap-super-admin": {},
		"/api/auth/refresh":               {},
		"/api/auth/logout":                {},
		// Catalog browsing is public so the marketplace landing page can
		// show services before a visitor registers or logs in. GET here is
		// public -- JWTAuth matches by exact path only, not method, so
		// admin writes MUST live on a different path (/admin/...) rather
		// than POST on the same "/categories"/"/services" path, or they
		// would accidentally become unauthenticated too.
		"/api/catalog/categories": {},
		"/api/catalog/services":   {},
		"/health/live":            {},
		"/health/ready":           {},
		"/metrics":                {},
	}

	r.Route("/api", func(api chi.Router) {
		api.Use(middleware.JWTAuth(deps.TokenService, publicPaths))

		api.Route("/auth", func(authRouter chi.Router) {
			authRouter.Mount("/", deps.AuthHandler.Routes())
		})

		api.Route("/catalog", func(catalogRouter chi.Router) {
			catalogRouter.Mount("/", deps.CatalogHandler.Routes())
		})

		api.Route("/master", func(masterRouter chi.Router) {
			masterRouter.Mount("/", deps.MasterHandler.Routes())
		})

		api.Route("/requests", func(requestRouter chi.Router) {
			requestRouter.Mount("/", deps.RequestHandler.Routes())
		})

		api.Route("/reviews", func(reviewRouter chi.Router) {
			reviewRouter.Mount("/", deps.ReviewHandler.Routes())
		})

		// PaymentHandler and ChatHandler each mount their own
		// "/requests/{id}/..." sub-resource, so they ride directly under
		// /api rather than under /api/requests -- the request domain's own
		// router already owns everything at /api/requests/{id}/* that
		// belongs to it.
		api.Mount("/", deps.PaymentHandler.Routes())
		api.Mount("/", deps.ChatHandler.Routes())

		// Cross-cutting admin dashboards live under one shared /api/admin
		// namespace rather than nested inside each domain's own mount --
		// each sub-router still gates on SUPER_ADMIN itself.
		api.Route("/admin", func(adminRouter chi.Router) {
			adminRouter.Mount("/", deps.CatalogHandler.AdminRoutes())
			adminRouter.Mount("/masters", deps.MasterHandler.AdminRoutes())
			adminRouter.Mount("/requests", deps.RequestHandler.AdminRoutes())
			adminRouter.Mount("/reviews", deps.ReviewHandler.AdminRoutes())
			adminRouter.Mount("/payments", deps.PaymentHandler.AdminRoutes())
		})
	})

	return r
}
