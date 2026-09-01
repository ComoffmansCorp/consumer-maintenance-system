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
	"math/rand"
	"os"
	"strings"

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
		ctx: ctx, logger: logger, pool: pool,
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
	pool    *platformdb.Pool
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

	if err := s.seedUsersAndRequests(serviceIDsByName); err != nil {
		return err
	}

	s.logger.Info("seed complete")
	s.logger.Info("login with", "usernames", "admin/master1..master10/client1..client12", "password", password)
	return nil
}

type masterSeed struct {
	username, fullName, city, bio string
	specializations               []string // service names, must match keys in serviceIDsByName
}

// masterRoster spans every top-level category in marketplaceCatalog with
// some deliberate overlap (master1/master8 both cover meters, master2/master9
// both cover electrics) so an OPEN request in a popular category collects
// more than one competing offer, same as a real bidding marketplace.
var masterRoster = []masterSeed{
	{"master1", "Виктор Мастеров", "Самара", "10 лет опыта обслуживания приборов учёта",
		[]string{"Осмотр счётчика", "Замена счётчика", "Поверка счётчика", "Опломбировка счётчика", "Установка нового счётчика"}},
	{"master2", "Елена Электрикова", "Самара", "Электрик, 7 лет практики",
		[]string{"Замена проводки", "Установка розетки", "Установка люстры", "Устранение короткого замыкания"}},
	{"master3", "Сергей Сантехников", "Тольятти", "Мастер по сантехнике широкого профиля",
		[]string{"Устранение течи", "Прочистка канализации", "Замена труб"}},
	{"master4", "Андрей Трубин", "Самара", "Монтаж сантехнического оборудования",
		[]string{"Установка смесителя", "Установка унитаза", "Установка водонагревателя"}},
	{"master5", "Наталья Чистова", "Самара", "Клининговая служба, выезд по городу",
		[]string{"Генеральная уборка квартиры", "Уборка после ремонта", "Мытьё окон"}},
	{"master6", "Игорь Бытов", "Тольятти", "Ремонт бытовой техники на дому",
		[]string{"Ремонт стиральной машины", "Ремонт холодильника", "Установка кондиционера"}},
	{"master7", "Павел Грузов", "Самара", "Переезды и грузоперевозки любой сложности",
		[]string{"Квартирный переезд", "Грузчики", "Вывоз мусора"}},
	{"master8", "Марина Счётчикова", "Новокуйбышевск", "Специалист по водосчётчикам",
		[]string{"Опломбировка счётчика", "Установка нового счётчика", "Осмотр счётчика"}},
	{"master9", "Дмитрий Проводов", "Самара", "Электромонтажные работы",
		[]string{"Замена проводки", "Установка розетки", "Устранение короткого замыкания"}},
	{"master10", "Ольга Ремонтова", "Тольятти", "Универсальный мастер по сантехнике",
		[]string{"Устранение течи", "Установка смесителя", "Установка унитаза", "Замена труб"}},
}

var clientRoster = []string{
	"client1:Ольга Заказчикова", "client2:Дмитрий Клиентов", "client3:Анна Петрова",
	"client4:Михаил Соколов", "client5:Екатерина Волкова", "client6:Алексей Морозов",
	"client7:Юлия Новикова", "client8:Роман Кузнецов", "client9:Светлана Попова",
	"client10:Артём Лебедев", "client11:Виктория Козлова", "client12:Николай Григорьев",
}

var demoAddresses = []string{
	"г. Самара, ул. Ленина, д. 10, кв. 5", "г. Самара, ул. Гагарина, д. 22, кв. 41",
	"г. Самара, ул. Молодогвардейская, д. 15, кв. 3", "г. Самара, ул. Победы, д. 7, кв. 12",
	"г. Тольятти, ул. Автостроителей, д. 44, кв. 9", "г. Тольятти, ул. Ленинградская, д. 3, кв. 22",
	"г. Самара, ул. Стара-Загора, д. 88, кв. 15", "г. Новокуйбышевск, ул. Дзержинского, д. 5, кв. 7",
	"г. Самара, ул. Луначарского, д. 30, кв. 18", "г. Самара, ул. Мичурина, д. 12, кв. 44",
}

