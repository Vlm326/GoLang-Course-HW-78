# Домашнее задание №2

Распределенная система для получения информации о репозиториях GitHub.

## Что реализовано

Проект состоит из двух сервисов:

- `collector`:
  - gRPC-сервер;
  - получает данные о репозитории из GitHub API;
  - инкапсулирует бизнес-логику и работу с внешним HTTP API.
- `api-gateway`:
  - REST-сервер;
  - принимает внешние HTTP-запросы;
  - вызывает `collector` по gRPC;
  - отдает OpenAPI-спецификацию и страницу Swagger UI.

## Архитектура

Использована структура, близкая к Clean Architecture:

- `domain`:
  - доменные сущности и ошибки;
- `usecase`:
  - бизнес-логика сценариев;
- `adapter`:
  - интеграции с внешними системами;
- `handler`:
  - transport-слой (`grpc` и `http`);
- `shared`:
  - общий gRPC-контракт и codec для взаимодействия сервисов.

## Стек

- Go 1.24
- `net/http` для REST и GitHub API
- `google.golang.org/grpc` для межсервисного взаимодействия
- статический OpenAPI и Swagger UI без дополнительных framework-зависимостей

## Структура проекта

```text
cmd/
  api-gateway/
  collector/
internal/
  collector/
    adapter/
    domain/
    handler/
    usecase/
  gateway/
    handler/
    usecase/
  shared/
    grpcjson/
    repositoryrpc/
```

## Переменные окружения

### Collector

- `COLLECTOR_ADDRESS`
  - по умолчанию `:50051`

### API Gateway

- `API_GATEWAY_ADDRESS`
  - по умолчанию `:8080`
- `COLLECTOR_GRPC_ADDRESS`
  - по умолчанию `localhost:50051`

## Запуск локально

### 1. Установить зависимости

```bash
go mod tidy
```

### 2. Запустить collector

```bash
go run ./cmd/collector
```

### 3. Запустить api-gateway

```bash
go run ./cmd/api-gateway
```

## Запуск через Docker Compose

```bash
docker compose up --build
```

## REST API

### Получить информацию о репозитории

```http
GET /api/v1/repositories/{owner}/{repo}
```

Пример:

```bash
curl http://localhost:8080/api/v1/repositories/golang/go
```

Пример ответа:

```json
{
  "name": "go",
  "description": "The Go programming language",
  "stars": 123456,
  "forks": 23456,
  "created_at": "2009-11-10T23:00:00Z"
}
```

## OpenAPI и Swagger UI

- OpenAPI JSON:

```bash
curl http://localhost:8080/openapi.json
```

- Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

## Обработка ошибок

- `400 Bad Request`
  - некорректный ввод;
- `404 Not Found`
  - репозиторий не найден;
- `502 Bad Gateway`
  - collector недоступен;
- `504 Gateway Timeout`
  - превышено время ожидания upstream-сервиса;
- `500 Internal Server Error`
  - прочие ошибки.

На уровне `collector` ошибки маппятся в корректные gRPC status codes, а на уровне `api-gateway` переводятся в HTTP status codes.

## Проверка

```bash
go test ./...
```

## Техническая заметка

В текущем окружении отсутствует `protoc`, поэтому gRPC-контракт реализован вручную через `grpc.ServiceDesc` и JSON codec. Межсервисное взаимодействие при этом остается gRPC, а внешнее API остается REST.
