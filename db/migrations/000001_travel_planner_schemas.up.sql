CREATE EXTENSION "uuid-ossp";

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