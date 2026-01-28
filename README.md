# Cat Cafe API

Backend API for the Cat Cafe application, built with Go (Fiber), PostgreSQL, and Redis.

## 🚀 Tech Stack

- **Language:** Go 1.24
- **Framework:** [Fiber v2](https://gofiber.io/)
- **Database:** PostgreSQL
- **Caching:** Redis
- **Authentication:** JWT (JSON Web Tokens)
- **Containerization:** Docker & Docker Compose
- **Hot Reload:** Air

## 🛠 Prerequisites

Ensure you have the following installed:

- [Go 1.24+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) & Docker Compose
- [Make](https://www.gnu.org/software/make/) (optional, for running Makefile commands)
- [Air](https://github.com/cosmtrek/air) (optional, for local hot reload)

## ⚙️ Configuration

The application is configured using `config.yaml` and environment variables.

1.  **Environment Variables:**
    Copy the example environment file:
    ```bash
    cp env.example .env
    ```
    Update `.env` with your database credentials if running locally without Docker Compose.

2.  **Config File:**
    Copy the example config file:
    ```bash
    cp config.yaml.template config.yaml
    ```
    Update `config.yaml` to match your environment.

## 🐳 Running with Docker (Recommended)

The easiest way to run the application is using Docker Compose. This will set up the API, PostgreSQL, and Redis.

Start the services:
```bash
make docker-up
# OR
docker-compose up -d
```

Stop the services:
```bash
make docker-down
# OR
docker-compose down
```

View logs:
```bash
make docker-logs
```

The API will be available at `http://localhost:8888`.

## 💻 Running Locally

If you prefer to run the application locally (e.g., for development):

1.  **Start Dependencies:**
    You still need PostgreSQL and Redis running. You can start them via Docker:
    ```bash
    docker-compose up -d postgres redis
    ```

2.  **Run Migrations:**
    Ensure your database schema is up to date:
    ```bash
    make migrate-up
    ```

3.  **Start the Server:**
    Using Air for hot reload:
    ```bash
    make run
    ```
    Or typically with Go:
    ```bash
    go run cmd/http/main.go
    ```

## 🗄️ Database Migrations

Manage database schema changes using the following Makefile commands:

- **Create a new migration:**
    ```bash
    make migrate-create name=migration_name
    ```
- **Apply migrations (Up):**
    ```bash
    make migrate-up
    ```
- **Rollback migrations (Down):**
    ```bash
    make migrate-down
    ```
- **Force specific version:**
    ```bash
    make migrate-force v=version_number
    ```

## 🎮 Seeding Data

For testing purposes, you can seed the database:

- **Create a new seed:**
    ```bash
    make seed-create name=seed_name
    ```
- **Apply seeds:**
    ```bash
    make seed-up
    ```
- **Revert seeds:**
    ```bash
    make seed-down
    ```

## 📂 Project Structure

```
├── cmd/
│   └── http/           # Main entry point for the HTTP server
├── db/
│   ├── migrations/     # SQL migration files
│   └── seeds/          # SQL seed files
├── internal/           # Private application code
│   ├── config/         # Configuration logic
│   ├── handlers/       # HTTP handlers (Controllers)
│   ├── middleware/     # Fiber middleware
│   ├── models/         # Domain models
│   ├── repositories/   # Data access layer
│   └── usecase/        # Business logic
├── pkg/                # Public shared code
├── Dockerfile          # Docker build instructions
├── docker-compose.yaml # Docker Compose services
└── Makefile            # Command shortcuts
```
