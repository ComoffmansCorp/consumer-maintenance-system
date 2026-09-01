-- Marketplace-only schema. No tenants, no B2B tables: a person is either a
-- SUPER_ADMIN (curates the catalog, moderates), a CLIENT (posts requests) or
-- a MASTER (claims/offers on requests via specialization). "Заявка" is the
-- unit of work end to end: created -> offers -> assigned -> completed, with
-- payment and chat riding along the same request id.

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ux_users_username_lower ON users (LOWER(username));

ALTER TABLE users ADD CONSTRAINT ck_users_role
    CHECK (role IN ('SUPER_ADMIN', 'CLIENT', 'MASTER'));

CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX ux_refresh_tokens_token_hash ON refresh_tokens (token_hash);
CREATE INDEX ix_refresh_tokens_user_id ON refresh_tokens (user_id);

-- One level of nesting only: a top category may have subcategories, a
-- subcategory may not have children of its own.
CREATE TABLE service_categories (
    id BIGSERIAL PRIMARY KEY,
    parent_category_id BIGINT NULL REFERENCES service_categories(id),
    name VARCHAR(150) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_service_categories_parent ON service_categories (parent_category_id);

CREATE TABLE services (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES service_categories(id),
    name VARCHAR(150) NOT NULL,
    description VARCHAR(1000),
    price_from NUMERIC(12,2),
    price_to NUMERIC(12,2),
    unit VARCHAR(30),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_services_category ON services (category_id);

-- 1:1 extension of users for role=MASTER. Created lazily on first profile
-- update (not at registration) so the auth domain never has to depend on
-- master.
CREATE TABLE master_profiles (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    city VARCHAR(150),
    bio VARCHAR(1000),
    rating_avg NUMERIC(3,2) NOT NULL DEFAULT 0,
    rating_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Specialization at the service granularity (not category) -- what a master
-- is allowed to claim/offer on. Enforced server-side, not just filtered
-- client-side.
CREATE TABLE master_specializations (
    master_user_id BIGINT NOT NULL REFERENCES master_profiles(user_id) ON DELETE CASCADE,
    service_id BIGINT NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    PRIMARY KEY (master_user_id, service_id)
);

CREATE INDEX ix_master_specializations_service ON master_specializations (service_id);

CREATE TABLE service_requests (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES users(id),
    service_id BIGINT NOT NULL REFERENCES services(id),
    description VARCHAR(2000) NOT NULL,
    address_text VARCHAR(500) NOT NULL,
    -- Populated from the Yandex Suggest response when the client picks a
    -- suggested address -- optional, a free-typed address without a
    -- matching suggestion still works, just without coordinates.
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    master_id BIGINT NULL REFERENCES users(id),
    agreed_price NUMERIC(12,2),
    cancel_reason VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE service_requests ADD CONSTRAINT ck_service_requests_status
    CHECK (status IN ('OPEN', 'ASSIGNED', 'COMPLETED', 'CANCELED'));

CREATE INDEX ix_service_requests_client ON service_requests (client_id);
CREATE INDEX ix_service_requests_master ON service_requests (master_id);
CREATE INDEX ix_service_requests_status ON service_requests (status);
CREATE INDEX ix_service_requests_service ON service_requests (service_id);

-- A master's bid on an open request. Accepting one offer rejects every
-- other PENDING offer on the same request in the same transaction.
CREATE TABLE request_offers (
    id BIGSERIAL PRIMARY KEY,
    request_id BIGINT NOT NULL REFERENCES service_requests(id) ON DELETE CASCADE,
    master_id BIGINT NOT NULL REFERENCES users(id),
    price NUMERIC(12,2) NOT NULL,
    comment VARCHAR(1000),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (request_id, master_id)
);

ALTER TABLE request_offers ADD CONSTRAINT ck_request_offers_status
    CHECK (status IN ('PENDING', 'ACCEPTED', 'REJECTED', 'WITHDRAWN'));

CREATE INDEX ix_request_offers_request ON request_offers (request_id);
CREATE INDEX ix_request_offers_master ON request_offers (master_id);

-- Full audit trail of every status transition a request goes through.
CREATE TABLE request_status_history (
    id BIGSERIAL PRIMARY KEY,
    request_id BIGINT NOT NULL REFERENCES service_requests(id) ON DELETE CASCADE,
    from_status VARCHAR(20),
    to_status VARCHAR(20) NOT NULL,
    changed_by BIGINT NOT NULL REFERENCES users(id),
    comment VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_request_status_history_request ON request_status_history (request_id);

CREATE TABLE favorites (
    client_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    master_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (client_id, master_id)
);

CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,
    request_id BIGINT NOT NULL UNIQUE REFERENCES service_requests(id),
    client_id BIGINT NOT NULL REFERENCES users(id),
    master_id BIGINT NOT NULL REFERENCES users(id),
    rating SMALLINT NOT NULL,
    comment VARCHAR(2000),
    hidden BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE reviews ADD CONSTRAINT ck_reviews_rating CHECK (rating BETWEEN 1 AND 5);

CREATE INDEX ix_reviews_master ON reviews (master_id);

-- Escrow-style ledger row per request, driven entirely by broker events
-- published from the request domain (request.assigned/completed/canceled) --
-- no direct write path from a handler.
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    request_id BIGINT NOT NULL UNIQUE REFERENCES service_requests(id),
    amount NUMERIC(12,2) NOT NULL,
    platform_fee NUMERIC(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'HELD',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE payments ADD CONSTRAINT ck_payments_status
    CHECK (status IN ('HELD', 'RELEASED', 'REFUNDED'));

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    request_id BIGINT NOT NULL REFERENCES service_requests(id) ON DELETE CASCADE,
    sender_id BIGINT NOT NULL REFERENCES users(id),
    text VARCHAR(4000) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ NULL
);

CREATE INDEX ix_messages_request_created ON messages (request_id, created_at);
