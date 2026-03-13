CREATE TABLE IF NOT EXISTS tenants (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    code VARCHAR(100) NOT NULL,
    plan VARCHAR(30) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_tenants_name_lower ON tenants (LOWER(name));
CREATE UNIQUE INDEX IF NOT EXISTS ux_tenants_code_lower ON tenants (LOWER(code));

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(255) NOT NULL,
    tenant_id BIGINT NULL REFERENCES tenants(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_users_tenant_username_lower
    ON users (tenant_id, LOWER(username))
    WHERE tenant_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS ux_platform_users_username_lower
    ON users (LOWER(username))
    WHERE tenant_id IS NULL;

CREATE TABLE IF NOT EXISTS organizations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    tenant_id BIGINT NOT NULL REFERENCES tenants(id)
);

CREATE INDEX IF NOT EXISTS ix_organizations_tenant_name ON organizations (tenant_id, LOWER(name));

CREATE TABLE IF NOT EXISTS addresses (
    id BIGSERIAL PRIMARY KEY,
    street VARCHAR(255) NOT NULL,
    house VARCHAR(255) NOT NULL,
    building VARCHAR(255),
    apartment VARCHAR(255),
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    consumer_id BIGINT NULL REFERENCES organizations(id)
);

CREATE INDEX IF NOT EXISTS ix_addresses_tenant_street ON addresses (tenant_id, LOWER(street));

CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    address_id BIGINT NOT NULL REFERENCES addresses(id),
    status VARCHAR(255) NOT NULL,
    due_date DATE,
    assignee_id BIGINT NULL REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS ix_tasks_tenant_assignee ON tasks (tenant_id, assignee_id);
CREATE INDEX IF NOT EXISTS ix_tasks_tenant_status ON tasks (tenant_id, status);

CREATE TABLE IF NOT EXISTS inspection_acts (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT UNIQUE REFERENCES tasks(id),
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    address_id BIGINT NOT NULL REFERENCES addresses(id),
    inspection_date DATE,
    consumer_id BIGINT NULL REFERENCES organizations(id),
    inspection_type VARCHAR(255),
    notes VARCHAR(2000)
);

CREATE INDEX IF NOT EXISTS ix_inspection_acts_tenant ON inspection_acts (tenant_id);

CREATE TABLE IF NOT EXISTS replacement_acts (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT UNIQUE REFERENCES tasks(id),
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    address_id BIGINT NOT NULL REFERENCES addresses(id),
    account_number VARCHAR(255) NOT NULL,
    installation_date DATE,
    old_brand VARCHAR(255),
    old_serial_number VARCHAR(255),
    old_readings DOUBLE PRECISION,
    new_brand VARCHAR(255),
    new_serial_number VARCHAR(255),
    new_readings DOUBLE PRECISION
);

CREATE INDEX IF NOT EXISTS ix_replacement_acts_tenant ON replacement_acts (tenant_id);

CREATE TABLE IF NOT EXISTS meters (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    serial_number VARCHAR(255) NOT NULL,
    manufacture_year INTEGER,
    verification_date DATE,
    seal_state VARCHAR(255),
    transformation_ratio INTEGER,
    inspection_act_id BIGINT NULL REFERENCES inspection_acts(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS ix_meters_inspection_act ON meters (inspection_act_id);

CREATE TABLE IF NOT EXISTS photos (
    id BIGSERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    note VARCHAR(255),
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    inspection_act_id BIGINT NULL REFERENCES inspection_acts(id) ON DELETE CASCADE,
    replacement_act_id BIGINT NULL REFERENCES replacement_acts(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS ix_photos_tenant ON photos (tenant_id);
CREATE INDEX IF NOT EXISTS ix_photos_inspection_act ON photos (inspection_act_id);
CREATE INDEX IF NOT EXISTS ix_photos_replacement_act ON photos (replacement_act_id);