var clientRemarks = []string{
	"Хотелось бы сегодня или завтра.", "Соседи посоветовали обратиться на площадку.",
	"Времени в обрез, нужен оперативный выезд.", "Можно в будний день после 18:00.",
	"Живу на 5 этаже, лифта нет, учтите пожалуйста.", "Нужна консультация перед началом работ.",
	"Ранее уже пользовался подобной услугой, всё понравилось.", "",
}

var offerComments = []string{
	"Могу подъехать сегодня", "Есть опыт именно с такими случаями", "Возьму по указанной цене",
	"Свободен завтра с утра", "Работаю быстро и аккуратно", "",
}

var cancelReasons = []string{
	"Нашёл мастера через знакомых", "Передумал, отложил на потом",
	"Мастер долго не отвечал", "Решил сделать сам", "Изменились обстоятельства",
}

var positiveReviewComments = []string{
	"Отличная работа, всё аккуратно и быстро!", "Мастер приехал вовремя, всё сделал качественно.",
	"Рекомендую, всё понравилось.", "Хороший специалист, буду обращаться ещё.",
	"Сделано быстро и без лишних вопросов.",
}

var mixedReviewComments = []string{
	"В целом неплохо, но приехал с опозданием.", "Работа сделана, но пришлось напомнить пару раз.",
}

var clientChatLines = []string{
	"Здравствуйте! Когда сможете подъехать?", "Уточните, пожалуйста, стоимость.",
	"Буду ждать в указанное время.", "Спасибо за оперативность!",
}

var masterChatLines = []string{
	"Добрый день! Буду завтра к 15:00.", "Уже выезжаю.",
	"Работа выполнена, можно принимать.", "Пожалуйста, обращайтесь ещё!",
}

