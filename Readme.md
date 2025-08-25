PVZ Service — Сервис приёмки заказов в пунктах выдачи.

Сервис для управления ПВЗ, приёмками и товарами. Поддерживает роли `moderator` и `employee`, мониторинг через Prometheus
и Grafana.

## 📘 API Документация (OpenAPI)

Полная спецификация API описана в файле [`swagger.yml`](swagger.yml)

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

Для запуска требуется **Docker** и **Docker Compose**.

```bash
git clone https://github.com/ne4chelovek/PVZ
cd pvz-service
docker-compose up --build