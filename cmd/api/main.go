package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/auth"
	authdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/auth/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/catalog"
	catalogdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/catalog/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/chat"
	chatdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/chat/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/master"
	masterdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/master/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/payment"
	paymentdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/payment/db"
	platformauth "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/auth"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/broker"
	platformcache "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/cache"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/config"
	platformdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/wiring"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/request"
	requestdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/request/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/review"
	reviewdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/review/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := platformdb.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("connect database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Redis is a soft dependency: session revocation and the catalog cache
	// both degrade gracefully to "no cache" (see cache.Client usage in
	// internal/auth and internal/catalog) rather than requiring it to be up,
	// so a connect failure here logs and continues with a nil client instead
	// of os.Exit like the Postgres pool above.
	cacheClient, err := platformcache.Connect(ctx, cfg.RedisURL)
	if err != nil {
		logger.Error("connect redis, continuing without cache", "error", err)
		cacheClient = nil
	} else {
		defer func() { _ = cacheClient.Close() }()
	}

	txManager := platformdb.NewTxManager(pool)

	tokenService, err := platformauth.NewService(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	if err != nil {
		logger.Error("init jwt service", "error", err)
		os.Exit(1) //nolint:gocritic // fatal startup error; process exit reclaims the pool anyway
	}

	// --- auth (foundational, no dependencies on other domains) ---
	authRepo := auth.NewRepository(authdb.New(pool))
	authService := auth.NewService(authRepo, tokenService, cacheClient)
	authHandler := auth.NewHandler(authService)

	// --- catalog (categories + services, no dependencies) ---
	catalogRepo := catalog.NewRepository(catalogdb.New(pool))
	catalogService := catalog.NewService(catalogRepo, cacheClient, logger)
	catalogHandler := catalog.NewHandler(catalogService)
	catalogAdapter := wiring.NewCatalogAdapter(catalogService)

	// --- master (profiles + specializations, depends on catalog) ---
	masterRepo := master.NewRepository(masterdb.New(pool))
	masterService := master.NewService(masterRepo, catalogAdapter)
	masterHandler := master.NewHandler(masterService)
	masterAdapter := wiring.NewMasterAdapter(masterService)

	// --- events (request -> payment cross-domain effects) ---
	eventBus := broker.NewBus(logger)

	// --- request (service_requests + offers + history + favorites,
	// depends on catalog and master) ---
	requestRepo := request.NewRepository(requestdb.New(pool))
	requestService := request.NewService(requestRepo, catalogAdapter, masterAdapter, wiring.NewTxRunnerAdapter(txManager), eventBus)
	requestHandler := request.NewHandler(requestService)
	requestAdapter := wiring.NewRequestAdapter(requestService)

	// --- review (depends on request and master) ---
	reviewRepo := review.NewRepository(reviewdb.New(pool))
	reviewService := review.NewService(reviewRepo, requestAdapter, masterAdapter)
	reviewHandler := review.NewHandler(reviewService)

	// --- payment (event-driven off the request lifecycle, depends on
	// request only for the read-side participant check) ---
	paymentRepo := payment.NewRepository(paymentdb.New(pool))
	paymentService := payment.NewService(paymentRepo, requestAdapter)
	paymentService.RegisterHandlers(eventBus)
	paymentHandler := payment.NewHandler(paymentService)

	// --- chat (depends on request for participant/assignment checks) ---
	chatRepo := chat.NewRepository(chatdb.New(pool))
	chatService := chat.NewService(chatRepo, requestAdapter)
	chatHandler := chat.NewHandler(chatService)

	router := server.NewRouter(server.Dependencies{
		Logger:         logger,
		Pool:           pool,
		TokenService:   tokenService,
		CacheClient:    cacheClient,
		AuthHandler:    authHandler,
		CatalogHandler: catalogHandler,
		MasterHandler:  masterHandler,
		RequestHandler: requestHandler,
		ReviewHandler:  reviewHandler,
		PaymentHandler: paymentHandler,
		ChatHandler:    chatHandler,
		CORSOrigins:    cfg.CORSAllowedOrigins,
	})

	httpServer := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	go func() {
		logger.Info("starting api server", "addr", cfg.HTTPAddr, "env", cfg.AppEnv)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped", "error", err)
			os.Exit(1)
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-stopCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
	logger.Info("server stopped")
}
