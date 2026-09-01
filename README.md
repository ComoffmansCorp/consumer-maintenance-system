# Consumer Maintenance System — маркетплейс бытовых услуг

Маркетплейс в духе Профи.ру/YouDo: клиент публикует заявку на услугу,
свободные мастера откликаются с ценой (bidding), клиент выбирает
исполнителя. Бэкенд — Go (модульный монолит, ports & adapters, 6 доменов),
фронт — Vue 3 (публичный сайт + встроенная админ-панель), мобильный клиент
мастера — Android.

## Стек

- **Backend**: Go, chi (роутер), PostgreSQL + pgx, sqlc (типобезопасный SQL),
  golang-migrate (миграции), golang-jwt (access + refresh), slog (логи),
  Prometheus (метрики).
- **Frontend**: Vue 3 + Vite + TypeScript + Pinia + Tailwind CSS + Axios.
- **Android**: Kotlin, Retrofit2 + OkHttp4 + Gson, coroutines, Navigation
  Component (роль MASTER — клиентская роль веб-only).
- **Инфраструктура**: Docker, docker-compose.

## Архитектура

Модульный монолит: каждый домен (`internal/<domain>`) — это
`handler → service → repository`, с собственными `models.go`/`dto.go` и,
где нужно, `ports.go` — интерфейсами к соседним доменам. Реализации портов
связываются только в `cmd/api/main.go` (composition root); домены не
импортируют друг друга напрямую. Кросс-доменные эффекты (например,
«оффер принят» → «создать эскроу-платёж») идут через
`internal/platform/broker` — простой in-process pub/sub.

Домены: `auth` (пользователи/роли/JWT), `catalog` (категории с
подкатегориями, услуги), `master` (профиль и специализация мастера,
рейтинг), `request` (заявки, торги/отклики, история статусов, избранное),
`review` (отзывы, денормализованный рейтинг мастера), `payment`
(эскроу-симуляция: HELD → RELEASED/REFUNDED), `chat` (переписка внутри
заявки).

Роли: `SUPER_ADMIN` (каталог, модерация), `CLIENT` (публикует заявки),
`MASTER` (откликается, выполняет). Без мультитенантности — платформенный
уровень, `tenant_id` в схеме нет.

Жизненный цикл заявки: `OPEN` (принимает отклики от мастеров с подходящей
специализацией — проверка на сервере) → клиент выбирает оффер → `ASSIGNED`
(эскроу `HELD`) → мастер выполняет → `COMPLETED` (эскроу `RELEASED`, можно
оставить отзыв) либо `CANCELED` с любого нетерминального статуса
(`REFUNDED`).

## Быстрый запуск

### Через Docker — весь стек одной командой

```bash
cp .env.example .env
# впиши JWT_SECRET, например: openssl rand -base64 32
docker compose up --build -d
```

Поднимает Postgres → применяет миграции (`migrate`) → заполняет базу
демо-данными (`seed`) → стартует API (`app`, `http://localhost:8080`) →
стартует фронтенд за nginx (`frontend`, **`http://localhost:5173`**).
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

Сид создаёт полный каталог, 64 заявки во всех статусах с офферами/отзывами/
платежами/чатом, даты разбросаны за последние ~30 дней — сервис выглядит
так, будто уже работает месяц.

### Локально без Docker

```bash
export $(grep -v '^#' .env.example | xargs)   # или свой .env
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/curs3?sslmode=disable
export JWT_SECRET=$(openssl rand -base64 32)

make migrate-up   # требует установленный golang-migrate CLI
make run
```

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

## Миграции

Одна чистая миграция `migrations/000001_init.{up,down}.sql` — 13 таблиц:
`users`, `refresh_tokens`, `service_categories`, `services`,
`master_profiles`, `master_specializations`, `service_requests`,
`request_offers`, `request_status_history`, `favorites`, `reviews`,
`payments`, `messages`.

## Генерация кода из SQL

Запросы лежат в `sql/<domain>/*.sql`, sqlc генерирует типобезопасный Go-код в
`internal/<domain>/db/` согласно `sqlc.yaml`. После правки SQL:

```bash
make generate
```

## Структура репозитория

```
cmd/api/main.go             # composition root
cmd/seed/main.go            # демо-данные (каталог, мастера, клиенты, заявки)
internal/<domain>/          # auth, catalog, master, request, review, payment, chat
internal/platform/          # auth (jwt), db, broker, httpx, middleware,
                             # observability, config, wiring (adapters)
internal/server/router.go   # маршрутизация
migrations/                 # golang-migrate
sql/<domain>/*.sql          # sqlc-запросы
frontend/                   # Vue 3 + Vite — публичный сайт + админка
```
