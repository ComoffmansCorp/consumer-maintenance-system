-- name: CreateConsumer :one
INSERT INTO organizations (name, type, description, tenant_id, active)
VALUES (@name, @type, @description, @tenant_id, TRUE)
RETURNING id, name, type, description, tenant_id, active, created_at, updated_at;

-- name: GetConsumerByID :one
SELECT id, name, type, description, tenant_id, active, created_at, updated_at
FROM organizations
WHERE id = @id AND tenant_id = @tenant_id
LIMIT 1;

-- name: ListConsumers :many
SELECT id, name, type, description, tenant_id, active, created_at, updated_at
FROM organizations
WHERE tenant_id = @tenant_id
  AND (@search::text = '' OR name ILIKE '%' || @search::text || '%')
ORDER BY name
LIMIT @page_limit OFFSET @page_offset;

-- name: CountConsumers :one
SELECT COUNT(*)::bigint AS count
FROM organizations
WHERE tenant_id = @tenant_id
  AND (@search::text = '' OR name ILIKE '%' || @search::text || '%');

-- name: UpdateConsumer :one
UPDATE organizations
SET name = @name,
    type = @type,
    description = @description,
    updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id
RETURNING id, name, type, description, tenant_id, active, created_at, updated_at;

-- name: DeactivateConsumer :one
UPDATE organizations
SET active = FALSE, updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id
RETURNING id, name, type, description, tenant_id, active, created_at, updated_at;

-- name: ExistsConsumerByID :one
SELECT EXISTS(
    SELECT 1 FROM organizations WHERE id = @id AND tenant_id = @tenant_id AND active = TRUE
) AS exists;
