set dotenv-load := true

# List available commands.
default:
    @just --list

# Start PostgreSQL.
db-up:
    docker-compose up -d postgres

# Show PostgreSQL container and health status.
db-status:
    docker-compose ps

# Stop PostgreSQL without removing it.
db-stop:
    docker-compose stop postgres

# Restart a stopped PostgreSQL container.
db-start:
    docker-compose start postgres

# Remove the PostgreSQL container while retaining its data.
db-down:
    docker-compose down

# Remove the PostgreSQL container and all local database data.
db-reset:
    docker-compose down --volumes

# Apply all pending migrations.
migrate:
    PGHOST="$DB_HOST" PGPORT="$DB_PORT" PGUSER="$POSTGRES_USER" PGPASSWORD="$POSTGRES_PASSWORD" PGDATABASE="$POSTGRES_DB" GOOSE_DRIVER=postgres GOOSE_DBSTRING="sslmode=disable" GOOSE_MIGRATION_DIR=migrations goose up

# Show migration status.
migration-status:
    PGHOST="$DB_HOST" PGPORT="$DB_PORT" PGUSER="$POSTGRES_USER" PGPASSWORD="$POSTGRES_PASSWORD" PGDATABASE="$POSTGRES_DB" GOOSE_DRIVER=postgres GOOSE_DBSTRING="sslmode=disable" GOOSE_MIGRATION_DIR=migrations goose status

# Run the Go server with local environment variables loaded.
run:
    go run ./cmd/server

# Build the application binary.
build:
    mkdir -p bin
    go build -o bin/server ./cmd/server

# Run tests.
test:
    go test -v ./...

# Run formatting, static analysis, and tests.
check:
    test -z "$(gofmt -l .)"
    go vet ./...
    go test -race ./...
