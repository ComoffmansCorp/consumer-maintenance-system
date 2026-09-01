# Consumer Maintenance System

Production-grade система выездного обслуживания приборов учёта (счётчиков).
Диспетчер создаёт наряды на осмотр или замену счётчика, инспектор выполняет
их с мобильного/веб-интерфейса, система формирует юридически значимые PDF-акты
с фотофиксацией. Бэкенд — Go (модульный монолит, ports & adapters), админка —
Vue 3.

## Стек

- **Backend**: Go, chi (роутер), PostgreSQL + pgx, sqlc (типобезопасный SQL),
  golang-migrate (миграции), golang-jwt (access + refresh), slog (логи),
  Prometheus (метрики), signintech/gopdf (генерация PDF-актов с кириллицей).
- **Frontend**: Vue 3 + Vite + TypeScript + Pinia + Tailwind CSS + Axios.
- **Инфраструктура**: Docker, docker-compose.

## Архитектура

Модульный монолит: каждый домен (`internal/<domain>`) — это
`handler → service → repository`, с собственными `models.go`/`dto.go` и,
где нужно, `ports.go` — интерфейсами к соседним доменам. Реализации портов
связываются только в `cmd/api/main.go` (composition root); домены не
импортируют друг друга напрямую. Кросс-доменные эффекты (например, «наряд
назначен» → уведомление) идут через `internal/platform/broker` — простой
in-process pub/sub, который в будущем можно заменить на настоящую очередь без
изменений в доменном коде.

Домены: `auth` (пользователи/роли/JWT), `organization` (SaaS-тенанты,
`/api/platform/tenants`), `consumer` (клиенты-потребители — юрлица/адреса,
таблица `organizations`), `address`, `meter`, `task` (наряды и жизненный
цикл), `act` (акты осмотра и замены + PDF), `photo` (фотофиксация),
`notification` (in-app уведомления).

Роли: `SUPER_ADMIN` (платформа), `TENANT_ADMIN` (управляет сотрудниками и
справочниками своей компании), `DISPATCHER` (создаёт и назначает наряды),
`ELECTRICIAN` (инспектор, выполняет наряды и заполняет акты).

## Маркетплейс услуг (второй контур, в разработке)

Поверх B2B-системы (тенанты, диспетчер → инспектор) разрабатывается
независимый B2C-маркетплейс в духе Профи.ру: клиент оформляет заявку на
услугу, свободные мастера сами откликаются, клиент выбирает исполнителя.
Платформенный уровень, не привязан к тенантам (`tenant_id = NULL`, как у
`SUPER_ADMIN`). Домен `internal/marketplace`, роли `CLIENT`/`MASTER`,
публичные роуты `/marketplace/*` во фронтенде (`MarketplaceShell.vue`,
отдельная от B2B визуальная идентичность). B2B и маркетплейс сосуществуют
как два независимых контура одного продукта — этот раздел будет обновлён,
когда контур стабилизируется.

## Быстрый запуск

### Через Docker — весь стек одной командой

```bash
cp .env.example .env
# впиши JWT_SECRET, например: openssl rand -base64 32
make docker-up
```

Поднимает Postgres → применяет миграции (`migrate`) → заполняет базу демо-данными
(`seed`) → стартует API (`app`, `http://localhost:8080`) → стартует собранный
фронтенд за nginx (`frontend`, **`http://localhost:5173`**). Остановить:
`make docker-down`. Пересобрать с нуля (чистый volume БД): `make docker-fresh`.

`seed` идемпотентен — при повторном запуске (`docker compose up` без `-v`)
данные не дублируются, он просто пропускает шаг, если демо-тенант уже есть.

**Тестовый вход** (после `make docker-up` открой `http://localhost:5173`):

| Роль | Логин | Пароль | Код компании |
|---|---|---|---|
| Администратор | `admin` | `Demo12345` | `demo` |
| Диспетчер | `dispatcher` | `Demo12345` | `demo` |
| Инспектор (для справки, веб не используется) | `inspector1` / `inspector2` | `Demo12345` | `demo` |

В демо-данных: 4 потребителя, 6 адресов, 7 нарядов во всех статусах (ожидает
без инспектора / ожидает с инспектором / в работе / выполнен ×2 с заполненным
актом, счётчиком и фото / отменён с причиной) — готово к осмотру сразу после
логина.

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
| `FILE_UPLOAD_DIR` | директория для загруженных фото |

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

`migrations/` (golang-migrate, `up`/`down` пары). Схема: `tenants`, `users`,
`organizations` (клиенты-потребители), `addresses`, `tasks`,
`inspection_acts`, `replacement_acts`, `meters`, `photos`,
`refresh_tokens`, `notifications`.

## Генерация кода из SQL

Запросы лежат в `sql/<domain>/*.sql`, sqlc генерирует типобезопасный Go-код в
`internal/<domain>/db/` согласно `sqlc.yaml`. После правки SQL:

```bash
make generate
```

## PDF-акты

`internal/act/pdf.go` рендерит акт осмотра/замены в PDF (`signintech/gopdf`)
с встроенным кириллическим шрифтом (DejaVu Sans, `internal/act/assets/fonts`),
таблицей приборов учёта и вложенными фотографиями. Эндпоинты:
`GET /api/acts/inspection/{id}/pdf`, `GET /api/acts/replacement/{id}/pdf`.

## Структура репозитория

```
cmd/api/main.go            # composition root
internal/<domain>/         # auth, organization, consumer, address, meter,
                            # task, act, photo, notification
internal/platform/         # auth (jwt), db, broker, httpx, middleware,
                            # observability, config, wiring (adapters)
internal/server/router.go  # маршрутизация
migrations/                 # golang-migrate
sql/<domain>/*.sql          # sqlc-запросы
frontend/                   # Vue 3 + Vite админка
```
