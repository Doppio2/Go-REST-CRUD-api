# Go REST CRUD: Лабораторный Инвентарь

## Обзор

Это небольшой full-stack CRUD-проект на Go.

Приложение работает с двумя основными сущностями:

- `equipment`
- `experiments`

Также поддерживается связь many-to-many между ними: один эксперимент может использовать несколько единиц оборудования.

В проект входят:

- REST API на стандартной библиотеке Go
- SQLite как база данных
- SQL-миграции через `goose`
- встроенный фронтенд на обычных HTML, CSS и JavaScript
- экспорт отчетов в CSV

## Возможности

- создание, просмотр, редактирование и удаление equipment
- создание, просмотр, редактирование и удаление experiments
- привязка equipment к experiments
- удаление связей equipment из experiments
- экспорт общего списка equipment в CSV
- экспорт общего списка experiments в CSV
- экспорт оборудования конкретного эксперимента в CSV
- переключение светлой и темной темы на фронтенде

## Технологии

- Go
- `net/http`
- SQLite
- `github.com/glebarez/go-sqlite`
- `github.com/pressly/goose/v3`
- обычные HTML, CSS и JavaScript
- `testify` для интеграционных тестов

## Структура проекта

```text
cli/
  api_service/        точка входа HTTP-сервиса
  migrate/            CLI для миграций
internal/
  config/             runtime-конфигурация
  entity/             доменные сущности
  frontend/           встроенный фронтенд и обработчики
  handler/            HTTP-хендлеры
  repo/               интерфейсы репозиториев
  repo/sqlite/        реализация SQLite и миграции
```

## Требования

- Go 1.24+

## Конфигурация

Переменные окружения:

- `APP_PORT` HTTP-порт, по умолчанию: `8080`
- `SQLITE_PATH` путь к SQLite-базе, по умолчанию: `db/go_rest_crud.db`
- `AUTO_MIGRATE` применять ли миграции при старте API, по умолчанию: `true`

## Запуск приложения

Запуск API:

```bash
go run ./cli/api_service/main.go
```

После этого открой:

```text
http://localhost:8080
```

## Миграции базы данных

Применить новые миграции:

```bash
go run ./cli/migrate up
```

Посмотреть статус миграций:

```bash
go run ./cli/migrate status
```

Откатить последнюю миграцию:

```bash
go run ./cli/migrate down
```

## Формат API

Успешные JSON-ответы:

```json
{
  "data": {}
}
```

Ошибки:

```json
{
  "error": {
    "code": "validation_error",
    "message": "name must not be empty"
  }
}
```

`DELETE`-эндпоинты при успехе возвращают `204 No Content`.

## Эндпоинты для CSV-экспорта

- `GET /equipment?format=csv`
- `GET /experiments?format=csv`
- `GET /experiments/{id}/equipment?format=csv`

## Тесты

Запуск всех тестов:

```bash
go test ./...
```

## Состояние MVP

Текущую версию уже можно считать MVP и проектом для портфолио:

- backend CRUD реализован
- база данных и миграции настроены
- фронтенд рабочий
- основные сценарии покрыты интеграционными тестами

Что можно улучшить позже:

- заменить редактирование через `prompt` на нормальные формы или модальное окно
- отдавать CSV потоком, а не через временные файлы
- добавить больше тестов на фронтенд
- усилить валидацию и доменные ограничения
