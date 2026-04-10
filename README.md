# Merch Shop API

Высоконагруженный REST-сервис для внутреннего магазина компании.  
Позволяет сотрудникам покупать мерч за монеты и переводить монеты друг другу.

---

## Стек технологий

- **Go** 1.21+
- **PostgreSQL** 15
- **Docker** / **Docker Compose**
- **JWT** для авторизации
- **chi** — роутер
- **testify** — тестирование

---

## Функциональность

| Эндпоинт | Метод | Описание | Авторизация |
|----------|-------|----------|-------------|
| `/api/auth` | POST | Регистрация / вход | Нет |
| `/api/buy/{item}` | GET | Покупка мерча | Да |
| `/api/sendCoin` | POST | Перевод монет | Да |
| `/api/info` | GET | Информация о пользователе | Да |

---

## Запуск проекта

### 1. Клонировать репозиторий

```bash
git clone https://github.com/ваш-аккаунт/go-merch-shop.git
cd go-merch-shop
```

### 2. Настроить переменные окружения
```bash
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=merch-shop
DB_PORT=5432
DB_HOST=db
DB_SSLMODE=disable
SERVER_PORT=8080
JWT_SECRET=your-super-secret-key
```

### 3. Запустить через Docker Compose
```bash
docker-compose up --build
```
Сервер будет доступен по адресу: http://localhost:8080

## Примеры запросов
### 1. Регистрация / Вход
```bash
curl -X POST http://localhost:8080/api/auth \
  -H "Content-Type: application/json" \
  -d '{"username":"alex","password":"123"}'
```
#### Ответ
```bash
{"token":"eyJhbGciOiJIUzI1NiIs..."}
```

### 2. Покупка мерча
```bash
curl -X GET http://localhost:8080/api/buy/t-shirt \
  -H "Authorization: Bearer <токен>"
```
#### Ответ
```bash
"мерч успешно приобретен"
```

### 3. Перевод монет
```bash
curl -X POST http://localhost:8080/api/sendCoin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <токен>" \
  -d '{"toUser":"ivan","amount":100}'
```
#### Ответ
```bash
"монеты успешно отправлены"
```

### 4. Получения информации
```bash
curl -X GET http://localhost:8080/api/info \
  -H "Authorization: Bearer <токен>"
```
#### Ответ
```bash
{
    "coins": 850,
    "inventory": [
        {"type": "t-shirt", "quantity": 2}
    ],
    "coinHistory": {
        "received": [
            {"fromUser": "ivan", "amount": 50}
        ],
        "sent": [
            {"toUser": "ivan", "amount": 100}
        ]
    }
}
```
## Тестирование
### Юнит-тесты (покрытие 66% - бизнес-логика)
```bash
go test -cover ./internal/service/...
```

## Архитектура
```text
handlers → services → repository → PostgreSQL
```
- **Handlers — HTTP-слой (валидация, таймауты)**
- **Services — бизнес-логика** 
- **Repository — работа с БД (транзакции, индексы)**