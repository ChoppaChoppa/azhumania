# Структурная схема проекта Azhumania

## Обзор проекта
**Azhumania** - это Telegram бот, построенный на архитектуре Clean Architecture с использованием Go. Проект использует PostgreSQL для постоянного хранения данных и Redis для кэширования.

## Архитектурные слои

```
┌─────────────────────────────────────────────────────────────┐
│                        PRESENTATION LAYER                   │
├─────────────────────────────────────────────────────────────┤
│  cmd/main.go                                                │
│  └── Точка входа в приложение                              │
│      ├── Инициализация сервисов                            │
│      ├── Настройка конфигурации                            │
│      └── Запуск Telegram бота                              │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                        BOT LAYER                            │
├─────────────────────────────────────────────────────────────┤
│  internal/bot/telegram/                                     │
│  ├── server.go                                              │
│  │   ├── TelegramBot struct                                │
│  │   ├── New() - конструктор бота                          │
│  │   └── Интеграция с Telegram API                         │
│  └── listen.go                                              │
│      ├── Listen() - основной цикл обработки сообщений      │
│      └── Обработка входящих сообщений                      │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                      SERVICE LAYER                          │
├─────────────────────────────────────────────────────────────┤
│  internal/service/                                           │
│  ├── service.go                                             │
│  │   ├── IService interface                                 │
│  │   ├── service struct                                     │
│  │   ├── New() - конструктор сервиса                       │
│  │   └── Handle() - обработка сообщений                    │
│  ├── handler.go                                             │
│  │   └── Логика обработки команд                           │
│  ├── users.go                                               │
│  │   └── Бизнес-логика работы с пользователями             │
│  └── models/                                                │
│      ├── azhumania.go                                       │
│      └── users.go                                           │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                    REPOSITORY LAYER                         │
├─────────────────────────────────────────────────────────────┤
│  internal/repository/                                        │
│  ├── models/                                                │
│  │   ├── azhumania.go                                       │
│  │   └── users.go                                           │
│  ├── database/psql/                                         │
│  │   ├── irepository.go                                     │
│  │   │   ├── IDatabase interface                           │
│  │   │   ├── IUsersDatabase interface                      │
│  │   │   └── IAzhumaniaDatabase interface                  │
│  │   ├── repository.go                                      │
│  │   │   └── Реализация интерфейсов БД                     │
│  │   ├── azhumania.go                                       │
│  │   │   └── Операции с данными Azhumania                  │
│  │   └── users.go                                           │
│  │       └── Операции с пользователями                     │
│  └── cache/redis/                                           │
│      ├── icache.go                                          │
│      │   ├── ICache interface                              │
│      │   ├── IUsersCache interface                         │
│      │   └── IAzhumaniaCache interface                     │
│      ├── cache.go                                           │
│      │   └── Реализация Redis кэша                         │
│      ├── azhumania.go                                       │
│      │   └── Кэширование данных Azhumania                  │
│      └── users.go                                           │
│          └── Кэширование пользователей                     │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                      DATA LAYER                             │
├─────────────────────────────────────────────────────────────┤
│  PostgreSQL Database                                        │
│  ├── Таблица users                                          │
│  └── Таблица azhumania                                      │
│                                                             │
│  Redis Cache                                                │
│  ├── Кэш пользователей                                      │
│  └── Кэш данных azhumania                                   │
└─────────────────────────────────────────────────────────────┘
```

## Поток данных

```
1. Telegram Message
   ↓
2. cmd/main.go (инициализация)
   ↓
3. internal/bot/telegram/listen.go (получение сообщения)
   ↓
4. internal/service/service.go (обработка бизнес-логики)
   ↓
5. internal/repository/cache/redis/ (проверка кэша)
   ↓
6. internal/repository/database/psql/ (если нет в кэше)
   ↓
7. Response back through layers
```

## Основные компоненты

### 1. **cmd/main.go**
- Точка входа в приложение
- Инициализация всех зависимостей
- Конфигурация подключений к БД и Redis
- Запуск Telegram бота

### 2. **Bot Layer (internal/bot/telegram/)**
- **server.go**: Создание и настройка Telegram бота
- **listen.go**: Основной цикл обработки сообщений

### 3. **Service Layer (internal/service/)**
- **service.go**: Основной сервис с бизнес-логикой
- **handler.go**: Обработка различных команд
- **users.go**: Логика работы с пользователями
- **models/**: DTO для сервисного слоя

### 4. **Repository Layer (internal/repository/)**
- **models/**: Доменные модели
- **database/psql/**: Работа с PostgreSQL
- **cache/redis/**: Работа с Redis кэшем

## Зависимости

### Внешние библиотеки:
- `github.com/go-telegram-bot-api/telegram-bot-api/v5` - Telegram Bot API
- `github.com/jmoiron/sqlx` - Расширение для работы с SQL
- `github.com/redis/go-redis/v9` - Redis клиент
- `github.com/rs/zerolog` - Логирование
- `github.com/Masterminds/squirrel` - SQL Query Builder

### Базы данных:
- **PostgreSQL** (порт 5431) - основное хранилище данных
- **Redis** (порт 55000) - кэширование

## Принципы архитектуры

1. **Clean Architecture** - разделение на слои с четкими границами
2. **Dependency Injection** - внедрение зависимостей через интерфейсы
3. **Interface Segregation** - разделение интерфейсов по функциональности
4. **Repository Pattern** - абстракция доступа к данным
5. **Caching Strategy** - использование Redis для кэширования

## Конфигурация

```go
const (
    apiKey = "7887768155:AAHrGzYl8Qic0mygjmzFTfvsonY8NIuS9dg"
    psql_dsn = "host=localhost port=5431 user=azhumania password=Asdflkjh12 dbname=azhumania"
    redis_host = "localhost:55000"
    redis_username = "default"
    redis_password = "redispw"
    redis_db = 0
)
```
