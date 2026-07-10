// Command seed populates a fresh database with demo data for manual
// testing: a tenant, staff of every role, consumers, addresses, and tasks
// spanning every status (including a completed one with a filled inspection
// act, meters and a photo). It goes through the same domain services the
// API uses — not raw SQL — so passwords are hashed correctly and every
// business rule (task state machine, act ownership) is exercised exactly as
// in production. Safe to run repeatedly: it skips seeding if the demo
// tenant already exists.
package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log/slog"
	"os"
	"time"

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
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/wiring"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/photo"
	photodb "github.com/myurbondarchuk/consumer-maintenance-system/internal/photo/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/task"
	taskdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/task/db"
)

const (
	tenantCode = "demo"
	password   = "Demo12345"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

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

	orgRepo := organization.NewRepository(orgdb.New(pool))

	authRepo := auth.NewRepository(authdb.New(pool))
	authService := auth.NewService(authRepo, wiring.NewTenantAdapter(orgRepo), tokenService, wiring.NewTxRunnerAdapter(txManager))

	consumerRepo := consumer.NewRepository(consumerdb.New(pool))
	consumerService := consumer.NewService(consumerRepo)

	addressRepo := address.NewRepository(addressdb.New(pool))
	addressService := address.NewService(addressRepo, wiring.NewConsumerAdapter(consumerService))

	actRef := &wiring.ActRef{}

	meterRepo := meter.NewRepository(meterdb.New(pool))
	meterService := meter.NewService(meterRepo, actRef)

	photoRepo := photo.NewRepository(photodb.New(pool))
	photoService := photo.NewService(photoRepo, actRef, cfg.FileUploadDir)

	eventBus := broker.NewBus(logger)

	taskRepo := task.NewRepository(taskdb.New(pool))
	taskService := task.NewService(taskRepo, wiring.NewTaskAddressAdapter(addressService), wiring.NewTaskUserAdapter(authService), actRef, eventBus)

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
	actRef.Set(actService)

	notificationRepo := notification.NewRepository(notificationdb.New(pool))
	notificationService := notification.NewService(notificationRepo, authService)
	notificationService.RegisterHandlers(eventBus)

	if _, err := orgRepo.GetByCode(ctx, tenantCode); err == nil {
		logger.Info("demo tenant already exists, skipping seed", "tenantCode", tenantCode)
		return
	}

	s := &seeder{
		ctx: ctx, logger: logger,
		auth: authService, consumers: consumerService, addresses: addressService,
		tasks: taskService, acts: actService, meters: meterService, photos: photoService,
	}
	if err := s.run(); err != nil {
		logger.Error("seed failed", "error", err)
		os.Exit(1)
	}
}

type seeder struct {
	ctx       context.Context
	logger    *slog.Logger
	auth      *auth.Service
	consumers *consumer.Service
	addresses *address.Service
	tasks     *task.Service
	acts      *act.Service
	meters    *meter.Service
	photos    *photo.Service
}

