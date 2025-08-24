PVZ Service — Сервис приёмки заказов в пунктах выдачи.

Сервис для управления ПВЗ, приёмками и товарами. Поддерживает роли `moderator` и `employee`, мониторинг через Prometheus
и Grafana.

---

## 🧰 Технологии

- Go 1.24
- PostgreSQL
- Redis (для JWT blacklisting)
- Gin (HTTP-фреймворк)
- Prometheus + Grafana (мониторинг)
- Docker & Docker Compose

---

## 🛠️ Функционал

- Создание ПВЗ (только `moderator`)
- Создание приёмки (только `employee`)
- Добавление и удаление товаров
- Закрытие приёмки
- Получение списка ПВЗ с фильтрацией
- Авторизация через JWT
- Мониторинг: метрики Prometheus (технические и бизнес-логика)
- Запуск через Docker Compose

---

## 🔐 Авторизация

Сервис использует JWT для управления доступом. Поддерживаются роли:

- `moderator`
- `employee`

### 🧪 Тестовый вход: dummyLogin

Для упрощения тестирования реализован эндпоинт:

```http
POST /dummyLogin
```

---

## 🚀 Запуск проекта

### 1. Установите зависимости

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### 2. Настройте переменные окружения

Создайте файл `.env` в корне проекта:

```env
PG_HOST=localhost
PG_PORT=5432
PG_USER=your_user
PG_PASSWORD=your_password
PG_DATABASE_NAME=pvz_db

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0

JWT_SECRET=your_jwt_secret_key
JWT_EXPIRATION=24h

SERVER_PORT=:8080