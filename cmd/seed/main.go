// Command seed populates a fresh database with demo data for manual
// testing: an admin, a broad service catalog, a couple of masters with
// specializations, a couple of clients, and requests spanning every status
// and every new feature -- an OPEN request with PENDING offers, a fully
// walked-through COMPLETED request with a RELEASED payment, a review and
// messages, and a CANCELED request with a REFUNDED payment. It goes through
// the same domain services the API uses -- not raw SQL -- so passwords are
// hashed correctly and every business rule (specialization gate, status
// transitions) is exercised exactly as in production. Safe to run
// repeatedly: it skips seeding if the catalog already exists.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/config"
	platformdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/wiring"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/request"
	requestdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/request/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/review"
	reviewdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/review/db"
)

const password = "Demo12345"

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

	authRepo := auth.NewRepository(authdb.New(pool))
	authService := auth.NewService(authRepo, tokenService)

	catalogRepo := catalog.NewRepository(catalogdb.New(pool))
	catalogService := catalog.NewService(catalogRepo)
	catalogAdapter := wiring.NewCatalogAdapter(catalogService)

	masterRepo := master.NewRepository(masterdb.New(pool))
	masterService := master.NewService(masterRepo, catalogAdapter)
	masterAdapter := wiring.NewMasterAdapter(masterService)

	eventBus := broker.NewBus(logger)

	requestRepo := request.NewRepository(requestdb.New(pool))
	requestService := request.NewService(requestRepo, catalogAdapter, masterAdapter, wiring.NewTxRunnerAdapter(txManager), eventBus)
	requestAdapter := wiring.NewRequestAdapter(requestService)

	reviewRepo := review.NewRepository(reviewdb.New(pool))
	reviewService := review.NewService(reviewRepo, requestAdapter, masterAdapter)

	paymentRepo := payment.NewRepository(paymentdb.New(pool))
	paymentService := payment.NewService(paymentRepo, requestAdapter)
	paymentService.RegisterHandlers(eventBus)

	chatRepo := chat.NewRepository(chatdb.New(pool))
	chatService := chat.NewService(chatRepo, requestAdapter)

	s := &seeder{
		ctx: ctx, logger: logger,
		auth: authService, catalog: catalogService, master: masterService,
		request: requestService, review: reviewService, payment: paymentService, chat: chatService,
	}

	if err := s.run(); err != nil {
		logger.Error("seed failed", "error", err)
		os.Exit(1)
	}
}

type seeder struct {
	ctx     context.Context
	logger  *slog.Logger
	auth    *auth.Service
	catalog *catalog.Service
	master  *master.Service
	request *request.Service
	review  *review.Service
	payment *payment.Service
	chat    *chat.Service
}

type serviceSeed struct {
	name        string
	description string
	priceFrom   float64
	priceTo     float64
	unit        string
}

type categorySeed struct {
	name          string
	services      []serviceSeed
	subcategories []categorySeed
}

