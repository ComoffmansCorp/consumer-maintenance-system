-- name: GetTenantByID :one
SELECT id, name, code, plan, active, created_at
FROM tenants
WHERE id = @id
LIMIT 1;

-- name: GetTenantByCode :one
SELECT id, name, code, plan, active, created_at
FROM tenants
WHERE LOWER(code) = LOWER(@code)
LIMIT 1;

-- name: ExistsTenantByName :one
SELECT EXISTS(
    SELECT 1 FROM tenants WHERE LOWER(name) = LOWER(@name)
) AS exists;

-- name: ExistsTenantByCode :one
SELECT EXISTS(
    SELECT 1 FROM tenants WHERE LOWER(code) = LOWER(@code)
) AS exists;

-- name: CreateTenant :one
INSERT INTO tenants (name, code, plan, active)
VALUES (@name, @code, @plan, @active)
RETURNING id, name, code, plan, active, created_at;

-- name: ListTenants :many
SELECT id, name, code, plan, active, created_at
FROM tenants
ORDER BY created_at DESC;

-- name: UpdateTenantPlan :one
UPDATE tenants
SET plan = @plan
WHERE id = @id
RETURNING id, name, code, plan, active, created_at;

-- name: DeactivateTenant :one
UPDATE tenants
SET active = FALSE
WHERE id = @id
RETURNING id, name, code, plan, active, created_at;
