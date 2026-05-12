# TimeSnap

[![Go Version](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

REST API для управления дедлайнами на Go + Fiber + GORM (SQLite).

## Возможности

- CRUD для дедлайнов:
  - `GET /deadlines/`
  - `GET /deadlines/:id`
  - `POST /deadlines/`
  - `PATCH /deadlines/:id`
  - `DELETE /deadlines/:id`
- Валидация входных данных.
- Swagger-документация по адресу `GET /swagger/index.html`.
- Фоновый воркер, который закрывает просроченные `active` дедлайны.

## Стек

- Go `1.26+`
- Fiber v2
- GORM
- SQLite
- Testify

## Структура проекта

- [`cmd/main.go`](cmd/main.go) — точка входа приложения.
- [`internal/config/config.go`](internal/config/config.go) — загрузка конфигурации.
- [`internal/database/db.go`](internal/database/db.go) — инициализация БД.
- [`internal/modules/deadline`](internal/modules/deadline) — модуль дедлайнов (handler/service/repository).
- [`internal/pkg/worker/worker.go`](internal/pkg/worker/worker.go) — фоновый воркер.

## Конфигурация

Конфиг читается из:

- [`config.dev.yml`](config.dev.yml) — режим `dev` (по умолчанию)
- [`config.prod.yml`](config.prod.yml) — режим `prod`

Переключение режима через переменную окружения:

```bash
APP_MODE=prod
```

Основные параметры:

- `server.port`
- `database.path`
- `database.migrate`

## Запуск локально

```bash
go mod download
task dev
```

После запуска:

- API: `http://localhost:8000`
- Swagger: `http://localhost:8000/swagger/index.html`

## Docker

```bash
docker compose up --build
```

## Пример создания дедлайна

```bash
curl -X POST http://localhost:8000/deadlines/ \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Закрыть задачу",
    "status": "active",
    "priority": "high",
    "due_date": "2026-12-31T18:00:00Z"
  }'
```

## Тесты

Запуск всех тестов:

```bash
task test
```
