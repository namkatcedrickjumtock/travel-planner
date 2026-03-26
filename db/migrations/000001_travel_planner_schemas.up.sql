CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE trips (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID        NOT NULL,
    title       VARCHAR NOT NULL,
    destination VARCHAR NOT NULL,
start_date  TIMESTAMPTZ NOT NULL,
    end_date    TIMESTAMPTZ NOT NULL,
status      VARCHAR  NOT NULL DEFAULT 'planning',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Index speeds up per-user trip lookups and destination-based searches.
CREATE INDEX idx_trips_user_id     ON trips (user_id);
CREATE INDEX idx_trips_destination ON trips (destination);

CREATE TABLE hotels (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
name            VARCHAR    NOT NULL,
location        VARCHAR    NOT NULL,
    price_per_night NUMERIC(10, 2)  NOT NULL,
    rating          NUMERIC(3, 2)   CHECK (rating >= 0 AND rating <= 5),
    available_from  TIMESTAMPTZ,
    available_to    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
-- Index speeds up location-based hotel searches.
CREATE INDEX idx_hotels_location ON hotels (location);
CREATE TABLE flights (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    airline          VARCHAR   NOT NULL,
    origin           VARCHAR(100)   NOT NULL,
    destination      VARCHAR(100)   NOT NULL,
    departure_time   TIMESTAMPTZ    NOT NULL,
    arrival_time     TIMESTAMPTZ    NOT NULL,
    price            NUMERIC(10, 2) NOT NULL,
    seats_available  INT            NOT NULL DEFAULT 0 CHECK (seats_available >= 0),
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);
-- Compound index supports the common "flights from A to B" query pattern.
CREATE INDEX idx_flights_origin_destination ON flights (origin, destination);
CREATE TABLE activities (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
name           VARCHAR   NOT NULL,
location       VARCHAR   NOT NULL,
description    TEXT,
    price          NUMERIC(10, 2) NOT NULL,
    duration_hours NUMERIC(5, 2)  NOT NULL CHECK (duration_hours > 0),
    available_date TIMESTAMPTZ,
    created_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);
-- Index speeds up location-based activity searches.
CREATE INDEX idx_activities_location ON activities (location);
CREATE TABLE bookings (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id      UUID           NOT NULL REFERENCES trips (id) ON DELETE CASCADE,
type         VARCHAR    NOT NULL CHECK (type IN ('hotel', 'flight', 'activity')),
    reference_id UUID           NOT NULL,
status       VARCHAR    NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'cancelled')),
    total_price  NUMERIC(10, 2) NOT NULL CHECK (total_price > 0),
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);
-- Index supports the frequent "get all bookings for a trip" query.
CREATE INDEX idx_bookings_trip_id ON bookings (trip_id);

-- SEED DATA — GopherCon Europe 2026, Berlin, Germany

-- Stable UUIDs so that bookings.trip_id FK references resolve correctly.
-- trips
INSERT INTO trips (id, user_id, title, destination, start_date, end_date, status) VALUES
    ('a1b2c3d4-0001-4000-8000-000000000001', uuid_generate_v4(), 'GopherCon Europe 2026 — Berlin', 'Berlin, Germany', '2026-07-06 08:00:00+00', '2026-07-12 20:00:00+00', 'confirmed'),
    ('a1b2c3d4-0002-4000-8000-000000000002', uuid_generate_v4(), 'Munich Oktoberfest Weekend',     'Munich, Germany', '2026-09-19 10:00:00+00', '2026-09-22 18:00:00+00', 'planning');

-- hotels
INSERT INTO hotels (id, name, location, price_per_night, rating, available_from, available_to) VALUES
    (uuid_generate_v4(), 'Hotel Adlon Kempinski',     'Berlin, Germany', 450.00, 4.90, '2026-07-01 14:00:00+00', '2026-07-31 12:00:00+00'),
    (uuid_generate_v4(), 'Radisson Blu Hotel Berlin', 'Berlin, Germany', 180.00, 4.30, '2026-07-01 14:00:00+00', '2026-07-31 12:00:00+00');

-- flights
INSERT INTO flights (id, airline, origin, destination, departure_time, arrival_time, price, seats_available) VALUES
    (uuid_generate_v4(), 'Lufthansa', 'Douala, Cameroon (DLA)', 'Berlin, Germany (BER)', '2026-07-05 21:00:00+00', '2026-07-06 11:30:00+00', 920.00,  42),
    (uuid_generate_v4(), 'Eurowings', 'London, UK (LHR)',        'Berlin, Germany (BER)', '2026-07-06 06:00:00+00', '2026-07-06 08:45:00+00', 195.00, 120);

-- activities
INSERT INTO activities (id, name, location, description, price, duration_hours, available_date) VALUES
    (uuid_generate_v4(), 'GopherCon Europe 2026 Conference Pass', 'Berlin, Germany', 'Full conference access including workshops, keynotes and social events for Go developers.',                  299.00, 8.00, '2026-07-07 09:00:00+00'),
    (uuid_generate_v4(), 'Berlin Wall & History Walking Tour',    'Berlin, Germany', 'Guided walk through Cold War history: the Brandenburg Gate, Checkpoint Charlie and the East Side Gallery.',  25.00, 3.50, '2026-07-08 10:00:00+00');

-- bookings
INSERT INTO bookings (id, trip_id, type, reference_id, status, total_price) VALUES
    (uuid_generate_v4(), 'a1b2c3d4-0001-4000-8000-000000000001', 'flight',   uuid_generate_v4(), 'confirmed', 820.00),
    (uuid_generate_v4(), 'a1b2c3d4-0001-4000-8000-000000000001', 'hotel',    uuid_generate_v4(), 'confirmed', 2250.00),
    (uuid_generate_v4(), 'a1b2c3d4-0001-4000-8000-000000000001', 'activity', uuid_generate_v4(), 'confirmed', 299.00),
    (uuid_generate_v4(), 'a1b2c3d4-0002-4000-8000-000000000002', 'flight',   uuid_generate_v4(), 'pending',   145.00),
    (uuid_generate_v4(), 'a1b2c3d4-0002-4000-8000-000000000002', 'activity', uuid_generate_v4(), 'pending',    65.00); -- was 0003, fixed to 0002