package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/act"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/address"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/auth"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/consumer"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/meter"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/notification"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/organization"
	platformauth "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/auth"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/middleware"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/observability"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/photo"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/task"
)

type Dependencies struct {
	Logger              *slog.Logger
	Pool                *db.Pool
	TokenService        *platformauth.Service
	AuthHandler         *auth.Handler
	OrgHandler          *organization.Handler
	ConsumerHandler     *consumer.Handler
	AddressHandler      *address.Handler
	TaskHandler         *task.Handler
	ActHandler          *act.Handler
	MeterHandler        *meter.Handler
	PhotoHandler        *photo.Handler
	NotificationHandler *notification.Handler
	CORSOrigins         []string
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

	publicPaths := map[string]struct{}{
		"/api/auth/login":                 {},
		"/api/auth/register-company":      {},
		"/api/auth/bootstrap-super-admin": {},
		"/api/auth/refresh":               {},
		"/api/auth/logout":                {},
		"/health/live":                    {},
		"/health/ready":                   {},
		"/metrics":                        {},
	}

	r.Route("/api", func(api chi.Router) {
		api.Use(middleware.JWTAuth(deps.TokenService, publicPaths))

		api.Route("/auth", func(authRouter chi.Router) {
			authRouter.Mount("/", deps.AuthHandler.Routes())
		})

		api.Route("/users", func(users chi.Router) {
			users.Mount("/", deps.AuthHandler.UserRoutes())
		})

		api.Route("/platform/tenants", func(tenants chi.Router) {
			tenants.Use(middleware.RequireRoles(string(auth.RoleSuperAdmin)))
			tenants.Mount("/", deps.OrgHandler.Routes())
		})

		api.Route("/consumers", func(consumers chi.Router) {
			consumers.Mount("/", deps.ConsumerHandler.Routes())
		})

		api.Route("/addresses", func(addresses chi.Router) {
			addresses.Mount("/", deps.AddressHandler.Routes())
		})

		api.Route("/tasks", func(tasks chi.Router) {
			tasks.Mount("/", deps.TaskHandler.Routes())
		})

		api.Route("/acts", func(acts chi.Router) {
			acts.Mount("/", deps.ActHandler.Routes())
			// Nested under the same {id} segment act's own inspection routes
			// use, so meters ride along as a sub-resource of one act.
			acts.Route("/inspection/{id}/meters", func(meters chi.Router) {
				meters.Mount("/", deps.MeterHandler.Routes())
			})
		})

		api.Route("/photos", func(photos chi.Router) {
			photos.Mount("/", deps.PhotoHandler.Routes())
		})

		api.Route("/notifications", func(notifications chi.Router) {
			notifications.Mount("/", deps.NotificationHandler.Routes())
		})
	})

	return r
}