func (s *seeder) run() error {
	ctx := s.ctx

	authResp, err := s.auth.RegisterCompany(ctx, auth.CompanyRegistrationRequest{
		TenantName: "Энергосбыт Демо",
		TenantCode: tenantCode,
		Plan:       "BUSINESS",
		Username:   "admin",
		Password:   password,
		FullName:   "Анна Смирнова",
	})
	if err != nil {
		return fmt.Errorf("register demo company: %w", err)
	}
	tenantID := *authResp.TenantID
	s.logger.Info("created tenant admin", "username", "admin", "password", password)

	dispatcher, err := s.auth.CreateTenantUser(ctx, tenantID, auth.CreateTenantUserRequest{
		Username: "dispatcher", Password: password, FullName: "Игорь Дежуров", Role: auth.RoleDispatcher,
	})
	if err != nil {
		return fmt.Errorf("create dispatcher: %w", err)
	}
	s.logger.Info("created dispatcher", "username", dispatcher.Username, "password", password)

	electrician1, err := s.auth.CreateTenantUser(ctx, tenantID, auth.CreateTenantUserRequest{
		Username: "inspector1", Password: password, FullName: "Пётр Инспекторов", Role: auth.RoleElectrician,
	})
	if err != nil {
		return fmt.Errorf("create electrician1: %w", err)
	}
	electrician2, err := s.auth.CreateTenantUser(ctx, tenantID, auth.CreateTenantUserRequest{
		Username: "inspector2", Password: password, FullName: "Анна Морозова", Role: auth.RoleElectrician,
	})
	if err != nil {
		return fmt.Errorf("create electrician2: %w", err)
	}
	s.logger.Info("created inspectors", "usernames", []string{electrician1.Username, electrician2.Username}, "password", password)

	consumers := []struct {
		name string
		typ  consumer.Type
		desc string
	}{
		{"ООО Ромашка", consumer.TypeCommercial, "Розничная сеть, 3 магазина"},
		{"МУП Водоканал", consumer.TypeGovernment, "Муниципальное предприятие"},
		{"ТСЖ Северное", consumer.TypeResidential, "Жилой комплекс, 6 корпусов"},
		{"ИП Фермер А.А.", consumer.TypeCommercial, "Фермерское хозяйство"},
	}
	consumerIDs := make([]int64, 0, len(consumers))
	for _, c := range consumers {
		dto, err := s.consumers.Create(ctx, tenantID, consumer.CreateRequest{Name: c.name, Type: c.typ, Description: c.desc})
		if err != nil {
			return fmt.Errorf("create consumer %s: %w", c.name, err)
		}
		consumerIDs = append(consumerIDs, dto.ID)
	}
	s.logger.Info("created consumers", "count", len(consumerIDs))

	addresses := []struct {
		street, house, building, apartment string
		consumerIdx                        int
	}{
		{"ул. Ленина", "10", "", "5", 0},
		{"ул. Ленина", "10", "", "12", 0},
		{"пр. Мира", "45", "2", "", 1},
		{"ул. Садовая", "3", "", "1", 2},
		{"ул. Садовая", "3", "", "14", 2},
		{"пер. Полевой", "7", "", "", 3},
	}
	addressIDs := make([]int64, 0, len(addresses))
	for _, a := range addresses {
		cid := consumerIDs[a.consumerIdx]
		dto, err := s.addresses.Create(ctx, tenantID, address.CreateRequest{
			Street: a.street, House: a.house, Building: a.building, Apartment: a.apartment, ConsumerID: &cid,
		})
		if err != nil {
			return fmt.Errorf("create address %s %s: %w", a.street, a.house, err)
		}
		addressIDs = append(addressIDs, dto.ID)
	}
	s.logger.Info("created addresses", "count", len(addressIDs))

	dueDate := httpx.Date{Time: mustParseDate("2026-07-20")}

	// Two unassigned pending tasks — dispatcher's queue.
	if _, err := s.tasks.Create(ctx, tenantID, task.CreateRequest{Type: task.TypeInspection, AddressID: addressIDs[3], DueDate: &dueDate}); err != nil {
		return fmt.Errorf("create pending task 1: %w", err)
	}
	if _, err := s.tasks.Create(ctx, tenantID, task.CreateRequest{Type: task.TypeReplacement, AddressID: addressIDs[4], DueDate: &dueDate}); err != nil {
		return fmt.Errorf("create pending task 2: %w", err)
	}

	// Assigned-but-not-started task.
	assignedID := electrician2.ID
	pendingAssigned, err := s.tasks.Create(ctx, tenantID, task.CreateRequest{
		Type: task.TypeInspection, AddressID: addressIDs[5], DueDate: &dueDate, AssigneeID: &assignedID,
	})
	if err != nil {
		return fmt.Errorf("create assigned task: %w", err)
	}
	_ = pendingAssigned

	// In-progress task (electrician1 took it, hasn't filled the act yet).
	inProgress, err := s.tasks.Create(ctx, tenantID, task.CreateRequest{Type: task.TypeReplacement, AddressID: addressIDs[1], DueDate: &dueDate})
	if err != nil {
		return fmt.Errorf("create in-progress task: %w", err)
	}
	if _, err := s.tasks.Assign(ctx, inProgress.ID, tenantID, electrician1.ID); err != nil {
		return fmt.Errorf("assign in-progress task: %w", err)
	}
	if _, err := s.tasks.Start(ctx, inProgress.ID, tenantID, electrician1.ID); err != nil {
		return fmt.Errorf("start in-progress task: %w", err)
	}

	// Canceled task.
	canceled, err := s.tasks.Create(ctx, tenantID, task.CreateRequest{Type: task.TypeInspection, AddressID: addressIDs[2], DueDate: &dueDate})
	if err != nil {
		return fmt.Errorf("create canceled task: %w", err)
	}
	if _, err := s.tasks.Cancel(ctx, canceled.ID, tenantID, "Потребитель перенёс визит на следующий месяц"); err != nil {
		return fmt.Errorf("cancel task: %w", err)
	}

	// Completed inspection task with a filled act, a meter and a photo.
	if err := s.seedCompletedInspection(tenantID, electrician1.ID, consumerIDs[0], addressIDs[0], dueDate); err != nil {
		return fmt.Errorf("seed completed inspection: %w", err)
	}

	// Completed replacement task with a filled act and a photo.
	if err := s.seedCompletedReplacement(tenantID, electrician2.ID, addressIDs[3], dueDate); err != nil {
		return fmt.Errorf("seed completed replacement: %w", err)
	}

	s.logger.Info("seed complete", "tenantCode", tenantCode)
	s.logger.Info("login with", "tenantCode", tenantCode, "username", "admin/dispatcher/inspector1/inspector2", "password", password)
	return nil
}