// marketplaceCatalog is a broad, realistic catalog spanning the kind of
// activities an actual home-services marketplace covers. Deliberately a
// data table, not dozens of individual CreateService calls, so it stays
// readable. One level of subcategory nesting, matching the schema.
var marketplaceCatalog = []categorySeed{
	{name: "Приборы учёта", subcategories: []categorySeed{
		{name: "Электросчётчики", services: []serviceSeed{
			{"Осмотр счётчика", "Плановый осмотр прибора учёта", 500, 1000, "визит"},
			{"Замена счётчика", "Замена прибора учёта на новый", 1500, 3000, "шт"},
			{"Поверка счётчика", "Проверка точности показаний", 800, 1500, "шт"},
		}},
		{name: "Водосчётчики", services: []serviceSeed{
			{"Опломбировка счётчика", "Установка контрольной пломбы", 300, 600, "шт"},
			{"Установка нового счётчика", "Монтаж прибора учёта с нуля", 1200, 2500, "шт"},
		}},
	}},
	{name: "Сантехника", subcategories: []categorySeed{
		{name: "Ремонт", services: []serviceSeed{
			{"Устранение течи", "Устранение течи в трубах или сантехнике", 800, 2500, "визит"},
			{"Прочистка канализации", "Устранение засора", 1000, 3000, "визит"},
			{"Замена труб", "Замена водопроводных или канализационных труб", 2000, 8000, "п.м."},
		}},
		{name: "Установка", services: []serviceSeed{
			{"Установка смесителя", "Замена или монтаж смесителя", 1000, 2000, "шт"},
			{"Установка унитаза", "Демонтаж старого, монтаж нового", 2000, 4000, "шт"},
			{"Установка водонагревателя", "Монтаж и подключение бойлера", 2500, 5000, "шт"},
		}},
	}},
	{name: "Электрика", services: []serviceSeed{
		{"Замена проводки", "Полная или частичная замена электропроводки", 5000, 20000, "объект"},
		{"Установка розетки", "Монтаж новой розетки или блока розеток", 500, 1200, "шт"},
		{"Установка люстры", "Монтаж и подключение светильника", 800, 2000, "шт"},
		{"Устранение короткого замыкания", "Поиск и устранение неисправности", 1500, 4000, "визит"},
	}},
	{name: "Клининг", services: []serviceSeed{
		{"Генеральная уборка квартиры", "Полная уборка всех помещений", 3000, 8000, "объект"},
		{"Уборка после ремонта", "Удаление строительной пыли и мусора", 4000, 10000, "объект"},
		{"Мытьё окон", "Мытьё окон и рам, включая труднодоступные", 1500, 4000, "объект"},
	}},
	{name: "Бытовая техника", services: []serviceSeed{
		{"Ремонт стиральной машины", "Диагностика и ремонт на месте", 1000, 5000, "визит"},
		{"Ремонт холодильника", "Диагностика и ремонт на месте", 1200, 6000, "визит"},
		{"Установка кондиционера", "Монтаж внутреннего и внешнего блоков", 4000, 12000, "шт"},
	}},
	{name: "Переезды и грузоперевозки", services: []serviceSeed{
		{"Квартирный переезд", "Перевозка вещей при переезде", 3000, 15000, "заказ"},
		{"Грузчики", "Погрузочно-разгрузочные работы", 500, 1000, "час"},
		{"Вывоз мусора", "Вывоз строительного и бытового мусора", 1500, 5000, "заказ"},
	}},
}

