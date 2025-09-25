## Blogo - Blog API in Go (Echo, sqlc, pgx)

This is a backend REST API for a simple blog platform, built with Go and the Echo framework. It demonstrates clean application structure, dependency injection, cookie-based JWT authentication, and a robust database layer using pgx + sqlc with migrations.

### Key Features
- User authentication (signup/login) with JWT stored in HttpOnly cookies
- Blog CRUD (public read; authenticated create/update/delete)
- Auth middleware that reads and verifies JWT from cookies
- CORS configured for a frontend at `http://localhost:3000`
- PostgreSQL with `pgx` pool, `sqlc` generated queries, `golang-migrate` for migrations
- Dependency Injection via a simple container for config, db, and queries
- Tests using `testify`; mock DB via `mockery`

### 🚀 Tech Stack
- Language: Go (1.22+ recommended)
- Web framework: Echo (`github.com/labstack/echo/v4`)
- Database driver/pool: `pgx` (`github.com/jackc/pgx/v5/pgxpool`)
- Query generation: `sqlc`
- Migrations: `golang-migrate` CLI
- Auth: `github.com/golang-jwt/jwt/v5` (JWT in HttpOnly, Secure cookies)
- Validation: Echo validator in `internal/services/model`
- Testing: `testify` + `mockery` for mocks

### Project Structure
- `cmd/server/main.go`: app entrypoint, DI container wiring, Echo setup and route registration
- `internal/config`: environment-driven configuration
- `internal/container`: simple DI container for `db`, `config`, and `sqlc.Querier`
- `internal/db`:
  - `connect.go`: creates `pgxpool.Pool`
  - `migration/`: SQL migrations (used by `golang-migrate`)
  - `query/`: `.sql` files for `sqlc`
  - `sqlc/`: generated Go code (models, interfaces, queries)
- `internal/handlers/u`: user handlers (signup, login, me, update, etc.)
- `internal/handlers/blog`: blog handlers (create, update, delete, get)
- `internal/middleware`: JWT cookie middleware
- `internal/routes`: route registration (user, blog, health)
- `internal/services/model`: request models and validation
- `internal/services/response`: common success response helper
- `docker/`: Dockerfile and docker-compose for Postgres
- `Makefile`: common commands (run server, docker up, migrate, test)

### Setup
1) Clone and install deps
```bash
git clone https://github.com/dhanushd-27/blog_go.git
cd blog_go
go mod download
```

2) Start Postgres with Docker
```bash
make docker-up
# Postgres will be available on localhost:5434 (per docker/docker-compose.yaml)
```

3) Configure environment
Create `.env` at repo root. All keys are required by `internal/config/config.go`.
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5434
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=blogo
SSL_MODE=disable
DB_MAX_CONN=10
DB_MIN_CONN=1
DB_CONN_MAX_LIFETIME=30m
DB_CONN_MAX_IDLE_LIFETIME=5m
DB_HEALTH_CHECK_PERIOD=30s
CONNECT_TIMEOUT=5s
JWT_SECRET=change-me
```

4) Run migrations
```bash
make migrate-up
# Uses internal/db/migration and connects to postgres://postgres:postgres@localhost:5434/blogo?sslmode=disable
```

5) Start the server
```bash
make server-up
# or
go run cmd/server/main.go
```

6) Run tests
```bash
make test
# or
go test -v ./...
```

### Auth Model
- Login issues a JWT signed with `JWT_SECRET` and sets cookie `token` with:
  - HttpOnly: true, Secure: true, SameSite=None, Path=/, 24h expiry
- Middleware reads JWT from cookie and places claims (e.g., `user_id`) on the Echo context
- CORS allows origin `http://localhost:3000` and credentials

### API Endpoints

Health
```http
GET /health  -> 200 {"status":"ok"}
```

User
```http
POST /signup
Body: { "name": string, "email": string, "password": string }
-> 200 success

POST /login
Body: { "email": string, "password": string }
-> 200 success (sets HttpOnly cookie "token")

GET /me  (auth required)
Cookie: token=<jwt>
-> 200 { id, name, email }

GET /user/all  (auth required)
-> 200

PUT /user/:id  (auth required; updates authenticated user name)
Body: { "name": string }
-> 200 success

DELETE /user/:id  (auth required)
-> 200 success
```

Blog
```http
GET /blogs
-> 200 list of blogs

GET /blogs/:id
-> 200 blog by id

POST /blogs  (auth required)
Body: { "title": string, "content": string }
-> 200 created blog

PATCH /blogs/:id  (auth required)
Body: { "title"?: string, "content"?: string } (at least one required)
-> 200 updated blog

DELETE /blogs/:id  (auth required)
-> 200 success
```

### Development Notes
- Echo server and routes are initialized in `cmd/server/main.go` with DI from `internal/container`.
- Database pool via `pgxpool` is configured in `internal/db/connect.go` from `.env`.
- `sqlc` generates type-safe queries into `internal/db/sqlc`. Update `.sql` files in `internal/db/query/` and re-run `sqlc` per `sqlc.yaml`.
- Migrations live in `internal/db/migration`. Use the `Makefile` targets `migrate-up`/`migrate-down` (requires `migrate` CLI installed).
- JWT cookie auth is enforced on protected routes via middleware in `internal/middleware/auth.go`.

### Makefile Targets
```bash
make docker-up        # start postgres
make docker-down      # stop postgres
make migrate-up       # run migrations up
make migrate-down     # rollback migrations
make server-up        # run the API server
make server-down      # kill the API server process
make test             # run tests
```

### Docker
- `docker/docker-compose.yaml` launches Postgres 17 on host port 5434.
- `docker/Dockerfile` can be used to containerize the API (to be filled as needed).