# Consumer Maintenance System — маркетплейс бытовых услуг

Маркетплейс в духе Профи.ру/YouDo: клиент публикует заявку на услугу,
свободные мастера откликаются с ценой (bidding), клиент выбирает
исполнителя, дальше — чат, эскроу-расчёт и отзыв. Бэкенд — Go (модульный
монолит, ports & adapters, 7 доменов), фронт — Vue 3 (публичный сайт,
личный кабинет, публичные профили мастеров, встроенная админ-панель),
мобильный клиент мастера — Android. Поверх монолита — полноценный
highload/observability-слой: Redis, API Gateway, Prometheus/Grafana,
централизованные логи, S3-совместимое хранилище и нагрузочное
тестирование, проверенное вживую на 400 000 RPM.

## Стек

- **Backend**: Go, chi (роутер), PostgreSQL + pgx, sqlc (типобезопасный SQL),
  golang-migrate (миграции), golang-jwt (access + refresh, ревокация через
  Redis), gorilla/websocket (real-time чат), slog (структурные JSON-логи),
  Prometheus (кастомные метрики по каждой ручке).
- **Frontend**: Vue 3 + Vite + TypeScript + Pinia + Tailwind CSS + Axios,
  нативный WebSocket с graceful fallback на поллинг.
- **Android**: Kotlin, Retrofit2 + OkHttp4 + Gson, Coil (изображения),
  coroutines, Navigation Component (роль MASTER — клиентская роль
  веб-only).
- **Инфраструктура**: Docker Compose (14 сервисов), nginx (API Gateway —
  rate limiting, keepalive, WebSocket upgrade), Redis (кэш каталога +
  ревокация сессий), MinIO (S3-совместимое хранилище фото/аватаров),
  Prometheus + Grafana (метрики), ClickHouse + Vector (централизованные
  логи), GitHub Actions (CI: build/vet/lint/миграции + фронтенд build).

## Архитектура

Модульный монолит: каждый домен (`internal/<domain>`) — это
`handler → service → repository`, с собственными `models.go`/`dto.go` и,
где нужно, `ports.go` — интерфейсами к соседним доменам. Реализации портов
связываются только в `cmd/api/main.go` (composition root); домены не
импортируют друг друга напрямую. Кросс-доменные эффекты (например,
«оффер принят» → «создать эскроу-платёж», «сообщение отправлено» →
«пуш в WebSocket-хаб») идут через `internal/platform/broker` — простой
in-process pub/sub.

Домены: `auth` (пользователи/роли/JWT), `catalog` (категории с
подкатегориями, услуги), `master` (профиль, специализация, рейтинг,
публичный каталог мастеров), `request` (заявки, торги/отклики, история
статусов, избранное), `review` (отзывы, денормализованный рейтинг
мастера), `payment` (эскроу-симуляция: HELD → RELEASED/REFUNDED), `chat`
(переписка внутри заявки — REST + WebSocket, непрочитанные диалоги).

Роли: `SUPER_ADMIN` (каталог, модерация), `CLIENT` (публикует заявки),
`MASTER` (откликается, выполняет). Без мультитенантности — платформенный
уровень, `tenant_id` в схеме нет.

Жизненный цикл заявки: `OPEN` (принимает отклики от мастеров с подходящей
специализацией — проверка на сервере) → клиент выбирает оффер → `ASSIGNED`
(эскроу `HELD`) → мастер выполняет → `COMPLETED` (эскроу `RELEASED`, можно
оставить отзыв) либо `CANCELED` с любого нетерминального статуса
(`REFUNDED`).

## Highload и эксплуатация

- **Redis** — две независимые роли: cache-aside для публичного каталога
  (`ListCategories`/`ListServices`, инвалидация на запись) и ревокация
  сессий (logout кладёт `revoked:<jti>` с TTL, `JWTAuth`-middleware
  проверяет перед каждым запросом; fail-open, если Redis недоступен).