// seedUsersAndRequests registers the master/client roster, wires
// specializations, then generates a large, status-diverse batch of requests
// (mostly OPEN with live offers) and fully walks the completed/canceled ones
// through their real lifecycle via the domain services -- exactly like the
// original hand-written scenarios, just at volume. Every generated request
// is then backdated with direct SQL (see backdateRequest) so the platform
// reads like it has a month of real history instead of everything having
// been created in the same second.
func (s *seeder) seedUsersAndRequests(serviceIDsByName map[string]int64) error {
	ctx := s.ctx
	rng := rand.New(rand.NewSource(42)) //nolint:gosec // deterministic demo data, not security-sensitive

	serviceNameByID := make(map[int64]string, len(serviceIDsByName))
	allServiceIDs := make([]int64, 0, len(serviceIDsByName))
	for name, id := range serviceIDsByName {
		serviceNameByID[id] = name
		allServiceIDs = append(allServiceIDs, id)
	}

	masterUserIDs := make([]int64, len(masterRoster))
	serviceToMasters := map[int64][]int64{}
	for i, ms := range masterRoster {
		resp, err := s.auth.RegisterMaster(ctx, auth.RegisterRequest{Username: ms.username, Password: password, FullName: ms.fullName})
		if err != nil {
			return fmt.Errorf("register %s: %w", ms.username, err)
		}
		masterUserIDs[i] = resp.UserID

		specIDs := make([]int64, 0, len(ms.specializations))
		for _, name := range ms.specializations {
			id, ok := serviceIDsByName[name]
			if !ok {
				return fmt.Errorf("master %s: unknown service %q", ms.username, name)
			}
			specIDs = append(specIDs, id)
			serviceToMasters[id] = append(serviceToMasters[id], resp.UserID)
		}
		if _, err := s.master.UpdateProfile(ctx, resp.UserID, master.UpdateProfileRequest{
			City: ms.city, Bio: ms.bio, SpecializationIDs: specIDs,
		}); err != nil {
			return fmt.Errorf("update %s profile: %w", ms.username, err)
		}
		if err := s.backdateUser(resp.UserID, 35+rng.Intn(15)); err != nil {
			return fmt.Errorf("backdate %s registration: %w", ms.username, err)
		}
	}

	servicesWithMasters := make([]int64, 0, len(serviceToMasters))
	for id := range serviceToMasters {
		servicesWithMasters = append(servicesWithMasters, id)
	}

	clientUserIDs := make([]int64, len(clientRoster))
	for i, entry := range clientRoster {
		username, fullName, _ := strings.Cut(entry, ":")
		resp, err := s.auth.RegisterClient(ctx, auth.RegisterRequest{Username: username, Password: password, FullName: fullName})
		if err != nil {
			return fmt.Errorf("register %s: %w", username, err)
		}
		clientUserIDs[i] = resp.UserID
		if err := s.backdateUser(resp.UserID, 30+rng.Intn(15)); err != nil {
			return fmt.Errorf("backdate %s registration: %w", username, err)
		}
	}
	s.logger.Info("registered marketplace roster", "masters", len(masterUserIDs), "clients", len(clientUserIDs), "password", password)

	type reqPlan struct {
		status                 string
		daysAgoMin, daysAgoMax int
	}
	var plans []reqPlan
	for i := 0; i < 8; i++ {
		plans = append(plans, reqPlan{"OPEN_NO_OFFERS", 0, 2})
	}
	for i := 0; i < 8; i++ {
		plans = append(plans, reqPlan{"OPEN_OFFERS", 0, 3})
	}
	for i := 0; i < 8; i++ {
		plans = append(plans, reqPlan{"ASSIGNED", 1, 7})
	}
	for i := 0; i < 32; i++ {
		plans = append(plans, reqPlan{"COMPLETED", 3, 30})
	}
	for i := 0; i < 8; i++ {
		plans = append(plans, reqPlan{"CANCELED", 3, 20})
	}
	rng.Shuffle(len(plans), func(i, j int) { plans[i], plans[j] = plans[j], plans[i] })

	counts := map[string]int{}
	for _, plan := range plans {
		client := clientUserIDs[rng.Intn(len(clientUserIDs))]
		address := demoAddresses[rng.Intn(len(demoAddresses))]

		var serviceID int64
		if plan.status == "OPEN_NO_OFFERS" {
			serviceID = allServiceIDs[rng.Intn(len(allServiceIDs))]
		} else {
			serviceID = servicesWithMasters[rng.Intn(len(servicesWithMasters))]
		}
		remark := clientRemarks[rng.Intn(len(clientRemarks))]
		description := serviceNameByID[serviceID]
		if remark != "" {
			description += ". " + remark
		}

		req, err := s.request.CreateRequest(ctx, client, request.CreateRequestRequest{
			ServiceID: serviceID, Description: description, AddressText: address,
		})
		if err != nil {
			return fmt.Errorf("create request (%s): %w", plan.status, err)
		}

		compatibleMasters := serviceToMasters[serviceID]

		switch plan.status {
		case "OPEN_NO_OFFERS":
			// nothing further -- demonstrates a request nobody can bid on
			// yet (or a service with no specialized master at all).

		case "OPEN_OFFERS":
			offerCount := 1 + rng.Intn(min(3, len(compatibleMasters)))
			offerers := rng.Perm(len(compatibleMasters))[:offerCount]
			for _, idx := range offerers {
				price := float64(500 + rng.Intn(4500))
				if _, err := s.request.SubmitOffer(ctx, req.ID, compatibleMasters[idx], request.SubmitOfferRequest{
					Price: price, Comment: offerComments[rng.Intn(len(offerComments))],
				}); err != nil {
					return fmt.Errorf("submit offer: %w", err)
				}
			}

		case "ASSIGNED", "COMPLETED", "CANCELED":
			masterID := compatibleMasters[rng.Intn(len(compatibleMasters))]
			price := float64(500 + rng.Intn(4500))
			offer, err := s.request.SubmitOffer(ctx, req.ID, masterID, request.SubmitOfferRequest{
				Price: price, Comment: offerComments[rng.Intn(len(offerComments))],
			})
			if err != nil {
				return fmt.Errorf("submit offer: %w", err)
			}
			if _, err := s.request.AcceptOffer(ctx, req.ID, offer.ID, client); err != nil {
				return fmt.Errorf("accept offer: %w", err)
			}

			if plan.status == "CANCELED" {
				reason := cancelReasons[rng.Intn(len(cancelReasons))]
				if _, err := s.request.CancelRequest(ctx, req.ID, client, reason); err != nil {
					return fmt.Errorf("cancel request: %w", err)
				}
				break
			}

			msgCount := 2 + rng.Intn(3)
			for i := 0; i < msgCount; i++ {
				sender, line := client, clientChatLines[rng.Intn(len(clientChatLines))]
				if i%2 == 1 {
					sender, line = masterID, masterChatLines[rng.Intn(len(masterChatLines))]
				}
				if _, err := s.chat.SendMessage(ctx, req.ID, sender, chat.SendMessageRequest{Text: line}); err != nil {
					return fmt.Errorf("send chat message: %w", err)
				}
			}

			if plan.status == "COMPLETED" {
				if _, err := s.request.CompleteRequest(ctx, req.ID, masterID); err != nil {
					return fmt.Errorf("complete request: %w", err)
				}
				if rng.Intn(100) < 85 { // most, not all, completed requests get reviewed
					rating, comment := 5, positiveReviewComments[rng.Intn(len(positiveReviewComments))]
					if rng.Intn(100) < 15 {
						rating, comment = 3+rng.Intn(2), mixedReviewComments[rng.Intn(len(mixedReviewComments))]
					}
					if _, err := s.review.CreateReview(ctx, client, review.CreateReviewRequest{
						RequestID: req.ID, Rating: rating, Comment: comment,
					}); err != nil {
						return fmt.Errorf("create review: %w", err)
					}
				}
			}
		}

		daysAgo := plan.daysAgoMin + rng.Intn(plan.daysAgoMax-plan.daysAgoMin+1)
		if err := s.backdateRequest(req.ID, daysAgo); err != nil {
			return err
		}
		counts[plan.status]++
	}
	s.logger.Info("seeded requests", "open_no_offers", counts["OPEN_NO_OFFERS"], "open_with_offers", counts["OPEN_OFFERS"],
		"assigned", counts["ASSIGNED"], "completed", counts["COMPLETED"], "canceled", counts["CANCELED"])

	// A handful of favorites -- clients who've had a good experience saving
	// a master for next time.
	favoriteCount := 0
	for i := 0; i < 8 && i < len(clientUserIDs); i++ {
		masterID := masterUserIDs[i%len(masterUserIDs)]
		if err := s.request.AddFavorite(ctx, clientUserIDs[i], masterID); err != nil {
			return fmt.Errorf("add favorite: %w", err)
		}
		favoriteCount++
	}
	s.logger.Info("seeded favorites", "count", favoriteCount)

	return nil
}

