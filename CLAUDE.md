# CLAUDE.md — booking-service (Backend, Go)

Backend of the QueueUp booking/appointment system. Read before writing or reviewing code.

## What this service does
REST API for auth + managing services + booking slots, plus an admin area for staff.
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
  middleware/            # JWT auth + role guard
  model/                 # domain types
migrations/              # golang-migrate .sql
queries/                 # .sql for sqlc
```
Keep layers separate: handler -> service -> repository. Handlers must never run SQL directly.

## API Endpoints
```
# Auth
POST   /api/auth/register          create account -> 201, NO session cookie (user then logs in)
POST   /api/auth/login             -> 200, set httpOnly cookie (customer or admin; role comes from the user record)
POST   /api/auth/logout            clear the cookie
POST   /api/admin/register         create an ADMIN account; body carries invite_code (gated, see below)

# Services (public read)
GET    /api/services               list VISIBLE services (is_active = true)
GET    /api/services/:id/slots?date=  available slots for that day

# Bookings
POST   /api/bookings               [auth]
GET    /api/bookings/me            [auth]
DELETE /api/bookings/:id           [auth]  owner only, else 403

# Admin [auth + role = admin]
POST   /api/admin/services         create a service
PATCH  /api/admin/services/:id     edit (name, duration, price, description, is_active)
DELETE /api/admin/services/:id     delete a service
GET    /api/admin/bookings         view all bookings
```

## Auth behavior (must follow)
- **Customer register:** email already exists -> 409 "this email is already registered". Otherwise create the user (role = 'customer'), return **201 with NO cookie** — the frontend flow is register -> go to login, so don't start a session here.
- **Login:** a wrong email OR a wrong password -> **401 with the same generic message** "email or password is incorrect". NEVER reveal whether the email exists (no account enumeration). On success set the JWT in an httpOnly + Secure cookie. The same endpoint serves admins; the JWT/claims carry the role.
- **Admin register:** verify `invite_code` against the `ADMIN_INVITE_CODE` env value, compared in constant time; missing/wrong -> 403. Email exists -> 409. Otherwise create the user with role = 'admin', return 201. Prefer seeding the first admin via migration/seed; keep this endpoint gated.
- Passwords stored with bcrypt only; email is unique (case-insensitive).

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

## DB Schema (summary)
- **users:** id, email (unique), password_hash, role ('customer' | 'admin'), created_at
- **services:** id, name, duration_minutes, price, description, is_active (bool, default true), created_at
- **bookings:** per the SQL above (has the exclusion constraint)
> `GET /api/services` returns only `is_active = true`; the admin list returns all.

## Conventions (review criteria)
- **Always validate on the server**, even though the client (zod) already validates — client validation is for UX, not security
- One error response shape: `{ "error": "message" }` + correct HTTP status (400/401/403/404/409)
- NEVER log or return `password_hash` / tokens / the admin invite code in a response or logs
- Time: store as `TIMESTAMPTZ` (UTC) in the DB
- JWT in an httpOnly + Secure cookie; don't return the token in the body
- Every `[auth]` endpoint must pass the JWT middleware; `[admin]` must pass a role guard (role = 'admin', else 403)
- Another user's booking: must return 403 (users can only view/cancel their own)
- Don't commit secrets/`.env` (incl. `ADMIN_INVITE_CODE`); use env vars

## What reviewers should watch for
- booking/cancel code not wrapped in a transaction, or not mapping `23P01` to 409
- login that distinguishes "user not found" from "wrong password" (account enumeration)
- customer register that starts a session / sets a cookie (flow is register -> login)
- admin endpoints missing the role guard, or the invite code compared non-constant-time / logged
- `GET /api/services` leaking hidden (is_active = false) services to customers
- SQL at risk of injection (must be parameterized / go through sqlc)
- endpoints missing auth/ownership checks
- sensitive data leaking into a response
