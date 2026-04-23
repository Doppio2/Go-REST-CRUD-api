# Go REST CRUD: Lab Inventory

## Overview

This repository contains a small full-stack CRUD application written in Go.

The project manages two main entities:

- `equipment`
- `experiments`

It also supports a many-to-many relationship between them, so a single experiment can reference multiple equipment items.

The application includes:

- a REST API built with the Go standard library
- SQLite persistence
- SQL migrations with `goose`
- a lightweight embedded frontend written in plain HTML, CSS, and JavaScript
- CSV export for reports

## Features

- Create, read, update, and delete equipment
- Create, read, update, and delete experiments
- Attach equipment to experiments
- Remove equipment from experiments
- Export equipment and experiment data to CSV
- Export equipment assigned to a specific experiment
- Light and dark theme toggle in the frontend

## Tech Stack

- Go
- `net/http`
- SQLite
- `github.com/glebarez/go-sqlite`
- `github.com/pressly/goose/v3`
- plain HTML, CSS, and JavaScript
- `testify` for integration tests

## Project Structure

```text
cli/
  api_service/        HTTP entrypoint
  migrate/            migration CLI
internal/
  config/             runtime configuration
  entity/             domain entities
  frontend/           embedded frontend assets and handlers
  handler/            HTTP handlers
  repo/               repository interfaces
  repo/sqlite/        SQLite implementation and migrations
```

## Requirements

- Go 1.24+

## Configuration

Environment variables:

- `APP_PORT` HTTP port, default: `8080`
- `SQLITE_PATH` SQLite database path, default: `db/go_rest_crud.db`
- `AUTO_MIGRATE` apply pending migrations on API startup, default: `true`

## Running the App

Start the API server:

```bash
go run ./cli/api_service/main.go
```

Then open:

```text
http://localhost:8080
```

## Database Migrations

Apply pending migrations:

```bash
go run ./cli/migrate up
```

Check migration status:

```bash
go run ./cli/migrate status
```

Rollback the latest migration:

```bash
go run ./cli/migrate down
```

## API Notes

Successful JSON responses use this shape:

```json
{
  "data": {}
}
```

Error responses use this shape:

```json
{
  "error": {
    "code": "validation_error",
    "message": "name must not be empty"
  }
}
```

`DELETE` endpoints return `204 No Content` on success.

## CSV Export Endpoints

- `GET /equipment?format=csv`
- `GET /experiments?format=csv`
- `GET /experiments/{id}/equipment?format=csv`

## Testing

Run all tests:

```bash
go test ./...
```

## MVP Status

The current version is suitable as an MVP and portfolio project:

- backend CRUD is implemented
- persistence and migrations are in place
- the frontend is functional
- integration tests cover the main flows

Future improvements could include:

- replacing prompt-based editing with inline forms or a modal
- streaming CSV directly instead of writing temporary files
- adding more frontend test coverage
- refining API validation and domain rules