- **API Gateway** — nginx перед `app`: rate limiting, статический
  `upstream` с keepalive-пулом (осознанный выбор — динамический resolver
  ломает keepalive и даёт массовые 502 под нагрузкой, см. комментарий в
  `observability/nginx/nginx.conf`), WebSocket upgrade для чата.
- **Нагрузочное тестирование** (`loadtest/highload.js`, k6) — ступенчатый
  разгон до **400 000 RPM** (7 полок по 30 с), независимо подтверждено
  через Prometheus: **0% ошибок** на всём диапазоне, p95-латентность на
  пике — 1.96 мс.
- **Observability** — Prometheus снимает `http_requests_total`/
  `http_request_duration_seconds` с лейблами `method`/`path`/`status` по
  каждой ручке; Grafana-дашборд `API Overview` (точная агрегатная RPM,
  request rate/error rate/latency по route, goroutines/GC) провижнится
  автоматически при старте. ClickHouse + Vector — централизованные логи,
  читаются прямо в Grafana (таблица + Explore), TTL 7 дней.
- **Изображения** — MinIO (self-hosted S3), фото услуг и аватары мастеров
  раздаются через `/media/` на обоих nginx-хопах (шлюз и фронтенд).
- **Чат в реальном времени** — WebSocket-хаб на `broker.Bus`, с честным
  graceful degradation: если сокет не подключился, фронтенд продолжает
  работать на 5-секундном поллинге, без ошибок пользователю.
- **Эксплуатация** — `make backup`/`make restore` (pg_dump/pg_restore,
  проверено round-trip), `make scale-app N=<n>` (масштабирование `app` +
  автоматический перезапуск шлюза одной командой), CI на GitHub Actions
  (build, vet, golangci-lint, прогон всех миграций на чистой БД, сборка
  фронтенда). Методология безопасных миграций схемы — `docs/zero-downtime-migrations.md`.

Известные ограничения (честно, не скрыты): open-source nginx не
перебалансирует трафик динамически между репликами `app`
(`--scale app=N` требует ручного `restart` шлюза, см. `make scale-app`);
WebSocket-хаб — in-process, не переживает горизонтальное масштабирование
`app` без Redis pub/sub в качестве шины; эскроу — симуляция, реального
платёжного провайдера нет.

## Быстрый запуск

### Через Docker — весь стек одной командой

```bash
cp .env.example .env
# впиши JWT_SECRET, например: openssl rand -base64 32
docker compose up --build -d
```

Поднимает Postgres, Redis, MinIO, ClickHouse → применяет миграции
(`migrate`) → заполняет базу демо-данными и картинками (`seed`,
`minio-init`) → стартует API (`app`) → стартует шлюз (`gateway`,
**`http://localhost:8000`** — единая точка входа для API, включая
WebSocket) → стартует фронтенд (`frontend`, **`http://localhost:5173`**).

| Что | URL |
|---|---|
| Публичный сайт | http://localhost:5173 |
| API через шлюз | http://localhost:8000 |
| Grafana (`admin`/`admin`) | http://localhost:3000 |
| Prometheus | http://localhost:9090 |
| MinIO-консоль (`minioadmin`/`minioadmin`) | http://localhost:9001 |

Остановить: `docker compose down`. Пересобрать с нуля (чистый volume БД,
**обязательно** после смены схемы миграций): `docker compose down -v &&
docker compose up --build -d`.

`seed` идемпотентен — при повторном запуске без `-v` данные не дублируются.

**Тестовый вход** (после запуска открой `http://localhost:5173`), пароль
один на всех — `Demo12345`:

| Роль | Логины |
|---|---|
| Администратор платформы | `admin` |
| Мастер | `master1` … `master10` (разные специализации по всем категориям каталога) |
| Клиент | `client1` … `client12` |

Сид создаёт полный каталог (24 услуги, 10 категорий), реалистичную историю
заявок во всех статусах с офферами/отзывами/платежами/чатом, реальные
изображения услуг и аватары мастеров в MinIO.

### Локально без Docker

