CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role          TEXT NOT NULL DEFAULT 'customer' CHECK (role IN ('customer', 'admin')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- CREATE TABLE services (
--     id               BIGSERIAL PRIMARY KEY,
--     name             TEXT NOT NULL,
--     duration_minutes INT NOT NULL CHECK (duration_minutes > 0),
--     price            NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (price >= 0),
--     created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
-- );

-- CREATE TABLE bookings (
--     id         BIGSERIAL PRIMARY KEY,
--     user_id    BIGINT NOT NULL REFERENCES users(id),
--     service_id BIGINT NOT NULL REFERENCES services(id),
--     time_range TSTZRANGE NOT NULL,
--     status     TEXT NOT NULL DEFAULT 'confirmed' CHECK (status IN ('confirmed', 'cancelled')),
--     created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
--     EXCLUDE USING gist (time_range WITH &&) WHERE (status = 'confirmed')
-- );

-- CREATE INDEX idx_bookings_user_id ON bookings (user_id);
-- CREATE INDEX idx_bookings_service_id ON bookings (service_id);
