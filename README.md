# Subscription Service

REST сервис для агрегации данных об онлайн подписках пользователей.

## Технологии

- Go 1.25
- PostgreSQL 15
- Chi Router
- Docker Compose
- Swagger
- golang-migrate
- slog

## Функциональность

- Создание подписки
- Получение подписки по ID
- Обновление подписки
- Удаление подписки
- Получение списка подписок с фильтрацией
- Подсчёт суммарной стоимости подписок за период

## Запуск

### 1. Клонировать репозиторий

```bash
git clone https://github.com/GulzhanKarakul/subscription-service
cd subscription-service
```

### 2. Создать .env файл

```bash
cp .env.example .env
```

### 3. Запустить сервис

```bash
docker compose up --build
```

### 4. Применить миграции

```bash
migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" up
```

## Swagger

После запуска доступен по адресу:

http://localhost:8080/swagger/index.html