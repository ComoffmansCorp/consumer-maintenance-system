-- Enum-style CHECK constraints for string columns introduced before domains existed.
ALTER TABLE users ADD CONSTRAINT ck_users_role
    CHECK (role IN ('SUPER_ADMIN', 'TENANT_ADMIN', 'DISPATCHER', 'ELECTRICIAN'));

ALTER TABLE organizations ADD CONSTRAINT ck_organizations_type
    CHECK (type IN ('COMMERCIAL', 'GOVERNMENT', 'RESIDENTIAL'));

ALTER TABLE tasks ADD CONSTRAINT ck_tasks_type
    CHECK (type IN ('INSPECTION', 'REPLACEMENT'));

ALTER TABLE tasks ADD CONSTRAINT ck_tasks_status
    CHECK (status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'CANCELED'));

ALTER TABLE inspection_acts ADD CONSTRAINT ck_inspection_acts_type
    CHECK (inspection_type IN ('SCHEDULED', 'UNSCHEDULED'));

ALTER TABLE meters ADD CONSTRAINT ck_meters_type
    CHECK (type IN ('SINGLE_PHASE', 'THREE_PHASE_DIRECT', 'THREE_PHASE_TRANSFORMER'));

ALTER TABLE meters ADD CONSTRAINT ck_meters_seal_state
    CHECK (seal_state IS NULL OR seal_state IN ('INTACT', 'BROKEN', 'MISSING'));

ALTER TABLE photos ADD CONSTRAINT ck_photos_exactly_one_act
    CHECK (
        (inspection_act_id IS NOT NULL AND replacement_act_id IS NULL) OR
        (inspection_act_id IS NULL AND replacement_act_id IS NOT NULL)
    );

-- Bookkeeping timestamps missing from the initial schema.
ALTER TABLE users ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE organizations ADD COLUMN active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE organizations ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE organizations ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE addresses ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE addresses ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE tasks ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE tasks ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE tasks ADD COLUMN completed_at TIMESTAMPTZ NULL;
ALTER TABLE tasks ADD COLUMN canceled_at TIMESTAMPTZ NULL;
ALTER TABLE tasks ADD COLUMN cancel_reason VARCHAR(500) NULL;

ALTER TABLE inspection_acts ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE inspection_acts ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE replacement_acts ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE replacement_acts ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE meters ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Photo metadata needed for validated uploads and proper downloads.
ALTER TABLE photos ADD COLUMN original_filename VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE photos ADD COLUMN content_type VARCHAR(100) NOT NULL DEFAULT 'application/octet-stream';
ALTER TABLE photos ADD COLUMN size_bytes BIGINT NOT NULL DEFAULT 0;
ALTER TABLE photos ADD COLUMN uploaded_by BIGINT NULL REFERENCES users(id);
ALTER TABLE photos ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE photos ALTER COLUMN original_filename DROP DEFAULT;
ALTER TABLE photos ALTER COLUMN content_type DROP DEFAULT;

-- In-app notifications, delivered via the internal event broker.
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message VARCHAR(1000) NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    read_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_notifications_user_unread ON notifications (user_id, read_at) WHERE read_at IS NULL;
CREATE INDEX ix_notifications_tenant_created ON notifications (tenant_id, created_at DESC);
