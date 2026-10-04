# CLAUDE.md — booking-service (Backend, Go)

Backend of the QueueUp booking/appointment system. Read before writing or reviewing code.

## What this service does
REST API for auth + managing services + booking slots.
**The part that must be correct above all:** preventing double-booking + handling race conditions at the database level.

## Tech Stack
- Go + chi (router)
- PostgreSQL + sqlc (hand-written SQL -> type-safe generated Go)
- golang-migrate (migrations)
- JWT (httpOnly cookie) + bcrypt (password hashing)
- Docker Compose (Postgres for dev)

## Structure
```
cmd/api/main.go          # entrypoint
internal/
  handler/               # HTTP handlers (parse req, return response)
  service/               # business logic
  repository/            # sqlc generated + query wrappers
  middleware/            # JWT auth
  model/                 # domain types
migrations/              # golang-migrate .sql
queries/                 # .sql for sqlc
```
Keep layers separate: handler -> service -> repository. Handlers must never run SQL directly.

## API Endpoints
```
POST   /api/auth/register
POST   /api/auth/login            -> set httpOnly cookie
POST   /api/auth/logout
GET    /api/services
GET    /api/services/:id/slots?date=
POST   /api/bookings              [auth]
GET    /api/bookings/me           [auth]
DELETE /api/bookings/:id          [auth]
POST   /api/admin/services        [auth+admin]
```

## Critical rule: prevent double-booking (must follow)
**DB layer — Postgres EXCLUDE constraint (source of truth):**
```sql
CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE TABLE bookings (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    service_id  BIGINT NOT NULL REFERENCES services(id),
    time_range  TSTZRANGE NOT NULL,
    status      TEXT NOT NULL DEFAULT 'confirmed',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    EXCLUDE USING gist (time_range WITH &&) WHERE (status = 'confirmed')
);
```
**Go layer:**
- Before insert, check `start_time > now()` -> otherwise 400 "cannot book in the past"
- Insert inside a transaction; if the DB returns an exclusion violation (SQLSTATE `23P01`) -> 409 "this time slot is already booked"
- DON'T use a bare `SELECT to check availability then INSERT` pattern (race condition)

## Conventions (review criteria)
- **Always validate on the server**, even though the client (zod) already validates — client validation is for UX, not security
- One error response shape: `{ "error": "message" }` + correct HTTP status (400/401/403/404/409)
- NEVER log or return `password_hash` / tokens in a response
- Store passwords with bcrypt only; email is unique
- Time: store as `TIMESTAMPTZ` (UTC) in the DB
- JWT in an httpOnly + Secure cookie; don't return the token in the body
- Every `[auth]` endpoint must pass the JWT middleware; `[admin]` must check role
- Another user's booking: must return 403 (users can only view/cancel their own)
- Don't commit secrets/`.env`; use env vars

## What reviewers should watch for
- booking/cancel code not wrapped in a transaction, or not mapping `23P01` to 409
- SQL at risk of injection (must be parameterized / go through sqlc)
- endpoints missing auth/ownership checks
- sensitive data leaking into a response
