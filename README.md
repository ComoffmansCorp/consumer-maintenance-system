# ⚡ Meter Management System
Современная система для управления задачами, актами осмотра/замены и учётом электросчётчиков.  
Проект построен на базе Java Spring Boot и предоставляет полноценный REST API для работы с объектами учета.

---

## 🚀 Возможности системы

### 👤 Пользователи
- Авторизация по JWT
- Роли (admin / worker / user)
- Хранение профилей и доступов

### 🏢 Организации и потребители
- Учёт организаций
- Адреса объектов обслуживания
- Привязка к потребителям

### 📦 Задачи
- Создание задач (осмотр или замена счётчика)
- Назначение исполнителей
- Контроль статусов
- Привязка к адресам и объектам

### 📝 Акты
#### 🔍 Inspection Acts
- Данные осмотра
- Фотофиксация
- Дата, тип, комментарии

#### 🔧 Replacement Acts
- Замена счётчика
- Старые/новые показания
- Даты установки, номера, марки

### 🔌 Учёт счётчиков
- Производитель
- Тип
- Серийный номер
- Дата поверки и установки
- Пломбы и коэффициенты

---

## 🏗️ Технологический стек

Backend:
- Java 17
- Spring Boot (Web, Security, Data JPA)
- Lombok
- JWT
- Validation API

Database:
- PostgreSQL

DevOps:
- Docker
- Maven
- Git

---

## 📂 Структура проекта

meter-management-system/  
├── src/  
│   ├── main/  
│   │   ├── java/...  
│   │   └── resources/  
│   └── test/  
├── db/  
│   ├── schema.sql  
│   └── data.sql  
├── docker-compose.yml  
├── .gitignore  
├── README.md  
└── pom.xml

---

## 🗄️ Структура базы данных

### ER-диаграмма
(сюда можно вставить картинку из docs/db_schema.png)

---

### Таблицы

#### users
- id
- username
- password
- full_name
- role

#### addresses
- id
- street
- house
- building
- apartment
- consumer_id

#### tasks
- id
- type
- status
- due_date
- address_id
- assignee_id

#### inspection_acts
- id
- task_id
- address_id
- inspection_date
- inspection_type
- consumer_id
- notes

#### replacement_acts
- id
- task_id
- address_id
- installation_date
- old_brand
- old_serial_number
- old_readings
- new_brand
- new_serial_number
- new_readings

#### meters
- id
- type
- serial_number
- manufacturer_year
- verification_date
- seal_state
- transformation_ratio
- consumer_id

#### photos
- id
- filename
- note
- inspection_act_id
- replacement_act_id

---

## 🐳 Запуск через Docker

### PostgreSQL
Создай файл docker-compose.yml:

version: "3.9"  
services:  
database:  
image: postgres:16  
container_name: meter-db  
environment:  
POSTGRES_DB: meter  
POSTGRES_USER: postgres  
POSTGRES_PASSWORD: postgres  
ports:  
- "5432:5432"  
volumes:  
- ./db:/docker-entrypoint-initdb.d

Запуск:

docker-compose up -d

---

## 🔧 Запуск backend

mvn spring-boot:run

или

mvn clean package  
java -jar target/meter-management-system.jar

---

## 🔑 Конфигурация

Используй переменные окружения:

DB_URL  
DB_USER  
DB_PASSWORD  
JWT_SECRET

Пример в application.properties:

spring.datasource.url=${DB_URL}  
spring.datasource.username=${DB_USER}  
spring.datasource.password=${DB_PASSWORD}  
jwt.secret=${JWT_SECRET}

---

## 🧪 Тестирование API

Можно тестировать через:
- Postman
- Insomnia
- Swagger (если добавить)

---

## 📘 Планы на развитие

- Swagger/OpenAPI
- Dockerfile для backend
- Flyway миграции
- Интеграционные тесты
- Оптимизация запросов
- Кэширование
- Frontend на React

---

## 👨‍💻 Автор
ComoffmansCorp