```bash
export $(grep -v '^#' .env.example | xargs)   # или свой .env
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/curs3?sslmode=disable
export JWT_SECRET=$(openssl rand -base64 32)

make migrate-up   # требует установленный golang-migrate CLI
make run
```

Для чата/кэша потребуется локальный Redis (`REDIS_URL` в `.env.example`);
без него приложение всё равно стартует — session revocation и cache-aside
просто fail-open.

### Frontend (dev-сервер)

```bash
cd frontend
npm install
npm run dev
```

Vite-сервер на `http://localhost:5173` уже настроен на проксирование
`/api` на бэкенд (`vite.config.ts`).

## Переменные окружения

См. `.env.example`. Ключевые:

| Переменная | Назначение |
|---|---|
| `DATABASE_URL` | строка подключения к Postgres |
| `REDIS_URL` | строка подключения к Redis (кэш + ревокация сессий) |
| `JWT_SECRET` | секрет подписи JWT — **обязателен**, без дефолта |
| `JWT_ACCESS_TTL_MINUTES` / `JWT_REFRESH_TTL_DAYS` | время жизни токенов |
| `CORS_ALLOWED_ORIGINS` | список разрешённых origin через запятую |

Секреты никогда не коммитятся — `.env` в `.gitignore`, в репозитории только
`.env.example`.

## Команды (Makefile)

- `make run` — локальный запуск API.
- `make build` — сборка бинарника в `bin/api`.
- `make test` — `go test ./...`.
- `make lint` — `golangci-lint run ./...`.
- `make generate` — перегенерировать sqlc-код из `sql/` и `migrations/`.
- `make migrate-up` / `make migrate-down` — миграции БД.
- `make docker-up` / `make docker-down` / `make docker-fresh` — Docker Compose.
- `make scale-app N=<n>` — масштабировать `app` и перезапустить шлюз одной командой.
- `make backup` / `make restore FILE=...` — бэкап/восстановление БД (pg_dump/pg_restore).

## Миграции

`migrations/000001_init.{up,down}.sql` — 13 таблиц: `users`,
`refresh_tokens`, `service_categories`, `services`, `master_profiles`,
`master_specializations`, `service_requests`, `request_offers`,
`request_status_history`, `favorites`, `reviews`, `payments`, `messages`.

`migrations/000002_add_images.{up,down}.sql` — `services.image_url`,
`master_profiles.avatar_url` (nullable-колонки, безопасно на живой
системе — см. `docs/zero-downtime-migrations.md`).

## Генерация кода из SQL

Запросы лежат в `sql/<domain>/*.sql`, sqlc генерирует типобезопасный Go-код в
`internal/<domain>/db/` согласно `sqlc.yaml`. После правки SQL:

```bash
make generate
```

## CI

`.github/workflows/ci.yml` — два джоба на каждый push/PR: **backend**
(`go build`, `go vet`, `golangci-lint`, применение всех миграций к чистой
Postgres в контейнере) и **frontend** (`vue-tsc` + `vite build`).

## Структура репозитория

```
cmd/api/main.go              # composition root
cmd/seed/main.go             # демо-данные (каталог, мастера, клиенты, заявки, изображения)
internal/<domain>/           # auth, catalog, master, request, review, payment, chat
internal/chat/hub.go         # WebSocket-хаб (fan-out через broker.Bus)
internal/chat/ws_handler.go  # WebSocket upgrade + авторизация по ?token=
internal/platform/           # auth (jwt), cache (redis), db, broker, httpx,
                              # middleware, observability, config, wiring (adapters)
internal/server/router.go    # маршрутизация
migrations/                  # golang-migrate
sql/<domain>/*.sql           # sqlc-запросы
frontend/                    # Vue 3 + Vite — публичный сайт, кабинет, админка
observability/               # nginx (gateway), prometheus, grafana, clickhouse, vector
loadtest/highload.js         # k6-сценарий нагрузочного теста
docs/                        # zero-downtime-migrations.md и т.д.
.github/workflows/ci.yml     # CI
```
