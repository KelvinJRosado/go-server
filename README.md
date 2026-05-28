# Go Server

A basic Go HTTP API for a Chirpy-style social app. The server supports user registration and login, JWT authentication, refresh tokens, creating and reading short posts called "chirps", deleting a user's own chirps, simple admin metrics, and a webhook endpoint for upgrading users.

The app listens on `http://localhost:8080` and serves the web app under `/app/`.

## Tech Stack

- **Go 1.26**
- **net/http** for routing and HTTP handlers
- **PostgreSQL** for persistence
- **sqlc** for generating type-safe database access code
- **goose** for database migrations
- **JWT** authentication with `github.com/golang-jwt/jwt/v5`
- **Argon2id** password hashing with `github.com/alexedwards/argon2id`

## Prerequisites

Install the following before running the project:

- Go 1.26 or newer
- PostgreSQL
- `goose` for migrations
- `sqlc` if you need to regenerate database code

Example tool installs:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## Environment Variables

Create a `.env` file in the project root with the values needed by the server:

```env
DB_URL=postgres://user:password@localhost:5432/go_server?sslmode=disable
JWT_SECRET=replace-with-a-secret
PLATFORM=dev
POLKA_KEY=replace-with-polka-key
```

`DB_URL` is required for the app and migration commands. `JWT_SECRET` is required for signing and validating access tokens.

## Install and Run

1. Install Go dependencies:

   ```sh
   go mod download
   ```

2. Create a PostgreSQL database that matches your `DB_URL`.

3. Run database migrations:

   ```sh
   make migrate-up
   ```

4. Start the server:

   ```sh
   make run
   ```

5. Open the app or health endpoint:

   ```sh
   open http://localhost:8080/app/
   curl http://localhost:8080/api/healthz
   ```

## Useful Commands

```sh
make run             # Run the app
make build           # Build the go-server binary
make test            # Run tests
make ci              # Run formatting, vet, tests, build, and tidy checks
make migrate-up      # Apply database migrations
make migrate-down    # Roll back the latest migration
make sqlc            # Regenerate database code from SQL queries
make clean           # Remove build and coverage artifacts
```
