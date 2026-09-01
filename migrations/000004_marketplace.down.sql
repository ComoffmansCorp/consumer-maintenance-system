DROP TABLE IF EXISTS service_requests;
DROP TABLE IF EXISTS master_specializations;
DROP TABLE IF EXISTS master_profiles;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS service_categories;

ALTER TABLE users DROP CONSTRAINT ck_users_role;
ALTER TABLE users ADD CONSTRAINT ck_users_role
    CHECK (role IN ('SUPER_ADMIN', 'TENANT_ADMIN', 'DISPATCHER', 'ELECTRICIAN'));