---

## ⭐ Поставь звёздочку, если проект понравился!  
# ⚡ Meter Management System
Современная система для управления задачами, актами осмотра/замены и учётом электросчётчиков.  
Проект построен на базе Java Spring Boot и предоставляет полноценный REST API для работы с объектами учета.

---

## 🚀 Возможности системы

### 👤 Пользователи
- Авторизация по JWT
- Роли (admin / worker / user)
- Хранение профилей и доступов

### 🏢 Организации и потребители
- Учёт организаций
- Адреса объектов обслуживания
- Привязка к потребителям

### 📦 Задачи
- Создание задач (осмотр или замена счётчика)
- Назначение исполнителей
- Контроль статусов
- Привязка к адресам и объектам

### 📝 Акты
#### 🔍 Inspection Acts
- Данные осмотра
- Фотофиксация
- Дата, тип, комментарии

#### 🔧 Replacement Acts
- Замена счётчика
- Старые/новые показания
- Даты установки, номера, марки

### 🔌 Учёт счётчиков
- Производитель
- Тип
- Серийный номер
- Дата поверки и установки
- Пломбы и коэффициенты

---

## 🏗️ Технологический стек

Backend:
- Java 17
- Spring Boot (Web, Security, Data JPA)
- Lombok
- JWT
- Validation API

Database:
- PostgreSQL

DevOps:
- Docker
- Maven
- Git

---

## 📂 Структура проекта

meter-management-system/  
├── src/  
│   ├── main/  
│   │   ├── java/...  
│   │   └── resources/  
│   └── test/  
├── db/  
│   ├── schema.sql  
│   └── data.sql  
├── docker-compose.yml  
├── .gitignore  
├── README.md  
└── pom.xml

---

## 🗄️ Структура базы данных

---

### Таблицы

#### users
- id
- username
- password
- full_name
- role

#### addresses
- id
- street
- house
- building
- apartment
- consumer_id

#### tasks
- id
- type
- status
- due_date
- address_id
- assignee_id

#### inspection_acts
- id
- task_id
- address_id
- inspection_date
- inspection_type
- consumer_id
- notes

#### replacement_acts
- id
- task_id
- address_id
- installation_date
- old_brand
- old_serial_number
- old_readings
- new_brand
- new_serial_number
- new_readings

#### meters
- id
- type
- serial_number
- manufacturer_year
- verification_date
- seal_state
- transformation_ratio
- consumer_id

#### photos
- id
- filename
- note
- inspection_act_id
- replacement_act_id

---

## 🐳 Запуск через Docker

### PostgreSQL
Создай файл docker-compose.yml:

version: "3.9"  
services:  
database:  
image: postgres:16  
container_name: meter-db  
environment:  
POSTGRES_DB: meter  
POSTGRES_USER: postgres  
POSTGRES_PASSWORD: postgres  
ports:  
- "5432:5432"  
volumes:  
- ./db:/docker-entrypoint-initdb.d

Запуск:

docker-compose up -d

---

## 🔧 Запуск backend

mvn spring-boot:run

или

mvn clean package  
java -jar target/meter-management-system.jar

---

## 🔑 Конфигурация

Используй переменные окружения:

DB_URL  
DB_USER  
DB_PASSWORD  
JWT_SECRET

Пример в application.properties:

spring.datasource.url=${DB_URL}  
spring.datasource.username=${DB_USER}  
spring.datasource.password=${DB_PASSWORD}  
jwt.secret=${JWT_SECRET}

---

## 🧪 Тестирование API

Можно тестировать через:
- Postman
- Insomnia
- Swagger (если добавить)

---

## 📘 Планы на развитие

- Swagger/OpenAPI
- Dockerfile для backend
- Flyway миграции
- Интеграционные тесты
- Оптимизация запросов
- Кэширование
- Frontend на React

---

## 👨‍💻 Автор
ComoffmansCorp

---

## ⭐ Поставь звёздочку, если проект понравился!  

