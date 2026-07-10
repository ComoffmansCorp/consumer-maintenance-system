DROP TABLE IF EXISTS notifications;

ALTER TABLE photos DROP COLUMN IF EXISTS created_at;
ALTER TABLE photos DROP COLUMN IF EXISTS uploaded_by;
ALTER TABLE photos DROP COLUMN IF EXISTS size_bytes;
ALTER TABLE photos DROP COLUMN IF EXISTS content_type;
ALTER TABLE photos DROP COLUMN IF EXISTS original_filename;

ALTER TABLE meters DROP COLUMN IF EXISTS created_at;

ALTER TABLE replacement_acts DROP COLUMN IF EXISTS updated_at;
ALTER TABLE replacement_acts DROP COLUMN IF EXISTS created_at;

ALTER TABLE inspection_acts DROP COLUMN IF EXISTS updated_at;
ALTER TABLE inspection_acts DROP COLUMN IF EXISTS created_at;

ALTER TABLE tasks DROP COLUMN IF EXISTS cancel_reason;
ALTER TABLE tasks DROP COLUMN IF EXISTS canceled_at;
ALTER TABLE tasks DROP COLUMN IF EXISTS completed_at;
ALTER TABLE tasks DROP COLUMN IF EXISTS updated_at;
ALTER TABLE tasks DROP COLUMN IF EXISTS created_at;

ALTER TABLE addresses DROP COLUMN IF EXISTS updated_at;
ALTER TABLE addresses DROP COLUMN IF EXISTS created_at;

ALTER TABLE organizations DROP COLUMN IF EXISTS updated_at;
ALTER TABLE organizations DROP COLUMN IF EXISTS created_at;
ALTER TABLE organizations DROP COLUMN IF EXISTS active;

ALTER TABLE users DROP COLUMN IF EXISTS created_at;

ALTER TABLE photos DROP CONSTRAINT IF EXISTS ck_photos_exactly_one_act;
ALTER TABLE meters DROP CONSTRAINT IF EXISTS ck_meters_seal_state;
ALTER TABLE meters DROP CONSTRAINT IF EXISTS ck_meters_type;
ALTER TABLE inspection_acts DROP CONSTRAINT IF EXISTS ck_inspection_acts_type;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS ck_tasks_status;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS ck_tasks_type;
ALTER TABLE organizations DROP CONSTRAINT IF EXISTS ck_organizations_type;
ALTER TABLE users DROP CONSTRAINT IF EXISTS ck_users_role;