// backdateRequest shifts a request and everything hanging off it (offers,
// status history, payment, review, chat messages) back by the same
// interval, preserving their relative ordering while spreading requests
// across roughly the last month -- so the platform doesn't look like it
// was all created in the same instant. Plain fmt-built INTERVAL is safe
// here: daysAgo always comes from this file's own rng, never user input.
func (s *seeder) backdateRequest(requestID int64, daysAgo int) error {
	if daysAgo <= 0 {
		return nil
	}
	interval := fmt.Sprintf("%d days", daysAgo)
	stmts := []string{
		fmt.Sprintf("UPDATE service_requests SET created_at = created_at - INTERVAL '%s', updated_at = updated_at - INTERVAL '%s' WHERE id = $1", interval, interval),
		fmt.Sprintf("UPDATE request_offers SET created_at = created_at - INTERVAL '%s', updated_at = updated_at - INTERVAL '%s' WHERE request_id = $1", interval, interval),
		fmt.Sprintf("UPDATE request_status_history SET created_at = created_at - INTERVAL '%s' WHERE request_id = $1", interval),
		fmt.Sprintf("UPDATE payments SET created_at = created_at - INTERVAL '%s', updated_at = updated_at - INTERVAL '%s' WHERE request_id = $1", interval, interval),
		fmt.Sprintf("UPDATE reviews SET created_at = created_at - INTERVAL '%s' WHERE request_id = $1", interval),
		fmt.Sprintf("UPDATE messages SET created_at = created_at - INTERVAL '%s' WHERE request_id = $1", interval),
	}
	for _, stmt := range stmts {
		if _, err := s.pool.Exec(s.ctx, stmt, requestID); err != nil {
			return fmt.Errorf("backdate request %d: %w", requestID, err)
		}
	}
	return nil
}

func (s *seeder) backdateUser(userID int64, daysAgo int) error {
	if daysAgo <= 0 {
		return nil
	}
	interval := fmt.Sprintf("%d days", daysAgo)
	stmt := fmt.Sprintf("UPDATE users SET created_at = created_at - INTERVAL '%s', updated_at = updated_at - INTERVAL '%s' WHERE id = $1", interval, interval)
	if _, err := s.pool.Exec(s.ctx, stmt, userID); err != nil {
		return fmt.Errorf("backdate user %d: %w", userID, err)
	}
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