func (s *seeder) run() error {
	ctx := s.ctx

	existing, err := s.catalog.ListCategories(ctx)
	if err != nil {
		return fmt.Errorf("check existing catalog: %w", err)
	}
	if len(existing) > 0 {
		s.logger.Info("catalog already exists, skipping seed")
		return nil
	}

	adminResp, err := s.auth.BootstrapSuperAdmin(ctx, auth.BootstrapSuperAdminRequest{
		Username: "admin", Password: password, FullName: "Администратор Платформы",
	})
	if err != nil {
		return fmt.Errorf("bootstrap super admin: %w", err)
	}
	s.logger.Info("created super admin", "username", "admin", "password", password, "userId", adminResp.UserID)

	serviceIDsByName := make(map[string]int64)
	totalServices, totalCategories := 0, 0
	for _, cat := range marketplaceCatalog {
		if err := s.seedCategory(cat, nil, serviceIDsByName, &totalCategories, &totalServices); err != nil {
			return err
		}
	}
	s.logger.Info("seeded catalog", "categories", totalCategories, "services", totalServices)

	master1, err := s.auth.RegisterMaster(ctx, auth.RegisterRequest{Username: "master1", Password: password, FullName: "Виктор Мастеров"})
	if err != nil {
		return fmt.Errorf("register master1: %w", err)
	}
	master2, err := s.auth.RegisterMaster(ctx, auth.RegisterRequest{Username: "master2", Password: password, FullName: "Елена Электрикова"})
	if err != nil {
		return fmt.Errorf("register master2: %w", err)
	}
	client1, err := s.auth.RegisterClient(ctx, auth.RegisterRequest{Username: "client1", Password: password, FullName: "Ольга Заказчикова"})
	if err != nil {
		return fmt.Errorf("register client1: %w", err)
	}
	client2, err := s.auth.RegisterClient(ctx, auth.RegisterRequest{Username: "client2", Password: password, FullName: "Дмитрий Клиентов"})
	if err != nil {
		return fmt.Errorf("register client2: %w", err)
	}
	s.logger.Info("created marketplace demo users", "usernames", []string{"master1", "master2", "client1", "client2"}, "password", password)

	// master1 specializes in meter services, master2 in electrics -- two
	// disjoint specialties so the hard server-side specialization gate is
	// exercised on offer submission.
	meterServiceID := serviceIDsByName["Осмотр счётчика"]
	meterReplacementID := serviceIDsByName["Замена счётчика"]
	electricServiceID := serviceIDsByName["Установка розетки"]
	plumbingServiceID := serviceIDsByName["Устранение течи"]

	if _, err := s.master.UpdateProfile(ctx, master1.UserID, master.UpdateProfileRequest{
		City: "Самара", Bio: "10 лет опыта обслуживания приборов учёта",
		SpecializationIDs: []int64{meterServiceID, meterReplacementID},
	}); err != nil {
		return fmt.Errorf("update master1 profile: %w", err)
	}
	if _, err := s.master.UpdateProfile(ctx, master2.UserID, master.UpdateProfileRequest{
		City: "Самара", Bio: "Электрик, 7 лет практики",
		SpecializationIDs: []int64{electricServiceID},
	}); err != nil {
		return fmt.Errorf("update master2 profile: %w", err)
	}

	// --- OPEN request with a couple of PENDING offers from both masters ---
	openReq, err := s.request.CreateRequest(ctx, client1.UserID, request.CreateRequestRequest{
		ServiceID: meterServiceID, Description: "Счётчик мигает, подозрение на неисправность",
		AddressText: "г. Самара, ул. Ленина, д. 10, кв. 5",
	})
	if err != nil {
		return fmt.Errorf("create open request: %w", err)
	}
	if _, err := s.request.SubmitOffer(ctx, openReq.ID, master1.UserID, request.SubmitOfferRequest{
		Price: 900, Comment: "Могу подъехать завтра утром",
	}); err != nil {
		return fmt.Errorf("submit offer 1 on open request: %w", err)
	}

	// A second request in the same specialization, so a second PENDING
	// offer exists in the pool without violating the one-offer-per-master
	// unique constraint on the first request.
	openReq2, err := s.request.CreateRequest(ctx, client2.UserID, request.CreateRequestRequest{
		ServiceID: meterServiceID, Description: "Нужен плановый осмотр перед продажей квартиры",
		AddressText: "г. Самара, ул. Гагарина, д. 22, кв. 41",
	})
	if err != nil {
		return fmt.Errorf("create open request 2: %w", err)
	}
	if _, err := s.request.SubmitOffer(ctx, openReq2.ID, master1.UserID, request.SubmitOfferRequest{
		Price: 700, Comment: "Свободен в эти выходные",
	}); err != nil {
		return fmt.Errorf("submit offer on open request 2: %w", err)
	}

	// A plumbing request neither seeded master can bid on -- master1 is
	// specialized in meters, master2 in electrics -- demonstrating the
	// hard specialization gate: it stays OPEN with zero offers.
	if _, err := s.request.CreateRequest(ctx, client1.UserID, request.CreateRequestRequest{
		ServiceID: plumbingServiceID, Description: "Течёт труба под раковиной на кухне",
		AddressText: "г. Самара, ул. Ленина, д. 10, кв. 5",
	}); err != nil {
		return fmt.Errorf("create unclaimable plumbing request: %w", err)
	}

	// --- fully completed lifecycle: created -> offer -> accepted ->
	// assigned -> completed, with a RELEASED payment, a review and chat
	// messages ---
	completedReq, err := s.request.CreateRequest(ctx, client2.UserID, request.CreateRequestRequest{
		ServiceID: electricServiceID, Description: "Нужно установить блок розеток на кухне",
		AddressText: "г. Самара, ул. Молодогвардейская, д. 15, кв. 3",
	})
	if err != nil {
		return fmt.Errorf("create completed request: %w", err)
	}
	offer, err := s.request.SubmitOffer(ctx, completedReq.ID, master2.UserID, request.SubmitOfferRequest{
		Price: 1100, Comment: "Есть все материалы с собой",
	})
	if err != nil {
		return fmt.Errorf("submit offer on completed request: %w", err)
	}
	if _, err := s.request.AcceptOffer(ctx, completedReq.ID, offer.ID, client2.UserID); err != nil {
		return fmt.Errorf("accept offer on completed request: %w", err)
	}
	if _, err := s.chat.SendMessage(ctx, completedReq.ID, client2.UserID, chat.SendMessageRequest{Text: "Здравствуйте! Когда сможете подъехать?"}); err != nil {
		return fmt.Errorf("send client message: %w", err)
	}
	if _, err := s.chat.SendMessage(ctx, completedReq.ID, master2.UserID, chat.SendMessageRequest{Text: "Добрый день! Буду завтра к 15:00."}); err != nil {
		return fmt.Errorf("send master message: %w", err)
	}
	if _, err := s.request.CompleteRequest(ctx, completedReq.ID, master2.UserID); err != nil {
		return fmt.Errorf("complete request: %w", err)
	}
	if _, err := s.review.CreateReview(ctx, client2.UserID, review.CreateReviewRequest{
		RequestID: completedReq.ID, Rating: 5, Comment: "Отличная работа, всё аккуратно и быстро!",
	}); err != nil {
		return fmt.Errorf("create review: %w", err)
	}

	// --- canceled lifecycle: created -> offer -> accepted -> assigned ->
	// canceled, with a REFUNDED payment ---
	canceledReq, err := s.request.CreateRequest(ctx, client1.UserID, request.CreateRequestRequest{
		ServiceID: meterReplacementID, Description: "Истёк срок поверки, нужна замена",
		AddressText: "г. Самара, ул. Ленина, д. 10, кв. 5",
	})
	if err != nil {
		return fmt.Errorf("create canceled request: %w", err)
	}
	canceledOffer, err := s.request.SubmitOffer(ctx, canceledReq.ID, master1.UserID, request.SubmitOfferRequest{Price: 2000})
	if err != nil {
		return fmt.Errorf("submit offer on canceled request: %w", err)
	}
	if _, err := s.request.AcceptOffer(ctx, canceledReq.ID, canceledOffer.ID, client1.UserID); err != nil {
		return fmt.Errorf("accept offer on canceled request: %w", err)
	}
	if _, err := s.request.CancelRequest(ctx, canceledReq.ID, client1.UserID, "Нашёл мастера через знакомых"); err != nil {
		return fmt.Errorf("cancel request: %w", err)
	}

	// --- favorites: client1 favorites master1 after a good experience ---
	if err := s.request.AddFavorite(ctx, client1.UserID, master1.UserID); err != nil {
		return fmt.Errorf("add favorite: %w", err)
	}

	s.logger.Info("seed complete")
	s.logger.Info("login with", "usernames", "admin/master1/master2/client1/client2", "password", password)
	return nil
}

