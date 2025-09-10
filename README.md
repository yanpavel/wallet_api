---
# Wallet API

REST API сервис для управления кошельками пользователей.  
Реализован на Go, используется PostgreSQL в качестве базы данных.  

Проект сделан как тестовое задание и демонстрирует работу с:
- `chi` для роутинга
- `database/sql` + `lib/pq` для работы с Postgres
- валидацией входных данных
- структурой приложения по слоям (`cmd/api`, `internal/store`, `internal/db` и т.д.)

---

## 🚀 Быстрый старт

### Требования
- Docker + Docker Compose
- Go 1.22+ (если запускать без контейнера)

### Запуск через Docker Compose
```bash
docker-compose up --build
```
---

После успешного старта:

API доступен на http://localhost:8080

PostgreSQL на localhost:6404 (пользователь admin, пароль adminpassword, база walletdb)

---

Пример строки подключения к БД
```bash
postgres://admin:adminpassword@localhost:6404/walletdb?sslmode=disable
```

---
📂 Структура проекта
```bash
├── cmd/api/         # Точка входа, HTTP-хендлеры
├── internal/
│   ├── db/          # Подключение и конфигурация базы данных
│   ├── store/       # Хранилища (UsersStore, WalletStore и др.)
│   └── validator/   # Валидация входных данных
├── migrations/      # SQL-миграции
└── docker-compose.yml
```
application — основная структура, которая хранит конфиг и сторы.

Stores (UsersStore, WalletStore) — слой работы с базой.

Handlers (users.go, wallet.go) — обработка HTTP-запросов.

---
🔑 Основные эндпоинты
Пользователи

POST /register — регистрация

POST /login — авторизация

Кошельки

GET /wallet/{id} — получить кошелёк

POST /wallet — операции пополнения / списания

operationType:

0 — депозит
1 — снятие средств

---
