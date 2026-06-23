# Consumer Maintenance System

Java Spring Boot backend для управления задачами, актами осмотра и замены, учётом организаций, адресов и счётчиков.

## Что есть в проекте

- `Dockerfile` — сборка образа приложения на Java 17.
- `docker-compose.yml` — запуск PostgreSQL и приложения вместе.
- `Makefile` — простые цели для сборки, запуска, тестов и Docker.
- `src/main/resources/db/migration/V1__init.sql` — SQL-миграции для создания таблиц.
- `src/main/resources/static/test-client.html` — простой тестовый фронтенд для проверки API.
- `uploads/` — директория для загружаемых файлов.

## Быстрый запуск

### Через Docker

```bash
make docker-up
```

После запуска приложение доступно по адресу:

- `http://localhost:8080`
- `http://localhost:8080/test-client.html` — тестовая страница для API.

Остановить контейнеры:

```bash
make docker-down
```

### Локально без Docker

```bash
make run
```

Команда `run` сначала выполняет `./mvnw clean package`, а затем запускает приложение с переменными окружения:

- `DB_URL=jdbc:postgresql://localhost:5432/curs3`
- `JWT_SECRET=5enjbH5/KpbWF7r4fZ6/ChiRbEjXfnOBr85Xxq9mBRo=`

Если вы хотите запустить приложение локально, убедитесь, что PostgreSQL доступен на `localhost:5432`.

## Makefile targets

- `make build` — сборка проекта Maven.
- `make run` — локальный запуск приложения.
- `make docker-build` — сборка Docker-образа `curs3projectback`.
- `make docker-up` — запуск приложения и PostgreSQL через `docker compose`.
- `make docker-down` — остановка контейнеров.
- `make clean` — очистка Maven-артефактов.
- `make test` — запуск тестов.

## Docker Compose

`docker-compose.yml` настраивает два сервиса:

- `db` — PostgreSQL 16 с базой `curs3` и пользователем `postgres`.
- `app` — приложение, которое собирается из `Dockerfile` и подключается к базе.

Переменные окружения для сервиса `app`:

- `DB_URL=jdbc:postgresql://db:5432/curs3`
- `DB_USERNAME=postgres`
- `DB_PASSWORD=postgres`
- `JWT_SECRET=5enjbH5/KpbWF7r4fZ6/ChiRbEjXfnOBr85Xxq9mBRo=`
- `FILE_UPLOAD_DIR=uploads`

## Миграции базы данных

Главная миграция хранится в:

- `src/main/resources/db/migration/V1__init.sql`

Это создаёт все таблицы и индексы для проекта, включая `tenants`, `users`, `organizations`, `addresses`, `tasks`, `inspection_acts`, `replacement_acts`, `meters` и `photos`.

## Конфигурация приложения

Файл конфигурации:

- `src/main/resources/application.properties`

Основные параметры:

- `spring.datasource.url=${DB_URL:jdbc:postgresql://localhost:5432/curs3}`
- `spring.datasource.username=${DB_USERNAME:postgres}`
- `spring.datasource.password=${DB_PASSWORD:postgres}`
- `file.upload-dir=${FILE_UPLOAD_DIR:uploads}`
- `jwt.secret` и `jwt.expiration`

> В текущей конфигурации `spring.jpa.hibernate.ddl-auto=create-drop`, поэтому схема создаётся и удаляется при запуске.

## Тестовый фронтенд

Откройте `http://localhost:8080/test-client.html` после запуска приложения. Это простой интерфейс для проверки:

- авторизации и регистрации
- создания задач
- работы с актами осмотра и замены
- админских эндпоинтов

## Полезные пути

- `src/main/java` — Java-код сервиса.
- `src/main/resources/static/test-client.html` — статический тестовый клиент.
- `src/main/resources/db/migration` — миграции базы.
- `docker-compose.yml` — Docker Compose конфигурация.
- `Dockerfile` — Docker-сборка приложения.
- `Makefile` — удобные команды.

## Запуск тестов

```bash
make test
```

## Контакты

Если нужно добавить документацию по API или фронтенду — просто дополни `README.md` новыми разделами.
