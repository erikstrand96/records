# records

A Go-based application for managing and storing information about your vinyl record collection.

## Overview

**records** helps you keep track of your vinyl records in an organized, searchable, and easily accessible way. The backend is built with Go for performance and reliability, uses SQLite for persistent storage, and comes with Docker support for easy deployment.

## Features

- Store vinyl records with key details: title, artist, year, and notes
- SQLite database with automatic schema migrations on startup
- RESTful JSON API
- Containerized deployment with Docker and Docker Compose

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.26 or later
- [Docker](https://www.docker.com/get-started) (optional, for containerized deployment)

### Clone the Repository

```bash
git clone https://github.com/erikstrand96/records.git
cd records
```

### Run Locally

```bash
go run .
```

The server starts on `0.0.0.0:8080` by default. A SQLite database file (`records.db`) will be created automatically in the working directory.

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `HOST` | `0.0.0.0` | Host the server binds to |
| `RECORDS_PORT` | `8080` | Port the server listens on |
| `RECORDS_DB_PATH` | `records.db` | Path to the SQLite database file |

### Run with Docker Compose

```bash
docker compose up --build
```

This builds the image and starts the server on port `56473`. The SQLite database is persisted in a named Docker volume (`records_data`).

### Build and Run with Docker (standalone)

```bash
docker build -t records .
docker run -p 8080:8080 records
```

## Database

The application uses SQLite with the [modernc.org/sqlite](https://github.com/ReneKroon/ttlcache/v2) pure-Go driver. Schema migrations are managed by [golang-migrate](https://github.com/golang-migrate/migrate) and run automatically on startup.

Migration files live in `internal/database/schema/` and follow the golang-migrate naming convention (e.g., `1_create_records_table.up.sql`).

## API Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/` | Returns a welcome message |
| `GET` | `/records` | Returns all records in the collection |

All endpoints return JSON responses.

## Project Structure

```
├── main.go                  # Entry point, HTTP server setup
├── Dockerfile               # Multi-stage Docker build
├── compose.yml              # Docker Compose configuration
├── go.mod / go.sum          # Go module dependencies
├── internal/
│   ├── config/              # Application configuration
│   └── database/
│       ├── database.go      # Database connection and migrations
│       └── schema/          # Migration SQL files
└── records.db               # SQLite database (generated at runtime)
```

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

Specify your license here (e.g., MIT, Apache-2.0).