// seedCategory creates a category (and, recursively, one level of
// subcategories) plus its services, collecting every created service id by
// name so the lifecycle-seeding code below can look them up without
// hardcoding ids that would break if this table were reordered.
func (s *seeder) seedCategory(cat categorySeed, parentID *int64, serviceIDsByName map[string]int64, totalCategories, totalServices *int) error {
	ctx := s.ctx
	created, err := s.catalog.CreateCategory(ctx, catalog.CreateCategoryRequest{Name: cat.name, ParentCategoryID: parentID})
	if err != nil {
		return fmt.Errorf("create category %q: %w", cat.name, err)
	}
	*totalCategories++

	for _, svc := range cat.services {
		priceFrom, priceTo := svc.priceFrom, svc.priceTo
		createdSvc, err := s.catalog.CreateService(ctx, catalog.CreateServiceRequest{
			CategoryID: created.ID, Name: svc.name, Description: svc.description,
			PriceFrom: &priceFrom, PriceTo: &priceTo, Unit: svc.unit,
		})
		if err != nil {
			return fmt.Errorf("create service %q: %w", svc.name, err)
		}
		serviceIDsByName[svc.name] = createdSvc.ID
		*totalServices++
	}

	for _, sub := range cat.subcategories {
		if err := s.seedCategory(sub, &created.ID, serviceIDsByName, totalCategories, totalServices); err != nil {
			return err
		}
	}
	return nil
}
