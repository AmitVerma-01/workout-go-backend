# Basic REST APIs for workout app ✅

A small, opinionated Go REST API for tracking workouts and workout entries. Built with Go (>=1.24), PostgreSQL, Chi router, and Goose for migrations — intended as a learning project and a minimal starting point for CRUD + token-based auth.

---

## Table of contents
- Features
- Quickstart
- Configuration
- Database & Migrations
- Run the app
- API Reference
- Testing
- Project structure
- Contributing
- Notes / Troubleshooting
- License

---

## Features
- CRUD for Users and Workouts (with entries)
- Token-based authentication (Bearer token, generated on login)
- PostgreSQL persistence with Goose migrations (embedded)
- Small, testable store layer with unit tests for the workout store

---

## Quickstart

1. Install prerequisites:
   - Go >= 1.24
   - Docker & docker-compose (recommended)
   - make (optional)

2. Start a local Postgres DB:
   - Option A — preferred: modify `docker-compose.yml` to set `POSTGRES_DB: workout_db` and run:
     ```sh
     docker compose up -d
     ```
   - Option B — use the provided `postgres_db` container and create the DB after it starts:
     ```sh
     docker compose up -d
     docker exec -it postgres_db psql -U postgres -c "CREATE DATABASE workout_db;"
     ```

3. Run the app:
   - From project root:
     ```sh
     make run
     ```
     or:
     ```sh
     go run main.go -port 8080
     ```

4. Visit health check:
   ```sh
   curl http://localhost:8080/health
   ```

---

## Configuration

The DB connection is currently hard-coded in `internals/store/database.go`:

```
host=localhost user=postgres password=postgres dbname=workout_db port=5432 sslmode=disable
```

You can either:
- Change `POSTGRES_DB` in `docker-compose.yml` to `workout_db`, or
- Adjust the connection string in `Open()` to match your environment.

(Enhancement idea: move DB config to environment variables / `.env`.)

---

## Database & Migrations

Migrations live in the `migrations/` directory and are applied automatically on app start via Goose and an embedded FS in `migrations`.

Important migrations:
- `00001_user.sql` — users table
- `00002_workout.sql` — workouts table
- `00003_workout_entries.sql` — workout entries table
- `00004_token.sql` — tokens
- `00005_user_id_alter.sql` — adds `user_id` to `workouts`

If you need to run migrations manually, you can install and use goose (or rely on app startup which applies migrations).

---

## Run the app

- Development:
  - `go run main.go` or `make run`
- Options:
  - `go run main.go -port 8080` to change the listen port

The application prints logs to stdout and serves HTTP endpoints defined in `internals/routes`.

---

## API Reference 📡

Base: `http://localhost:8080`

Authentication:
- Login returns a token object. Use the plaintext token in the Authorization header:
  ```http
  Authorization: Bearer <token_plaintext>
  ```

Endpoints:

- Health
  - `GET /health` — Response: `All good!`

- Users
  - `POST /user` — Body: `{ "name", "email", "password", "bio" }` — Creates a user
  - `GET /user/{id}` — Get user by id
  - `PATCH /user/{id}` — Update user (name/email/bio)
  - `DELETE /user/{id}` — Delete user

- Authentication
  - `POST /login` — Body: `{ "email", "password" }` — Response: `{ "token": { "plaintext", "expiry", ... } }`

- Workouts (require auth)
  - `POST /workout` — Create workout (see example below)
  - `GET /workout/{id}` — Get workout by id
  - `PATCH /workout/{id}` — Update workout
  - `DELETE /workout/{id}` — Delete workout
  - `GET /workouts` — List workouts

Create workout example:
```json
{
  "title": "Morning run",
  "description": "Park loop",
  "duration_minutes": 30,
  "calories_burned": 300,
  "entries": [
    {"exercise_name":"Running","sets":1,"reps":300,"order_index":1}
  ]
}
```

Example login & create workout:
1. Create user:
```sh
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{"name":"alice","email":"alice@example.com","password":"secret123","bio":"loves running"}'
```
2. Login:
```sh
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'
```
Copy returned `token.plaintext`
3. Create workout:
```sh
curl -X POST http://localhost:8080/workout \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_here>" \
  -d '{ ... }'
```

---

## Testing 🧪

Unit tests exist (e.g., `internals/store/workout_store_test.go`). They expect a test Postgres instance on port `5433` with DB name `postgres_test` by default.

Run tests:
1. Start a test Postgres:
```sh
docker run --name postgres_test -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres_test -p 5433:5432 -d postgres:13.1-alpine
```
2. Run:
```sh
go test ./... -v
```
Note: Tests will apply migrations and truncate tables as part of setup.

---

## Project structure 📁

- `main.go` — app entry
- `migrations/` — SQL migrations (goose)
- `internals/`
  - `api/` — HTTP handlers
  - `app/` — application bootstrap
  - `middleware/` — auth & user middleware
  - `routes/` — route wiring
  - `store/` — DB access layer (users, workouts, tokens)
  - `tokens/` — token generation & model
- `utils/` — helpers (JSON, ID read, regex)
- `docker-compose.yml` — Postgres service (dev)
- `makefile` — convenience run target

---

## Contributing 🤝

- Open issues for bugs or feature requests.
- Follow existing patterns for handlers and stores.
- Add tests for new functionality and migrations for schema changes.

---

## Notes / Troubleshooting ⚠️

> - If migrations fail on startup, ensure DB exists and the connection string matches your Postgres instance.
> - If using the provided `docker-compose`, make sure `POSTGRES_DB` matches `dbname` in `internals/store/database.go` or create the required DB manually.
> - Consider moving DB config to env vars for easier local/CI setup.

---

## License 📜

No license file included in the repository. If you want to open-source this, add a `LICENSE` (e.g., MIT) or specify your preferred license.

---

If you'd like, I can also:
- Add environment variable support for DB config, or
- Commit this `readme.md` and open a PR with the change.

Let me know which option you prefer.