package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/act"
	actdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/act/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/address"
	addressdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/address/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/auth"
	authdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/auth/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/consumer"
	consumerdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/consumer/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/meter"
	meterdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/meter/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/notification"
	notificationdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/notification/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/organization"
	orgdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/organization/db"
	platformauth "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/auth"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/broker"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/config"
	platformdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/wiring"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/photo"
	photodb "github.com/myurbondarchuk/consumer-maintenance-system/internal/photo/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/server"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/task"
	taskdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/task/db"
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

	txManager := platformdb.NewTxManager(pool)

	tokenService, err := platformauth.NewService(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	if err != nil {
		logger.Error("init jwt service", "error", err)
		os.Exit(1) //nolint:gocritic // fatal startup error; process exit reclaims the pool anyway
	}

	// --- organization (platform tenant management) ---
	orgRepo := organization.NewRepository(orgdb.New(pool))
	orgService := organization.NewService(orgRepo)
	orgHandler := organization.NewHandler(orgService)

	// --- auth ---
	authRepo := auth.NewRepository(authdb.New(pool))
	authService := auth.NewService(
		authRepo,
		wiring.NewTenantAdapter(orgRepo),
		tokenService,
		wiring.NewTxRunnerAdapter(txManager),
	)
	authHandler := auth.NewHandler(authService)

	// --- consumer ---
	consumerRepo := consumer.NewRepository(consumerdb.New(pool))
	consumerService := consumer.NewService(consumerRepo)
	consumerHandler := consumer.NewHandler(consumerService)

	// --- address ---
	addressRepo := address.NewRepository(addressdb.New(pool))
	addressService := address.NewService(addressRepo, wiring.NewConsumerAdapter(consumerService))
	addressHandler := address.NewHandler(addressService)

	// actRef breaks the act <-> {task, meter, photo} construction cycle: it
	// is handed out now and populated once act.Service exists below, before
	// the HTTP server starts accepting requests.
	actRef := &wiring.ActRef{}

	// --- meter ---
	meterRepo := meter.NewRepository(meterdb.New(pool))
	meterService := meter.NewService(meterRepo, actRef)
	meterHandler := meter.NewHandler(meterService)

	// --- photo ---
	photoRepo := photo.NewRepository(photodb.New(pool))
	photoService := photo.NewService(photoRepo, actRef, cfg.FileUploadDir)
	photoHandler := photo.NewHandler(photoService)

	// --- events ---
	eventBus := broker.NewBus(logger)

	// --- task ---
	taskRepo := task.NewRepository(taskdb.New(pool))
	taskService := task.NewService(
		taskRepo,
		wiring.NewTaskAddressAdapter(addressService),
		wiring.NewTaskUserAdapter(authService),
		actRef,
		eventBus,
	)
	taskHandler := task.NewHandler(taskService)

	// --- act ---
	actRepo := act.NewRepository(actdb.New(pool))
	actService := act.NewService(
		actRepo,
		wiring.NewActTaskAdapter(taskService),
		wiring.NewActAddressAdapter(addressService),
		wiring.NewActConsumerAdapter(consumerService),
		wiring.NewActMeterAdapter(meterService),
		wiring.NewActPhotoAdapter(photoService),
		wiring.NewActUserAdapter(authService),
		wiring.NewActTenantAdapter(orgRepo),
	)
	actHandler := act.NewHandler(actService)
	actRef.Set(actService)

	// --- notification ---
	notificationRepo := notification.NewRepository(notificationdb.New(pool))
	notificationService := notification.NewService(notificationRepo, authService)
	notificationService.RegisterHandlers(eventBus)
	notificationHandler := notification.NewHandler(notificationService)

	router := server.NewRouter(server.Dependencies{
		Logger:              logger,
		Pool:                pool,
		TokenService:        tokenService,
		AuthHandler:         authHandler,
		OrgHandler:          orgHandler,
		ConsumerHandler:     consumerHandler,
		AddressHandler:      addressHandler,
		TaskHandler:         taskHandler,
		ActHandler:          actHandler,
		MeterHandler:        meterHandler,
		PhotoHandler:        photoHandler,
		NotificationHandler: notificationHandler,
		CORSOrigins:         cfg.CORSAllowedOrigins,
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
