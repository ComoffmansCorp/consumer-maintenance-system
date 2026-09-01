-- Marketplace flow (platform-level, NOT tenant-scoped): a person picks a
-- service and posts a request, any master whose specialization matches can
-- claim it from the open pool. Independent of the existing B2B dispatcher/
-- electrician flow — "наряд" stays tenant-scoped, "заявка" here does not.

ALTER TABLE users DROP CONSTRAINT ck_users_role;
ALTER TABLE users ADD CONSTRAINT ck_users_role
    CHECK (role IN ('SUPER_ADMIN', 'TENANT_ADMIN', 'DISPATCHER', 'ELECTRICIAN', 'MASTER', 'CLIENT'));

CREATE TABLE service_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ux_service_categories_name_lower ON service_categories (LOWER(name));

CREATE TABLE services (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES service_categories(id),
    name VARCHAR(150) NOT NULL,
    description VARCHAR(500),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_services_category ON services (category_id);

-- 1:1 extension of users for role=MASTER. Created lazily on first profile
-- update (not at registration) to avoid the auth domain depending on
-- marketplace — auth stays a foundational, dependency-free domain.
CREATE TABLE master_profiles (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    city VARCHAR(150),
    bio VARCHAR(1000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- What a master is allowed to claim. Enforced server-side on the claim
-- endpoint, not just filtered client-side in the "open requests" list.
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
    -- suggested address (not required -- free-typed addresses without a
    -- matching suggestion still work, just without coordinates). Not used
    -- for a visual map yet, but cheap to capture now for future geo
    -- filtering/sorting without another migration.
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    master_id BIGINT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    claimed_at TIMESTAMPTZ NULL,
    completed_at TIMESTAMPTZ NULL,
    canceled_at TIMESTAMPTZ NULL,
    cancel_reason VARCHAR(500)
);

ALTER TABLE service_requests ADD CONSTRAINT ck_service_requests_status
    CHECK (status IN ('OPEN', 'IN_PROGRESS', 'COMPLETED', 'CANCELED'));

CREATE INDEX ix_service_requests_status ON service_requests (status);
CREATE INDEX ix_service_requests_client ON service_requests (client_id);
CREATE INDEX ix_service_requests_master ON service_requests (master_id);
CREATE INDEX ix_service_requests_service ON service_requests (service_id);