func (s *seeder) seedCompletedInspection(tenantID, inspectorID, consumerID, addressID int64, dueDate httpx.Date) error {
	ctx := s.ctx
	t, err := s.tasks.Create(ctx, tenantID, task.CreateRequest{Type: task.TypeInspection, AddressID: addressID, DueDate: &dueDate})
	if err != nil {
		return err
	}
	if _, err := s.tasks.Assign(ctx, t.ID, tenantID, inspectorID); err != nil {
		return err
	}
	if _, err := s.tasks.Start(ctx, t.ID, tenantID, inspectorID); err != nil {
		return err
	}

	inspectionDate := httpx.Date{Time: mustParseDate("2026-07-05")}
	consumerIDCopy := consumerID
	actDTO, err := s.acts.CreateInspection(ctx, tenantID, inspectorID, act.CreateInspectionRequest{
		TaskID:         t.ID,
		InspectionDate: &inspectionDate,
		ConsumerID:     &consumerIDCopy,
		InspectionType: act.InspectionScheduled,
		Notes:          "Плановый осмотр, показания сняты, пломба не нарушена.",
	})
	if err != nil {
		return err
	}

	verification := httpx.Date{Time: mustParseDate("2025-03-12")}
	year := int32(2019)
	if _, err := s.meters.Create(ctx, actDTO.ID, tenantID, meter.CreateRequest{
		Type: meter.TypeSinglePhase, SerialNumber: "SN-778812", ManufactureYear: &year,
		VerificationDate: &verification, SealState: meter.SealIntact,
	}); err != nil {
		return err
	}

	if err := s.attachSamplePhoto(tenantID, inspectorID, actDTO.ID, nil, "Показания на момент осмотра"); err != nil {
		return err
	}

	if _, err := s.tasks.Complete(ctx, t.ID, tenantID, inspectorID); err != nil {
		return err
	}
	return nil
}

func (s *seeder) seedCompletedReplacement(tenantID, inspectorID, addressID int64, dueDate httpx.Date) error {
	ctx := s.ctx
	t, err := s.tasks.Create(ctx, tenantID, task.CreateRequest{Type: task.TypeReplacement, AddressID: addressID, DueDate: &dueDate})
	if err != nil {
		return err
	}
	if _, err := s.tasks.Assign(ctx, t.ID, tenantID, inspectorID); err != nil {
		return err
	}
	if _, err := s.tasks.Start(ctx, t.ID, tenantID, inspectorID); err != nil {
		return err
	}

	installDate := httpx.Date{Time: mustParseDate("2026-07-02")}
	oldReadings := 4821.5
	newReadings := 0.0
	actDTO, err := s.acts.CreateReplacement(ctx, tenantID, inspectorID, act.CreateReplacementRequest{
		TaskID:           t.ID,
		AccountNumber:    "AC-100245",
		InstallationDate: &installDate,
		OldBrand:         "Меркурий 201",
		OldSerialNumber:  "M201-55231",
		OldReadings:      &oldReadings,
		NewBrand:         "Меркурий 230",
		NewSerialNumber:  "M230-90112",
		NewReadings:      &newReadings,
	})
	if err != nil {
		return err
	}

	replActID := actDTO.ID
	if err := s.attachSamplePhoto(tenantID, inspectorID, 0, &replActID, "Новый прибор после установки"); err != nil {
		return err
	}

	if _, err := s.tasks.Complete(ctx, t.ID, tenantID, inspectorID); err != nil {
		return err
	}
	return nil
}

func (s *seeder) attachSamplePhoto(tenantID, uploaderID, inspectionActID int64, replacementActID *int64, note string) error {
	img := generateSamplePhoto()
	var inspID *int64
	if replacementActID == nil {
		inspID = &inspectionActID
	}
	_, err := s.photos.Upload(s.ctx, photo.UploadInput{
		TenantID: tenantID, UploadedBy: uploaderID,
		InspectionActID: inspID, ReplacementActID: replacementActID,
		Note: note, OriginalFilename: "meter.jpg", ContentType: "image/jpeg",
		Size: int64(len(img)), Content: bytes.NewReader(img),
	})
	return err
}

// generateSamplePhoto renders a small placeholder JPEG (a stylized meter
// display) so seeded acts have something visible in the photo gallery
// without shipping a binary asset in the repo.
func generateSamplePhoto() []byte {
	const w, h = 480, 360
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	bg := color.RGBA{30, 41, 59, 255}   // ink
	panel := color.RGBA{13, 148, 136, 255} // primary
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, bg)
		}
	}
	for y := h/2 - 40; y < h/2+40; y++ {
		for x := 40; x < w-40; x++ {
			img.Set(x, y, panel)
		}
	}
	buf := &bytes.Buffer{}
	_ = jpeg.Encode(buf, img, &jpeg.Options{Quality: 80})
	return buf.Bytes()
}

func mustParseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}
