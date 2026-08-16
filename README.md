# Reliable Email Delivery

This project demonstrates durable acceptance and asynchronous email delivery with Go and PostgreSQL. The current foundation provides a PostgreSQL container, versioned Goose migrations, and separate process-liveness and database-readiness probes.

## Prerequisites

- Go 1.26.3
- Docker with Colima
- Docker Compose (`docker-compose`)
- [Goose](https://github.com/pressly/goose)
- [just](https://github.com/casey/just)

## Start From A Clean Clone

Create local configuration and replace the example password with a local development password:

```sh
cp .env.example .env
openssl rand -hex 24
```

Copy the generated password into `POSTGRES_PASSWORD` in `.env`.

Start PostgreSQL and confirm that it is healthy:

```sh
just db-up
just db-status
```

Apply the migrations and inspect their status:

```sh
just migrate
just migration-status
```

Start the service:

```sh
just run
```

In another terminal, inspect the probes:

```sh
curl -i http://localhost:8080/livez
curl -i http://localhost:8080/readyz
```

Both probes return `204 No Content` while PostgreSQL is available.

## Verify Database Readiness

Stop PostgreSQL while leaving the service running:

```sh
just db-stop
curl -i http://localhost:8080/livez
curl -i http://localhost:8080/readyz
```

Liveness remains `204 No Content`; readiness changes to `503 Service Unavailable`.

Restart PostgreSQL and wait for it to become healthy:

```sh
just db-start
just db-status
curl -i http://localhost:8080/readyz
```

Readiness returns to `204 No Content` without restarting the service.

## Development Commands

```sh
just test              # Run unit tests
just check             # Run formatting, vet, and tests
just build             # Build bin/server
just db-down           # Remove the PostgreSQL container but keep its data
just db-reset          # Remove the PostgreSQL container and all local data
```

After `just db-reset`, run `just db-up` and `just migrate` to initialize a new database.
